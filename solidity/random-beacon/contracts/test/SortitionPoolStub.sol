// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "../RandomBeacon.sol";

// Stub contract used in tests
contract SortitionPoolStub is ISortitionPool {
    mapping(address => bool) public operators;
    mapping(address => bool) public eligibleOperators;

    event OperatorsRemoved(address[] operators);

    function insertOperator(address operator) external override {
        operators[operator] = true;
    }

    function removeOperator(address operator) external override {
        delete operators[operator];
        delete eligibleOperators[operator];
    }

    function updateOperatorStatus(address operator) external override {
        // no-op
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
}
