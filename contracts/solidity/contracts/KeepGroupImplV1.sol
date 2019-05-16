pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/cryptography/ECDSA.sol";
import "./StakingProxy.sol";
import "./TokenStaking.sol";
import "./utils/UintArrayUtils.sol";
import "./utils/AddressArrayUtils.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";


contract KeepGroupImplV1 is Ownable {

    using SafeMath for uint256;
    using BytesLib for bytes;
    using ECDSA for bytes32;

    event OnGroupRegistered(bytes groupPubKey);

    // TODO: Rename to DkgResultSubmittedEvent
    // TODO: Add memberIndex
    event DkgResultPublishedEvent(uint256 requestId, bytes groupPubKey);

    // TODO: Remove requestId once Keep Client DKG is refactored to
    // use randomBeaconValue as unique id.
    event GroupSelectionStarted(uint256 randomBeaconValue, uint256 requestId, uint256 seed);

    uint256 internal _groupThreshold;
    uint256 internal _groupSize;
    uint256 internal _minStake;
    address internal _stakingProxy;
    address internal _randomBeacon;

    uint256 internal _timeoutInitial;
    uint256 internal _timeoutSubmission;
    uint256 internal _timeoutChallenge;
    uint256 internal _timeDKG;
    uint256 internal _resultPublicationBlockStep;
    uint256 internal _ticketSubmissionStartBlock;
    uint256 internal _randomBeaconValue;

    uint256[] internal _tickets;
    bytes[] internal _submissions;

    // Store whether DKG result was published for the corresponding requestID.
    mapping (uint256 => bool) internal _dkgResultPublished;

    bool internal _groupSelectionInProgress;

    struct Proof {
        address sender;
        uint256 stakerValue;
        uint256 virtualStakerIndex;
    }

    mapping(uint256 => Proof) internal _proofs;

    // _activeGroupsThreshold is the minimal number of groups that should not
    // expire to protect the minimal network throughput.
    // It should be at least 1.
    uint256 internal _activeGroupsThreshold;
 
    // _activeTime is the time in block after which a group expires
    uint256 internal _activeTime;
 
    // _expiredOffset is pointing to the first active group, it is also the
    // expired groups counter
    uint256 internal _expiredOffset = 0;

    struct Group {
        bytes groupPubKey;
        uint registrationBlockHeight;
    }

    Group[] internal _groups;
    
    mapping (bytes => address[]) internal _groupMembers;

    mapping (string => bool) internal _initialized;


    /**
     * @dev Checks if submitter is eligible to submit.
     * @param submitterMemberIndex The claimed index of the submitter.
     */
    modifier onlyEligibleSubmitter(uint256 submitterMemberIndex) {
        uint256[] memory selected = selectedTickets();
        require(submitterMemberIndex > 0, "Submitter member index must be greater than 0.");
        require(_proofs[selected[submitterMemberIndex - 1]].sender == msg.sender, "Submitter member index does not match sender address.");
        uint T_init = _ticketSubmissionStartBlock + _timeoutChallenge + _timeDKG;
        require(block.number >= (T_init + (submitterMemberIndex-1) * _resultPublicationBlockStep), "Submitter is not eligible to submit at the current block.");
        _;
    }

    /**
     * @dev Reverts if ticket challenge period is not over.
     */
    modifier whenTicketChallengeIsOver() {
        require(
            block.number >= _ticketSubmissionStartBlock + _timeoutChallenge,
            "Ticket submission challenge period must be over."
        );
        _;
    }

    /**
     * @dev Triggers the selection process of a new candidate group.
     */
    function runGroupSelection(uint256 randomBeaconValue, uint256 requestId, uint256 seed) public {
        require(msg.sender == _randomBeacon);

        uint256 latestDKGSubmission = _ticketSubmissionStartBlock + _timeoutChallenge + _timeDKG + _groupSize * _resultPublicationBlockStep;

        if (!_groupSelectionInProgress || block.number > latestDKGSubmission) {
            cleanup();
            _ticketSubmissionStartBlock = block.number;
            _randomBeaconValue = randomBeaconValue;
            _groupSelectionInProgress = true;
            emit GroupSelectionStarted(randomBeaconValue, requestId, seed);
        }
    }

    // TODO: replace with a secure authorization protocol (addressed in RFC 4).
    address internal _stakingContract;

    function authorizeStakingContract(address stakingContract) public onlyOwner {
        _stakingContract = stakingContract;
    }

    /**
     * @dev Submits ticket to request to participate in a new candidate group.
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

        if (block.number > _ticketSubmissionStartBlock + _timeoutSubmission) {
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
    function selectedTickets() public view whenTicketChallengeIsOver returns (uint256[] memory) {

        uint256[] memory ordered = orderedTickets();

        require(
            ordered.length >= _groupSize,
            "The number of submitted tickets is less than specified group size."
        );

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
    function selectedParticipants() public view whenTicketChallengeIsOver returns (address[] memory) {

        uint256[] memory ordered = orderedTickets();

        require(
            ordered.length >= _groupSize,
            "The number of submitted tickets is less than specified group size."
        );

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
     * @dev Checks if member is disqualified.
     * @param dqBytes bytes representing disqualified members.
     * @param memberIndex position of the member to check.
     * @return true if staker is disqualified, false otherwise.
     */
    function _isDisqualified(bytes memory dqBytes, uint256 memberIndex) internal pure returns (bool){
        return dqBytes[memberIndex] != 0x00;
    }

     /**
     * @dev Checks if member is inactive.
     * @param iaBytes bytes representing inactive members.
     * @param memberIndex position of the member to check.
     * @return true if staker is inactive, false otherwise.
     */
    function _isInactive(bytes memory iaBytes, uint256 memberIndex) internal pure returns (bool){
        return iaBytes[memberIndex] != 0x00;
    }

    /**
     * @dev Submits result of DKG protocol. It is on-chain part of phase 14 of the protocol.
     * @param submitterMemberIndex Claimed index of the staker. We pass this for gas efficiency purposes.
     * @param requestId Relay request ID associated with DKG protocol execution.
     * @param groupPubKey Group public key generated as a result of protocol execution.
     * @param disqualified bytes representing disqualified group members; 1 at the specific index
     * means that the member has been disqualified. Indexes reflect positions of members in the
     * group, as outputted by the group selection protocol.
     * @param inactive bytes representing inactive group members; 1 at the specific index means
     * that the member has been marked as inactive. Indexes reflect positions of members in the
     * group, as outputted by the group selection protocol.
     * @param signatures Concatenation of signed resultHashes collected off-chain.
     * @param signingMembersIndexes indices of members corresponding to each signature.
     */
    function submitDkgResult(
        uint256 requestId,
        uint256 submitterMemberIndex,
        bytes memory groupPubKey,
        bytes memory disqualified,
        bytes memory inactive,
        bytes memory signatures,
        uint[] memory signingMembersIndexes
    ) public onlyEligibleSubmitter(submitterMemberIndex) {

        require(
            disqualified.length == _groupSize && inactive.length == _groupSize,
            "Inactive and disqualified bytes arrays don't match the group size."
        );

        require(
            !_dkgResultPublished[requestId], 
            "DKG result for this request ID already published."
        );

        bytes32 resultHash = keccak256(abi.encodePacked(groupPubKey, disqualified, inactive));
        verifySignatures(signatures, signingMembersIndexes, resultHash);
        address[] memory members = selectedParticipants();

        for (uint i = 0; i < _groupSize; i++) {
            if(!_isInactive(inactive, i) && !_isDisqualified(disqualified, i)) {
                _groupMembers[groupPubKey].push(members[i]);
            }
        }

        _groups.push(Group(groupPubKey, block.number));
        // TODO: punish/reward logic
        cleanup();
        _dkgResultPublished[requestId] = true;
        emit DkgResultPublishedEvent(requestId, groupPubKey);

        _groupSelectionInProgress = false;
    }

    /**
    * @dev Verifies that provided members signatures of the DKG result were produced
    * by the members stored previously on-chain in the order of their ticket values
    * and returns indices of members with a boolean value of their signature validity.
    * @param signatures Concatenation of user-generated signatures.
    * @param resultHash The result hash signed by the users.
    * @param signingMemberIndices Indices of members corresponding to each signature.
    * @return Array of member indices with a boolean value of their signature validity.
    */
    function verifySignatures(
        bytes memory signatures,
        uint256[] memory signingMemberIndices,
        bytes32 resultHash
    ) internal view returns (bool) {

        uint256 signaturesCount = signatures.length / 65;
        require(signatures.length >= 65, "Signatures bytes array is too short.");
        require(signatures.length % 65 == 0, "Signatures in the bytes array should be 65 bytes long.");
        require(signaturesCount == signingMemberIndices.length, "Number of signatures and indices don't match.");
        require(signaturesCount >= _groupThreshold, "Number of signatures is below honest majority threshold.");

        bytes memory current; // Current signature to be checked.
        uint256[] memory selected = selectedTickets();

        for(uint i = 0; i < signaturesCount; i++){
            require(signingMemberIndices[i] > 0, "Index should be greater than zero.");
            require(signingMemberIndices[i] <= selected.length, "Provided index is out of acceptable tickets bound.");
            current = signatures.slice(65*i, 65);
            address recoveredAddress = resultHash.toEthSignedMessageHash().recover(current);

            require(
                _proofs[selected[signingMemberIndices[i] - 1]].sender == recoveredAddress,
                "Invalid signature. Signer and recovered address at provided index don't match."
            );
        }

        return true;
    }

    /**
     * @dev Checks if DKG protocol result has been already published for the
     * specific relay request ID associated with the protocol execution. 
     */
    function isDkgResultSubmitted(uint256 requestId) public view returns(bool) {
        return _dkgResultPublished[requestId];
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
     * @param timeDKG Timeout in blocks after DKG result is complete and ready to be published.
     * @param resultPublicationBlockStep Time in blocks after which member with the given index is eligible
     * to submit DKG result.
     * @param activeGroupsThreshold is the minimal number of groups that cannot be marked as expired and
     * needs to be greater than 0.
     * @param activeTime is the time in block after which a group expires.
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
        uint256 timeDKG,
        uint256 resultPublicationBlockStep,
        uint256 activeGroupsThreshold,
        uint256 activeTime
    ) public onlyOwner {
        require(!initialized(), "Contract is already initialized.");
        require(stakingProxy != address(0x0), "Staking proxy address can't be zero.");
        require(activeGroupsThreshold > 0, "The minumum number of active groups needs to be more than zero.");
        _initialized["KeepGroupImplV1"] = true;
        _stakingProxy = stakingProxy;
        _randomBeacon = randomBeacon;
        _minStake = minStake;
        _groupSize = groupSize;
        _groupThreshold = groupThreshold;
        _timeoutInitial = timeoutInitial;
        _timeoutSubmission = timeoutSubmission;
        _timeoutChallenge = timeoutChallenge;
        _timeDKG = timeDKG;
        _resultPublicationBlockStep = resultPublicationBlockStep;
        _activeGroupsThreshold = activeGroupsThreshold;
        _activeTime = activeTime;
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
     * @dev resultPublicationBlockStep is the duration (in blocks) after which
     * member with the given index is eligible to submit DKG result.
     */
    function resultPublicationBlockStep() public view returns (uint256) {
        return _resultPublicationBlockStep;
    }

    /**
     * @dev ticketSubmissionStartBlock block number at which current group
     * selection started.
     */
    function ticketSubmissionStartBlock() public view returns (uint256) {
        return _ticketSubmissionStartBlock;
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
     * @dev Gets number of active groups.
     */
    function numberOfGroups() public view returns(uint256) {
        return _groups.length - _expiredOffset;
    }

    /**
     * @dev Gets the random beacon value for the group selection currently in progress.
     */
    function randomBeaconValue() public view returns(uint256) {
        return _randomBeaconValue;
    }

    /**
     * @dev Returns true if the group selection is in progress.
     */
    function groupSelectionInProgress() public view returns(bool) {
        return _groupSelectionInProgress;
    }

    /**
     * @dev Checks if a group with the given public key is registered.
     */
    function isGroupRegistered(bytes memory groupPubKey) public view returns(bool) {
        for (uint i = 0; i < numberOfGroups(); i++) {
            if (_groups[i].groupPubKey.equalStorage(groupPubKey)) {
                return true;
            }
        }

        return false;
    }

    /**
     * @dev Returns public key of a group from active groups using modulo operator.
     * @param previousEntry Previous random beacon value.
     */
    function selectGroup(uint256 previousEntry) public returns(bytes memory) {
        uint256 numberOfActiveGroups = _groups.length - _expiredOffset;
        uint256 selectedGroup = previousEntry % numberOfActiveGroups;

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
        if (numberOfActiveGroups > _activeGroupsThreshold) {
            while (_groups[_expiredOffset + selectedGroup].registrationBlockHeight + _activeTime < block.number) {
                /**
                * We do -1 to see how many groups are available after the potential removal.
                * For example:
                * _groups = [EEEAAAA]
                * - assuming selectedGroup = 0, then we'll have 4-0-1=3 groups after the removal: [EEEEAAA]
                * - assuming selectedGroup = 1, then we'll have 4-1-1=2 groups after the removal: [EEEEEAA]
                * - assuming selectedGroup = 2, then, we'll have 4-2-1=1 groups after the removal: [EEEEEEA]
                * - assuming selectedGroup = 3, then, we'll have 4-3-1=0 groups after the removal: [EEEEEEE]
                */
                if (numberOfActiveGroups - selectedGroup - 1 > _activeGroupsThreshold) {
                    selectedGroup++;
                    _expiredOffset += selectedGroup;
                    numberOfActiveGroups -= selectedGroup;
                    selectedGroup = previousEntry % numberOfActiveGroups;
                } else {
                    /* Number of groups that did not expire is less or equal _activeGroupsThreshold
                    * and we have more groups than _activeGroupsThreshold (including those expired) groups.
                    * Hence, we maintain the minimum _activeGroupsThreshold of active groups and
                    * do not let any other groups to expire
                    */
                    _expiredOffset = _groups.length - _activeGroupsThreshold;
                    numberOfActiveGroups = _activeGroupsThreshold;
                    selectedGroup = previousEntry % numberOfActiveGroups;
                    break;
                }
            }
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
