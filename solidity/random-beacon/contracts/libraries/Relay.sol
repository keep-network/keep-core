// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

// TODO: Documentation
// TODO: Unit tests
library Relay {
    struct Request {
        uint256 id;
        uint256 startBlock;
        uint256 groupIndex;
        bytes previousEntry;
    }

    struct Data {
        uint256 requestCount;
        Request currentRequest;
        uint256 groupSize;
        uint256 entrySubmissionBlockStep;
        uint256 bleedingPeriod;
    }

    function requestEntry(
        Data storage self,
        uint256 groupIndex,
        bytes calldata previousEntry
    ) internal {
        require(!isRequestInProgress(self), "Another relay request in progress");

        self.requestCount++;

        self.currentRequest = Request(
            self.requestCount,
            block.number,
            groupIndex,
            previousEntry
        );

        // TODO: Emit event.
    }

    function submitEntry(
        Data storage self,
        bytes calldata entry
    ) internal {
        require(isRequestInProgress(self), "No relay request in progress");
        require(!isRequestTimedOut(self), "Relay request timed out");

        // TODO: Verify entry and submitter eligibility.
        // TODO: Implement bleeding.

        delete self.currentRequest;

        // TODO: Emit event.
    }

    function isRequestInProgress(
        Data storage self
    ) internal view returns (bool) {
        return self.currentRequest.id != 0;
    }

    function isRequestTimedOut(
        Data storage self
    ) internal view returns (bool) {
        // TODO: Bleeding period should be counted in.
        return isRequestInProgress(self) &&
            block.number > self.groupSize * self.entrySubmissionBlockStep;
    }
}