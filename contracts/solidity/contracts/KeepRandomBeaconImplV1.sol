pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./StakingProxy.sol";


/**
 * @title KeepRandomBeaconImplV1
 * @dev Initial version of implementation contract that works under Keep Random
 * Beacon proxy and allows upgradability. The purpose of the contract is to have
 * up-to-date logic for threshold random number generation. Updated contracts
 * must inherit from this contract and have to be initialized under updated version name
 */
contract KeepRandomBeaconImplV1 is Ownable {

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 requestID, uint256 payment, uint256 blockReward, uint256 seed, uint blockNumber); 
    event RelayEntryGenerated(uint256 requestID, uint256 requestResponse, uint256 requestGroupID, uint256 previousEntry, uint blockNumber);
    event RelayResetEvent(uint256 lastValidRelayEntry, uint256 lastValidRelayTxHash, uint256 lastValidRelayBlock);
    event SubmitGroupPublicKeyEvent(byte[] groupPublicKey, uint256 requestID, uint256 activationBlockHeight);

    uint256 internal _seq;
    uint256 internal _minPayment;
    uint256 internal _minStake;
    address internal _stakingProxy;
    uint256 internal _withdrawalDelay;
    uint256 internal _pendingWithdrawal;

    mapping (string => bool) internal _initialized;
    mapping (uint256 => address) internal _requestPayer;
    mapping (uint256 => uint256) internal _requestPayment;
    mapping (uint256 => uint256) internal _blockReward;
    mapping (uint256 => uint256) internal _requestGroup;

    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() public payable {
        revert("Can not call contract without explicitly calling a function.");
    }

    /**
     * @dev Initialize Keep Random Beacon implementaion contract with a linked staking proxy contract.
     * @param stakingProxy Address of a staking proxy contract that will be linked to this contract.
     * @param minPayment Minimum amount of ether (in wei) that allows anyone to request a random number.
     * @param minStake Minimum amount in KEEP that allows KEEP network client to participate in a group.
     * @param withdrawalDelay Delay before the owner can withdraw ether from this contract.
     */
    function initialize(address stakingProxy, uint256 minPayment, uint256 minStake, uint256 withdrawalDelay)
        public
        onlyOwner
    {
        require(!initialized(), "Contract is already initialized.");
        require(stakingProxy != address(0x0), "Staking proxy address can't be zero.");
        _stakingProxy = stakingProxy;
        _minStake = minStake;
        _minPayment = minPayment;
        _initialized["KeepRandomBeaconImplV1"] = true;
        _withdrawalDelay = withdrawalDelay;
        _pendingWithdrawal = 0;
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return _initialized["KeepRandomBeaconImplV1"];
    }

    /**
     * @dev Checks that the specified user has an appropriately large stake.
     * @param staker Specifies the identity of the random beacon client.
     * @return True if staked enough to participate in the group, false otherwise.
     */
    function hasMinimumStake(address staker) public view returns(bool) {
        uint256 balance;
        balance = StakingProxy(_stakingProxy).balanceOf(staker);
        return (balance >= _minStake);
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

        emit RelayEntryRequested(requestID, msg.value, blockReward, seed, block.number);
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
    function finishWithdrawal() public onlyOwner {
        require(_pendingWithdrawal > 0, "Pending withdrawal timestamp must be set and be greater than zero.");
        require(block.timestamp >= _pendingWithdrawal, "The current time must pass the pending withdrawal timestamp.");

        // Reset pending withdrawal before sending to prevent re-entrancy attacks
        _pendingWithdrawal = 0;
        owner().transfer(address(this).balance);
    }

    /**
     * @dev Set the minimum payment that is required before a relay entry occurs.
     * @param minPayment is the value in wei that is required to be payed for the process to start.
     */
    function setMinimumPayment(uint256 minPayment) public onlyOwner {
        _minPayment = minPayment;
    }

    /**
     * @dev Set the minimum amount of KEEP that allows a Keep network client to participate in a group.
     * @param minStake Amount in KEEP.
     */
    function setMinimumStake(uint256 minStake) public onlyOwner {
        _minStake = minStake;
    }

    /**
     * @dev Get the minimum payment that is required before a relay entry occurs.
     */
    function minimumPayment() public view returns(uint256) {
        return _minPayment;
    }

    /**
     * @dev Get the minimum amount in KEEP that allows KEEP network client to participate in a group.
     */
    function minimumStake() public view returns(uint256) {
        return _minStake;
    }

    /**
     * @dev Creates a new relay entry and stores the associated data on the chain.
     * @param requestID The request that started this generation - to tie the results back to the request.
     * @param groupSignature The generated random number.
     * @param groupID Public key of the group that generated the threshold signature.
     */
    function relayEntry(uint256 requestID, uint256 groupSignature, uint256 groupID, uint256 previousEntry) public {
        _requestGroup[requestID] = groupID;

        emit RelayEntryGenerated(requestID, groupSignature, groupID, previousEntry, block.number);
    }

    /**
     * @dev Takes a generated key and place it on the blockchain. Creates an event.
     * @param groupPublicKey Group public key.
     * @param requestID Request ID.
     */
    function submitGroupPublicKey(byte[] groupPublicKey, uint256 requestID) public {
        uint256 activationBlockHeight = block.number;

        // TODO -- lots of stuff - don't know yet.
        emit SubmitGroupPublicKeyEvent(groupPublicKey, requestID, activationBlockHeight);
    }

    /**
     * @dev Gets version of the current implementation.
    */
    function version() public pure returns (string) {
        return "V1";
    }
}
