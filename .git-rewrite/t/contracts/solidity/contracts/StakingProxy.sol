pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./utils/AddressArrayUtils.sol";

interface authorizedStakingContract {
    function stakeBalanceOf(address addr) external view returns (uint256);
}


/**
 * @title Staking Proxy Contract
 * @dev An ownable staking proxy contract to provide upgradable staking.
 * Upgraded contracts are added to authorizedContracts list. The staking
 * contracts must call "emitStakedEvent" and "emitUnstakedEvent" functions on
 * this contract.
 */
contract StakingProxy is Ownable {

    using AddressArrayUtils for address[];

    /**
     * @dev Only authorized contracts can invoke functions with this modifier.
     */
    modifier onlyAuthorized {
        require(isAuthorized(msg.sender), "Sender is not authorized.");
        _;
    }

    address[] public authorizedContracts;
    address[] public deauthorizedContracts;

    event Staked(address indexed staker, uint256 amount);
    event Unstaked(address indexed staker, uint256 amount);
    event AuthorizedContractAdded(address indexed contractAddress);
    event AuthorizedContractRemoved(address indexed contractAddress);

    /**
     * @dev Gets the sum of all staking balances of the specified staker address.
     * @param _staker The address to query the balance of.
     * @return An uint256 representing the amount staked by the passed address.
     */
    function balanceOf(address _staker)
        public
        view
        returns (uint256)
    {
        require(_staker != address(0), "Staker address can't be zero.");
        uint256 balance = 0;
        for (uint i = 0; i < authorizedContracts.length; i++) {
            balance = balance + authorizedStakingContract(authorizedContracts[i]).stakeBalanceOf(_staker);
        }
        return balance;
    }

    /**
     * @dev Authorize contract. Owner of this proxy can add a staking contract
     * into the authorized list and added staking contract will be accounted
     * for the total staker's balance and corresponding stake/unstake events.
     * @param _contract The address of a staking contract.
     */
    function authorizeContract(address _contract) 
        public
        onlyOwner
    {
        require(_contract != address(0), "Contract address can't be zero.");
        require(!isAuthorized(_contract), "Contract is already authorized.");
        require(!isDeauthorized(_contract), "Contract was deauthorized.");

        authorizedContracts.push(_contract);
        emit AuthorizedContractAdded(_contract);
    }

    /**
     * @dev Deauthorize contract. Owner of this proxy can remove a staking
     * contract from the authorized list and removed staking contract will be
     * excluded from the total staker's balance and corresponding stake/unstake
     * events are not going to be broadcasted. Removed contract is also added to
     * the deauthorized list for easier tracking of legacy contracts and
     * to prevent repeated authorization of a legacy contract.
     * @param _contract The address of a staking contract.
     */
    function deauthorizeContract(address _contract) 
        public
        onlyOwner
    {
        require(_contract != address(0), "Contract address can't be zero.");
        require(isAuthorized(_contract), "Contract is already authorized.");
        require(!isDeauthorized(_contract), "Contract was deauthorized.");

        authorizedContracts.removeAddress(_contract);
        deauthorizedContracts.push(_contract);

        emit AuthorizedContractRemoved(_contract);
    }

    /**
     * @dev Emit staked event. This function is called by every authorized
     * staking contract where staking occurs so the network clients can have
     * a single point to listen to the events across multiple staking contracts.
     * @param _staker The address of the staker.
     * @param _amount The staked amount.
     */
    function emitStakedEvent(address _staker, uint256 _amount)
        public
        onlyAuthorized
    {
        emit Staked(_staker, _amount);
    }

    /**
     * @dev Emit unstaked event. This function is called by every authorized
     * staking contract where unstaking occurs so the network clients can have
     * a single point to listen to the events across multiple staking contracts.
     * @param _staker The address of the staker.
     * @param _amount The unstaked amount.
     */
    function emitUnstakedEvent(address _staker, uint256 _amount)
        public
        onlyAuthorized
    {
        emit Unstaked(_staker, _amount);
    }

    /**
     * @dev Check if a staking contract is authorized to work with this
     * contract otherwise it's not allowed to call "emit events" methods on this
     * contract and it's balance is not inlcuded into the total staking balance.
     * @param _address The address of a staking contract.
     * @return True if staking contract is authorized, false otherwise.
     */
    function isAuthorized(address _address)
        public
        view
        returns (bool)
    {
        return authorizedContracts.contains(_address);
    }

    /**
     * @dev Check if a staking contract is deauthorized. If it's deauthorized
     * it won't be possible to authorize it again and as a non authorized
     * contract it's not allowed to call "emit events" methods on this contract
     * and it's balance is not inlcuded into the total staking balance.
     * @param _address The address of a staking contract.
     * @return True if staking contract is deauthorized, false otherwise.
     */
    function isDeauthorized(address _address)
        public
        view
        returns (bool)
    {
        return deauthorizedContracts.contains(_address);
    }
}
