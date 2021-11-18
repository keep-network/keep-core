// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "../RandomBeacon.sol";

// Stub contract used in tests
contract StakingStub is IRandomBeaconStaking {
    mapping(address => uint256) public stakedTokens;

    bool public shouldFail;

    event Slashed(uint256 amount, address[] operators);

    event Seized(
        uint256 amount,
        uint256 rewardMultiplier,
        address notifier,
        address[] operators
    );

    function slash(uint256 amount, address[] memory operators)
        external
        override
    {
        if (shouldFail) {
            revert("error");
        }

        if (amount > 0 && operators.length > 0) {
            emit Slashed(amount, operators);
        }
    }

    function seize(
        uint256 amount,
        uint256 rewardMultiplier,
        address notifier,
        address[] memory operators
    ) external override {
        if (shouldFail) {
            revert("error");
        }

        if (amount > 0 && operators.length > 0) {
            emit Seized(amount, rewardMultiplier, notifier, operators);
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

    function setFailureFlag(bool _shouldFail) external {
        shouldFail = _shouldFail;
    }
}
