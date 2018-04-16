pragma solidity ^0.4.18;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";


/**
 * @title Keep Random Beacon
 * @dev A proxy contract to provide upgradable Random Beacon functionality.
 * Owner can do upgrades by updating implementation state variable to
 * the address of the upgraded contract. All calls to this proxy contract
 * are delegated to the implementation contract.
 */
contract KeepRandomBeacon is Ownable {

    // Current implementation contract address.
    address public implementation;

    // Current implementation version.
    string public version;

    event Upgraded(string version, address indexed implementation);

    function KeepRandomBeacon(string _version, address _implementation) {
        require(_implementation != address(0));
        version = _version;
        implementation = _implementation;
    }

    /**
     * @dev Delegate call to the current implementation contract.
     */
    function() payable {
        assembly {
            let ptr := mload(0x40)
            calldatacopy(ptr, 0, calldatasize)
            let result := delegatecall(gas, implementation, ptr, calldatasize, 0, 0)
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
        require(_version != version);
        version = _version;
        implementation = _implementation;
        Upgraded(version, implementation);
    }
}
