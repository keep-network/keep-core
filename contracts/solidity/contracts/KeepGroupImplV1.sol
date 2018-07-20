pragma solidity ^0.4.24;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";
import "./EternalStorage.sol";
import "./KeepRandomBeaconImplV1.sol";


contract KeepGroupImplV1 is Ownable, EternalStorage {

    event GroupExistsEvent(bytes32 groupPubKey, bool exists);
    event GroupStartedEvent(bytes32 groupPubKey);
    event GroupCompleteEvent(bytes32 groupPubKey);
    event GroupErrorCode(uint8 code);

	bytes32 esKeepGroupImplV1;
	bytes32 esKeepRandomBeaconAddress;
	bytes32 esGroupThreshold;
	bytes32 esGroupSize;
	bytes32 esGroupsCount;
   	bytes32 esGroupIndexToGroupPubKey;
	bytes32 esMemberIndexToMemberPubKey;
	bytes32 esGroupExists;
	bytes32 esGroupComplete;
	bytes32 esMembersCount;
    bytes32 esGroup;
	bytes32 esListOfGroupMemberIDs; 
	bytes32 esNoOfListOfGroupMemberIDs; 

	constructor () {
		esKeepGroupImplV1 = keccak256("KeepGroupImplV1");
		esKeepRandomBeaconAddress = keccak256("keepRandomBeaconAddress");
		esGroupThreshold = keccak256("groupThreshold");
		esGroupSize = keccak256("groupSize");
		esGroupsCount = keccak256("groupsCount");
   		esGroupIndexToGroupPubKey = keccak256("groupIndexToGroupPubKey");
		esMemberIndexToMemberPubKey = keccak256("memberIndexToMemberPubKey");
		esGroupExists = keccak256("groupExists");
		esGroupComplete = keccak256("groupComplete");
		esMembersCount = keccak256("membersCount");
    	esGroup = keccak256("group");
		esListOfGroupMemberIDs =  keccak256("listOfGroupMemberIDs"); 
		esNoOfListOfGroupMemberIDs =  keccak256("noOfListOfGroupMemberIDs"); 
	}

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
		esKeepGroupImplV1 = keccak256("KeepGroupImplV1");
		esKeepRandomBeaconAddress = keccak256("keepRandomBeaconAddress");
		esGroupThreshold = keccak256("groupThreshold");
		esGroupSize = keccak256("groupSize");
		esGroupsCount = keccak256("groupsCount");
   		esGroupIndexToGroupPubKey = keccak256("groupIndexToGroupPubKey");
		esMemberIndexToMemberPubKey = keccak256("memberIndexToMemberPubKey");
		esGroupExists = keccak256("groupExists");
		esGroupComplete = keccak256("groupComplete");
		esMembersCount = keccak256("membersCount");
    	esGroup = keccak256("group");
		esListOfGroupMemberIDs =  keccak256("listOfGroupMemberIDs"); 
		esNoOfListOfGroupMemberIDs =  keccak256("noOfListOfGroupMemberIDs"); 
        boolStorage[esKeepGroupImplV1] = true;
        addressStorage[esKeepRandomBeaconAddress] = _keepRandomBeaconAddress;
        uintStorage[esGroupThreshold] = _groupThreshold;
        uintStorage[esGroupSize] = _groupSize;
        uintStorage[esGroupsCount] = 0;
    }

	// temp: testing functions to verify stuff is getting set correctly.
	function getValue1() public view returns(bytes32) {
		return ( esKeepGroupImplV1 );
	}
	function getValue2() public view returns(bytes32) {
		return ( esGroup );
	}

    /**
     * @dev Allows owner to chagne the groupSize and Threshold.
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
        return bytes32StorageMap[esGroupIndexToGroupPubKey][_groupIndex];
    }

    /**
     * @dev Gets group index number.
     * @param _groupPubKey Group public key.
     */
    function getGroupIndex(bytes32 _groupPubKey) public view returns(uint) {
        for (uint i = 0; i < uintStorage[esGroupsCount]; i++) {
            if (bytes32StorageMap[esGroupIndexToGroupPubKey][i] == _groupPubKey) {
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
        return bytes32StorageMap[esMemberIndexToMemberPubKey][_memberIndex ^ uint256(getGroupPubKey(_groupIndex))];		// this is a problem!!! PJS xyzzy
    }

    /**
     * @dev Emits events with group status, whether it exists or not.
     * @param _groupPubKey Group public key.
     */
    function emitEventGroupExists(bytes32 _groupPubKey) public {
        if (boolStorageMap2[esGroupExists][_groupPubKey]) {
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

        if (boolStorageMap2[esGroupExists][_groupPubKey] == true) {
            emit GroupErrorCode(20);
            return false;
        }

        boolStorageMap2[esGroupExists][_groupPubKey] = true;
        boolStorageMap2[esGroupComplete][_groupPubKey] = false;
        uintStorageMap2[esMembersCount][_groupPubKey] = 0;

        uint256 lastIndex = uintStorage[esGroupsCount];
        bytes32StorageMap[esGroupIndexToGroupPubKey][lastIndex] = _groupPubKey;
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

        if (boolStorageMap2[esGroupExists][_groupPubKey] != true) {
            emit GroupErrorCode(10);
            return false;
        }

        for (uint i = 0; i < uintStorageMap[esMembersCount][uint256(_groupPubKey)]; i++) {
            delete bytes32StorageMap[esMemberIndexToMemberPubKey][ i ^ uint256(_groupPubKey)];		// Problem again xyzzy
        }

        delete uintStorageMap2[esMembersCount][_groupPubKey];
        delete boolStorageMap2[esGroupExists][_groupPubKey];
        delete boolStorageMap2[esGroupComplete][_groupPubKey];

        uint _groupIndex = getGroupIndex(_groupPubKey);
        delete bytes32StorageMap[esGroupIndexToGroupPubKey][_groupIndex];

        // Get last group _groupPubKey and move it into released index
        uint lastIndex = uintStorage[esGroupsCount];
        bytes32 lastGroup = bytes32Storage[esGroupIndexToGroupPubKey][lastIndex];
        bytes32StorageMap[esGroup][_groupIndex] = lastGroup;
        uintStorage[esGroupsCount]--;
    }

    /**
     * @dev Checks if member is part of the group.
     * @param _groupPubKey Group public key.
     * @param _memberPubKey Member public key.
     * @return True if member is part of the group, false otherwise.
     */
    function isMember(bytes32 _groupPubKey, bytes32 _memberPubKey) public view returns(bool) {
        for (uint i = 0; i < uintStorageMap[esMembersCount][uint256(_groupPubKey)]; i++) {
            if (bytes32StorageMap[esMemberIndexToMemberPubKey][ i ^ uint256(_groupPubKey)] == _memberPubKey) {		// Problem again xyzzy
                return true;
            }
        }
        return false;
    }

	// Temporary Code for Milestone 1 follows

    event OnStakerAdded(uint32 index, bytes32 groupMemberID);
	bytes32[] listOfGroupMemberIDs; 

    /**
     * @dev Testing for M1 - create a staker.
     * @param _groupMemberID the ID of the member that is being added.
     */
    function addStaker(bytes32 _groupMemberID) public onlyOwner {
		// TODO save some info at this point - this is only for use in Milestone 1.
		listOfGroupMemberIDs.push( _groupMemberID );
		uint32 index = uint32(listOfGroupMemberIDs.length - 1);
    	emit OnStakerAdded( index, _groupMemberID );
	}

    /**
     * @dev Testing for M1 - return true if the staker at _index is _groupMemberID
     * @param _index Index where to find the member.
     * @param _groupMemberID the ID of the member that is being tested for.
     */
    function isGroupMemberStaker(uint32 _index, bytes32 _groupMemberID) public view returns (bool) {
        require( _index >= 0 && _index < listOfGroupMemberIDs.length );
		return ( listOfGroupMemberIDs[_index] == _groupMemberID );
	}

    /**
     * @dev Testing for M1 - return the groupMemberID for the _index staker.
     * @param _index Index where to add the member.
     */
    function getStaker(uint32 _index) public view returns ( bytes32 ) {
        require( _index >= 0 && _index < listOfGroupMemberIDs.length );
		return ( listOfGroupMemberIDs[_index] );
	}

    /**
     * @dev Testing for M1 - return the number of stakers
     */
    function getNStaker() public view returns ( uint256 ) {
		return ( listOfGroupMemberIDs.length );
	}

    /**
     * @dev Testing for M1 - for testing - reset the array to 0 length.
	 */
    function resetStaker() public onlyOwner {
		delete(listOfGroupMemberIDs);
	}

}
