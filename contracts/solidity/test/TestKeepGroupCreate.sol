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
		// I don't know if we can do this - since .initialize() checks to see if it is already been called.
        keepGroupContract.initialize(2, 3, address(keepRandomBeacon));
    }

    function testCreateGroup() public {
        keepGroupContract.createGroup(groupOnePubKey);
        Assert.equal(keepGroupContract.getGroupIndex(groupOnePubKey), 0, "Should get index of a group by its public key.");
        Assert.equal(keepGroupContract.getGroupPubKey(0), groupOnePubKey, "Should get public key of a group by its index.");
        Assert.equal(keepGroupContract.numberOfGroups(), 1, "Should get number of groups.");
    }

    function testGroupNotFound() public {
        // http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests
        ThrowProxy throwProxy = new ThrowProxy(address(keepGroupContract));

        // Prime the proxy
        KeepGroupImplV1(address(throwProxy)).getGroupIndex(groupTwoPubKey);

        // Execute the call that is supposed to throw.
        // r will be false if it threw and true if it didn't.
        bool r = throwProxy.execute.gas(200000)();
        Assert.isFalse(r, "Should fail to get index of a group that doesn't exist.");
    }

}
