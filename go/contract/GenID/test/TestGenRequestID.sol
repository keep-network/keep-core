pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "../contracts/GenRequestID.sol";

contract TestGenRequestID {	
	
	GenRequestID idSeq = new GenRequestID();

	function testAnID() public {
		uint256 id0 = idSeq.GenerateNextRequestID();
		uint256 id1 = idSeq.GenerateNextRequestID();
		id1 = id1 - 1;
		Assert.equal(id0, id1, "should have generated a unique id, did not.");
	}

}
