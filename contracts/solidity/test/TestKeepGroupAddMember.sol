pragma solidity ^0.4.21;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "./helpers/ThrowProxy.sol";
import "../contracts/KeepGroupImplV1.sol";


contract KeepRandomBeaconMock {
    function hasMinimumStake(address _staker) public view returns(bool) {
        return true;
    }
}


contract TestKeepGroupAddMember {
    // Create KEEP random beacon contract mock
    KeepRandomBeaconMock keepRandomBeacon = new KeepRandomBeaconMock();

    // Create KEEP Group Contract
    KeepGroupImplV1 keepGroupContract = new KeepGroupImplV1();

    bytes32 public groupOnePubKey = hex"0100";
    bytes32 public memberOnePubKey = hex"0200";
    bytes32 public memberTwoPubKey = hex"0300";
    bytes32 public memberThreePubKey = hex"0400";

    function beforeAll() public {
        keepGroupContract.initialize(2, 3, address(keepRandomBeacon));
        keepGroupContract.createGroup(groupOnePubKey);
    }

    function testAddGroupMember() public {
        keepGroupContract.addMemberToGroup(groupOnePubKey, memberOnePubKey);
        Assert.equal(keepGroupContract.isMember(groupOnePubKey, memberOnePubKey), true, "Should be true if member is part of the group.");
        Assert.equal(keepGroupContract.isMember(groupOnePubKey, memberTwoPubKey), false, "Should be false if member is not part of the group.");
        Assert.equal(keepGroupContract.getGroupMemberPubKey(0, 0), memberOnePubKey, "Should get public key of a member by its and group index.");
    }

    function testGroupComplete() public {
        Assert.equal(keepGroupContract.groupIsComplete(groupOnePubKey), false, "Should be false if group is not complete");
        keepGroupContract.addMemberToGroup(groupOnePubKey, memberOnePubKey);
        keepGroupContract.addMemberToGroup(groupOnePubKey, memberTwoPubKey);
        keepGroupContract.addMemberToGroup(groupOnePubKey, memberThreePubKey);
        Assert.equal(keepGroupContract.groupIsComplete(groupOnePubKey), true, "Should be true if group is complete");
    }
}
