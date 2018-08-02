pragma solidity ^0.4.24;


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
    mapping(bytes32 => mapping(bytes32 => uint256)) uintBytes32StorageMap;
    mapping(bytes32 => mapping(uint256 => address)) addressUintStorageMap;
    mapping(bytes32 => mapping(uint256 => bool)) boolUintStorageMap;
    mapping(bytes32 => mapping(bytes32 => bool)) boolBytes32StorageMap;
    mapping(bytes32 => mapping(uint256 => bytes32)) bytes32UintStorageMap;
    mapping(bytes32 => mapping(bytes32 => bytes32)) bytes32StorageMap;

    mapping(bytes32 => bytes32[]) bytes32StorageArray;
}
