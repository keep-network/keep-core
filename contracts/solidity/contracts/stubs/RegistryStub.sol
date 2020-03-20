pragma solidity ^0.5.4;

import "../Registry.sol";

contract RegistryStub is Registry {

    function getGovernance() public view returns (address) {
        return governance;
    }

    function getRegistryKeeper() public view returns (address) {
        return registryKeeper;
    }

    function getPanicButton() public view returns (address) {
        return panicButton;
    }
}