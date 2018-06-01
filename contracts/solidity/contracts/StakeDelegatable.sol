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

    /**
     * @dev Only not null addresses can be passed to the functions with this modifier.
     */
    modifier notNull(address _address) {
        require(_address != address(0));
        _;
    }

    /**
     * @dev Only non staker addresses can be passed to the functions with this modifier.
     */
    modifier notStaker(address _address) {
        require(stakeBalances[_address] == 0);
        _;
    }

    mapping(address => uint256) public stakeBalances;
    mapping(address => address) public stakerToOperator;
    mapping(address => address) public operatorToStaker;

    /**
     * @dev Gets the stake balance of the specified address.
     * @param _address The address to query the balance of.
     * @return The amount staked by the passed address.
     */
    function stakeBalanceOf(address _address)
        public
        view
        notNull(_address)
        returns (uint256)
    {
        address delegator = stakerToOperator[_address];
        if (delegator != address(0) && operatorToStaker[delegator] == _address) {
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
    function getStakerOrOperator(address _address)
        public
        view
        notNull(_address)
        returns (address)
    {
        address operator = operatorToStaker[_address];
        if (operator != address(0) && stakerToOperator[operator] == _address) {
            return operator;
        }
        return _address;
    }

    /**
     * @dev Approves address to operate your stake balance. You can only
     * have one operator address. Operator must also request to operate
     * your stake by calling requestOperateFor() method.
     * @param _address Address to where you want to delegate your balance.
     */
    function approveOperatorAt(address _address)
        public
        notNull(_address)
        notStaker(_address)
    {
        operatorToStaker[msg.sender] = _address;
    }

    /**
     * @dev Requests to operate stake for a specified address.
     * Staker address must approve you to operate by calling
     * approveOperatorAt() method.
     * @param _address Address for which you request to operate.
     */
    function requestOperateFor(address _address)
        public
        notNull(_address)
    {
        stakerToOperator[msg.sender] = _address;
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
