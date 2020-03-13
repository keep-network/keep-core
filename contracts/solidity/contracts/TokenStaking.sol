pragma solidity ^0.5.4;

import "./StakeDelegatable.sol";
import "./utils/UintArrayUtils.sol";
import "./Registry.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";


/**
 * @title TokenStaking
 * @dev A token staking contract for a specified standard ERC20Burnable token.
 * A holder of the specified token can stake delegate its tokens to this contract
 * and recover the stake after undelegation period is over.
 */
contract TokenStaking is StakeDelegatable {

    using UintArrayUtils for uint256[];
    using SafeERC20 for ERC20Burnable;

    event Staked(address indexed from, uint256 value);
    event Undelegated(address indexed operator, uint256 undelegatedAt);
    event RecoveredStake(address operator, uint256 recoveredAt);

    // Registry contract with a list of approved operator contracts and upgraders.
    Registry public registry;

    // Authorized operator contracts.
    mapping(address => mapping (address => bool)) internal authorizations;

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
     * @param _initializationPeriod To avoid certain attacks on work selection, recently created
     * operators must wait for a specific period of time before being eligible for work selection.
     * @param _undelegationPeriod The staking contract guarantees that an undelegated operator’s
     * stakes will stay locked for a number of blocks after undelegation, and thus available as
     * collateral for any work the operator is engaged in.
     */
    constructor(
        address _tokenAddress,
        address _registry,
        uint256 _initializationPeriod,
        uint256 _undelegationPeriod
    ) public {
        require(_tokenAddress != address(0x0), "Token address can't be zero.");
        token = ERC20Burnable(_tokenAddress);
        registry = Registry(_registry);
        initializationPeriod = _initializationPeriod;
        undelegationPeriod = _undelegationPeriod;
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
        require(operators[operator].owner == address(0), "Operator address is already in use.");
        address authorizer = _extraData.toAddress(40);

        // Transfer tokens to this contract.
        token.safeTransferFrom(_from, address(this), _value);

        operators[operator] = Operator(
            OperatorParams.pack(_value, block.number, 0),
            _from,
            magpie,
            authorizer
        );
        ownerOperators[_from].push(operator);

        emit Staked(operator, _value);
    }

    /**
     * @notice Cancels stake of tokens within the operator initialization period
     * without being subjected to the token lockup for the undelegation period.
     * This can be used to undo mistaken delegation to the wrong operator address.
     * @param _operator Address of the stake operator.
     */
    function cancelStake(address _operator) public {
        address owner = operators[_operator].owner;
        require(
            msg.sender == _operator ||
            msg.sender == owner, "Only operator or the owner of the stake can cancel the delegation."
        );
        uint256 operatorParams = operators[_operator].packedParams;

        require(
            block.number <= operatorParams.getCreationBlock().add(initializationPeriod),
            "Initialization period is over"
        );

        uint256 amount = operatorParams.getAmount();
        operators[_operator].packedParams = operatorParams.setAmount(0);

        token.safeTransfer(owner, amount);
    }

    /**
     * @notice Undelegates staked tokens. You will be able to recover your stake by calling
     * `recoverStake()` with operator address once undelegation period is over.
     * @param _operator Address of the stake operator.
     */
    function undelegate(address _operator) public {
        address owner = operators[_operator].owner;
        require(
            msg.sender == _operator ||
            msg.sender == owner, "Only operator or the owner of the stake can undelegate."
        );
        uint256 oldParams = operators[_operator].packedParams;
        uint256 newParams = oldParams.setUndelegationBlock(block.number);
        operators[_operator].packedParams = newParams;
        emit Undelegated(_operator, block.number);
    }

    /**
     * @notice Recovers staked tokens and transfers them back to the owner. Recovering
     * tokens can only be performed when the operator is finished undelegating.
     * @param _operator Operator address.
     */
    function recoverStake(address _operator) public {
        uint256 operatorParams = operators[_operator].packedParams;
        require(
            block.number > operatorParams.getUndelegationBlock().add(undelegationPeriod),
            "Can not recover stake before undelegation period is over."
        );
        address owner = operators[_operator].owner;
        uint256 amount = operatorParams.getAmount();

        operators[_operator].packedParams = operatorParams.setAmount(0);

        token.safeTransfer(owner, amount);
        emit RecoveredStake(_operator, block.number);
    }

    /**
     * @dev Gets stake delegation info for the given operator.
     * @param _operator Operator address.
     * @return amount The amount of tokens the given operator delegated.
     * @return createdAt The time when the stake has been delegated.
     * @return undelegatedAt The time when undelegation has been requested.
     * If undelegation has not been requested, 0 is returned.
     */
    function getDelegationInfo(address _operator)
    public view returns (uint256 amount, uint256 createdAt, uint256 undelegatedAt) {
        return operators[_operator].packedParams.unpack();
    }

    /**
     * @dev Slash provided token amount from every member in the misbehaved
     * operators array and burn 100% of all the tokens.
     * @param amountToSlash Token amount to slash from every misbehaved operator.
     * @param misbehavedOperators Array of addresses to seize the tokens from.
     */
    function slash(uint256 amountToSlash, address[] memory misbehavedOperators)
        public
        onlyApprovedOperatorContract(msg.sender) {

        uint256 totalAmountToBurn = 0;
        for (uint i = 0; i < misbehavedOperators.length; i++) {
            address operator = misbehavedOperators[i];
            require(authorizations[msg.sender][operator], "Not authorized");

            uint256 operatorParams = operators[operator].packedParams;
            require(
                block.number > operatorParams.getCreationBlock().add(initializationPeriod),
                "Operator stakes must be active"
            );

            uint256 currentAmount = operatorParams.getAmount();

            if (currentAmount < amountToSlash) {
                totalAmountToBurn = totalAmountToBurn.add(currentAmount);

                uint256 newAmount = 0;
                operators[operator].packedParams = operatorParams.setAmount(newAmount);
            } else {
                totalAmountToBurn = totalAmountToBurn.add(amountToSlash);

                uint256 newAmount = currentAmount.sub(amountToSlash);
                operators[operator].packedParams = operatorParams.setAmount(newAmount);
            }
        }

        token.burn(totalAmountToBurn);
    }

    /**
     * @dev Seize provided token amount from every member in the misbehaved
     * operators array. The tattletale is rewarded with 5% of the total seized
     * amount scaled by the reward adjustment parameter and the rest 95% is burned.
     * @param amountToSeize Token amount to seize from every misbehaved operator.
     * @param rewardMultiplier Reward adjustment in percentage. Min 1% and 100% max.
     * @param tattletale Address to receive the 5% reward.
     * @param misbehavedOperators Array of addresses to seize the tokens from.
     */
    function seize(
        uint256 amountToSeize,
        uint256 rewardMultiplier,
        address tattletale,
        address[] memory misbehavedOperators
    ) public onlyApprovedOperatorContract(msg.sender) {
        uint256 totalAmountToBurn = 0;
        for (uint i = 0; i < misbehavedOperators.length; i++) {
            address operator = misbehavedOperators[i];
            uint256 operatorParams = operators[operator].packedParams;

            require(authorizations[msg.sender][operator], "Not authorized");
            require(
                block.number > operatorParams.getCreationBlock().add(initializationPeriod),
                "Operator stakes must be active"
            );

            uint256 currentAmount = operatorParams.getAmount();

            if (currentAmount < amountToSeize) {
                totalAmountToBurn = totalAmountToBurn.add(currentAmount);

                uint256 newAmount = 0;
                operators[operator].packedParams = operatorParams.setAmount(newAmount);
            } else {
                totalAmountToBurn = totalAmountToBurn.add(amountToSeize);

                uint256 newAmount = currentAmount.sub(amountToSeize);
                operators[operator].packedParams = operatorParams.setAmount(newAmount);
            }
        }

        uint256 tattletaleReward = (totalAmountToBurn.mul(5).div(100)).mul(rewardMultiplier).div(100);

        token.safeTransfer(tattletale, tattletaleReward);
        token.burn(totalAmountToBurn.sub(tattletaleReward));
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
     * @dev Checks if operator contract has access to the staked token balance of
     * the provided operator.
     * @param _operator address of stake operator.
     * @param _operatorContract address of operator contract.
     */
    function isAuthorizedForOperator(address _operator, address _operatorContract) public view returns (bool) {
        return authorizations[_operatorContract][_operator];
    }

    /**
     * @dev Gets the eligible stake balance of the specified address.
     * An eligible stake is a stake that passed the initialization period
     * and is not currently undelegating. Also, the operator had to approve
     * the specified operator contract.
     *
     * Operator with a minimum required amount of eligible stake can join the
     * network and participate in new work selection.
     *
     * @param _operator address of stake operator.
     * @param _operatorContract address of operator contract.
     * @return an uint256 representing the eligible stake balance.
     */
    function eligibleStake(
        address _operator,
        address _operatorContract
    ) public view returns (uint256 balance) {
        bool isAuthorized = authorizations[_operatorContract][_operator];

        uint256 operatorParams = operators[_operator].packedParams;
        uint256 createdAt = operatorParams.getCreationBlock();
        uint256 undelegatedAt = operatorParams.getUndelegationBlock();

        bool isActive = block.number > createdAt.add(initializationPeriod);
        bool isUndelegating = (undelegatedAt > 0) && (block.number > undelegatedAt);

        if (isAuthorized && isActive && !isUndelegating) {
            balance = operatorParams.getAmount();
        }
    }

    /**
     * @dev Gets the active stake balance of the specified address.
     * An active stake is a stake that passed the initialization period.
     * Also, the operator had to approve the specified operator contract.
     *
     * The difference between eligible stake is that active stake does not make
     * the operator eligible for work selection but it may be still finishing
     * earlier work during undelegation period. Operator with a minimum required
     * amount of active stake can join the network but cannot be selected to any
     * new work.
     *
     * @param _operator address of stake operator.
     * @param _operatorContract address of operator contract.
     * @return an uint256 representing the eligible stake balance.
     */
    function activeStake(
        address _operator,
        address _operatorContract
    ) public view returns (uint256 balance) {
        bool isAuthorized = authorizations[_operatorContract][_operator];

        uint256 operatorParams = operators[_operator].packedParams;
        uint256 createdAt = operatorParams.getCreationBlock();

        bool isActive = block.number > createdAt.add(initializationPeriod);

        if (isAuthorized && isActive) {
            balance = operatorParams.getAmount();
        }
    }
}
