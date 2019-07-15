pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

/**
 * @title Keep Random Beacon service
 * @dev A proxy contract to provide upgradable Random Beacon functionality.
 * Owner can do upgrades by updating implementation state variable to
 * the address of the upgraded contract. All calls to this proxy contract
 * are delegated to the implementation contract.
 */
contract KeepRandomBeaconService is Ownable {

    // Storage position of the address of the current implementation
    bytes32 private constant implementationPosition = keccak256("network.keep.randombeacon.proxy.implementation");

    event Upgraded(address implementation);

    constructor(address _implementation) public {
        require(_implementation != address(0), "Implementation address can't be zero.");
        setImplementation(_implementation);
    }

    /**
     * @dev Gets the address of the current implementation.
     * @return address of the current implementation.
    */
    function implementation() public view returns (address _implementation) {
        bytes32 position = implementationPosition;
        /* solium-disable-next-line */
        assembly {
            _implementation := sload(position)
        }
    }

    /**
     * @dev Sets the address of the current implementation.
     * @param _implementation address representing the new implementation to be set.
    */
    function setImplementation(address _implementation) internal {
        bytes32 position = implementationPosition;
        /* solium-disable-next-line */
        assembly {
            sstore(position, _implementation)
        }
    }

    /**
     * @dev Delegate call to the current implementation contract.
     */
    function() external payable {
        address _impl = implementation();
        /* solium-disable-next-line */
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
     * @param _implementation Address of the new implementation contract.
     */
    function upgradeTo(address _implementation)
        public
        onlyOwner
    {
        address currentImplementation = implementation();
        require(_implementation != address(0), "Implementation address can't be zero.");
        require(_implementation != currentImplementation, "Implementation address must be different from the current one.");
        setImplementation(_implementation);
        emit Upgraded(_implementation);
    }
}
