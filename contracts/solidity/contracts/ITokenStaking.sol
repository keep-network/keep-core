pragma solidity ^0.5.4;

/**
 * @title Keep Network Token Staking
 *
 * @notice Provides an information about eligible stake of network operators.
 * The Keep network uses staking of tokens to enforce correct behavior.
 * Anyone with tokens can stake them, setting them aside as collateral for
 * network operations. Staked tokens are delegated to an operator address who
 * performs work for operator contracts. Operators can earn rewards from
 * contributing to the network, but if they misbehave their collateral can be
 * taken away (stake slashing) as punishment.
 */
interface ITokenStaking {

    /**
     * @dev Gets the eligible stake balance of the specified operator.
     * An eligible stake is a stake that passed the initialization period
     * and is not currently undelegating. Also, the operator had to approve
     * the specified operator contract.
     *
     * Operator with a minimum required amount of eligible stake can join the
     * network and participate in new work selection.
     *
     * @param _operator Address of stake operator.
     * @param _operatorContract Address of operator contract.
     * @return Eligible stake balance.
     */
    function eligibleStake(
        address _operator,
        address _operatorContract
    ) external view returns (uint256 balance);
}