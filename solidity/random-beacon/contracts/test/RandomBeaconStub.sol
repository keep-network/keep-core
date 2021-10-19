pragma solidity ^0.8.6;

import "../RandomBeacon.sol";
import "../libraries/DKG.sol";

contract RandomBeaconStub is RandomBeacon {
    constructor(ISortitionPool _sortitionPool) RandomBeacon(_sortitionPool) {}

    function getDkgData() external view returns (DKG.Data memory) {
        return dkg;
    }
}
