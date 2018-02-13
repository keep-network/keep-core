pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "../contracts/KStart.sol";

contract TestKStart {	
	
	KStart ks = new KStart();

	function testRequestRelay() public {
		uint256 rid = ks.requestRelay(12,12,12);
		rid = rid;
		// Assert.equal(id0, id1, "should have generated a unique id, did not.");
	}

}
