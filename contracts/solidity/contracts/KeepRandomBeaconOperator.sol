pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/cryptography/ECDSA.sol";
import "./TokenStaking.sol";
import "./KeepRandomBeaconOperatorGroups.sol";
import "./utils/UintArrayUtils.sol";
import "./utils/AddressArrayUtils.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";
import "./cryptography/BLS.sol";

interface ServiceContract {
    function entryCreated(uint256 requestId, uint256 entry, address payable submitter) external;
    function fundRequestSubsidyFeePool() payable external;
}

/**
 * @title KeepRandomBeaconOperator
 * @dev Keep client facing contract for random beacon security-critical operations.
 * Handles group creation and expiration, BLS signature verification and incentives.
 * The contract is not upgradeable. New functionality can be implemented by deploying
 * new versions following Keep client update and re-authorization by the stakers.
 */
contract KeepRandomBeaconOperator {
    using SafeMath for uint256;
    using BytesLib for bytes;
    using ECDSA for bytes32;
    using AddressArrayUtils for address[];

    event OnGroupRegistered(bytes groupPubKey);

    // TODO: Rename to DkgResultSubmittedEvent
    // TODO: Add memberIndex
    event DkgResultPublishedEvent(bytes groupPubKey);

    // These are the public events that are used by clients
    event SignatureRequested(uint256 payment, uint256 previousEntry, uint256 seed, bytes groupPublicKey);
    event SignatureSubmitted(uint256 requestResponse, bytes requestGroupPubKey, uint256 previousEntry, uint256 seed);

    event GroupSelectionStarted(uint256 newEntry);

    // Contract owner.
    address public owner;

    address[] public serviceContracts;

    // TODO: replace with a secure authorization protocol (addressed in RFC 11).
    TokenStaking public stakingContract;

    KeepRandomBeaconOperatorGroups public groupContract;

    // Size of a group in the threshold relay.
    uint256 public groupSize = 5;

    // Minimum number of group members needed to interact according to the
    // protocol to produce a relay entry.
    uint256 public groupThreshold = 3;

    // Minimum amount of KEEP that allows sMPC cluster client to participate in
    // the Keep network.
    uint256 public minimumStake = 200000 * 1e18;

    // Timeout in blocks after the initial ticket submission is finished.
    uint256 public ticketInitialSubmissionTimeout = 4;

    // Timeout in blocks after the reactive ticket submission is finished.
    uint256 public ticketReactiveSubmissionTimeout = 4;

    // Timeout in blocks after the period where tickets can be challenged is
    // finished.
    uint256 public ticketChallengeTimeout = 4;

    // Time in blocks after which the next group member is eligible
    // to submit the result.
    uint256 public resultPublicationBlockStep = 3;

    // Time in blocks after DKG result is complete and ready to be published
    // by clients.
    uint256 public timeDKG = 7*(3+1);

    // Timeout in blocks for a relay entry to appear on the chain. Blocks are
    // counted from the moment relay request occur.
    //
    // Timeout is never shorter than the time needed by clients to generate
    // relay entry and the time it takes for the last group member to become
    // eligible to submit the result plus at least one block to submit it.
    uint256 public relayEntryTimeout = 24;

    uint256 public signingGasEstimate = 1240000; // TODO: Update once alt_bn128 gas costs reduction is implemented.
    uint256 public createGroupGasEstimate = 2260000;

    uint256 public dkgSubmitterReward;

    struct Proof {
        address sender;
        uint256 stakerValue;
        uint256 virtualStakerIndex;
    }

    mapping(uint256 => Proof) public proofs;

    bool public groupSelectionInProgress;

    uint256 public ticketSubmissionStartBlock;
    uint256 public groupSelectionRelayEntry;
    uint256[] public tickets;
    bytes[] public submissions;

    struct SigningRequest {
        uint256 relayRequestId;
        uint256 payment;
        uint256 signingFee;
        uint256 callbackFee;
        uint256 groupIndex;
        uint256 previousEntry;
        uint256 seed;
        address serviceContract;
    }

    uint256 public currentEntryStartBlock;
    SigningRequest internal signingRequest;

    bool internal entryInProgress;

    // Seed value used for the genesis group selection.
    // https://www.wolframalpha.com/input/?i=pi+to+78+digits
    uint256 internal _genesisGroupSeed = 31415926535897932384626433832795028841971693993751058209749445923078164062862;

    /**
     * @dev Triggers the first group selection. Genesis can be called only when
     * there are no groups on the operator contract.
     */
    function genesis() public {
        require(numberOfGroups() == 0, "There can be no groups.");
        startGroupSelection(_genesisGroupSeed, 0);
    }

    /**
     * @dev Throws if called by any account other than the owner.
     */
    modifier onlyOwner() {
        require(owner == msg.sender, "Caller is not the owner.");
        _;
    }

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
     * @dev Checks if sender is authorized.
     */
    modifier onlyServiceContract() {
        require(
            serviceContracts.contains(msg.sender),
            "Only authorized service contract can call this method."
        );
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
     * @dev Initializes the contract with service and staking contract addresses and
     * the deployer as the contract owner.
     */
    constructor(address _serviceContract, address _stakingContract, address _groupContract) public {
        require(_serviceContract != address(0), "Service contract address can't be zero.");
        require(_stakingContract != address(0), "Staking contract address can't be zero.");
        serviceContracts.push(_serviceContract);
        stakingContract = TokenStaking(_stakingContract);
        groupContract = KeepRandomBeaconOperatorGroups(_groupContract);
        owner = msg.sender;
    }

    /**
     * @dev Adds service contract
     * @param serviceContract Address of the service contract.
     */
    function addServiceContract(address serviceContract) public onlyOwner {
        serviceContracts.push(serviceContract);
    }

    /**
     * @dev Removes service contract
     * @param serviceContract Address of the service contract.
     */
    function removeServiceContract(address serviceContract) public onlyOwner {
        serviceContracts.removeAddress(serviceContract);
    }

    /**
     * @dev Triggers the selection process of a new candidate group.
     * @param _newEntry New random beacon value that stakers will use to
     * generate their tickets.
     */
    function createGroup(uint256 _newEntry) public payable onlyServiceContract {
        startGroupSelection(_newEntry, msg.value);
    }

    function startGroupSelection(uint256 _newEntry, uint256 _payment) internal {
        // dkgTimeout is the time after key generation protocol is expected to
        // be complete plus the expected time to submit the result.
        uint256 dkgTimeout = ticketSubmissionStartBlock +
            ticketChallengeTimeout +
            timeDKG +
            groupSize * resultPublicationBlockStep;

        require(!groupSelectionInProgress || block.number > dkgTimeout, "Group selection is in progress.");
        cleanup();
        ticketSubmissionStartBlock = block.number;
        groupSelectionRelayEntry = _newEntry;
        groupSelectionInProgress = true;
        emit GroupSelectionStarted(_newEntry);
        dkgSubmitterReward = _payment;
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
            stakingContract.authorizedTransferFrom(msg.sender, address(this), minimumStake);
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
        uint256 expected = uint256(keccak256(abi.encodePacked(groupSelectionRelayEntry, stakerValue, virtualStakerIndex)));
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

        bytes32 resultHash = keccak256(abi.encodePacked(groupPubKey, disqualified, inactive));
        verifySignatures(signatures, signingMembersIndexes, resultHash);
        address[] memory members = selectedParticipants();

        for (uint i = 0; i < groupSize; i++) {
            if(!_isInactive(inactive, i) && !_isDisqualified(disqualified, i)) {
                groupContract.addGroupMember(groupPubKey, members[i]);
            }
        }

        groupContract.addGroup(groupPubKey);

        // TODO: punish/reward logic
        uint256 rewards = dkgSubmitterReward;
        dkgSubmitterReward = 0;
        stakingContract.magpieOf(msg.sender).transfer(rewards);
        cleanup();
        emit DkgResultPublishedEvent(groupPubKey);

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
     * @param signingFee Amount in wei included in msg.value to cover the cost of signing to produce relay entry.
     * @param callbackFee Amount in wei included in msg.value to cover the cost of the callback contract call.
     */
    function sign(
        uint256 requestId,
        uint256 seed,
        uint256 previousEntry,
        uint256 signingFee,
        uint256 callbackFee
    ) public payable onlyServiceContract {
        signRelayEntry(requestId, seed, previousEntry, msg.sender, msg.value, signingFee, callbackFee);
    }

    function signRelayEntry(
        uint256 requestId,
        uint256 seed,
        uint256 previousEntry,
        address serviceContract,
        uint256 payment,
        uint256 signingFee,
        uint256 callbackFee
    ) internal {
        require(!entryInProgress || hasEntryTimedOut(), "Relay entry is in progress.");

        currentEntryStartBlock = block.number;
        entryInProgress = true;

        uint256 groupIndex = groupContract.selectGroup(previousEntry);
        bytes memory groupPubKey = groupContract.getGroupPublicKey(groupIndex);

        signingRequest = SigningRequest(
            requestId,
            payment,
            signingFee,
            callbackFee,
            groupIndex,
            previousEntry,
            seed,
            serviceContract
        );

        emit SignatureRequested(payment, previousEntry, seed, groupPubKey);
    }

    /**
     * @dev Creates a new relay entry and stores the associated data on the chain.
     * @param _groupSignature Group BLS signature over the concatenation of the
     * previous entry and seed.
     */
    function relayEntry(uint256 _groupSignature) public {
        bytes memory groupPubKey = groupContract.getGroupPublicKey(signingRequest.groupIndex);

        require(
            BLS.verify(
                groupPubKey,
                abi.encodePacked(signingRequest.previousEntry, signingRequest.seed),
                bytes32(_groupSignature)
            ),
            "Group signature failed to pass BLS verification."
        );
        
        emit SignatureSubmitted(
            _groupSignature,
            groupPubKey,
            signingRequest.previousEntry,
            signingRequest.seed
        );

        ServiceContract(signingRequest.serviceContract).entryCreated(
            signingRequest.relayRequestId,
            _groupSignature,
            msg.sender
        );

        entryInProgress = false;

        // Calculate each group member reward = baseReward * delayFactor / groupSize
        uint256 decimals = 1e4; // Adding 4 decimals to perform float division.
        uint256 profitMargin = signingRequest.payment.sub(signingRequest.signingFee).sub(signingRequest.callbackFee);
        uint256 baseReward = profitMargin.div(groupSize);
        uint256 entryTimeout = currentEntryStartBlock.add(relayEntryTimeout);
        uint256 delayFactor = entryTimeout.sub(block.number).mul(decimals).div(relayEntryTimeout.sub(1))**2;
        uint256 delayFactorInverse = uint256(1).mul(decimals**2).sub(delayFactor);
        uint256 groupReward = baseReward.mul(delayFactor).div(decimals**2);
        uint256 delayPenalty = baseReward.mul(delayFactorInverse).div(decimals**2);
        
        for (uint i = 0; i < groupSize; i++) {
            address payable operator = address(uint160(groupContract.getGroupMember(groupPubKey, i)));
            stakingContract.magpieOf(operator).transfer(groupReward);
        }

        // The submitter reward consists of:
        // The callback gas expenditure (reimbursed by the service contract)
        // The entry verification fee to cover the cost of verifying the submission
        // Extra reward - 5% of the delay penalties of the entire group
        uint256 extraReward = delayPenalty.mul(groupSize).mul(5).div(100); 
        stakingContract.magpieOf(msg.sender).transfer(signingRequest.signingFee.add(extraReward));

        // Rewards not paid out to the operators are paid out to requesters to subsidize new requests.
        uint256 subsidy = profitMargin.sub(groupReward.mul(groupSize)).sub(extraReward);
        if (subsidy > 0) {
            ServiceContract(signingRequest.serviceContract).fundRequestSubsidyFeePool.value(subsidy)();
        }
    }

    /**
     * @dev Returns true if the currently ongoing new relay entry generation
     * operation timed out. There is a certain timeout for a new relay entry
     * to be produced, see `relayEntryTimeout` value.
     */
    function hasEntryTimedOut() internal view returns (bool) {
        return entryInProgress && block.number > currentEntryStartBlock + relayEntryTimeout;
    }

    /**
     * @dev Function used to inform about the fact the currently ongoing
     * new relay entry generation operation timed out. As a result, the group
     * which was supposed to produce a new relay entry is immediatelly
     * terminated and a new group is selected to produce a new relay entry.
     */
    function reportRelayEntryTimeout() public {
        require(hasEntryTimedOut(), "Current relay entry did not time out");

        groupContract.terminateGroup(signingRequest.groupIndex);

        // We could terminate the last active group. If that's the case,
        // do not try to execute signing again because there is no group
        // which can handle it.
        if (numberOfGroups() > 0) {
            signRelayEntry(
                signingRequest.relayRequestId,
                signingRequest.seed,
                signingRequest.previousEntry,
                signingRequest.serviceContract,
                signingRequest.payment,
                signingRequest.signingFee,
                signingRequest.callbackFee
            );
        }
    }

    /**
     * @dev Checks that the specified user has enough stake.
     * @param staker Specifies the identity of the staker.
     * @return True if staked enough to participate in the group, false otherwise.
     */
    function hasMinimumStake(address staker) public view returns(bool) {
        return stakingContract.balanceOf(staker) >= minimumStake;
    }

    /**
     * @dev Gets staking weight.
     * @param staker Specifies the identity of the staker.
     * @return Number of how many virtual stakers can staker represent.
     */
    function stakingWeight(address staker) public view returns(uint256) {
        return stakingContract.balanceOf(staker).div(minimumStake);
    }

    /**
     * @dev Checks if group with the given public key is registered.
     */
    function isGroupRegistered(bytes memory groupPubKey) public view returns(bool) {
        return groupContract.isGroupRegistered(groupPubKey);
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
        return groupContract.isStaleGroup(groupPubKey);
    }

    /**
     * @dev Gets the number of active groups. Expired and terminated groups are
     * not counted as active.
     */
    function numberOfGroups() public view returns(uint256) {
        return groupContract.numberOfGroups();
    }
}
