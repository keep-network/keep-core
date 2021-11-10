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
        uint8[] calldata misbehavedMembrsIndices
    ) external {
        groups.addCandidateGroup(groupPubKey, members, misbehavedMembrsIndices);
    }

    function popCandidateGroup() external {
        groups.popCandidateGroup();
    }

    function activateCandidateGroup() external {
        groups.activateCandidateGroup();
    }

    function selectGroup(uint256 seed) external returns (uint64) {
        return groups.selectGroup(seed);
    }

    function setGroupLifetime(uint256 groupLifetime) external {
        return groups.setGroupLifetime(groupLifetime);
    }

    function setRelayEntryTimeout(uint256 timeout) external {
        return groups.setRelayEntryTimeout(timeout);
    }

    function getGroupsRegistry() external view returns (bytes32[] memory) {
        return groups.groupsRegistry;
    }

    function numberOfActiveGroups() external view returns (uint64) {
        return groups.numberOfActiveGroups();
    }

    function getGroup(bytes memory groupPubKey)
        external
        view
        returns (Groups.Group memory)
    {
        return groups.getGroup(groupPubKey);
    }

    // group id is an index in the groups.groupsRegistry array
    function getGroupById(uint64 groupId)
        external
        view
        returns (Groups.Group memory)
    {
        return groups.getGroup(groupId);
    }

    function isStaleGroup(bytes memory groupPubKey)
        external
        view
        returns (bool)
    {
        return groups.isStaleGroup(groupPubKey);
    }

    function isStaleGroupById(uint64 groupId) external view returns (bool) {
        return groups.isStaleGroup(groupId);
    }
}
