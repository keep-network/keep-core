// SPDX-License-Identifier: MIT
//
// ▓▓▌ ▓▓ ▐▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▄
// ▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
//   ▓▓▓▓▓▓    ▓▓▓▓▓▓▓▀    ▐▓▓▓▓▓▓    ▐▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
//   ▓▓▓▓▓▓▄▄▓▓▓▓▓▓▓▀      ▐▓▓▓▓▓▓▄▄▄▄         ▓▓▓▓▓▓▄▄▄▄         ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
//   ▓▓▓▓▓▓▓▓▓▓▓▓▓▀        ▐▓▓▓▓▓▓▓▓▓▓         ▓▓▓▓▓▓▓▓▓▓         ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
//   ▓▓▓▓▓▓▀▀▓▓▓▓▓▓▄       ▐▓▓▓▓▓▓▀▀▀▀         ▓▓▓▓▓▓▀▀▀▀         ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▀
//   ▓▓▓▓▓▓   ▀▓▓▓▓▓▓▄     ▐▓▓▓▓▓▓     ▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌
// ▓▓▓▓▓▓▓▓▓▓ █▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓
// ▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓
//
//

pragma solidity ^0.8.9;

import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "@threshold-network/solidity-contracts/contracts/staking/TokenStaking.sol";

library Authorization {
    struct AuthorizationDecrease {
        uint96 decreasingTo; // amount
        uint64 decreasingAt; // timestamp
    }

    struct Parameters {
        uint96 minimumAuthorization;
        uint64 authorizationDecreaseDelay;
    }

    struct Data {
        Parameters parameters;
        mapping(address => AuthorizationDecrease) authorizationDecreaseRequests;
    }

    function setMinimumAuthorization(
        Data storage self,
        uint96 _minimumAuthorization
    ) internal {
        self.parameters.minimumAuthorization = _minimumAuthorization;
    }

    function setAuthorizationDecreaseDelay(
        Data storage self,
        uint64 _authorizationDecreaseDelay
    ) internal {
        self
            .parameters
            .authorizationDecreaseDelay = _authorizationDecreaseDelay;
    }

    /// @notice Used by T staking contract to inform the random beacon that the
    ///         authorized amount for the given operator increased. Can only be
    ///         called when the sortition pool is not locked. Increases in-pool
    ///         weight and rewards weight in the pool proportionally to the
    ///         authorized stake amount immediatelly. Reverts if the sortition
    ///         pool is locked or if the authorization amount is below the
    ///         minimum.
    /// @dev Should only be callable by T staking contract.
    function authorizationIncreased(
        Data storage self,
        SortitionPool sortitionPool,
        address operator,
        uint96 toAmount
    ) external {
        require(!sortitionPool.isLocked(), "Sortition pool is locked");
        require(
            toAmount >= self.parameters.minimumAuthorization,
            "Below minimum authorization amount"
        );

        if (sortitionPool.isOperatorInPool(operator)) {
            sortitionPool.updateOperatorStatus(operator, uint256(toAmount));
        } else {
            sortitionPool.insertOperator(operator, uint256(toAmount));
        }
    }

   
}
