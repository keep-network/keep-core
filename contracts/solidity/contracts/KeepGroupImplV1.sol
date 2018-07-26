pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./EternalStorage.sol";
import "./KeepRandomBeaconImplV1.sol";


contract KeepGroupImplV1 is Ownable, EternalStorage {

    event GroupExistsEvent(bytes32 groupPubKey, bool exists);
    event GroupStartedEvent(bytes32 groupPubKey);
    event GroupCompleteEvent(bytes32 groupPubKey);
    event GroupErrorCode(uint8 code);

    bytes32 private constant esKeepGroupImplV1 = keccak256("KeepGroupImplV1");
    bytes32 private constant esKeepRandomBeaconAddress = keccak256("keepRandomBeaconAddress");
    bytes32 private constant esGroupThreshold = keccak256("groupThreshold");
    bytes32 private constant esGroupSize = keccak256("groupSize");
    bytes32 private constant esGroupsCount = keccak256("groupsCount");
    bytes32 private constant esGroupIndexToGroupPubKey = keccak256("groupIndexToGroupPubKey");
    bytes32 private constant esMemberIndexToMemberPubKey = keccak256("memberIndexToMemberPubKey");
    bytes32 private constant esGroupExists = keccak256("groupExists");
    bytes32 private constant esGroupComplete = keccak256("groupComplete");
    bytes32 private constant esMembersCount = keccak256("membersCount");
    bytes32 private constant esGroup = keccak256("group");
    bytes32 private constant esListOfGroupMemberIDs = keccak256("listOfGroupMemberIDs");
    bytes32 private constant esNoOfListOfGroupMemberIDs = keccak256("noOfListOfGroupMemberIDs");
    // Temporary Code for Milestone 1 follows
    bytes32 private constant esListOfGroupMembersIDs = keccak256("ListOfGroupMembersIDs");
    bytes32 private constant esListOfGroupMembersIDsCount = keccak256("ListOfGroupMembersIDsCount");
    // End Temporary Code for Milestone 1

    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() public payable {
        revert("Can not call contract without explicitly calling a function.");
    }

    /**
     * @dev Initialize Keep Group implementation contract with a linked Keep Random Beacon contract.
     * @param _keepRandomBeaconAddress Address of Keep Random Beacon that will be linked to this contract.
     * @param _groupThreshold Max number of bad members in a group that we can detect as well as “number
     * of good members needed to produce a relay entry”.
     * @param _groupSize Minimum number of members in a group - to form a group.
     */
    function initialize(uint256 _groupThreshold, uint256 _groupSize, address _keepRandomBeaconAddress) public onlyOwner {
        require(!initialized(), "Contract is already initialized.");
        require(_keepRandomBeaconAddress != address(0x0), "Random Beacon address can't be zero.");
        boolStorage[esKeepGroupImplV1] = true;
        addressStorage[esKeepRandomBeaconAddress] = _keepRandomBeaconAddress;
        uintStorage[esGroupThreshold] = _groupThreshold;
        uintStorage[esGroupSize] = _groupSize;
        uintStorage[esGroupsCount] = 0;
		// Temporary Code for Milestone 1 follows
        uintStorage[esListOfGroupMembersIDsCount] = 0;
    	// End Temporary Code for Milestone 1
    }

    /**
     * @dev Allows owner to change the groupSize and Threshold.
     */
    function setGroupSizeThreshold ( uint256 _groupSize, uint256 _groupThreshold ) public onlyOwner {
        uintStorage[esGroupThreshold] = _groupThreshold;
        uintStorage[esGroupSize] = _groupSize;
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return boolStorage[esKeepGroupImplV1];
    }

    /**
     * @dev Gets the threshold size for groups.
     */
    function groupThreshold() public view returns(uint256) {
        return uintStorage[esGroupThreshold];
    }

    /**
     * @dev Gets the minimum number of members in a group.
     */
    function groupSize() public view returns(uint256) {
        return uintStorage[esGroupSize];
    }

    /**
     * @dev Gets number of active groups.
     */
    function numberOfGroups() public view returns(uint256) {
        return uintStorage[esGroupsCount];
    }

    /**
     * @dev Gets public key of a group.
     * @param _groupIndex Index number of a group.
     */
    function getGroupPubKey(uint256 _groupIndex) public view returns(bytes32) {
        return bytes32UintStorageMap[esGroupIndexToGroupPubKey][_groupIndex];
    }

    /**
     * @dev Gets group index number.
     * @param _groupPubKey Group public key.
     */
    function getGroupIndex(bytes32 _groupPubKey) public view returns(uint) {
        for (uint i = 0; i < uintStorage[esGroupsCount]; i++) {
            if (bytes32UintStorageMap[esGroupIndexToGroupPubKey][i] == _groupPubKey) {
                return i;
            }
        }
        revert("Group index is not found.");
    }

    /**
     * @dev Gets member public key with group and member index numbers.
     * @param _groupIndex Index number of a group.
     * @param _memberIndex Index number of a member.
     */
    function getGroupMemberPubKey(uint256 _groupIndex, uint256 _memberIndex) public view returns(bytes32) {
        return bytes32StorageMap[esMemberIndexToMemberPubKey][keccak256(
            abi.encodePacked(_memberIndex, uint256(getGroupPubKey(_groupIndex))))];
    }

    /**
     * @dev Emits events with group status, whether it exists or not.
     * @param _groupPubKey Group public key.
     */
    function emitEventGroupExists(bytes32 _groupPubKey) public {
        if (boolBytes32StorageMap[esGroupExists][_groupPubKey]) {
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

        if (boolBytes32StorageMap[esGroupExists][_groupPubKey] == true) {
            emit GroupErrorCode(20);
            return false;
        }

        boolBytes32StorageMap[esGroupExists][_groupPubKey] = true;
        boolBytes32StorageMap[esGroupComplete][_groupPubKey] = false;
        uintBytes32StorageMap[esMembersCount][_groupPubKey] = 0;

        uint256 lastIndex = uintStorage[esGroupsCount];
        bytes32UintStorageMap[esGroupIndexToGroupPubKey][lastIndex] = _groupPubKey;
        uintStorage[esGroupsCount]++;

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

        if (boolBytes32StorageMap[esGroupExists][_groupPubKey] != true) {
            emit GroupErrorCode(10);
            return false;
        }

        for (uint i = 0; i < uintBytes32StorageMap[esMembersCount][_groupPubKey]; i++) {
            delete bytes32StorageMap[esMemberIndexToMemberPubKey][keccak256(abi.encodePacked(i, uint256(_groupPubKey)))];
        }

        delete uintBytes32StorageMap[esMembersCount][_groupPubKey];
        delete boolBytes32StorageMap[esGroupExists][_groupPubKey];
        delete boolBytes32StorageMap[esGroupComplete][_groupPubKey];

        uint _groupIndex = getGroupIndex(_groupPubKey);
        delete bytes32UintStorageMap[esGroupIndexToGroupPubKey][_groupIndex];

        // Get last group _groupPubKey and move it into released index
        uint lastIndex = uintStorage[esGroupsCount];
        bytes32 lastGroup = bytes32Storage[esGroupIndexToGroupPubKey][lastIndex];
        bytes32UintStorageMap[esGroup][_groupIndex] = lastGroup;
        uintStorage[esGroupsCount]--;
    }

    /**
     * @dev Checks if member is part of the group.
     * @param _groupPubKey Group public key.
     * @param _memberPubKey Member public key.
     * @return True if member is part of the group, false otherwise.
     */
    function isMember(bytes32 _groupPubKey, bytes32 _memberPubKey) public view returns(bool) {
        for (uint i = 0; i < uintBytes32StorageMap[esMembersCount][_groupPubKey]; i++) {
            if (bytes32StorageMap[esMemberIndexToMemberPubKey][keccak256(
                abi.encodePacked(i, uint256(_groupPubKey)))] == _memberPubKey) {
                return true;
            }
        }
        return false;
    }

    // Temporary Code for Milestone 1 follows

    event OnStakerAdded(uint32 index, bytes32 groupMemberID);

    /**
     * @dev Testing for M1 - create a staker.
     * @param _groupMemberID the ID of the member that is being added.
     */
    function addStaker(bytes32 _groupMemberID) public onlyOwner {
        // TODO save some info at this point - this is only for use in Milestone 1 and will
        // not need to be added to the "forever" storage.
        uint32 count = uint32(uintStorage[esListOfGroupMembersIDsCount]);
     	bytes32StorageArray[esListOfGroupMembersIDs].push(_groupMemberID);
		count = count + 1;
        uintStorage[esListOfGroupMembersIDsCount] = count;
        emit OnStakerAdded(count, _groupMemberID);
    }

    /**
     * @dev Testing for M1 - return true if the staker at _index is _groupMemberID
     * @param _index Index where to find the member.
     * @param _groupMemberID the ID of the member that is being tested for.
     */
    function isGroupMemberStaker(uint32 _index, bytes32 _groupMemberID) public view returns (bool) {
        require(_index >= 0 && _index <= uintStorage[esListOfGroupMembersIDsCount], "Index must be within the length of Group member's array.");
     	return ( bytes32StorageArray[esListOfGroupMembersIDs][_index] == _groupMemberID);
    }

    /**
     * @dev Testing for M1 - return the groupMemberID for the _index staker.
     * @param _index Index where to add the member.
     */
    function getStaker(uint32 _index) public view returns (bytes32) {
        require(_index >= 0 && _index <= uintStorage[esListOfGroupMembersIDsCount], "Index must be within the length of Group member's array.");
        return (bytes32StorageArray[esListOfGroupMembersIDs][_index]);
    }

    /**
     * @dev Testing for M1 - return the number of stakers
     */
    function getNStaker() public view returns (uint256) {
        return (uintStorage[esListOfGroupMembersIDsCount]);
    }

    /**
     * @dev Testing for M1 - for testing - reset the array to 0 length.
     */
    function resetStaker() public onlyOwner {
        uintStorage[esListOfGroupMembersIDsCount] = 0;
        delete( bytes32StorageArray[esListOfGroupMembersIDs] );
    }

}
