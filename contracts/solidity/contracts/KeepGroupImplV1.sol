pragma solidity ^0.4.21;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";
import "./EternalStorage.sol";
import "./KeepRandomBeaconImplV1.sol";


/**
 * @dev Interface for checking minimum stake balance.
 */
// interface keepRandomBeacon {
//     function hasMinimumStake(address _staker) external view returns(bool);
// }


contract KeepGroupImplV1 is Ownable, EternalStorage {

    event GroupExistsEvent(bytes32 groupPubKey, bool exists);
    event GroupStartedEvent(bytes32 groupPubKey);
    event GroupCompleteEvent(bytes32 groupPubKey);
    event GroupErrorCode(uint8 code);

    // TODO: make sure we know staker eth address so we can check its minimum stake
    // modifier hasMinimumStake(bytes32 _staker) {
    //     //keepRandomBeacon beacon = keepRandomBeacon(addressStorage[keccak256("keepRandomBeaconAddress")]);
    //     //require(beacon.hasMinimumStake(_staker));
    //     _;
    // }

    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() public payable {
        revert();
    }

    /**
     * @dev Initialize Keep Group implementaion contract with a linked Keep Random Beacon contract.
     * @param _keepRandomBeaconAddress Address of Keep Random Beacon that will be linked to this contract.
     * @param _groupThreshold Max number of bad members in a group that we can detect as well as “number
     * of good members needed to produce a relay entry”.
     * @param _groupSize Minimum number of members in a group - to form a group.
     */
    function initialize(uint256 _groupThreshold, uint256 _groupSize, address _keepRandomBeaconAddress) public onlyOwner {
        require(!initialized());
        require(_keepRandomBeaconAddress != address(0x0));
        boolStorage[keccak256("KeepGroupImplV1")] = true;
        addressStorage[keccak256("keepRandomBeaconAddress")] = _keepRandomBeaconAddress;
        uintStorage[keccak256("groupThreshold")] = _groupThreshold;
        uintStorage[keccak256("groupSize")] = _groupSize;
        uintStorage[keccak256("groupsCount")] = 0;
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return boolStorage[keccak256("KeepGroupImplV1")];
    }

    /**
     * @dev Sets new threshold size for groups.
     */
    function setGroupThreshold(uint256 _groupThreshold) public onlyOwner {
        uintStorage[keccak256("groupThreshold")] = _groupThreshold;
        /// TODO: determine if size decreased, then partially complete groups may now be complete.  Iterate over groups. Find
    }

    /**
     * @dev Gets the threshold size for groups.
     */
    function groupThreshold() public view returns(uint256) {
        return uintStorage[keccak256("groupThreshold")];
    }

    /**
     * @dev Sets the minimum number of members in a group.
     */
    function setGroupSize(uint256 _groupSize) public onlyOwner {
        uintStorage[keccak256("groupSize")] = _groupSize;
    }

    /**
     * @dev Gets the minimum number of members in a group.
     */
    function groupSize() public view returns(uint256) {
        return uintStorage[keccak256("groupSize")];
    }

    /**
     * @dev Gets number of active groups.
     */
    function numberOfGroups() public view returns(uint256) {
        return uintStorage[keccak256("groupsCount")];
    }

    /**
     * @dev Gets public key of a group.
     * @param _groupIndex Index number of a group.
     */
    function getGroupPubKey(uint256 _groupIndex) public view returns(bytes32) {
        return bytes32Storage[keccak256("groupIndexToGroupPubKey", _groupIndex)];
    }

    /**
     * @dev Gets group index number.
     * @param _groupPubKey Group public key.
     */
    function getGroupIndex(bytes32 _groupPubKey) public view returns(uint) {
        for (uint i = 0; i < uintStorage[keccak256("groupsCount")]; i++) {
            if (bytes32Storage[keccak256("groupIndexToGroupPubKey", i)] == _groupPubKey) {
                return i;
            }
        }
        revert();
    }

    /**
     * @dev Gets member public key with group and member index numbers.
     * @param _groupIndex Index number of a group.
     * @param _memberIndex Index number of a member.
     */
    function getGroupMemberPubKey(uint256 _groupIndex, uint256 _memberIndex) public view returns(bytes32) {
        return bytes32Storage[keccak256("memberIndexToMemberPubKey", _memberIndex, getGroupPubKey(_groupIndex))];
    }

    /**
     * @dev Emits events with group status, whether it exists or not.
     * @param _groupPubKey Group public key.
     */
    function emitEventGroupExists(bytes32 _groupPubKey) public {
        if (boolStorage[keccak256("groupExists", _groupPubKey)]) {
            emit GroupExistsEvent(_groupPubKey, true);
        } else {
            emit GroupExistsEvent(_groupPubKey, false);
        }
    }

    /**
     * @dev Checks if group is complete.
     * @param _groupPubKey Group public key.
     */
    function groupIsComplete(bytes32 _groupPubKey) public view returns(bool) {
        return boolStorage[keccak256("groupComplete", _groupPubKey)];
    }

    /**
     * @dev Creates a new group with provided public key.
     * @param _groupPubKey Group public key.
     * @return True if group was created, false otherwise.
     */
    function createGroup(bytes32 _groupPubKey) public returns(bool) {

        if (boolStorage[keccak256("groupExists", _groupPubKey)] == true) {
            emit GroupErrorCode(20);
            return false;
        }

        boolStorage[keccak256("groupExists", _groupPubKey)] = true;
        boolStorage[keccak256("groupComplete", _groupPubKey)] = false;
        uintStorage[keccak256("membersCount", _groupPubKey)] = 0;

        uint256 lastIndex = uintStorage[keccak256("groupsCount")];
        bytes32Storage[keccak256("groupIndexToGroupPubKey", lastIndex)] = _groupPubKey;
        uintStorage[keccak256("groupsCount")]++;

        emit GroupStartedEvent(_groupPubKey);
        return true;
    }

    /**
     * @dev Removes a group and the list of its members. Last group public
     * key is moved into the released index and the total group list count
     * is reduced accordingly.
     * @param _groupPubKey Group public key.
     * @return True if group was removed, false otherwise.
     */
    function dissolveGroup(bytes32 _groupPubKey) public onlyOwner returns(bool) {

        // TODO: refine the logic to dissolve a group automaticaly after a certain point
        // More info at https://github.com/keep-network/keep-core/issues/140

        if (boolStorage[keccak256("groupExists", _groupPubKey)] != true) {
            emit GroupErrorCode(10);
            return false;
        }

        for (uint i = 0; i < uintStorage[keccak256("membersCount", _groupPubKey)]; i++) {
            delete bytes32Storage[keccak256("memberIndexToMemberPubKey", i, _groupPubKey)];
        }

        delete uintStorage[keccak256("membersCount", _groupPubKey)];
        delete boolStorage[keccak256("groupExists", _groupPubKey)];
        delete boolStorage[keccak256("groupComplete", _groupPubKey)];

        uint _groupIndex = getGroupIndex(_groupPubKey);
        delete bytes32Storage[keccak256("groupIndexToGroupPubKey", _groupIndex)];

        // Get last group _groupPubKey and move it into released index
        uint lastIndex = uintStorage[keccak256("groupsCount")];
        bytes32 lastGroup = bytes32Storage[keccak256("groupIndexToGroupPubKey", lastIndex)];
        bytes32Storage[keccak256("group", _groupIndex)] = lastGroup;
        uintStorage[keccak256("groupsCount")]--;
    }

    /**
     * @dev Checks if member is part of the group.
     * @param _groupPubKey Group public key.
     * @param _memberPubKey Member public key.
     * @return True if member is part of the group, false otherwise.
     */
    function isMember(bytes32 _groupPubKey, bytes32 _memberPubKey) public view returns(bool) {
        for (uint i = 0; i < uintStorage[keccak256("membersCount", _groupPubKey)]; i++) {
            if (bytes32Storage[keccak256("memberIndexToMemberPubKey", i, _groupPubKey)] == _memberPubKey) {
                return true;
            }
        }
        return false;
    }

    /**
     * @dev Adds member to the group.
     * @param _groupPubKey Group public key.
     * @param _memberPubKey Member public key.
     * @return True if member was added to the group, false otherwise
     * along with emitting corresponding error code.
     */
    function addMemberToGroup(bytes32 _groupPubKey, bytes32 _memberPubKey)
        public
        // hasMinimumStake(_memberPubKey)
        returns(bool)
    {
        // Group does not exist.
        if (boolStorage[keccak256("groupExists", _groupPubKey)] != true) {
            emit GroupErrorCode(3);
            return false;
        }

        // Group is not accepting new members.
        if (boolStorage[keccak256("groupComplete", _groupPubKey)] == true) {
            emit GroupErrorCode(2);
            return false;
        }

        // Member already exists in the group.
        if (isMember(_groupPubKey, _memberPubKey)) {
            emit GroupErrorCode(1);
            return false;
        }

        uint256 lastIndex = uintStorage[keccak256("membersCount", _groupPubKey)];
        bytes32Storage[keccak256("memberIndexToMemberPubKey", lastIndex, _groupPubKey)] = _memberPubKey;
        uintStorage[keccak256("membersCount", _groupPubKey)]++;

        // If the group has passed the threshold size, it is formed.
        if (lastIndex >= uintStorage[keccak256("groupThreshold")]) {
            boolStorage[keccak256("groupComplete", _groupPubKey)] = true;
            emit GroupCompleteEvent(_groupPubKey);
        }
    }
}
