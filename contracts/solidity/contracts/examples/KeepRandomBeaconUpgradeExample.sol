pragma solidity ^0.5.4;

import "../KeepRandomBeaconImplV1.sol";


/**
 * @title KeepRandomBeaconUpgradeExample
 * @dev Example version of a new implementation contract to test upgradability
 * under Keep Random Beacon proxy.
 */
contract KeepRandomBeaconUpgradeExample is KeepRandomBeaconImplV1 {

    uint256 internal _newVar;

    /**
     * @dev Example of overriding existing function.
     * Reference http://solidity.readthedocs.io/en/develop/contracts.html#inheritance
     * >Functions can be overridden by another function with the same name and the
     * same number/types of inputs.
     */
    function initialize(
        uint256 _minPayment, uint256 _withdrawalDelay, uint256 _genesisEntry, 
        bytes memory _genesisGroupPubKey, address _groupContract, uint256 _relayRequestTimeout)
        public
        onlyOwner
    {
        super.initialize(_minPayment, _withdrawalDelay, _genesisEntry, _genesisGroupPubKey, _groupContract, 
        _relayRequestTimeout);
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
