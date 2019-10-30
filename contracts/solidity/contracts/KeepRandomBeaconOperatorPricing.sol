pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";

library KeepRandomBeaconOperatorPricing {

    using SafeMath for uint256;

    struct Options {
        // Size of a group in the threshold relay.
        uint256 groupSize;

        // Time in blocks it takes off-chain cluster to generate a new relay entry
        // and be ready to submit it to the chain.
        uint256 relayEntryGenerationTime;

        // Timeout in blocks for a relay entry to appear on the chain.
        uint256 relayEntryTimeout;

        // Each signing group member reward expressed in wei.
        uint256 groupMemberBaseReward;
    }

    function init(
        Options storage self,
        uint256 _groupSize,
        uint256 _relayEntryGenerationTime,
        uint256 _relayEntryTimeout,
        uint256 _groupMemberBaseReward
    ) public {
        self.groupSize = _groupSize;
        self.relayEntryGenerationTime = _relayEntryGenerationTime;
        self.relayEntryTimeout = _relayEntryTimeout;
        self.groupMemberBaseReward = _groupMemberBaseReward;
    }

    /**
     * @dev Get rewards breakdown in wei for successful entry for the current signing request.
     */
    function newEntryRewardsBreakdown(
        Options storage options,
        uint256 currentEntryStartBlock,
        uint256 entryVerificationAndProfitFee
    ) public view returns(uint256 groupMemberReward, uint256 submitterReward, uint256 subsidy) {
        uint256 decimals = 1e16; // Adding 16 decimals to perform float division.

        uint256 groupProfitFee = options.groupMemberBaseReward.mul(options.groupSize);

        uint256 delayFactor = getDelayFactor(options, currentEntryStartBlock);
        groupMemberReward = options.groupMemberBaseReward.mul(delayFactor).div(decimals);

        // delay penalty = base reward * (1 - delay factor)
        uint256 groupMemberDelayPenalty = options.groupMemberBaseReward.sub(
            options.groupMemberBaseReward.mul(delayFactor).div(decimals)
        );

        // The submitter reward consists of:
        // The callback gas expenditure (reimbursed by the service contract)
        // The entry verification fee to cover the cost of verifying the submission,
        // paid regardless of their gas expenditure
        // Submitter extra reward - 5% of the delay penalties of the entire group
        uint256 submitterExtraReward = groupMemberDelayPenalty.mul(options.groupSize).mul(5).div(100);
        uint256 entryVerificationFee = entryVerificationAndProfitFee.sub(groupProfitFee);
        submitterReward = entryVerificationFee.add(submitterExtraReward);

        // Rewards not paid out to the operators are paid out to requesters to subsidize new requests.
        subsidy = groupProfitFee.sub(groupMemberReward.mul(options.groupSize)).sub(submitterExtraReward);
    }

    /**
     * @dev Gets delay factor for rewards calculation.
     * @return Integer representing floating-point number with 16 decimals places.
     */
    function getDelayFactor(
        Options storage options,
        uint256 currentEntryStartBlock
    ) public view returns(uint256 delayFactor) {
        uint256 decimals = 1e16; // Adding 16 decimals to perform float division.

        // T_deadline is the earliest block when no submissions are accepted
        // and an entry timed out. The last block the entry can be published in is
        //     currentEntryStartBlock + relayEntryTimeout
        // and submission are no longer accepted from block
        //     currentEntryStartBlock + relayEntryTimeout + 1.
        uint256 deadlineBlock = currentEntryStartBlock.add(options.relayEntryTimeout).add(1);

        // T_begin is the earliest block the result can be published in.
        // It takes relayEntryGenerationTime to generate a new entry, so it can
        // be published at block relayEntryGenerationTime + 1 the earliest.
        uint256 submissionStartBlock = currentEntryStartBlock.add(options.relayEntryGenerationTime).add(1);

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
}