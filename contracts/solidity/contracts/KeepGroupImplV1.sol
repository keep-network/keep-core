pragma solidity ^0.4.21;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";
import "./EternalStorage.sol";
import "./KeepRandomBeaconImplV1.sol";


contract KeepGroupImplV1 is Ownable, EternalStorage {

    event GroupExistsEvent(bytes32 groupPubKey, bool exists);
    event GroupStartedEvent(bytes32 groupPubKey);
    event GroupCompleteEvent(bytes32 groupPubKey);
    event GroupErrorCode(uint8 code);
    event OnStakerAdded(uint32 index, bytes32 groupMemberID);

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
     * @dev Gets the threshold size for groups.
     */
    function groupThreshold() public view returns(uint256) {
        return uintStorage[keccak256("groupThreshold")];
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
     * @dev Testing for M1 - create a staker.
     * @param _index Index where to add the member.
     * @param _groupMemberID the ID of the member that is being added.
     */
    function addStaker(uint32 _index, bytes32 _groupMemberID) public {
		// TODO save some info at this point - this is only for use in Milestone 1 and will
		// not need to be added to the "forever" storage.
    	emit OnStakerAdded(_index, _groupMemberID);
	}

}
