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

    function getSlashingFactor() external view returns (uint256) {
        return relay.getSlashingFactor();
    }
}
