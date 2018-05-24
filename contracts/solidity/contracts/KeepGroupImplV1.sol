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

    uint256 public groupSize;
    address public keepRandomBeaconAddress; // address of contract for determining if isStaked.

    struct GroupMemberStruct {
        uint8 memberExists;
        bytes32 pubKey;
    }

    struct GroupStruct {
        uint8 groupExists;
        bool groupComplete;
        bytes32 groupPubKey;
        mapping (bytes32 => GroupMemberStruct) groupMembers;
        uint32 nMembers;
        bytes32[] listOfMembers; // list of hashes in group
    }

    mapping (bytes32 => GroupStruct) public groupIdMap;
    uint256 public nGroups;
    bytes32[] listOfGroupIDs; // array of keys into groupIdMap - so can get back a list of groups

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
    function initialize(uint256 _initialGroupSize, address _keepRandomBeaconAddress) public {
        keepRandomBeaconAddress = _keepRandomBeaconAddress;
        groupSize = _initialGroupSize;
        nGroups = 0;
    }

    /// @dev set size of groups 
    function setGroupSize(uint256 _groupSize) public onlyOwner {
        groupSize = _groupSize;
        /// TODO: determine if size decreased, then partially complete groups may now be complete.  Iterate over groups. Find
    }

    /// @dev return number of groups.
    function getNGroups() public view returns(uint256) {
        return nGroups;
    }

    /// @dev fetch back the number of members in a group.
    function getGroupNMembers(uint256 _no) public view returns(uint256) {
        bytes32 gID;
        if (_no >= 0 && _no < listOfGroupIDs.length && _no < nGroups) {
            gID = listOfGroupIDs[_no];
            GroupStruct storage ag = groupIdMap[gID];
            if (ag.groupExists == 1) {
                return (ag.nMembers);
            }
        }
        revert();
    }

    /// @dev fetch back the group pubkey.
    function getGroupPubKey(uint256 _no) public view returns(bytes32) {
        bytes32 gID;
        if (_no >= 0 && _no < listOfGroupIDs.length && _no < nGroups) {
            gID = listOfGroupIDs[_no];
            GroupStruct storage ag = groupIdMap[gID];
            if (ag.groupExists == 1) {
                return (ag.groupPubKey);
            }
        }
        revert();
    }

    /// @dev get the public key for the _no gorup and the member _memberNo
    function getGroupMemberPubKey(uint256 _no, uint256 _memberNo) public view returns(bytes32) {
        bytes32 gID;
        if (_no >= 0 && _no < listOfGroupIDs.length && _no < nGroups) {
            gID = listOfGroupIDs[_no];
            GroupStruct storage ag = groupIdMap[gID];
            if (ag.groupExists == 1) {
                if (_memberNo >= 0 && _memberNo < ag.listOfMembers.length && _memberNo < ag.nMembers) {
                    bytes32 mKey = ag.listOfMembers[_memberNo];
                    GroupMemberStruct storage gm = ag.groupMembers[mKey];
                    if (gm.memberExists == 1) {
                        return (gm.pubKey);
                    }
                }
            }
        }
        revert();
    }

    /// @dev find out if a group already exists generating an event.
    function groupExists(bytes32 _groupPubKey) public {
        GroupStruct storage ag = groupIdMap[_groupPubKey];
        if (ag.groupExists == 1) {
            GroupExistsEvent(_groupPubKey, true);
        } else {
            GroupExistsEvent(_groupPubKey, false);
        }
    }

    /// @dev return true if group is complete (has sufficient members)
    function groupIsComplete(bytes32 _groupPubKey) public view returns(bool) {
        GroupStruct storage ag = groupIdMap[_groupPubKey];
        return (ag.groupComplete);
    }

    /// @dev function to check if group exists in a contract. Returns true if group exits.
    function groupExistsView(bytes32 _groupPubKey) public view returns(bool) {
        GroupStruct storage ag = groupIdMap[_groupPubKey];
        if (ag.groupExists == 1) {
            return(true);
        } 
        return(false);
    }

    /// @dev start a new group, save the group bublic key.
    function createGroup(bytes32 _groupPubKey) public returns(bool) {
        GroupStruct storage ag = groupIdMap[_groupPubKey];
        if (ag.groupExists != 1) {
            ag.groupExists = 1; // marker for later so we can see if in map
            ag.groupComplete = false;
            ag.groupPubKey = _groupPubKey; // save the groups public key
            groupIdMap[_groupPubKey] = ag;
            nGroups++;
            listOfGroupIDs.push(_groupPubKey);
            GroupStartedEvent(_groupPubKey);
            return (true);
        }
        GroupErrorCode(20);
        return (false);
    }

    /// @dev discard/delete a group if it exists.
    /// @param _groupPubKey is the public key that identifies the group.
    function disolveGroup(bytes32 _groupPubKey) public onlyOwner returns(bool) {
        GroupStruct storage ag = groupIdMap[_groupPubKey];	
        if (ag.groupExists == 1) {
            delete groupIdMap[_groupPubKey];
            bool done = false;
            for (uint index = 0; !done && index < nGroups; index++) {
                if (listOfGroupIDs[index] == _groupPubKey) {
                    done = true;
                    for (uint i = index; i < listOfGroupIDs.length-1; i++) {
                        listOfGroupIDs[i] = listOfGroupIDs[i+1];
                    }
                    delete listOfGroupIDs[listOfGroupIDs.length-1];
                    listOfGroupIDs.length--;
                }
            }
            nGroups--;
            return (true);
        } 
        GroupErrorCode(10);
        return (false);
    }

    /// @dev add the transaction sender to the group specified by _groupPubKey using the public key for the member of _memberPubKey.
    function addMemberToGroup(bytes32 _groupPubKey, bytes32 _memberPubKey) public isStaked returns(bool) {
        uint8 rejectedReasonCode = 0;
        GroupStruct storage ag = groupIdMap[_groupPubKey];	// fetch back the gorup.
        if (ag.groupExists == 1) {
            if (!ag.groupComplete) {	// if the group is still accepting new members
                // check for unique entry in group?
                GroupMemberStruct storage gm = ag.groupMembers[_memberPubKey];
                if (gm.memberExists == 0) { // if the MemberPubKey is not in the group
                    gm.memberExists = 1; // Mark the MemberPubKey to be a part of the group
                    gm.pubKey = _memberPubKey; // Save the MemberPubKey in the gorup
                    ag.nMembers++;
                    ag.listOfMembers.push(_memberPubKey);
                    if (ag.nMembers >= groupSize) { // if the group has passed the threshold size, it is formed.
                        ag.groupComplete = true;
                        GroupCompleteEvent(_groupPubKey);
                    }
                } else {
                    rejectedReasonCode = 1;
                }
            } else {
                rejectedReasonCode = 2;
            }
            // groupIdMap[gID] = ag; // xyzzy - is this necessary??
        } else {
            rejectedReasonCode = 3;
        }
        if (rejectedReasonCode != 0) {
            GroupErrorCode(rejectedReasonCode);
            return(false);
        }
        return(true);
    }
}
