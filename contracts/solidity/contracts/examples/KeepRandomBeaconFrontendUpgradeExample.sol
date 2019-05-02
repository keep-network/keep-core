pragma solidity ^0.5.4;

import "../KeepRandomBeaconFrontendImplV1.sol";


/**
 * @title KeepRandomBeaconFrontendUpgradeExample
 * @dev Example version of a new implementation contract to test upgradability
 * under Keep Random Beacon proxy.
 */
contract KeepRandomBeaconFrontendUpgradeExample is KeepRandomBeaconFrontendImplV1 {

    uint256 internal _newVar;

    /**
     * @dev Example of overriding existing function.
     * Reference http://solidity.readthedocs.io/en/develop/contracts.html#inheritance
     * >Functions can be overridden by another function with the same name and the
     * same number/types of inputs.
     */
    function initialize(uint256 _minPayment, uint256 withdrawalDelay, address _backendContract)
        public
        onlyOwner
    {
        withdrawalDelay; // Silence unused var
        super.initialize(_minPayment, _withdrawalDelay, _backendContract);
        _initialized["KeepRandomBeaconImplV2"] = true;
        // Example of adding new data to the existing storage.
        _newVar = 1234;
    }

    /**
     * @dev Example of overriding initialized function.
     */
    function initialized() public view returns (bool) {
        return _initialized["KeepRandomBeaconImplV2"];
    }

    /**
     * @dev Example of adding a new function.
     */
    function getNewVar() public view returns (uint256) {
        return _newVar;
    }

    function version() public pure returns (string memory) {
        return "V2";
    }
}
