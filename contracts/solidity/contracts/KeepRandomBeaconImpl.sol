pragma solidity ^0.4.18;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";
import "./StakingProxy.sol";
import "./EternalStorage.sol";


/**
 * @title KeepRandomBeaconImpl
 * @dev Implementation contract that works under Keep Random Beacon proxy to
 * allow upgradability. The purpose of the contract is to have up-to-date logic
 * for random threshold number generation.
 */
contract KeepRandomBeaconImpl is Ownable, EternalStorage {

    StakingProxy public stakingProxy; // Staking proxy contract that is used to check stake balances against.

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 requestID, uint256 payment, uint256 blockReward, uint256 seed, uint blockNumber); 
    event RelayEntryGenerated(uint256 requestID, uint256 requestResponse, uint256 requestGroupID, uint256 previousEntry, uint blockNumber); 
    event RelayResetEvent(uint256 lastValidRelayEntry, uint256 lastValidRelayTxHash, uint256 lastValidRelayBlock);
    event SubmitGroupPublicKeyEvent(byte[] groupPublicKey, uint256 requestID, uint256 groupCount, uint256 activationBlockHeight);

    /**
     * @dev Creates Keep Random Beacon implementaion contract with a linked staking proxy contract.
     * @param _stakingProxy Address of a staking proxy contract that will be linked to this contract.
     * @param _minPayment Minimum amount of ether (in wei) that allows anyone to request a random number.
     * @param _minStake Minimum amount in KEEP that allows KEEP network client to participate in a group.
     */
    function KeepRandomBeaconImpl(address _stakingProxy, uint256 _minPayment, uint256 _minStake) public {
        require(_stakingProxy != address(0x0));
        stakingProxy = StakingProxy(_stakingProxy);
        uintStorage[keccak256("minStake")] = _minStake;
        uintStorage[keccak256("minPayment")] = _minPayment;
        uintStorage[keccak256("groupCountSequence")] = 0;
    }

    /// @dev Accept payments
    function () public payable {
    }

    /**
     * @dev Checks that the specified user has an appropriately large stake. 
     * @param _staker Specifies the identity of the random beacon client.
     * @return True if staked enough to participate in the group, false otherwise.
     */
    function isStaked(address _staker) public view returns(bool) {
        uint256 balance;
        balance = stakingProxy.balanceOf(_staker);
        return (balance >= uintStorage[keccak256("minStake")]);
    }

    /**
     * @dev Creates a request to generate a signature (random number)
     * @param _blockReward The value in keep for generating the signature.
     * @param _seed Initial seed random value from the client. It should be a cryptographically generated random value.
     * @return An uint256 representing uniquely generated ID. It is also returned as part of the event.
     */
    function requestRelay(uint256 _blockReward, uint256 _seed) public payable returns ( uint256 requestID ) {
        require(msg.value >= uintStorage[keccak256("minPayment")]); // Prevents payments that are too small in wei

        requestID = nextID();

        addressStorage[keccak256("requestPayer", requestID)] = msg.sender;
        uintStorage[keccak256("requestPayment", requestID)] = msg.value;
        uintStorage[keccak256("blockReward", requestID)] = _blockReward; // TODO - who decides the block reward? is it in KEEP?

        // Generate an event at this point, just return instead, RandomNumberRequest.
        RelayEntryRequested(requestID, msg.value, _blockReward, _seed, block.number);
    }

    /**
     * @dev Transfer 'msg.value' of funds directly from this contract to Keep multiwallet
     */
    function widthdrawAmount() public payable onlyOwner {
        owner.transfer(msg.value);
    }

    /**
     * @dev Set the minimum payment that is required before a relay entry occurs.
     * @param _minPayment is the value in wei that is required to be payed for the process to start.
     */
    function setMinimumPayment(uint256 _minPayment) public onlyOwner {
        uintStorage[keccak256("minPayment")] = _minPayment;
    }

    /**
     * @dev Takes the resulting signature and puts it onto the chain.
     * @param _requestID The request that started this generation - to tie the results back to the request.
     * @param _groupSignature The generated random number.
     * @param _groupID Public key of the group that generated the threshold signature.
     */
    function relayEntry(uint256 _requestID, uint256 _groupSignature, uint256 _groupID, uint256 _previousEntry) public {
        uintStorage[keccak256("requestGroupID", _requestID)] = _groupID;

        RelayEntryGenerated(_requestID, _groupSignature, _groupID, _previousEntry, block.number);
    }

    /**
     * @dev Makes an accusation that the relay entry has been falsified.
     * @param _lastValidRelayTxHash Last valid relay TX hash.
     * @param _lastValidRelayBlock Last valid relay block.
     */
    function relayEntryAccusation(uint256 _lastValidRelayTxHash, uint256 _lastValidRelayBlock) public {
        uint256 lastValidRelayEntry;
        lastValidRelayEntry = 1010101010;     // Some arbitrary number for testing.
        // TODO -- really need to understand what is needed at this point.
        // validate accusation by performing the checks in this code (slow/expensive)
        // raise event if accusation is shown to be true
        // penalty for false accusations - msg.sender? gets docked/rewarded?
        RelayResetEvent(lastValidRelayEntry, _lastValidRelayTxHash, _lastValidRelayBlock);
    }

    /**
     * @dev Takes a generated key and place it on the blockchain. Creates an event.
     * @param _groupPublicKey Group public key.
     * @param _requestID Request ID.
     */
    function submitGroupPublicKey(byte[] _groupPublicKey, uint256 _requestID) public {
        uint256 activationBlockHeight = block.number;
        uintStorage[keccak256("groupCountSequence")]++;

        // TODO -- lots of stuff - don't know yet.
        SubmitGroupPublicKeyEvent(_groupPublicKey, _requestID, uintStorage[keccak256("groupCountSequence")], activationBlockHeight);
    }

    /**
     * @dev Resets the group count to 0. Can only be called by the owner of the contract.
     */
    function resetGroupCount() public onlyOwner {
        uintStorage[keccak256("groupCountSequence")] = 0;
    }

    /**
     * @dev Generates a unique ID
     * @return An uint256 representing uniquely generated ID.
     */
    function nextID() private returns(uint256 requestID) {
        requestID = (block.timestamp ^ uint256(msg.sender)) + uintStorage[keccak256("seq")];
        uintStorage[keccak256("seq")]++;
        return (requestID);
    }
}
