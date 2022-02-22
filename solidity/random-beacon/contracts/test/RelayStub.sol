// SPDX-License-Identifier: MIT
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
        relay.setRelayEntrySoftTimeout(relayEntrySoftTimeout);
        relay.setRelayEntryHardTimeout(relayEntryHardTimeout);
    }

    function setCurrentRequestStartBlock() external {
        relay.currentRequestStartBlock = uint64(block.number);
    }

    function setRelayEntrySubmissionFailureSlashingAmount(
        uint256 relayEntrySubmissionFailureSlashingAmount
    ) external {
        relay.setRelayEntrySubmissionFailureSlashingAmount(
            relayEntrySubmissionFailureSlashingAmount
        );
    }

    function calculateSlashingAmount() external returns (uint256) {
        return relay.calculateSlashingAmount();
    }
}
