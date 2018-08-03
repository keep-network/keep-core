pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";


/**
 * @title Stake Delegatable
 * @dev A contract that allows delegate your stake balance to any address
 * that is not already a staker. Delegator refers to a staker who has
 * delegated their stake to another address, while delegate refers to an
 * address that has had a stake delegated to it.
 */
contract StakeDelegatable {

    /**
     * @dev Only not null addresses can be passed to the functions with this modifier.
     */
    modifier notNull(address _address) {
        require(_address != address(0), "Provided address can not be zero.");
        _;
    }

    /**
     * @dev Only non staker addresses can be passed to the functions with this modifier.
     */
    modifier notStaker(address _address) {
        require(stakeBalances[_address] == 0, "Provided address is not a staker.");
        _;
    }

    mapping(address => uint256) public stakeBalances;
    mapping(address => address) public delegatorToDelegate;
    mapping(address => address) public delegateToDelegator;

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
        // If provided address is a delegate return its delegator balance.
        address delegator = delegatorToDelegate[_address];
        if (delegator != address(0) && delegateToDelegator[delegator] == _address) {
            return stakeBalances[delegator];
        }

        // If provided address is a delegator return zero balance since the 
        // balance is delegated to a delegate.
        address delegate = delegateToDelegator[_address];
        if (delegate != address(0) && delegatorToDelegate[delegate] == _address) {
            return 0;
        }

        return stakeBalances[_address];
    }

    /**
     * @dev Returns address of a delegate if it exists for the
     * provided address or the provided address otherwise.
     * @param _address The address to check.
     * @return Delegate address or provided address.
     */
    function getDelegatorOrDelegate(address _address)
        public
        view
        notNull(_address)
        returns (address)
    {
        address delegate = delegateToDelegator[_address];
        if (delegate != address(0) && delegatorToDelegate[delegate] == _address) {
            return delegate;
        }
        return _address;
    }

    /**
     * @dev Approves address to delegate your stake balance. You can only
     * have one delegate address. Delegate must also request to operate
     * your stake by calling requestDelegateFor() method.
     * @param _address Address to where you want to delegate your balance.
     */
    function approveDelegateAt(address _address)
        public
        notNull(_address)
        notStaker(_address)
    {
        delegateToDelegator[msg.sender] = _address;
    }

    /**
     * @dev Requests to delegate stake for a specified address.
     * An address must approve you first to delegate by calling
     * requestDelegateFor() method.
     * @param _address Address for which you request to delegate.
     */
    function requestDelegateFor(address _address)
        public
        notNull(_address)
    {
        delegatorToDelegate[msg.sender] = _address;
    }

    /**
     * @dev Removes delegate for your stake balance.
     */
    function removeDelegate() public {
        address delegate = delegateToDelegator[msg.sender];
        delete delegatorToDelegate[delegate];
        delete delegateToDelegator[msg.sender];
    }

    /**
     * @dev Revert if a delegate try to stake.
     * @param _address The address to check.
     */
    function revertIfDelegateStakes(address _address) internal {
        address delegator = delegatorToDelegate[_address];
        if (delegator != address(0)) {
            revert("Provided address can not stake since it has stake delegated to it.");
        }
    }
}
