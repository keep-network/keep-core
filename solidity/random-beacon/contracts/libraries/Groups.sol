// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./BytesLib.sol";

library Groups {
    using BytesLib for bytes;

    event CandidateGroupRegistered(bytes indexed groupPubKey);

    event CandidateGroupRemoved(bytes indexed groupPubKey);

    event GroupActivated(uint64 indexed groupId, bytes indexed groupPubKey);

    struct Group {
        bytes groupPubKey;
        uint256 activationTimestamp;
        address[] members;
    }

    struct Data {
        mapping(bytes32 => Group) groupsData;
        bytes32[] groupsRegistry;
        // TODO: Remember about decreasing the counter in case of expiration or termination.
        uint64 activeGroupsCount;
    }

    /// @notice Adds a new group.
    function addCandidateGroup(
        Data storage self,
        bytes calldata groupPubKey,
        address[] memory members,
        bytes memory misbehaved
    ) internal {
        bytes32 groupPubKeyHash = keccak256(groupPubKey);

        require(
            self.groupsData[groupPubKeyHash].activationTimestamp == 0,
            "group was already activated"
        );

        require(
            self.groupsRegistry.length <= type(uint64).max,
            "max number of registered groups reached"
        );

        Group memory group;
        group.groupPubKey = groupPubKey;
        self.groupsData[groupPubKeyHash] = group;

        // FIXME: We can end up with multiple candidate groups registered with the same public key,
        // and later being activated. Although it shouldn't happen when this library
        // is used with Random Beacon DKG as it's handling just one result/group
        // at a time.
        self.groupsRegistry.push(groupPubKeyHash);

        setGroupMembers(self.groupsData[groupPubKeyHash], members, misbehaved);

        emit CandidateGroupRegistered(groupPubKey);
    }

    // TODO: This function should be optimized for members storing.
    // See https://github.com/keep-network/keep-core/pull/2666/files#r732629138
    function setGroupMembers(
        Group storage group,
        address[] memory members,
        bytes memory misbehaved
    ) private {
        group.members = members;

        // Iterate misbehaved array backwards, replace misbehaved
        // member with the last element and reduce array length
        uint256 i = misbehaved.length;
        while (i > 0) {
            // group member indexes start from 1, so we need to -1 on misbehaved
            uint256 memberArrayPosition = misbehaved.toUint8(i - 1) - 1;
            group.members[memberArrayPosition] = group.members[
                group.members.length - 1
            ];
            group.members.pop();
            i--;
        }
    }

    function popCandidateGroup(Data storage self) internal {
        bytes32 groupPubKeyHash = self.groupsRegistry[
            self.groupsRegistry.length - 1
        ];

        require(
            self.groupsData[groupPubKeyHash].activationTimestamp == 0,
            "the latest registered group was already activated"
        );

        // Pop the latest registered group. Group data are not deleted from the
        // `groupsData` map to reduce transaction costs.
        self.groupsRegistry.pop();

        emit CandidateGroupRemoved(
            self.groupsData[groupPubKeyHash].groupPubKey
        );
    }

    function activateGroup(Data storage self) internal {
        Group storage group = self.groupsData[
            self.groupsRegistry[self.groupsRegistry.length - 1]
        ];

        require(group.activationTimestamp == 0, "group was already activated");

        group.activationTimestamp = block.timestamp;

        self.activeGroupsCount++;

        emit GroupActivated(
            uint64(self.groupsRegistry.length - 1),
            group.groupPubKey
        );
    }

    // TODO: Add group termination and expiration

    /// @notice Gets the number of active groups. Candidate, expired and terminated
    /// groups are not counted as active.
    function numberOfActiveGroups(Data storage self)
        internal
        view
        returns (uint64)
    {
        // TODO: Revisit and include candidate, terminated and expired groups
        return self.activeGroupsCount;
        // TODO: Subtract expired and terminated groups
        // .sub(self.expiredGroupOffset).sub(
        //     self.activeTerminatedGroups.length)
    }

    function selectGroup(Data storage self, uint256 seed)
        internal
        view
        returns (uint64)
    {
        // TODO: Implement.
        return uint64(self.groupsRegistry.length - 1);
    }

    function getGroup(Data storage self, uint64 groupId)
        internal
        view
        returns (Group memory)
    {
        return self.groupsData[self.groupsRegistry[groupId]];
    }

    function getGroup(Data storage self, bytes memory groupPubKey)
        internal
        view
        returns (Group memory)
    {
        return self.groupsData[keccak256(groupPubKey)];
    }
}
