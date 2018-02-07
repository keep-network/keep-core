pragma solidity ^0.4.18;

/// @title Interface contract for accessing random threshold number generation.
/// @author Philip Schlump

contract KStart { 
    /* Public variables of this - used for testing */
    string public name;

    /* creates arrays with all relevant data */
    mapping (uint256 => uint256) public payment;
    mapping (uint256 => uint256) public blockReward;
    mapping (uint256 => uint256) public seed;
    mapping (uint256 => uint256) public number;

    uint256 public RequestIDSequence;

    /* This generates a public event on the blockchain that will notify clients */
    event RelayRequestReady(uint256 RequestID, uint256 payment, uint256 blockReward, uint256 seed);
    event RandomNumberReady(uint256 RequestID, bool Valid, uint256 RandomNumber);

    /* Constructor */
    function KStart(string testName) public {
        name = testName;                                // Set the name for display purposes     
    	RequestIDSequence = 1;							// Start at 1, count up
    }

	/// @notice checks that the specified user has an appropriately large stake.   Returns true if staked.
	/// @param _UserPublicKey specifies the user.
	/// @dev check in the staking registry for this user.   Must be able to look up user based on Pub-Key.
	///   Q: If user is trying to run 3 nodes but only has a stake for 2 nodes how is this determined?
	///   Q: If how is the amount of :keep: specified per-node - where?
	///   TODO: Find the contract that we are supposed to call - call it.
	/// For the moment just return true so we can test with this.
    function isStaked(uint256 /*_UserPublicKey*/) pure public returns(bool) {
		return true;
	}

    // Inputs: payment(eth), blockReward(keep), seed(number) 
	// The "RequestID" is generated sequence number and returned.
	// 	RequestID is definitely an output 
	// 	RequestID Monotonically Increasing 
    function requestRandomNumber(uint256 _payment, uint256 _blockReward, uint256 _seed) public returns ( uint256 RequestID ) {
    	RequestIDSequence = RequestIDSequence + 1;
		RequestID = RequestIDSequence;

        requestBy[RequestID] = RequestID;
        payment[RequestID] = _payment ;
        blockReward[RequestID] = _blockReward ;
        seed[RequestID] = _seed ;

		// generate an event at this point, just return instead, RandomNumberRequest
    	RelayRequestReady( RequestID, _payment, _blockReward, _seed);
    }

	// Must include "group" that is responable:  ATL
	// threshold relay has a number, or failed and passes _Valid as false
    function randomNumberComplete(uint256 _RequestID, uint256 _theRandomNumber) public {
		number[_RequestID] = _theRandomNumber;
    	RandomNumberReady(_RequestID, _theRandomNumber);
	}

	// Missing some sort of failed-validation, or group broke call?
	
	// Function that can be called to get your random number
	function getRandomNumber( uint256 _RequestID ) view public returns ( uint256 theRandomNumber ) {
		theRandomNumber = number[_RequestID];
	}

    /* This unnamed function is called whenever someone tries to send ether to it */
    function () public {
        revert(); // Prevents accidental sending of ether
    }        
} 
