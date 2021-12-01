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
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "./BLS.sol";
import "./Groups.sol";
import "./Submission.sol";

library Relay {
    using SafeERC20 for IERC20;

    struct Request {
        // Request identifier.
        uint64 id;
        // Identifier of group responsible for signing.
        uint64 groupId;
        // Request start block.
        uint128 startBlock;
    }

    struct Data {
        // Total count of all requests.
        uint64 requestCount;
        // Previous entry value.
        bytes previousEntry;
        // Data of current request.
        Request currentRequest;
        // Fee paid by the relay requester.
        uint256 relayRequestFee;
        // The number of blocks it takes for a group member to become
        // eligible to submit the relay entry.
        uint256 relayEntrySubmissionEligibilityDelay;
        // Hard timeout in blocks for a group to submit the relay entry.
        uint256 relayEntryHardTimeout;
        // Slashing amount for not submitting relay entry
        uint256 relayEntrySubmissionFailureSlashingAmount;
    }

    /// @notice Target DKG group size in the threshold relay. A group has
    ///         the target size if all their members behaved properly during
    ///         group formation. Actual group size can be lower in groups
    ///         with proven misbehaved members.
    uint256 public constant dkgGroupSize = 64;

    /// @notice Seed used as the first relay entry value.
    /// It's a G1 point G * PI =
    /// G * 31415926535897932384626433832795028841971693993751058209749445923078164062862
    /// Where G is the generator of G1 abstract cyclic group.
    bytes public constant relaySeed =
        hex"15c30f4b6cf6dbbcbdcc10fe22f54c8170aea44e198139b776d512d8f027319a1b9e8bfaf1383978231ce98e42bafc8129f473fc993cf60ce327f7d223460663";

    event RelayEntryRequested(
        uint256 indexed requestId,
        uint64 groupId,
        bytes previousEntry
    );

    event RelayEntrySubmitted(
        uint256 indexed requestId,
        address submitter,
        bytes entry
    );

    event RelayEntryTimedOut(
        uint256 indexed requestId,
        uint64 terminatedGroupId
    );

    /// @notice Initializes the very first `previousEntry` with an initial
    ///         `relaySeed` value. Can be performed only once.
    function initSeedEntry(Data storage self) internal {
        require(
            self.previousEntry.length == 0,
            "Seed entry already initialized"
        );
        self.previousEntry = relaySeed;
    }

    /// @notice Creates a request to generate a new relay entry, which will
    ///         include a random number (by signing the previous entry's
    ///         random number).
    /// @param groupId Identifier of the group chosen to handle the request.
    function requestEntry(Data storage self, uint64 groupId) internal {
        require(
            !isRequestInProgress(self),
            "Another relay request in progress"
        );

        uint64 currentRequestId = ++self.requestCount;

        self.currentRequest = Request(
            currentRequestId,
            groupId,
            uint128(block.number)
        );

        emit RelayEntryRequested(currentRequestId, groupId, self.previousEntry);
    }

    /// @notice Creates a new relay entry.
    /// @param entry Group BLS signature over the previous entry.
    /// @param groupPubKey Public key of the group which signed the relay entry.
    /// @return slashingAmount Amount by which group members should be slashed
    ///         in case the relay entry was submitted after the soft timeout.
    ///         The value is zero if entry was submitted on time.
    function submitEntry(
        Data storage self,
        bytes calldata entry,
        bytes storage groupPubKey
    ) internal returns (uint256 slashingAmount) {
        require(isRequestInProgress(self), "No relay request in progress");
        require(!hasRequestTimedOut(self), "Relay request timed out");

        require(
            BLS.verify(groupPubKey, self.previousEntry, entry),
            "Invalid entry"
        );

        // If the soft timeout has been exceeded apply stake slashing for
        // all group members. Note that `getSlashingFactor` returns the
        // factor multiplied by 1e18 to avoid precision loss. In that case
        // the final result needs to be divided by 1e18.
        slashingAmount =
            (getSlashingFactor(self, dkgGroupSize) *
                self.relayEntrySubmissionFailureSlashingAmount) /
            1e18;

        emit RelayEntrySubmitted(self.currentRequest.id, msg.sender, entry);

        self.previousEntry = entry;
        delete self.currentRequest;

        return slashingAmount;
    }

    /// @notice Set relayRequestFee parameter.
    /// @param newRelayRequestFee New value of the parameter.
    function setRelayRequestFee(Data storage self, uint256 newRelayRequestFee)
        internal
    {
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

        self
            .relayEntrySubmissionEligibilityDelay = newRelayEntrySubmissionEligibilityDelay;
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

    /// @notice Set relayEntrySubmissionFailureSlashingAmount parameter.
    /// @param newRelayEntrySubmissionFailureSlashingAmount New value of
    ///        the parameter.
    function setRelayEntrySubmissionFailureSlashingAmount(
        Data storage self,
        uint256 newRelayEntrySubmissionFailureSlashingAmount
    ) internal {
        require(!isRequestInProgress(self), "Relay request in progress");

        self
            .relayEntrySubmissionFailureSlashingAmount = newRelayEntrySubmissionFailureSlashingAmount;
    }

    /// @notice Retries the current relay request in case a relay entry
    ///         timeout was reported.
    /// @param newGroupId ID of the group chosen to retry the current request.
    function retryOnEntryTimeout(Data storage self, uint64 newGroupId)
        internal
    {
        require(hasRequestTimedOut(self), "Relay request did not time out");

        Request memory currentRequest = self.currentRequest;
        uint64 previousGroupId = currentRequest.groupId;

        emit RelayEntryTimedOut(currentRequest.id, previousGroupId);

        self.currentRequest = Request(
            currentRequest.id,
            newGroupId,
            uint128(block.number)
        );

        emit RelayEntryRequested(
            currentRequest.id,
            newGroupId,
            self.previousEntry
        );
    }

    /// @notice Cleans up the current relay request in case a relay entry
    ///         timeout was reported.
    function cleanupOnEntryTimeout(Data storage self) internal {
        require(hasRequestTimedOut(self), "Relay request did not time out");

        emit RelayEntryTimedOut(
            self.currentRequest.id,
            self.currentRequest.groupId
        );

        delete self.currentRequest;
    }

    /// @notice Returns whether a relay entry request is currently in progress.
    /// @return True if there is a request in progress. False otherwise.
    function isRequestInProgress(Data storage self)
        internal
        view
        returns (bool)
    {
        return self.currentRequest.id != 0;
    }

    /// @notice Returns whether the current relay request has timed out.
    /// @return True if the request timed out. False otherwise.
    function hasRequestTimedOut(Data storage self)
        internal
        view
        returns (bool)
    {
        uint256 _relayEntryTimeout = (dkgGroupSize *
            self.relayEntrySubmissionEligibilityDelay) +
            self.relayEntryHardTimeout;

        return
            isRequestInProgress(self) &&
            block.number > self.currentRequest.startBlock + _relayEntryTimeout;
    }

    /// @notice Computes the slashing factor which should be used during
    ///         slashing of the group which exceeded the soft timeout.
    /// @dev This function doesn't use the constant `groupSize` directly and
    ///      use a `_groupSize` parameter instead to facilitate testing.
    ///      Big group sizes in tests make readability worse and dramatically
    ///      increase the time of execution.
    /// @param _groupSize _groupSize Group size.
    /// @return A slashing factor represented as a fraction multiplied by 1e18
    ///         to avoid precision loss. When using this factor during slashing
    ///         amount computations, the final result should be divided by
    ///         1e18 to obtain a proper result. The slashing factor is
    ///         always in range <0, 1e18>.
    function getSlashingFactor(Data storage self, uint256 _groupSize)
        internal
        view
        returns (uint256)
    {
        uint256 softTimeoutBlock = self.currentRequest.startBlock +
            (_groupSize * self.relayEntrySubmissionEligibilityDelay);

        if (block.number > softTimeoutBlock) {
            uint256 submissionDelay = block.number - softTimeoutBlock;
            uint256 slashingFactor = (submissionDelay * 1e18) /
                self.relayEntryHardTimeout;
            return slashingFactor > 1e18 ? 1e18 : slashingFactor;
        }

        return 0;
    }
}
