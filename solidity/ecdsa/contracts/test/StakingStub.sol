// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import {IWalletStaking} from "../WalletRegistry.sol";

// TODO: get rid of this contract; use T staking contract for tests
// Stub contract used in tests
contract StakingStub is IWalletStaking {
    mapping(address => uint256) public stakedTokens;

    event Seized(
        uint256 amount,
        uint256 rewardMultiplier,
        address notifier,
        address[] stakingProviders
    );

    function authorizedStake(
        address stakingProvider,
        address // application
    ) external view override returns (uint256) {
        return stakedTokens[stakingProvider];
    }

    function increaseAuthorization(
        address stakingProvider,
        address, // application
        uint96 amount
    ) public {
        stakedTokens[stakingProvider] += amount;
    }

    function requestAuthorizationDecrease(address stakingProvider) external {
        stakedTokens[stakingProvider] = 0;
    }

    function seize(
        uint96 amount,
        uint256 rewardMultiplier,
        address notifier,
        address[] memory stakingProviders
    ) external {
        if (amount > 0 && stakingProviders.length > 0) {
            for (uint256 i = 0; i < stakingProviders.length; i++) {
                stakedTokens[stakingProviders[i]] -= amount;
            }
            emit Seized(amount, rewardMultiplier, notifier, stakingProviders);
        }
    }
}
