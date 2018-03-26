pragma solidity ^0.4.18;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";

interface authorizedStakingContract {
    function stakeBalanceOf(address addr) public returns (uint256);
}


/**
 * @title Staking Proxy Contract
 * @dev An ownable staking proxy contract to provide upgradable staking.
 * Upgraded contract can be authorized as an active contract by the owner.
 * The old contracts are put into legacyContracts array.
 * The staking contracts must call emit events on this contract.
 */
contract StakingProxy is Ownable {
    /** 
     * @dev Only authorized contracts can invoke functions with this modifier.
     */
    modifier onlyAuthorized {
        require(isAuthorized(msg.sender));
        _;
    }

    address public activeContract;
    address[] public legacyContracts;

    event Staked(address indexed user, uint256 amount);
    event Unstaked(address indexed user, uint256 amount);

    /**
     * @dev Gets the sum of all staking balances of the specified staker address.
     * @param _staker The address to query the balance of.
     * @return An uint256 representing the amount staked by the passed address.
     */
    function balanceOf(address _staker)
        public
        constant
        returns (uint256 _balance)
    {
        uint256 balance = authorizedStakingContract(activeContract).stakeBalanceOf(_staker);
        for (uint i = 0; i < legacyContracts.length; i++) {
            balance = balance + authorizedStakingContract(legacyContracts[i]).stakeBalanceOf(_staker);
        }
        return balance;
    }

    /**
     * @dev Update active contract.
     * @param _contract The address of a staking contract.
     */
    function updateActiveContract(address _contract) 
        public
        onlyOwner
    {
        require(_contract != address(0));
        require(_contract != activeContract);
        legacyContracts.push(activeContract);
        activeContract = _contract;
    }

    /**
     * @dev Emit staked event.
     * @param _staker The address of the staker.
     * @param _amount The staked amount.
     */
    function emitStakedEvent(address _staker, uint256 _amount)
        public
        onlyAuthorized
    {
        Staked(_staker, _amount);
    }

    /**
     * @dev Emit unstaked event.
     * @param _staker The address of the staker.
     * @param _amount The unstaked amount.
     */
    function emitUnstakedEvent(address _staker, uint256 _amount)
        public
        onlyAuthorized
    {
        Unstaked(_staker, _amount);
    }

    /**
     * @dev Check if a staking contract is authorized to work with this contract.
     * @param _address The address of a staking contract.
     * @return A bool wether it's authorized.
     */
    function isAuthorized(address _address) 
        public
        returns (bool) 
    {
        if (_address == activeContract) {
            return true;
        }
        for (uint i = 0; i < legacyContracts.length; i++) {
            if (legacyContracts[i] == _address) {
                return true;
            }
        }
        return false;
    }
}
