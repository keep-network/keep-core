pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenGrant.sol";
import "./helpers/ThrowProxy.sol";

contract TestTokenGrantRevoke {  
  
  // Create KEEP token.
  KeepToken t = new KeepToken();

  // Create token grant contract with 30 days withdrawal delay.
  TokenGrant c = new TokenGrant(t, 30 days);

  uint id;
  address beneficiary = 0xf17f52151EbEF6C7334FAD080c5704D77216b732;

  // Grant owner can revoke revocable token grant.
  function testCanFullyRevokeGrant() {
    uint balance = t.balanceOf(address(this));
  
    // Create revocable token grant.
    t.approve(address(c), 100);
    id = c.grant(100, beneficiary, 10 days, now, 0, true);
    
    Assert.equal(t.balanceOf(address(this)), balance - 100, "Amount should be taken out from grant creator main balance.");
    Assert.equal(c.balanceOf(beneficiary), 100, "Amount should be added to beneficiary's granted balance.");
    
    c.revoke(id);

    Assert.equal(t.balanceOf(address(this)), balance, "Amount should be returned to token grant owner.");
    Assert.equal(c.balanceOf(beneficiary), 0, "Amount should be removed from beneficiary's grant balance.");
  }

  // Token grant creator can revoke the grant but no amount 
  // is refunded since duration of the vesting is over.
  function testCanZeroRevokeGrant() {
    uint balance = t.balanceOf(address(this));
  
    // Create revocable token grant with 0 duration.
    t.approve(address(c), 100);
    id = c.grant(100, beneficiary, 0, now, 0, true);
    
    Assert.equal(t.balanceOf(address(this)), balance - 100, "Amount should be removed from grant creator main balance.");
    Assert.equal(c.balanceOf(beneficiary), 100, "Amount should be added to beneficiary's granted balance.");
    
    c.revoke(id);

    Assert.equal(t.balanceOf(address(this)), balance - 100, "No amount to be returned to grant creator since vesting duration is over.");
    Assert.equal(c.balanceOf(beneficiary), 100, "Amount should stay at beneficiary's grant balance.");
  }

  // Test can not revoke token grant.
  function testCanNotRevokeGrant() {

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
