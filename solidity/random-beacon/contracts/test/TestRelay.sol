// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../libraries/Relay.sol";

contract TestRelay {
    using Relay for Relay.Data;

    uint256 public constant groupSize = 8;

    Relay.Data internal relay;

    constructor() {
        relay.setRelayEntrySubmissionEligibilityDelay(10);
    }

    function isEligible(uint256 submitterIndex, bytes calldata entry)
        external
        view
        returns (bool)
    {
        (uint256 firstEligibleIndex, uint256 lastEligibleIndex) = relay
            .getEligibilityRange(entry, groupSize);

        return
            relay.isEligible(
                submitterIndex,
                firstEligibleIndex,
                lastEligibleIndex,
                groupSize
            );
    }
}
