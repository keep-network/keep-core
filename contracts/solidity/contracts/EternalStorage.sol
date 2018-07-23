pragma solidity ^0.4.21;


/**
 * @title EternalStorage
 * @dev This contract holds all the necessary state variables to carry out the storage of any contract.
 */
contract EternalStorage {
    mapping(bytes32 => uint256) internal uintStorage;
    mapping(bytes32 => string) internal stringStorage;
    mapping(bytes32 => address) internal addressStorage;
    mapping(bytes32 => bytes) internal bytesStorage;
    mapping(bytes32 => bytes32) internal bytes32Storage;
    mapping(bytes32 => bool) internal boolStorage;
    mapping(bytes32 => int256) internal intStorage;

	mapping(bytes32 => mapping(uint256 => uint256)) uintStorageMap;
	mapping(bytes32 => mapping(bytes32 => uint256)) uintStorageMap2;
	mapping(bytes32 => mapping(uint256 => address)) addressStorageMap;
	mapping(bytes32 => mapping(uint256 => bool)) boolStorageMap;
	mapping(bytes32 => mapping(bytes32 => bool)) boolStorageMap2;
	mapping(bytes32 => mapping(uint256 => bytes32)) bytes32StorageMap;
	mapping(bytes32 => mapping(bytes32 => bytes32)) bytes32bytes32StorageMap;
}
