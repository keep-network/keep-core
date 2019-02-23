pragma solidity ^0.5.4;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenStaking.sol";
import "./helpers/ThrowProxy.sol";


contract TestStake {

  // Create KEEP token
  KeepToken token = new KeepToken();

  // Create staking contract with 30 days withdrawal delay
  TokenStaking stakingContract = new TokenStaking(address(token), address(0), 30 days);

  uint withdrawalId;

  function testTotalSupply() public {
    uint expected = token.INITIAL_SUPPLY();
    Assert.equal(token.balanceOf(address(this)), expected, "Owner should have all tokens initially.");
  }

  // Token holder should be able to stake it's tokens
  function testCanStake() public {
    uint balance = token.balanceOf(address(this));

    token.approveAndCall(address(stakingContract), 100, "0x00");
    
    Assert.equal(token.balanceOf(address(this)), balance - 100, "Stake amount should be taken out from token holder's main balance.");
    Assert.equal(stakingContract.stakeBalanceOf(address(this)), 100, "Stake amount should be added to token holder's stake balance.");
  }

  // Token holder should be able to initiate unstake of it's tokens
  function testCanInitiateUnstake() public {
    uint balance = token.balanceOf(address(this));

    withdrawalId = stakingContract.initiateUnstake(100);

    // Inspect created withdrawal request
    address owner;
    uint amount;
    uint start;
    (owner, amount, start) = stakingContract.withdrawals(withdrawalId);
    Assert.equal(owner, address(this), "Withdrawal request should maintain a record of the owner.");
    Assert.equal(amount, 100, "Withdrawal request should maintain a record of the amount.");
    Assert.equal(start, now, "Withdrawal request should maintain a record of when it was initiated.");

    Assert.equal(stakingContract.stakeBalanceOf(address(this)), 0, "Unstake amount should be taken out from token holder's stake balance."); 
    Assert.equal(token.balanceOf(address(this)), balance, "Unstake amount should not be added to token holder main balance.");
  }

  // Should not be able to finish unstake when withdrawal delay is not over
  function testCannotFinishUnstake() public {

    // http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests
    ThrowProxy throwProxy = new ThrowProxy(address(stakingContract));

    // Prime the proxy
    TokenStaking(address(throwProxy)).finishUnstake(withdrawalId);

    // Execute the call that is supposed to throw.
    // r will be false if it threw and true if it didn't.
    bool r = throwProxy.execute.gas(200000)();
    Assert.isFalse(r, "Should throw when trying to unstake when withdrawal delay is not over.");
    Assert.equal(stakingContract.stakeBalanceOf(address(this)), 0, "Stake balance should stay unchanged.");
  }

  // Token holder should not be able to stake without providing correct stakingContract address.
  function testCanNotStakeWithWrongRecipient() public {
    bytes memory callData = abi.encodeWithSignature("approveAndCall(address,uint256,bytes)", address(0), 100, "0x1234");
    (bool result, ) = address(token).call(callData);
    Assert.isFalse(result, "Should throw when trying to stake.");
  }
}
