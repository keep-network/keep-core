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

import "../contracts/KeepRelayBeacon.sol";

contract KeepGroup { 

	uint256 public groupSize;
    address public contractOwner = msg.sender;
    address public keepRelayBeaconAddress;	// address of contract for determining if isStaked.

	struct GroupMemberStruct {
		uint8 memberExists;
		bytes32 PubKey;
	}
	struct GroupStruct {
		uint8 groupExists;
		bool groupComplete;
		bytes32 GroupPubKey;
		mapping (bytes32 => GroupMemberStruct) GroupMembers;
		uint32 nMembers;
		bytes32[] listOfMembers; // list of hashes in group
	}
    mapping (bytes32 => GroupStruct) public groupIdMap;
	uint256 public nGroups;
	bytes32[] listOfGroupIDs; // array of keys into groupIdMap - so can get back a list of groups

    event GroupExistsEvent(bytes32 GroupID, bool Exists);
    event GroupStartedEvent(bytes32 GroupID);
    event GroupCompleteEvent(bytes32 GroupID);
    event GroupErrorCode(uint8 Code);

    /// @dev validates that a call to the function will only succeed if the owner of the contract made the call.
    modifier onlyOwner() {
        require(msg.sender == contractOwner);
        _;
    }

    modifier isStaked() {
		// -- testing - commented out -- //  require(KeepRelayBeacon(keepRelayBeaconAddress).isStaked(msg.sender));
        _;
    }

	/// @dev constructor - set initial group size
    function KeepGroup(uint256 _initialGroupSize, address _keepRelayBeaconAddress) public {
		keepRelayBeaconAddress = _keepRelayBeaconAddress;
		groupSize = _initialGroupSize;
		nGroups = 0;
	}

	/// @dev set size of groups 
    function setGroupSize (uint256 _GroupSize) public onlyOwner {
		groupSize = _GroupSize;
		/// TODO: determine if size decreased, then partially complete groups may now be complete.  Iterate over groups. Find
	}

	/// @dev return number of groups.
    function getNGroups () view public returns(uint256) {
		return nGroups;
	}

	/// @dev fetch back the number of members in a group.
    function getGroupNMembers (uint256 _No) view public returns(uint256) {
		bytes32 gID;
		if ( _No >= 0 && _No < listOfGroupIDs.length && _No < nGroups ) {
			gID = listOfGroupIDs[_No];
			GroupStruct storage ag = groupIdMap[gID];
			if ( ag.groupExists == 1 ) {
				return ( ag.nMembers );
			}
		}
		revert();
	}

	/// @dev fetch back the group pubkey.
    function getGroupPubKey (uint256 _No) view public returns(bytes32) {
		bytes32 gID;
		if ( _No >= 0 && _No < listOfGroupIDs.length && _No < nGroups ) {
			gID = listOfGroupIDs[_No];
			GroupStruct storage ag = groupIdMap[gID];
			if ( ag.groupExists == 1 ) {
				return ( ag.GroupPubKey );
			}
		}
		revert();
	}

	/// @dev get the public key for the _No gorup and the member _MemberNo
    function getGroupMemberPubKey (uint256 _No, uint256 _MemberNo) view public returns(bytes32) {
		bytes32 gID;
		if ( _No >= 0 && _No < listOfGroupIDs.length && _No < nGroups ) {
			gID = listOfGroupIDs[_No];
			GroupStruct storage ag = groupIdMap[gID];
			if ( ag.groupExists == 1 ) {
				if ( _MemberNo >= 0 && _MemberNo < ag.listOfMembers.length && _MemberNo < ag.nMembers ) {
					bytes32 mKey = ag.listOfMembers[_MemberNo];
					GroupMemberStruct storage gm = ag.GroupMembers[mKey];
					if ( gm.memberExists == 1 ) {
						return ( gm.PubKey );
					}
				}
			}
		}
		revert();
	}

	/// @dev find out if a group already exists generating an event.
    function groupExists (bytes32 _gID) public {
		GroupStruct storage ag = groupIdMap[_gID];	
		if ( ag.groupExists == 1 ) {
			GroupExistsEvent(_gID, true);
		} else {
			GroupExistsEvent(_gID, false);
		}
	}

	/// @dev return true if group is complete (has sufficient members)
    function groupIsComplete (bytes32 _gID) view public returns(bool) {
		GroupStruct storage ag = groupIdMap[_gID];	
		return ( ag.groupComplete );
	}

	/// @dev function to check if group exists in a contract.  Returns true if group exits.
    function groupExistsView (bytes32 _gID) view public returns(bool) {
		GroupStruct storage ag = groupIdMap[_gID];	
		if ( ag.groupExists == 1 ) {
			return(true);
		} 
		return(false);
	}

	/// @dev start a new group, save the group bublic key.
    function createGroup (bytes32 _gID) public returns(bool) {
		GroupStruct storage ag = groupIdMap[_gID];	
		if ( ag.groupExists != 1 ) {
			ag.groupExists = 1;				// marker for later so we can see if in map
			ag.groupComplete = false;
			ag.GroupPubKey = _gID;	// save the groups public key
			groupIdMap[_gID] = ag;
			nGroups++;
			listOfGroupIDs.push(_gID);
			GroupStartedEvent(_gID);
			return ( true );
		}
		GroupErrorCode( 20 );
		return ( false );
	}

	/// @dev discard/delete a group if it exists.
	/// @param _gID is the public key that identifies the group.
    function disolveGroup (bytes32 _gID) public onlyOwner returns(bool) {
		GroupStruct storage ag = groupIdMap[_gID];	
		if ( ag.groupExists  == 1 ) {
			delete groupIdMap[_gID];	
			bool done = false;
			for ( uint index = 0; !done && index < nGroups; index++ ) {
				if ( listOfGroupIDs[index] == _gID ) {
					done = true;
					for (uint i = index; i < listOfGroupIDs.length-1; i++){
						listOfGroupIDs[i] = listOfGroupIDs[i+1];
					}
					delete listOfGroupIDs[listOfGroupIDs.length-1];
					listOfGroupIDs.length--;
				}
			}
			nGroups--;
			return ( true );
		} 
		GroupErrorCode( 10 );
		return ( false );
	}

	/// @dev add the transaction sender to the group specified by _gID using the public key for the member of _MemberPubKey.
    function addMemberToGroup (bytes32 _gID, bytes32 _MemberPubKey) public isStaked returns(bool) {
		uint8 rejectedReasonCode = 0;
		GroupStruct storage ag = groupIdMap[_gID];	// fetch back the gorup.
		if ( ag.groupExists == 1 ) {
			if ( ! ag.groupComplete ) {	// if the group is still accepting new members
				// check for unique entry in group?
				GroupMemberStruct storage gm = ag.GroupMembers[_MemberPubKey];
				if ( gm.memberExists == 0 ) {	// if the MemberPubKey is not in the group
					gm.memberExists = 1;	// Mark the MemberPubKey to be a part of the group
					gm.PubKey = _MemberPubKey;	// Save the MemberPubKey in the gorup
					ag.nMembers++;
					ag.listOfMembers.push( _MemberPubKey );	//
					if ( ag.nMembers >= groupSize ) { // if the group has passed the threshold size, it is formed.
						ag.groupComplete = true;
						GroupCompleteEvent(_gID);
					}
				} else {
					rejectedReasonCode = 1;
				}
			} else {
				rejectedReasonCode = 2;
			}
			// groupIdMap[gID] = ag;	// xyzzy - is this necessary??
		} else {
			rejectedReasonCode = 3;
		}
		if ( rejectedReasonCode != 0 ) {
    		GroupErrorCode( rejectedReasonCode );
			return(false);
		}
		return(true);
	}
	   
    /// @dev Accept payments
    function () public payable {
    }
}

