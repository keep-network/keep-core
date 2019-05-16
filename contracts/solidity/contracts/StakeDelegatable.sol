pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/cryptography/ECDSA.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";
import "./StakingProxy.sol";


/**
 * @title Stake Delegatable
 * @dev A base contract to allow stake delegation for staking contracts.
 */
contract StakeDelegatable {
    using SafeMath for uint256;
    using SafeERC20 for ERC20;
    using BytesLib for bytes;
    using ECDSA for bytes32;

    ERC20 public token;
    StakingProxy public stakingProxy;

    uint256 public stakeWithdrawalDelay;

    // Stake balances.
    mapping(address => uint256) public stakeBalances;

    // Stake balances.
    mapping(address => uint256) public initialStakeBalances;

    // Stake delegation mappings.
    mapping(address => address) public operatorToOwner;
    mapping(address => address) public magpieToOwner;

    // List of operators for the stake owner.
    mapping(address => address[]) public ownerOperators;

    /**
     * @dev Gets the stake balance of the specified address.
     * @param _address The address to query the balance of.
     * @return An uint256 representing the amount staked by the passed address.
     */
    function stakeBalanceOf(address _address) public view returns (uint256 balance) {
        return stakeBalances[_address];
    }

    /**
     * @dev Gets the list of operators of the specified address.
     * @return An array of addresses.
     */
    function operatorsOf(address _address) public view returns (address[] memory) {
        return ownerOperators[_address];
    }
}
