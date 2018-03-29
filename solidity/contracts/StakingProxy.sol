pragma solidity ^0.4.18;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";

interface authorizedStakingContract {
    function stakeBalanceOf(address addr) public returns (uint256);
}


/**
 * @title Staking Proxy Contract
 * @dev An ownable staking proxy contract to provide upgradable staking.
 * Upgraded contracts are added to authorizedContracts list.
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

    address[] public authorizedContracts;
    address[] public deauthorizedContracts;

    event Staked(address indexed user, uint256 amount);
    event Unstaked(address indexed user, uint256 amount);
    event AuthorizedContractAdded(address indexed contractAddress);
    event AuthorizedContractRemoved(address indexed contractAddress);

    /**
     * @dev Gets the sum of all staking balances of the specified staker address.
     * @param _staker The address to query the balance of.
     * @return balance An uint256 representing the amount staked by the passed address.
     */
    function balanceOf(address _staker)
        public
        constant
        returns (uint256 balance)
    {
        uint256 balance = 0;
        for (uint i = 0; i < authorizedContracts.length; i++) {
            balance = balance + authorizedStakingContract(authorizedContracts[i]).stakeBalanceOf(_staker);
        }
        return balance;
    }

    /**
     * @dev Authorize contract.
     * @param _contract The address of a staking contract.
     */
    function authorizeContract(address _contract) 
        public
        onlyOwner
    {
        require(_contract != address(0));
        require(!isAuthorized(_contract));
        authorizedContracts.push(_contract);
        AuthorizedContractAdded(_contract);
    }

    /**
     * @dev Deauthorize contract.
     * @param _contract The address of a staking contract.
     */
    function deauthorizeContract(address _contract) 
        public
        onlyOwner
    {
        require(_contract != address(0));

        // Find and remove contract address from authorizedContracts array.
        for (uint i = 0; i < authorizedContracts.length; i++) {
            // If contract is found in array.
            if (_contract == authorizedContracts[i]) {
                // Delete element at index and shift array.
                for (uint j = i; j < authorizedContracts.length-1; j++) {
                    authorizedContracts[j] = authorizedContracts[j+1];
                }
                delete authorizedContracts[authorizedContracts.length-1];
                authorizedContracts.length--;
                deauthorizedContracts.push(_contract);
                AuthorizedContractRemoved(_contract);
            }
        }
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
        for (uint i = 0; i < contracts.length; i++) {
            if (contracts[i] == _address) {
                return true;
            }
        }
        return false;
    }
}
