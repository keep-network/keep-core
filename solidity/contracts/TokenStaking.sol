pragma solidity ^0.4.18;

import 'zeppelin-solidity/contracts/token/ERC20/StandardToken.sol';
import 'zeppelin-solidity/contracts/token/ERC20/SafeERC20.sol';
import 'zeppelin-solidity/contracts/math/SafeMath.sol';

/**
 * @title TokenStaking
 * @dev A token staking contract for a specified standard ERC20 token. 
 * A holder of the specified token can stake its tokens to this contract
 * and unstake after withdrawal delay is over.
 */
contract TokenStaking {
  using SafeMath for uint256;
  using SafeERC20 for StandardToken;

  StandardToken public token;

  event ReceivedApproval(uint256 _value);
  event Staked(address indexed from, uint256 value);
  event InitiatedUnstake(uint256 id);
  event FinishedUnstake();

  struct Withdrawal {
    address staker;
    uint256 amount;
    uint256 createdAt;
  }

  uint256 public withdrawalDelay;
  uint256 public numWithdrawals;

  mapping(address => uint256) public balances;
  mapping(address => uint256[]) public withdrawalIndices;
  mapping(uint256 => Withdrawal) public withdrawals;

  /**
   * @dev Creates a token staking contract for a provided Standard ERC20 token.
   * @param _tokenAddress Address of a token that will be linked to this contract.
   * @param _delay Withdrawal delay for unstake.
   */
  function TokenStaking(address _tokenAddress, uint256 _delay) {
    require(_tokenAddress != address(0x0));
    token = StandardToken(_tokenAddress);
    withdrawalDelay = _delay;
  }

  /** 
   * @notice Receives approval of token transfer and stakes the approved ammount.
   * @dev Makes sure provided token contract is the same one linked to this contract. 
   * @param _from The owner of the tokens who approved them to transfer.
   * @param _value Approved amount for the transfer and stake.
   * @param _token Token contract address.
   * @param _extraData Any extra data.
  */
  function receiveApproval(address _from, uint256 _value, address _token, bytes _extraData) {
    ReceivedApproval(_value);

    // Make sure provided token contract is the same one linked to this contract.
    require(StandardToken(_token) == token);

    // Make sure sender has enough tokens.
    require(_value <= token.balanceOf(_from));

    // Transfer tokens to this contract.
    token.transferFrom(_from, this, _value);

    // Maintain a record of the stake amount by the sender.
    balances[_from] = balances[_from].add(_value);
    Staked(_from, _value);
  }

  /**
   * @notice Initiates unstake of staked tokens and returns withdrawal request ID.
   * You will be able to call `finishUnstake()` with this ID and finish 
   * unstake once withdrawal delay is over. 
   * @param _value The amount to be unstaked.
   */
  function initiateUnstake(uint256 _value) public returns (uint256 id) {
    require(_value <= balances[msg.sender]);

    balances[msg.sender] = balances[msg.sender].sub(_value);
    
    id = numWithdrawals++;
    withdrawals[id] = Withdrawal(msg.sender, _value, now);
    withdrawalIndices[msg.sender].push(id);
    InitiatedUnstake(id);
    return id;
  }

  /**
   * @notice Finishes unstake of the tokens of provided withdrawal request. 
   * You can only finish unstake once withdrawal delay is over for the request, 
   * otherwise the function will fail and remaining gas is returned.
   * @param _id Withdrawal ID.
   */
  function finishUnstake(uint256 _id) public {
    require(now >= withdrawals[_id].createdAt.add(withdrawalDelay));

    // No need to call approve since msg.sender will be this staking contract.
    token.safeTransfer(withdrawals[_id].staker, withdrawals[_id].amount);

    // Cleanup.
    delete withdrawals[_id];

    FinishedUnstake();
  }

  /**
   * @dev Gets the stake balance of the specified address.
   * @param _staker The address to query the balance of.
   * @return An uint256 representing the amount owned by the passed address.
   */
  function balanceOf(address _staker) public constant returns (uint256 balance) {
    return balances[_staker];
  }

  /**
   * @dev Gets withdrawal request by ID.
   * @param _id ID of withdrawal request.
   * @return staker, amount, createdAt.
   */
  function getWithdrawal(uint256 _id) public constant returns (address, uint256, uint256) {
    return (withdrawals[_id].staker, withdrawals[_id].amount, withdrawals[_id].createdAt);
  }

  /**
   * @dev Gets withdrawal ids of the specified address.
   * @param _staker The address to query.
   * @return An uint256 array of withdrawal IDs.
   */
  function getWithdrawals(address _staker) public constant returns (uint256[]) {
    return withdrawalIndices[_staker];
  }
}
