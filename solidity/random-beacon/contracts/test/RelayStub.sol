// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../libraries/Relay.sol";
import "../libraries/Groups.sol";

contract RelayStub {
    using Relay for Relay.Data;

    Relay.Data internal relay;

    constructor() {
        relay.setRelayEntrySubmissionEligibilityDelay(10);
        relay.setRelayEntryHardTimeout(5760);
    }

    function setCurrentRequestStartBlock() external {
        relay.currentRequestStartBlock = uint64(block.number);
    }

    function getSlashingFactor(uint256 groupSize)
        external
        view
        returns (uint256)
    {
        return relay.getSlashingFactor(groupSize);
    }
}
