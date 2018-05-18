pragma solidity ^0.4.21;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenGrant.sol";
import "./helpers/ThrowProxy.sol";

contract TestTokenGrantReleaseThrow {
  
  // Create KEEP token.
  KeepToken t = new KeepToken();

  // Create token grant contract with 30 days withdrawal delay.
  TokenGrant c = new TokenGrant(t, 0, 30 days);

  uint id;
  address beneficiary = 0xf17f52151EbEF6C7334FAD080c5704D77216b732;
  uint start = now;
  uint duration = 10 days;
  uint cliff = 0;

  function testCannotReleaseGrantedAmmount() public {

    t.approve(address(c), 100);
    id = c.grant(100, beneficiary, duration, start, cliff, false);

    Assert.equal(c.unreleasedAmount(id), 0, "Unreleased token grant amount should be 0.");

    // http://truffleframework.com/tutorials/testing-for-throws-in-solidity-tests
    ThrowProxy throwProxy = new ThrowProxy(address(c));

    // Prime the proxy
    TokenGrant(address(throwProxy)).release(id);

    // Execute the call that is supposed to throw.
    // r will be false if it threw and true if it didn't.
    bool r = throwProxy.execute.gas(200000)();
    Assert.isFalse(r, "Should throw when trying to release token grant.");

    Assert.equal(t.balanceOf(beneficiary), 0, "Released balance should not be added to beneficiary main balance.");
  }

}
