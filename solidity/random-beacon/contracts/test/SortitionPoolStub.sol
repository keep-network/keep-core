// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "@keep-network/sortition-pools/contracts/SortitionTree.sol";
import "../RandomBeacon.sol";

// Stub contract used in tests
contract SortitionPoolStub is ISortitionPool, SortitionTree {
    mapping(address => bool) public operators;
    mapping(address => bool) public eligibleOperators;

    constructor() SortitionTree() {
        // Fill operators IDs array so identifiers start with bigger number.
        uint256 i = 100;
        while (i > 0) {
            allocateOperatorID(address(0));
            i--;
        }
    }

    function joinPool(address operator) external override {
        operators[operator] = true;

        allocateOperatorID(operator);
    }

    function isOperatorInPool(address operator)
        external
        view
        override
        returns (bool)
    {
        return operators[operator];
    }

    // Helper function, it does not exist in the sortition pool
    function setOperatorEligibility(address operator, bool eligibility) public {
        eligibleOperators[operator] = eligibility;
    }

    function isOperatorEligible(address operator)
        public
        view
        override
        returns (bool)
    {
        return eligibleOperators[operator];
    }

    // TODO: Fix sortition pool public API to accept/return uint32 for IDs
    function getIDOperator(uint32 id) public view override returns (address) {
        return SortitionTree.getIDOperator(id);
    }
}
