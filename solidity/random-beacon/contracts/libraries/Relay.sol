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
//                           Trust math, not hardware.

pragma solidity ^0.8.6;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "./BLS.sol";
import "./Groups.sol";

// TODO: Documentation
library Relay {
    using SafeERC20 for IERC20;

    struct Request {
        uint256 id;
        uint256 startBlock;
        Groups.Group group;
        bytes previousEntry;
    }

    struct Data {
        uint256 requestCount;
        Request currentRequest;

        IERC20 tToken;
        uint256 relayRequestFee;

        uint256 groupSize;
        uint256 relayEntrySubmissionEligibilityDelay;
        uint256 relayEntryHardTimeout;
    }

    event RelayEntryRequested(
        uint256 indexed requestId,
        bytes groupPublicKey,
        bytes previousEntry
    );
    event RelayEntrySubmitted(uint256 indexed requestId, bytes entry);

    function requestEntry(
        Data storage self,
        Groups.Group memory group,
        bytes calldata previousEntry
    ) internal {
        require(!isRequestInProgress(self), "Another relay request in progress");

        self.tToken.safeTransferFrom(msg.sender, address(this), self.relayRequestFee);

        uint256 currentRequestId = ++self.requestCount;

        // TODO: Accepting and storing the whole Group object is not efficient
        //       as a lot of data is copied. Revisit once `Groups` library is
        //       ready.
        self.currentRequest = Request(
            currentRequestId,
            block.number,
            group,
            previousEntry
        );

        emit RelayEntryRequested(currentRequestId, group.groupPubKey, previousEntry);
    }

    function submitEntry(
        Data storage self,
        uint256 submitterIndex,
        bytes calldata entry
    ) internal {
        require(isRequestInProgress(self), "No relay request in progress");
        // TODO: Add timeout reporting.
        require(!isRequestTimedOut(self), "Relay request timed out");

        require(
            submitterIndex > 0 && submitterIndex <= self.groupSize,
            "Invalid submitter index"
        );
        require(
            self.currentRequest.group.members[submitterIndex - 1] == msg.sender,
            "Unexpected submitter index"
        );

        require(
            BLS.verify(
                self.currentRequest.group.groupPubKey,
                self.currentRequest.previousEntry,
                entry
            ),
            "Invalid entry"
        );

        (uint256 firstEligibleIndex, uint256 lastEligibleIndex) =
            getEligibilityRange(self, entry);
        require(
            isEligible(self, submitterIndex, firstEligibleIndex, lastEligibleIndex),
            "Submitter is not eligible"
        );

        // TODO: Use submitterIndex, firstEligibleIndex and lastEligibleIndex
        //       to prepare an array of addresses which should be kicked from
        //       the sortition pool for 2 weeks.

        // TODO: If soft timeout has elapsed, take bleeding into account
        //       and slash all members appropriately.

        delete self.currentRequest;

        emit RelayEntrySubmitted(self.requestCount, entry);
    }

    function isRequestInProgress(
        Data storage self
    ) internal view returns (bool) {
        return self.currentRequest.id != 0;
    }

    function isRequestTimedOut(
        Data storage self
    ) internal view returns (bool) {
        uint256 relayEntryTimeout =
            (self.groupSize * self.relayEntrySubmissionEligibilityDelay) +
            self.relayEntryHardTimeout;

        return isRequestInProgress(self) &&
            block.number > self.currentRequest.startBlock + relayEntryTimeout;
    }

    function getEligibilityRange(
        Data storage self,
        bytes calldata entry
    ) internal view returns (uint256 firstEligibleIndex, uint256 lastEligibleIndex) {
        uint256 groupSize = self.groupSize;

        // Modulo `groupSize` will give indexes in range <0, groupSize-1>
        // We count member indexes from `1` so we need to add `1` to the result.
        firstEligibleIndex = (uint256(keccak256(entry)) % groupSize) + 1;

        // Shift is computed by leveraging Solidity integer division which is
        // equivalent to floored division. That gives the desired result.
        // Shift value should be in range <0, groupSize-1> so we must cap
        // it explicitly.
        uint256 shift = (block.number - self.currentRequest.startBlock) /
            self.relayEntrySubmissionEligibilityDelay;
        shift = shift > groupSize - 1 ? groupSize - 1 : shift;

        // Last eligible index must be wrapped if their value is bigger than
        // the group size.
        lastEligibleIndex = firstEligibleIndex + shift;
        lastEligibleIndex = lastEligibleIndex > groupSize ?
            lastEligibleIndex - groupSize : lastEligibleIndex;

        return (firstEligibleIndex, lastEligibleIndex);
    }

    function isEligible(
        Data storage self,
        uint256 submitterIndex,
        uint256 firstEligibleIndex,
        uint256 lastEligibleIndex
    ) internal view returns (bool) {
        if (firstEligibleIndex <= lastEligibleIndex) {
            // First eligible index is equal or smaller than the last.
            // We just need to make sure the submitter index is in range
            // <firstEligibleIndex, lastEligibleIndex>.
            return firstEligibleIndex <= submitterIndex &&
                submitterIndex <= lastEligibleIndex;
        } else {
            // First eligible index is bigger than the last. We need to deal
            // with wrapped range and check whether the submitter index is
            // either in range <1, lastEligibleIndex> or
            // <firstEligibleIndex, groupSize>.
            return (1 <= submitterIndex && submitterIndex <= lastEligibleIndex) ||
                (firstEligibleIndex <= submitterIndex && submitterIndex <= self.groupSize);
        }
    }
}