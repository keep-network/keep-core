pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/cryptography/ECDSA.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";
import "../StakingProxy.sol";


/**
 * @title Stake Delegatable
 * @dev A base contract to allow stake delegation for staking contracts.
 */
contract StakeDelegatable {
    using SafeMath for uint256;
    using BytesLib for bytes;
    using ECDSA for bytes32;

    StakingProxy public stakingProxy;

    // Stake delegation mappings.
    mapping(address => address) public operatorToOwner;
    mapping(address => address) public magpieToOwner;

    // List of operators for the stake owner.
    mapping(address => address[]) public ownerOperators;

    /**
     * @dev Checks if sender is eligible to call.
     * @param _staker Address of the stake owner or operator.
     */
    modifier onlyOwnerOrOperator(address _staker) {
        require(
            msg.sender == _staker ||
            msg.sender == operatorToOwner[_staker],
            "Only stake owner or operator can call this function."
        );
        _;
    }

    /**
     * @dev Gets the list of operators of the specified address.
     * @return An array of addresses.
     */
    function operatorsOf(address _address) public view returns (address[] memory) {
        return ownerOperators[_address];
    }

    /**
     * @dev Gets magpie and operator address from the delegation data bytes array.
     * @param _from The owner of the tokens.
     * @param _delegationData Data for stake delegation. This byte array must have the
     * following values concatenated: Magpie address (20 bytes) where the rewards for participation
     * are sent and the operator's ECDSA (65 bytes) signature of the address of the stake owner.
     */
    function _extractDelegationData(address _from, bytes memory _delegationData) internal returns (address magpie, address operator) {
        require(_delegationData.length == 85, "Stake delegation data must be provided.");
        operator = keccak256(abi.encodePacked(_from)).toEthSignedMessageHash().recover(_delegationData.slice(20, 65));
        return (_delegationData.toAddress(0), operator);
    }

    /**
     * @dev Delegates stake.
     * @param _from The owner of the tokens.
     * @param _value Amount to stake.
     * @param _magpie Address where the rewards for participation are sent.
     * @param _operator Address of the stake operator.
     */
    function _delegateStake(address _from, uint256 _value, address _magpie, address _operator) internal {
        require(operatorToOwner[_operator] == address(0), "Operator address is already in use.");
        operatorToOwner[_operator] = _from;
        magpieToOwner[_magpie] = _from;
        ownerOperators[_from].push(_operator);
    }
}
