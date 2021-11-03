// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "../RandomBeacon.sol";

// Stub contract used in tests
contract StakingStub is IStaking {
    mapping(address => uint256) public stakedTokens;

    event Slashed(uint256 amount, address[] operators);

    function slash(uint256 amount, address[] memory operators)
        external
        override
    {
        if (amount > 0 && operators.length > 0) {
            emit Slashed(amount, operators);
        }
    }

    function eligibleStake(
        address operator,
        address // operatorContract
    ) external view returns (uint256) {
        return stakedTokens[operator];
    }

    function setStake(address operator, uint256 stake) public {
        stakedTokens[operator] = stake;
    }
}
