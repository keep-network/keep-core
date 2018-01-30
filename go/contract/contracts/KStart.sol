pragma solidity ^0.4.19;

// From: https://gist.github.com/karalabe/08f4b780e01c8452d989
// From: https://github.com/ethereum/go-ethereum/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts
// Modified from token.sol

contract KStart { 
    /* Public variables of the token */
    string public name;

    /* This creates an array with all balances */
    mapping (address => uint256) public requestID;
    mapping (address => uint256) public payment;
    mapping (address => uint256) public blockReward;
    mapping (address => uint256) public seed;

    /* This generates a public event on the blockchain that will notify clients */
    event RelayRequestReady(uint256 RequestID, uint256 payment, uint256 blockReward, uint256 seed);

    /* Constructor */
    function KStart(string testName) public {
        // balanceOf[msg.sender] = initialSupply;              // Give the creator all initial tokens                    
        name = testName;                                   // Set the name for display purposes     
    }

    /* Inputs: RequestID(user random), payment(eth), blockReward(keep), seed(number) */
    /* Outputs: RequestID */
	// This would be better if the "RequestID" was a generated sequence number and returnd. -- How do I do that?
    function requestRandomNumber(uint256 _RequestID, uint256 _payment, uint256 _blockReward, uint256 _seed) public {
        requestID[msg.sender] = _RequestID ;
        payment[msg.sender] = _payment ;
        blockReward[msg.sender] = _blockReward ;
        seed[msg.sender] = _seed ;

		// generate an event at this point, just return instead, RandomNumberRequest
    	RelayRequestReady( _RequestID, _payment, _blockReward, _seed);
    }

    /* This unnamed function is called whenever someone tries to send ether to it */
    function () public {
        revert(); // Prevents accidental sending of ether
    }        
} 
