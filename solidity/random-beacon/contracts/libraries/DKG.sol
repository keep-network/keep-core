// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

// TODO: Consider which functions can be internal or internal. What are implications
// to security.

/// @title DKG library
library DKG {
    struct Data {
        uint256 seed;
        uint256 groupSize;
        uint256 timeoutDuration;
        uint256 startTimestamp;
    }

    /// DKG is in an invalid state. Expected in progress: `expectedInProgress`,
    /// but actual in progress: `actualInProgress`.
    error InvalidInProgressState(
        bool expectedInProgress,
        bool actualInProgress
    );

    /// DKG haven't timed out yet. Expected timestamp: `expectedTimeoutTimestamp`,
    /// but actual timestamp: `actualTimestamp`.
    error NotTimedOut(
        uint256 expectedTimeoutTimestamp,
        uint256 actualTimestamp
    );

    event DkgStarted(uint256 seed, uint256 groupSize, uint256 timeoutDuration);
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
        delete self.timeoutDuration;
        delete self.startTimestamp;
    }

    function isInProgress(Data storage self) public view returns (bool) {
        return self.startTimestamp > 0;
    }

    function start(
        Data storage self,
        uint256 seed,
        uint256 groupSize,
        uint256 timeoutDuration
    ) internal assertInProgress(self, false) {
        self.seed = seed;
        self.groupSize = groupSize;
        self.timeoutDuration = timeoutDuration;
        self.startTimestamp = block.timestamp;

        emit DkgStarted(seed, groupSize, timeoutDuration);
    }

    function notifyTimeout(Data storage self)
        internal
        assertInProgress(self, true)
        cleanup(self)
    {
        if (block.timestamp < self.startTimestamp + self.timeoutDuration)
            revert NotTimedOut(
                self.startTimestamp + self.timeoutDuration,
                block.timestamp
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
