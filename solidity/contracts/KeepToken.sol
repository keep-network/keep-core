pragma solidity ^0.4.18;

import 'zeppelin-solidity/contracts/token/MintableToken.sol';

contract KeepToken is MintableToken {
  string public name = "KEEP Token";
  string public symbol = "KEEP";
  uint256 public decimals = 18;

  mapping(address => uint256) stakeBalances;

  event StakeIn(address indexed from, uint256 value);
  event StakeOut(address indexed from, uint256 value);

  /**
  * @dev Stake in tokens
  * @param _value The amount to be staked in
  */
  function stakeIn(uint256 _value) public returns (bool) {
    require(_value <= balances[msg.sender]);

    balances[msg.sender] = balances[msg.sender].sub(_value);
    stakeBalances[msg.sender] = stakeBalances[msg.sender].add(_value);
    StakeIn(msg.sender, _value);
    return true;
  }

  /**
  * @dev Stake out tokens
  * @param _value The amount to be staked out
  */
  function stakeOut(uint256 _value) public returns (bool) {
    require(_value <= stakeBalances[msg.sender]);

    stakeBalances[msg.sender] = stakeBalances[msg.sender].sub(_value);
    balances[msg.sender] = balances[msg.sender].add(_value);
    StakeOut(msg.sender, _value);
    return true;
  }


  /**
  * @dev Gets the stake balance of the specified address.
  * @param _owner The address to query the the balance of.
  * @return An uint256 representing the amount owned by the passed address.
  */
  function stakeBalanceOf(address _owner) public constant returns (uint256 stakeBalance) {
    return stakeBalances[_owner];
  }
}
