pragma solidity ^0.5.4;

import "./StakeDelegatable.sol";
import "./utils/UintArrayUtils.sol";
import "./Registry.sol";


/**
 * @title TokenStaking
 * @dev A token staking contract for a specified standard ERC20Burnable token.
 * A holder of the specified token can stake its tokens to this contract
 * and unstake after withdrawal delay is over.
 */
contract TokenStaking is StakeDelegatable {

    using UintArrayUtils for uint256[];

    event Staked(address indexed from, uint256 value);
    event InitiatedUnstake(address indexed operator, uint256 value, uint256 createdAt);
    event FinishedUnstake(address operator);

    struct Withdrawal {
        uint256 amount;
        uint256 createdAt;
    }

    // Registry contract with a list of approved operator contracts and upgraders.
    Registry public registry;

    // Authorized operator contracts.
    mapping(address => mapping (address => bool)) internal authorizations;

    mapping(address => Withdrawal) public withdrawals;

    modifier onlyApprovedOperatorContract(address operatorContract) {
        require(
            registry.isApprovedOperatorContract(operatorContract),
            "Operator contract is not approved"
        );
        _;
    }

    /**
     * @dev Creates a token staking contract for a provided Standard ERC20Burnable token.
     * @param _tokenAddress Address of a token that will be linked to this contract.
     * @param _registry Address of a keep registry that will be linked to this contract.
     * @param _delay Withdrawal delay for unstake.
     */
    constructor(address _tokenAddress, address _registry, uint256 _delay) public {
        require(_tokenAddress != address(0x0), "Token address can't be zero.");
        token = ERC20Burnable(_tokenAddress);
        registry = Registry(_registry);
        stakeWithdrawalDelay = _delay;
    }

    /**
     * @notice Receives approval of token transfer and stakes the approved amount.
     * @dev Makes sure provided token contract is the same one linked to this contract.
     * @param _from The owner of the tokens who approved them to transfer.
     * @param _value Approved amount for the transfer and stake.
     * @param _token Token contract address.
     * @param _extraData Data for stake delegation. This byte array must have the
     * following values concatenated: Magpie address (20 bytes) where the rewards for participation
     * are sent, operator's (20 bytes) address, authorizer (20 bytes) address.
     */
    function receiveApproval(address _from, uint256 _value, address _token, bytes memory _extraData) public {
        require(ERC20Burnable(_token) == token, "Token contract must be the same one linked to this contract.");
        require(_value <= token.balanceOf(_from), "Sender must have enough tokens.");
        require(_extraData.length == 60, "Stake delegation data must be provided.");

        address payable magpie = address(uint160(_extraData.toAddress(0)));
        address operator = _extraData.toAddress(20);
        require(operatorToOwner[operator] == address(0), "Operator address is already in use.");
        address authorizer = _extraData.toAddress(40);

        operatorToOwner[operator] = _from;
        operatorToMagpie[operator] = magpie;
        operatorToAuthorizer[operator] = authorizer;
        ownerOperators[_from].push(operator);

        // Transfer tokens to this contract.
        token.transferFrom(_from, address(this), _value);

        // Maintain a record of the stake amount by the sender.
        stakeBalances[operator] = stakeBalances[operator].add(_value);
        emit Staked(operator, _value);
    }

    /**
     * @notice Initiates unstake of staked tokens and returns withdrawal request ID.
     * You will be able to call `finishUnstake()` with this ID and finish
     * unstake once withdrawal delay is over.
     * @param _value The amount to be unstaked.
     * @param _operator Address of the stake operator.
     */
    function initiateUnstake(uint256 _value, address _operator) public {
        address owner = operatorToOwner[_operator];
        require(
            msg.sender == _operator ||
            msg.sender == owner, "Only operator or the owner of the stake can initiate unstake.");
        require(_value <= stakeBalances[_operator], "Staker must have enough tokens to unstake.");

        stakeBalances[_operator] = stakeBalances[_operator].sub(_value);
        uint256 createdAt = now;
        withdrawals[_operator] = Withdrawal(withdrawals[_operator].amount.add(_value), createdAt);

        emit InitiatedUnstake(_operator, _value, createdAt);
    }

    /**
     * @notice Finishes unstake of the tokens of provided withdrawal request.
     * You can only finish unstake once withdrawal delay is over for the request,
     * otherwise the function will fail and remaining gas is returned.
     * @param _operator Operator address.
     */
    function finishUnstake(address _operator) public {
        require(now >= withdrawals[_operator].createdAt.add(stakeWithdrawalDelay), "Can not finish unstake before withdrawal delay is over.");
        address owner = operatorToOwner[_operator];

        // No need to call approve since msg.sender will be this staking contract.
        token.safeTransfer(owner, withdrawals[_operator].amount);

        // Cleanup withdrawal record.
        delete withdrawals[_operator];

        // Release operator only when the stake is depleted
        if (stakeBalances[_operator] <= 0) {
            operatorToOwner[_operator] = address(0);
            ownerOperators[owner].removeAddress(_operator);
        }

        emit FinishedUnstake(_operator);
    }

    /**
     * @dev Gets withdrawal request by Operator.
     * @param _operator address of withdrawal request.
     * @return amount The amount the given operator will be able to withdraw
     *                once the withdrawal delay has passed.
     * @return createdAt The initiation time of the withdrawal request for the
     *                   given operator, used to determine when the withdrawal
     *                   delay has passed.
     */
    function getWithdrawal(address _operator) public view returns (uint256 amount, uint256 createdAt) {
        return (withdrawals[_operator].amount, withdrawals[_operator].createdAt);
    }

    /**
     * @dev Slash provided token amount from every member in the misbehaved
     * operators array and burn 100% of all the tokens.
     * @param amount Token amount to slash from every misbehaved operator.
     * @param misbehavedOperators Array of addresses to seize the tokens from.
     */
    function slash(uint256 amount, address[] memory misbehavedOperators) 
        public
        onlyApprovedOperatorContract(msg.sender) {
        for (uint i = 0; i < misbehavedOperators.length; i++) {
            address operator = misbehavedOperators[i];
            require(authorizations[msg.sender][operator], "Not authorized");
            stakeBalances[operator] = stakeBalances[operator].sub(amount);
        }

        token.burn(misbehavedOperators.length.mul(amount));
    }

    /**
     * @dev Seize provided token amount from every member in the misbehaved
     * operators array. The tattletale is rewarded with 5% of the total seized
     * amount scaled by the reward adjustment parameter and the rest 95% is burned.
     * @param amount Token amount to seize from every misbehaved operator.
     * @param rewardMultiplier Reward adjustment in percentage. Min 1% and 100% max.
     * @param tattletale Address to receive the 5% reward.
     * @param misbehavedOperators Array of addresses to seize the tokens from.
     */
    function seize(
        uint256 amount,
        uint256 rewardMultiplier,
        address tattletale,
        address[] memory misbehavedOperators
    ) public onlyApprovedOperatorContract(msg.sender) {
        for (uint i = 0; i < misbehavedOperators.length; i++) {
            address operator = misbehavedOperators[i];
            require(authorizations[msg.sender][operator], "Not authorized");
            stakeBalances[operator] = stakeBalances[operator].sub(amount);
        }

        uint256 total = misbehavedOperators.length.mul(amount);
        uint256 tattletaleReward = (total.mul(5).div(100)).mul(rewardMultiplier).div(100);

        token.transfer(tattletale, tattletaleReward);
        token.burn(total.sub(tattletaleReward));
    }

    /**
     * @dev Authorizes operator contract to access staked token balance of
     * the provided operator. Can only be executed by stake operator authorizer.
     * @param _operator address of stake operator.
     * @param _operatorContract address of operator contract.
     */
    function authorizeOperatorContract(address _operator, address _operatorContract)
        public
        onlyOperatorAuthorizer(_operator)
        onlyApprovedOperatorContract(_operatorContract) {
        authorizations[_operatorContract][_operator] = true;
    }

    /**
     * @dev Checks if operator contract has been authorized for the provided operator.
     * @param _operator address of stake operator.
     * @param _operatorContract address of operator contract.
     * @return Returns True if operator contract has been authorized for the provided operator.
     */
    function isAuthorized(address _operator, address _operatorContract) public view returns (bool) {
        return authorizations[_operatorContract][_operator];
    }
}
