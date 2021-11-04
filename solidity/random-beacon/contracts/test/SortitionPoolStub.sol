// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "./SortitionTreeStub.sol";
import "../RandomBeacon.sol";

// Stub contract used in tests
contract SortitionPoolStub is ISortitionPool {
    SortitionTreeStub internal sortitionTree;

    mapping(address => bool) public operators;
    uint256 public operatorsCount;

    mapping(address => bool) public eligibleOperators;

    event OperatorStatusUpdated(uint32 id);

    constructor() {
        sortitionTree = new SortitionTreeStub();
    }

    function insertOperator(address operator) external override {
        operators[operator] = true;
        operatorsCount++;

        sortitionTree.publicAllocateOperatorID(operator);
    }

    function removeOperators(uint32[] calldata ids) external override {
        address[] memory _operators = getIDOperators(ids);

        for (uint256 i = 0; i < _operators.length; i++) {
            delete operators[_operators[i]];
            delete eligibleOperators[_operators[i]];
            operatorsCount--;
        }
    }

    function updateOperatorStatus(uint32 id) external override {
        emit OperatorStatusUpdated(id);
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

    function getIDOperator(uint32 id) public view override returns (address) {
        return sortitionTree.getIDOperator(id);
    }

    function getIDOperators(uint32[] calldata ids)
        public
        view
        override
        returns (address[] memory)
    {
        address[] memory operators = new address[](ids.length);

        for (uint256 i = 0; i < ids.length; i++) {
            operators[i] = sortitionTree.getIDOperator(ids[i]);
        }

        return operators;
    }

    function getOperatorID(address operator)
        public
        view
        override
        returns (uint32)
    {
        return sortitionTree.getOperatorID(operator);
    }

    function transferOwnership(address newOwner) public {
        // no-op
    }
}
