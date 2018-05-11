pragma solidity ^0.4.21;

import "./StakingProxy.sol";


/**
 * @title Staking Delegate Contract.
 * @dev An optional contract for staking proxy to allow staker
 * to delegate it's balance to an arbitrary operator address.
 */
contract StakingDelegate {

    /**
     * @dev Only authorized contract can invoke functions with this modifier.
     */
    modifier onlyAuthorized {
        require(msg.sender == stakingProxy);
        _;
    }

    StakingProxy public stakingProxy;

    mapping(address => address) internal delegatorFor;
    mapping(address => address) internal operatorFor;

    /**
     * @dev Creates a staking delegate contract.
     * @param _stakingProxy Address of a staking proxy that will be linked to this contract.
     */
    function StakingDelegate(address _stakingProxy) public {
        require(_stakingProxy != address(0x0));
        stakingProxy = StakingProxy(_stakingProxy);
    }

    /**
     * @dev Gets delegated staking balance of address.
     * @param _address The address to query the delegated staking balance of.
     */
    function delegatedBalanceOf(address _address)
        public
        view
        onlyAuthorized
        returns (uint256)
    {

        // Get actual stake balance for the address.
        uint256 balance = stakingProxy.totalBalanceOf(_address);

        // If the provided address is an operator then get it's delegator stake balance.
        address delegator = delegatorFor[_address];
        if (balance == 0 && delegator != address(0)) {
            balance = stakingProxy.totalBalanceOf(delegator);
        }

        // If the provided address is a delegator we assume it has no stake balance.
        address operator = operatorFor[_address];
        if (operator != address(0)) {
            return 0;
        }

        return balance;
    }

    /**
     * @dev Delegates your stake balance to a specified address.
     * @param _operator Address to where you you want to delegate your balance.
     */
    function delegateStakeTo(address _operator) public {
        require(_operator != address(0));

        // Operator must not be a staker.
        require(stakingProxy.totalBalanceOf(_operator) == 0);

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
     * @dev Remove delegate for your stake balance.
     */
    function removeDelegate() public {
        address operator = operatorFor[msg.sender];
        delete delegatorFor[operator];
        delete operatorFor[msg.sender];
    }
}
