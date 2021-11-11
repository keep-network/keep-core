// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "@keep-network/sortition-pools/contracts/SortitionTree.sol";
import "../RandomBeacon.sol";

// Stub contract used in tests
contract SortitionPoolStub is ISortitionPool, SortitionTree {
    mapping(address => bool) public operators;
    mapping(address => bool) public eligibleOperators;

    mapping(bytes32 => uint32[]) public groups;

    event OperatorsRemoved(address[] operators);

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

    function getIDOperators(uint32[] calldata ids)
        public
        view
        override
        returns (address[] memory)
    {
        address[] memory operators = new address[](ids.length);

        for (uint256 i = 0; i < ids.length; i++) {
            operators[i] = SortitionTree.getIDOperator(ids[i]);
        }

        return operators;
    }

    function setSelectGroupResult(bytes32 seed, uint32[] calldata members)
        public
    {
        groups[seed] = members;
    }

    function selectGroup(uint256 groupSize, bytes32 seed)
        external
        view
        override
        returns (uint32[] memory members)
    {
        members = groups[seed];
        require(groupSize == members.length, "Wrong group size");
        return members;
    }

    function removeOperators(uint32[] calldata ids) external override {
        address[] memory _operators = getIDOperators(ids);

        for (uint256 i = 0; i < _operators.length; i++) {
            delete operators[_operators[i]];
            delete eligibleOperators[_operators[i]];
        }

        if (_operators.length > 0) {
            emit OperatorsRemoved(_operators);
        }
    }
}
