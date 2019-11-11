pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";
import "./cryptography/AltBn128.sol";


interface OperatorContract {
    function relayEntryTimeout() external view returns(uint256);
}

/**
 * @title KeepRandomBeaconOperatorGroups
 * @dev A helper contract for operator contract to store groups and 
 * perform logic to expire and terminate groups.
 */
contract KeepRandomBeaconOperatorGroups {
    using SafeMath for uint256;
    using BytesLib for bytes;

    // Contract owner.
    address public owner;

    // Operator contract that is linked to this contract.
    address public operatorContract;

    // The minimal number of groups that should not expire to protect the
    // minimal network throughput.
    uint256 public activeGroupsThreshold = 5;
 
    // Time in blocks after which a group expires.
    uint256 public groupActiveTime = 3000;

    // Duplicated constant from operator contract to avoid extra call.
    // The value is set when the operator contract is added.
    uint256 public relayEntryTimeout;

    struct Group {
        bytes groupPubKey;
        uint registrationBlockHeight;
    }

    Group[] internal groups;
    uint256[] internal terminatedGroups;
    mapping (bytes => address[]) internal groupMembers;

    // Sum of all group member rewards earned so far. The value is the same for
    // all group members. Submitter reward and reimbursement is paid immediately
    // and is not included here. Each group member can withdraw no more than
    // this value.
    mapping (bytes => uint256) internal groupMemberRewards;

    // expiredGroupOffset is pointing to the first active group, it is also the
    // expired groups counter
    uint256 internal expiredGroupOffset = 0;

    /**
     * @dev Throws if called by any account other than the owner.
     */
    modifier onlyOwner() {
        require(owner == msg.sender, "Caller is not the owner.");
        _;
    }

    /**
     * @dev Throws if called by any account other than the authorized address.
     */
    modifier onlyOperatorContract() {
        require(operatorContract == msg.sender, "Caller is not authorized.");
        _;
    }

    /**
     * @dev Initializes the contract with deployer as the contract owner.
     */
    constructor() public {
        owner = msg.sender;
    }

    /**
     * @dev Sets operator contract.
     */
    function setOperatorContract(address _operatorContract) public onlyOwner {
        require(operatorContract == address(0), "Operator contract can only be set once.");
        operatorContract = _operatorContract;
        relayEntryTimeout = OperatorContract(operatorContract).relayEntryTimeout();
    }

    /**
     * @dev Adds group.
     */
    function addGroup(bytes memory groupPubKey) public onlyOperatorContract {
        groups.push(Group(groupPubKey, block.number));
    }

    /**
     * @dev Adds group member.
     */
    function addGroupMember(bytes memory groupPubKey, address member) public onlyOperatorContract {
        groupMembers[groupPubKey].push(member);
    }

    /**
     * @dev Adds group member reward per group so the accumulated amount can be withdrawn later.
     */
    function addGroupMemberReward(bytes memory groupPubKey, uint256 amount) public onlyOperatorContract {
        groupMemberRewards[groupPubKey] = groupMemberRewards[groupPubKey].add(amount);
    }

    /**
     * @dev Returns accumulated group member rewards for provided group.
     */
    function getGroupMemberRewards(bytes memory groupPubKey) public view returns (uint256) {
        return groupMemberRewards[groupPubKey];
    }

    /**
     * @dev Gets group public key.
     */
    function getGroupPublicKey(uint256 groupIndex) public view returns (bytes memory) {
        return groups[groupIndex].groupPubKey;
    }

    /**
     * @dev Gets group public key in a compressed form.
     */
    function getGroupPublicKeyCompressed(uint256 groupIndex) public view returns (bytes memory) {
        return AltBn128.g2Compress(AltBn128.g2Unmarshal(groups[groupIndex].groupPubKey));
    }

    /**
     * @dev Gets group member.
     */
    function getGroupMember(bytes memory groupPubKey, uint256 memberIndex) public view returns (address) {
        return groupMembers[groupPubKey][memberIndex];
    }

    /**
     * @dev Terminates group.
     */
    function terminateGroup(uint256 groupIndex) public onlyOperatorContract {
        terminatedGroups.push(groupIndex);
    }

    /**
     * @dev Checks if group with the given public key is registered.
     */
    function isGroupRegistered(bytes memory groupPubKey) public view returns(bool) {
        for (uint i = 0; i < groups.length; i++) {
            if (groups[i].groupPubKey.equalStorage(groupPubKey)) {
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
    function groupActiveTimeOf(Group memory group) internal view returns(uint256) {
        return group.registrationBlockHeight.add(groupActiveTime);
    }

    /**
     * @dev Gets the cutoff time in blocks after which the given group is
     * considered as stale. Stale group is an expired group which is no longer
     * performing any operations.
     */
    function groupStaleTime(Group memory group) internal view returns(uint256) {
        return groupActiveTimeOf(group).add(relayEntryTimeout);
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
    function isStaleGroup(bytes memory groupPubKey) public view returns(bool) {
        for (uint i = 0; i < groups.length; i++) {
            if (groups[i].groupPubKey.equalStorage(groupPubKey)) {
                bool isExpired = expiredGroupOffset > i;
                bool isStale = groupStaleTime(groups[i]) < block.number;
                return isExpired && isStale;
            }
        }

        return true; // no group found, consider it as a stale group
    }

    /**
     * @dev Gets list of indices of staled groups.
     */
    function getStaleGroupsIndices() public view returns(uint256[] memory indices) {
        uint256 counter;
        for (uint i = 0; i < groups.length; i++) {
            if (isStaleGroup(groups[i].groupPubKey)) {
                counter++;
            }
        }

        indices = new uint256[](counter);
        counter = 0;
        for (uint i = 0; i < groups.length; i++) {
            if (isStaleGroup(groups[i].groupPubKey)) {
                indices[counter] = i;
                counter++;
            }
        }
    }

    /**
     * @dev Gets all indices in the provided group for a member.
     */
    function getGroupMemberIndices(bytes memory groupPubKey, address member) public view returns (uint256[] memory indices) {
        uint256 counter;
        for (uint i = 0; i < groupMembers[groupPubKey].length; i++) {
            if (groupMembers[groupPubKey][i] == member) {
                counter++;
            }
        }

        indices = new uint256[](counter);
        counter = 0;
        for (uint i = 0; i < groupMembers[groupPubKey].length; i++) {
            if (groupMembers[groupPubKey][i] == member) {
                indices[counter] = i;
                counter++;
            }
        }
    }

    /**
     * @dev Gets group member rewards available to withdraw for a member.
     */
    function availableRewards(address groupMember) public view returns (uint256 rewards) {
        uint256[] memory staleGroupsIndices = getStaleGroupsIndices();
        for (uint i = 0; i < staleGroupsIndices.length; i++) {
            bytes memory groupPublicKey = getGroupPublicKey(staleGroupsIndices[i]);
            uint256[] memory groupMemberIndices = getGroupMemberIndices(groupPublicKey, groupMember);
            for (uint j = 0; j < groupMemberIndices.length; j++) {
                rewards = rewards.add(groupMemberRewards[groupPublicKey]);
            }
        }
    }

    /**
     * @dev Withdraw accumulated group member rewards for a member.
     */
    function withdraw(address groupMember) public returns (uint256 rewards) {
        uint256[] memory staleGroupsIndices = getStaleGroupsIndices();
        for (uint i = 0; i < staleGroupsIndices.length; i++) {
            bytes memory groupPublicKey = getGroupPublicKey(staleGroupsIndices[i]);
            uint256[] memory groupMemberIndices = getGroupMemberIndices(groupPublicKey, groupMember);
            for (uint j = 0; j < groupMemberIndices.length; j++) {
                delete groupMembers[groupPublicKey][groupMemberIndices[j]];
                rewards = rewards.add(groupMemberRewards[groupPublicKey]);
            }
        }
    }

    /**
     * @dev Gets the number of active groups. Expired and terminated groups are
     * not counted as active.
     */
    function numberOfGroups() public view returns(uint256) {
        return groups.length.sub(expiredGroupOffset).sub(terminatedGroups.length);
    }

    /**
     * @dev Goes through groups starting from the oldest one that is still
     * active and checks if it hasn't expired. If so, updates the information
     * about expired groups so that all expired groups are marked as such.
     * It does not mark more than `activeGroupsThreshold` active groups as
     * expired.
     */
    function expireOldGroups() internal {
        // move expiredGroupOffset as long as there are some groups that should
        // be marked as expired and we are above activeGroupsThreshold of
        // active groups.
        while(
            groupActiveTimeOf(groups[expiredGroupOffset]) < block.number &&
            numberOfGroups() > activeGroupsThreshold
        ) {
            expiredGroupOffset++;
        }

        // Go through all terminatedGroups and if some of the terminated
        // groups are expired, remove them from terminatedGroups collection.
        // This is needed because we evaluate the shift of selected group index
        // based on how many non-expired groups has been terminated.
        for (uint i = 0; i < terminatedGroups.length; i++) {
            if (expiredGroupOffset > terminatedGroups[i]) {
                terminatedGroups[i] = terminatedGroups[terminatedGroups.length - 1];
                terminatedGroups.length--;
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
    function selectGroup(uint256 seed) public onlyOperatorContract returns(uint256) {
        require(numberOfGroups() > 0, "At least one active group required");

        expireOldGroups();
        uint256 selectedGroup = seed % numberOfGroups();
        return shiftByTerminatedGroups(shiftByExpiredGroups(selectedGroup));
    }

    /**
     * @dev Evaluates the shift of selected group index based on the number of
     * expired groups.
     */
    function shiftByExpiredGroups(uint256 selectedIndex) internal view returns(uint256) {
        return expiredGroupOffset.add(selectedIndex);
    }

    /**
     * @dev Evaluates the shift of selected group index based on the number of
     * non-expired, terminated groups.
     */
    function shiftByTerminatedGroups(uint256 selectedIndex) internal view returns(uint256) {
        uint256 shiftedIndex = selectedIndex;
        for (uint i = 0; i < terminatedGroups.length; i++) {
            if (terminatedGroups[i] <= shiftedIndex) {
                shiftedIndex++;
            }
        }

        return shiftedIndex;
    }

}
