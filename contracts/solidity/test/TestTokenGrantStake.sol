pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenGrant.sol";
import "./helpers/ThrowProxy.sol";

contract TestTokenGrantStake {  
  
  // Create KEEP token
  KeepToken t = new KeepToken();

  // Create token grant contract with 30 days withdrawal delay.
  TokenGrant c = new TokenGrant(t, 30 days);

  uint id;
  address beneficiary = address(this); // For test simplicity set beneficiary the same as sender.
  uint start = now;
  uint duration = 10 days;
  uint cliff = 0;

  // Token grant beneficiary should be able to stake unreleased granted balance.
  function testCanStakeTokenGrant() {

    // Approve transfer of tokens to the token grant contract.
    t.approve(address(c), 100);
    // Create new token grant.
    id = c.grant(100, beneficiary, duration, start, cliff, false);
    // Stake token grant.
    c.stake(id);

    Assert.equal(c.stakeBalances(beneficiary), 100, "Token grant balance should be added to beneficiary grant stake balance.");

    var (_owner, _beneficiary, _locked, _revoked, _revocable, _amount, _duration, _start, _cliff, _released) = c.grants(id);
    Assert.equal(_locked, true, "Token grant should become locked.");
  }

  // Token grant beneficiary should be able to initiate unstake of the token grant
  function testCanInitiateUnstakeTokenGrant() {
    c.initiateUnstake(id);
    Assert.equal(c.stakeWithdrawalStart(id), now, "Stake withdrawal start should be set.");
    Assert.equal(c.stakeBalances(beneficiary), 100, "Stake balance should stay unchanged.");
  }

  // Token grant beneficiary can not finish unstake of the grant until delay is over
  function testCannotFinishUnstake() {
  
    // http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests
    ThrowProxy throwProxy = new ThrowProxy(address(c));

    // Prime the proxy.
    TokenGrant(address(throwProxy)).finishUnstake(id);

    // Execute the call that is supposed to throw.
    // r will be false if it threw and true if it didn't.
    bool r = throwProxy.execute.gas(200000)();
    Assert.isFalse(r, "Should throw when trying to unstake when delay is not over.");
    Assert.equal(c.stakeBalances(beneficiary), 100, "Stake balance should stay unchanged.");
  }
}
