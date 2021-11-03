// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "@keep-network/sortition-pools/contracts/SortitionTree.sol";

// Stub contract used in tests
contract SortitionTreeStub is SortitionTree {
    function publicAllocateOperatorID(address operator)
        public
        returns (uint256)
    {
        return allocateOperatorID(operator);
    }
}
