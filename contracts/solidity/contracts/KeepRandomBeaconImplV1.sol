pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./BLS.sol";


interface GroupContract {
    function runGroupSelection(uint256 randomBeaconValue) external;
}


/**
 * @title KeepRandomBeaconImplV1
 * @dev Initial version of implementation contract that works under Keep Random
 * Beacon proxy and allows upgradability. The purpose of the contract is to have
 * up-to-date logic for threshold random number generation. Updated contracts
 * must inherit from this contract and have to be initialized under updated version name
 */
contract KeepRandomBeaconImplV1 is Ownable {

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 requestID, uint256 payment, uint256 blockReward, uint256 seed); 
    event RelayEntryGenerated(uint256 requestID, uint256 requestResponse, bytes requestGroupPubKey, uint256 previousEntry, uint256 seed);

    uint256 internal _seq;
    uint256 internal _minPayment;
    uint256 internal _withdrawalDelay;
    uint256 internal _pendingWithdrawal;
    address internal _groupContract;
    uint256 internal _previousEntry;

    mapping (string => bool) internal _initialized;
    mapping (uint256 => address) internal _requestPayer;
    mapping (uint256 => uint256) internal _requestPayment;
    mapping (uint256 => uint256) internal _blockReward;
    mapping (uint256 => bytes) internal _requestGroup;

    mapping (uint256 => bool) internal _relayEntryRequested;

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
     * @param genesisEntry Initial relay entry to create first group.
     * @param groupContract Group contract linked to this contract.
     */
    function initialize(uint256 minPayment, uint256 withdrawalDelay, uint256 genesisEntry, address groupContract)
        public
        onlyOwner
    {
        require(!initialized(), "Contract is already initialized.");
        _minPayment = minPayment;
        _initialized["KeepRandomBeaconImplV1"] = true;
        _withdrawalDelay = withdrawalDelay;
        _pendingWithdrawal = 0;
        _previousEntry = genesisEntry;
        _groupContract = groupContract;
        GroupContract(_groupContract).runGroupSelection(_previousEntry);
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return _initialized["KeepRandomBeaconImplV1"];
    }

    /**
     * @dev Creates a request to generate a new relay entry, which will include a
     * random number (by signing the previous entry's random number).
     * @param blockReward The value in KEEP for generating the signature.
     * @param seed Initial seed random value from the client. It should be a cryptographically generated random value.
     * @return An uint256 representing uniquely generated ID. It is also returned as part of the event.
     */
    function requestRelayEntry(uint256 blockReward, uint256 seed) public payable returns (uint256 requestID) {
        require(
            msg.value >= _minPayment,
            "Payment is less than required minimum."
        );

        requestID = _seq++;
        _requestPayer[requestID] = msg.sender;
        _requestPayment[requestID] = msg.value;
        _blockReward[requestID] = blockReward;

        emit RelayEntryRequested(requestID, msg.value, blockReward, seed);
        return requestID;
    }

    /**
     * @dev Return the current RequestID - used in testing.
     */
    function getRequestId() public view returns (uint256 requestID) {
        return _seq;
    }

    /**
     * @dev Initiate withdrawal of this contract balance to the owner.
     */
    function initiateWithdrawal() public onlyOwner {
        _pendingWithdrawal = block.timestamp + _withdrawalDelay;
    }

    /**
     * @dev Finish withdrawal of this contract balance to the owner.
     */
    function finishWithdrawal(address payable payee) public onlyOwner {
        require(_pendingWithdrawal > 0, "Pending withdrawal timestamp must be set and be greater than zero.");
        require(block.timestamp >= _pendingWithdrawal, "The current time must pass the pending withdrawal timestamp.");

        // Reset pending withdrawal before sending to prevent re-entrancy attacks
        _pendingWithdrawal = 0;
        payee.transfer(address(this).balance);
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
     * @dev Creates a new relay entry and stores the associated data on the chain.
     * @param requestID The request that started this generation - to tie the results back to the request.
     * @param groupSignature The generated random number.
     * @param groupPubKey Public key of the group that generated the threshold signature.
     */
    function relayEntry(uint256 requestID, uint256 groupSignature, bytes memory groupPubKey, uint256 previousEntry, uint256 seed) public {    
        // Temporary solution for M2. Every group member submits a new relay entry
        // with the same request ID and we filter out duplicates here. 
        // This behavior will change post-M2 when we'll integrate phase 14 and/or 
        // implement relay requests.
        if (_relayEntryRequested[requestID]) {
            return;
        }
        _relayEntryRequested[requestID] = true;

        require(BLS.verify(groupPubKey, abi.encodePacked(previousEntry), bytes32(groupSignature)));

        _requestGroup[requestID] = groupPubKey;
        emit RelayEntryGenerated(requestID, groupSignature, groupPubKey, previousEntry, seed);
        GroupContract(_groupContract).runGroupSelection(groupSignature);
    }

    /**
     * @dev Gets version of the current implementation.
    */
    function version() public pure returns (string memory) {
        return "V1";
    }
}
