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

    function roughlyAddGroup(
        bytes calldata groupPubKey,
        uint32[] calldata members
    ) external {
        bytes32 groupPubKeyHash = keccak256(groupPubKey);

        Groups.Group memory group;
        group.groupPubKey = groupPubKey;
        group.members = members;
        /* solhint-disable-next-line not-rely-on-time */
        group.activationBlockNumber = block.number;

        groups.groupsData[groupPubKeyHash] = group;
        groups.groupsRegistry.push(groupPubKeyHash);
    }

    function groupLifetimeOf(bytes32 groupPubKeyHash)
        external
        view
        returns (uint256)
    {
        return
            groups.groupsData[groupPubKeyHash].activationBlockNumber +
            groups.groupLifetime;
    }

    function roughlyTerminateGroup(uint64 groupId) public {
        groups.groupsData[groups.groupsRegistry[groupId]].terminated = true;
        // just add groupId without sorting for simplicity
        groups.activeTerminatedGroups.push(groupId);
    }

    function isGroupTerminated(uint64 groupId) external view returns (bool) {
        bytes32 groupPubKeyHash = groups.groupsRegistry[groupId];
        Groups.Group memory group = groups.groupsData[groupPubKeyHash];

        return group.terminated;
    }

    function publicDkgLockState() external {
        dkgLockState();
    }
}
