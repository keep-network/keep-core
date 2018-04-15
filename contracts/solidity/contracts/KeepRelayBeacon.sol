pragma solidity ^0.4.18;

/// @title Interface contract for accessing random threshold number generation.
/// @author Philip Schlump

import "zeppelin-solidity/contracts/ownership/Ownable.sol";
import './TokenStaking.sol';

contract KeepRelayBeacon is Ownable { 

    uint256 minPayment;         // value in wei
    uint256 minStake;    // Minimum amount in KEEP that is allowed for a client to participate in a group
    uint256 public groupCountSequence;
    uint256 seq = 1;
    TokenStaking public staking;

    mapping (uint256 => address) public requestPayer;        // Payment from
    mapping (uint256 => uint256) public requestPayment;        // Payment amount to generate *signature*
    mapping (uint256 => uint256) public blockReward;        // 
    mapping (uint256 => uint256) public seed;                // Input Seed
    mapping (uint256 => uint256) public requestResponse;            // The randomly generated number
    mapping (uint256 => uint256) public requestGroupID;            // What gorup generated the signatre

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 requestID, uint256 payment, uint256 blockReward, uint256 seed, uint blockNumber); 
    event RelayEntryGenerated(uint256 requestID, uint256 requestResponse, uint256 requestGroupID, uint256 previousEntry, uint blockNumber); 
    event RelayResetEvent(uint256 lastValidRelayEntry, uint256 lastValidRelayTxHash, uint256 lastValidRelayBlock);
    event SubmitGroupPublicKeyEvent(byte[] groupPublicKey, uint256 requestID, uint256 groupCount, uint256 activationBlockHeight);

    // Constructor 
    function KeepRelayBeacon(address _stakingAddress, uint256 _minKeep) public {
        require(_stakingAddress != address(0x0));
        staking = TokenStaking(_stakingAddress);
        minpayment = 1;
        groupCountSequence = 0;
        minStake = _minKeep;    // Minimum amount in KEEP that is allowed for a client to participate in a group
    }

    /// @dev Accept payments
    function () public payable {
    }

    /// @dev checks that the specified user has an appropriately large stake.   Returns true if staked.
    /// @param _staker specifies the identity of the random beacon client.
    function isStaked(address _staker) public view returns(bool) {
        uint256 balance;
        balance = staking.balanceOf(_staker);
        return (balance >= minStake);
    }

    /// @dev make a request to generate a signature (random number)
    /// @param _blockReward is the value in keep??? for generating the signature.
    /// @param _seed is an initial seed random value from the client.  It should be a cryptographically generated random value.
    /// @dev The "RequestID" is generated unique ID. It is returned and part of the event.
    function requestRelay(uint256 _blockReward, uint256 _seed) public payable returns ( uint256 requestID ) {
        require(msg.value >= minPayment); // Prevents payments that are too small in wei

        requestID = nextID();

        requestPayer[requestID] = msg.sender;
        requestPayment[requestID] = msg.value;

        blockReward[requestID] = _blockReward;        // TODO - who decides the block reward?  is it in KEEP?
        seed[requestID] = _seed;                    // TODO - is it a security risk to save the "seed" as a public value?

        // generate an event at this point, just return instead, RandomNumberRequest
        RelayEntryRequested(requestID, msg.value, _blockReward, _seed, block.number);
    }

    // @dev Transfer 'msg.value' of funds directly from this contract to Keep multiwallet
    function widthdrawAmount() public payable onlyOwner {
        owner.transfer(msg.value);
    }

    /// @dev Set the minimum payment that is required before a relay entry occurs.
    /// @param _minPayment is the value in wei that is required to be payed for the process to start.
    function setMinimumPayment( uint256 _minPayment ) public onlyOwner {
        minPayment = _minPayment;
    }

    /// @dev take the resulting signature and put it onto the chain.
    /// @param _requestID the request that started this generation - to tie the results back to the request.
    /// @param _groupSignature is the generated random number
    /// @param _groupID is the public key of the gorup that generated the threshold signature
    function relayEntry(uint256 _requestID, uint256 _groupSignature, uint256 _groupID, uint256 _previousEntry) public {
        requestResponse[_requestID] = _groupSignature;
        requestGroupID[_requestID] = _groupID;

        RelayEntryGenerated(_requestID, _groupSignature, _groupID, _previousEntry, block.number);
    }

    /// @dev make an accusation that the relay entry has been falsified. 
    function relayEntryAccusation( uint256 _lastValidRelayTxHash, uint256 _lastValidRelayBlock) public {
        uint256 lastValidRelayEntry;
        lastValidRelayEntry = 1010101010;     // Some arbitrary number for testing.
        // TODO -- really need to understand what is needed at this point.
        // validate accusation by performing the checks in this code (slow/expensive)
        // raise event if accusation is shown to be true
        // penalty for false accusations - msg.sender? gets docked/rewarded?
        RelayResetEvent(lastValidRelayEntry, _lastValidRelayTxHash, _lastValidRelayBlock);    
    }

    /// @dev take a generated key and place it on the blockchain.  Create an event.
    function submitGroupPublicKey (byte[] groupPublicKey, uint256 _requestID) public {
        uint256 activationBlockHeight = block.number;
        // uint256 public groupCountSequence;
        groupCountSequence = groupCountSequence + 1;

        // TODO -- lots of stuff - don't know yet.
        SubmitGroupPublicKeyEvent(groupPublicKey, _requestID, groupCountSequence, activationBlockHeight);
    }

    /// @dev resets the group count to 0.  Can only be called by the owner of the contract.
    function resetGroupCount() public onlyOwner {
        groupCountSequence = 0;
    }

    /// @dev get the next RequestID 
    function nextID() private returns(uint256 requestID) {
        requestID = (block.timestamp ^ uint256(msg.sender)) + seq;
        seq = seq + 1;
        return (requestID);
    }
}
