pragma solidity ^0.4.18;

/// @title Interface contract for accessing random threshold number generation.
/// @author Philip Schlump

import './TokenStaking.sol';

contract KeepRelayBeacon { 

	address public contractOwner = msg.sender;
	uint256 minPayment; 		// value in wei
	uint256 minStakingBalance;	// Minimum amount in KEEP that is allowed for a client to participate in a group
	uint256 public GroupCountSequence;
	uint256 seq = 1;
	TokenStaking public staking;

    mapping (uint256 => address) public payment_from;		// Payment from
    mapping (uint256 => uint256) public payment_amount;		// Payment amount to generate *signature*
    mapping (uint256 => uint256) public blockReward;		// 
    mapping (uint256 => uint256) public seed;				// Input Seed
    mapping (uint256 => uint256) public signature;			// The randomly generated number
    mapping (uint256 => uint256) public groupID;			// What gorup generated the signatre

    // These are the public events that are used by clients
    event RelayEntryRequested(uint256 RequestID, uint256 Payment, uint256 BlockReward, uint256 Seed, uint blockNumber); 
    event RelayEntryGenerated(uint256 RequestID, uint256 Signature, uint256 GroupID, uint256 PreviousEntry, uint blockNumber); 
    event RelayResetEvent(uint256 LastValidRelayEntry, uint256 LastValidRelayTxHash, uint256 LastValidRelayBlock);
    event SubmitGroupPublicKeyEvent(byte[] _PK_G_i, uint256 RequestID, uint256 GroupCount, uint256 ActivationBlockHeight);

    // Constructor 
    function KeepRelayBeacon(address _stakingAddress, uint256 _minKeep) public {
		require(_stakingAddress != address(0x0));
        staking = TokenStaking(_stakingAddress);
		minPayment = 1;
		GroupCountSequence = 0;
		minStakingBalance = _minKeep;	// Minimum amount in KEEP that is allowed for a client to participate in a group
    }

	/// @notice get the next RequestID 
    function nextID() private returns(uint256 RequestID) {
		RequestID = ( block.timestamp ^ uint256(msg.sender) ) + seq;
		seq = seq + 1;
		return ( RequestID );
	}

	/// @notice checks that the specified user has an appropriately large stake.   Returns true if staked.
	/// @param _staker specifies the identity of the random beacon client.
    function isStaked(address _staker) view public returns(bool) {
		uint256 balance;
		balance = staking.balanceOf(_staker);
		return ( balance >= minStakingBalance );
	}

	/// @notice make a request to generate a signature (random number)
    /// @param _blockReward is the value in keep??? for generating the signature.
    /// @param _seed is an initial seed random value from the client.  It should be a cryptographically generated random value.
	/// @dev The "RequestID" is generated unique ID. It is returned and part of the event.
    function requestRelay(uint256 _blockReward, uint256 _seed) public payable returns ( uint256 RequestID ) {
		require( msg.value >= minPayment ); // Prevents payments that are too small in wei

		RequestID = nextID();

		payment_from[RequestID] = msg.sender;
		payment_amount[RequestID] = msg.value;

        blockReward[RequestID] = _blockReward ;		// TODO - who decides the block reward?  is it in KEEP?
        seed[RequestID] = _seed ;					// TODO - is it a security risk to save the "seed" as a public value?

		// generate an event at this point, just return instead, RandomNumberRequest
     	RelayEntryRequested( RequestID, msg.value, _blockReward, _seed, block.number);
    }

	/// @notice validates that a call to the function will only succeed if the owner of the contract made the call.
	modifier onlyOwner() {
		require(msg.sender == contractOwner);
		_;
	}

	// @notice Transfer 'msg.value' of funds directly from this contract to Keep multiwallet
	function widthdrawAmount() onlyOwner public payable {
		contractOwner.transfer(msg.value);
	}

	/// @notice Set the minimum payment that is required before a relay entry occurs.
	/// @param _minPayment is the value in wei that is required to be payed for the process to start.
	function setMinimumPayment( uint256 _minPayment ) onlyOwner public {
		minPayment = _minPayment;
	}

	/// @notice take the resulting signature and put it onto the chain.
	/// @param _RequestID the request that started this generation - to tie the results back to the request.
	/// @param _groupSignature is the generated random number
	/// @param _groupID is the public key of the gorup that generated the threshold signature
    function relayEntry(uint256 _RequestID, uint256 _groupSignature, uint256 _groupID, uint256 _previousEntry) public {
		signature[_RequestID] = _groupSignature;
		groupID[_RequestID] = _groupID;

     	RelayEntryGenerated(_RequestID, _groupSignature, _groupID, _previousEntry, block.number);
	}

	/// @notice make an accusation that the relay entry has been falsified. 
    function relayEntryAccusation( uint256 _LastValidRelayTxHash, uint256 _LastValidRelayBlock) public {
		uint256 LastValidRelayEntry;
		LastValidRelayEntry = 1010101010;	 // Some arbitrary number for testing.
		// TODO -- really need to understand what is needed at this point.
		// validate accusation by performing the checks in this code (slow/expensive)
		// raise event if accusation is shown to be true
		// penalty for false accusations - msg.sender? gets docked/rewarded?
		RelayResetEvent(LastValidRelayEntry, _LastValidRelayTxHash, _LastValidRelayBlock);	
	}

	/// @notice hmm...
    function submitGroupPublicKey (byte[] _PK_G_i, uint256 _RequestID) public {
		uint256 ActivationBlockHeight = block.number;
		// uint256 public GroupCountSequence;
		GroupCountSequence = GroupCountSequence + 1;

		// TODO -- lots of stuff - don't know yet.
		SubmitGroupPublicKeyEvent(_PK_G_i, _RequestID, GroupCountSequence, ActivationBlockHeight);
	}

	/// @dev resets the group count to 0.  Can only be called by the owner of the contract.
	function resetGroupCount() onlyOwner public {
		GroupCountSequence = 0;
	}
	
	function () public payable {
	}
} 

