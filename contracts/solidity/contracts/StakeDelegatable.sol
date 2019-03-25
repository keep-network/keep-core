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

    // Stake delegation mappings.
    mapping(address => address) public operatorToOwner;
    mapping(address => address) public magpieToOwner;

    /**
     * @dev Gets the stake balance of the specified address.
     * @param _address The address to query the balance of.
     * @return An uint256 representing the amount staked by the passed address.
     */
    function stakeBalanceOf(address _address) public view returns (uint256 balance) {
        return stakeBalances[_address];
    }
}
