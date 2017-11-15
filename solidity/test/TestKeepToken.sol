pragma solidity ^0.4.18;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/KeepToken.sol";

contract TestKeepToken {
  KeepToken keepToken = KeepToken(DeployedAddresses.KeepToken());

  function testTotalSupply() {
    uint returned = keepToken.totalSupply();

    uint expected = 0;

    Assert.equal(returned, expected, "Should start with a totalSupply of 0.");
  }
}
