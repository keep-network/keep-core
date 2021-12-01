// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "@keep-network/sortition-pools/contracts/SortitionPool.sol";

library Authorization {
    struct AuthorizationDecrease {
        uint96 decreasingTo; // amount
        uint64 decreasingAt; // timestamp
    }

    struct Parameters {
        uint96 minimumAuthorization;
        uint64 authorizationDecreaseDelay;
    }

    struct Update {
        address operator;
        uint96 authorization;
    }

    struct Queue {
        uint128 length;
        uint128 start;
        mapping(uint256 => Update) updates;
    }

    struct Data {
        Parameters parameters;
        Queue queue;
        // Address of the Sortition Pool contract.
        SortitionPool sortitionPool;
        mapping(address => AuthorizationDecrease) authorizationDecreaseRequests;
    }

    /// @notice Initializes the sortitionPool parameter. Can be performed only once.
    /// @param _sortitionPool Value of the parameter.
    function initSortitionPool(Data storage self, SortitionPool _sortitionPool)
        internal
    {
        require(
            address(self.sortitionPool) == address(0),
            "Sortition Pool address already set"
        );

        self.sortitionPool = _sortitionPool;
    }

    function setMinimumAuthorization(
        Data storage self,
        uint256 _minimumAuthorization
    ) internal {
        self.parameters.minimumAuthorization = uint96(_minimumAuthorization);
    }

    function setAuthorizationDecreaseDelay(
        Data storage self,
        uint256 _authorizationDecreaseDelay
    ) internal {
        self.parameters.authorizationDecreaseDelay = uint64(
            _authorizationDecreaseDelay
        );
    }

    /// @notice Used by T staking contract to inform the random beacon that the
    ///         authorized amount for the given operator increased.
    /// @dev ONLY T STAKING
    function authorizationIncreased(
        Data storage self,
        address operator,
        uint96 amount
    ) external {
        require(
            amount >= self.parameters.minimumAuthorization,
            "Below minimum authorization amount"
        );

        updateAuthorization(self, operator, amount);
    }

    // pass sortition pool as a parameter
    function authorizationDecreaseRequested(
        Data storage self,
        address operator,
        uint96 amount
    ) external {
        require(
            amount == 0 || amount >= self.parameters.minimumAuthorization,
            "Authorization amount should be 0 or above the minimum"
        );

        self.authorizationDecreaseRequests[operator] = AuthorizationDecrease(
            amount,
            // solhint-disable-next-line not-rely-on-time
            uint64(block.timestamp)
        );
    }

    function approveAuthorizationDecrease(Data storage self, address operator)
        internal
    {
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

        updateAuthorization(self, operator, request.decreasingTo);
        delete self.authorizationDecreaseRequests[operator];
    }

    function updateAuthorization(
        Data storage self,
        address operator,
        uint96 authorization
    ) internal {
        SortitionPool pool = self.sortitionPool;

        // When pool is locked, add the update to the queue.
        // When pool is unlocked, process the oldest item from the queue
        // before processing this update.
        if (pool.isLocked()) {
            addToQueue(self.queue, operator, authorization);
        } else {
            incrementQueue(self.queue, pool);
            updatePool(pool, operator, authorization);
        }
    }

    function addToQueue(
        Queue storage q,
        address operator,
        uint96 authorization
    ) internal {
        uint128 newLength = q.length + 1;
        q.updates[uint256(newLength)] = Update(operator, authorization);
        q.length = newLength;
    }

    function incrementQueue(Queue storage q, SortitionPool pool) internal {
        uint256 qStart = uint256(q.start);
        uint256 qLen = uint256(q.length);
        if (qLen > qStart) {
            Update memory u = q.updates[qStart];
            updatePool(pool, u.operator, u.authorization);
            delete q.updates[qStart];
            q.start = uint128(qStart + 1);
        }
    }

    function processQueue(Data storage self, uint256 times) internal {
        for (uint256 i = 0; i < times; i++) {
            incrementQueue(self.queue, self.sortitionPool);
        }
    }

    function updatePool(
        SortitionPool pool,
        address operator,
        uint96 authorization
    ) internal {
        if (pool.isOperatorInPool(operator)) {
            pool.updateOperatorStatus(operator, uint256(authorization));
        } else {
            pool.insertOperator(operator, uint256(authorization));
        }
    }
}
