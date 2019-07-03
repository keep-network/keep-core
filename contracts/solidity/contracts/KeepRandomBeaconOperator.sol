pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/cryptography/ECDSA.sol";
import "./StakingProxy.sol";
import "./TokenStaking.sol";
import "./utils/UintArrayUtils.sol";
import "./utils/AddressArrayUtils.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";
import "./cryptography/BLS.sol";

interface ServiceContract {
    function entryCreated(uint256 requestId, uint256 entry) external;
}

/**
 * @title KeepRandomBeaconOperator
 * @dev Keep client facing contract for random beacon security-critical operations.
 * Handles group creation and expiration, BLS signature verification and incentives.
 * The contract is not upgradeable. New functionality can be implemented by deploying
 * new versions following Keep client update and re-authorization by the stakers.
 */
contract KeepRandomBeaconOperator is Ownable {

    using SafeMath for uint256;
    using BytesLib for bytes;
    using ECDSA for bytes32;
    using AddressArrayUtils for address[];

    event OnGroupRegistered(bytes groupPubKey);

    // TODO: Rename to DkgResultSubmittedEvent
    // TODO: Add memberIndex
    event DkgResultPublishedEvent(uint256 signingId, bytes groupPubKey);

    // These are the public events that are used by clients
    event SignatureRequested(uint256 signingId, uint256 payment, uint256 previousEntry, uint256 seed, bytes groupPublicKey);
    event SignatureSubmitted(uint256 signingId, uint256 requestResponse, bytes requestGroupPubKey, uint256 previousEntry, uint256 seed);

    // TODO: Remove signingId once Keep Client DKG is refactored to
    // use groupSelectionSeed as unique id.
    event GroupSelectionStarted(uint256 groupSelectionSeed, uint256 signingId, uint256 seed);

    uint256 public signingRequestCounter;
    uint256 public groupThreshold;
    uint256 public groupSize;
    uint256 public minimumStake;
    address public stakingProxy;
    address[] public serviceContracts;

    uint256 public ticketInitialSubmissionTimeout;
    uint256 public ticketReactiveSubmissionTimeout;
    uint256 public ticketChallengeTimeout;
    uint256 public timeDKG;
    uint256 public resultPublicationBlockStep;
    uint256 public ticketSubmissionStartBlock;
    uint256 public groupSelectionSeed;

    uint256[] public tickets;
    bytes[] public submissions;

    // Store whether DKG result was published for the corresponding signingId.
    mapping (uint256 => bool) public dkgResultPublished;

    bool public groupSelectionInProgress;

    struct Proof {
        address sender;
        uint256 stakerValue;
        uint256 virtualStakerIndex;
    }

    mapping(uint256 => Proof) public proofs;

    // activeGroupsThreshold is the minimal number of groups that should not
    // expire to protect the minimal network throughput.
    // It should be at least 1.
    uint256 public activeGroupsThreshold;
 
    // groupActiveTime is the time in block after which a group expires
    uint256 public groupActiveTime;

    // Timeout in blocks for a relay entry to appear on the chain. Blocks are
    // counted from the moment relay request occur.
    uint256 public relayRequestTimeout;

    // expiredGroupOffset is pointing to the first active group, it is also the
    // expired groups counter
    uint256 public expiredGroupOffset = 0;

    struct Group {
        bytes groupPubKey;
        uint registrationBlockHeight;
    }

    Group[] public groups;

    mapping (bytes => address[]) internal groupMembers;

    bool public initialized;

    struct SigningRequest {
        uint256 requestId;
        uint256 payment;
        bytes groupPubKey;
        address serviceContract;
    }

    mapping(uint256 => SigningRequest) public signingRequests;

    /**
     * @dev Checks if submitter is eligible to submit.
     * @param submitterMemberIndex The claimed index of the submitter.
     */
    modifier onlyEligibleSubmitter(uint256 submitterMemberIndex) {
        uint256[] memory selected = selectedTickets();
        require(submitterMemberIndex > 0, "Submitter member index must be greater than 0.");
        require(proofs[selected[submitterMemberIndex - 1]].sender == msg.sender, "Submitter member index does not match sender address.");
        uint T_init = ticketSubmissionStartBlock + ticketChallengeTimeout + timeDKG;
        require(block.number >= (T_init + (submitterMemberIndex-1) * resultPublicationBlockStep), "Submitter is not eligible to submit at the current block.");
        _;
    }

    /**
     * @dev Reverts if ticket challenge period is not over.
     */
    modifier whenTicketChallengeIsOver() {
        require(
            block.number >= ticketSubmissionStartBlock + ticketChallengeTimeout,
            "Ticket submission challenge period must be over."
        );
        _;
    }

    /**
     * @dev Initialize the contract with a linked Staking proxy contract.
     * @param _stakingProxy Address of a staking proxy contract that will be linked to this contract.
     * @param _serviceContract Address of a random beacon service contract that will be linked to this contract.
     * @param _minimumStake Minimum amount in KEEP that allows KEEP network client to participate in a group.
     * @param _groupSize Size of a group in the threshold relay.
     * @param _groupThreshold Minimum number of interacting group members needed to produce a relay entry.
     * @param _ticketInitialSubmissionTimeout Timeout in blocks after the initial ticket submission is finished.
     * @param _ticketReactiveSubmissionTimeout Timeout in blocks after the reactive ticket submission is finished.
     * @param _ticketChallengeTimeout Timeout in blocks after the period where tickets can be challenged is finished.
     * @param _timeDKG Timeout in blocks after DKG result is complete and ready to be published.
     * @param _resultPublicationBlockStep Time in blocks after which member with the given index is eligible
     * @param _genesisEntry Initial relay entry to create first group.
     * @param _genesisGroupPubKey Group to respond to the initial relay entry request.
     * to submit DKG result.
     * @param _activeGroupsThreshold is the minimal number of groups that cannot be marked as expired and
     * needs to be greater than 0.
     * @param _groupActiveTime is the time in block after which a group expires.
     * @param _relayRequestTimeout Timeout in blocks for a relay entry to appear on the chain.
     * Blocks are counted from the moment relay request occur.
     */
    function initialize(
        address _stakingProxy,
        address _serviceContract,
        uint256 _minimumStake,
        uint256 _groupThreshold,
        uint256 _groupSize,
        uint256 _ticketInitialSubmissionTimeout,
        uint256 _ticketReactiveSubmissionTimeout,
        uint256 _ticketChallengeTimeout,
        uint256 _timeDKG,
        uint256 _resultPublicationBlockStep,
        uint256 _activeGroupsThreshold,
        uint256 _groupActiveTime,
        uint256 _relayRequestTimeout,
        uint256 _genesisEntry,
        bytes memory _genesisGroupPubKey
    ) public onlyOwner {
        require(!initialized, "Contract is already initialized.");
        require(_stakingProxy != address(0x0), "Staking proxy address can't be zero.");
        initialized = true;
        stakingProxy = _stakingProxy;
        serviceContracts.push(_serviceContract);
        minimumStake = _minimumStake;
        groupSize = _groupSize;
        groupThreshold = _groupThreshold;
        ticketInitialSubmissionTimeout = _ticketInitialSubmissionTimeout;
        ticketReactiveSubmissionTimeout = _ticketReactiveSubmissionTimeout;
        ticketChallengeTimeout = _ticketChallengeTimeout;
        timeDKG = _timeDKG;
        resultPublicationBlockStep = _resultPublicationBlockStep;
        activeGroupsThreshold = _activeGroupsThreshold;
        groupActiveTime = _groupActiveTime;
        relayRequestTimeout = _relayRequestTimeout;
        groupSelectionSeed = _genesisEntry;

        // Create initial relay entry request. This will allow relayEntry to be called once
        // to trigger the creation of the first group. Requests are removed on successful
        // entries so genesis entry can only be called once.
        signingRequestCounter++;
        signingRequests[signingRequestCounter] = SigningRequest(0, 0, _genesisGroupPubKey, _serviceContract);
    }

    /**
     * @dev Triggers the selection process of a new candidate group.
     * @param _groupSelectionSeed Random value that stakers will use to generate their tickets.
     * @param _signingId Relay request ID associated with DKG protocol execution.
     * @param _seed Random value from the client. It should be a cryptographically generated random value.
     */
    function createGroup(uint256 _groupSelectionSeed, uint256 _signingId, uint256 _seed) private {
        // dkgTimeout is the time after DKG is expected to be complete plus the expected period to submit the result.
        uint256 dkgTimeout = ticketSubmissionStartBlock + ticketChallengeTimeout + timeDKG + groupSize * resultPublicationBlockStep;

        if (!groupSelectionInProgress || block.number > dkgTimeout) {
            cleanup();
            ticketSubmissionStartBlock = block.number;
            groupSelectionSeed = _groupSelectionSeed;
            groupSelectionInProgress = true;
            emit GroupSelectionStarted(_groupSelectionSeed, _signingId, _seed);
        }
    }

    // TODO: replace with a secure authorization protocol (addressed in RFC 4).
    address public stakingContract;

    function authorizeStakingContract(address _stakingContract) public onlyOwner {
        stakingContract = _stakingContract;
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

        if (block.number > ticketSubmissionStartBlock + ticketReactiveSubmissionTimeout) {
            revert("Ticket submission period is over.");
        }

        // Invalid tickets are rejected and their senders penalized.
        if (!cheapCheck(msg.sender, stakerValue, virtualStakerIndex)) {
            // TODO: replace with a secure authorization protocol (addressed in RFC 4).
            TokenStaking _stakingContract = TokenStaking(stakingContract);
            _stakingContract.authorizedTransferFrom(msg.sender, address(this), minimumStake);
        } else {
            tickets.push(ticketValue);
            proofs[ticketValue] = Proof(msg.sender, stakerValue, virtualStakerIndex);
        }
    }

    /**
     * @dev Gets submitted tickets in ascending order.
     */
    function orderedTickets() public view returns (uint256[] memory) {
        return UintArrayUtils.sort(tickets);
    }

    /**
     * @dev Gets selected tickets in ascending order.
     */
    function selectedTickets() public view whenTicketChallengeIsOver returns (uint256[] memory) {

        uint256[] memory ordered = orderedTickets();

        require(
            ordered.length >= groupSize,
            "The number of submitted tickets is less than specified group size."
        );

        uint256[] memory selected = new uint256[](groupSize);

        for (uint i = 0; i < groupSize; i++) {
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
            Proof memory proof = proofs[ordered[i]];
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
            ordered.length >= groupSize,
            "The number of submitted tickets is less than specified group size."
        );

        address[] memory selected = new address[](groupSize);

        for (uint i = 0; i < groupSize; i++) {
            Proof memory proof = proofs[ordered[i]];
            selected[i] = proof.sender;
        }

        return selected;
    }

    /**
     * @dev Gets ticket proof.
     */
    function getTicketProof(uint256 ticketValue) public view returns (address sender, uint256 stakerValue, uint256 virtualStakerIndex) {
        return (
            proofs[ticketValue].sender,
            proofs[ticketValue].stakerValue,
            proofs[ticketValue].virtualStakerIndex
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
        uint256 expected = uint256(keccak256(abi.encodePacked(groupSelectionSeed, stakerValue, virtualStakerIndex)));
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
     * @param signingId Relay request ID associated with DKG protocol execution.
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
        uint256 signingId,
        uint256 submitterMemberIndex,
        bytes memory groupPubKey,
        bytes memory disqualified,
        bytes memory inactive,
        bytes memory signatures,
        uint[] memory signingMembersIndexes
    ) public onlyEligibleSubmitter(submitterMemberIndex) {

        require(
            disqualified.length == groupSize && inactive.length == groupSize,
            "Inactive and disqualified bytes arrays don't match the group size."
        );

        require(
            !dkgResultPublished[signingId], 
            "DKG result for this request ID already published."
        );

        bytes32 resultHash = keccak256(abi.encodePacked(groupPubKey, disqualified, inactive));
        verifySignatures(signatures, signingMembersIndexes, resultHash);
        address[] memory members = selectedParticipants();

        for (uint i = 0; i < groupSize; i++) {
            if(!_isInactive(inactive, i) && !_isDisqualified(disqualified, i)) {
                groupMembers[groupPubKey].push(members[i]);
            }
        }

        groups.push(Group(groupPubKey, block.number));
        // TODO: punish/reward logic
        cleanup();
        dkgResultPublished[signingId] = true;
        emit DkgResultPublishedEvent(signingId, groupPubKey);

        groupSelectionInProgress = false;
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
        require(signaturesCount >= groupThreshold, "Number of signatures is below honest majority threshold.");

        bytes memory current; // Current signature to be checked.
        uint256[] memory selected = selectedTickets();

        for(uint i = 0; i < signaturesCount; i++){
            require(signingMemberIndices[i] > 0, "Index should be greater than zero.");
            require(signingMemberIndices[i] <= selected.length, "Provided index is out of acceptable tickets bound.");
            current = signatures.slice(65*i, 65);
            address recoveredAddress = resultHash.toEthSignedMessageHash().recover(current);

            require(
                proofs[selected[signingMemberIndices[i] - 1]].sender == recoveredAddress,
                "Invalid signature. Signer and recovered address at provided index don't match."
            );
        }

        return true;
    }

    /**
     * @dev Checks if DKG protocol result has been already published for the
     * specific relay request ID associated with the protocol execution. 
     */
    function isDkgResultSubmitted(uint256 signingId) public view returns(bool) {
        return dkgResultPublished[signingId];
    }


    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() external payable {
        revert("Can not call contract without explicitly calling a function.");
    }

    /**
     * @dev Checks that the specified user has enough stake.
     * @param staker Specifies the identity of the staker.
     * @return True if staked enough to participate in the group, false otherwise.
     */
    function hasMinimumStake(address staker) public view returns(bool) {
        return StakingProxy(stakingProxy).balanceOf(staker) >= minimumStake;
    }

    /**
     * @dev Gets staking weight.
     * @param staker Specifies the identity of the staker.
     * @return Number of how many virtual stakers can staker represent.
     */
    function stakingWeight(address staker) public view returns(uint256) {
        return StakingProxy(stakingProxy).balanceOf(staker)/minimumStake;
    }

    /**
     * @dev Set the minimum amount of KEEP that allows a Keep network client to participate in a group.
     * @param _minimumStake Amount in KEEP.
     */
    function setMinimumStake(uint256 _minimumStake) public onlyOwner {
        minimumStake = _minimumStake;
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
        return groupSize.mul(space.div(tokenSupply().div(minimumStake)));
    }

    /**
     * @dev Gets number of active groups.
     */
    function numberOfGroups() public view returns(uint256) {
        return groups.length - expiredGroupOffset;
    }

    /**
     * @dev Gets the cutoff time in blocks until which the given group is
     * considered as an active group. The group may not be marked as expired
     * even though its active time has passed if one of the rules inside
     * `selectGroup` function are not met (e.g. minimum active group threshold).
     * Hence, this value informs when the group may no longer be considered
     * as active but it does not mean that the group will be immediatelly
     * considered not as such.
     */
    function groupActiveTimeOf(Group memory group) internal view returns(uint256) {
        return group.registrationBlockHeight + groupActiveTime;
    }

    /**
     * @dev Gets the cutoff time in blocks after which the given group is
     * considered as stale. Stale group is an expired group which is no longer
     * performing any operations.
     */
    function groupStaleTime(Group memory group) internal view returns(uint256) {
        return groupActiveTimeOf(group) + relayRequestTimeout;
    }

    /**
     * @dev Checks if a group with the given public key is a stale group.
     * Stale group is an expired group which is no longer performing any
     * operations. It is important to understand that an expired group may
     * still perform some operations for which it was selected when it was still
     * active. We consider a group to be stale when it's expired and when its
     * expiration time and potentially executed operation timeout are both in
     * the past.
     */
    function isStaleGroup(bytes memory groupPubKey) public view returns(bool) {
        for (uint i = 0; i < groups.length; i++) {
            if (groups[i].groupPubKey.equalStorage(groupPubKey)) {
                bool isExpired = expiredGroupOffset > i;
                bool isStale = groupStaleTime(groups[i]) < block.number;
                return isExpired && isStale;
            }
        }

        return true; // no group found, consider it as a stale group
    }

    /**
     * @dev Returns public key of a group from active groups using modulo operator.
     * @param seed Signing group selection seed.
     */
    function selectGroup(uint256 seed) public returns(bytes memory) {
        uint256 numberOfActiveGroups = groups.length - expiredGroupOffset;
        uint256 selectedGroup = seed % numberOfActiveGroups;

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
        if (numberOfActiveGroups > activeGroupsThreshold) {
            while (groupActiveTimeOf(groups[expiredGroupOffset + selectedGroup]) < block.number) {
                /**
                * We do -1 to see how many groups are available after the potential removal.
                * For example:
                * groups = [EEEAAAA]
                * - assuming selectedGroup = 0, then we'll have 4-0-1=3 groups after the removal: [EEEEAAA]
                * - assuming selectedGroup = 1, then we'll have 4-1-1=2 groups after the removal: [EEEEEAA]
                * - assuming selectedGroup = 2, then, we'll have 4-2-1=1 groups after the removal: [EEEEEEA]
                * - assuming selectedGroup = 3, then, we'll have 4-3-1=0 groups after the removal: [EEEEEEE]
                */
                if (numberOfActiveGroups - selectedGroup - 1 > activeGroupsThreshold) {
                    selectedGroup++;
                    expiredGroupOffset += selectedGroup;
                    numberOfActiveGroups -= selectedGroup;
                    selectedGroup = seed % numberOfActiveGroups;
                } else {
                    /* Number of groups that did not expire is less or equal activeGroupsThreshold
                    * and we have more groups than activeGroupsThreshold (including those expired) groups.
                    * Hence, we maintain the minimum activeGroupsThreshold of active groups and
                    * do not let any other groups to expire
                    */
                    expiredGroupOffset = groups.length - activeGroupsThreshold;
                    numberOfActiveGroups = activeGroupsThreshold;
                    selectedGroup = seed % numberOfActiveGroups;
                    break;
                }
            }
        }
        return groups[expiredGroupOffset + selectedGroup].groupPubKey;
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

        for (uint i = 0; i < tickets.length; i++) {
            delete proofs[tickets[i]];
        }

        delete tickets;

        // TODO: cleanup DkgResults
    }

    /**
     * @dev Creates a request to generate a new relay entry, which will include a
     * random number (by signing the previous entry's random number).
     * @param requestId Request Id trackable by service contract.
     * @param seed Initial seed random value from the client. It should be a cryptographically generated random value.
     * @param previousEntry Previous relay entry that is used to select a signing group for this request.
     */
    function sign(uint256 requestId, uint256 seed, uint256 previousEntry) public payable {

        require(
            serviceContracts.contains(msg.sender),
            "Only authorized service contract can request relay entry."
        );

        require(
            numberOfGroups() > 0,
            "At least one group needed to serve the request."
        );

        bytes memory groupPubKey = selectGroup(previousEntry);

        signingRequestCounter++;

        signingRequests[signingRequestCounter] = SigningRequest(requestId, msg.value, groupPubKey, msg.sender);

        emit SignatureRequested(signingRequestCounter, msg.value, previousEntry, seed, groupPubKey);
    }

    /**
     * @dev Creates a new relay entry and stores the associated data on the chain.
     * @param _signingId The request that started this generation - to tie the results back to the request.
     * @param _groupSignature The generated random number.
     * @param _groupPubKey Public key of the group that generated the threshold signature.
     */
    function relayEntry(uint256 _signingId, uint256 _groupSignature, bytes memory _groupPubKey, uint256 _previousEntry, uint256 _seed) public {

        require(signingRequests[_signingId].groupPubKey.equalStorage(_groupPubKey), "Provided group was not selected to produce entry for this request.");
        require(BLS.verify(_groupPubKey, abi.encodePacked(_previousEntry, _seed), bytes32(_groupSignature)), "Group signature failed to pass BLS verification.");

        address serviceContract = signingRequests[_signingId].serviceContract;
        uint256 requestId = signingRequests[_signingId].requestId;
        delete signingRequests[_signingId];

        emit SignatureSubmitted(_signingId, _groupSignature, _groupPubKey, _previousEntry, _seed);

        ServiceContract(serviceContract).entryCreated(requestId, _groupSignature);
        createGroup(_groupSignature, _signingId, _seed);
    }
}
