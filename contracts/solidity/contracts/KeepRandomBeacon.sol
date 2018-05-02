pragma solidity ^0.4.21;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";
import "./EternalStorage.sol";

/**
 * @title Keep Random Beacon
 * @dev A proxy contract to provide upgradable Random Beacon functionality.
 * Owner can do upgrades by updating implementation state variable to
 * the address of the upgraded contract. All calls to this proxy contract
 * are delegated to the implementation contract.
 */
contract KeepRandomBeacon is Ownable, EternalStorage {

    // Current implementation contract address.
    address public implementation;

    // Current implementation version.
    string public version;

    event Upgraded(string version, address implementation);

    // Mirror events from the implementation contract
    event RelayEntryRequested(uint256 requestID, uint256 payment, uint256 blockReward, uint256 seed, uint blockNumber); 
    event RelayEntryGenerated(uint256 requestID, uint256 requestResponse, uint256 requestGroupID, uint256 previousEntry, uint blockNumber); 
    event SubmitGroupPublicKeyEvent(byte[] groupPublicKey, uint256 requestID, uint256 activationBlockHeight);

    function KeepRandomBeacon(string _version, address _implementation) {
        require(_implementation != address(0));
        version = _version;
        implementation = _implementation;
    }

    /**
     * @dev Delegate call to the current implementation contract.
     */
    function() payable {
        address _impl = implementation;
        assembly {
            let ptr := mload(0x40)
            calldatacopy(ptr, 0, calldatasize)
            let result := delegatecall(gas, _impl, ptr, calldatasize, 0, 0)
            let size := returndatasize
            returndatacopy(ptr, 0, size)

            switch result
            case 0 { revert(ptr, size) }
            default { return(ptr, size) }
        }
    }

    /**
     * @dev Upgrade current implementation.
     * @param _version Version name for the new implementation.
     * @param _implementation Address of the new implementation contract.
     */
    function upgradeTo(string _version, address _implementation)
        public
        onlyOwner
    {
        require(_implementation != address(0));
        require(_implementation != implementation);
        require(keccak256(_version) != keccak256(version));
        version = _version;
        implementation = _implementation;
        emit Upgraded(version, implementation);
    }
}
