pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";

library KeepRandomBeaconOperatorPricing {

    using SafeMath for uint256;

    struct Options {
        // Time in blocks it takes off-chain cluster to generate a new relay entry
        // and be ready to submit it to the chain.
        uint256 relayEntryGenerationTime;

        // Timeout in blocks for a relay entry to appear on the chain.
        uint256 relayEntryTimeout;
    }

    function init(
        Options storage self,
        uint256 _relayEntryGenerationTime,
        uint256 _relayEntryTimeout
    ) public {
        self.relayEntryGenerationTime = _relayEntryGenerationTime;
        self.relayEntryTimeout = _relayEntryTimeout;
    }

    /**
     * @dev Gets delay factor for rewards calculation.
     * @return Integer representing floating-point number with 16 decimals places.
     */
    function getDelayFactor(
        Options storage options,
        uint256 currentEntryStartBlock
    ) internal view returns(uint256 delayFactor) {
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