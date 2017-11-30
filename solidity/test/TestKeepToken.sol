pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";

contract TestKeepToken {
  KeepToken keepToken = KeepToken(DeployedAddresses.KeepToken());

  function testTotalSupply() {
    // uint expected = 10000000000000000000000;
    // Assert.equal(keepToken.balanceOf(tx.origin), expected, "Owner should have all KEEP tokens initially");
  }

  function testCanStake() {
  }

  function testCanInitiateUnstake() {
  }

  function testCanFinishUnstakeAfterDelayIsOver() {
  }

  function testCannotFinishUnstakeBeforeDelayIsOver() {
  }

  function testCannotWithdrawStakeBeforeDelayIsOver() {
  }

}
