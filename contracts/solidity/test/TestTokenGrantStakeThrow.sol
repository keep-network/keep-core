pragma solidity ^0.4.21;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenGrant.sol";
import "./helpers/ThrowProxy.sol";

contract TestTokenGrantStakeThrow {
  
  // Create KEEP token
  KeepToken t = new KeepToken();

  // Create token grant contract with 30 days withdrawal delay.
  TokenGrant c = new TokenGrant(t, 0, 30 days);

  uint id;
  address beneficiary = address(this); // For test simplicity set beneficiary the same as sender.
  uint start = now;
  uint duration = 10 days;
  uint cliff = 0;

  // Token grant beneficiary can not finish unstake of the grant until delay is over
  function testCannotFinishUnstake() public {
    
    // Approve transfer of tokens to the token grant contract.
    t.approve(address(c), 100);
    // Create new token grant.
    id = c.grant(100, beneficiary, duration, start, cliff, false);
    // Stake token grant.
    c.stake(id);

    c.initiateUnstake(id);
  
    // http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests
    ThrowProxy throwProxy = new ThrowProxy(address(c));

    // Prime the proxy.
    TokenGrant(address(throwProxy)).finishUnstake(id);

    // Execute the call that is supposed to throw.
    // r will be false if it threw and true if it didn't.
    bool r = throwProxy.execute.gas(200000)();
    Assert.isFalse(r, "Should throw when trying to unstake when delay is not over.");
    Assert.equal(c.stakeBalances(beneficiary), 0, "Stake balance should stay unchanged.");
  }
}
