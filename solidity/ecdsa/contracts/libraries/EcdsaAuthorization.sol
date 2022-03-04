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
import "@threshold-network/solidity-contracts/contracts/staking/IStaking.sol";

/// @notice Library managing the state of stake authorizations for ECDSA
///         operator contract and the presence of operators in the sortition
///         pool based on the stake authorized for them.
library EcdsaAuthorization {
    struct Parameters {
        // The minimum authorization required by ECDSA application so that
        // operator can join the sortition pool and do the work.
        uint96 minimumAuthorization;
        // Authorization decrease delay in seconds between the time
        // authorization decrease is requested and the time the authorization
        // decrease can be approved. It is always the same value, no matter if
        // authorization decrease amount is small, significant, or if it is
        // a decrease to zero.
        uint64 authorizationDecreaseDelay;
    }

    struct AuthorizationDecrease {
        uint96 decreasingBy; // amount
        uint64 decreasingAt; // timestamp
    }

    struct Data {
        Parameters parameters;
        mapping(address => address) stakingProviderToOperator;
        mapping(address => address) operatorToStakingProvider;
        mapping(address => AuthorizationDecrease) pendingDecreases;
    }

    event OperatorRegistered(
        address indexed stakingProvider,
        address indexed operator
    );

    event AuthorizationIncreaseRequested(
        address indexed stakingProvider,
        address indexed operator,
        uint96 toAmount
    );

    event AuthorizationDecreaseRequested(
        address indexed stakingProvider,
        address indexed operator,
        uint96 decreasingBy,
        uint64 decreasingAt
    );

    event AuthorizationDecreaseApproved(address indexed stakingProvider);

    /// @notice Sets the minimum authorization for ECDSA application. Without
    ///         at least the minimum authorization, staking provider is not
    ///         eligible to join and operate in the network.
    function setMinimumAuthorization(
        Data storage self,
        uint96 _minimumAuthorization
    ) internal {
        self.parameters.minimumAuthorization = _minimumAuthorization;
    }

    /// @notice Sets the authorization decrease delay. It is the time in seconds
    ///         that needs to pass between the time authorization decrease is
    ///         requested and the time the authorization decrease can be
    ///         approved, no matter the authorization decrease amount.
    function setAuthorizationDecreaseDelay(
        Data storage self,
        uint64 _authorizationDecreaseDelay
    ) internal {
        self
            .parameters
            .authorizationDecreaseDelay = _authorizationDecreaseDelay;
    }

    /// @notice Used by staking provider to set operator address that will
    ///         operate ECDSA node. The given staking provider can set operator
    ///         address only one time. The operator address can not be changed
    ///         and must be unique. Reverts if the operator is already set for
    ///         the staking provider or if the operator address is already in
    ///         use. Reverts if there is a pending authorization decrease for
    ///         the staking provider.
    function registerOperator(Data storage self, address operator) internal {
        address stakingProvider = msg.sender;

        require(
            self.stakingProviderToOperator[stakingProvider] == address(0),
            "Operator already set for the staking provider"
        );
        require(
            self.operatorToStakingProvider[operator] == address(0),
            "Operator address already in use"
        );

        // Authorization request for a staking provider who has not yet
        // registered their operator can be approved immediately.
        // We need to make sure that the approval happens before operator
        // is registered to do not let the operator join the sortition pool
        // with an unresolved authorization decrease request that can be
        // approved at any point.
        AuthorizationDecrease storage decrease = self.pendingDecreases[
            stakingProvider
        ];
        require(
            decrease.decreasingAt == 0,
            "There is a pending authorization decrease request"
        );

        self.stakingProviderToOperator[stakingProvider] = operator;
        self.operatorToStakingProvider[operator] = stakingProvider;

        emit OperatorRegistered(stakingProvider, operator);
    }

    /// @notice Used by T staking contract to inform the application that the
    ///         authorized stake amount for the given staking provider increased.
    ///
    ///         Reverts if the authorization amount is below the minimum.
    ///
    ///         The function is not updating the sortition pool. Sortition pool
    ///         state needs to be updated by the operator with a call to
    ///         `joinSortitionPool` or `updateOperatorStatus`.
    ///
    /// @dev Should only be callable by T staking contract.
    function authorizationIncreased(
        Data storage self,
        address stakingProvider,
        uint96 toAmount
    ) internal {
        require(
            toAmount >= self.parameters.minimumAuthorization,
            "Authorization below the minimum"
        );

        address operator = self.stakingProviderToOperator[stakingProvider];
        emit AuthorizationIncreaseRequested(
            stakingProvider,
            operator,
            toAmount
        );
    }

    /// @notice Used by T staking contract to inform the application that the
    ///         authorization decrease for the given staking provider has been
    ///         requested.
    ///
    ///         Reverts if the amount after deauthorization would be non-zero
    ///         and lower than the minimum authorization.
    ///
    ///         If the operator is not known (`registerOperator` was not called)
    ///         it lets to `approveAuthorizationDecrease` immediatelly. If the
    ///         operator is known (`registerOperator` was called), the operator
    ///         needs to update state of the sortition pool with a call to
    ///         `joinSortitionPool` or `updateOperatorStatus`. After the
    ///         sortition pool state is in sync, authorization decrease delay
    ///         starts.
    ///
    ///         After authorization decrease delay passes, authorization
    ///         decrease request needs to be approved with a call to
    ///         `approveAuthorizationDecrease` function.
    ///
    ///         If there is a pending authorization decrease request, it is
    ///         overwritten.
    ///
    /// @dev Should only be callable by T staking contract.
    function authorizationDecreaseRequested(
        Data storage self,
        address stakingProvider,
        uint96 fromAmount,
        uint96 toAmount
    ) internal {
        require(
            toAmount == 0 || toAmount >= self.parameters.minimumAuthorization,
            "Authorization amount should be 0 or above the minimum"
        );

        address operator = self.stakingProviderToOperator[stakingProvider];

        uint64 decreasingAt;

        if (operator == address(0)) {
            // Operator is not known. It means `registerOperator` was not
            // called yet, and there is no chance the operator could
            // call `joinSortitionPool`. We can let to approve authorization
            // decrease immediately because that operator was never in the
            // sortition pool.

            // solhint-disable-next-line not-rely-on-time
            decreasingAt = uint64(block.timestamp);
        } else {
            // Operator is known. It means that this operator is or was in
            // the sortition pool. Before authorization decrease delay starts,
            // the operator needs to update the state of the sortition pool
            // with a call to `joinSortitionPool` or `updateOperatorStatus`.
            // For now, we set `decreasingAt` as "never decreasing" and let
            // it be updated by `joinSortitionPool` or `updateOperatorStatus`
            // once we know the sortition pool is in sync.

            // solhint-disable-next-line not-rely-on-time
            decreasingAt = type(uint64).max;
        }

        uint96 decreasingBy = fromAmount - toAmount;

        self.pendingDecreases[stakingProvider] = AuthorizationDecrease(
            decreasingBy,
            decreasingAt
        );

        emit AuthorizationDecreaseRequested(
            stakingProvider,
            operator,
            decreasingBy,
            decreasingAt
        );
    }

    /// @notice Approves the previously registered authorization decrease
    ///         request. Reverts if authorization decrease delay have not passed
    ///         yet or if the auhorization decrease was not requested for the
    ///         given staking provider.
    function approveAuthorizationDecrease(
        Data storage self,
        IStaking tokenStaking,
        address stakingProvider
    ) internal {
        AuthorizationDecrease storage decrease = self.pendingDecreases[
            stakingProvider
        ];
        require(
            decrease.decreasingAt > 0,
            "Authorization decrease not requested"
        );
        require(
            decrease.decreasingAt != type(uint64).max,
            "Authorization decrease request not activated"
        );
        require(
            // solhint-disable-next-line not-rely-on-time
            block.timestamp > decrease.decreasingAt,
            "Authorization decrease delay not passsed"
        );

        tokenStaking.approveAuthorizationDecrease(stakingProvider);
        delete self.pendingDecreases[stakingProvider];

        emit AuthorizationDecreaseApproved(stakingProvider);
    }

    /// @notice Lets the operator join the sortition pool. The operator address
    ///         must be known - before calling this function, it has to be
    ///         appointed by the staking provider by calling `registerOperator`.
    ///         Also, the operator must have the minimum authorization required
    ///         by ECDSA. Function reverts if there is no minimum stake
    ///         authorized or if the operator is not known. If there was an
    ///         authorization decrease requested, it is activated by starting
    ///         the authorization decrease delay.
    function joinSortitionPool(
        Data storage self,
        IStaking tokenStaking,
        SortitionPool sortitionPool
    ) internal {
        address operator = msg.sender;

        address stakingProvider = self.operatorToStakingProvider[operator];
        require(stakingProvider != address(0), "Unknown operator");

        uint96 authorizedStake = tokenStaking.authorizedStake(
            stakingProvider,
            address(this)
        );

        AuthorizationDecrease storage decrease = self.pendingDecreases[
            stakingProvider
        ];

        uint96 eligibleStake = authorizedStake - decrease.decreasingBy;

        require(
            eligibleStake >= self.parameters.minimumAuthorization,
            "Authorization below the minimum"
        );

        sortitionPool.insertOperator(operator, eligibleStake);

        if (decrease.decreasingAt == type(uint64).max) {
            decrease.decreasingAt =
                // solhint-disable-next-line not-rely-on-time
                uint64(block.timestamp) +
                self.parameters.authorizationDecreaseDelay;
        }
    }

    /// @notice Updates status of the operator in the sortition pool. If there
    ///         was an authorization decrease requested, it is activated by
    ///         starting the authorization decrease delay.
    ///         Function reverts if the operator is not known.
    function updateOperatorStatus(
        Data storage self,
        IStaking tokenStaking,
        SortitionPool sortitionPool,
        address operator
    ) internal {
        address stakingProvider = self.operatorToStakingProvider[operator];
        require(stakingProvider != address(0), "Unknown operator");

        uint96 authorizedStake = tokenStaking.authorizedStake(
            stakingProvider,
            address(this)
        );

        AuthorizationDecrease storage decrease = self.pendingDecreases[
            stakingProvider
        ];

        if (sortitionPool.isOperatorInPool(operator)) {
            uint96 eligibleStake = authorizedStake - decrease.decreasingBy;
            sortitionPool.updateOperatorStatus(operator, eligibleStake);
        }

        if (decrease.decreasingAt == type(uint64).max) {
            decrease.decreasingAt =
                // solhint-disable-next-line not-rely-on-time
                uint64(block.timestamp) +
                self.parameters.authorizationDecreaseDelay;
        }
    }

    /// @notice Returns the remaining time in seconds that needs to pass before
    ///         the requested authorization decrease can be approved.
    ///         If the sortition pool state was not updated yet by the operator
    ///         after requesting the authorization decrease, returns
    ///         `type(uint64).max`.
    function remainingAuthorizationDecreaseDelay(
        Data storage self,
        address stakingProvider
    ) internal view returns (uint64) {
        AuthorizationDecrease storage decrease = self.pendingDecreases[
            stakingProvider
        ];

        if (decrease.decreasingAt == type(uint64).max) {
            return type(uint64).max;
        }

        // solhint-disable-next-line not-rely-on-time
        uint64 _now = uint64(block.timestamp);
        return _now > decrease.decreasingAt ? 0 : decrease.decreasingAt - _now;
    }

    // TODO: involuntaryAuthorizationDecrease
    // TODO: isOperatorUpToDate
}
