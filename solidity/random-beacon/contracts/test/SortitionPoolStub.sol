// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "../RandomBeacon.sol";

// Stub contract used in tests
contract SortitionPoolStub is ISortitionPool {
    mapping(address => bool) public operators;
    uint256 public operatorsCount;

    mapping(address => bool) public eligibleOperators;

    event OperatorStatusUpdated(address operator);

    function insertOperator(address operator) external override {
        operators[operator] = true;
        operatorsCount++;
    }

    function removeOperator(address operator) external override {
        delete operators[operator];
        delete eligibleOperators[operator];
        operatorsCount--;
    }

    function updateOperatorStatus(address operator) external override {
        emit OperatorStatusUpdated(operator);
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
