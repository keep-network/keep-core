pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenVesting.sol";
import "./helpers/ThrowProxy.sol";

contract TestVestingRevoke {  
  
  // Create KEEP token
  KeepToken token = new KeepToken();

  // Create vesting contract with 30 days withdrawal delay
  TokenVesting vestingContract = new TokenVesting(token, 30 days);

  uint vestingId;
  address vestingBeneficiary = 0xf17f52151EbEF6C7334FAD080c5704D77216b732;

  // Vesting owner can revoke revocable vesting
  function testCanRevokeVestingFull() {
    uint balance = token.balanceOf(address(this));
  
    // create revocable vesting
    token.approve(address(vestingContract), 100);
    vestingId = vestingContract.vest(100, vestingBeneficiary, 10 days, now, 0, true);
    
    Assert.equal(token.balanceOf(address(this)), balance-100, "Amount should be removed from vesting owner balance.");
    Assert.equal(vestingContract.vestingBalanceOf(vestingBeneficiary), 100, "Amount should be added to beneficiary's vesting balance.");
    
    vestingContract.revoke(vestingId);

    Assert.equal(token.balanceOf(address(this)), balance, "Vesting amount should be returned to vesting owner.");
    Assert.equal(vestingContract.vestingBalanceOf(vestingBeneficiary), 0, "Amount should be removed from beneficiary's vesting balance");
  }

  // Vesting owner can revoke but no amount is refunded since duration of vesting is over
  function testCanRevokeVestingZeroRefund() {
    uint balance = token.balanceOf(address(this));
  
    // create revocable vesting with 0 duration
    token.approve(address(vestingContract), 100);
    vestingId = vestingContract.vest(100, vestingBeneficiary, 0, now, 0, true);
    
    Assert.equal(token.balanceOf(address(this)), balance-100, "Amount should be removed from vesting owner balance.");
    Assert.equal(vestingContract.vestingBalanceOf(vestingBeneficiary), 100, "Amount should be added to beneficiary's vesting balance.");
    
    vestingContract.revoke(vestingId);

    Assert.equal(token.balanceOf(address(this)), balance-100, "No amount to be returned to vesting owner since vesting duration is over.");
    Assert.equal(vestingContract.vestingBalanceOf(vestingBeneficiary), 100, "Amount should stay at beneficiary's vesting balance");
  }

  // Can not revoke vesting
  function testCanNotRevokeVesting() {

    token.approve(address(vestingContract), 100);
    // create non revocable vesting
    vestingId = vestingContract.vest(100, vestingBeneficiary, 10 days, now, 0, false);

    // http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests
    ThrowProxy throwProxy = new ThrowProxy(address(vestingContract));

    // Prime the proxy
    TokenVesting(address(throwProxy)).revoke(vestingId);

    // Execute the call that is supposed to throw.
    // r will be false if it threw and true if it didn't.
    bool r = throwProxy.execute.gas(200000)();
    Assert.isFalse(r, "Should throw when trying to revoke vesting");
  }
}
