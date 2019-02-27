pragma solidity ^0.5.4;


/**
 * @title KeepGroupStub
 * @dev A simplified Keep Group contract to help local development.
 */
contract KeepGroupStub {

    uint256 internal _randomBeaconValue;

    /**
     * @dev Triggers the selection process of a new candidate group.
     */
    function runGroupSelection(uint256 randomBeaconValue) public {
        _randomBeaconValue = randomBeaconValue;
    }

}
