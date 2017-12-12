pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenStaking.sol";

contract TestStakeNoDelay {
  // Create KEEP token
  KeepToken token = new KeepToken();

  // Create staking contract with no withdrawal delay
  TokenStaking stakingContract = new TokenStaking(token, 0);

  uint withdrawalId;

  // Token holder should be able to stake it's tokens
  function testCanStake() {
    uint balance = token.balanceOf(address(this));

    token.approve(address(stakingContract), 100);
    stakingContract.stake(100);
    
    Assert.equal(token.balanceOf(address(this)), balance - 100, "Stake amount should be taken out from token holder's main balance");
    Assert.equal(stakingContract.stakeBalanceOf(address(this)), 100, "Stake amount should be added to token holder's stake balance");
  }


  // Token holder should be able to initiate unstake of it's tokens
  function testCanInitiateUnstake() {
    uint balance = token.balanceOf(address(this));

    withdrawalId = stakingContract.initiateUnstake(100);

    // Inspect created withdrawal request
    var (owner, amount, start, released) = stakingContract.stakeWithdrawals(withdrawalId);
    Assert.equal(owner, address(this), "Withdrawal request should keep record of the owner");
    Assert.equal(amount, 100, "Withdrawal request should keep record of the amount");
    Assert.equal(start, now, "Withdrawal request should keep record of when it was initiated");
    Assert.equal(released, false, "Withdrawal request should not be marked as released");

    Assert.equal(stakingContract.stakeBalanceOf(address(this)), 0, "Unstake amount should be taken out from token holder's stake balance"); 
    Assert.equal(token.balanceOf(address(this)), balance, "Unstake amount should not be added to token holder main balance");
  }

  // Should be able to finish unstake of it's tokens when withdrawal delay is over
  function testCanFinishUnstake() {
    uint balance = token.balanceOf(address(this));

    stakingContract.finishUnstake(withdrawalId);
    Assert.equal(token.balanceOf(address(this)), balance + 100, "Unstake amount should be added to token holder main balance");
    Assert.equal(stakingContract.stakeBalanceOf(address(this)), 0, "Stake balance should be empty");

    // Inspect changes in withdrawal request
    var (owner, amount, start, released) = stakingContract.stakeWithdrawals(withdrawalId);
    Assert.equal(released, true, "Withdrawal request should be marked as released");
  }
}
