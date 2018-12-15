pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./StakingProxy.sol";
import "./TokenStaking.sol";


contract KeepGroupImplV1 is Ownable {

    event GroupExistsEvent(bytes32 groupPubKey, bool exists);
    event GroupStartedEvent(bytes32 groupPubKey);
    event GroupCompleteEvent(bytes32 groupPubKey);
    event GroupErrorCode(uint8 code);

    uint256 internal _groupThreshold;
    uint256 internal _groupSize;
    uint256 internal _groupsCount;
    uint256 internal _minStake;
    address internal _stakingProxy;

    uint256 internal _timeoutInitial;
    uint256 internal _timeoutSubmission;
    uint256 internal _timeoutChallenge;
    uint256 internal _submissionStart;

    uint256 internal _randomBeaconValue;

    uint256[] internal _tickets;

    struct Proof {
        address sender;
        uint256 stakerValue;
        uint256 virtualStakerIndex;
    }

    mapping(uint256 => Proof) internal _proofs;

    mapping (uint256 => bytes32) internal _groupIndexToGroupPubKey;
    mapping (bytes32 => mapping (uint256 => bytes32)) internal _memberIndexToMemberPubKey;
    mapping (bytes32 => bool) internal _groupExists;
    mapping (bytes32 => bool) internal _groupComplete;
    mapping (bytes32 => uint256) internal _membersCount;
    mapping (string => bool) internal _initialized;

    /**
     * @dev Triggers the selection process of a new candidate group.
     */
    function runGroupSelection(uint256 randomBeaconValue) public onlyOwner {
        _submissionStart = block.number;
        _randomBeaconValue = randomBeaconValue;
    }

    // TODO: replace with a secure authorization protocol (addressed in RFC 4).
    address internal _stakingContract;

    function authorizeStakingContract(address stakingContract) public onlyOwner {
        _stakingContract = stakingContract;
    }

    /**
     * @dev Submit ticket to request to participate in a new candidate group.
     * @param ticketValue Result of a pseudorandom function with input values of
     * random beacon output, staker-specific 'stakerValue' and virtualStakerIndex.
     * @param stakerValue Staker-specific value.
     * @param virtualStakerIndex Number within a range of 1 to staker's weight.
     */
    function submitTicket(
        uint256 ticketValue,
        uint256 stakerValue,
        uint256 virtualStakerIndex
    ) public {

        if (block.number > _submissionStart + _timeoutSubmission) {
            revert("Ticket submission period is over.");
        }

        if (block.number > _submissionStart + _timeoutInitial && _tickets.length > _groupSize) {
            revert("Initial submission period is over with enough tickets received.");
        }
 
        _tickets.push(ticketValue);
        _proofs[ticketValue] = Proof(msg.sender, stakerValue, virtualStakerIndex);

    }

    /**
     * @dev Gets ticket proof.
     */
    function getTicketProof(uint256 ticketValue) public view returns (uint256, uint256) {
        return (
            _proofs[ticketValue].stakerValue,
            _proofs[ticketValue].virtualStakerIndex
        );
    }

    /**
     * @dev Performs surface-level validation of the ticket.
     * @param staker Address of the staker.
     * @param stakerValue Staker-specific value.
     * @param virtualStakerIndex Number within a range of 1 to staker's weight.
     */
    function cheapCheck(
        address staker,
        uint256 stakerValue,
        uint256 virtualStakerIndex
    ) public view returns(bool) {
        bool isVirtualStakerIndexValid = virtualStakerIndex > 0 && virtualStakerIndex <= stakingWeight(staker);
        bool isStakerValueValid = uint256(staker) == stakerValue;
        return isVirtualStakerIndexValid && isStakerValueValid;
    }

    /**
     * @dev Performs full verification of the ticket.
     * @param staker Address of the staker.
     * @param ticketValue Result of a pseudorandom function with input values of
     * random beacon output, staker-specific 'stakerValue' and virtualStakerIndex.
     * @param stakerValue Staker-specific value.
     * @param virtualStakerIndex Number within a range of 1 to staker's weight.
     */
    function costlyCheck(
        address staker,
        uint256 ticketValue,
        uint256 stakerValue,
        uint256 virtualStakerIndex
    ) public view returns(bool) {
        bool passedCheapCheck = cheapCheck(staker, stakerValue, virtualStakerIndex);
        uint256 expected = uint256(keccak256(abi.encodePacked(_randomBeaconValue, stakerValue, virtualStakerIndex)));
        return passedCheapCheck && ticketValue == expected;
    }

    function challenge(
        uint256 ticketValue
    ) public {

        Proof memory proof = _proofs[ticketValue];
        require(proof.sender != 0, "Ticket must be published.");

        // TODO: replace with a secure authorization protocol (addressed in RFC 4).
        TokenStaking stakingContract = TokenStaking(_stakingContract);
        if (costlyCheck(
            proof.sender,
            ticketValue,
            proof.stakerValue,
            proof.virtualStakerIndex
        )) {
            // Slash challenger's stake balance.
            stakingContract.authorizedTransferFrom(msg.sender, this, _minStake);
        } else {
            // Slash invalid ticket sender stake balance and reward the challenger.
            stakingContract.authorizedTransferFrom(proof.sender, msg.sender, _minStake);
        }
    }

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
     * @dev Initialize Keep Group implementation contract with a linked Staking proxy contract.
     * @param stakingProxy Address of a staking proxy contract that will be linked to this contract.
     * @param minStake Minimum amount in KEEP that allows KEEP network client to participate in a group.
     * @param groupThreshold Max number of bad members in a group that we can detect as well as “number
     * of good members needed to produce a relay entry”.
     * @param groupSize Minimum number of members in a group - to form a group.
     * @param timeoutInitial Timeout in blocks after the initial ticket submission is finished.
     * @param timeoutSubmission Timeout in blocks after the reactive ticket submission is finished.
     * @param timeoutChallenge Timeout in blocks after the period where tickets can be challenged is finished.
     */
    function initialize(
        address stakingProxy,
        uint256 minStake,
        uint256 groupThreshold,
        uint256 groupSize,
        uint256 timeoutInitial,
        uint256 timeoutSubmission,
        uint256 timeoutChallenge
    ) public onlyOwner {
        require(!initialized(), "Contract is already initialized.");
        require(stakingProxy != address(0x0), "Staking proxy address can't be zero.");
        _initialized["KeepGroupImplV1"] = true;
        _stakingProxy = stakingProxy;
        _minStake = minStake;
        _groupThreshold = groupThreshold;
        _groupSize = groupSize;
        _groupsCount = 0;
        _timeoutInitial = timeoutInitial;
        _timeoutSubmission = timeoutSubmission;
        _timeoutChallenge = timeoutChallenge;
    }

    /**
     * @dev Gets staking weight.
     * @param staker Specifies the identity of the staker.
     * @return Number of how many virtual stakers can staker represent.
     */
    function stakingWeight(address staker) public view returns(uint256) {
        return StakingProxy(_stakingProxy).balanceOf(staker)/_minStake;
    }

    /**
     * @dev Set the minimum amount of KEEP that allows a Keep network client to participate in a group.
     * @param minStake Amount in KEEP.
     */
    function setMinimumStake(uint256 minStake) public onlyOwner {
        _minStake = minStake;
    }

    /**
     * @dev Get the minimum amount in KEEP that allows KEEP network client to participate in a group.
     */
    function minimumStake() public view returns(uint256) {
        return _minStake;
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
    function addStaker(bytes32 groupMemberID) public onlyOwner {
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
