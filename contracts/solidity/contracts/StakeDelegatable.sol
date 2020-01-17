pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";
import "./utils/AddressArrayUtils.sol";


/**
 * @title Stake Delegatable
 * @dev A base contract to allow stake delegation for staking contracts.
 */
contract StakeDelegatable {
    using SafeMath for uint256;
    using SafeERC20 for ERC20;
    using BytesLib for bytes;
    using AddressArrayUtils for address[];

    ERC20 public token;

    uint256 public stakeWithdrawalDelay;

    // Stake balances.
    mapping(address => uint256) public stakeBalances;

    // Stake delegation mappings.
    mapping(address => address) public operatorToOwner;
    mapping(address => address payable) public operatorToMagpie;

    // List of operators for the stake owner.
    mapping(address => address[]) public ownerOperators;

    /**
     * @dev Gets the stake balance of the specified address.
     * @param _address The address to query the balance of.
     * @return An uint256 representing the amount staked by the passed address.
     */
    function balanceOf(address _address) public view returns (uint256 balance) {
        return stakeBalances[_address];
    }

    /**
     * @dev Gets the list of operators of the specified address.
     * @return An array of addresses.
     */
    function operatorsOf(address _address) public view returns (address[] memory) {
        return ownerOperators[_address];
    }

    /**
     * @dev Gets the stake owner for the specified operator address.
     * @return Stake owner address.
     */
    function ownerOf(address _operator) public view returns (address) {
        return operatorToOwner[_operator];
    }

    /**
     * @dev Gets the magpie for the specified operator address.
     * @return Magpie address.
     */
    function magpieOf(address _operator) public view returns (address payable) {
        return operatorToMagpie[_operator];
    }

    /**
     * @dev Gets the array of magpie addresses for the specified array of operator addresses.
     * @return Array of magpie addresses.
     */
    function magpiesOf(address[] memory _operators) public view returns (address payable[] memory magpies) {
        magpies = new address payable[](_operators.length);
        for (uint i = 0; i < _operators.length; i++) {
            magpies[i] = operatorToMagpie[_operators[i]];
        }
    }
}
