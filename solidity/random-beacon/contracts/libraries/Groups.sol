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
        group.groupPubKey = "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d01659dc18b57722ecf6a4beb4d04dfe780a660c4c3bb2b165ab8486114c464c621bf37ecdba226629c20908c7f475c5b3a7628ce26d696436eab0b0148034dfcd";
        return group;
    }
}
