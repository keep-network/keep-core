pragma solidity ^0.5.4;
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";
import "../../cryptography/AltBn128.sol";

library Groups {
    using SafeMath for uint256;
    using BytesLib for bytes;

    struct Group {
        bytes groupPubKey;
        uint registrationBlockHeight;
    }

    struct Storage {
        // The minimal number of groups that should not expire to protect the
        // minimal network throughput.
        uint256 activeGroupsThreshold;
    
        // Time in blocks after which a group expires.
        uint256 groupActiveTime;

        // Duplicated constant from operator contract to avoid extra call.
        // The value is set when the operator contract is added.
        uint256 relayEntryTimeout;

        Group[] groups;
        uint256[] terminatedGroups;
        mapping (bytes => address[]) groupMembers;

        // Sum of all group member rewards earned so far. The value is the same for
        // all group members. Submitter reward and reimbursement is paid immediately
        // and is not included here. Each group member can withdraw no more than
        // this value.
        mapping (bytes => uint256) groupMemberRewards;

        // expiredGroupOffset is pointing to the first active group, it is also the
        // expired groups counter
        uint256 expiredGroupOffset;
    }

    /**
     * @dev Adds group.
     */
    function addGroup(
        Storage storage self,
        bytes memory groupPubKey
    ) internal {
        self.groups.push(Group(groupPubKey, block.number));
    }

    /**
     * @dev Adds group member.
     */
    function addGroupMember(
        Storage storage self,
        bytes memory groupPubKey,
        address member
    ) internal {
        self.groupMembers[groupPubKey].push(member);
    }

    /**
     * @dev Adds group member reward per group so the accumulated amount can be withdrawn later.
     */
    function addGroupMemberReward(
        Storage storage self,
        bytes memory groupPubKey,
        uint256 amount
    ) internal {
        self.groupMemberRewards[groupPubKey] = self.groupMemberRewards[groupPubKey].add(amount);
    }

    /**
     * @dev Returns accumulated group member rewards for provided group.
     */
    function getGroupMemberRewards(
        Storage storage self,
        bytes memory groupPubKey
    ) internal view returns (uint256) {
        return self.groupMemberRewards[groupPubKey];
    }

    /**
     * @dev Gets group public key.
     */
    function getGroupPublicKey(
        Storage storage self,
        uint256 groupIndex
    ) internal view returns (bytes memory) {
        return self.groups[groupIndex].groupPubKey;
    }

    /**
     * @dev Gets group member.
     */
    function getGroupMember(
        Storage storage self,
        bytes memory groupPubKey,
        uint256 memberIndex
    ) internal view returns (address) {
        return self.groupMembers[groupPubKey][memberIndex];
    }

    /**
     * @dev Gets all indices in the provided group for a member.
     */
    function getGroupMemberIndices(
        Storage storage self,
        bytes memory groupPubKey,
        address member
    ) public view returns (uint256[] memory indices) {
        uint256 counter;
        for (uint i = 0; i < self.groupMembers[groupPubKey].length; i++) {
            if (self.groupMembers[groupPubKey][i] == member) {
                counter++;
            }
        }

        indices = new uint256[](counter);
        counter = 0;
        for (uint i = 0; i < self.groupMembers[groupPubKey].length; i++) {
            if (self.groupMembers[groupPubKey][i] == member) {
                indices[counter] = i;
                counter++;
            }
        }
    }

    /**
     * @dev Terminates group.
     */
    function terminateGroup(
        Storage storage self,
        uint256 groupIndex
    ) internal {
        self.terminatedGroups.push(groupIndex);
    }

    /**
     * @dev Checks if group with the given public key is registered.
     */
    function isGroupRegistered(
        Storage storage self,
        bytes memory groupPubKey
    ) internal view returns(bool) {
        for (uint i = 0; i < self.groups.length; i++) {
            if (self.groups[i].groupPubKey.equalStorage(groupPubKey)) {
                return true;
            }
        }
        return false;
    }

    /**
     * @dev Gets the cutoff time in blocks until which the given group is
     * considered as an active group assuming it hasn't been terminated before.
     * The group may not be marked as expired even though its active
     * time has passed if one of the rules inside `selectGroup` function are not
     * met (e.g. minimum active group threshold). Hence, this value informs when
     * the group may no longer be considered as active but it does not mean that
     * the group will be immediatelly considered not as such.
     */
    function groupActiveTimeOf(
        Storage storage self,
        Group memory group
    ) internal view returns(uint256) {
        return group.registrationBlockHeight.add(self.groupActiveTime);
    }

    /**
     * @dev Gets the cutoff time in blocks after which the given group is
     * considered as stale. Stale group is an expired group which is no longer
     * performing any operations.
     */
    function groupStaleTime(
        Storage storage self,
        Group memory group
    ) internal view returns(uint256) {
        return groupActiveTimeOf(self, group).add(self.relayEntryTimeout);
    }

    /**
     * @dev Checks if a group with the given public key is a stale group.
     * Stale group is an expired group which is no longer performing any
     * operations. It is important to understand that an expired group may
     * still perform some operations for which it was selected when it was still
     * active. We consider a group to be stale when it's expired and when its
     * expiration time and potentially executed operation timeout are both in
     * the past.
     */
    function isStaleGroup(
        Storage storage self,
        bytes memory groupPubKey
    ) public view returns(bool) {
        for (uint i = 0; i < self.groups.length; i++) {
            if (self.groups[i].groupPubKey.equalStorage(groupPubKey)) {
                bool isExpired = self.expiredGroupOffset > i;
                bool isStale = groupStaleTime(self, self.groups[i]) < block.number;
                return isExpired && isStale;
            }
        }

        revert("Group does not exist");
    }

    /**
     * @dev Gets the number of active groups. Expired and terminated groups are
     * not counted as active.
     */
    function numberOfGroups(
        Storage storage self
    ) internal view returns(uint256) {
        return self.groups.length.sub(self.expiredGroupOffset).sub(self.terminatedGroups.length);
    }

    /**
     * @dev Goes through groups starting from the oldest one that is still
     * active and checks if it hasn't expired. If so, updates the information
     * about expired groups so that all expired groups are marked as such.
     * It does not mark more than `activeGroupsThreshold` active groups as
     * expired.
     */
    function expireOldGroups(
        Storage storage self
    ) internal {
        // move expiredGroupOffset as long as there are some groups that should
        // be marked as expired and we are above activeGroupsThreshold of
        // active groups.
        while(
            groupActiveTimeOf(self, self.groups[self.expiredGroupOffset]) < block.number &&
            numberOfGroups(self) > self.activeGroupsThreshold
        ) {
            self.expiredGroupOffset++;
        }

        // Go through all terminatedGroups and if some of the terminated
        // groups are expired, remove them from terminatedGroups collection.
        // This is needed because we evaluate the shift of selected group index
        // based on how many non-expired groups has been terminated.
        for (uint i = 0; i < self.terminatedGroups.length; i++) {
            if (self.expiredGroupOffset > self.terminatedGroups[i]) {
                self.terminatedGroups[i] = self.terminatedGroups[self.terminatedGroups.length - 1];
                self.terminatedGroups.length--;
            }
        }
    }

    /**
     * @dev Returns an index of a randomly selected active group. Terminated and
     * expired groups are not considered as active.
     * Before new group is selected, information about expired groups
     * is updated. At least one active group needs to be present for this
     * function to succeed.
     * @param seed Random number used as a group selection seed.
     */
    function selectGroup(
        Storage storage self,
        uint256 seed
    ) public returns(uint256) {
        require(numberOfGroups(self) > 0, "At least one active group required");

        expireOldGroups(self);
        uint256 selectedGroup = seed % numberOfGroups(self);
        return shiftByTerminatedGroups(self, shiftByExpiredGroups(self, selectedGroup));
    }

    /**
     * @dev Evaluates the shift of selected group index based on the number of
     * expired groups.
     */
    function shiftByExpiredGroups(
        Storage storage self,
        uint256 selectedIndex
    ) internal view returns(uint256) {
        return self.expiredGroupOffset.add(selectedIndex);
    }

    /**
     * @dev Evaluates the shift of selected group index based on the number of
     * non-expired, terminated groups.
     */
    function shiftByTerminatedGroups(
        Storage storage self,
        uint256 selectedIndex
    ) internal view returns(uint256) {
        uint256 shiftedIndex = selectedIndex;
        for (uint i = 0; i < self.terminatedGroups.length; i++) {
            if (self.terminatedGroups[i] <= shiftedIndex) {
                shiftedIndex++;
            }
        }

        return shiftedIndex;
    }

    /**
     * @dev Withdraws accumulated group member rewards for msg.sender
     * using the provided group index and member indices. Once the
     * accumulated reward is withdrawn from the selected group, member is
     * removed from it. Rewards can be withdrawn only from stale group.
     *
     * @param groupIndex Group index.
     * @param groupMemberIndices Array of member indices for the group member.
     */
    function withdrawFromGroup(
        Storage storage self,
        uint256 groupIndex,
        uint256[] memory groupMemberIndices
    ) public returns (uint256 rewards) {
        for (uint i = 0; i < groupMemberIndices.length; i++) {
            bool isExpired = self.expiredGroupOffset > groupIndex;
            bool isStale = groupStaleTime(self, self.groups[groupIndex]) < block.number;

            bytes memory groupPublicKey = getGroupPublicKey(self, groupIndex);
            if (isExpired && isStale && msg.sender == self.groupMembers[groupPublicKey][groupMemberIndices[i]]) {
                delete self.groupMembers[groupPublicKey][groupMemberIndices[i]];
                rewards = rewards.add(self.groupMemberRewards[groupPublicKey]);
            }
        }
    }
}
