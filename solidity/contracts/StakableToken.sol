pragma solidity ^0.4.18;

import 'zeppelin-solidity/contracts/token/StandardToken.sol';

/**
 * @title StakableToken
 * @dev Adds staking functionality to the token. 
 * Staking balance can be withdrawn by the token holder.
 * 
 */
contract StakableToken is StandardToken {
  using SafeMath for uint256;

  event Stake(address indexed from, uint256 value);
  event InitiateUnstake(uint256 id);
  event FinishUnstake();

  struct StakeWithdrawal {
    address owner;
    uint256 amount;
    uint256 start;
    bool released;
  }
  
  uint256 public stakeWithdrawalDelay;
  uint256 public numWithdrawals;

  // Stake balances
  mapping(address => uint256) public stakeBalances;

  // Stake withdrawals
  mapping(uint256 => StakeWithdrawal) public stakeWithdrawals;
  
  function StakableToken(uint256 _withdrawalDelay) {
    stakeWithdrawalDelay = _withdrawalDelay;
  }

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
  * @dev Initiate unstake of the tokens
  * @param _value The amount to be unstaked
  */
  function initiateUnstake(uint256 _value) public returns (uint256 id) {
    require(_value <= stakeBalances[msg.sender]);

    stakeBalances[msg.sender] = stakeBalances[msg.sender].sub(_value);
    
    // Create new stake withdrawal request
    id = numWithdrawals++;
    stakeWithdrawals[id] = StakeWithdrawal(msg.sender, _value, now, false);
    InitiateUnstake(id);
    return id;
  }

  /**
  * @dev Finish unstake of the tokens
  * @param _id Stake withdrawal ID
  */
  function finishUnstake(uint256 _id) public {
    require(!stakeWithdrawals[_id].released);
    require(now >= stakeWithdrawals[_id].start.add(stakeWithdrawalDelay));
    stakeWithdrawals[_id].released = true;
    balances[msg.sender] = balances[msg.sender].add(stakeWithdrawals[_id].amount);
    FinishUnstake();
  }

  /**
  * @dev Gets the stake balance of the specified address.
  * @param _owner The address to query the the balance of.
  * @return An uint256 representing the amount owned by the passed address.
  */
  function stakeBalanceOf(address _owner) public constant returns (uint256 stakeBalance) {
    return stakeBalances[_owner];
  }

  /**
  * @dev Gets withdrawal request by ID.
  * @param _id ID of withdrawal request.
  * @return owner, amount, start, released.
  */
  function getWithdrawal(uint256 _id) public constant returns (address, uint256, uint256, bool) {
    return (stakeWithdrawals[_id].owner, stakeWithdrawals[_id].amount, stakeWithdrawals[_id].start, stakeWithdrawals[_id].released);
  }
}
