pragma solidity ^0.4.18;

/// @title Interface contract for accessing random threshold number generation.
/// @author Philip Schlump

// https://ethereum.stackexchange.com/questions/20750/error-calling-a-function-from-another-contract-member-not-found-or-not-visi
// https://ethereum.stackexchange.com/questions/730/attaching-an-address-for-a-contract-to-call-another-contract
// https://github.com/pipermerriam/ethereum-uuid
// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-55.md
// https://github.com/ethereum/EIPs/issues/55

// https://blog.golemproject.net/how-to-find-10m-by-just-reading-blockchain-6ae9d39fcd95

// import 'GenRequestID.sol';

contract GenRequestID { 
    function GenerateNextRequestID() public returns(uint256 RequestID);
    function () public;
} 

contract KStart { 
    /* creates arrays with all relevant data */
    mapping (uint256 => uint256) public payment;
    mapping (uint256 => uint256) public blockReward;
    mapping (uint256 => uint256) public seed;
    mapping (uint256 => uint256) public number;

	GenRequestID public GenRequestIDSequence;

    /* This generates a public event on the blockchain that will notify clients */
    event RelayRequestReady(uint256 RequestID, uint256 payment, uint256 blockReward, uint256 seed);
    event RandomNumberReady(uint256 RequestID, uint256 RandomNumber);

    /* Constructor */
    function KStart() public {
		// GenRequestIDSequence = GenRequestID(GEN_REQUEST_ID_ADDR);
		// GenRequestIDSequence = GenRequestID(0x9fbda871d559710256a2502a2517b794b482db40);
		GenRequestIDSequence = GenRequestID(0x9FBDa871d559710256a2502A2517b794B482Db40);
    }

	/// @notice checks that the specified user has an appropriately large stake.   Returns true if staked.
	/// @param _UserPublicKey specifies the user.
	/// @dev check in the staking registry for this user.   Must be able to look up user based on Pub-Key.
	///   Q: If user is trying to run 3 nodes but only has a stake for 2 nodes how is this determined?
	///   Q: If how is the amount of :keep: specified per-node - where?
	///   TODO: Find the contract that we are supposed to call - call it.
	/// For the moment just return true so we can test with this.
    function isStaked(uint256 _UserPublicKey) pure public returns(bool) {
		_UserPublicKey = _UserPublicKey;	// to make it used
		return true;
	}

    // Inputs: payment(eth), blockReward(keep), seed(number) 
	// The "RequestID" is generated sequence number and returned.
	// 	RequestID is definitely an output 
	// 	RequestID Monotonically Increasing 
    function requestRandomNumber(uint256 _payment, uint256 _blockReward, uint256 _seed) public returns ( uint256 RequestID ) {
		RequestID = GenRequestIDSequence.GenerateNextRequestID();

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
