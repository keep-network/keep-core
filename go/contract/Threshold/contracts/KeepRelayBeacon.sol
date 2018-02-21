pragma solidity ^0.4.18;



/// @title Interface contract for accessing random threshold number generation.
/// @author Philip Schlump


// Contract Interface for GenRequestID contract from ./GenRequestID.sol
contract GenRequestID { 
    function GenerateNextRequestID() public returns(uint256 RequestID);
    function () public;
} 

contract KeepRelayBeacon { 

    /* creates arrays with all relevant data */
    mapping (uint256 => uint256) public payment;
    mapping (uint256 => uint256) public blockReward;
    mapping (uint256 => uint256) public seed;			// Input Seed
    mapping (uint256 => uint256) public signature;		// The randomly generated number
    mapping (uint256 => uint256) public groupID;		// What gorup generated the signatre

	GenRequestID public GenRequestIDSequence;

    /* This generates a public event on the blockchain that will notify clients */
    event RelayEntryRequested(uint256 RequestID, uint256 Payment, uint256 BlockReward, uint256 Seed); 
    event RelayEntryGenerated(uint256 RequestID, uint256 Signature, uint256 GroupID, uint256 PreviousEntry ); // xyzzy - RelayEntryGenerated.
    event RelayResetEvent(uint256 LastValidRelayEntry, uint256 LastValidRelayTxHash, uint256 LastValidRelayBlock);	// xyzzy - data types on TxHash, Block
    event SubmitGroupPublicKeyEvent(uint256 _PK_G_i, uint256 _id, uint256 _activationBlockHeight);

    /* Constructor */
    function KStart() public {
		GenRequestIDSequence = GenRequestID(0x1CEdd10a7D1CBeC5D807B81aB24e542cc0F6BE31);
    }

	// get the next id from the generator contract
    function nextID() private returns(uint256 RequestID) {
		RequestID = GenRequestIDSequence.GenerateNextRequestID();
		return ( RequestID );
	}

	/// @notice checks that the specified user has an appropriately large stake.   Returns true if staked.
	/// @param _UserPublicKey specifies the user.
	/// @dev check in the staking registry for this user.   Must be able to look up user based on Pub-Key.
	///   Q: If user is trying to run 3 nodes but only has a stake for 2 nodes how is this determined?
	///   Q: If how is the amount of :keep: specified per-node - where?
	///   TODO: Find the contract that we are supposed to call - call it.
	/// For the moment just return true so we can test with this.
    function isStaked(uint256 _UserPublicKey) pure public returns(bool) {
		_UserPublicKey = _UserPublicKey;	// to make it used so Solidity will not complain - temporary
		return true;
	}

    // Inputs: payment(eth), blockReward(keep), seed(number) 
	// The "RequestID" is generated sequence number and returned.
	// 	RequestID is definitely an output 
	// 	RequestID Monotonically Increasing 
	// This will down-streem from event result in a SignatureShareBroadcast on the KEEP p2p network.
    function requestRelay(uint256 _payment, uint256 _blockReward, uint256 _seed) public returns ( uint256 RequestID ) {
		RequestID = nextID();

        payment[RequestID] = _payment ;				// TODO - validation on these values?
        blockReward[RequestID] = _blockReward ;		// TODO - validation on these values?
        seed[RequestID] = _seed ;

		// generate an event at this point, just return instead, RandomNumberRequest
     	RelayEntryRequested( RequestID, _payment, _blockReward, _seed);
    }

	/// @dev Must include "group" that is responable:  ATL
	/// 		threshold relay has a number, or failed and passes _Valid as false
	///
	/// @param _RequestID the request that started this generation - to tie the results back to the request.
	/// @param _groupSignature is the generated random number
	/// @param _groupID is the public key of the gorup that generated the threshold signature
    function relayEntry(uint256 _RequestID, uint256 _groupSignature, uint256 _groupID, uint256 _previousEntry) public {
		signature[_RequestID] = _groupSignature;
		groupID[_RequestID] = _groupID;

     	RelayEntryGenerated(_RequestID, _groupSignature, _groupID, _previousEntry);
	}

    function relayEntryAccusation( uint256 _LastValidRelayTxHash, uint256 _LastValidRelayBlock) public {
		uint256 LastValidRelayEntry;
		// xyzzy  / TODO
		// validate acusation by performaing the checks in this code (slow/expensive)
		// raise event if acusation is shown to be true
		// penalty for false acusations - msg.sender? gets docked/rewareded?
		if ( 0 == 1 ) {
			RelayResetEvent(LastValidRelayEntry, _LastValidRelayTxHash, _LastValidRelayBlock);	
		}
	}

    function submitGroupPublicKey (uint256 _PK_G_i, uint256 _id) public {
		uint256 ActivationBlockHeight = block.number;
		// xyzzy  / TODO
		if ( 0 == 1 ) {
			SubmitGroupPublicKeyEvent(_PK_G_i, _id, ActivationBlockHeight);
		}
	}
	
	// Function that can be called to get your random number
	function getRandomNumber( uint256 _RequestID ) view public returns ( uint256 theRandomNumber ) {
		theRandomNumber = signature[_RequestID];
	}

    // This unnamed function is called whenever someone tries to send ether to it.
    function () public {
        revert(); // Prevents accidental sending of ether
    }        
} 
