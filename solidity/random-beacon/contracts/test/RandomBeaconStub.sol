pragma solidity ^0.8.6;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "../RandomBeacon.sol";
import "../libraries/DKG.sol";
import "../libraries/Callback.sol";
import "../libraries/Groups.sol";

contract RandomBeaconStub is RandomBeacon {
    constructor(
        SortitionPool _sortitionPool,
        IERC20 _tToken,
        IRandomBeaconStaking _staking
    ) RandomBeacon(_sortitionPool, _tToken, _staking) {}

    function getDkgData() external view returns (DKG.Data memory) {
        return dkg;
    }

    function getCallbackData() external view returns (Callback.Data memory) {
        return callback;
    }

    function incrementActiveGroupsCount() external {
        groups.activeGroupsCount++;
    }

    function publicDkgLockState() external {
        dkgLockState();
    }
}
