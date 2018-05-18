pragma solidity ^0.4.21;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";
import "../contracts/TokenGrant.sol";
import "./helpers/ThrowProxy.sol";

contract TestTokenGrantRelease {  
  
  // Create KEEP token.
  KeepToken t = new KeepToken();

  // Create token grant contract with 30 days withdrawal delay.
  TokenGrant c = new TokenGrant(t, 0, 30 days);

  uint id;
  address beneficiary = 0xf17f52151EbEF6C7334FAD080c5704D77216b732;
  uint start = now;
  uint duration = 10 days;
  uint cliff = 0;

  // Token holder should be able to grant it's tokens to a beneficiary.
  function testCanGrant() public {
    uint balance = t.balanceOf(address(this));

    t.approve(address(c), 100);
    id = c.grant(100, beneficiary, duration, start, cliff, false);
    Assert.equal(t.balanceOf(address(this)), balance - 100, "Amount should be taken out from grant creator main balance.");
    Assert.equal(c.balanceOf(beneficiary), 100, "Amount should be added to beneficiary's granted balance.");
  }

  // Token grant with 0 duration, release available immediately.
  function testCanReleaseGrantedAmount() public {
    
    t.approve(address(c), 100);
    id = c.grant(100, beneficiary, 0, now, 0, false);

    c.release(id);
    Assert.equal(t.balanceOf(beneficiary), 100, "Released balance should be added to beneficiary main balance.");
    Assert.equal(c.unreleasedAmount(id), 0, "Unreleased granted amount should be 0 after release.");
  }
}
