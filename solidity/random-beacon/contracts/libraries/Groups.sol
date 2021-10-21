// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./BytesLib.sol";

library Groups {
    using BytesLib for bytes;

    // The index of a group is flagged with the most significant bit set,
    // to distinguish the group `0` from null.
    // The flag is toggled with bitwise XOR (`^`)
    // which keeps all other bits intact but flips the flag bit.
    // The flag should be set before writing to `groupIndices`,
    // and unset after reading from `groupIndices`
    // before using the value.
    uint256 constant GROUP_INDEX_FLAG = 1 << 255;

    event PendingGroupRegistered(
        uint64 indexed groupId,
        bytes indexed groupPubKey
    );

    struct Group {
        bytes groupPubKey;
        uint256 activationTimestamp;
        address[] members;
    }

    struct Data {
        // Mapping of `groupPubKey` to flagged `groupIndex`
        mapping(bytes => uint256) groupIndices;
        // mapping(bytes32 => bytes) gr
        Group[] groups;
        // TODO: Remember about decreasing the counter in case of expiration or termination.
        uint64 activeGroupsCount;
    }

    /// @notice Adds a new group.
    function addPendingGroup(
        Data storage self,
        bytes calldata groupPubKey,
        address[] memory members,
        bytes memory misbehaved
    ) internal {
        if (self.groupIndices[groupPubKey] != 0) {
            require(
                !wasGroupActivated(self, groupPubKey),
                "group was already activated"
            );
        }

        require(
            self.groups.length <= type(uint64).max,
            "max number of groups reached"
        );

        self.groupIndices[groupPubKey] = (self.groups.length ^
            GROUP_INDEX_FLAG);

        Group memory group;
        group.groupPubKey = groupPubKey;
        self.groups.push(group);

        setGroupMembers(_getGroup((self), groupPubKey), members, misbehaved);

        emit PendingGroupRegistered(
            uint64(self.groups.length - 1),
            groupPubKey
        );
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

    // TODO: Could we further optimize this library and don't require groupPubKey
    // to be passed for group activation and removal? Could we assume that
    // the most recent group in the groups stack is a pending group? If so we could
    // also remove storing groupPubKey in the DKG library.

    function activateGroup(Data storage self, bytes memory groupPubKey)
        internal
    {
        require(
            !wasGroupActivated(self, groupPubKey),
            "group was already activated"
        );

        Group storage group = _getGroup(self, groupPubKey);
        group.activationTimestamp = block.timestamp;

        self.activeGroupsCount++;
    }

    // TODO: Add group termination and expiration

    /// @notice Gets the number of active groups. Pending, expired and terminated
    /// groups are not counted as active.
    function numberOfActiveGroups(Data storage self)
        internal
        view
        returns (uint64)
    {
        // TODO: Revisit and include pending, terminated and expired groups
        return self.activeGroupsCount;
        // TODO: Subtract expired and terminated groups
        // .sub(self.expiredGroupOffset).sub(
        //     self.activeTerminatedGroups.length)
    }

    function wasGroupActivated(Data storage self, bytes memory groupPubKey)
        internal
        view
        returns (bool)
    {
        Group memory group = _getGroup(self, groupPubKey);

        return group.activationTimestamp > 0;
    }

    function selectGroup(Data storage self, uint256 seed)
        internal
        view
        returns (uint64)
    {
        // TODO: Implement.
        return uint64(self.groups.length - 1);
    }

    function getGroup(Data storage self, uint64 groupId)
        internal
        view
        returns (Group memory)
    {
        return self.groups[groupId];
    }

    function getGroup(Data storage self, bytes memory groupPubKey)
        internal
        view
        returns (Group memory)
    {
        return _getGroup(self, groupPubKey);
    }

    function _getGroup(Data storage self, bytes memory groupPubKey)
        private
        view
        returns (Group storage)
    {
        uint256 flaggedIndex = self.groupIndices[groupPubKey];
        require(flaggedIndex != 0, "Group does not exist");

        uint256 groupId = flaggedIndex ^ GROUP_INDEX_FLAG;

        return self.groups[groupId];
    }
}
