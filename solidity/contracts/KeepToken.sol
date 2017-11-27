pragma solidity ^0.4.18;

import './StakableToken.sol';
import './VestableToken.sol';

contract KeepToken is StakableToken, VestableToken {
  string public name = "KEEP Token";
  string public symbol = "KEEP";
  uint256 public decimals = 18;

  // TODO
  // configurable withdrawal delay for staking,.
  // initiateUnstake, finishUnstake... or similar.
}
