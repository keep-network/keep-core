pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./utils/AddressArrayUtils.sol";
import "./DelayedWithdrawal.sol";


interface OperatorContract {
    function sign(uint256 requestId, uint256 seed, uint256 previousEntry) payable external;
    // TODO: Add createGroup() when it's implemented on Operator contract.
}


/**
 * @title KeepRandomBeaconServiceImplV1
 * @dev Initial version of service contract that works under Keep Random
 * Beacon proxy and allows upgradability. The purpose of the contract is to have
 * up-to-date logic for threshold random number generation. Updated contracts
 * must inherit from this contract and have to be initialized under updated version name
 */
contract KeepRandomBeaconServiceImplV1 is Ownable, DelayedWithdrawal {

    using AddressArrayUtils for address[];

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 requestId);
    event RelayEntryGenerated(uint256 requestId, uint256 entry);

    uint256 internal _minPayment;
    uint256 internal _previousEntry;

    // Each service contract tracks its own requests and these are independent
    // from operator contracts which track signing requests instead.
    uint256 internal _requestCounter;

    struct Callback {
        address callbackContract;
        string callbackMethod;
    }

    mapping(uint256 => Callback) internal _callbacks;

    address[] internal _operatorContracts;

    // Mapping to store new implementation versions that inherit from this contract.
    mapping (string => bool) internal _initialized;

    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() external payable {
        revert("Can not call contract without explicitly calling a function.");
    }

    /**
     * @dev Initialize Keep Random Beacon service contract implementation.
     * @param minPayment Minimum amount of ether (in wei) that allows anyone to request a random number.
     * @param withdrawalDelay Delay before the owner can withdraw ether from this contract.
     * @param operatorContract Operator contract linked to this contract.
     */
    function initialize(uint256 minPayment, uint256 withdrawalDelay, address operatorContract)
        public
        onlyOwner
    {
        require(!initialized(), "Contract is already initialized.");
        _minPayment = minPayment;
        _initialized["KeepRandomBeaconServiceImplV1"] = true;
        _withdrawalDelay = withdrawalDelay;
        _pendingWithdrawal = 0;
        _operatorContracts.push(operatorContract);
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return _initialized["KeepRandomBeaconServiceImplV1"];
    }

    /**
     * @dev Selects an operator contract from the available list using modulo operation
     * with previous entry weighted by the number of active groups on each operator contract.
     * @return Address of operator contract.
     */
    function selectOperatorContract() public view returns (address) {
        // TODO: Implement logic
        return _operatorContracts[0];
    }

    /**
     * @dev Creates a request to generate a new relay entry, which will include a
     * random number (by signing the previous entry's random number).
     * @param seed Initial seed random value from the client. It should be a cryptographically generated random value.
     * @return An uint256 representing uniquely generated entry Id.
     */
    function requestRelayEntry(uint256 seed) public payable returns (uint256) {
        return requestRelayEntry(seed, address(0), "");
    }

    /**
     * @dev Creates a request to generate a new relay entry, which will include a
     * random number (by signing the previous entry's random number).
     * @param seed Initial seed random value from the client. It should be a cryptographically generated random value.
     * @param callbackContract Callback contract address. Callback is called once a new relay entry has been generated.
     * @param callbackMethod Callback contract method signature. String representation of your method with a single
     * uint256 input parameter i.e. "relayEntryCallback(uint256)".
     * @return An uint256 representing uniquely generated relay request ID. It is also returned as part of the event.
     */
    function requestRelayEntry(uint256 seed, address callbackContract, string memory callbackMethod) public payable returns (uint256) {
        require(
            msg.value >= _minPayment,
            "Payment is less than required minimum."
        );

        _requestCounter++;

        // TODO: Figure out pricing, if we decide to pass payment to the backed use this instead:
        // OperatorContract(selectOperatorContract()).sign.value(msg.value)(_requestCounter, seed, _previousEntry);
        OperatorContract(selectOperatorContract()).sign(_requestCounter, seed, _previousEntry);

        if (callbackContract != address(0)) {
            _callbacks[_requestCounter] = Callback(callbackContract, callbackMethod);
        }

        emit RelayEntryRequested(_requestCounter);
        return _requestCounter;
    }

    /**
     * @dev Store valid entry returned by operator contract and call customer specified callback if required.
     * @param requestId Request id tracked internally by this contract.
     * @param entry The generated random number.
     */
    function entryCreated(uint256 requestId, uint256 entry) public {
        require(
            _operatorContracts.contains(msg.sender),
            "Only authorized operator contract can call relay entry."
        );

        _previousEntry = entry;
        emit RelayEntryGenerated(requestId, entry);

        if (_callbacks[requestId].callbackContract != address(0)) {
            _callbacks[requestId].callbackContract.call(abi.encodeWithSignature(_callbacks[requestId].callbackMethod, entry));
            delete _callbacks[requestId];
        }
    }

    /**
     * @dev Set the minimum payment that is required before a relay entry occurs.
     * @param minPayment is the value in wei that is required to be payed for the process to start.
     */
    function setMinimumPayment(uint256 minPayment) public onlyOwner {
        _minPayment = minPayment;
    }

    /**
     * @dev Get the minimum payment that is required before a relay entry occurs.
     */
    function minimumPayment() public view returns(uint256) {
        return _minPayment;
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
