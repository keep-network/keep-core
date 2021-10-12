// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

library Groups {
    struct Group {
        bytes groupPubKey;
        uint256 activationTimestamp;
        address[] members;
    }
}
