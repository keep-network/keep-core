pragma solidity ^0.4.18;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";
import "./Utils/AddressArrayUtils.sol";

interface authorizedStakingContract {
    function stakeBalanceOf(address addr) public constant returns (uint256);
}


/**
 * @title Staking Proxy Contract
 * @dev An ownable staking proxy contract to provide upgradable staking.
 * Upgraded contracts are added to authorizedContracts list.
 * The staking contracts must call emit events on this contract.
 */
contract StakingProxy is Ownable {

    using AddressArrayUtils for address[];

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
        returns (uint256)
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

        // Must not be already authorized
        require(!isAuthorized(_contract));

        // Must not be deauthorized
        require(!isDeauthorized(_contract));

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

        // Must be authorized previously
        require(isAuthorized(_contract));

        // Must not be already deauthorized
        require(!isDeauthorized(_contract));

        // Find and remove contract address from authorizedContracts array.
        authorizedContracts.removeAddress(_contract);

        AuthorizedContractRemoved(_contract);
        deauthorizedContracts.push(_contract);
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
     * @return A bool whether it's authorized.
     */
    function isAuthorized(address _address) 
        public
        constant
        returns (bool) 
    {
        return authorizedContracts.isFound(_address);
    }

    /**
     * @dev Check if a staking contract is deauthorized.
     * @param _address The address of a staking contract.
     * @return A bool whether it's deauthorized.
     */
    function isDeauthorized(address _address) 
        public
        constant
        returns (bool) 
    {
        return deauthorizedContracts.isFound(_address);
    }
}
