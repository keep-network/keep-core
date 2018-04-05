pragma solidity ^0.4.18;

import "../contracts/KeepGroup.sol";
import "../contracts/KeepRelayBeacon.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenStaking.sol";

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "./helpers/ThrowProxy.sol";

contract TestKeepGroup {	
	
	// Create KEEP token
	KeepToken token = new KeepToken();

	// Create staking contract with 30 days withdrawal delay
	TokenStaking stakingContract = new TokenStaking(token, 30 days);
//
//	// function KeepRelayBeacon(address _stakingAddress, uint256 _minKeep) public {
//	KeepRelayBeacon relayBeaconContract = new KeepRelayBeacon(address(stakingContract), 3);
//
//	// Create KEEP Group Contract
//	KeepGroup keepGroupContract = new KeepGroup(10, address(relayBeaconContract));
	KeepGroup keepGroupContract = new KeepGroup(10, address(0x0));

	// test
	// 		function setGroupSize (uint256 _GroupSize) public onlyOwner {
	// 		function getNGroups () view public returns(uint256) {
	function testGroupSize2() public {
		keepGroupContract.setGroupSize(3);
		uint got = keepGroupContract.groupSize();
		Assert.equal(3, got, "Should have 3 in the gorup.");
	}

	// test
	// 		function createGroup (bytes32 _gID) public returns(bool) {
	function testCreateGroup00() public {
		bytes32 gID ;
		uint n;

		n = keepGroupContract.getNGroups();
		Assert.equal(n,0,"should have 0 groups.");

		gID = hex"0100";
		keepGroupContract.createGroup(gID);
		n = keepGroupContract.getNGroups();
		Assert.equal(n,1,"should have 1 groups.");
	}

	// test
	// 		function addMemberToGroup (bytes32 _gID, bytes32 _MemberPubKey) public isStaked returns(bool) {
	function testAddMemberToGroup01() public {
		bytes32 gID ;
		bytes32 memberPubKey ;
		bool ex;
		gID = hex"0100";
		memberPubKey = hex"0201";
		ex = keepGroupContract.addMemberToGroup (gID, memberPubKey);
		Assert.equal(ex,true,"should have success in adding member to group.");
	}

	// test
	// 		function addMemberToGroup (bytes32 _gID, bytes32 _MemberPubKey) public isStaked returns(bool) {
	function testAddMemberToGroup02() public {
		bytes32 gID ;
		bytes32 memberPubKey ;
		bool ex;
		gID = hex"0100";
		memberPubKey = hex"0202";
		ex = keepGroupContract.addMemberToGroup (gID, memberPubKey);
		Assert.equal(ex,true,"should have success in adding member to group.");
	}

	// test
    // 		function groupIsComplete (bytes32 _gID) public {
	function testIsGroupComplete00() public {
		bytes32 gID ;
		bool ex;
		gID = hex"0100";
		ex = keepGroupContract.groupIsComplete(gID);
		Assert.equal(ex,false,"group should be complete at this point.");
	}

	// test
	// 		function addMemberToGroup (bytes32 _gID, bytes32 _MemberPubKey) public isStaked returns(bool) {
	function testAddMemberToGroup03() public {
		bytes32 gID ;
		bytes32 memberPubKey ;
		bool ex;
		gID = hex"0100";
		memberPubKey = hex"0203";
		ex = keepGroupContract.addMemberToGroup (gID, memberPubKey);
		Assert.equal(ex,true,"should have success in adding member to group.");
	}

	// test
    // 		function groupIsComplete (bytes32 _gID) public {
	function testIsGroupComplete01() public {
		bytes32 gID ;
		bool ex;
		gID = hex"0100";
		ex = keepGroupContract.groupIsComplete(gID);
		Assert.equal(ex,true,"group should not be complete at this point.");
	}

}
