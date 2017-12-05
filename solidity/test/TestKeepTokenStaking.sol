pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";

contract TestKeepTokenStaking {  
  
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
    withdrawalId = token.initiateUnstake(100);

    // Withdrawal request
    var (owner, amount, start, released) = token.stakeWithdrawals(withdrawalId);
    Assert.equal(owner, address(this), "Withdrawal request should keep record of the owner");
    Assert.equal(amount, 100, "Withdrawal request should keep record of the amount");
    Assert.equal(start, now, "Withdrawal request should keep record of when it was initiated");
    Assert.equal(released, false, "Withdrawal request should be marked as non released");

    uint balance = token.balanceOf(address(this));
    Assert.equal(token.stakeBalanceOf(address(this)), 0, "Unstake amount should be taken out from token holder's stake balance"); 
    Assert.equal(token.balanceOf(address(this)), balance, "Unstake amount should not be added to token holder main balance");
  }
}

contract TestKeepTokenNoDelay {
  // Create KEEP token with no withdrawal delay
  KeepToken token = new KeepToken(0);
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
    withdrawalId = token.initiateUnstake(100);
    uint balance = token.balanceOf(address(this));
    Assert.equal(token.stakeBalanceOf(address(this)), 0, "Unstake amount should be taken out from token holder's stake balance"); 
    Assert.equal(token.balanceOf(address(this)), balance + 100, "Unstake amount should be added to token holder balance immediately");

    var (owner, amount, start, released) = token.stakeWithdrawals(withdrawalId);
    Assert.equal(owner, address(this), "Withdrawal request should keep record of the owner");
    Assert.equal(amount, 100, "Withdrawal request should keep record of the amount");
    Assert.equal(start, now, "Withdrawal request should keep record of when it was initiated");
    Assert.equal(released, true, "Withdrawal request should be marked as released");
  }
}
