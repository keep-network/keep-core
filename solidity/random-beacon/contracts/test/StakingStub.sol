// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "../RandomBeacon.sol";

// Stub contract used in tests
contract StakingStub is IStaking {
    event Slashed(uint256 amount, address[] operators);

    function slash(uint256 amount, address[] memory operators)
        external
        override
    {
        if (amount > 0 && operators.length > 0) {
            emit Slashed(amount, operators);
        }
    }
}
