pragma solidity ^0.4.18;

import 'zeppelin-solidity/contracts/token/StandardToken.sol';

contract StakableToken is StandardToken {

  mapping(address => uint256) public stakeBalances;
  
  event Stake(address indexed from, uint256 value);
  event Unstake(address indexed from, uint256 value);

  /**
  * @dev Stake tokens
  * @param _value The amount to be staked
  */
  function stake(uint256 _value) public returns (bool) {
    require(_value <= balances[msg.sender]);

    balances[msg.sender] = balances[msg.sender].sub(_value);
    stakeBalances[msg.sender] = stakeBalances[msg.sender].add(_value);
    Stake(msg.sender, _value);
    return true;
  }

  /**
  * @dev Unstake tokens
  * @param _value The amount to be unstaked
  */
  function unstake(uint256 _value) public returns (bool) {
    require(_value <= stakeBalances[msg.sender]);

    stakeBalances[msg.sender] = stakeBalances[msg.sender].sub(_value);
    balances[msg.sender] = balances[msg.sender].add(_value);
    Unstake(msg.sender, _value);
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
