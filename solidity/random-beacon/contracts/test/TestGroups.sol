// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../libraries/Groups.sol";

contract TestGroups {
    using Groups for Groups.Data;

    Groups.Data internal groups;

    event PendingGroupRegistered(
        uint64 indexed groupId,
        bytes indexed groupPubKey
    );

    function addPendingGroup(
        bytes calldata groupPubKey,
        address[] memory members,
        bytes memory misbehaved
    ) external {
        groups.addPendingGroup(groupPubKey, members, misbehaved);
    }

    function activateGroup(bytes memory groupPubKey) external {
        groups.activateGroup(groupPubKey);
    }

    function getGroup(bytes memory groupPubKey)
        external
        view
        returns (Groups.Group memory)
    {
        return groups.getGroup(groupPubKey);
    }

    function getFlaggedGroupIndex(bytes memory groupPubKey)
        external
        view
        returns (uint256)
    {
        return groups.groupIndices[groupPubKey];
    }

    function getGroups() external view returns (Groups.Group[] memory) {
        return groups.groups;
    }
}
