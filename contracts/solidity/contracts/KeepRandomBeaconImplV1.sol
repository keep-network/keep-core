pragma solidity ^0.4.21;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";
import "./StakingProxy.sol";
import "./EternalStorage.sol";


/**
 * @title KeepRandomBeaconImplV1
 * @dev Initial version of implementation contract that works under Keep Random
 * Beacon proxy and allows upgradability. The purpose of the contract is to have
 * up-to-date logic for random threshold number generation. Updated contracts
 * must inherit from this contract and have to be initialized under updated version name
 */
contract KeepRandomBeaconImplV1 is Ownable, EternalStorage {

    StakingProxy public stakingProxy; // Staking proxy contract that is used to check stake balances against.

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 requestID, uint256 payment, uint256 blockReward, uint256 seed, uint blockNumber); 
    event RelayEntryGenerated(uint256 requestID, uint256 requestResponse, uint256 requestGroupID, uint256 previousEntry, uint blockNumber);
    event RelayResetEvent(uint256 lastValidRelayEntry, uint256 lastValidRelayTxHash, uint256 lastValidRelayBlock);
    event SubmitGroupPublicKeyEvent(byte[] groupPublicKey, uint256 requestID, uint256 activationBlockHeight);

    /**
     * @dev Initialize Keep Random Beacon implementaion contract with a linked staking proxy contract.
     * @param _stakingProxy Address of a staking proxy contract that will be linked to this contract.
     * @param _minPayment Minimum amount of ether (in wei) that allows anyone to request a random number.
     * @param _minStake Minimum amount in KEEP that allows KEEP network client to participate in a group.
     * @param _withdrawalDelay Delay before the owner can withdraw ether from this contract.
     */
    function initialize(address _stakingProxy, uint256 _minPayment, uint256 _minStake, uint256 _withdrawalDelay)
        public
        onlyOwner
    {
        require(!initialized());
        require(_stakingProxy != address(0x0));
        addressStorage[keccak256("stakingProxy")] = _stakingProxy;
        uintStorage[keccak256("minStake")] = _minStake;
        uintStorage[keccak256("minPayment")] = _minPayment;
        boolStorage[keccak256("KeepRandomBeaconImplV1")] = true;
        uintStorage[keccak256("withdrawalDelay")] = _withdrawalDelay;
        uintStorage[keccak256("pendingWithdrawal")] = 0;
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return boolStorage[keccak256("KeepRandomBeaconImplV1")];
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
        stakingProxy = StakingProxy(addressStorage[keccak256("stakingProxy")]);
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
        emit RelayEntryRequested(requestID, msg.value, _blockReward, _seed, block.number);
    }

    /**
     * @dev Initiate withdraw of this contract balance to the owner.
     */
    function initiateWithdraw() public onlyOwner {
        uint256 withdrawalDelay = uintStorage[keccak256("withdrawalDelay")];
        uintStorage[keccak256("pendingWithdrawal")] = block.timestamp + withdrawalDelay;
    }

    /**
     * @dev Finish withdraw of this contract balance to the owner.
     */
    function finishWithdraw() public onlyOwner {
        uint pendingWithdrawal = uintStorage[keccak256("pendingWithdrawal")];

        require(pendingWithdrawal > 0);
        require(block.timestamp >= pendingWithdrawal);

        // Reset pending withdrawal before sending to prevent re-entrancy attacks
        uintStorage[keccak256("pendingWithdrawal")] = 0;
        owner.transfer(this.balance);
    }

    /**
     * @dev Set the minimum payment that is required before a relay entry occurs.
     * @param _minPayment is the value in wei that is required to be payed for the process to start.
     */
    function setMinimumPayment(uint256 _minPayment) public onlyOwner {
        uintStorage[keccak256("minPayment")] = _minPayment;
    }

    /**
     * @dev Set the minimum amount in KEEP that allows KEEP network client to participate in a group.
     * @param _minStake Amount in KEEP.
     */
    function setMinimumStake(uint256 _minStake) public onlyOwner {
        uintStorage[keccak256("minStake")] = _minStake;
    }

    /**
     * @dev Get the minimum payment that is required before a relay entry occurs.
     */
    function minimumPayment() public view returns(uint256) {
        return uintStorage[keccak256("minPayment")];
    }

    /**
     * @dev Get the minimum amount in KEEP that allows KEEP network client to participate in a group.
     */
    function minimumStake() public view returns(uint256) {
        return uintStorage[keccak256("minStake")];
    }

    /**
     * @dev Takes the resulting signature and puts it onto the chain.
     * @param _requestID The request that started this generation - to tie the results back to the request.
     * @param _groupSignature The generated random number.
     * @param _groupID Public key of the group that generated the threshold signature.
     */
    function relayEntry(uint256 _requestID, uint256 _groupSignature, uint256 _groupID, uint256 _previousEntry) public {
        uintStorage[keccak256("requestGroupID", _requestID)] = _groupID;

        emit RelayEntryGenerated(_requestID, _groupSignature, _groupID, _previousEntry, block.number);
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
        emit RelayResetEvent(lastValidRelayEntry, _lastValidRelayTxHash, _lastValidRelayBlock);
    }

    /**
     * @dev Takes a generated key and place it on the blockchain. Creates an event.
     * @param _groupPublicKey Group public key.
     * @param _requestID Request ID.
     */
    function submitGroupPublicKey(byte[] _groupPublicKey, uint256 _requestID) public {
        uint256 activationBlockHeight = block.number;

        // TODO -- lots of stuff - don't know yet.
        emit SubmitGroupPublicKeyEvent(_groupPublicKey, _requestID, activationBlockHeight);
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
