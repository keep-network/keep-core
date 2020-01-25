pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/utils/ReentrancyGuard.sol";
import "./TokenStaking.sol";
import "./cryptography/BLS.sol";
import "./utils/AddressArrayUtils.sol";
import "./libraries/operator/GroupSelection.sol";
import "./libraries/operator/Groups.sol";
import "./libraries/operator/DKGResultVerification.sol";

interface ServiceContract {
    function entryCreated(uint256 requestId, bytes calldata entry, address payable submitter) external;
    function fundRequestSubsidyFeePool() external payable;
    function fundDkgFeePool() external payable;
}

/**
 * @title KeepRandomBeaconOperator
 * @dev Keep client facing contract for random beacon security-critical operations.
 * Handles group creation and expiration, BLS signature verification and incentives.
 * The contract is not upgradeable. New functionality can be implemented by deploying
 * new versions following Keep client update and re-authorization by the stakers.
 */
contract KeepRandomBeaconOperator is ReentrancyGuard {
    using SafeMath for uint256;
    using AddressArrayUtils for address[];
    using GroupSelection for GroupSelection.Storage;
    using Groups for Groups.Storage;
    using DKGResultVerification for DKGResultVerification.Storage;

    event OnGroupRegistered(bytes groupPubKey);

    // TODO: Rename to DkgResultSubmittedEvent
    // TODO: Add memberIndex
    event DkgResultPublishedEvent(bytes groupPubKey);

    event RelayEntryRequested(bytes previousEntry, bytes groupPublicKey);
    event RelayEntrySubmitted();

    event GroupSelectionStarted(uint256 newEntry);

    event GroupMemberRewardsWithdrawn(address indexed beneficiary, address operator, uint256 amount, uint256 groupIndex);

    GroupSelection.Storage groupSelection;
    Groups.Storage groups;
    DKGResultVerification.Storage dkgResultVerification;

    // Contract owner.
    address internal owner;

    address[] internal serviceContracts;

    // TODO: replace with a secure authorization protocol (addressed in RFC 11).
    TokenStaking internal stakingContract;

    // Minimum amount of KEEP that allows sMPC cluster client to participate in
    // the Keep network. Expressed as number with 18-decimal places.
    uint256 public minimumStake = 200000 * 1e18;

    // Each signing group member reward expressed in wei.
    uint256 public groupMemberBaseReward = 145*1e11; // 14500 Gwei, 10% of operational cost

    // The price feed estimate is used to calculate the gas price for reimbursement
    // next to the actual gas price from the transaction. We use both values to
    // defend against malicious miner-submitters who can manipulate transaction
    // gas price. Expressed in wei.
    uint256 public priceFeedEstimate = 20*1e9; // (20 Gwei = 20 * 10^9 wei)

    // Fluctuation margin to cover the immediate rise in gas price.
    // Expressed in percentage.
    uint256 public fluctuationMargin = 50; // 50%

    // Size of a group in the threshold relay.
    uint256 public groupSize = 64;

    // Minimum number of group members needed to interact according to the
    // protocol to produce a relay entry.
    uint256 public groupThreshold = 33;

    // Time in blocks after which the next group member is eligible
    // to submit the result.
    uint256 public resultPublicationBlockStep = 3;

    // Time in blocks it takes off-chain cluster to generate a new relay entry
    // and be ready to submit it to the chain.
    uint256 public relayEntryGenerationTime = (1+3);

    // Timeout in blocks for a relay entry to appear on the chain. Blocks are
    // counted from the moment relay request occur.
    //
    // Timeout is never shorter than the time needed by clients to generate
    // relay entry and the time it takes for the last group member to become
    // eligible to submit the result plus at least one block to submit it.
    uint256 public relayEntryTimeout = relayEntryGenerationTime.add(groupSize.mul(resultPublicationBlockStep));

    // Gas required to verify BLS signature and produce successful relay
    // entry. Excludes callback and DKG gas.
    uint256 public entryVerificationGasEstimate = 300000;

    // Gas required to submit DKG result.
    uint256 public dkgGasEstimate = 1740000;

    // Reimbursement for the submitter of the DKG result.
    // This value is set when a new DKG request comes to the operator contract.
    // It contains a full payment for DKG multiplied by the fluctuation margin.
    // When submitting DKG result, the submitter is reimbursed with the actual cost
    // and some part of the fee stored in this field may be returned to the service
    // contract.
    uint256 public dkgSubmitterReimbursementFee;

    // Service contract that triggered current group selection.
    ServiceContract internal groupSelectionStarterContract;

    struct SigningRequest {
        uint256 relayRequestId;
        uint256 entryVerificationAndProfitFee;
        uint256 groupIndex;
        bytes previousEntry;
        address serviceContract;
    }

    uint256 internal currentEntryStartBlock;
    SigningRequest internal signingRequest;

    // Seed value used for the genesis group selection.
    // https://www.wolframalpha.com/input/?i=pi+to+78+digits
    uint256 internal _genesisGroupSeed = 31415926535897932384626433832795028841971693993751058209749445923078164062862;

    /**
     * @dev Triggers the first group selection. Genesis can be called only when
     * there are no groups on the operator contract.
     */
    function genesis() public payable {
        require(numberOfGroups() == 0, "Groups exist");
        // Set latest added service contract as a group selection starter to receive any DKG fee surplus.
        groupSelectionStarterContract = ServiceContract(serviceContracts[serviceContracts.length.sub(1)]);
        startGroupSelection(_genesisGroupSeed, msg.value);
    }

    /**
     * @dev Throws if called by any account other than the owner.
     */
    modifier onlyOwner() {
        require(owner == msg.sender, "Caller is not the owner");
        _;
    }

    /**
     * @dev Checks if sender is authorized.
     */
    modifier onlyServiceContract() {
        require(
            serviceContracts.contains(msg.sender),
            "Caller is not an authorized contract"
        );
        _;
    }

    constructor(address _serviceContract, address _stakingContract) public {
        owner = msg.sender;

        serviceContracts.push(_serviceContract);
        stakingContract = TokenStaking(_stakingContract);

        groupSelection.ticketSubmissionTimeout = 12;
        groupSelection.groupSize = groupSize;
        groups.activeGroupsThreshold = 5;
        groups.groupActiveTime = 3000;

        dkgResultVerification.timeDKG = 6*(1+5);
        dkgResultVerification.resultPublicationBlockStep = resultPublicationBlockStep;
        dkgResultVerification.groupSize = groupSize;
        // TODO: For now, the required number of signatures is equal to group
        // threshold. This should be updated to keep a safety margin for
        // participants misbehaving during signing.
        dkgResultVerification.signatureThreshold = groupThreshold;
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
     * @dev Set the gas price in wei for calculating reimbursements.
     * @param _priceFeedEstimate is the gas price for calculating reimbursements.
     */
    function setPriceFeedEstimate(uint256 _priceFeedEstimate) public onlyOwner {
        priceFeedEstimate = _priceFeedEstimate;
    }

    /**
     * @dev Adds a safety margin for gas price fluctuations to the current gas price.
     * The gas price for DKG or relay entry is set when the request is processed
     * but the result submission transaction will be sent later. We add a safety
     * margin that should be sufficient for getting requests processed within a
     * a deadline under all circumstances.
     * @param gasPrice Gas price in wei.
     */
    function gasPriceWithFluctuationMargin(uint256 gasPrice) internal view returns (uint256) {
        return gasPrice.add(gasPrice.mul(fluctuationMargin).div(100));
    }

    /**
     * @dev Triggers the selection process of a new candidate group.
     * @param _newEntry New random beacon value that stakers will use to
     * generate their tickets.
     */
    function createGroup(uint256 _newEntry) public payable onlyServiceContract {
        groupSelectionStarterContract = ServiceContract(msg.sender);
        startGroupSelection(_newEntry, msg.value);
    }

    function startGroupSelection(uint256 _newEntry, uint256 _payment) internal {
        require(
            _payment >= gasPriceWithFluctuationMargin(priceFeedEstimate).mul(dkgGasEstimate),
            "Insufficient DKG fee"
        );

        require(isGroupSelectionPossible(), "Group selection in progress");

        // If previous group selection failed and there is reimbursement left
        // return it to the DKG fee pool.
        if (dkgSubmitterReimbursementFee > 0) {
            uint256 surplus = dkgSubmitterReimbursementFee;
            dkgSubmitterReimbursementFee = 0;
            ServiceContract(msg.sender).fundDkgFeePool.value(surplus)();
        }

        groupSelection.start(_newEntry);
        emit GroupSelectionStarted(_newEntry);
        dkgSubmitterReimbursementFee = _payment;
    }

    function isGroupSelectionPossible() public view returns (bool) {
        if (!groupSelection.inProgress) {
            return true;
        }

        // dkgTimeout is the time after key generation protocol is expected to
        // be complete plus the expected time to submit the result.
        uint256 dkgTimeout = groupSelection.ticketSubmissionStartBlock +
        groupSelection.ticketSubmissionTimeout +
        dkgResultVerification.timeDKG +
        groupSize * resultPublicationBlockStep;

        return block.number > dkgTimeout;
    }

    /**
     * @dev Submits ticket to request to participate in a new candidate group.
     * @param ticket Bytes representation of a ticket that holds the following:
     * - ticketValue: first 8 bytes of a result of keccak256 cryptography hash
     *   function on the combination of the group selection seed (previous
     *   beacon output), staker-specific value (address) and virtual staker index.
     * - stakerValue: a staker-specific value which is the address of the staker.
     * - virtualStakerIndex: 4-bytes number within a range of 1 to staker's weight;
     *   has to be unique for all tickets submitted by the given staker for the
     *   current candidate group selection.
     */
    function submitTicket(bytes32 ticket) public {
        uint64 ticketValue;
        uint160 stakerValue;
        uint32 virtualStakerIndex;

        bytes memory ticketBytes = abi.encodePacked(ticket);
        /* solium-disable-next-line */
        assembly {
            // ticket value is 8 bytes long
            ticketValue := mload(add(ticketBytes, 8))
            // staker value is 20 bytes long
            stakerValue := mload(add(ticketBytes, 28))
            // virtual staker index is 4 bytes long
            virtualStakerIndex := mload(add(ticketBytes, 32))
        }

        uint256 stakingWeight = stakingContract.balanceOf(msg.sender).div(minimumStake);
        groupSelection.submitTicket(ticketValue, uint256(stakerValue), uint256(virtualStakerIndex), stakingWeight);
    }

    /**
     * @dev Gets the timeout in blocks after which group candidate ticket
     * submission is finished.
     */
    function ticketSubmissionTimeout() public view returns (uint256) {
        return groupSelection.ticketSubmissionTimeout;
    }

    /**
     * @dev Gets the number of submitted group candidate tickets so far.
     */
    function submittedTicketsCount() public view returns (uint256) {
        return groupSelection.tickets.length;
    }

    /**
     * @dev Gets selected participants in ascending order of their tickets.
     */
    function selectedParticipants() public view returns (address[] memory) {
        return groupSelection.selectedParticipants();
    }

    /**
     * @dev Submits result of DKG protocol. It is on-chain part of phase 14 of
     * the protocol.
     *
     * @param submitterMemberIndex Claimed submitter candidate group member index
     * @param groupPubKey Generated candidate group public key
     * @param disqualified Bytes representing disqualified group members;
     * 1 at the specific index means that the member has been disqualified.
     * Indexes reflect positions of members in the group, as outputted by the
     * group selection protocol.
     * @param inactive Bytes representing inactive group members;
     * 1 at the specific index means that the member has been marked as inactive.
     * Indexes reflect positions of members in the group, as outputted by the
     * group selection protocol.
     * @param signatures Concatenation of signatures from members supporting the
     * result.
     * @param signingMembersIndexes Indices of members corresponding to each
     * signature.
     */
    function submitDkgResult(
        uint256 submitterMemberIndex,
        bytes memory groupPubKey,
        bytes memory disqualified,
        bytes memory inactive,
        bytes memory signatures,
        uint[] memory signingMembersIndexes
    ) public {
        address[] memory members = selectedParticipants();

        dkgResultVerification.verify(
            submitterMemberIndex,
            groupPubKey,
            disqualified,
            inactive,
            signatures,
            signingMembersIndexes,
            members,
            groupSelection.ticketSubmissionStartBlock + groupSelection.ticketSubmissionTimeout
        );

        for (uint i = 0; i < groupSize; i++) {
            // Check member was neither marked as inactive nor as disqualified
            if(inactive[i] == 0x00 && disqualified[i] == 0x00) {
                groups.addGroupMember(groupPubKey, members[i]);
            }
        }

        groups.addGroup(groupPubKey);
        reimburseDkgSubmitter();
        emit DkgResultPublishedEvent(groupPubKey);
        groupSelection.stop();
    }

    /**
     * @dev Compare the reimbursement fee calculated based on the current transaction gas
     * price and the current price feed estimate with the DKG reimbursement fee calculated
     * and paid at the moment when the DKG was requested. If there is any surplus, it will
     * be returned to the DKG fee pool of the service contract which triggered the DKG.
     */
    function reimburseDkgSubmitter() internal {
        uint256 gasPrice = priceFeedEstimate;
        // We need to check if tx.gasprice is non-zero as a workaround to a bug
        // in go-ethereum:
        // https://github.com/ethereum/go-ethereum/pull/20189
        if (tx.gasprice > 0 && tx.gasprice < priceFeedEstimate) {
            gasPrice = tx.gasprice;
        }

        uint256 reimbursementFee = dkgGasEstimate.mul(gasPrice);
        address payable magpie = stakingContract.magpieOf(msg.sender);

        if (reimbursementFee < dkgSubmitterReimbursementFee) {
            uint256 surplus = dkgSubmitterReimbursementFee.sub(reimbursementFee);
            dkgSubmitterReimbursementFee = 0;
            // Reimburse submitter with actual DKG cost.
            (bool success, ) = magpie.call.value(reimbursementFee)("");
            require(success, "Failed reimburse actual DKG cost");

            // Return surplus to the contract that started DKG.
            groupSelectionStarterContract.fundDkgFeePool.value(surplus)();
        } else {
            // If submitter used higher gas price reimburse only dkgSubmitterReimbursementFee max.
            reimbursementFee = dkgSubmitterReimbursementFee;
            dkgSubmitterReimbursementFee = 0;
            (bool success, ) = magpie.call.value(reimbursementFee)("");
            require(success, "Failed reimburse DKG fee");
        }
    }

    /**
     * @dev Set the minimum amount of KEEP that allows a Keep network client to participate in a group.
     * @param _minimumStake Amount in KEEP.
     */
    function setMinimumStake(uint256 _minimumStake) public onlyOwner {
        minimumStake = _minimumStake;
    }

    /**
     * @dev Creates a request to generate a new relay entry, which will include a
     * random number (by signing the previous entry's random number).
     * @param requestId Request Id trackable by service contract
     * @param previousEntry Previous relay entry
     */
    function sign(
        uint256 requestId,
        bytes memory previousEntry
    ) public payable onlyServiceContract {
        require(
            msg.value >= groupProfitFee().add(entryVerificationGasEstimate.mul(gasPriceWithFluctuationMargin(priceFeedEstimate))),
            "Insufficient new entry fee"
        );
        signRelayEntry(requestId, previousEntry, msg.sender, msg.value);
    }

    function signRelayEntry(
        uint256 requestId,
        bytes memory previousEntry,
        address serviceContract,
        uint256 entryVerificationAndProfitFee
    ) internal {
        require(!isEntryInProgress() || hasEntryTimedOut(), "Beacon is busy");

        currentEntryStartBlock = block.number;

        uint256 groupIndex = groups.selectGroup(uint256(keccak256(previousEntry)));
        signingRequest = SigningRequest(
            requestId,
            entryVerificationAndProfitFee,
            groupIndex,
            previousEntry,
            serviceContract
        );

        bytes memory groupPubKey = groups.getGroupPublicKey(groupIndex);
        emit RelayEntryRequested(previousEntry, groupPubKey);
    }

    /**
     * @dev Creates a new relay entry and stores the associated data on the chain.
     * @param _groupSignature Group BLS signature over the concatenation of the
     * previous entry and seed.
     */
    function relayEntry(bytes memory _groupSignature) public nonReentrant {
        require(isEntryInProgress(), "Entry was submitted");
        require(!hasEntryTimedOut(), "Entry timed out");

        bytes memory groupPubKey = groups.getGroupPublicKey(signingRequest.groupIndex);

        require(
            BLS.verify(
                groupPubKey,
                signingRequest.previousEntry,
                _groupSignature
            ),
            "Invalid signature"
        );

        emit RelayEntrySubmitted();

        ServiceContract(signingRequest.serviceContract).entryCreated(
            signingRequest.relayRequestId,
            _groupSignature,
            msg.sender
        );

        (uint256 groupMemberReward, uint256 submitterReward, uint256 subsidy) = newEntryRewardsBreakdown();
        groups.addGroupMemberReward(groupPubKey, groupMemberReward);

        (bool success, ) = stakingContract.magpieOf(msg.sender).call.value(submitterReward)("");
        require(success, "Failed send relay submitter reward");

        if (subsidy > 0) {
            ServiceContract(signingRequest.serviceContract).fundRequestSubsidyFeePool.value(subsidy)();
        }

        currentEntryStartBlock = 0;
    }

    /**
     * @dev Get rewards breakdown in wei for successful entry for the current signing request.
     */
    function newEntryRewardsBreakdown() internal view returns(uint256 groupMemberReward, uint256 submitterReward, uint256 subsidy) {
        uint256 decimals = 1e16; // Adding 16 decimals to perform float division.

        uint256 delayFactor = getDelayFactor();
        groupMemberReward = groupMemberBaseReward.mul(delayFactor).div(decimals);

        // delay penalty = base reward * (1 - delay factor)
        uint256 groupMemberDelayPenalty = groupMemberBaseReward.mul(decimals.sub(delayFactor));

        // The submitter reward consists of:
        // The callback gas expenditure (reimbursed by the service contract)
        // The entry verification fee to cover the cost of verifying the submission,
        // paid regardless of their gas expenditure
        // Submitter extra reward - 5% of the delay penalties of the entire group
        uint256 submitterExtraReward = groupMemberDelayPenalty.mul(groupSize).mul(5).div(100).div(decimals);
        uint256 entryVerificationFee = signingRequest.entryVerificationAndProfitFee.sub(groupProfitFee());
        submitterReward = entryVerificationFee.add(submitterExtraReward);

        // Rewards not paid out to the operators are paid out to requesters to subsidize new requests.
        subsidy = groupProfitFee().sub(groupMemberReward.mul(groupSize)).sub(submitterExtraReward);
    }

    /**
     * @dev Gets delay factor for rewards calculation.
     * @return Integer representing floating-point number with 16 decimals places.
     */
    function getDelayFactor() internal view returns(uint256 delayFactor) {
        uint256 decimals = 1e16; // Adding 16 decimals to perform float division.

        // T_deadline is the earliest block when no submissions are accepted
        // and an entry timed out. The last block the entry can be published in is
        //     currentEntryStartBlock + relayEntryTimeout
        // and submission are no longer accepted from block
        //     currentEntryStartBlock + relayEntryTimeout + 1.
        uint256 deadlineBlock = currentEntryStartBlock.add(relayEntryTimeout).add(1);

        // T_begin is the earliest block the result can be published in.
        // It takes relayEntryGenerationTime to generate a new entry, so it can
        // be published at block relayEntryGenerationTime + 1 the earliest.
        uint256 submissionStartBlock = currentEntryStartBlock.add(relayEntryGenerationTime).add(1);

        // Use submissionStartBlock block as entryReceivedBlock if entry submitted earlier than expected.
        uint256 entryReceivedBlock = block.number <= submissionStartBlock ? submissionStartBlock:block.number;

        // T_remaining = T_deadline - T_received
        uint256 remainingBlocks = deadlineBlock.sub(entryReceivedBlock);

        // T_deadline - T_begin
        uint256 submissionWindow = deadlineBlock.sub(submissionStartBlock);

        // delay factor = [ T_remaining / (T_deadline - T_begin)]^2
        //
        // Since we add 16 decimal places to perform float division, we do:
        // delay factor = [ T_temaining * decimals / (T_deadline - T_begin)]^2 / decimals =
        //    = [T_remaining / (T_deadline - T_begin) ]^2 * decimals
        delayFactor = ((remainingBlocks.mul(decimals).div(submissionWindow))**2).div(decimals);
    }

    /**
     * @dev Returns true if generation of a new relay entry is currently in
     * progress.
     */
    function isEntryInProgress() internal view returns (bool) {
        return currentEntryStartBlock != 0;
    }

    /**
     * @dev Returns true if the currently ongoing new relay entry generation
     * operation timed out. There is a certain timeout for a new relay entry
     * to be produced, see `relayEntryTimeout` value.
     */
    function hasEntryTimedOut() internal view returns (bool) {
        return currentEntryStartBlock != 0 && block.number > currentEntryStartBlock + relayEntryTimeout;
    }

    /**
     * @dev Function used to inform about the fact the currently ongoing
     * new relay entry generation operation timed out. As a result, the group
     * which was supposed to produce a new relay entry is immediatelly
     * terminated and a new group is selected to produce a new relay entry.
     */
    function reportRelayEntryTimeout() public {
        require(hasEntryTimedOut(), "Entry did not time out");

        groups.terminateGroup(signingRequest.groupIndex);

        // We could terminate the last active group. If that's the case,
        // do not try to execute signing again because there is no group
        // which can handle it.
        if (numberOfGroups() > 0) {
            signRelayEntry(
                signingRequest.relayRequestId,
                signingRequest.previousEntry,
                signingRequest.serviceContract,
                signingRequest.entryVerificationAndProfitFee
            );
        }
    }

    /**
     * @dev Gets group profit fee expressed in wei.
     */
    function groupProfitFee() public view returns(uint256) {
        return groupMemberBaseReward.mul(groupSize);
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
     * @dev Checks if group with the given public key is registered.
     */
    function isGroupRegistered(bytes memory groupPubKey) public view returns(bool) {
        return groups.isGroupRegistered(groupPubKey);
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
        return groups.isStaleGroup(groupPubKey);
    }

    /**
     * @dev Gets the number of active groups. Expired and terminated groups are
     * not counted as active.
     */
    function numberOfGroups() public view returns(uint256) {
        return groups.numberOfGroups();
    }

    /**
     * @dev Returns accumulated group member rewards for provided group.
     */
    function getGroupMemberRewards(bytes memory groupPubKey) public view returns (uint256) {
        return groups.groupMemberRewards[groupPubKey];
    }

    /**
     * @dev Gets all indices in the provided group for a member.
     */
    function getGroupMemberIndices(bytes memory groupPubKey, address member) public view returns (uint256[] memory indices) {
        return groups.getGroupMemberIndices(groupPubKey, member);
    }

    /**
     * @dev Withdraws accumulated group member rewards for operator
     * using the provided group index and member indices. Once the
     * accumulated reward is withdrawn from the selected group, member is
     * removed from it. Rewards can be withdrawn only from stale group.
     * @param operator Operator address.
     * @param groupIndex Group index.
     * @param groupMemberIndices Array of member indices for the group member.
     */
    function withdrawGroupMemberRewards(address operator, uint256 groupIndex, uint256[] memory groupMemberIndices) public nonReentrant {
        uint256 accumulatedRewards = groups.withdrawFromGroup(operator, groupIndex, groupMemberIndices);
        (bool success, ) = stakingContract.magpieOf(operator).call.value(accumulatedRewards)("");
        require(success, "Failed withdraw rewards");
        emit GroupMemberRewardsWithdrawn(stakingContract.magpieOf(operator), operator, accumulatedRewards, groupIndex);
    }

    /**
    * @dev Gets the index of the first active group.
    */
    function getFirstActiveGroupIndex() public view returns (uint256) {
        return groups.expiredGroupOffset;
    }

    /**
    * @dev Gets group public key.
    */
    function getGroupPublicKey(uint256 groupIndex) public view returns (bytes memory) {
        return groups.getGroupPublicKey(groupIndex);
    }

    /**
     * @dev Returns members of the given group by group public key.
     */
    function getGroupMembers(bytes memory groupPubKey) public view returns (address[] memory members) {
        return groups.getGroupMembers(groupPubKey);
    }
}
