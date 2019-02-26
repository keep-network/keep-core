pragma solidity ^0.5.4;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenGrant.sol";
import "./helpers/ThrowProxy.sol";

contract TestTokenGrantNonRevocable {  
    
    // Create KEEP token.
    KeepToken t = new KeepToken();

    // Create token grant contract with 30 days withdrawal delay.
    TokenGrant c = new TokenGrant(address(t), address(0), 30 days);

    uint id;
    address beneficiary = 0xf17f52151EbEF6C7334FAD080c5704D77216b732;

    // Test can not revoke token grant.
    function testCanNotRevokeGrant() public {

        // Approve amount on the token.
        t.approve(address(c), 100);

        // Create non revocable token grant.
        id = c.grant(100, beneficiary, 10 days, now, 0, false);

        // http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests
        ThrowProxy throwProxy = new ThrowProxy(address(c));

        // Prime the proxy.
        TokenGrant(address(throwProxy)).revoke(id);

        // Execute the call that is supposed to throw.
        // r will be false if it threw and true if it didn't.
        bool r = throwProxy.execute.gas(200000)();
        Assert.isFalse(r, "Should throw when trying to revoke token grant.");
    }
}
