pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";
import "./utils/UintArrayUtils.sol";
import "./StakingProxy.sol";
import "./mixins/TokenStaking.sol";


/**
 * @title Staking
 * @dev A staking contract for a specified standard ERC20 token.
 * A holder of the specified token can stake its tokens to this contract
 * and unstake after withdrawal delay is over.
 */
contract Staking is TokenStaking {

    using UintArrayUtils for uint256[];

    event InitiatedUnstake(uint256 id);
    event FinishedUnstake(uint256 id);

    StakingProxy public stakingProxy;

    uint256 public stakeWithdrawalDelay;

    struct Withdrawal {
        address staker;
        uint256 amount;
        uint256 createdAt;
    }

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
        token = ERC20(_tokenAddress);
        stakingProxy = StakingProxy(_stakingProxy);
        stakeWithdrawalDelay = _delay;
    }

    /**
     * @dev Gets the stake balance of the specified address.
     * @param _address The address to query the balance of.
     * @return An uint256 representing the amount staked by the passed address.
     */
    function stakeBalanceOf(address _address) public view returns (uint256 balance) {
        return stakedBalances[_address];
    }

    /**
     * @notice Initiates unstake of staked tokens and returns withdrawal request ID.
     * You will be able to call `finishTokensUnstake()` with this ID and finish
     * unstake once withdrawal delay is over.
     * @param _value The amount to be unstaked.
     * @param _staker Address of the stake owner or its operator.
     */
    function initiateUnstake(uint256 _value, address _staker)
        public
        onlyOwnerOrOperator(_staker)
        returns (uint256 id)
    {

        require(_value <= stakeBalanceOf(_staker), "Staker must have enough tokens to unstake.");

        id = numWithdrawals++;
        withdrawals[id] = Withdrawal(_staker, _value, now);
        withdrawalIndices[_staker].push(id);
        emit InitiatedUnstake(id);
        if (address(stakingProxy) != address(0)) {
            stakingProxy.emitUnstakedEvent(_staker, _value);
        }
        return id;
    }

    /**
     * @notice Finishes unstake of the tokens of provided withdrawal request.
     * You can only finish unstake once withdrawal delay is over for the request,
     * otherwise the function will fail and remaining gas is returned.
     * @param _id Withdrawal ID.
     * @param _value Amount to withdraw.
     */
    function finishTokensUnstake(uint256 _id, uint256 _value)
        public
        onlyOwnerOrOperator(withdrawals[_id].staker)
    {
        require(now >= withdrawals[_id].createdAt.add(stakeWithdrawalDelay), "Can not finish unstake before withdrawal delay is over.");
        require(_value <= withdrawals[_id].amount, "Can not withdraw more than unstaked amount.");

        address staker = withdrawals[_id].staker;
        _transferUnstakedTokens(staker, _value);

        if (_value == withdrawals[_id].amount) {
            // Cleanup withdrawal records.
            withdrawalIndices[staker].removeValue(_id);
            delete withdrawals[_id];
        }

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
    function getWithdrawals(address _staker) public view returns (uint256[] memory) {
        return withdrawalIndices[_staker];
    }

    // TODO: replace with a secure authorization protocol (addressed in RFC 4).
    function authorizedTransferFrom(address from, address to, uint256 amount) public {
        stakedBalances[from] = stakedBalances[from].sub(amount);
        stakedBalances[to] = stakedBalances[to].add(amount);
    }

}
