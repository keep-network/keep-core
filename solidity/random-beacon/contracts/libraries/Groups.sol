// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

/// @notice This library is used as a registry of created groups.
/// @dev This library should be used along with DKG library that ensures linear
///      groups creation (only one group creation happens at a time). A candidate
///      group has to be popped or activated before adding a new candidate group.
library Groups {
    event CandidateGroupRegistered(bytes indexed groupPubKey);

    event CandidateGroupRemoved(bytes indexed groupPubKey);

    event GroupActivated(uint64 indexed groupId, bytes indexed groupPubKey);

    struct Group {
        bytes groupPubKey;
        uint256 activationTimestamp;
        uint32[] members;
    }

    struct Data {
        // Mapping of keccak256 hashes of group public keys to groups details.
        mapping(bytes32 => Group) groupsData;
        // Holds keccak256 hashes of group public keys in the order of registration.
        bytes32[] groupsRegistry;
        // TODO: Remember about decreasing the counter in case of expiration or termination.
        uint64 activeGroupsCount;
    }

    /// @notice Adds a new candidate group. The group is stored with group public
    ///         key and group members, but is not yet activated.
    /// @dev The group members list is stored with all misbehaved members filtered out.
    /// @param groupPubKey Generated candidate group public key
    /// @param members Addresses of candidate group members as outputted by the
    ///        group selection protocol.
    /// @param misbehavedMembersIndices Array of misbehaved (disqualified or
    ///        inactive) group members indices; Indices reflect positions of
    ///        members in the group, as outputted by the group selection
    ///        protocol.
    function addCandidateGroup(
        Data storage self,
        bytes calldata groupPubKey,
        uint32[] calldata members,
        uint8[] calldata misbehavedMembersIndices
    ) internal {
        bytes32 groupPubKeyHash = keccak256(groupPubKey);

        require(
            self.groupsData[groupPubKeyHash].activationTimestamp == 0,
            "group with this public key was already activated"
        );

        require(
            self.groupsRegistry.length <= type(uint64).max,
            "max number of registered groups reached"
        );

        // We use group from storage that is assumed to be a struct set to the
        // default values. We need to remember to overwrite fields in case a
        // candidate group was already registered before and popped.
        // To ensure we modify only a new group or a popped one, we check the
        // group has not been activated.
        Group storage group = self.groupsData[groupPubKeyHash];
        require(
            group.activationTimestamp == 0,
            "Cannot add group that has been activated"
        );

        group.groupPubKey = groupPubKey;
        setGroupMembers(group, members, misbehavedMembersIndices);

        self.groupsRegistry.push(groupPubKeyHash);

        emit CandidateGroupRegistered(groupPubKey);
    }

    /// @notice Sets addresses of members for the group eliminating members at
    ///         positions pointed by the misbehavedMembersIndices array.
    ///
    ///         NOTE THAT THIS FUNCTION CHANGES ORDER OF MEMBERS IN THE GROUP
    ///         IF THERE IS AT LEAST ONE MISBEHAVED MEMBER
    ///
    ///         The final group members indexes should be obtained post-DKG
    ///         and they may differ from the ones outputted by the group
    ///         selection protocol.
    /// @param group The group storage.
    /// @param members Group member addresses as outputted by the group selection
    ///        protocol.
    /// @param misbehavedMembersIndices Array of misbehaved (disqualified or
    ///        inactive) group members. Indices reflect positions
    ///        of members in the group as outputted by the group selection
    ///        protocol.
    function setGroupMembers(
        Group storage group,
        uint32[] calldata members,
        uint8[] calldata misbehavedMembersIndices
    ) private {
        group.members = members;

        // Iterate misbehaved array backwards, replace misbehaved
        // member with the last element and reduce array length
        uint256 i = misbehavedMembersIndices.length;
        while (i > 0) {
            // group member indices start from 1, so we need to -1 on misbehaved
            uint8 memberArrayPosition = misbehavedMembersIndices[i - 1] - 1;
            group.members[memberArrayPosition] = group.members[
                group.members.length - 1
            ];
            group.members.pop();
            i--;
        }
    }

    /// @notice Removes the latest candidate group.
    /// @dev To optimize gas usage it doesn't delete group details from the
    ///      `groupsData` mapping. The data will be overwritten in case a new
    ///      candidate group gets registered.
    function popCandidateGroup(Data storage self) internal {
        bytes32 groupPubKeyHash = self.groupsRegistry[
            self.groupsRegistry.length - 1
        ];

        require(
            self.groupsData[groupPubKeyHash].activationTimestamp == 0,
            "the latest registered group was already activated"
        );

        self.groupsRegistry.pop();

        emit CandidateGroupRemoved(
            self.groupsData[groupPubKeyHash].groupPubKey
        );
    }

    /// @notice Activates the latest candidate group.
    function activateCandidateGroup(Data storage self) internal {
        Group storage group = self.groupsData[
            self.groupsRegistry[self.groupsRegistry.length - 1]
        ];

        require(
            group.activationTimestamp == 0,
            "the latest registered group was already activated"
        );

        // solhint-disable-next-line not-rely-on-time
        group.activationTimestamp = block.timestamp;

        self.activeGroupsCount++;

        emit GroupActivated(
            uint64(self.groupsRegistry.length - 1),
            group.groupPubKey
        );
    }

    // TODO: Add group termination and expiration

    /// @notice Gets the number of active groups. Candidate, expired and terminated
    ///         groups are not counted as active.
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
        return uint64(seed % self.groupsRegistry.length);
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
