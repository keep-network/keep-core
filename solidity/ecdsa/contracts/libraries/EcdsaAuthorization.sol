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

    struct Data {
        Parameters parameters;
        mapping(address => address) stakingProviderToOperator;
        mapping(address => address) operatorToStakingProvider;
    }

    event OperatorRegistered(
        address indexed stakingProvider,
        address indexed operator
    );

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
    ///         and must be unique.
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

        self.stakingProviderToOperator[stakingProvider] = operator;
        self.operatorToStakingProvider[operator] = stakingProvider;

        emit OperatorRegistered(stakingProvider, operator);
    }

    /// @notice Used by T staking contract to inform the application that the
    ///         authorized amount for the given operator increased. Can only be
    ///         called when the sortition pool is not locked. Increases in-pool
    ///         weight and rewards weight in the pool proportionally to the
    ///         authorized stake amount immediatelly. Reverts if the sortition
    ///         pool is locked or if the authorization amount is below the
    ///         minimum. If the operator is not known (`registerOperator`) was
    ///         not called, or if the operator is not in the sortition pool,
    ///         function is not executing any updates.
    /// @dev Should only be callable by T staking contract.
    function authorizationIncreased(
        Data storage self,
        SortitionPool sortitionPool,
        address stakingProvider,
        uint96,
        uint96 toAmount
    ) internal {
        require(
            toAmount >= self.parameters.minimumAuthorization,
            "Authorization below the minimum"
        );

        address operator = self.stakingProviderToOperator[stakingProvider];
        if (
            operator != address(0) && sortitionPool.isOperatorInPool(operator)
        ) {
            sortitionPool.updateOperatorStatus(operator, toAmount);
        }
    }

    /// @notice Lets the operator join the sortition pool. The operator address
    ///         must be known - before calling this function, it has to be
    ///         appointed by the staking provider by calling `registerOperator`.
    ///         Also, the operator must have the minimum authorization required
    ///         by ECDSA. Function reverts if there is no minimum stake
    ///         authorized or if the operator is not known.
    function joinSortitionPool(
        Data storage self,
        IStaking tokenStaking,
        SortitionPool sortitionPool
    ) internal {
        address operator = msg.sender;
        address stakingProvider = self.operatorToStakingProvider[operator];

        require(stakingProvider != address(0), "Unknown operator");

        uint256 authorizedStake = tokenStaking.authorizedStake(
            stakingProvider,
            address(this)
        );
        require(
            authorizedStake >= self.parameters.minimumAuthorization,
            "Authorization below the minimum"
        );

        sortitionPool.insertOperator(operator, authorizedStake);
    }
}
