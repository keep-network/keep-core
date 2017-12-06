pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "./helpers/ThrowProxy.sol";

contract TestVestingStakeNoDelay {  
  
  // Create KEEP token with no withdrawal delay
  KeepToken token = new KeepToken(0);

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

  // Vesting beneficiary can finish unstake of vesting when delay is over
  function testCanFinishUnstakeVesting() {
    token.finishUnstakeVesting(vestingId);
    Assert.equal(token.vestingStakeBalances(vestingBeneficiary), 0, "Vesting balance should change.");
    var (owner, beneficiary, locked, revoked, revocable, amount, duration, start, cliff, released) = token.vestings(vestingId);
    Assert.equal(locked, false, "Vesting should become unlocked");
  }
}
      