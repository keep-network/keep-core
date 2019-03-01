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

    /**
     * @dev Returns a group from available groups using modulo operator.
     * @param i Any uint256 number.
     */
    function modSelectGroup(uint256 i) public pure returns(bytes memory) {
        i;
        // Compressed public key (G2 point) generated with Go client using secret key 123
        return hex"1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d0";
    }

}
