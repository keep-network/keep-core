pragma solidity ^0.4.18;

/// @title Generate the RequestID 
/// @author Philip Schlump
/// @notice this contract should have a fixed address and never get reloaded.
//    If it is reloaded then the constant in the constructor should be updated to a number
//	  larger than the largest RequestID that has ever been generated.

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

    /// @notice This unnamed function is called whenever someone tries to send ether to it.
    function () public {
        revert(); // Prevents accidental sending of ether
    }        
} 
