pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./utils/AddressArrayUtils.sol";
import "./DelayedWithdrawal.sol";


interface OperatorContract {
    function entryVerificationGasEstimate() external view returns(uint256);
    function dkgGasEstimate() external view returns(uint256);
    function groupProfitFee() external view returns(uint256);
    function sign(
        uint256 requestId,
        uint256 seed,
        uint256 previousEntry
    ) payable external;
    function numberOfGroups() external view returns(uint256);
    function createGroup(uint256 newEntry) payable external;
}

/**
 * @title KeepRandomBeaconServiceImplV1
 * @dev Initial version of service contract that works under Keep Random
 * Beacon proxy and allows upgradability. The purpose of the contract is to have
 * up-to-date logic for threshold random number generation. Updated contracts
 * must inherit from this contract and have to be initialized under updated version name
 */
contract KeepRandomBeaconServiceImplV1 is Ownable, DelayedWithdrawal {
    using SafeMath for uint256;
    using AddressArrayUtils for address[];

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 requestId);
    event RelayEntryGenerated(uint256 requestId, uint256 entry);

    // Minimum gas price for relay entry request.
    uint256 internal _minGasPrice;

    // Fluctuation safety factor to cover the immediate rise in gas fees during DKG execution.
    // Must be presented as a big number with 18 decimals i.e. 1.5% as 1.5*1e18.
    uint256 internal _fluctuationMargin;

    // Fraction in % of the estimated cost of DKG that is included
    // in relay request fee. Must be presented as a big number with
    // 18 decimals i.e. 1.5% as 1.5*1e18.
    uint256 internal _dkgContributionMargin;

    // Every relay request payment includes DKG contribution that is added to
    // the DKG fee pool, once the pool amount reaches DKG cost estimate the relay
    // entry will trigger the creation of a new group.
    uint256 internal _dkgFeePool;

    // Rewards not paid out to the operators are sent to request subsidy pool to
    // subsidize new requests: 1% is returned to the requester's surplus address.
    uint256 internal _requestSubsidyFeePool;

    uint256 internal _previousEntry;

    // Each service contract tracks its own requests and these are independent
    // from operator contracts which track signing requests instead.
    uint256 internal _requestCounter;

    struct Callback {
        address callbackContract;
        string callbackMethod;
        uint256 callbackFee;
        uint256 callbackGas;
        address payable surplusRecipient;
    }

    mapping(uint256 => Callback) internal _callbacks;

    address[] internal _operatorContracts;

    // Mapping to store new implementation versions that inherit from this contract.
    mapping (string => bool) internal _initialized;

    // Seed used as the first random beacon value.
    // It is a signature over 78 digits of PI and 78 digits of Euler's number
    // using BLS private key 123.
    uint256 constant internal _beaconSeed = 10920102476789591414949377782104707130412218726336356788412941355500907533021;

    /**
     * @dev Initialize Keep Random Beacon service contract implementation.
     * @param minGasPrice Minimum gas price for relay entry request.
     * @param fluctuationMargin Fluctuation safety factor to cover the immediate rise in gas fees during 
     * DKG execution. Must be presented as a big number with 18 decimals i.e. 1.5% as 1.5*1e18.
     * @param dkgContributionMargin Fraction in % of the estimated cost of DKG that is included in relay
     * request fee. Must be presented as a big number with 18 decimals i.e. 1.5% as 1.5*1e18.
     * @param withdrawalDelay Delay before the owner can withdraw ether from this contract.
     * @param operatorContract Operator contract linked to this contract.
     */
    function initialize(
        uint256 minGasPrice,
        uint256 fluctuationMargin,
        uint256 dkgContributionMargin,
        uint256 withdrawalDelay,
        address operatorContract
    )
        public
        onlyOwner
    {
        require(!initialized(), "Contract is already initialized.");
        _initialized["KeepRandomBeaconServiceImplV1"] = true;
        _minGasPrice = minGasPrice;
        _fluctuationMargin = fluctuationMargin;
        _dkgContributionMargin = dkgContributionMargin;
        _withdrawalDelay = withdrawalDelay;
        _pendingWithdrawal = 0;
        _operatorContracts.push(operatorContract);

        _previousEntry = _beaconSeed;
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return _initialized["KeepRandomBeaconServiceImplV1"];
    }

    /**
     * @dev Adds operator contract
     * @param operatorContract Address of the operator contract.
     */
    function addOperatorContract(address operatorContract) public onlyOwner {
        _operatorContracts.push(operatorContract);
    }

    /**
     * @dev Removes operator contract
     * @param operatorContract Address of the operator contract.
     */
    function removeOperatorContract(address operatorContract) public onlyOwner {
        _operatorContracts.removeAddress(operatorContract);
    }

    /**
     * @dev Add funds to DKG fee pool.
     */
    function fundDkgFeePool() public payable {
        _dkgFeePool += msg.value;
    }

    /**
     * @dev Add funds to request subsidy fee pool.
     */
    function fundRequestSubsidyFeePool() public payable {
        _requestSubsidyFeePool += msg.value;
    }

    /**
     * @dev Selects an operator contract from the available list using modulo operation
     * with seed value weighted by the number of active groups on each operator contract.
     * @param seed Cryptographically generated random value.
     * @return Address of operator contract.
     */
    function selectOperatorContract(uint256 seed) public view returns (address) {

        uint256 totalNumberOfGroups;

        for (uint i = 0; i < _operatorContracts.length; i++) {
            totalNumberOfGroups += OperatorContract(_operatorContracts[i]).numberOfGroups();
        }

        require(totalNumberOfGroups > 0, "Total number of groups must be greater than zero.");

        uint256 selectedIndex = seed % totalNumberOfGroups;

        uint256 selectedContract;
        uint256 indexByGroupCount;

        for (uint256 i = 0; i < _operatorContracts.length; i++) {
            indexByGroupCount += OperatorContract(_operatorContracts[i]).numberOfGroups();
            if (selectedIndex < indexByGroupCount) {
                return _operatorContracts[selectedContract];
            }
            selectedContract++;
        }

        return _operatorContracts[selectedContract];
    }

    /**
     * @dev Creates a request to generate a new relay entry, which will include a
     * random number (by signing the previous entry's random number).
     * @param seed Initial seed random value from the client. It should be a cryptographically generated random value.
     * @return An uint256 representing uniquely generated entry Id.
     */
    function requestRelayEntry(uint256 seed) public payable returns (uint256) {
        return requestRelayEntry(seed, address(0), "", 0);
    }

    /**
     * @dev Creates a request to generate a new relay entry, which will include a
     * random number (by signing the previous entry's random number).
     * @param seed Initial seed random value from the client. It should be a cryptographically generated random value.
     * @param callbackContract Callback contract address. Callback is called once a new relay entry has been generated.
     * @param callbackMethod Callback contract method signature. String representation of your method with a single
     * uint256 input parameter i.e. "relayEntryCallback(uint256)".
     * @param callbackGas Gas required for the callback.
     * The customer needs to ensure they provide a sufficient callback gas
     * to cover the gas fee of executing the callback. Any surplus is returned
     * to the customer. If the callback gas amount turns to be not enough to
     * execute the callback, callback execution is skipped.
     * @return An uint256 representing uniquely generated relay request ID. It is also returned as part of the event.
     */
    function requestRelayEntry(
        uint256 seed,
        address callbackContract,
        string memory callbackMethod,
        uint256 callbackGas
    ) public payable returns (uint256) {
        require(
            msg.value >= entryFeeEstimate(callbackGas),
            "Payment is less than required minimum."
        );

        (uint256 entryVerificationFee, uint256 dkgContributionFee, uint256 groupProfitFee) = entryFeeBreakdown();
        uint256 callbackFee = msg.value.sub(entryVerificationFee).sub(dkgContributionFee).sub(groupProfitFee);

        _dkgFeePool += dkgContributionFee;

        _requestCounter++;
        uint256 requestId = _requestCounter;

        OperatorContract(selectOperatorContract(_previousEntry)).sign.value(
            entryVerificationFee.add(groupProfitFee)
        )(requestId, seed, _previousEntry);

        if (callbackContract != address(0)) {
            _callbacks[requestId] = Callback(callbackContract, callbackMethod, callbackFee, callbackGas, msg.sender);
        }

        // Send 1% of the request subsidy pool to the requestor.
        if (_requestSubsidyFeePool >= 100) {
            uint256 amount = _requestSubsidyFeePool.div(100);
            _requestSubsidyFeePool -= amount;
            msg.sender.transfer(amount);
        }

        emit RelayEntryRequested(requestId);
        return requestId;
    }

    /**
     * @dev Store valid entry returned by operator contract and call customer specified callback if required.
     * @param requestId Request id tracked internally by this contract.
     * @param entry The generated random number.
     * @param submitter Relay entry submitter.
     */
    function entryCreated(uint256 requestId, uint256 entry, address payable submitter) public {
        require(
            _operatorContracts.contains(msg.sender),
            "Only authorized operator contract can call relay entry."
        );

        _previousEntry = entry;
        emit RelayEntryGenerated(requestId, entry);

        if (_callbacks[requestId].callbackContract != address(0)) {
            executeEntryCreatedCallback(requestId, entry, submitter);
        }

        triggerDkgIfApplicable(entry);
    }

    /**
     * @dev Executes customer specified callback for the relay entry request.
     * @param requestId Request id tracked internally by this contract.
     * @param entry The generated random number.
     * @param submitter Relay entry submitter.
     */
    function executeEntryCreatedCallback(uint256 requestId, uint256 entry, address payable submitter) internal {
        bool success; // Store status of external contract call.
        bytes memory data; // Store result data of external contract call.

        uint256 gasBeforeCallback = gasleft();
        (success, data) = _callbacks[requestId].callbackContract.call.gas(_callbacks[requestId].callbackGas)(abi.encodeWithSignature(_callbacks[requestId].callbackMethod, entry));
        uint256 gasSpent = gasBeforeCallback.sub(gasleft()).add(21000); // Also reimburse 21000 gas (ethereum transaction minimum gas)

        // Obtain the actual callback gas expenditure and refund the surplus.
        uint256 callbackSurplus = 0;
        uint256 callbackFee = gasSpent.mul(tx.gasprice);

        // If we spent less on the callback than the customer transferred for the
        // callback execution, we need to reimburse the difference.
        if (callbackFee < _callbacks[requestId].callbackFee) {
            callbackSurplus = _callbacks[requestId].callbackFee.sub(callbackFee);
            // Reimburse submitter with his actual callback cost.
            submitter.transfer(callbackFee);
            // Return callback surplus to the requestor.
            _callbacks[requestId].surplusRecipient.transfer(callbackSurplus);
        } else {
            // Reimburse submitter with the callback payment sent by the requestor.
            submitter.transfer(_callbacks[requestId].callbackFee);
        }

        delete _callbacks[requestId];
    }

    /**
     * @dev Triggers the selection process of a new candidate group.
     * @param entry The generated random number.
     */
    function triggerDkgIfApplicable(uint256 entry) internal {
        bool success; // Store status of external contract call.
        bytes memory data; // Store result data of external contract call.

        address latestOperatorContract = _operatorContracts[_operatorContracts.length.sub(1)];
        uint256 dkgFeeEstimate = _minGasPrice.mul(OperatorContract(latestOperatorContract).dkgGasEstimate()).mul(_fluctuationMargin).div(1e18);
        if (_dkgFeePool >= dkgFeeEstimate) {
            _dkgFeePool = _dkgFeePool.sub(dkgFeeEstimate);
            (success, data) = latestOperatorContract.call.value(dkgFeeEstimate)(abi.encodeWithSignature("createGroup(uint256)", entry));
        }
    }

    /**
     * @dev Set the minimum gas price in wei for estimating relay entry request payment.
     * @param minGasPrice is the minimum gas price required for estimating relay entry request payment.
     */
    function setMinimumGasPrice(uint256 minGasPrice) public onlyOwner {
        _minGasPrice = minGasPrice;
    }

    /**
     * @dev Get the minimum gas price in wei that is used to estimate relay entry request payment.
     */
    function minimumGasPrice() public view returns(uint256) {
        return _minGasPrice;
    }

    /**
     * @dev Get the minimum payment in wei for relay entry callback.
     * @param callbackGas Gas required for the callback.
     */
    function minimumCallbackFee(uint256 callbackGas) public view returns(uint256) {
        return callbackGas.mul(_minGasPrice).mul(_fluctuationMargin).div(1e18);
    }

    /**
     * @dev Get the entry fee estimate in wei for relay entry request.
     * @param callbackGas Gas required for the callback.
     */
    function entryFeeEstimate(uint256 callbackGas) public view returns(uint256) {
        (uint256 entryVerificationFee, uint256 dkgContributionFee, uint256 groupProfitFee) = entryFeeBreakdown();
        return entryVerificationFee.add(dkgContributionFee).add(groupProfitFee).add(minimumCallbackFee(callbackGas));
    }

    /**
     * @dev Get the entry fee breakdown in wei for relay entry request.
     */
    function entryFeeBreakdown() public view returns(
        uint256 entryVerificationFee,
        uint256 dkgContributionFee,
        uint256 groupProfitFee
    ) {
        uint256 entryVerificationGas;

        // Use most expensive operator contract for estimated entry verification gas value and group profit fee.
        for (uint i = 0; i < _operatorContracts.length; i++) {
            OperatorContract operator = OperatorContract(_operatorContracts[i]);

            if (operator.numberOfGroups() > 0) {
                entryVerificationGas = operator.entryVerificationGasEstimate() > entryVerificationGas ? operator.entryVerificationGasEstimate():entryVerificationGas;
                groupProfitFee = operator.groupProfitFee() > groupProfitFee ? operator.groupProfitFee():groupProfitFee;
            }
        }

        // Use DKG gas estimate from the latest operator contract since it will be used for the next group creation.
        address latestOperatorContract = _operatorContracts[_operatorContracts.length.sub(1)];
        uint256 dkgGas = OperatorContract(latestOperatorContract).dkgGasEstimate();

        return (
            entryVerificationGas.mul(_minGasPrice),
            dkgGas.mul(_minGasPrice.mul(_fluctuationMargin).div(1e18)).mul(_dkgContributionMargin).div(100).div(1e18),
            groupProfitFee
        );
    }

    /**
     * @dev Gets the previous relay entry value.
     */
    function previousEntry() public view returns(uint256) {
        return _previousEntry;
    }

    /**
     * @dev Gets version of the current implementation.
     */
    function version() public pure returns (string memory) {
        return "V1";
    }
}
