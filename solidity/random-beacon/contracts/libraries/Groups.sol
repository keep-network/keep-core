// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

/// @notice This library is used as a registry of created groups.
/// @dev This library should be used along with DKG library that ensures linear
///      groups creation (only one group creation happens at a time). A candidate
///      group has to be popped or activated before adding a new candidate group.
library Groups {
    struct Group {
        bytes groupPubKey;
        uint256 activationBlockNumber;
        uint32[] members;
        // When selected group does not create a relay entry on-time it should
        // be marked as terminated.
        bool terminated;
    }

    struct Data {
        // Mapping of keccak256 hashes of group public keys to groups details.
        mapping(bytes32 => Group) groupsData;
        // Holds keccak256 hashes of group public keys in the order of registration.
        bytes32[] groupsRegistry;
        // Group ids that were active but failed creating a relay entry. When an
        // active-terminated group qualifies to become 'expired', then it will
        // be removed from this array.
        uint64[] activeTerminatedGroups;
        // TODO: Remember about decreasing the counter in case of expiration or termination.
        uint64 activeGroupsCount;
        // Points to the first active group, it is also the expired groups counter.
        uint64 expiredGroupOffset;
        // Group lifetime in blocks. When a group reached its lifetime, it
        // is no longer selected for new relay requests but may still be
        // responsible for submitting relay entry if relay request assigned
        // to that group is still pending.
        uint256 groupLifetime;
        // Calculated in the Relay.sol lib.
        uint256 relayEntryTimeout;
    }

    event CandidateGroupRegistered(bytes indexed groupPubKey);

    event CandidateGroupRemoved(bytes indexed groupPubKey);

    event GroupActivated(uint64 indexed groupId, bytes indexed groupPubKey);

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
            self.groupsData[groupPubKeyHash].activationBlockNumber == 0,
            "group with this public key was already activated"
        );

        require(
            self.groupsRegistry.length <= type(uint64).max,
            "max number of registered groups reached"
        );

        // We use group from storage that is assumed to be a struct set to the
        // default values. We need to remember to overwrite fields in case a
        // candidate group was already registered before and popped.
        Group storage group = self.groupsData[groupPubKeyHash];
        group.groupPubKey = groupPubKey;

        setGroupMembers(group, members, misbehavedMembersIndices);

        self.groupsRegistry.push(groupPubKeyHash);

        emit CandidateGroupRegistered(groupPubKey);
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
            self.groupsData[groupPubKeyHash].activationBlockNumber == 0,
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
            group.activationBlockNumber == 0,
            "the latest registered group was already activated"
        );

        // solhint-disable-next-line not-rely-on-time
        group.activationBlockNumber = block.number;

        self.activeGroupsCount++;

        emit GroupActivated(
            uint64(self.groupsRegistry.length - 1),
            group.groupPubKey
        );
    }

    /// @notice Goes through groups starting from the oldest one that is still
    ///         active and checks if it hasn't expired. If so, updates the information
    ///         about expired groups so that all expired groups are marked as such.
    function expireOldGroups(Data storage self) internal {
        // Move expiredGroupOffset as long as there are some groups that should
        // be marked as expired. It is possible that expired group offset will
        // move out of the groups array by one position. It means that all groups
        // are expired (it points to the first active group) and that place in
        // groups array - currently empty - will be possibly filled later by
        // a new group.
        while (
            self.expiredGroupOffset < self.groupsRegistry.length &&
            groupLifetimeOf(
                self,
                self.groupsRegistry[self.expiredGroupOffset]
            ) <
            block.number
        ) {
            self.expiredGroupOffset++;
            self.activeGroupsCount--;
        }

        // Go through all activeTerminatedGroups and if some of the terminated
        // groups are expired, remove them from activeTerminatedGroups collection
        // and rearrange the array to preserve the original order.
        // This is needed because we evaluate the shift of selected group index
        // based on how many non-expired groups have been terminated. Hence it is
        // important that a number of terminated groups matches the length of
        // activeTerminatedGroups[].
        for (uint256 i = 0; i < self.activeTerminatedGroups.length; i++) {
            if (self.expiredGroupOffset > self.activeTerminatedGroups[i]) {
                // When 'i'th group qualifies for expiration, we need to remove
                // it from the activeTerminatedGroups array manually by rearranging
                // the order starting from 'i'th group.
                for (
                    uint256 j = i;
                    j < self.activeTerminatedGroups.length - 1;
                    j++
                ) {
                    self.activeTerminatedGroups[j] = self
                        .activeTerminatedGroups[j + 1];
                }
                // Resizing the array length by 1. The last element was copied
                // over in the loop above to an index "second to last". This is
                // why we can safely remove it from here.
                self.activeTerminatedGroups.pop();
                // At this point the length of activeTerminatedGroups[] was shrinked
                // by 1. We need to adjust 'i'th counter by 1 as well, otherwise
                // it will overflow this array and throw an error.
                i--;
            }
        }
    }

    /// @notice Terminates group with the provided index. Reverts if the group
    ///         is already terminated.
    /// @param  groupId Index in the groupRegistry array.
    function terminateGroup(Data storage self, uint64 groupId) internal {
        require(
            !isGroupTerminated(self, groupId),
            "Group has been already terminated"
        );
        self.groupsData[self.groupsRegistry[groupId]].terminated = true;
        self.activeTerminatedGroups.push();
        self.activeGroupsCount--;

        // Sorting activeTerminatedGroups in ascending order so a non-terminated
        // group is properly selected.
        uint256 i;
        for (
            i = self.activeTerminatedGroups.length - 1;
            i > 0 && self.activeTerminatedGroups[i - 1] > groupId;
            i--
        ) {
            self.activeTerminatedGroups[i] = self.activeTerminatedGroups[i - 1];
        }
        self.activeTerminatedGroups[i] = groupId;
    }

    /// @notice Returns an index of a randomly selected active group. Terminated
    ///         and expired groups are not considered as active.
    ///         Before new group is selected, information about expired groups
    ///         is updated. At least one active group needs to be present for this
    ///         function to succeed.
    /// @param seed Random number used as a group selection seed.
    function selectGroup(Data storage self, uint256 seed)
        internal
        returns (uint64)
    {
        expireOldGroups(self);

        require(numberOfActiveGroups(self) > 0, "No active groups");

        uint64 selectedGroup = uint64(seed % numberOfActiveGroups(self));
        uint64 result = shiftByTerminatedGroups(
            self,
            shiftByExpiredGroups(self, selectedGroup)
        );
        return result;
    }

    /// @notice Setter for relay entry timeout.
    /// @param timeout Relay entry timout calculated in Relay.sol lib.
    function setRelayEntryTimeout(Data storage self, uint256 timeout) internal {
        self.relayEntryTimeout = timeout;
    }

    /// @notice Setter for group lifetime.
    /// @param lifetime Lifetime of a group in blocks.
    function setGroupLifetime(Data storage self, uint256 lifetime) internal {
        self.groupLifetime = lifetime;
    }

    /// @notice Checks if group with the given index is terminated.
    function isGroupTerminated(Data storage self, uint64 groupId)
        internal
        view
        returns (bool)
    {
        return self.groupsData[self.groupsRegistry[groupId]].terminated;
    }

    /// @notice Gets the cutoff time until which the given group is considered
    ///         to be active assuming it hasn't been terminated before.
    function groupLifetimeOf(Data storage self, bytes32 groupPubKeyHash)
        internal
        view
        returns (uint256)
    {
        return
            self.groupsData[groupPubKeyHash].activationBlockNumber +
            self.groupLifetime;
    }

    /// @notice Gets the cutoff time in blocks after which the given group is
    ///         considered as stale. Stale group is an expired group which is no
    ///         longer performing any operations.
    function groupStaleTime(Data storage self, bytes32 groupPubKeyHash)
        internal
        view
        returns (uint256)
    {
        return groupLifetimeOf(self, groupPubKeyHash) + self.relayEntryTimeout;
    }

    /// @notice Checks if a group with the given index is a stale group.
    ///         Stale group is an expired group which is no longer performing any
    ///         operations. It is important to understand that an expired group
    ///         may still perform some operations for which it was selected when
    ///         it was still active. We consider a group to be stale when it's
    ///         expired and when its expiration time and potentially executed
    ///         operation timeout are both in the past.
    function isStaleGroup(Data storage self, bytes memory groupPubKey)
        internal
        view
        returns (bool)
    {
        // TODO: Can a group be considered as "stale" if "expiredGroupOffset" was
        //       not moved forward which means that group selection wasn't triggered?
        //       In other words can it be stale when it's not officially expired?
        return groupStaleTime(self, keccak256(groupPubKey)) < block.number;
    }

    /// @notice Checks if a group with the given index is a stale group.
    ///         Stale group is an expired group which is no longer performing any
    ///         operations. It is important to understand that an expired group
    ///         may still perform some operations for which it was selected when
    ///         it was still active. We consider a group to be stale when it's
    ///         expired and when its expiration time and potentially executed
    ///         operation timeout are both in the past.
    function isStaleGroup(Data storage self, uint64 groupId)
        internal
        view
        returns (bool)
    {
        // TODO: Can a group be considered as "stale" if "expiredGroupOffset" was
        //       not moved forward (group selection wasn't triggered)?
        //       In other words can it be stale when it's not officially expired?
        return
            groupStaleTime(self, self.groupsRegistry[groupId]) < block.number;
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

    /// @notice Gets the number of active groups. Candidate, expired and terminated
    ///         groups are not counted as active.
    function numberOfActiveGroups(Data storage self)
        internal
        view
        returns (uint64)
    {
        return self.activeGroupsCount;
    }

    /// @notice Evaluates the shift of a selected group index based on the number
    ///         of expired groups.
    function shiftByExpiredGroups(Data storage self, uint64 selectedIndex)
        internal
        view
        returns (uint64)
    {
        return self.expiredGroupOffset + selectedIndex;
    }

    /// @notice Evaluates the shift of a selected group index based on the number
    ///         of non-expired but terminated groups.
    function shiftByTerminatedGroups(Data storage self, uint64 selectedIndex)
        internal
        view
        returns (uint64)
    {
        uint64 shiftedIndex = selectedIndex;
        for (uint64 i = 0; i < self.activeTerminatedGroups.length; i++) {
            if (self.activeTerminatedGroups[i] <= shiftedIndex) {
                shiftedIndex++;
            }
        }

        return shiftedIndex;
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
}
