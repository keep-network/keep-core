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


contract TestKeepGroupCreate {
    // Create KEEP random beacon contract mock
    KeepRandomBeaconMock keepRandomBeacon = new KeepRandomBeaconMock();

    // Create KEEP Group Contract
    KeepGroupImplV1 keepGroupContract = new KeepGroupImplV1();

    bytes32 public groupOnePubKey = hex"0100";
    bytes32 public groupTwoPubKey = hex"0200";

    function beforeAll() public {
        keepGroupContract.initialize(2, 3, address(keepRandomBeacon));
        keepGroupContract.createGroup(groupOnePubKey);
    }

    function testCreateGroup() public {
        Assert.equal(keepGroupContract.getGroupIndex(groupOnePubKey), 0, "Should get index of a group by its public key.");
        Assert.equal(keepGroupContract.getGroupPubKey(0), groupOnePubKey, "Should get public key of a group by its index.");
        Assert.equal(keepGroupContract.numberOfGroups(), 1, "Should get number of groups.");
    }

    function testFindNonExistingGroup() public {
        bytes4 methodId = bytes4(keccak256("getGroupIndex(bytes32)"));
        Assert.isTrue(address(keepGroupContract).call(methodId, groupOnePubKey), "Should succeed to call to find existing group index.");
    }

    function testFindExistingGroup() public {
        bytes4 methodId = bytes4(keccak256("getGroupIndex(bytes32)"));
        Assert.isFalse(address(keepGroupContract).call(methodId, groupTwoPubKey), "Should fail to call to find non existing group index.");
    }

}
