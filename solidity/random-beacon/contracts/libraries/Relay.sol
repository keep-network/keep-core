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

library Relay {
    using SafeERC20 for IERC20;

    struct Request {
        // Request identifier.
        uint256 id;
        // Request start block.
        uint256 startBlock;
        // Group responsible for signing as part of the request.
        Groups.Group group;
        // Previous entry value which should be signed as part of the request.
        bytes previousEntry;
    }

    struct Data {
        // Total count of all requests.
        uint256 requestCount;
        // Data of current request.
        Request currentRequest;
        // Address of the T token contract.
        IERC20 tToken;
        // Fee paid by the relay requester.
        uint256 relayRequestFee;
        // Size of the group performing signing.
        uint256 groupSize;
        // The number of blocks it takes for a group member to become
        // eligible to submit the relay entry.
        uint256 relayEntrySubmissionEligibilityDelay;
        // Hard timeout in blocks for a group to submit the relay entry.
        uint256 relayEntryHardTimeout;
    }

    event RelayEntryRequested(
        uint256 indexed requestId,
        bytes groupPublicKey,
        bytes previousEntry
    );
    event RelayEntrySubmitted(uint256 indexed requestId, bytes entry);

    /// @notice Creates a request to generate a new relay entry, which will
    ///         include a random number (by signing the previous entry's
    ///         random number).
    /// @param group Group chosen to handle the request.
    /// @param previousEntry Previous relay entry.
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

    /// @notice Creates a new relay entry.
    /// @param submitterIndex Index of the entry submitter.
    /// @param entry Group BLS signature over the previous entry.
    function submitEntry(
        Data storage self,
        uint256 submitterIndex,
        bytes calldata entry
    ) internal {
        require(isRequestInProgress(self), "No relay request in progress");
        // TODO: Add timeout reporting.
        require(!hasRequestTimedOut(self), "Relay request timed out");

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

    /// @notice Set relayRequestFee parameter.
    /// @param newRelayRequestFee New value of the parameter.
    function setRelayRequestFee(
        Data storage self,
        uint256 newRelayRequestFee
    ) internal {
        require(!isRequestInProgress(self), "Relay request in progress");

        self.relayRequestFee = newRelayRequestFee;
    }

    /// @notice Set relayEntrySubmissionEligibilityDelay parameter.
    /// @param newRelayEntrySubmissionEligibilityDelay New value of the parameter.
    function setRelayEntrySubmissionEligibilityDelay(
        Data storage self,
        uint256 newRelayEntrySubmissionEligibilityDelay
    ) internal {
        require(!isRequestInProgress(self), "Relay request in progress");

        self.relayEntrySubmissionEligibilityDelay =
            newRelayEntrySubmissionEligibilityDelay;
    }

    /// @notice Set relayEntryHardTimeout parameter.
    /// @param newRelayEntryHardTimeout New value of the parameter.
    function setRelayEntryHardTimeout(
        Data storage self,
        uint256 newRelayEntryHardTimeout
    ) internal {
        require(!isRequestInProgress(self), "Relay request in progress");

        self.relayEntryHardTimeout = newRelayEntryHardTimeout;
    }

    /// @notice Returns whether a relay entry request is currently in progress.
    /// @return True if there is a request in progress. False otherwise.
    function isRequestInProgress(
        Data storage self
    ) internal view returns (bool) {
        return self.currentRequest.id != 0;
    }

    /// @notice Returns whether the current relay request has timed out.
    /// @return True if the request timed out. False otherwise.
    function hasRequestTimedOut(
        Data storage self
    ) internal view returns (bool) {
        uint256 relayEntryTimeout =
            (self.groupSize * self.relayEntrySubmissionEligibilityDelay) +
            self.relayEntryHardTimeout;

        return isRequestInProgress(self) &&
            block.number > self.currentRequest.startBlock + relayEntryTimeout;
    }

    /// @notice Determines the eligibility range for given relay entry basing on
    ///         current block number.
    /// @param entry Entry value for which the eligibility range should be
    ///        determined.
    /// @return firstEligibleIndex Index of the first member which is eligible
    ///         to submit the relay entry.
    /// @return lastEligibleIndex Index of the last member which is eligible
    ///         to submit the relay entry.
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

    /// @notice Returns whether the given submitter index is eligible to submit
    ///         a relay entry within given eligibility range.
    /// @param submitterIndex Index of the submitter whose eligibility is checked.
    /// @param firstEligibleIndex First index of the given eligibility range.
    /// @param lastEligibleIndex Last index of the given eligibility range.
    /// @return True if eligible. False otherwise.
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