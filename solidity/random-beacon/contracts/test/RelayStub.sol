// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../libraries/Relay.sol";
import "../libraries/Groups.sol";

contract RelayStub {
    using Relay for Relay.Data;

    Relay.Data internal relay;

    bytes public groupPubKey;

    constructor() {
        relay.setRelayEntrySubmissionEligibilityDelay(10);
        relay.setRelayEntryHardTimeout(100);
        relay.setRelayEntrySubmissionFailureSlashingAmount(1000e18);
    }

    function setCurrentRequestStartBlock() external {
        relay.currentRequestStartBlock = uint64(block.number);
    }

    function calculateSlashingAmount() external returns (uint256) {
        return relay.calculateSlashingAmount();
    }

    function submitEntry(bytes calldata entry, bytes calldata _groupPubKey)
        external
    {
        groupPubKey = _groupPubKey;
        relay.submitEntry(entry, groupPubKey);
    }
}
