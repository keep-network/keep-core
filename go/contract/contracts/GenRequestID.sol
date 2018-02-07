pragma solidity ^0.4.18;

/// @title Generate the RequestID 
/// @author Philip Schlump

contract GenRequestID { 
	// The current id
    uint256 public RequestIDSequence;

    // Constructor 
    function GenRequestID() public {
    	RequestIDSequence = 1;							// Start at 1, count up
    }

	/// @notice Get the next id and return it.
    function GenerateNextRequestID() public returns(uint256 RequestID) {
    	RequestIDSequence = RequestIDSequence + 1;
		RequestID = RequestIDSequence;
		return (RequestID);
	}

    /* This unnamed function is called whenever someone tries to send ether to it */
    function () public {
        revert(); // Prevents accidental sending of ether
    }        
} 
