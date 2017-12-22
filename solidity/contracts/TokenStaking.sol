pragma solidity ^0.4.18;

import 'zeppelin-solidity/contracts/token/BasicToken.sol';
import 'zeppelin-solidity/contracts/token/StandardToken.sol';
import 'zeppelin-solidity/contracts/token/SafeERC20.sol';
import 'zeppelin-solidity/contracts/math/SafeMath.sol';

/**
 * @title TokenStaking
 * @dev A token staking contract for a specified standard ERC20 token. 
 * A holder of the specified token can stake its tokens to this contract
 * and unstake after withdrawal delay is over.
 */
contract TokenStaking is BasicToken {
  using SafeMath for uint256;
  using SafeERC20 for StandardToken;

  StandardToken public token;

  event Stake(address indexed from, uint256 value);
  event InitiateUnstake(uint256 id);
  event FinishUnstake();

  struct StakeWithdrawal {
    address staker;
    uint256 amount;
    uint256 createdAt;
    bool released;
  }

  uint256 public stakeWithdrawalDelay;
  uint256 public numWithdrawals;

  mapping(address => uint256) public stakeBalances;
  mapping(uint256 => StakeWithdrawal) public stakeWithdrawals;

  /**
   * @dev Creates a token staking contract for a provided Standard ERC20 token.
   * @param _tokenAddress address of a token that will be linked to this contract.
   * @param _delay withdrawal delay for unstake.
   */
  function TokenStaking(address _tokenAddress, uint256 _delay) {
    require(_tokenAddress != address(0x0));
    token = StandardToken(_tokenAddress);
    stakeWithdrawalDelay = _delay;
  }

  /**
   * @notice Stakes provided token amount to this contract. You must approve the amount on the token contract first.
   * @dev Transfers tokens from sender balance to this staking contract balance.
   * Sender should approve the amount first by calling approve() on the token.
   * @param _value The amount to be staked.
   */
  function stake(uint256 _value) public {
    // Make sure sender has enough tokens.
    require(_value <= token.balanceOf(msg.sender));

    token.transferFrom(msg.sender, this, _value);

    // Keep record of the stake amount by the sender.
    stakeBalances[msg.sender] = stakeBalances[msg.sender].add(_value);
    Stake(msg.sender, _value);
  }

  /**
   * @notice Initiates unstake of staked tokens.
   * @dev Creates a new stake withdrawal request that 
   * can be checked later and processed if withdrawal delay is over.
   * @param _value The amount to be unstaked.
   */
  function initiateUnstake(uint256 _value) public returns (uint256 id) {
    require(_value <= stakeBalances[msg.sender]);

    stakeBalances[msg.sender] = stakeBalances[msg.sender].sub(_value);
    
    id = numWithdrawals++;
    stakeWithdrawals[id] = StakeWithdrawal(msg.sender, _value, now, false);
    InitiateUnstake(id);
    return id;
  }

  /**
   * @notice Finishes unstake of the tokens.
   * @dev Transfers tokens from this staking contract balance 
   * to the staker token balance if the withdrawal delay is over.
   * @param _id Stake withdrawal ID.
   */
  function finishUnstake(uint256 _id) public {
    require(!stakeWithdrawals[_id].released);
    require(now >= stakeWithdrawals[_id].createdAt.add(stakeWithdrawalDelay));
    stakeWithdrawals[_id].released = true;
    
    // No need to call approve since msg.sender will be this staking contract.
    token.safeTransfer(stakeWithdrawals[_id].staker, stakeWithdrawals[_id].amount);

    FinishUnstake();
  }

  /**
   * @dev Gets the stake balance of the specified address.
   * @param _staker The address to query the the balance of.
   * @return An uint256 representing the amount owned by the passed address.
   */
  function stakeBalanceOf(address _staker) public constant returns (uint256 stakeBalance) {
    return stakeBalances[_staker];
  }

  /**
   * @dev Gets withdrawal request by ID.
   * @param _id ID of withdrawal request.
   * @return staker, amount, createdAt, released.
   */
  function getWithdrawal(uint256 _id) public constant returns (address, uint256, uint256, bool) {
    return (stakeWithdrawals[_id].staker, stakeWithdrawals[_id].amount, stakeWithdrawals[_id].createdAt, stakeWithdrawals[_id].released);
  }
}
