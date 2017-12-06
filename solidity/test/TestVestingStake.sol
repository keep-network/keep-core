pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "./helpers/ThrowProxy.sol";

contract TestVestingStake {  
  
  // Create KEEP token with 30 days withdrawal delay
  KeepToken token = new KeepToken(30 days);

  uint vestingId;

  // Vesting params
  address vestingBeneficiary = address(this); // For test simplicity set beneficiary the same as sender
  uint vestingStart = now;
  uint vestingDuration = 10 days;
  uint vestingCliff = 0;

  // Vesting beneficiary should be able to stake unreleased vested balance
  function testCanStakeVesting() {

    // Create new vesting
    vestingId = token.vest(100, vestingBeneficiary, vestingDuration, vestingStart, vestingCliff, false);
    token.stakeVesting(vestingId);

    Assert.equal(token.vestingStakeBalances(vestingBeneficiary), 100, "Vesting balance should be added to beneficiary vesting stake balance.");

    var (owner, beneficiary, locked, revoked, revocable, amount, duration, start, cliff, released) = token.vestings(vestingId);
    Assert.equal(locked, true, "Vesting should become locked");
  }

  // Vesting beneficiary should be able to initiate unstake of vesting
  function testCanInitiateUnstakeVesting() {
    token.initiateUnstakeVesting(vestingId);
    Assert.equal(token.vestingStakeWithdrawalStart(vestingId), now, "Vesting withdrawals start should be set");
    Assert.equal(token.vestingStakeBalances(vestingBeneficiary), 100, "Vesting balance should stay unchanged.");
  }

  // Vesting beneficiary can not finish unstake of vesting until delay is over
  function testCannotFinishUnstakeVesting() {
  
    // http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests
    ThrowProxy throwProxy = new ThrowProxy(address(token));

    // Prime the proxy
    KeepToken(address(throwProxy)).finishUnstakeVesting(vestingId);

    // Execute the call that is supposed to throw.
    // r will be false if it threw and true if it didn't.
    bool r = throwProxy.execute.gas(200000)();
    Assert.isFalse(r, "Should throw when trying to unstake vesting when delay is not over");
    Assert.equal(token.vestingStakeBalances(vestingBeneficiary), 100, "Vesting balance should stay unchanged.");
  }
}
      