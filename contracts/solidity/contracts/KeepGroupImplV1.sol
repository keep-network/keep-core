pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./KeepRandomBeaconImplV1.sol";


contract KeepGroupImplV1 is Ownable {

    struct DkgResult {
        bool success;
        bytes32 groupPubKey;
        bytes disqualified;
        bytes inactive;
    }

    event DkgResultPublishedEvent(uint256 requestId);
    event GroupExistsEvent(bytes32 groupPubKey, bool exists);
    event GroupStartedEvent(bytes32 groupPubKey);
    event GroupCompleteEvent(bytes32 groupPubKey);
    event GroupErrorCode(uint8 code);

    address internal _keepRandomBeaconAddress;
    uint256 internal _groupThreshold;
    uint256 internal _groupSize;
    uint256 internal _groupsCount;

    mapping (uint256 => DkgResult) internal _requestIdToDkgResult;
    mapping (uint256 => bool) internal _dkgResultPublished;
    mapping (uint256 => bytes32) internal _groupIndexToGroupPubKey;
    mapping (bytes32 => mapping (uint256 => bytes32)) internal _memberIndexToMemberPubKey;
    mapping (bytes32 => bool) internal _groupExists;
    mapping (bytes32 => bool) internal _groupComplete;
    mapping (bytes32 => uint256) internal _membersCount;
    mapping (string => bool) internal _initialized;

    // Temporary Code for Milestone 1 follows
    bytes32[] private _listOfGroupMembersIDs;
    // End Temporary Code for Milestone 1

    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() public payable {
        revert("Can not call contract without explicitly calling a function.");
    }

    /**
     * @dev Initialize Keep Group implementation contract with a linked Keep Random Beacon contract.
     * @param keepRandomBeaconAddress Address of Keep Random Beacon that will be linked to this contract.
     * @param groupThreshold Max number of bad members in a group that we can detect as well as “number
     * of good members needed to produce a relay entry”.
     * @param groupSize Minimum number of members in a group - to form a group.
     */
    function initialize(uint256 groupThreshold, uint256 groupSize, address keepRandomBeaconAddress) public onlyOwner {
        require(!initialized(), "Contract is already initialized.");
        require(keepRandomBeaconAddress != address(0x0), "Random Beacon address can't be zero.");
        _initialized["KeepGroupImplV1"] = true;
        _keepRandomBeaconAddress = keepRandomBeaconAddress;
        _groupThreshold = groupThreshold;
        _groupSize = groupSize;
        _groupsCount = 0;
    }

    /**
     * @dev Allows owner to change the groupSize and Threshold.
     */
    function setGroupSizeThreshold(uint256 groupSize, uint256 groupThreshold) public onlyOwner {
        _groupThreshold = groupThreshold;
        _groupSize = groupSize;
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return _initialized["KeepGroupImplV1"];
    }

    /**
     * @dev Gets the threshold size for groups.
     */
    function groupThreshold() public view returns(uint256) {
        return _groupThreshold;
    }

    /**
     * @dev Gets the minimum number of members in a group.
     */
    function groupSize() public view returns(uint256) {
        return _groupSize;
    }

    /**
     * @dev Gets number of active groups.
     */
    function numberOfGroups() public view returns(uint256) {
        return _groupsCount;
    }

    /**
     * @dev Gets public key of a group.
     * @param groupIndex Index number of a group.
     */
    function getGroupPubKey(uint256 groupIndex) public view returns(bytes32) {
        return _groupIndexToGroupPubKey[groupIndex];
    }

    /**
     * @dev Gets group index number.
     * @param groupPubKey Group public key.
     */
    function getGroupIndex(bytes32 groupPubKey) public view returns(uint) {
        for (uint i = 0; i < _groupsCount; i++) {
            if (_groupIndexToGroupPubKey[i] == groupPubKey) {
                return i;
            }
        }
        revert("Group index is not found.");
    }

    /**
     * @dev Gets member public key with group and member index numbers.
     * @param groupIndex Index number of a group.
     * @param memberIndex Index number of a member.
     */
    function getGroupMemberPubKey(uint256 groupIndex, uint256 memberIndex) public view returns(bytes32) {
        return _memberIndexToMemberPubKey[getGroupPubKey(groupIndex)][memberIndex];
    }

    /**
     * @dev Emits events with group status, whether it exists or not.
     * @param groupPubKey Group public key.
     */
    function emitEventGroupExists(bytes32 groupPubKey) public {
        if (_groupExists[groupPubKey]) {
            emit GroupExistsEvent(groupPubKey, true);
        } else {
            emit GroupExistsEvent(groupPubKey, false);
        }
    }

    /**
     * @dev Submits result of DKG protocol. It is on-chain part of phase 13 of the protocol.
     * @param requestId Relay request ID assosciated with DKG protocol execution.
     * @param success Result of DKG protocol execution; true if success, false otherwise.
     * @param groupPubKey Group public key generated as a result of protocol execution.
     * @param disqualified bytes representing disqualified group members; 1 at the specific index 
     * means that the member has been disqualified. Indexes reflect positions of members in the
     * group, as outputted by the group selection protocol.
     * @param inactive bytes representing inactive group members; 1 at the specific index means
     * that the member has been marked as inactive. Indexes reflect positions of members in the
     * group, as outputted by the group selection protocol.
     */
    function submitDkgResult(
        uint256 requestId,
        bool success, 
        bytes32 groupPubKey,
        bytes disqualified,
        bytes inactive
    ) public {
        _requestIdToDkgResult[requestId] = DkgResult(success, groupPubKey, disqualified, inactive);
        _dkgResultPublished[requestId] = true;
  
        emit DkgResultPublishedEvent(requestId);
    }

    /**
     * @dev Checks if DKG protocol result has been already published for the
     * specific relay request ID associated with the protocol execution. 
     */
    function isDkgResultSubmitted(uint256 requestId) public view returns(bool) {
        return _dkgResultPublished[requestId];
    }

    /**
     * @dev Creates a new group with provided public key.
     * @param groupPubKey Group public key.
     * @return True if group was created, false otherwise.
     */
    function createGroup(bytes32 groupPubKey) public returns(bool) {

        if (_groupExists[groupPubKey] == true) {
            emit GroupErrorCode(20);
            return false;
        }

        _groupExists[groupPubKey] = true;
        _groupComplete[groupPubKey] = false;
        _membersCount[groupPubKey] = 0;

        _groupIndexToGroupPubKey[_groupsCount] = groupPubKey;
        _groupsCount++;

        emit GroupStartedEvent(groupPubKey);
        return true;
    }

    /**
     * @dev Removes a group and the list of its members. Last group public
     * key is moved into the released index and the total group list count
     * is reduced accordingly.
     * @param groupPubKey Group public key.
     * @return True if group was removed, false otherwise.
     */
    function dissolveGroup(bytes32 groupPubKey) public onlyOwner returns(bool) {

        if (_groupExists[groupPubKey] != true) {
            emit GroupErrorCode(10);
            return false;
        }

        for (uint i = 0; i < _membersCount[groupPubKey]; i++) {
            delete _memberIndexToMemberPubKey[groupPubKey][i];
        }

        delete _membersCount[groupPubKey];
        delete _groupExists[groupPubKey];
        delete _groupComplete[groupPubKey];

        uint groupIndex = getGroupIndex(groupPubKey);
        delete _groupIndexToGroupPubKey[groupIndex];

        // Get last group _groupPubKey and move it into released index
        uint lastIndex = _groupsCount;
        bytes32 lastGroup = _groupIndexToGroupPubKey[lastIndex];
        _groupIndexToGroupPubKey[groupIndex] = lastGroup;
        _groupsCount--;
    }

    /**
     * @dev Checks if member is part of the group.
     * @param groupPubKey Group public key.
     * @param memberPubKey Member public key.
     * @return True if member is part of the group, false otherwise.
     */
    function isMember(bytes32 groupPubKey, bytes32 memberPubKey) public view returns(bool) {
        for (uint i = 0; i < _membersCount[groupPubKey]; i++) {
            if (_memberIndexToMemberPubKey[groupPubKey][i] == memberPubKey) {
                return true;
            }
        }
        return false;
    }

    // Temporary Code for Milestone 1 follows

    event OnStakerAdded(uint32 index, bytes32 groupMemberID);

    /**
     * @dev Testing for M1 - create a staker.
     * @param groupMemberID the ID of the member that is being added.
     */
    function addStaker(bytes32 groupMemberID) public {
        // TODO save some info at this point - this is only for use in Milestone 1 and will
        // not need to be added to the "forever" storage.
        _listOfGroupMembersIDs.push(groupMemberID);
        emit OnStakerAdded(uint32(_listOfGroupMembersIDs.length - 1), groupMemberID);
    }

    /**
     * @dev Testing for M1 - return true if the staker at _index is _groupMemberID
     * @param index Index where to find the member.
     * @param groupMemberID the ID of the member that is being tested for.
     */
    function isGroupMemberStaker(uint32 index, bytes32 groupMemberID) public view returns (bool) {
        require(
            index >= 0 && index <= _listOfGroupMembersIDs.length,
            "Index must be within the length of Group member's array."
        );
        return _listOfGroupMembersIDs[index] == groupMemberID;
    }

    /**
     * @dev Testing for M1 - return the groupMemberID for the _index staker.
     * @param index Index where to add the member.
     */
    function getStaker(uint32 index) public view returns (bytes32) {
        require(
            index >= 0 && index <= _listOfGroupMembersIDs.length,
            "Index must be within the length of Group member's array."
        );
        return _listOfGroupMembersIDs[index];
    }

    /**
     * @dev Testing for M1 - return the number of stakers
     */
    function getNStaker() public view returns (uint256) {
        return _listOfGroupMembersIDs.length;
    }

    /**
     * @dev Testing for M1 - for testing - reset the array to 0 length.
     */
    function resetStaker() public onlyOwner {
        delete _listOfGroupMembersIDs;
    }

    /**
     * @dev Gets version of the current implementation.
    */
    function version() public pure returns (string) {
        return "V1";
    }
}
