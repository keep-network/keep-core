// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../libraries/Relay.sol";
import "../libraries/Groups.sol";

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

    function getPunishedMembers(
        uint256 submitterIndex,
        uint256 firstEligibleIndex,
        address[] memory members
    ) external view returns (address[] memory) {
        Groups.Group memory group;
        group.members = members;

        return
            relay.getPunishedMembers(
                submitterIndex,
                firstEligibleIndex,
                group,
                groupSize
            );
    }
}
