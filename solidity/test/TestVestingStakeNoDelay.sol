pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenVesting.sol";
import "./helpers/ThrowProxy.sol";

contract TestVestingStakeNoDelay {  
  
  // Create KEEP token
  KeepToken token = new KeepToken();

  // Create vesting contract with no withdrawal delay
  TokenVesting vestingContract = new TokenVesting(token, 0);

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

  // Vesting beneficiary can finish unstake of vesting when delay is over
  function testCanFinishUnstakeVesting() {
    vestingContract.finishUnstakeVesting(vestingId);
    Assert.equal(vestingContract.vestingStakeBalances(vestingBeneficiary), 0, "Vesting balance should change.");
    var (owner, beneficiary, locked, revoked, revocable, amount, duration, start, cliff, released) = vestingContract.vestings(vestingId);
    Assert.equal(locked, false, "Vesting should become unlocked");
  }
}
      