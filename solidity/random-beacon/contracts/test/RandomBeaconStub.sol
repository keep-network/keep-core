pragma solidity ^0.8.6;

import "../RandomBeacon.sol";
import "../libraries/DKG.sol";
import "../libraries/Groups.sol";

contract RandomBeaconStub is RandomBeacon {
    constructor(ISortitionPool _sortitionPool) RandomBeacon(_sortitionPool) {}

    function getDkgData() external view returns (DKG.Data memory) {
        return dkg;
    }

    function getGroups() external view returns (Groups.Group[] memory) {
        return groups.groups;
    }
}
