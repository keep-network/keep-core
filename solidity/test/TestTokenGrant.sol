pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenGrant.sol";
import "./helpers/ThrowProxy.sol";

contract TestTokenGrant {  
  
  // Create KEEP token.
  KeepToken t = new KeepToken();

  // Create token grant contract with 30 days withdrawal delay.
  TokenGrant c = new TokenGrant(t, 30 days);

  uint id;
  address beneficiary = 0xf17f52151EbEF6C7334FAD080c5704D77216b732;
  uint start = now;
  uint duration = 10 days;
  uint cliff = 0;

  // Token holder should be able to grant it's tokens to a beneficiary.
  function testCanGrant() {
    uint balance = t.balanceOf(address(this));

    t.approve(address(c), 100);
    id = c.grant(100, beneficiary, duration, start, cliff, false);
    Assert.equal(t.balanceOf(address(this)), balance - 100, "Amount should be taken out from grant creator main balance.");
    Assert.equal(c.balanceOf(beneficiary), 100, "Amount should be added to beneficiary's granted balance.");
  }

  function testCanGetGrantByID() {
    var (_amount, _released, _locked, _revoked) = c.getGrant(id);
    Assert.equal(_amount, 100, "Grant should keep record of the granted amount.");
    Assert.equal(_released, 0, "Grant should have 0 amount released initially.");
    Assert.equal(_locked, false, "Grant should initially be unlocked.");
    Assert.equal(_revoked, false, "Grant should not be marked as revoked initially.");
  }

  function testCanGetGrantVestingScheduleByGrantID() {
    var (_owner, _duration, _start, _cliff) = c.getGrantVestingSchedule(id);
    Assert.equal(_owner, address(this), "Grant should keep record of the creator.");
    Assert.equal(_duration, duration, "Grant should have vesting schedule duration.");
    Assert.equal(_start, start, "Grant should have start time.");
    Assert.equal(_cliff, start+cliff, "Grant should have vesting schedule cliff duration.");
  }

  function testCannotReleaseGrantedAmmount() {
    Assert.equal(c.unreleasedAmount(id), 0, "Unreleased token grant amount should be 0.");

    // http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests
    ThrowProxy throwProxy = new ThrowProxy(address(c));

    // Prime the proxy
    TokenGrant(address(throwProxy)).release(id);

    // Execute the call that is supposed to throw.
    // r will be false if it threw and true if it didn't.
    bool r = throwProxy.execute.gas(200000)();
    Assert.isFalse(r, "Should throw when trying to release token grant.");

    Assert.equal(t.balanceOf(beneficiary), 0, "Released balance should not be added to beneficiary main balance.");
  }

  // TOken grant with 0 duration, release available immediately.
  function testCanReleaseGrantedAmount() {
    
    t.approve(address(c), 100);
    id = c.grant(100, beneficiary, 0, now, 0, false);

    c.release(id);
    Assert.equal(t.balanceOf(beneficiary), 100, "Released balance should be added to beneficiary main balance.");
    Assert.equal(c.unreleasedAmount(id), 0, "Unreleased granted amount should be 0 after release.");
  }
}
