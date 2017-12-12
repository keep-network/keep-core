pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenVesting.sol";
import "./helpers/ThrowProxy.sol";

contract TestVestingStake {  
  
  // Create KEEP token
  KeepToken token = new KeepToken();

  // Create vesting contract with 30 days withdrawal delay
  TokenVesting vestingContract = new TokenVesting(token, 30 days);

  uint vestingId;

  // Vesting params
  address vestingBeneficiary = address(this); // For test simplicity set beneficiary the same as sender
  uint vestingStart = now;
  uint vestingDuration = 10 days;
  uint vestingCliff = 0;

  // Vesting beneficiary should be able to stake unreleased vested balance
  function testCanStakeVesting() {

    // Create new vesting
    token.approve(address(vestingContract), 100);
    vestingId = vestingContract.vest(100, vestingBeneficiary, vestingDuration, vestingStart, vestingCliff, false);
    vestingContract.stakeVesting(vestingId);

    Assert.equal(vestingContract.vestingStakeBalances(vestingBeneficiary), 100, "Vesting balance should be added to beneficiary vesting stake balance.");

    var (owner, beneficiary, locked, revoked, revocable, amount, duration, start, cliff, released) = vestingContract.vestings(vestingId);
    Assert.equal(locked, true, "Vesting should become locked");
  }

  // Vesting beneficiary should be able to initiate unstake of vesting
  function testCanInitiateUnstakeVesting() {
    vestingContract.initiateUnstakeVesting(vestingId);
    Assert.equal(vestingContract.vestingStakeWithdrawalStart(vestingId), now, "Vesting withdrawals start should be set");
    Assert.equal(vestingContract.vestingStakeBalances(vestingBeneficiary), 100, "Vesting balance should stay unchanged.");
  }

  // Vesting beneficiary can not finish unstake of vesting until delay is over
  function testCannotFinishUnstakeVesting() {
  
    // http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests
    ThrowProxy throwProxy = new ThrowProxy(address(vestingContract));

    // Prime the proxy
    TokenVesting(address(throwProxy)).finishUnstakeVesting(vestingId);

    // Execute the call that is supposed to throw.
    // r will be false if it threw and true if it didn't.
    bool r = throwProxy.execute.gas(200000)();
    Assert.isFalse(r, "Should throw when trying to unstake vesting when delay is not over");
    Assert.equal(vestingContract.vestingStakeBalances(vestingBeneficiary), 100, "Vesting balance should stay unchanged.");
  }
}
      