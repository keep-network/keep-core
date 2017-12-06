pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "./helpers/ThrowProxy.sol";

contract TestStake {  
  
  // Create KEEP token with 30 days withdrawal delay
  KeepToken token = new KeepToken(30 days);

  uint withdrawalId;

  function testTotalSupply() {
    uint expected = token.INITIAL_SUPPLY();
    Assert.equal(token.balanceOf(address(this)), expected, "Owner should have all tokens initially");
  }

  // Token holder should be able to stake it's tokens
  function testCanStake() {
    uint balance = token.balanceOf(address(this));
    token.stake(100);
    Assert.equal(token.balanceOf(address(this)), balance - 100, "Stake amount should be taken out from token holder's main balance");
    Assert.equal(token.stakeBalanceOf(address(this)), 100, "Stake amount should be added to token holder's stake balance");
  }

  // Token holder should be able to initiate unstake of it's tokens
  function testCanInitiateUnstake() {
    uint balance = token.balanceOf(address(this));

    withdrawalId = token.initiateUnstake(100);

    // Inspect created withdrawal request
    var (owner, amount, start, released) = token.stakeWithdrawals(withdrawalId);
    Assert.equal(owner, address(this), "Withdrawal request should keep record of the owner");
    Assert.equal(amount, 100, "Withdrawal request should keep record of the amount");
    Assert.equal(start, now, "Withdrawal request should keep record of when it was initiated");
    Assert.equal(released, false, "Withdrawal request should not be marked as released");

    Assert.equal(token.stakeBalanceOf(address(this)), 0, "Unstake amount should be taken out from token holder's stake balance"); 
    Assert.equal(token.balanceOf(address(this)), balance, "Unstake amount should not be added to token holder main balance");
  }

  // Should not be able to finish unstake when withdrawal delay is not over
  function testCannotFinishUnstake() {

    // http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests
    ThrowProxy throwProxy = new ThrowProxy(address(token));

    // Prime the proxy
    KeepToken(address(throwProxy)).finishUnstake(withdrawalId);

    // Execute the call that is supposed to throw.
    // r will be false if it threw and true if it didn't.
    bool r = throwProxy.execute.gas(200000)();
    Assert.isFalse(r, "Should throw when trying to unstake when withdrawal delay is not over");
    Assert.equal(token.stakeBalanceOf(address(this)), 0, "Stake balance should stay unchanged");
  }
}
