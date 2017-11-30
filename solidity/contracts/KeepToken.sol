pragma solidity ^0.4.18;

import './StakableToken.sol';
import './VestableToken.sol';

/**
 * @title KEEP Token
 * @dev Combines functionality of "Stakable" and "Vestable" tokens
 * Token holder can stake both normal and vested tokens
 * Withdrawal is only possible after withdrawal delay period is over
 */
contract KeepToken is StakableToken, VestableToken {
  string public constant NAME = "KEEP Token";
  string public constant SYMBOL = "KEEP";
  uint256 public constant DECIMALS = 18;
  uint256 public constant INITIAL_SUPPLY = 10000 * (10 ** uint256(DECIMALS));

  /**
   * @dev Gives msg.sender all of existing tokens.
   */
  function KeepToken() {
    totalSupply = INITIAL_SUPPLY;
    balances[msg.sender] = INITIAL_SUPPLY;

    // Stakable token withdrawal delay
    stakeWithdrawalDelay = 1000;

    // Vestable token withdrawal delay
    vestingWithdrawalDelay = 1000;
  }
}
