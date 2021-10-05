// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

// TODO: Consider which functions can be internal or internal. What are implications
// to security.

/// @title DKG library
library DKG {
    struct Data {
        uint256 seed;
        uint256 groupSize;
        uint256 dkgSubmissionEligibilityDelay;
        uint256 startBlock;
    }

    /// DKG is in an invalid state. Expected in progress: `expectedInProgress`,
    /// but actual in progress: `actualInProgress`.
    error InvalidInProgressState(
        bool expectedInProgress,
        bool actualInProgress
    );

    /// DKG haven't timed out yet. Expected block: `expectedTimeoutBlock`,
    /// but actual block: `actualBlock`.
    error NotTimedOut(uint256 expectedTimeoutBlock, uint256 actualBlock);

    event DkgStarted(
        uint256 seed,
        uint256 groupSize,
        uint256 dkgSubmissionEligibilityDelay
    ); // TODO: Add more parameters
    event DkgTimedOut(uint256 seed);
    event DkgCompleted(uint256 seed);

    modifier assertInProgress(Data storage self, bool expectedValue) {
        if (isInProgress(self) != expectedValue)
            revert InvalidInProgressState(expectedValue, isInProgress(self));
        _;
    }

    modifier cleanup(Data storage self) {
        _;
        delete self.seed;
        delete self.groupSize;
        delete self.startBlock;
        delete self.dkgSubmissionEligibilityDelay;
    }

    function isInProgress(Data storage self) public view returns (bool) {
        return self.startBlock > 0;
    }

    function dkgTimeout(Data storage self) public view returns (uint256) {
        return self.groupSize * self.dkgSubmissionEligibilityDelay;
    }

    function start(
        Data storage self,
        uint256 seed,
        uint256 groupSize,
        uint256 dkgSubmissionEligibilityDelay
    ) internal assertInProgress(self, false) {
        self.seed = seed;
        self.groupSize = groupSize;
        self.dkgSubmissionEligibilityDelay = dkgSubmissionEligibilityDelay;
        self.startBlock = block.number;

        emit DkgStarted(seed, groupSize, dkgSubmissionEligibilityDelay);
    }

    function notifyTimeout(Data storage self)
        internal
        assertInProgress(self, true)
        cleanup(self)
    {
        if (block.number <= self.startBlock + dkgTimeout(self))
            revert NotTimedOut(
                self.startBlock + dkgTimeout(self) + 1,
                block.number
            );

        emit DkgTimedOut(self.seed);
    }

    function finish(Data storage self)
        internal
        assertInProgress(self, true)
        cleanup(self)
    {
        emit DkgCompleted(self.seed);
    }
}
