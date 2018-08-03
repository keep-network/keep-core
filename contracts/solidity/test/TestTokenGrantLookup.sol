pragma solidity ^0.4.21;

import "truffle/Assert.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenGrant.sol";


contract TestTokenGrantLookup {

    // Create KEEP token.
    KeepToken t = new KeepToken();

    // Create token grant contract with 30 days withdrawal delay.
    TokenGrant c = new TokenGrant(t, 0, 30 days);

    uint id;
    address beneficiary = 0xf17f52151EbEF6C7334FAD080c5704D77216b732;
    uint start = now;
    uint duration = 10 days;
    uint cliff = 0;

    function beforeAll() public {
        t.approve(address(c), 100);
        id = c.grant(100, beneficiary, duration, start, cliff, false);
    }

    function testCanGetGrantByID() public {
        uint _amount;
        uint _released;
        bool _locked;
        bool _revoked;
        (_amount, _released, _locked, _revoked) = c.getGrant(id);
        Assert.equal(_amount, 100, "Grant should maintain a record of the granted amount.");
        Assert.equal(_released, 0, "Grant should have 0 amount released initially.");
        Assert.equal(_locked, false, "Grant should initially be unlocked.");
        Assert.equal(_revoked, false, "Grant should not be marked as revoked initially.");
    }
}
