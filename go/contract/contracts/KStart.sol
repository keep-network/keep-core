pragma solidity ^0.4.19;

contract KStart { 
    /* Public variables of this - used for testing */
    string public name;

    /* creates arrays with all relevant data */
    mapping (address => uint256) public requestID;
    mapping (address => uint256) public payment;
    mapping (address => uint256) public blockReward;
    mapping (address => uint256) public seed;
    mapping (address => uint) public complete;
    mapping (address => uint256) public number;

    /* This generates a public event on the blockchain that will notify clients */
    event RelayRequestReady(uint256 RequestID, uint256 payment, uint256 blockReward, uint256 seed);
    event RandomNumberReady(uint256 RequestID, bool Valid, uint256 RandomNumber);

    /* Constructor */
    function KStart(string testName) public {
        name = testName;                                   // Set the name for display purposes     
    }

	// Check that the person has a valid account (function name is wrong)
    function isStaked(uint256 /*_UserPublicKey*/) pure public returns(bool) {
		// TODO: check that the user with this public key has a stake that is adequate.
		return true;
	}

    /* Inputs: RequestID(user random), payment(eth), blockReward(keep), seed(number) */
    /* Outputs: RequestID ??? Should this be an output instead of an input? */
	// This would be better if the "RequestID" was a generated sequence number and returnd. -- How do I do that?
    function requestRandomNumber(uint256 _RequestID, uint256 _payment, uint256 _blockReward, uint256 _seed) public {
        requestID[msg.sender] = _RequestID ;
        payment[msg.sender] = _payment ;
        blockReward[msg.sender] = _blockReward ;
        seed[msg.sender] = _seed ;
		complete[msg.sender] = 0;

		// generate an event at this point, just return instead, RandomNumberRequest
    	RelayRequestReady( _RequestID, _payment, _blockReward, _seed);
    }

	// threshold relay has a number, or failed and passes _Valid as false
    function randomNumberComplete(uint256 _RequestID, bool _Valid, uint256 _theRandomNumber) public {
		if ( _Valid ) {
			complete[msg.sender] = 1;
			number[msg.sender] = _theRandomNumber;
		} else {
			complete[msg.sender] = 2;
		}
    	RandomNumberReady(_RequestID, _Valid, _theRandomNumber);
	}

	// Missing some sort of failed-validation, or group broke call?
	
	// Function that can be called to get your random number
	function getRandomNumber() view public returns ( uint status, uint256 theRandomNumber) {
		status = complete[msg.sender];
		theRandomNumber = number[msg.sender];
	}

    /* This unnamed function is called whenever someone tries to send ether to it */
    function () public {
        revert(); // Prevents accidental sending of ether
    }        
} 
