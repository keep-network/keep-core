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
import {ISortitionPool, IStaking} from "../RandomBeacon.sol";

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
        // Address of the Sortition Pool contract.
        ISortitionPool sortitionPool;
        // Address of the T token contract.
        IERC20 tToken;
        // Address of the staking contract.
        IStaking staking;
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

    /// @notice Ideal size of a group in the threshold relay. A group has
    ///         an ideal size if all their members behaved properly during
    ///         group formation. Actual group size can be lower in groups
    ///         with proven misbehaved members.
    uint256 public constant idealGroupSize = 64;

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

    event RelayEntrySubmitted(uint256 indexed requestId, bytes entry);

    event RelayEntryTimedOut(uint256 indexed requestId);

    /// @notice Initializes the very first `previousEntry` with an initial
    ///         `relaySeed` value. Can be performed only once.
    function initSeedEntry(Data storage self) internal {
        require(
            self.previousEntry.length == 0,
            "Seed entry already initialized"
        );
        self.previousEntry = relaySeed;
    }

    /// @notice Initializes the sortitionPool parameter. Can be performed
    ///         only once.
    /// @param _sortitionPool Value of the parameter.
    function initSortitionPool(Data storage self, ISortitionPool _sortitionPool)
        internal
    {
        require(
            address(self.sortitionPool) == address(0),
            "Sortition Pool address already set"
        );

        self.sortitionPool = _sortitionPool;
    }

    /// @notice Initializes the tToken parameter. Can be performed only once.
    /// @param _tToken Value of the parameter.
    function initTToken(Data storage self, IERC20 _tToken) internal {
        require(
            address(self.tToken) == address(0),
            "T token address already set"
        );

        self.tToken = _tToken;
    }

    /// @notice Initializes the staking parameter. Can be performed
    ///         only once.
    /// @param _staking Value of the parameter.
    function initStaking(Data storage self, IStaking _staking) internal {
        require(
            address(self.staking) == address(0),
            "Staking address already set"
        );

        self.staking = _staking;
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

        // slither-disable-next-line reentrancy-events
        self.tToken.safeTransferFrom(
            msg.sender,
            address(this),
            self.relayRequestFee
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
    /// @param submitterIndex Index of the entry submitter.
    /// @param entry Group BLS signature over the previous entry.
    /// @param group Group data.
    /// @return inactiveMembers Array of members IDs which should be considered
    ///         inactive  for not submitting relay entry on their turn.
    function submitEntry(
        Data storage self,
        uint256 submitterIndex,
        bytes calldata entry,
        Groups.Group memory group
    ) internal returns (uint32[] memory inactiveMembers) {
        require(isRequestInProgress(self), "No relay request in progress");
        require(!hasRequestTimedOut(self), "Relay request timed out");

        uint256 actualGroupSize = group.members.length;

        require(
            submitterIndex > 0 && submitterIndex <= actualGroupSize,
            "Invalid submitter index"
        );
        require(
            self.sortitionPool.getIDOperator(
                group.members[submitterIndex - 1]
            ) == msg.sender,
            "Unexpected submitter index"
        );

        (
            uint256 firstEligibleIndex,
            uint256 lastEligibleIndex
        ) = getEligibilityRange(self, entry, actualGroupSize);
        require(
            isEligible(
                self,
                submitterIndex,
                firstEligibleIndex,
                lastEligibleIndex
            ),
            "Submitter is not eligible"
        );

        require(
            BLS.verify(group.groupPubKey, self.previousEntry, entry),
            "Invalid entry"
        );

        // Get the list of members IDs which should be considered inactive due
        // to not submitting the entry on their turn.
        inactiveMembers = getInactiveMembers(
            self,
            submitterIndex,
            firstEligibleIndex,
            group.members
        );

        // If the soft timeout has been exceeded apply stake slashing for
        // all group members. Note that `getSlashingFactor` returns the
        // factor multiplied by 1e18 to avoid precision loss. In that case
        // the final result needs to be divided by 1e18.
        uint256 slashingAmount = (getSlashingFactor(self, idealGroupSize) *
            self.relayEntrySubmissionFailureSlashingAmount) / 1e18;

        // TODO: This call will be removed from here in the follow-up PR.
        if (slashingAmount > 0) {
            // slither-disable-next-line reentrancy-events
            self.staking.slash(
                slashingAmount,
                self.sortitionPool.getIDOperators(group.members)
            );
        }

        self.previousEntry = entry;
        delete self.currentRequest;

        emit RelayEntrySubmitted(self.requestCount, entry);

        return inactiveMembers;
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
    /// @param groupId ID of the group chosen to retry the current request.
    function retryOnEntryTimeout(Data storage self, uint64 groupId) internal {
        require(hasRequestTimedOut(self), "Relay request did not time out");

        uint64 currentRequestId = self.currentRequest.id;

        emit RelayEntryTimedOut(currentRequestId);

        self.currentRequest = Request(
            currentRequestId,
            groupId,
            uint128(block.number)
        );

        emit RelayEntryRequested(currentRequestId, groupId, self.previousEntry);
    }

    /// @notice Cleans up the current relay request in case a relay entry
    ///         timeout was reported.
    function cleanupOnEntryTimeout(Data storage self) internal {
        require(hasRequestTimedOut(self), "Relay request did not time out");

        emit RelayEntryTimedOut(self.currentRequest.id);

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
        uint256 relayEntryTimeout = (idealGroupSize *
            self.relayEntrySubmissionEligibilityDelay) +
            self.relayEntryHardTimeout;

        return
            isRequestInProgress(self) &&
            block.number > self.currentRequest.startBlock + relayEntryTimeout;
    }

    /// @notice Determines the eligibility range for given relay entry basing on
    ///         current block number.
    /// @dev Parameters _entry and _groupSize are passed because the first
    ///      eligible index is computed as `_entry % _groupSize`. This function
    ///      doesn't use the constant `groupSize` directly to facilitate
    ///      testing. Big group sizes in tests make readability worse and
    ///      dramatically increase the time of execution.
    /// @param _entry Entry value for which the eligibility range should be
    ///        determined.
    /// @param _groupSize Group size for which eligibility range should be determined.
    /// @return firstEligibleIndex Index of the first member which is eligible
    ///         to submit the relay entry.
    /// @return lastEligibleIndex Index of the last member which is eligible
    ///         to submit the relay entry.
    function getEligibilityRange(
        Data storage self,
        bytes calldata _entry,
        uint256 _groupSize
    )
        internal
        view
        returns (uint256 firstEligibleIndex, uint256 lastEligibleIndex)
    {
        // Modulo `groupSize` will give indexes in range <0, groupSize-1>
        // We count member indexes from `1` so we need to add `1` to the result.
        firstEligibleIndex = (uint256(keccak256(_entry)) % _groupSize) + 1;

        // Shift is computed by leveraging Solidity integer division which is
        // equivalent to floored division. That gives the desired result.
        // Shift value should be in range <0, groupSize-1> so we must cap
        // it explicitly.
        uint256 shift = (block.number - self.currentRequest.startBlock) /
            self.relayEntrySubmissionEligibilityDelay;
        shift = shift > _groupSize - 1 ? _groupSize - 1 : shift;

        // Last eligible index must be wrapped if their value is bigger than
        // the group size. If wrapping occurs, the lastEligibleIndex is smaller
        // than the firstEligibleIndex. In that case, the eligibility queue
        // can look as follows: 1, 2 (last), 3, 4, 5, 6, 7 (first), 8.
        lastEligibleIndex = firstEligibleIndex + shift;
        lastEligibleIndex = lastEligibleIndex > _groupSize
            ? lastEligibleIndex - _groupSize
            : lastEligibleIndex;

        return (firstEligibleIndex, lastEligibleIndex);
    }

    /// @notice Returns whether the given submitter index is eligible to submit
    ///         a relay entry within given eligibility range.
    /// @param _submitterIndex Index of the submitter whose eligibility is checked.
    /// @param _firstEligibleIndex First index of the given eligibility range.
    /// @param _lastEligibleIndex Last index of the given eligibility range.
    /// @return True if eligible. False otherwise.
    function isEligible(
        /* solhint-disable-next-line no-unused-vars */
        Data storage self,
        uint256 _submitterIndex,
        uint256 _firstEligibleIndex,
        uint256 _lastEligibleIndex
    ) internal view returns (bool) {
        if (_firstEligibleIndex <= _lastEligibleIndex) {
            // First eligible index is equal or smaller than the last.
            // We just need to make sure the submitter index is in range
            // <firstEligibleIndex, lastEligibleIndex>.
            return
                _firstEligibleIndex <= _submitterIndex &&
                _submitterIndex <= _lastEligibleIndex;
        } else {
            // First eligible index is bigger than the last. We need to deal
            // with wrapped range and check whether the submitter index is
            // either in range <1, lastEligibleIndex> or
            // <firstEligibleIndex, groupSize>.
            return
                _submitterIndex <= _lastEligibleIndex ||
                _firstEligibleIndex <= _submitterIndex;
        }
    }

    /// @notice Determines a list of members which should be considered as
    ///         inactive due to not submitting a relay entry on their turn.
    ///         Inactive members are determined using the eligibility queue and
    ///         are taken from the <firstEligibleIndex, submitterIndex) range.
    ///         It also handles the `submitterIndex < firstEligibleIndex` case
    ///         and wraps the queue accordingly.
    /// @param _submitterIndex Index of the relay entry submitter.
    /// @param _firstEligibleIndex First index of the given eligibility range.
    /// @param _groupMembers IDs of the group members.
    /// @return An array of members IDs which should be  inactive due
    ///         to not submitting a relay entry on their turn.
    function getInactiveMembers(
        /* solhint-disable-next-line no-unused-vars */
        Data storage self,
        uint256 _submitterIndex,
        uint256 _firstEligibleIndex,
        uint32[] memory _groupMembers
    ) internal view returns (uint32[] memory) {
        uint256 _groupSize = _groupMembers.length;

        uint256 inactiveMembersCount = _submitterIndex >= _firstEligibleIndex
            ? _submitterIndex - _firstEligibleIndex
            : _groupSize - (_firstEligibleIndex - _submitterIndex);

        uint32[] memory inactiveMembersIDs = new uint32[](inactiveMembersCount);

        for (uint256 i = 0; i < inactiveMembersCount; i++) {
            uint256 memberIndex = _firstEligibleIndex + i;

            if (memberIndex > _groupSize) {
                memberIndex = memberIndex - _groupSize;
            }

            inactiveMembersIDs[i] = _groupMembers[memberIndex - 1];
        }

        return inactiveMembersIDs;
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
