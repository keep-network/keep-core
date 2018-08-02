pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC20/StandardToken.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./StakingProxy.sol";
import "./utils/UintArrayUtils.sol";
import "./StakeDelegatable.sol";


/**
 * @title TokenStaking
 * @dev A token staking contract for a specified standard ERC20 token.
 * A holder of the specified token can stake its tokens to this contract
 * and unstake after withdrawal delay is over.
 */
contract TokenStaking is StakeDelegatable {
    using SafeMath for uint256;
    using SafeERC20 for StandardToken;
    using UintArrayUtils for uint256[];

    StandardToken public token;
    StakingProxy public stakingProxy;

    event ReceivedApproval(uint256 _value);
    event Staked(address indexed from, uint256 value);
    event InitiatedUnstake(uint256 id);
    event FinishedUnstake(uint256 id);

    struct Withdrawal {
        address staker;
        uint256 amount;
        uint256 createdAt;
    }

    uint256 public withdrawalDelay;
    uint256 public numWithdrawals;

    mapping(address => uint256[]) public withdrawalIndices;
    mapping(uint256 => Withdrawal) public withdrawals;

    /**
     * @dev Creates a token staking contract for a provided Standard ERC20 token.
     * @param _tokenAddress Address of a token that will be linked to this contract.
     * @param _stakingProxy Address of a staking proxy that will be linked to this contract.
     * @param _delay Withdrawal delay for unstake.
     */
    constructor(address _tokenAddress, address _stakingProxy, uint256 _delay) public {
        require(_tokenAddress != address(0x0), "Token address can't be zero.");
        token = StandardToken(_tokenAddress);
        stakingProxy = StakingProxy(_stakingProxy);
        withdrawalDelay = _delay;
    }

    /**
     * @notice Receives approval of token transfer and stakes the approved ammount.
     * @dev Makes sure provided token contract is the same one linked to this contract.
     * @param _from The owner of the tokens who approved them to transfer.
     * @param _value Approved amount for the transfer and stake.
     * @param _token Token contract address.
     * @param extraData_ Any extra data.
     */
    function receiveApproval(address _from, uint256 _value, address _token, bytes extraData_) public {
        extraData_; // Suppress unused variable warning.
        emit ReceivedApproval(_value);

        require(StandardToken(_token) == token, "Token contract must be the same one linked to this contract.");
        require(_value <= token.balanceOf(_from), "Sender must have enough tokens.");

        revertIfDelegateStakes(_from);

        // Transfer tokens to this contract.
        token.transferFrom(_from, this, _value);

        // Maintain a record of the stake amount by the sender.
        stakeBalances[_from] = stakeBalances[_from].add(_value);

        // Emit staked event. Check if staker works via operator first.
        address delegatorOrDelegate = getDelegatorOrDelegate(_from);
        emit Staked(delegatorOrDelegate, _value);
        if (address(stakingProxy) != address(0)) {
            stakingProxy.emitStakedEvent(delegatorOrDelegate, _value);
        }
    }

    /**
     * @notice Initiates unstake of staked tokens and returns withdrawal request ID.
     * You will be able to call `finishUnstake()` with this ID and finish
     * unstake once withdrawal delay is over.
     * @param _value The amount to be unstaked.
     */
    function initiateUnstake(uint256 _value) public returns (uint256 id) {

        require(_value <= stakeBalances[msg.sender], "Staker must have enough tokens to unstake.");

        stakeBalances[msg.sender] = stakeBalances[msg.sender].sub(_value);

        id = numWithdrawals++;
        withdrawals[id] = Withdrawal(msg.sender, _value, now);
        withdrawalIndices[msg.sender].push(id);
        emit InitiatedUnstake(id);

        // Emit unstaked event. Check if staker delegated its balance first.
        address delegatorOrDelegate = getDelegatorOrDelegate(msg.sender);
        if (address(stakingProxy) != address(0)) {
            stakingProxy.emitUnstakedEvent(delegatorOrDelegate, _value);
        }
        return id;
    }

    /**
     * @notice Finishes unstake of the tokens of provided withdrawal request.
     * You can only finish unstake once withdrawal delay is over for the request,
     * otherwise the function will fail and remaining gas is returned.
     * @param _id Withdrawal ID.
     */
    function finishUnstake(uint256 _id) public {
        require(now >= withdrawals[_id].createdAt.add(withdrawalDelay), "Can not finish unstake before withdrawal delay is over.");

        address staker = withdrawals[_id].staker;

        // No need to call approve since msg.sender will be this staking contract.
        token.safeTransfer(staker, withdrawals[_id].amount);

        // Cleanup withdrawal index.
        withdrawalIndices[staker].removeValue(_id);

        // Cleanup withdrawal record.
        delete withdrawals[_id];

        emit FinishedUnstake(_id);
    }

    /**
     * @dev Gets withdrawal request by ID.
     * @param _id ID of withdrawal request.
     * @return staker, amount, createdAt.
     */
    function getWithdrawal(uint256 _id) public view returns (address, uint256, uint256) {
        return (withdrawals[_id].staker, withdrawals[_id].amount, withdrawals[_id].createdAt);
    }

    /**
     * @dev Gets withdrawal ids of the specified address.
     * @param _staker The address to query.
     * @return An uint256 array of withdrawal IDs.
     */
    function getWithdrawals(address _staker) public view returns (uint256[]) {
        return withdrawalIndices[_staker];
    }
}
