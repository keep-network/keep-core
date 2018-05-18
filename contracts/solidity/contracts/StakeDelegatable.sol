pragma solidity ^0.4.21;

import "zeppelin-solidity/contracts/math/SafeMath.sol";


/**
 * @title Stake Delegatable
 * @dev A contract that allows delegate your stake balance to any address
 * that is not already a staker.
 */
contract StakeDelegatable {

    mapping(address => uint256) public stakeBalances;
    mapping(address => address) public delegatorFor;
    mapping(address => address) public operatorFor;

    /**
     * @dev Gets the stake balance of the specified address.
     * @param _staker The address to query the balance of.
     * @return The amount staked by the passed address.
     */
    function stakeBalanceOf(address _staker) public view returns (uint256) {
        address delegator = delegatorFor[_staker];
        if (delegator != address(0)) {
            return stakeBalances[delegator];
        }
        return stakeBalances[_staker];
    }

    /**
     * @dev Delegates your stake balance to a specified address.
     * @param _operator Address to where you want to delegate your balance.
     */
    function delegateStakeTo(address _operator) public {
        require(_operator != address(0));

        // Operator must not be a staker.
        require(stakeBalances[_operator] == 0);

        // Revert if operator address was already used.
        address previousDelegator = delegatorFor[_operator];
        require(previousDelegator == address(0));

        // Release previous operator address when you delegate stake to a new one.
        address previousOperator = operatorFor[msg.sender];
        if (previousOperator != address(0)) {
            delete delegatorFor[previousOperator];
        }

        operatorFor[msg.sender] = _operator;
        delegatorFor[_operator] = msg.sender;
    }

    /**
     * @dev Removes delegate for your stake balance.
     */
    function removeDelegate() public {
        address operator = operatorFor[msg.sender];
        delete delegatorFor[operator];
        delete operatorFor[msg.sender];
    }

    /**
     * @dev Removes delegate for the address if it's an operator and staked.
     * @param _address The address to check.
     */
    function removeDelegateIfStakedAsOperator(address _address) internal {
        address delegator = delegatorFor[_address];
        if (delegator != address(0)) {
            delete delegatorFor[_address];
            delete operatorFor[delegator];
        }
    }
}
