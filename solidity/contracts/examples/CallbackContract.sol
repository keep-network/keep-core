pragma solidity 0.5.17;

import "../IRandomBeacon.sol";

/**
 * @title CallbackContract
 * @dev Example callback contract for Random Beacon.
 */
contract CallbackContract is IRandomBeaconConsumer {

    uint256 internal _lastEntry;

    function __beaconCallback(uint256 entry)
        public
    {
        _lastEntry = entry;
    }

    function lastEntry() public view returns (uint256)
    {
        return _lastEntry;
    }
}
