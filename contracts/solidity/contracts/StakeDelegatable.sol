pragma solidity ^0.4.21;

import "zeppelin-solidity/contracts/math/SafeMath.sol";


/**
 * @title Stake Delegatable
 * @dev A contract that allows delegate your stake balance to any address
 * that is not already a staker. Delegator refers to a staker who has
 * delegated their stake to another address, while operator refers to an
 * address that has had a stake delegated to it.
 */
contract StakeDelegatable {

    mapping(address => uint256) public stakeBalances;
    mapping(address => address) public stakerToOperator;
    mapping(address => address) public operatorToStaker;

    /**
     * @dev Gets the stake balance of the specified address.
     * @param _address The address to query the balance of.
     * @return The amount staked by the passed address.
     */
    function stakeBalanceOf(address _address) public view returns (uint256) {
        require(_address != address(0));

        address delegator = stakerToOperator[_address];
        if (delegator != address(0)) {
            return stakeBalances[delegator];
        }
        return stakeBalances[_address];
    }

    /**
     * @dev Returns address of an operator if it exists for the
     * provided staker address or the provided staker address otherwise.
     * @param _address The address to check.
     * @return Operator address or provided staker address.
     */
    function getStakerOrOperator(address _address) public view returns (address) {
        require(_address != address(0));

        address operator = operatorToStaker[_address];
        if (operator != address(0)) {
            return operator;
        }
        return _address;
    }

    /**
     * @dev Delegates your stake balance to a specified address.
     * An address can only have one operator address. You can delegate
     * stake to any ethereum address as long as it isn't currently staking
     * or operating someone else's stake. In the current implementation
     * delegating stake doesn't hide the stake balance on your main stake
     * address.
     * @param _operator Address to where you want to delegate your balance.
     */
    function delegateStakeTo(address _operator) public {
        require(_operator != address(0));

        // Operator must not be a staker.
        require(stakeBalances[_operator] == 0);

        // Revert if operator address was already used.
        address previousDelegator = stakerToOperator[_operator];
        require(previousDelegator == address(0));

        // Release previous operator address when you delegate stake to a new one.
        address previousOperator = operatorToStaker[msg.sender];
        if (previousOperator != address(0)) {
            delete stakerToOperator[previousOperator];
        }

        operatorToStaker[msg.sender] = _operator;
        stakerToOperator[_operator] = msg.sender;
    }

    /**
     * @dev Removes delegate for your stake balance.
     */
    function removeDelegate() public {
        address operator = operatorToStaker[msg.sender];
        delete stakerToOperator[operator];
        delete operatorToStaker[msg.sender];
    }

    /**
     * @dev Removes delegate for the address if it's an operator and staked.
     * @param _address The address to check.
     */
    function removeDelegateIfStakedAsOperator(address _address) internal {
        address delegator = stakerToOperator[_address];
        if (delegator != address(0)) {
            delete stakerToOperator[_address];
            delete operatorToStaker[delegator];
        }
    }
}
