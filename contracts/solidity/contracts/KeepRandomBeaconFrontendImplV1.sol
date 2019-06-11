pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./DelayedWithdrawal.sol";


interface OperatorContract {
    function requestRelayEntry(address from, uint256 seed, uint256 previousEntry) payable external returns (uint256 requestId);
    function numberOfGroups() external view returns(uint256);
    function selectGroup(uint256 previousEntry) external returns(bytes memory);
}


/**
 * @title KeepRandomBeaconFrontendImplV1
 * @dev Initial version of implementation contract that works under Keep Random
 * Beacon proxy and allows upgradability. The purpose of the contract is to have
 * up-to-date logic for threshold random number generation. Updated contracts
 * must inherit from this contract and have to be initialized under updated version name
 */
contract KeepRandomBeaconFrontendImplV1 is Ownable, DelayedWithdrawal {

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 requestID, uint256 payment, uint256 previousEntry, uint256 seed, bytes groupPublicKey); 
    event RelayEntryGenerated(uint256 requestID, uint256 requestResponse, bytes requestGroupPubKey, uint256 previousEntry, uint256 seed);

    uint256 internal _minPayment;
    address internal _operatorContract;
    uint256 internal _previousEntry;
    uint256 internal _relayRequestTimeout;

    mapping (string => bool) internal _initialized;

    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() external payable {
        revert("Can not call contract without explicitly calling a function.");
    }

    /**
     * @dev Initialize Keep Random Beacon implementaion contract.
     * @param minPayment Minimum amount of ether (in wei) that allows anyone to request a random number.
     * @param withdrawalDelay Delay before the owner can withdraw ether from this contract.
     * @param operatorContract Operator contract linked to this contract.
     * @param relayRequestTimeout Timeout in blocks for a relay entry to appear on the chain.
     * Blocks are counted from the moment relay request occur.
     */
    function initialize(uint256 minPayment, uint256 withdrawalDelay, address operatorContract, uint256 relayRequestTimeout)
        public
        onlyOwner
    {
        require(!initialized(), "Contract is already initialized.");
        _minPayment = minPayment;
        _initialized["KeepRandomBeaconFrontendImplV1"] = true;
        _withdrawalDelay = withdrawalDelay;
        _pendingWithdrawal = 0;
        _operatorContract = operatorContract;
        _relayRequestTimeout = relayRequestTimeout;
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return _initialized["KeepRandomBeaconFrontendImplV1"];
    }

    /**
     * @dev Creates a request to generate a new relay entry, which will include a
     * random number (by signing the previous entry's random number).
     * @param seed Initial seed random value from the client. It should be a cryptographically generated random value.
     * @return An uint256 representing uniquely generated relay request ID. It is also returned as part of the event.
     */
    function requestRelayEntry(uint256 seed) public payable returns (uint256) {
        require(
            msg.value >= _minPayment,
            "Payment is less than required minimum."
        );

        // TODO: Figure out pricing, if we decide to pass payment to the backed use this instead:
        // OperatorContract(_operatorContract).requestRelayEntry.value(msg.value)(msg.sender, seed, _previousEntry);
        return OperatorContract(_operatorContract).requestRelayEntry(msg.sender, seed, _previousEntry);
    }

    /**
     * @dev Creates a new relay entry and stores the associated data on the chain.
     * @param requestID The request that started this generation - to tie the results back to the request.
     * @param groupSignature The generated random number.
     * @param groupPubKey Public key of the group that generated the threshold signature.
     * @param previousEntry Previous relay entry value.
     * @param seed Initial seed random value from the client. It should be a cryptographically generated random value.
     */
    function relayEntry(uint256 requestID, uint256 groupSignature, bytes memory groupPubKey, uint256 previousEntry, uint256 seed) public {
        require(
            msg.sender == _operatorContract,
            "Only authorized operator contract can call relay entry."
        );

        _previousEntry = groupSignature;
        emit RelayEntryGenerated(requestID, groupSignature, groupPubKey, previousEntry, seed);
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
     * Gets the timeout in blocks for a relay entry to appear on the chain.
     */
    function relayRequestTimeout() public view returns(uint256) {
        return _relayRequestTimeout;
    }

    /**
     * @dev Gets version of the current implementation.
     */
    function version() public pure returns (string memory) {
        return "V1";
    }
}
