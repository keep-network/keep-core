// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../libraries/Groups.sol";

contract GroupsStub {
    using Groups for Groups.Data;

    Groups.Data internal groups;

    event CandidateGroupRegistered(bytes indexed groupPubKey);

    event CandidateGroupRemoved(bytes indexed groupPubKey);

    event GroupActivated(uint64 indexed groupId, bytes indexed groupPubKey);

    function addCandidateGroup(
        bytes calldata groupPubKey,
        uint32[] calldata members,
        uint8[] calldata misbehaved
    ) external {
        groups.addCandidateGroup(groupPubKey, members, misbehaved);
    }

    function popCandidateGroup() external {
        groups.popCandidateGroup();
    }

    function activateCandidateGroup() external {
        groups.activateCandidateGroup();
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
