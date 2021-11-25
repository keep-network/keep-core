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

library Submission {
    /// @notice Determines the submitter eligibility range for given parameters.
    /// @param seed Value for which the eligibility range should be determined.
    /// @param protocolStartBlock Starting block of the protocol the submitter
    ///        eligibility range is calculated for.
    /// @param eligibilityDelay Delay in blocks each group member needs to wait
    ///        before becoming eligible.
    /// @param groupSize Group size for which eligibility range should be
    ///        determined.
    /// @return firstEligibleIndex Index of the first member which is eligible
    ///         to submit.
    /// @return lastEligibleIndex Index of the last member which is eligible
    ///         to submit.
    function getEligibilityRange(
        uint256 seed,
        uint256 protocolStartBlock,
        uint256 eligibilityDelay,
        uint256 groupSize
    )
        external
        view
        returns (uint256 firstEligibleIndex, uint256 lastEligibleIndex)
    {
        // Modulo `groupSize` will give indexes in range <0, groupSize-1>
        // We count member indexes from `1` so we need to add `1` to the result.
        firstEligibleIndex = (seed % groupSize) + 1;

        // Shift is computed by leveraging Solidity integer division which is
        // equivalent to floored division. That gives the desired result.
        // Shift value should be in range <0, groupSize-1> so we must cap
        // it explicitly.
        uint256 shift = (block.number - protocolStartBlock) / eligibilityDelay;
        shift = shift > groupSize - 1 ? groupSize - 1 : shift;

        // Last eligible index must be wrapped if their value is bigger than
        // the group size. If wrapping occurs, the lastEligibleIndex is smaller
        // than the firstEligibleIndex. In that case, the eligibility queue
        // can look as follows: 1, 2 (last), 3, 4, 5, 6, 7 (first), 8.
        lastEligibleIndex = firstEligibleIndex + shift;
        lastEligibleIndex = lastEligibleIndex > groupSize
            ? lastEligibleIndex - groupSize
            : lastEligibleIndex;

        return (firstEligibleIndex, lastEligibleIndex);
    }

    /// @notice Returns whether the given submitter index is eligible to submit
    ///         within given eligibility range.
    /// @param submitterIndex Index of the submitter whose eligibility is checked.
    /// @param firstEligibleIndex First index of the given eligibility range.
    /// @param lastEligibleIndex Last index of the given eligibility range.
    /// @return True if eligible. False otherwise.
    function isEligible(
        uint256 submitterIndex,
        uint256 firstEligibleIndex,
        uint256 lastEligibleIndex
    ) external view returns (bool) {
        if (firstEligibleIndex <= lastEligibleIndex) {
            // First eligible index is equal or smaller than the last.
            // We just need to make sure the submitter index is in range
            // <firstEligibleIndex, lastEligibleIndex>.
            return
                firstEligibleIndex <= submitterIndex &&
                submitterIndex <= lastEligibleIndex;
        } else {
            // First eligible index is bigger than the last. We need to deal
            // with wrapped range and check whether the submitter index is
            // either in range <1, lastEligibleIndex> or
            // <firstEligibleIndex, groupSize>.
            return
                submitterIndex <= lastEligibleIndex ||
                firstEligibleIndex <= submitterIndex;
        }
    }

    /// @notice Determines a list of members which should be considered as
    ///         inactive due to not submitting on their turn.
    ///         Inactive members are determined using the eligibility queue and
    ///         are taken from the <firstEligibleIndex, submitterIndex) range.
    ///         It also handles the `submitterIndex < firstEligibleIndex` case
    ///         and wraps the queue accordingly.
    /// @param submitterIndex Index of the submitter.
    /// @param firstEligibleIndex First index of the given eligibility range.
    /// @param groupMembers IDs of the group members.
    /// @return An array of members IDs which should be  inactive due
    ///         to not submitting on their turn.
    function getInactiveMembers(
        uint256 submitterIndex,
        uint256 firstEligibleIndex,
        uint32[] memory groupMembers
    ) external view returns (uint32[] memory) {
        uint256 groupSize = groupMembers.length;

        uint256 inactiveMembersCount = submitterIndex >= firstEligibleIndex
            ? submitterIndex - firstEligibleIndex
            : groupSize - (firstEligibleIndex - submitterIndex);

        uint32[] memory inactiveMembersIDs = new uint32[](inactiveMembersCount);

        for (uint256 i = 0; i < inactiveMembersCount; i++) {
            uint256 memberIndex = firstEligibleIndex + i;

            if (memberIndex > groupSize) {
                memberIndex = memberIndex - groupSize;
            }

            inactiveMembersIDs[i] = groupMembers[memberIndex - 1];
        }

        return inactiveMembersIDs;
    }
}
