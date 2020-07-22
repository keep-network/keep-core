pragma solidity 0.5.17;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";

/// @notice TokenStaking contract library allowing to perform two-step stake
/// top-ups for existing delegations.
/// Top-up is a two-step process: it is initiated with a declared top-up value
/// and after waiting for at least the initialization period it can be
/// committed.
library TopUps {
    using SafeMath for uint256;

    struct TopUp {
        uint256 amount;
        uint256 createdAt;
    }

    struct Storage {
        // operator -> TopUp
        mapping(address => TopUp) topUps;
    }

    /// @notice Initiates top-up of the given value for tokens delegated to
    /// the provided operator. If there is an existing top-up still
    /// initializing, top-up values are summed up and initialization period
    /// is set to the current block timestamp.
    /// @param value Top-up value, the number of tokens added to the stake.
    /// @param operator Operator The operator with existing delegation to which
    /// the tokens should be added to.
    function initiate(
        Storage storage self,
        uint256 value,
        address operator
    ) public {
        TopUp memory awaiting = self.topUps[operator];
        self.topUps[operator] = TopUp(awaiting.amount.add(value), now);
    }

    /// @notice Commits the top-up if it passed the initialization period.
    /// Tokens are added to the stake once the top-up is committed.
    /// @param operator Operator The operator with a pending stake top-up.
    /// @param initializationPeriod Stake initialization period.
    function commit(
        Storage storage self,
        address operator,
        uint256 initializationPeriod
    ) public returns (uint256) {
        TopUp memory topUp = self.topUps[operator];
        require(topUp.amount > 0, "No top up to commit");
        require(
            now > topUp.createdAt.add(initializationPeriod),
            "Stake is initializing"
        );

        delete self.topUps[operator];
        return topUp.amount;
    }

    /// @notice Force-commits existing top-up without taking into account
    /// stake initialization period. If there is no pending top-up for the
    /// operator, function does nothing. This function should be used when
    /// the stake is recovered to return tokens from a pending top-up.
    /// @param operator Operator The operator from which the stake is recovered.
    function forceCommit(
        Storage storage self,
        address operator
    ) public returns (uint256) {
        TopUp memory topUp = self.topUps[operator];
        if (topUp.amount == 0) {
            return 0;
        }

        delete self.topUps[operator];
        return topUp.amount;
    }
}
