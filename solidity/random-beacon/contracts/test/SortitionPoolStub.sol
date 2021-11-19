// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "./SortitionTreeStub.sol";
import "../RandomBeacon.sol";

// Stub contract used in tests
//
// TODO: Deprecated. This stub should be eventually removed in favor of real
//       sortition pool.
contract SortitionPoolStub is ISortitionPool {
    SortitionTreeStub internal sortitionTree;
    bool internal locked;

    mapping(address => bool) public operators;
    uint256 public operatorsCount;

    mapping(address => bool) public eligibleOperators;

    mapping(bytes32 => uint32[]) public groups;

    event OperatorStatusUpdated(uint32 id);

    constructor() {
        sortitionTree = new SortitionTreeStub();
    }

    function lock() external override {
        locked = true;
    }

    function unlock() external override {
        locked = false;
    }

    function insertOperator(address operator) external override {
        require(!locked, "Pool is locked");

        operators[operator] = true;
        operatorsCount++;

        sortitionTree.publicAllocateOperatorID(operator);
    }

    function banRewards(uint32[] calldata ids, uint256 duration)
        external
        override
    {
        // no-op
    }

    function updateOperatorStatus(uint32 id) external override {
        require(!locked, "Pool is locked");

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

    function isLocked() public view override returns (bool) {
        return locked;
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
}
