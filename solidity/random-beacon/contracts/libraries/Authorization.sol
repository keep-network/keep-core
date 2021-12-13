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

    /// @notice Used by T staking contract to inform the random beacon that the
    ///         authorization decrease for the given operator has been regitered. 
    ///         Can only be called when the sortition pool is not locked.
    ///
    ///         If the authorization is decreased to the amount higher than the
    ///         minimum authorization, in-pool weight and rewarded stake are
    ///         reduced in the sortition pool immediately so that the operator,
    ///         during the autorization decrease delay period, can not be 
    ///         selected to new groups with stake authorization amount higher
    ///         than the one to which it is deauthorizing to. Authorized stake
    ///         amount remains the same until deauthorization request is
    ///         approved.
    ///
    ///         If the operator is requesting authorization decrease to zero,
    ///         in pool-weight is reduced to zero immediately and operator is
    ///         removed from the sortition pool. Rewarded stake stays at the
    ///         same level as before. Operator, during the authorization
    ///         decrease delay period, can not be selected to new groups but
    ///         it should still earn rewards based on their last authrized stake
    ///         until authorization decrease request is approved.
    ///
    ///         Overwrites pending authorization decrease request if the one
    ///         pending is for the amount higher than the minimum authorization.
    ///
    ///         Reverts if authorization decrease to 0 is pending.
    ///         Reverts if authorization decrease is requested to non-zero value
    ///         below the minimum authorization amount.
    ///         Reverts if the sortition pool is locked.
    function authorizationDecreaseRequested(
        Data storage self,
        SortitionPool sortitionPool,
        address operator,
        uint96 fromAmount,
        uint96 toAmount
    ) external {
        require(
            toAmount == 0 || toAmount >= self.parameters.minimumAuthorization,
            "Authorization amount should be 0 or above the minimum"
        );

        if (toAmount == 0) {
            // Update in-pool weight to 0 but keep rewards the same until
            // authorization decrease is not approved. We do not want that
            // operator to be selected to new groups.
            sortitionPool.updateOperatorStatusAndRewards(
                operator,
                0,
                fromAmount
            );
        } else {
            // Update in-pool weight and rewards weight to the new authorization
            // amount. The operator should be selected to new groups with the
            // new weight.
            sortitionPool.updateOperatorStatus(operator, toAmount);
        }

        // Register new authorization decrease request but do not decrease
        // authorization immediatelly. There might be still some groups that
        // might need slashing the higher amount.
        self.authorizationDecreaseRequests[operator] = AuthorizationDecrease(
            toAmount,
            // solhint-disable-next-line not-rely-on-time
            uint64(block.timestamp) + self.parameters.authorizationDecreaseDelay
        );
    }

    /// @notice Approves the previously registered authorization decrease
    ///         request. Reverts if authorization decrease delay have not passed
    ///         yet or if the auhorization decrease was not requested for the
    ///         given operator.
    function approveAuthorizationDecrease(
        Data storage self,
        TokenStaking tokenStaking,
        address operator
    ) external {
        AuthorizationDecrease memory request = self
            .authorizationDecreaseRequests[operator];
        require(
            request.decreasingAt > 0,
            "Authorization decrease not requested"
        );
        require(
            // solhint-disable-next-line not-rely-on-time
            request.decreasingAt > block.timestamp,
            "Authorization decrease delay not passsed"
        );

        tokenStaking.approveAuthorizationDecrease(operator);
        delete self.authorizationDecreaseRequests[operator];
    }

    // TODO: involuntaryAuthorizationDecrease

    // TODO: withdrawal of sortition pool rewards in RandomBeacon and the code
    //       allowing to calculate reward reduction based on the time passed
    //       between authorization decrase delay passed and the moment
    //       withdrawRewards is called.
}
