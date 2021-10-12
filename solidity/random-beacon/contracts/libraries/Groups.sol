// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

library Groups {
    struct Group {
        bytes groupPubKey;
        uint256 activationTimestamp;
        address[] members;
    }

    function selectGroup(Data storage self, uint256 seed) internal view returns (Groups.Group memory) {
        // TODO: Assert at least one group exists and implement selection logic.
        Groups.Group memory group;
        return group;
    }
}
