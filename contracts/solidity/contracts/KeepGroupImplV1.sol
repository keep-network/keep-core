pragma solidity ^0.4.18;
// pragma experimental ABIEncoderV2;


// Test Plan
//  1. Create it , pass it in the KeepRelayBeacon
//	2. Greate Group
//		a. Verify Event
//	3. Add to it
//		a. Check that member got added
//		b. 
// 

/// @title Group Management
/// @author Philip Schlump

// interface KeepRelayBeacon { 
//     function isStaked(address _staker) view public returns(bool);
// }

import "zeppelin-solidity/contracts/ownership/Ownable.sol";
import "./EternalStorage.sol";
import "./KeepRandomBeaconImplV1.sol";


contract KeepGroupImplV1 is Ownable {

    event GroupExistsEvent(bytes32 groupPubKey, bool exists);
    event GroupStartedEvent(bytes32 groupPubKey);
    event GroupCompleteEvent(bytes32 groupPubKey);
    event GroupErrorCode(uint8 code);

    modifier isStaked() {
        // -- testing - commented out -- //  require(KeepRelayBeacon(keepRandomBeaconAddress).isStaked(msg.sender));
        _;
    }

    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() public payable {
        revert();
    }

    /**
     * @dev Initialize Keep Keep Group implementaion contract with a linked Keep Random Beacon contract.
     * @param _keepRandomBeaconAddress Address of Keep Random Beacon that will be linked to this contract.
     * @param _initialGroupSize Initial group size.
     */
    function initialize(uint256 _groupThreshold, address _keepRandomBeaconAddress) public {
        addressStorage[keccak256("keepRandomBeaconAddress")] = _keepRandomBeaconAddress;
        uintStorage[keccak256("groupThreshold")] = _groupThreshold;
        uintStorage[keccak256("groupsCount")] = 0;
    }

    /// @dev set size of groups
    function setGroupThreshold(uint256 _groupThreshold) public onlyOwner {
        uintStorage[keccak256("groupThreshold")] = _groupThreshold;
        /// TODO: determine if size decreased, then partially complete groups may now be complete.  Iterate over groups. Find
    }

    /// @dev return number of groups.
    function getNumberOfGroups() public view returns(uint256) {
        return uintStorage[keccak256("groupsCount")];
    }

    /// @dev fetch back the number of members in a group.
    function getGroupNMembers(uint256 _i) public view returns(uint256) {
        return uintStorage[keccak256("membersCount", getGroupPubKey(_i)];
    }

    /// @dev fetch back the group pubkey.
    function getGroupPubKey(uint256 _i) public view returns(bytes32) {
        return byteStorage[keccak256("groupToIndex", _i)];
    }

    function getGroupNumber(bytes32 _groupPubKey) public view returns(uint) {
        for (uint i = 0; i < uintStorage[keccak256("groupsCount")]; i++) {
            if (bytesStorage[keccak256("groupToIndex", i)] == _groupPubKey) {
                return i;
            }
        }
        revert();
    }

    /// @dev get the public key for the _no gorup and the member _memberNo
    function getGroupMemberPubKey(uint256 _i, uint256 _j) public view returns(bytes32) {
        return bytesStorage[keccak256("memberToIndex", _j, getGroupPubKey(_i)]
    }

    /// @dev find out if a group already exists generating an event.
    function groupExists(bytes32 _groupPubKey) public {
        if (boolStorage[keccak256("groupExists", _groupPubKey)]) {
            GroupExistsEvent(_groupPubKey, true);
        } else {
            GroupExistsEvent(_groupPubKey, false);
        }
    }

    /// @dev return true if group is complete (has sufficient members)
    function groupIsComplete(bytes32 _groupPubKey) public view returns(bool) {
        return boolStorage[keccak256("groupComplete", _groupPubKey)]
    }

    /// @dev function to check if group exists in a contract. Returns true if group exits.
    function groupExistsView(bytes32 _groupPubKey) public view returns(bool) {
        return boolStorage[keccak256("groupExists", _groupPubKey)];
    }

    /// @dev start a new group, save the group bublic key.
    function createGroup(bytes32 _groupPubKey) public returns(bool) {

        if (boolStorage[keccak256("groupExists", _groupPubKey)] == true) {
            GroupErrorCode(20);
            return false
        }

        boolStorage[keccak256("groupExists", _groupPubKey)] = true;
        boolStorage[keccak256("groupComplete", _groupPubKey)] = false;
        uintStorage[keccak256("membersCount", _groupPubKey)] = 0;

        uintStorage[keccak256("groupsCount")]++;
        uint256 lastIndex = uintStorage[keccak256("groupsCount")];
        byteStorage[keccak256("groupToIndex", lastIndex)] == _groupPubKey;

        GroupStartedEvent(_groupPubKey);
        return true;
    }

    /// @dev discard/delete a group if it exists.
    /// @param _groupPubKey is the public key that identifies the group.
    function disolveGroup(bytes32 _groupPubKey) public onlyOwner returns(bool) {

        if (boolStorage[keccak256("groupExists", _groupPubKey)] != true) {
            GroupErrorCode(10);
            return false
        }

        for (uint256 index = 0; index < uintStorage[keccak256("membersCount", _groupPubKey)]; index++) {
            delete bytesStorage[keccak256("memberToIndex", index, _groupPubKey)];
        }

        delete uintStorage[keccak256("membersCount", _groupPubKey)];
        delete boolStorage[keccak256("groupExists", _groupPubKey)];
        delete boolStorage[keccak256("groupComplete", _groupPubKey)];

        uint i = getGroupNumber(_groupPubKey);
        delete bytesStorage[keccak256("groupToIndex", i);

        // Get last group _groupPubKey and move it into released index
        uint groupsCount = uintStorage[keccak256("groupsCount")];
        byte32 lastGroup = byteStorage[keccak256("groupToIndex", groupsCount)];
        bytesStorage[keccak256("group", i)] = lastGroup;
        uintStorage[keccak256("groupsCount")]--;
    }

    // Chek if member is part of a group
    function isMember(bytes32 _groupPubKey, bytes32 _memberPubKey) public view returns(bool) {
        for (uint i = 0; i < uintStorage[keccak256("membersCount", _groupPubKey)]; i++) {
            if (bytesStorage[keccak256("memberToIndex", i, _groupPubKey)] == _memberPubKey) {
                return true;
            }
        }
        return false;
    }

    /// @dev add the transaction sender to the group specified by _groupPubKey using the public key for the member of _memberPubKey.
    function addMemberToGroup(bytes32 _groupPubKey, bytes32 _memberPubKey) public isStaked returns(bool) {

        // Group does not exist.
        if (boolStorage[keccak256("groupExists", _groupPubKey)] != true) {
            GroupErrorCode(3);
            return false;
        }

        // Group is not accepting new members.
        if (boolStorage[keccak256("groupComplete", _groupPubKey)] != true) {
            GroupErrorCode(2);
            return false;
        }

        // Member already exists in the group.
        if (isMember(_groupPubKey, _memberPubKey) {
            GroupErrorCode(1);
            return false;
        }

        uintStorage[keccak256("membersCount", _groupPubKey)]++;
        uint256 lastIndex = uintStorage[keccak256("membersCount", _groupPubKey)];
        addressStorage[keccak256("memberToIndex", lastIndex, _groupPubKey)] == _memberPubKey;

        // If the group has passed the threshold size, it is formed.
        if (membersCount >= uintStorage[keccak256("groupThreshold")]) {
            boolStorage[keccak256("groupComplete", _groupPubKey)] = true;
            GroupCompleteEvent(_groupPubKey);
        }
    }
}
