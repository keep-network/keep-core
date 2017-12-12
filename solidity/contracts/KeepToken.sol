pragma solidity ^0.4.18;

import 'zeppelin-solidity/contracts/token/StandardToken.sol';

/**
 * @title KEEP Token
 * @dev Standard ERC20 token
 */
contract KeepToken is StandardToken {
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
  }

}
