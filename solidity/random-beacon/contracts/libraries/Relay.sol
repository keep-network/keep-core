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

pragma solidity ^0.8.9;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "./AltBn128.sol";
import "./BLS.sol";
import "./Groups.sol";
import "./Submission.sol";

library Relay {
    using SafeERC20 for IERC20;

    struct Data {
        // Total count of all requests.
        uint64 requestCount;
        // Data of current request.
        // Request identifier.
        uint64 currentRequestID;
        // Identifier of group responsible for signing.
        uint64 currentRequestGroupID;
        // Request start block.
        uint64 currentRequestStartBlock;
        // Previous entry value.
        AltBn128.G1Point previousEntry;
        // Fee paid by the relay requester.
        uint96 relayRequestFee;
        // The number of blocks it takes for a group member to become
        // eligible to submit the relay entry.
        uint32 relayEntrySubmissionEligibilityDelay;
        // Hard timeout in blocks for a group to submit the relay entry.
        uint32 relayEntryHardTimeout;
        // Slashing amount for not submitting relay entry
        uint96 relayEntrySubmissionFailureSlashingAmount;
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
            self.previousEntry.x == 0 && self.previousEntry.y == 0,
            "Seed entry already initialized"
        );
        self.previousEntry = AltBn128.g1Unmarshal(relaySeed);
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

        self.currentRequestID = currentRequestId;
        self.currentRequestGroupID = groupId;
        self.currentRequestStartBlock = uint64(block.number);

        emit RelayEntryRequested(
            currentRequestId,
            groupId,
            AltBn128.g1Marshal(self.previousEntry)
        );
    }

    /// @notice Creates a new relay entry.
    /// @param entry Group BLS signature over the previous entry.
    /// @param groupPubKey Public key of the group which signed the relay entry.
    function submitEntry(
        Data storage self,
        bytes calldata entry,
        bytes storage groupPubKey
    ) internal {
        require(isRequestInProgress(self), "No relay request in progress");
        require(
            BLS._verify(
                AltBn128.g2Unmarshal(groupPubKey),
                self.previousEntry,
                AltBn128.g1Unmarshal(entry)
            ),
            "Invalid entry"
        );

        emit RelayEntrySubmitted(self.currentRequestID, msg.sender, entry);

        self.previousEntry = AltBn128.g1Unmarshal(entry);
        self.currentRequestID = 0;
        self.currentRequestGroupID = 0;
        self.currentRequestStartBlock = 0;
    }

    /// @notice Calculates the slashing amount for all group members.
    /// @dev Must be used when a soft timeout was hit.
    /// @return Amount by which group members should be slashed
    ///         in case the relay entry was submitted after the soft timeout.
    function calculateSlashingAmount(Data storage self)
        internal
        returns (uint256)
    {
        uint256 softTimeoutBlock = self.currentRequestStartBlock +
            (dkgGroupSize * self.relayEntrySubmissionEligibilityDelay);

        uint256 submissionDelay = block.number - softTimeoutBlock;
        // The slashing amount is a result of a calculated portion of the submission
        // delay blocks. The max delay can be set up to relayEntryHardTimeout, which
        // in consequence sets the max slashing amount.
        if (submissionDelay > self.relayEntryHardTimeout) {
            submissionDelay = self.relayEntryHardTimeout;
        }

        return
            (submissionDelay * self.relayEntrySubmissionFailureSlashingAmount) /
            self.relayEntryHardTimeout;
    }

    /// @notice Set relayRequestFee parameter.
    /// @param newRelayRequestFee New value of the parameter.
    function setRelayRequestFee(Data storage self, uint256 newRelayRequestFee)
        internal
    {
        require(!isRequestInProgress(self), "Relay request in progress");

        self.relayRequestFee = uint96(newRelayRequestFee);
    }

    /// @notice Set relayEntrySubmissionEligibilityDelay parameter.
    /// @param newRelayEntrySubmissionEligibilityDelay New value of the parameter.
    function setRelayEntrySubmissionEligibilityDelay(
        Data storage self,
        uint256 newRelayEntrySubmissionEligibilityDelay
    ) internal {
        require(!isRequestInProgress(self), "Relay request in progress");

        self.relayEntrySubmissionEligibilityDelay = uint32(
            newRelayEntrySubmissionEligibilityDelay
        );
    }

    /// @notice Set relayEntryHardTimeout parameter.
    /// @param newRelayEntryHardTimeout New value of the parameter.
    function setRelayEntryHardTimeout(
        Data storage self,
        uint256 newRelayEntryHardTimeout
    ) internal {
        require(!isRequestInProgress(self), "Relay request in progress");

        self.relayEntryHardTimeout = uint32(newRelayEntryHardTimeout);
    }

    /// @notice Set relayEntrySubmissionFailureSlashingAmount parameter.
    /// @param newRelayEntrySubmissionFailureSlashingAmount New value of
    ///        the parameter.
    function setRelayEntrySubmissionFailureSlashingAmount(
        Data storage self,
        uint256 newRelayEntrySubmissionFailureSlashingAmount
    ) internal {
        require(!isRequestInProgress(self), "Relay request in progress");

        self.relayEntrySubmissionFailureSlashingAmount = uint96(
            newRelayEntrySubmissionFailureSlashingAmount
        );
    }

    /// @notice Retries the current relay request in case a relay entry
    ///         timeout was reported.
    /// @param newGroupId ID of the group chosen to retry the current request.
    function retryOnEntryTimeout(Data storage self, uint64 newGroupId)
        internal
    {
        require(hasRequestHardTimedOut(self), "Relay request did not time out");
        require(isRequestInProgress(self), "Relay request is not in progress");

        uint64 currentRequestId = self.currentRequestID;
        uint64 previousGroupId = self.currentRequestGroupID;

        emit RelayEntryTimedOut(currentRequestId, previousGroupId);

        self.currentRequestGroupID = newGroupId;
        self.currentRequestStartBlock = uint64(block.number);

        emit RelayEntryRequested(
            currentRequestId,
            newGroupId,
            AltBn128.g1Marshal(self.previousEntry)
        );
    }

    /// @notice Cleans up the current relay request in case a relay entry
    ///         timeout was reported.
    function cleanupOnEntryTimeout(Data storage self) internal {
        require(hasRequestHardTimedOut(self), "Relay request did not time out");
        require(isRequestInProgress(self), "Relay request is not in progress");

        emit RelayEntryTimedOut(
            self.currentRequestID,
            self.currentRequestGroupID
        );

        self.currentRequestID = 0;
        self.currentRequestGroupID = 0;
        self.currentRequestStartBlock = 0;
    }

    /// @notice Returns whether a relay entry request is currently in progress.
    /// @return True if there is a request in progress. False otherwise.
    function isRequestInProgress(Data storage self)
        internal
        view
        returns (bool)
    {
        return self.currentRequestID != 0;
    }

    /// @notice Returns whether the current relay request has hard timed out.
    /// @return True if the request hard timed out. False otherwise.
    function hasRequestHardTimedOut(Data storage self)
        internal
        view
        returns (bool)
    {
        uint256 hardTimeout = (dkgGroupSize *
            self.relayEntrySubmissionEligibilityDelay) +
            self.relayEntryHardTimeout;

        return block.number > self.currentRequestStartBlock + hardTimeout;
    }

    /// @notice Returns whether the current relay request has soft timed out.
    /// @return True if the request soft timed out. False otherwise.
    function hasRequestSoftTimedOut(Data storage self)
        internal
        view
        returns (bool)
    {
        uint256 softTimeout = dkgGroupSize *
            self.relayEntrySubmissionEligibilityDelay;

        return block.number > self.currentRequestStartBlock + softTimeout;
    }
}
