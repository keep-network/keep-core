pragma solidity 0.5.17;

import "../Registry.sol";


contract RegistryStub is Registry {
    function getGovernance() public view returns (address) {
        return governance;
    }

    function getRegistryKeeper() public view returns (address) {
        return registryKeeper;
    }

    function getDefaultPanicButton() public view returns (address) {
        return defaultPanicButton;
    }

    function getPanicButtonForContract(address operatorContract)
        public
        view
        returns (address)
    {
        return panicButtons[operatorContract];
    }
}
