pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./StakingProxy.sol";
import "./TokenStaking.sol";
import "./utils/UintArrayUtils.sol";
import "./utils/AddressArrayUtils.sol";


contract KeepGroupImplV1 is Ownable {

    using SafeMath for uint256;

    event OnGroupRegistered(bytes groupPubKey);

    struct DkgResult {
        bytes groupPubKey;
        bytes disqualified;
        bytes inactive;
    }

    // TODO: Rename to DkgResultSubmittedEvent
    // TODO: Add memberIndex
    event DkgResultPublishedEvent(uint256 requestId, bytes groupPubKey); 
    
    event DkgResultVoteEvent(uint256 requestId, uint256 memberIndex, bytes32 resultHash);

    // Legacy code moved from Random Beacon contract
    // TODO: refactor according to the Phase 14
    event SubmitGroupPublicKeyEvent(bytes groupPublicKey, uint256 requestID);

    event testAddEvent(uint256 c);                                                                                       //<----------

    uint256 internal _groupThreshold;
    uint256 internal _groupSize;
    uint256 internal _minStake;
    address internal _stakingProxy;
    address internal _randomBeacon;

    uint256 internal _timeoutInitial;
    uint256 internal _timeoutSubmission;
    uint256 internal _timeoutChallenge;
    uint256 internal _submissionStart;

    uint256 internal _randomBeaconValue;

    uint256[] internal _tickets;
    bytes[] internal _submissions;

    mapping (uint256 => DkgResult) internal _requestIdToDkgResult;
    mapping (uint256 => bool) internal _dkgResultPublished;
    mapping (bytes => uint256) internal _submissionVotes;
    mapping (address => mapping (bytes => bool)) internal _hasVoted;

    struct Proof {
        address sender;
        uint256 stakerValue;
        uint256 virtualStakerIndex;
    }

    mapping(uint256 => Proof) internal _proofs;

    // _numberOfActiveGroups is the minimal number of groups that should not
    // expired to protect the minimal network throughput.
    uint256 internal _numberOfActiveGroups;
 
    // _groupExpirationTimeout is the time in block after which a group expires.
    uint256 internal _groupExpirationTimeout;
 
    // _expiredOffset is pointing to the first active group, it is also the
    // expired groups counter.
    uint256 internal _expiredOffset = 0;

    // _deletedOffset is pointing to the first not deleted offset, it is also
    // the deleted groups counter.
    uint256 internal _deletedOffset = 0;

    // _expirationThreshold is the number after which batch of groups should be
    // deleted. It is used only by selectGroupV3.
    uint256 internal _expirationThreshold = 50;

    struct Group {
        bytes groupPubKey;
        uint registrationBlockHeight;
    }

    Group[] internal _groups;
    mapping (bytes => address[]) internal _groupMembers;

    mapping (string => bool) internal _initialized;

    /**
     * @dev Triggers the selection process of a new candidate group.
     */
    function runGroupSelection(uint256 randomBeaconValue) public {
        require(msg.sender == _randomBeacon);
        cleanup();
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
     * @param stakerValue Staker-specific value. Currently uint representation of staker address.
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

        // Invalid tickets are rejected and their senders penalized.
        if (!cheapCheck(msg.sender, stakerValue, virtualStakerIndex)) {
            // TODO: replace with a secure authorization protocol (addressed in RFC 4).
            TokenStaking stakingContract = TokenStaking(_stakingContract);
            stakingContract.authorizedTransferFrom(msg.sender, address(this), _minStake);
        } else {
            _tickets.push(ticketValue);
            _proofs[ticketValue] = Proof(msg.sender, stakerValue, virtualStakerIndex);
        }
    }

    /**
     * @dev Gets submitted tickets in ascending order.
     */
    function orderedTickets() public view returns (uint256[] memory) {
        return UintArrayUtils.sort(_tickets);
    }

    /**
     * @dev Gets selected tickets in ascending order.
     */
    function selectedTickets() public view returns (uint256[] memory) {

        require(
            block.number > _submissionStart + _timeoutChallenge,
            "Ticket submission challenge period must be over."
        );

        uint256[] memory ordered = orderedTickets();
        uint256[] memory selected = new uint256[](_groupSize);

        for (uint i = 0; i < _groupSize; i++) {
            selected[i] = ordered[i];
        }

        return selected;
    }

    /**
     * @dev Gets participants ordered by their lowest-valued ticket.
     */
    function orderedParticipants() public view returns (address[] memory) {

        uint256[] memory ordered = orderedTickets();
        address[] memory participants = new address[](ordered.length);

        for (uint i = 0; i < ordered.length; i++) {
            Proof memory proof = _proofs[ordered[i]];
            participants[i] = proof.sender;
        }

        return participants;
    }

    /**
     * @dev Gets selected participants in ascending order of their tickets.
     */
    function selectedParticipants() public view returns (address[] memory) {

        require(
            block.number > _submissionStart + _timeoutChallenge,
            "Ticket submission challenge period must be over."
        );

        uint256[] memory ordered = orderedTickets();
        address[] memory selected = new address[](_groupSize);

        for (uint i = 0; i < _groupSize; i++) {
            Proof memory proof = _proofs[ordered[i]];
            selected[i] = proof.sender;
        }

        return selected;
    }

    /**
     * @dev Gets ticket proof.
     */
    function getTicketProof(uint256 ticketValue) public view returns (address, uint256, uint256) {
        return (
            _proofs[ticketValue].sender,
            _proofs[ticketValue].stakerValue,
            _proofs[ticketValue].virtualStakerIndex
        );
    }

    /**
     * @dev Performs surface-level validation of the ticket.
     * @param staker Address of the staker.
     * @param stakerValue Staker-specific value. Currently uint representation of staker address.
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
     * @param stakerValue Staker-specific value. Currently uint representation of staker address.
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

    /**
     * @dev Submits result of DKG protocol. It is on-chain part of phase 13 of the protocol.
     * @param requestId Relay request ID assosciated with DKG protocol execution.
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
//        uint256 memberIndex, TODO: Add memberIndex 
        bytes memory groupPubKey,
        bytes memory disqualified,
        bytes memory inactive
    ) public {

        require(
            block.number > _submissionStart + _timeoutChallenge,
            "Ticket submission challenge period must be over."
        );

        require(
            _tickets.length >= _groupSize,
            "There should be enough valid tickets submitted to form a group."
        );

        _requestIdToDkgResult[requestId] = DkgResult(groupPubKey, disqualified, inactive);
        _dkgResultPublished[requestId] = true;
  
        emit DkgResultPublishedEvent(requestId, groupPubKey);
    }

    /**
     * @dev Checks if DKG protocol result has been already published for the
     * specific relay request ID associated with the protocol execution. 
     */
    function isDkgResultSubmitted(uint256 requestId) public view returns(bool) {
        return _dkgResultPublished[requestId];
    }

    /*
     * @dev Gets number of votes for each submitted DKG result hash. 
     * @param requestId Relay request ID assosciated with DKG protocol execution.
     * @return Hashes of submitted DKG results and number of votes for each hash.
     */
    function getDkgResultsVotes(uint256 requestId) public view returns (bytes32[] memory, uint256[] memory) {
        // TODO: Implement
        bytes32[] memory resultsHashes;
        uint256[] memory resultsVotes;

        return (resultsHashes, resultsVotes);
    }

    /*
     * @dev receives vote for provided resultHash.
     * @param index the claimed index of the user.
     * @param resultHash Hash of DKG result to vote for
     */
    function voteOnDkgResult(
        uint256 requestId,
        uint256 memberIndex,
        bytes32 resultHash
    ) public {
        // TODO: Implement
    }

    // Legacy code moved from Random Beacon contract		
    // TODO: refactor according to the Phase 14		
    function submitGroupPublicKey(bytes memory groupPublicKey, uint256 requestID) public {

        // TODO: Remove this section once dispute logic is implemented,
        // implement conflict resolution logic described in Phase 14,
        // make sure only valid members are stored.
        _groups.push(Group(groupPublicKey, block.number));
        address[] memory members = orderedParticipants();
        if (members.length > 0) {
            for (uint i = 0; i < _groupSize; i++) {
                _groupMembers[groupPublicKey].push(members[i]);
            }
        }
        emit OnGroupRegistered(groupPublicKey);
        emit SubmitGroupPublicKeyEvent(groupPublicKey, requestID);
    }		

    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() external payable {
        revert("Can not call contract without explicitly calling a function.");
    }

    /**
     * @dev Initialize Keep Group implementation contract with a linked Staking proxy contract.
     * @param stakingProxy Address of a staking proxy contract that will be linked to this contract.
     * @param randomBeacon Address of a random beacon contract that will be linked to this contract.
     * @param minStake Minimum amount in KEEP that allows KEEP network client to participate in a group.
     * @param groupSize Size of a group in the threshold relay.
     * @param groupThreshold Minimum number of interacting group members needed to produce a relay entry.
     * @param timeoutInitial Timeout in blocks after the initial ticket submission is finished.
     * @param timeoutSubmission Timeout in blocks after the reactive ticket submission is finished.
     * @param timeoutChallenge Timeout in blocks after the period where tickets can be challenged is finished.
     * @param numberOfActiveGroups is the minimal number of groups that cannot be marked as expired.
     * @param groupExpirationTimeout is the time in block after which a group expires.
     */
    function initialize(
        address stakingProxy,
        address randomBeacon,
        uint256 minStake,
        uint256 groupThreshold,
        uint256 groupSize,
        uint256 timeoutInitial,
        uint256 timeoutSubmission,
        uint256 timeoutChallenge,
        uint256 groupExpirationTimeout,
        uint256 numberOfActiveGroups
    ) public onlyOwner {
        require(!initialized(), "Contract is already initialized.");
        require(stakingProxy != address(0x0), "Staking proxy address can't be zero.");
        _initialized["KeepGroupImplV1"] = true;
        _stakingProxy = stakingProxy;
        _randomBeacon = randomBeacon;
        _minStake = minStake;
        _groupSize = groupSize;
        _groupThreshold = groupThreshold;
        _timeoutInitial = timeoutInitial;
        _timeoutSubmission = timeoutSubmission;
        _timeoutChallenge = timeoutChallenge;
        _groupExpirationTimeout = groupExpirationTimeout;
        _numberOfActiveGroups = numberOfActiveGroups;
    }

    /**
     * @dev Checks that the specified user has enough stake.
     * @param staker Specifies the identity of the staker.
     * @return True if staked enough to participate in the group, false otherwise.
     */
    function hasMinimumStake(address staker) public view returns(bool) {
        return StakingProxy(_stakingProxy).balanceOf(staker) >= _minStake;
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
     * @dev Allows owner to change the groupSize.
     */
    function setGroupSize(uint256 groupSize) public onlyOwner {
        _groupSize = groupSize;
    }

    /**
     * @dev ticketInitialSubmissionTimeout is the duration (in blocks) the
     * staker has to submit tickets that fall under the natural threshold
     * to satisfy the initial ticket timeout (see group selection, phase 2a).
     */
    function ticketInitialSubmissionTimeout() public view returns (uint256) {
        return _timeoutInitial;
    }

    /**
     * @dev ticketReactiveSubmissionTimeout is the duration (in blocks) the
     * staker has to submit any tickets that did not fall under the natural
     * threshold. This final chance to submit tickets is called reactive
     * ticket submission (defined in the group selection algorithm, 2b).
     */
    function ticketReactiveSubmissionTimeout() public view returns (uint256) {
        return _timeoutSubmission;
    }

    /**
     * @dev ticketChallengeTimeout is the duration (in blocks) the staker
     * has to submit any challenges for tickets that fail any checks.
     */
    function ticketChallengeTimeout() public view returns (uint256) {
        return _timeoutChallenge;
    }

    /**
     * @dev ticketSubmissionStartBlock block number at which current group
     * selection started.
     */
    function ticketSubmissionStartBlock() public view returns (uint256) {
        return _submissionStart;
    }

    /**
     * @dev Return total number of all tokens issued.
     */
    function tokenSupply() public pure returns (uint256) {
        return (10**9) * (10**18);
    }

    /**
     * @dev Return natural threshold, the value N virtual stakers' tickets would be expected
     * to fall below if the tokens were optimally staked, and the tickets' values were evenly 
     * distributed in the domain of the pseudorandom function.
     */
    function naturalThreshold() public view returns (uint256) {
        uint256 space = 2**256-1; // Space consisting of all possible tickets.
        return _groupSize.mul(space.div(tokenSupply().div(_minStake)));
    }

    /**
     * @dev Checks if this contract is initialized.
     */
    function initialized() public view returns (bool) {
        return _initialized["KeepGroupImplV1"];
    }

    /**
     * @dev Gets size of a group in the threshold relay.
     */
    function groupSize() public view returns(uint256) {
        return _groupSize;
    }

    /**
     * @dev Gets number of interacting group members needed to produce a relay entry.
    */
    function groupThreshold() public view returns(uint256) {
        return _groupThreshold;
    }

    /**
     * @dev Adds a new group based on groupPublicKey.
     * @param groupPublicKey is the identifier of the newly created group.
     */
    function groupAdd(bytes memory groupPublicKey) public {
        _groups.push(Group(groupPublicKey, block.number));
        address[] memory members = orderedParticipants();
        if (members.length > 0) {
            for (uint i = 0; i < _groupSize; i++) {
                _groupMembers[groupPublicKey].push(members[i]);
            }
        }
    }

    /**
     * @dev Gets number of active groups.
     */
    function numberOfGroups() public view returns(uint256) {
        return _groups.length - _expiredOffset;
    }

    /*
     * Unmodified selectGroup without deletion. Used for reference purposes.
     */
    function selectGroupV0(uint256 previousEntry) public returns(bytes memory) {
        uint256 activeGroupsNumber = _groups.length - _expiredOffset;
        uint256 selectedGroup = previousEntry % activeGroupsNumber;

        while (_groups[_expiredOffset + selectedGroup].registrationBlockHeight + _groupExpirationTimeout < block.number) {
            if (activeGroupsNumber > _numberOfActiveGroups) {
                if (selectedGroup == 0) {
                    _expiredOffset++;
                    activeGroupsNumber--;
                } else {
                    _expiredOffset += selectedGroup;
                    activeGroupsNumber -= selectedGroup;
                }
                selectedGroup = previousEntry % activeGroupsNumber;
            } else break;
        }

        return _groups[_expiredOffset + selectedGroup].groupPubKey;
    }

    /*
     * First selectGroup with deletion implementation. Expired groups are
     * deleted while marked expired.
     */
    function selectGroupV1(uint256 previousEntry) public returns(bytes memory) {
        uint256 activeGroupsNumber = _groups.length - _expiredOffset;
        uint256 selectedGroup = previousEntry % activeGroupsNumber;
     
        while (_groups[_expiredOffset + selectedGroup].registrationBlockHeight + _groupExpirationTimeout < block.number) {
            if (activeGroupsNumber > _numberOfActiveGroups) {
                if (selectedGroup == 0) {
                    delete _groups[_expiredOffset];
                    _expiredOffset++;
                    activeGroupsNumber--;
                } else {
                    for (uint i = 1; i <= selectedGroup; i++)
                        delete _groups[_expiredOffset++];
                    activeGroupsNumber -= selectedGroup;
                }
                selectedGroup = previousEntry % activeGroupsNumber;
            } else break;
        }

        return _groups[_expiredOffset + selectedGroup].groupPubKey;
    }

    /*
     * Second selectGroup with deletion implementation. Expired groups are
     * deleted in a single batch after an active group is found.
     */
    function selectGroupV2(uint256 previousEntry) public returns(bytes memory) {
        uint256 activeGroupsNumber = _groups.length - _expiredOffset;
        uint256 selectedGroup = previousEntry % activeGroupsNumber;
        uint256 oldOffset = _expiredOffset;

        while (_groups[_expiredOffset + selectedGroup].registrationBlockHeight + _groupExpirationTimeout < block.number) {
            if (activeGroupsNumber > _numberOfActiveGroups) {
                if (selectedGroup == 0) {
                    _expiredOffset++;
                    activeGroupsNumber--;
                } else {
                    _expiredOffset += selectedGroup;
                    activeGroupsNumber -= selectedGroup;
                }
                selectedGroup = previousEntry % activeGroupsNumber;
            } else break;
        }

        for (; oldOffset < _expiredOffset; oldOffset++)
            delete _groups[oldOffset];

        return _groups[_expiredOffset + selectedGroup].groupPubKey;
    }

    /*
     * Third selectGroup with deletion implementation. Expired groups are
     * deleted in a single batch after an active group is found and when an 
     * expirationThreshold is satisfied.
     */
    function selectGroupV3(uint256 previousEntry, uint256 expirationThreshold) public returns(bytes memory) {
        uint256 activeGroupsNumber = _groups.length - _expiredOffset;
        uint256 selectedGroup = previousEntry % activeGroupsNumber;

        while (_groups[_expiredOffset + selectedGroup].registrationBlockHeight + _groupExpirationTimeout < block.number) {
            if (activeGroupsNumber > _numberOfActiveGroups) {
                if (selectedGroup == 0) {
                    _expiredOffset++;
                    activeGroupsNumber--;
                } else {
                    _expiredOffset += selectedGroup;
                    activeGroupsNumber -= selectedGroup;
                }
                selectedGroup = previousEntry % activeGroupsNumber;
            } else break;
        }

        if (_expiredOffset - _deletedOffset > expirationThreshold)
            for (uint i = 0; i < expirationThreshold; i++)
                delete _groups[_deletedOffset++];

        return _groups[_expiredOffset + selectedGroup].groupPubKey;
    }

    /**
     * @dev Returns public key of a group from active groups using modulo operator.
     * @param previousEntry Previous random beacon value.
     */
    function selectGroup(uint256 previousEntry) public returns(bytes memory) {
        uint256 activeGroupsNumber = _groups.length - _expiredOffset;
        uint256 selectedGroup = previousEntry % activeGroupsNumber;

    /**
     * We selected a group based on the information about expired groups offset
     * from the previous call of the function. Now we need to check whether the
     * selected group did not expire in the meantime. To do that, we compare its
     * registration block height and group expiration timeout against the
     * current block number. If the group has expired we move the expired groups
     * offset to the position of the selected expired group and we try to select
     * the next group knowing that all groups before the one previously selected
     * are expired and should not be taken into account. We do this until we
     * find an active group or until we reach the minimum active groups
     * threshold.
     *
     * This approach is more efficient than traversing all groups one by one
     * starting from the previous value of expired groups offset since we can
     * mark expired groups in batches, in a fewer number of steps.
     */
        while (_groups[_expiredOffset + selectedGroup].registrationBlockHeight + _groupExpirationTimeout < block.number) {
            if (activeGroupsNumber > _numberOfActiveGroups) {
                if (selectedGroup == 0) {
                    _expiredOffset++;
                    activeGroupsNumber--;
                } else {
                    _expiredOffset += selectedGroup;
                    activeGroupsNumber -= selectedGroup;
                }
                selectedGroup = previousEntry % activeGroupsNumber;
            } else break;
        }

        return _groups[_expiredOffset + selectedGroup].groupPubKey;
    }

    /**
     * @dev Gets version of the current implementation.
    */
    function version() public pure returns (string memory) {
        return "V1";
    }

    /**
     * @dev Cleanup data of previous group selection.
     */
    function cleanup() private {

        for (uint i = 0; i < _tickets.length; i++) {
            delete _proofs[_tickets[i]];
        }

        delete _tickets;

        // TODO: cleanup DkgResults
    }

}
