pragma solidity ^0.4.18;

import "../contracts/KeepGroup.sol";
import "../contracts/KeepRelayBeacon.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenStaking.sol";

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "./helpers/ThrowProxy.sol";

contract TestKeepGroup02 {	
	
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

	// Run the simpelest test to verify that the test is getting run.
	function testGroupSize() public {
		uint got = keepGroupContract.groupSize();
		Assert.equal(10, got, "Should have 10 in the gorup.");
	}

	// test
	// 		function setGroupSize (uint256 _GroupSize) public onlyOwner {
	// 		function getNGroups () view public returns(uint256) {
	function testGroupSize2() public {
		keepGroupContract.setGroupSize(5);
		uint got = keepGroupContract.groupSize();
		Assert.equal(5, got, "Should have 5 in the gorup.");
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
	// 		function createGroup (bytes32 _gID) public returns(bool) {
	// 		function getNGroups () view public returns(uint256) {
	function testCreateGroup01() public {
		bytes32 gID ;
		uint n;

		gID = hex"0101";
		keepGroupContract.createGroup(gID);
		n = keepGroupContract.getNGroups();
		Assert.equal(n,2,"should have 2 groups.");
	}

	// test
	// 		function groupExists (bytes32 _gID) public {
	function testGroupExistsView01() public {
		bytes32 gID ;
		bool ex;
		gID = hex"0100";
		ex = keepGroupContract.groupExistsView(gID);
		Assert.equal(ex,true,"should have found group.");
	}

	// test
	// 		function groupExists (bytes32 _gID) public {
	function testGroupExistsView02() public {
		bytes32 gID ;
		bool ex;
		gID = hex"0104";
		ex = keepGroupContract.groupExistsView(gID);
		Assert.equal(ex,false,"should **not** have found group.");
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
    //		function getGroupNMembers (uint256 _No) view public returns(uint256) {
	function testGetGroupNMembers01() public {
		uint n;
		n = keepGroupContract.getGroupNMembers(0);
		Assert.equal(n,1,"should have 1 member in group.");
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
    //		function getGroupNMembers (uint256 _No) view public returns(uint256) {
	function testGetGroupNMembers02() public {
		uint n;
		n = keepGroupContract.getGroupNMembers(0);
		Assert.equal(n,2,"should have 2 member in group.");
	}

	// test
    // 		function groupIsComplete (bytes32 _gID) public {
	function testIsGroupComplete00() public {
		bytes32 gID ;
		bool ex;
		gID = hex"0100";
		ex = keepGroupContract.groupIsComplete(gID);
		Assert.equal(ex,false,"group should not be complete at this point.");
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

}
