// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../libraries/Relay.sol";

contract RelayStub {
    using Relay for Relay.Data;

    uint256 public constant groupSize = 8;

    Relay.Data internal relay;

    constructor() {
        relay.setRelayEntrySubmissionEligibilityDelay(10);
    }

    function setCurrentRequestStartBlock() external {
        relay.currentRequest.startBlock = uint128(block.number);
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
                lastEligibleIndex
            );
    }
}
