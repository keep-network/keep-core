pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20Burnable.sol";
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
    using SafeERC20 for ERC20Burnable;
    using BytesLib for bytes;
    using AddressArrayUtils for address[];

    ERC20Burnable public token;

    uint256 public initializationPeriod;
    uint256 public undelegationPeriod;

    // List of operators for the stake owner.
    mapping(address => address[]) public ownerOperators;

    struct Operator {
        uint256 amount;
        uint256 createdAt;
        uint256 undelegatedAt;
        address owner;
        address payable beneficiary;
        address authorizer;
    }

    mapping(address => Operator) public operators;

    modifier onlyOperatorAuthorizer(address _operator) {
        require(
            operators[_operator].authorizer == msg.sender,
            "Not operator authorizer"
        );
        _;
    }

    /**
     * @dev Gets the stake balance of the specified address.
     * @param _address The address to query the balance of.
     * @return An uint256 representing the amount staked by the passed address.
     */
    function balanceOf(address _address) public view returns (uint256 balance) {
        return operators[_address].amount;
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
        return operators[_operator].owner;
    }

    /**
     * @dev Gets the magpie for the specified operator address.
     * @return Magpie address.
     */
    function magpieOf(address _operator) public view returns (address payable) {
        return operators[_operator].beneficiary;
    }

    /**
     * @dev Gets the authorizer for the specified operator address.
     * @return Authorizer address.
     */
    function authorizerOf(address _operator) public view returns (address) {
        return operators[_operator].authorizer;
    }
}
