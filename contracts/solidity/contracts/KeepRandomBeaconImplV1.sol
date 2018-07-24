pragma solidity ^0.4.24;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";
import "./StakingProxy.sol";
import "./EternalStorage.sol";


/**
 * @title KeepRandomBeaconImplV1
 * @dev Initial version of implementation contract that works under Keep Random
 * Beacon proxy and allows upgradability. The purpose of the contract is to have
 * up-to-date logic for threshold random number generation. Updated contracts
 * must inherit from this contract and have to be initialized under updated version name
 */
contract KeepRandomBeaconImplV1 is Ownable, EternalStorage {

    StakingProxy public stakingProxy; // Staking proxy contract that is used to check stake balances against.

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 requestID, uint256 payment, uint256 blockReward, uint256 seed, uint blockNumber); 
    event RelayEntryGenerated(uint256 requestID, uint256 requestResponse, uint256 requestGroupID, uint256 previousEntry, uint blockNumber);
    event RelayResetEvent(uint256 lastValidRelayEntry, uint256 lastValidRelayTxHash, uint256 lastValidRelayBlock);
    event SubmitGroupPublicKeyEvent(byte[] groupPublicKey, uint256 requestID, uint256 activationBlockHeight);

    bytes32 private constant esSeq = keccak256("seq");
    bytes32 private constant esMinPayment = keccak256("minPayment");
    bytes32 private constant esMinStake = keccak256("minStake");
    bytes32 private constant esRequestPayer = keccak256("requestPayer");
    bytes32 private constant esRequestPayment = keccak256("requestPayment");
    bytes32 private constant esBlockReward = keccak256("blockReward");
    bytes32 private constant esRequestGroupID = keccak256("requestGroupID");
    bytes32 private constant esGroupThreshold = keccak256("groupThreshold");
    bytes32 private constant esStakingProxy = keccak256("stakingProxy");
    bytes32 private constant esKeepRandomBeaconImplV1 = keccak256("KeepRandomBeaconImplV1");
    bytes32 private constant esWithdrawalDelay = keccak256("withdrawalDelay");
    bytes32 private constant esPendingWithdrawal = keccak256("pendingWithdrawal");

    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() public payable {
        revert();
    }

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
        require(!initialized(), "Contract has already been initialized.");
        require(_stakingProxy != address(0x0), "Invalid 0 address for StakingProxy passed.");
        addressStorage[esStakingProxy] = _stakingProxy;
        uintStorage[esMinStake] = _minStake;
        uintStorage[esMinPayment] = _minPayment;
        boolStorage[esKeepRandomBeaconImplV1] = true;
        uintStorage[esWithdrawalDelay] = _withdrawalDelay;
        uintStorage[esPendingWithdrawal] = 0;
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return boolStorage[esKeepRandomBeaconImplV1];
    }

    /**
     * @dev Checks that the specified user has an appropriately large stake.
     * @param _staker Specifies the identity of the random beacon client.
     * @return True if staked enough to participate in the group, false otherwise.
     */
    function hasMinimumStake(address _staker) public view returns(bool) {
        uint256 balance;
        stakingProxy = StakingProxy(addressStorage[esStakingProxy]);
        balance = stakingProxy.balanceOf(_staker);
        return (balance >= uintStorage[esMinStake]);
    }

    /**
     * @dev Creates a request to generate a new relay entry, which will include a
     * random number (by signing the previous entry's random number).
     * @param _blockReward The value in KEEP for generating the signature.
     * @param _seed Initial seed random value from the client. It should be a cryptographically generated random value.
     * @return An uint256 representing uniquely generated ID. It is also returned as part of the event.
     */
    function requestRelayEntry(uint256 _blockReward, uint256 _seed) public payable returns (uint256 requestID) {
        require(msg.value >= uintStorage[esMinPayment], "Payment too small."); // Prevents payments that are too small in wei

        requestID = uintStorage[esSeq]++;
        addressStorageMap[esRequestPayer][requestID] = msg.sender;
        uintStorageMap[esRequestPayment][requestID] = msg.value;
        uintStorageMap[esBlockReward][requestID] = _blockReward;

        emit RelayEntryRequested(requestID, msg.value, _blockReward, _seed, block.number);
        return requestID;
    }

    /**
     * @dev Return the current RequestID - used in testing.
     */
    function getRequestId() public view returns (uint256 requestID) {
        requestID = uintStorage[esSeq];
        return requestID;
	}

    /**
     * @dev Initiate withdrawal of this contract balance to the owner.
     */
    function initiateWithdrawal() public onlyOwner {
        uint256 withdrawalDelay = uintStorage[esWithdrawalDelay];
        uintStorage[esPendingWithdrawal] = block.timestamp + withdrawalDelay;
    }

    /**
     * @dev Finish withdrawal of this contract balance to the owner.
     */
    function finishWithdrawal() public onlyOwner {
        uint pendingWithdrawal = uintStorage[esPendingWithdrawal];

        require(pendingWithdrawal > 0, "Pending Withdrawl must be larger than 0.");
        require(block.timestamp >= pendingWithdrawal);

        // Reset pending withdrawal before sending to prevent re-entrancy attacks
        uintStorage[esPendingWithdrawal] = 0;
        owner.transfer(address(this).balance);
    }

    /**
     * @dev Set the minimum payment that is required before a relay entry occurs.
     * @param _minPayment is the value in wei that is required to be payed for the process to start.
     */
    function setMinimumPayment(uint256 _minPayment) public onlyOwner {
        uintStorage[esMinPayment] = _minPayment;
    }

    /**
     * @dev Set the minimum amount of KEEP that allows a Keep network client to participate in a group.
     * @param _minStake Amount in KEEP.
     */
    function setMinimumStake(uint256 _minStake) public onlyOwner {
        uintStorage[esMinStake] = _minStake;
    }

    /**
     * @dev Get the minimum payment that is required before a relay entry occurs.
     */
    function minimumPayment() public view returns(uint256) {
        return uintStorage[esMinPayment];
    }

    /**
     * @dev Get the minimum amount in KEEP that allows KEEP network client to participate in a group.
     */
    function minimumStake() public view returns(uint256) {
        return uintStorage[esMinStake];
    }

    /**
     * @dev Gets the threshold size for groups.
     */
    function groupThreshold() public view returns(uint256) {
        return uintStorage[esGroupThreshold];
    }

    /**
     * @dev Creates a new relay entry and stores the associated data on the chain.
     * @param _requestID The request that started this generation - to tie the results back to the request.
     * @param _groupSignature The generated random number.
     * @param _groupID Public key of the group that generated the threshold signature.
     */
    function relayEntry(uint256 _requestID, uint256 _groupSignature, uint256 _groupID, uint256 _previousEntry) public {
        uintStorageMap[esRequestGroupID][_requestID] = _groupID;

        emit RelayEntryGenerated(_requestID, _groupSignature, _groupID, _previousEntry, block.number);
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
}
