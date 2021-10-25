// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../libraries/Groups.sol";

contract TestGroups {
    using Groups for Groups.Data;

    Groups.Data internal groups;

    event PendingGroupRegistered(bytes indexed groupPubKey);

    event PendingGroupRemoved(bytes indexed groupPubKey);

    event GroupActivated(uint64 indexed groupId, bytes indexed groupPubKey);

    function addPendingGroup(
        bytes calldata groupPubKey,
        address[] memory members,
        bytes memory misbehaved
    ) external {
        groups.addPendingGroup(groupPubKey, members, misbehaved);
    }

    function activateGroup() external {
        groups.activateGroup();
    }

    function getGroupsRegistry() external view returns (bytes32[] memory) {
        return groups.groupsRegistry;
    }

    function getGroup(bytes memory groupPubKey)
        external
        view
        returns (Groups.Group memory)
    {
        return groups.getGroup(groupPubKey);
    }

    function numberOfActiveGroups() external view returns (uint64) {
        return groups.numberOfActiveGroups();
    }
}
