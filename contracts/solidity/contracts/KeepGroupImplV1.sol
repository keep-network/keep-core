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

    uint256 public groupThreshold;
    uint256 public groupSize;
    uint256 public minimumStake;
    address public stakingProxy;
    address public randomBeacon;

    uint256 public ticketInitialSubmissionTimeout;
    uint256 public ticketReactiveSubmissionTimeout;
    uint256 public ticketChallengeTimeout;
    uint256 public timeDKG;
    uint256 public resultPublicationBlockStep;
    uint256 public ticketSubmissionStartBlock;
    uint256 public randomBeaconValue;

    uint256[] public tickets;
    bytes[] public submissions;

    // Store whether DKG result was published for the corresponding requestID.
    mapping (uint256 => bool) public dkgResultPublished;

    struct Proof {
        address sender;
        uint256 stakerValue;
        uint256 virtualStakerIndex;
    }

    mapping(uint256 => Proof) public proofs;

    struct Group {
        bytes groupPubKey;
        uint registrationTime;
    }

    Group[] public groups;

    mapping (bytes => address[]) private groupMembers;

    bool public initialized;


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
     * @dev Triggers the selection process of a new candidate group.
     */
    function runGroupSelection(uint256 _randomBeaconValue) public {
        require(msg.sender == randomBeacon);
        cleanup();
        ticketSubmissionStartBlock = block.number;
        randomBeaconValue = _randomBeaconValue;
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
    function getTicketProof(uint256 ticketValue) public view returns (address, uint256, uint256) {
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
        uint256 expected = uint256(keccak256(abi.encodePacked(randomBeaconValue, stakerValue, virtualStakerIndex)));
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
            disqualified.length == groupSize && inactive.length == groupSize,
            "Inactive and disqualified bytes arrays don't match the group size."
        );

        require(
            !dkgResultPublished[requestId], 
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
        dkgResultPublished[requestId] = true;
        emit DkgResultPublishedEvent(requestId, groupPubKey);
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
    function isDkgResultSubmitted(uint256 requestId) public view returns(bool) {
        return dkgResultPublished[requestId];
    }


    /**
     * @dev Prevent receiving ether without explicitly calling a function.
     */
    function() external payable {
        revert("Can not call contract without explicitly calling a function.");
    }

    /**
     * @dev Initialize Keep Group implementation contract with a linked Staking proxy contract.
     * @param _stakingProxy Address of a staking proxy contract that will be linked to this contract.
     * @param _randomBeacon Address of a random beacon contract that will be linked to this contract.
     * @param _minimumStake Minimum amount in KEEP that allows KEEP network client to participate in a group.
     * @param _groupSize Size of a group in the threshold relay.
     * @param _groupThreshold Minimum number of interacting group members needed to produce a relay entry.
     * @param _ticketInitialSubmissionTimeout Timeout in blocks after the initial ticket submission is finished.
     * @param _ticketReactiveSubmissionTimeout Timeout in blocks after the reactive ticket submission is finished.
     * @param _ticketChallengeTimeout Timeout in blocks after the period where tickets can be challenged is finished.
     * @param _timeDKG Timeout in blocks after DKG result is complete and ready to be published.
     * @param _resultPublicationBlockStep Time in blocks after which member with the given index is eligible
     * to submit DKG result.
     */
    function initialize(
        address _stakingProxy,
        address _randomBeacon,
        uint256 _minimumStake,
        uint256 _groupThreshold,
        uint256 _groupSize,
        uint256 _ticketInitialSubmissionTimeout,
        uint256 _ticketReactiveSubmissionTimeout,
        uint256 _ticketChallengeTimeout,
        uint256 _timeDKG,
        uint256 _resultPublicationBlockStep
    ) public onlyOwner {
        require(!initialized, "Contract is already initialized.");
        require(_stakingProxy != address(0x0), "Staking proxy address can't be zero.");
        initialized = true;
        stakingProxy = _stakingProxy;
        randomBeacon = _randomBeacon;
        minimumStake = _minimumStake;
        groupSize = _groupSize;
        groupThreshold = _groupThreshold;
        ticketInitialSubmissionTimeout = _ticketInitialSubmissionTimeout;
        ticketReactiveSubmissionTimeout = _ticketReactiveSubmissionTimeout;
        ticketChallengeTimeout = _ticketChallengeTimeout;
        timeDKG = _timeDKG;
        resultPublicationBlockStep = _resultPublicationBlockStep;
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
     * @dev Allows owner to change the groupSize.
     */
    function setGroupSize(uint256 _groupSize) public onlyOwner {
        groupSize = _groupSize;
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
        return groups.length;
    }

    /**
     * @dev Returns public key of a group from available groups using modulo operator.
     * @param previousEntry Previous random beacon value.
     */
    function selectGroup(uint256 previousEntry) public view returns(bytes memory) {
        return groups[previousEntry % groups.length].groupPubKey;
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

}
