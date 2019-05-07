pragma solidity ^0.5.4;

import "./StakeDelegatable.sol";
import "./utils/UintArrayUtils.sol";


/**
 * @title TokenStaking
 * @dev A token staking contract for a specified standard ERC20 token.
 * A holder of the specified token can stake its tokens to this contract
 * and unstake after withdrawal delay is over.
 */
contract Staking is StakeDelegatable {

    using UintArrayUtils for uint256[];

    event ReceivedApproval(uint256 _value);
    event Staked(address indexed from, uint256 value);
    event InitiatedUnstake(uint256 id);
    event FinishedUnstake(uint256 id);

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
     * @notice Receives approval of token transfer and stakes the approved ammount.
     * @dev Makes sure provided token contract is the same one linked to this contract.
     * @param _from The owner of the tokens who approved them to transfer.
     * @param _value Approved amount for the transfer and stake.
     * @param _token Token contract address.
     * @param _extraData Data for stake delegation. This byte array must have the
     * following values concatenated: Magpie address (20 bytes) where the rewards for participation
     * are sent and the operator's ECDSA (65 bytes) signature of the address of the stake owner.
     */
    function receiveApproval(address _from, uint256 _value, address _token, bytes memory _extraData) public {
        emit ReceivedApproval(_value);

        require(ERC20(_token) == token, "Token contract must be the same one linked to this contract.");
        require(_value <= token.balanceOf(_from), "Sender must have enough tokens.");
        require(_extraData.length == 85, "Stake delegation data must be provided.");

        address magpie = _extraData.toAddress(0);
        address operator = keccak256(abi.encodePacked(_from)).toEthSignedMessageHash().recover(_extraData.slice(20, 65));
        require(operatorToOwner[operator] == address(0), "Operator address is already in use.");

        operatorToOwner[operator] = _from;
        magpieToOwner[magpie] = _from;
        ownerOperators[_from].push(operator);

        // Transfer tokens to this contract.
        token.transferFrom(_from, address(this), _value);

        // Maintain a record of the stake amount by the sender.
        stakeBalances[operator] = stakeBalances[operator].add(_value);
        emit Staked(operator, _value);
        if (address(stakingProxy) != address(0)) {
            stakingProxy.emitStakedEvent(operator, _value);
        }
    }

    /**
     * @notice Initiates unstake of staked tokens and returns withdrawal request ID.
     * You will be able to call `finishUnstake()` with this ID and finish
     * unstake once withdrawal delay is over.
     * @param _value The amount to be unstaked.
     * @param _operator Address of the stake operator.
     */
    function initiateUnstake(uint256 _value, address _operator) public returns (uint256 id) {
        address owner = operatorToOwner[_operator];
        require(
            msg.sender == _operator ||
            msg.sender == owner, "Only operator or the owner of the stake can initiate unstake.");
        require(_value <= stakeBalances[_operator], "Staker must have enough tokens to unstake.");

        stakeBalances[_operator] = stakeBalances[_operator].sub(_value);

        id = numWithdrawals++;
        withdrawals[id] = Withdrawal(owner, _value, now);
        withdrawalIndices[owner].push(id);
        emit InitiatedUnstake(id);
        if (address(stakingProxy) != address(0)) {
            stakingProxy.emitUnstakedEvent(owner, _value);
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
        require(now >= withdrawals[_id].createdAt.add(stakeWithdrawalDelay), "Can not finish unstake before withdrawal delay is over.");

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
    function getWithdrawals(address _staker) public view returns (uint256[] memory) {
        return withdrawalIndices[_staker];
    }

    // TODO: replace with a secure authorization protocol (addressed in RFC 4).
    function authorizedTransferFrom(address from, address to, uint256 amount) public {
        stakeBalances[from] = stakeBalances[from].sub(amount);
        stakeBalances[to] = stakeBalances[to].add(amount);
    }

}
