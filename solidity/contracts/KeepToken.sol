pragma solidity ^0.4.18;

import 'zeppelin-solidity/contracts/token/MintableToken.sol';

contract KeepToken is MintableToken {
  string public name = "KEEP TOKEN";
  string public symbol = "KEEP";
  uint256 public decimals = 18;
}
