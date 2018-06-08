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


contract TestKeepGroupInitialize {
    // Create KEEP random beacon contract mock
    KeepRandomBeaconMock keepRandomBeacon = new KeepRandomBeaconMock();

    // Create KEEP Group Contract
    KeepGroupImplV1 keepGroupContract = new KeepGroupImplV1();

    function testCannotInitialize() public {

        // http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests
        ThrowProxy throwProxy = new ThrowProxy(address(keepGroupContract));

        // Prime the proxy
        KeepGroupImplV1(address(throwProxy)).initialize(2, 3, 0);

        // Execute the call that is supposed to throw.
        // r will be false if it threw and true if it didn't.
        bool r = throwProxy.execute.gas(200000)();
        Assert.isFalse(r, "Should fail to initialize without KEEP random beacon contract address.");
    }

    function testInitialize() public {
        keepGroupContract.initialize(2, 3, address(keepRandomBeacon));
        Assert.equal(keepGroupContract.initialized(), true, "Should be initialized.");
    }


}
