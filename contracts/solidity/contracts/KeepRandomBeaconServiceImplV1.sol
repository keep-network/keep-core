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
        bytes calldata previousEntry
    ) external payable;
    function numberOfGroups() external view returns(uint256);
    function createGroup(uint256 newEntry) external payable;
}

/**
 * @title KeepRandomBeaconServiceImplV1
 * @dev Initial version of service contract that works under Keep Random
 * Beacon proxy and allows upgradability. The purpose of the contract is to have
 * up-to-date logic for threshold random number generation. Updated contracts
 * must inherit from this contract and have to be initialized under updated version name
 * Warning: you can't set constants directly in the contract and must use initialize()
 * please see openzeppelin upgradeable contracts approach for more info.
 */
contract KeepRandomBeaconServiceImplV1 is Ownable, DelayedWithdrawal {
    using SafeMath for uint256;
    using AddressArrayUtils for address[];

    event RelayEntryRequested(uint256 requestId);
    event RelayEntryGenerated(uint256 requestId, uint256 entry);

    // The price feed estimate is used to calculate the gas price for reimbursement
    // next to the actual gas price from the transaction. We use both values to
    // defend against malicious miner-submitters who can manipulate transaction
    // gas price. Expressed in wei.
    uint256 internal _priceFeedEstimate;

    // Fluctuation margin to cover the immediate rise in gas price.
    // Expressed in percentage.
    uint256 internal _fluctuationMargin;

    // Fraction in % of the estimated cost of DKG that is included
    // in relay request fee.
    uint256 internal _dkgContributionMargin;

    // Every relay request payment includes DKG contribution that is added to
    // the DKG fee pool, once the pool value reaches the required minimum, a new
    // relay entry will trigger the creation of a new group. Expressed in wei.
    uint256 internal _dkgFeePool;

    // Rewards not paid out to the operators are sent to request subsidy pool to
    // subsidize new requests: 1% of the subsidy pool is returned to the requester's
    // surplus address. Expressed in wei.
    uint256 internal _requestSubsidyFeePool;

    // Each service contract tracks its own requests and these are independent
    // from operator contracts which track signing requests instead.
    uint256 internal _requestCounter;

    bytes internal _previousEntry;

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
    // It's a G1 point G * PI =
    // G * 31415926535897932384626433832795028841971693993751058209749445923078164062862
    // Where G is the generator of G1 abstract cyclic group.
    bytes constant internal _beaconSeed =
    hex"15c30f4b6cf6dbbcbdcc10fe22f54c8170aea44e198139b776d512d8f027319a1b9e8bfaf1383978231ce98e42bafc8129f473fc993cf60ce327f7d223460663";

    /**
     * @dev Initialize Keep Random Beacon service contract implementation.
     * @param priceFeedEstimate The price feed estimate is used to calculate the gas price for
     * reimbursement next to the actual gas price from the transaction. We use both values to defend
     * against malicious miner-submitters who can manipulate transaction gas price.
     * @param fluctuationMargin Fluctuation margin to cover the immediate rise in gas price.
     * Expressed in percentage.
     * @param dkgContributionMargin Fraction in % of the estimated cost of DKG that is included in relay
     * request fee.
     * @param withdrawalDelay Delay before the owner can withdraw ether from this contract.
     * @param operatorContract Operator contract linked to this contract.
     */
    function initialize(
        uint256 priceFeedEstimate,
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
        _priceFeedEstimate = priceFeedEstimate;
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
     * @dev Creates a request to generate a new relay entry, which will include
     * a random number (by signing the previous entry's random number).
     * @return An uint256 representing uniquely generated entry Id.
     */
    function requestRelayEntry() public payable returns (uint256) {
        return requestRelayEntry(address(0), "", 0);
    }

    /**
     * @dev Creates a request to generate a new relay entry, which will include
     * a random number (by signing the previous entry's random number).
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

        OperatorContract operatorContract = OperatorContract(
            selectOperatorContract(uint256(keccak256(_previousEntry)))
        );
        uint256 selectedOperatorContractFee = operatorContract.groupProfitFee().add(
            operatorContract.entryVerificationGasEstimate().mul(gasPriceWithFluctuationMargin(_priceFeedEstimate)));

        _requestCounter++;
        uint256 requestId = _requestCounter;

        operatorContract.sign.value(
            selectedOperatorContractFee
        )(requestId, _previousEntry);

        // If selected operator contract is cheaper than expected return the
        // surplus to the subsidy fee pool.
        // We do that instead of returning the surplus to the requestor to have
        // a consistent beacon pricing for customers without fluctuations caused
        // by different operator contracts being selected.
        uint256 surplus = entryVerificationFee.add(groupProfitFee).sub(selectedOperatorContractFee);
        _requestSubsidyFeePool = _requestSubsidyFeePool.add(surplus);

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
    function entryCreated(uint256 requestId, bytes memory entry, address payable submitter) public {
        require(
            _operatorContracts.contains(msg.sender),
            "Only authorized operator contract can call relay entry."
        );

        _previousEntry = entry;
        uint256 entryAsNumber = uint256(keccak256(entry));
        emit RelayEntryGenerated(requestId, entryAsNumber);

        if (_callbacks[requestId].callbackContract != address(0)) {
            executeEntryCreatedCallback(requestId, entryAsNumber, submitter);
            delete _callbacks[requestId];
        }

        triggerDkgIfApplicable(entryAsNumber);
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
        (success, data) = _callbacks[requestId].callbackContract.call.gas(
            _callbacks[requestId].callbackGas
        )(abi.encodeWithSignature(_callbacks[requestId].callbackMethod, entry));
        uint256 gasSpent = gasBeforeCallback.sub(gasleft()).add(21000); // Also reimburse 21000 gas (ethereum transaction minimum gas)

        uint256 gasPrice = _priceFeedEstimate;
        // We need to check if tx.gasprice is non-zero as a workaround to a bug
        // in go-ethereum:
        // https://github.com/ethereum/go-ethereum/pull/20189
        if (tx.gasprice > 0 && tx.gasprice < _priceFeedEstimate) {
            gasPrice = tx.gasprice;
        }

        // Obtain the actual callback gas expenditure and refund the surplus.
        uint256 callbackSurplus = 0;
        uint256 callbackFee = gasSpent.mul(gasPrice);

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
    }

    /**
     * @dev Triggers the selection process of a new candidate group if the DKG
     * fee pool equals or exceeds DKG cost estimate.
     * @param entry The generated random number.
     */
    function triggerDkgIfApplicable(uint256 entry) internal {
        bool success; // Store status of external contract call.
        bytes memory data; // Store result data of external contract call.

        address latestOperatorContract = _operatorContracts[_operatorContracts.length.sub(1)];
        uint256 dkgFeeEstimate = OperatorContract(latestOperatorContract).dkgGasEstimate().mul(
            gasPriceWithFluctuationMargin(_priceFeedEstimate)
        );
        if (_dkgFeePool >= dkgFeeEstimate) {
            _dkgFeePool = _dkgFeePool.sub(dkgFeeEstimate);
            // Disabling ethlint error message. No security implications, we're calling audited and trusted contract here.
            // solium-disable-next-line
            (success, data) = latestOperatorContract.call.value(dkgFeeEstimate)(abi.encodeWithSignature("createGroup(uint256)", entry));
        }
    }

    /**
     * @dev Set the gas price in wei for estimating relay entry request payment.
     * @param priceFeedEstimate is the gas price required for estimating relay entry request payment.
     */
    function setPriceFeedEstimate(uint256 priceFeedEstimate) public onlyOwner {
        _priceFeedEstimate = priceFeedEstimate;
    }

    /**
     * @dev Get the gas price in wei that is used to estimate relay entry request payment.
     */
    function priceFeedEstimate() public view returns(uint256) {
        return _priceFeedEstimate;
    }

    /**
     * @dev Adds a safety margin for gas price fluctuations to the current gas price.
     * The gas price for DKG or relay entry is set when the request is processed
     * but the result submission transaction will be sent later. We add a safety
     * margin that should be sufficient for getting requests processed within a
     * a deadline under all circumstances.
     * @param gasPrice Gas price in wei.
     */
    function gasPriceWithFluctuationMargin(uint256 gasPrice) internal view returns (uint256) {
        return gasPrice.add(gasPrice.mul(_fluctuationMargin).div(100));
    }

    /**
     * @dev Get the minimum payment in wei for relay entry callback.
     * The returned value includes safety margin for gas price fluctuations.
     * @param callbackGas Gas required for the callback.
     */
    function callbackFee(uint256 callbackGas) public view returns(uint256) {
        // We take the gas price from the price feed to not let malicious
        // miner-requestors manipulate the gas price when requesting relay entry
        // and underpricing expensive callbacks.
        return callbackGas.mul(gasPriceWithFluctuationMargin(_priceFeedEstimate));
    }

    /**
     * @dev Get the entry fee estimate in wei for relay entry request.
     * @param callbackGas Gas required for the callback.
     */
    function entryFeeEstimate(uint256 callbackGas) public view returns(uint256) {
        (uint256 entryVerificationFee, uint256 dkgContributionFee, uint256 groupProfitFee) = entryFeeBreakdown();
        return entryVerificationFee.add(dkgContributionFee).add(groupProfitFee).add(callbackFee(callbackGas));
    }

    /**
     * @dev Get the entry fee breakdown in wei for relay entry request.
     * Entry verification fee returned contains safety margin for gas price fluctuations.
     */
    function entryFeeBreakdown() public view returns(
        uint256 entryVerificationFee,
        uint256 dkgContributionFee,
        uint256 groupProfitFee
    ) {
        uint256 entryVerificationGas;

        // Select the most expensive entry verification from all the operator contracts
        // and the highest group profit fee from all the operator contracts. We do not
        // know what is going to be the gas price at the moment of submitting an entry,
        // thus we can't calculate at this point which contract is the most expensive
        // based on the entry verification gas and group profit fee. Hence, we need to
        // select maximum of both those values separately.
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
            entryVerificationGas.mul(gasPriceWithFluctuationMargin(_priceFeedEstimate)),
            dkgGas.mul(_priceFeedEstimate).mul(_dkgContributionMargin).div(100),
            groupProfitFee
        );
    }

    /**
     * @dev Gets the previous relay entry value.
     */
    function previousEntry() public view returns(bytes memory) {
        return _previousEntry;
    }

    /**
     * @dev Gets version of the current implementation.
     */
    function version() public pure returns (string memory) {
        return "V1";
    }
}
