// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.6;

import "../libraries/Relay.sol";
import "../libraries/Groups.sol";

contract RelayStub {
    using Relay for Relay.Data;

    Relay.Data internal relay;

    function setTimeouts(
        uint256 relayEntrySoftTimeout,
        uint256 relayEntryHardTimeout
    ) public {
        relay.setTimeouts(relayEntrySoftTimeout, relayEntryHardTimeout);
    }

    function setCurrentRequestStartBlock() external {
        relay.currentRequestStartBlock = uint64(block.number);
    }

    function setRelayEntrySubmissionFailureSlashingAmount(
        uint96 relayEntrySubmissionFailureSlashingAmount
    ) external {
        relay.setRelayEntrySubmissionFailureSlashingAmount(
            relayEntrySubmissionFailureSlashingAmount
        );
    }

    function calculateSlashingAmount() external view returns (uint96) {
        return relay.calculateSlashingAmount();
    }
}
