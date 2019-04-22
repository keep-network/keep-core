pragma solidity ^0.5.4;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/utils/ThrowProxy.sol";
import "../contracts/KeepGroupImplV1.sol";


contract StakingProxyMock {
    function balanceOf(address _staker) public pure returns(uint256) {
        _staker; // Suppress unused variable warning.
        return 200;
    }
}


contract TestKeepGroupInitialize {
    // Create Staking proxy contract mock
    StakingProxyMock stakingProxy = new StakingProxyMock();

    // Create KEEP Group Contract
    KeepGroupImplV1 keepGroupContract = new KeepGroupImplV1();

    function testCannotInitialize() public {

        // http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests
        ThrowProxy throwProxy = new ThrowProxy(address(keepGroupContract));

        // Prime the proxy
        KeepGroupImplV1(address(throwProxy)).initialize(address(0), address(0), 200, 150, 200, 1, 1, 1, 3);

        // Execute the call that is supposed to throw.
        // r will be false if it threw and true if it didn't.
        bool r = throwProxy.execute.gas(200000)();
        Assert.isFalse(r, "Should fail to initialize without Staking proxy address.");
    }

    function testInitialize() public {
        keepGroupContract.initialize(address(stakingProxy), address(0), 200, 150, 200, 1, 1, 1, 3);
        Assert.equal(keepGroupContract.initialized(), true, "Should be initialized.");
    }


}
