pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "../contracts/KeepRelayBeacon.sol";

contract TestKeepRelayBeacon {	

	address stakingAddress = address(0x51fcf39b2050c04a14152a5dd69884e811968128);
	KeepRelayBeacon ks = new KeepRelayBeacon( stakingAddress, 1 );

	function testRequestRelay() public {
		uint256 rid = ks.requestRelay(5,1212121);
		Assert.equal(rid, rid, "should have generated a unique request id, did not.");
	}

}
