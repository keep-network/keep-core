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
        relay.currentRequest.startBlock = uint128(block.number);
    }

    function isEligible(
        uint256 submitterIndex,
        bytes calldata entry,
        uint256 groupSize
    ) external view returns (bool) {
        (uint256 firstEligibleIndex, uint256 lastEligibleIndex) = relay
            .getEligibilityRange(entry, groupSize);

        return
            relay.isEligible(
                submitterIndex,
                firstEligibleIndex,
                lastEligibleIndex
            );
    }

    function getInactiveMembers(
        uint256 submitterIndex,
        uint256 firstEligibleIndex,
        uint32[] memory groupMembers
    ) external view returns (uint32[] memory) {
        return
            relay.getInactiveMembers(
                submitterIndex,
                firstEligibleIndex,
                groupMembers
            );
    }

    function getSlashingFactor(uint256 groupSize)
        external
        view
        returns (uint256)
    {
        return relay.getSlashingFactor(groupSize);
    }
}
