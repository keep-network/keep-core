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

// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

import "@openzeppelin/contracts/access/Ownable.sol";

/// @title Random beacon
/// @notice Random beacon contract which represents the on-chain part of the
///         random beacon functionality
/// @dev Should be owned by the random beacon governance contract
contract RandomBeacon is Ownable {
    /// @notice Relay request fee in T
    uint256 public relayRequestFee;

    /// @notice The number of blocks for a member to become eligible to submit
    ///         relay entry
    uint256 public relayEntrySubmissionEligibilityDelay;

    /// @notice Hard timeout for a relay entry
    uint256 public relayEntryHardTimeout;

    /// @notice Callback gas limit
    uint256 public callbackGasLimit;

    /// @notice The frequency of a new group creation
    uint256 public groupCreationFrequency;

    /// @notice Group lifetime
    uint256 public groupLifetime;

    /// @notice The number of blocks for which a DKG result can be challenged
    uint256 public dkgResultChallengePeriodLength;

    /// @notice The number of blocks for a member to become eligible to submit
    ///         DKG result
    uint256 public dkgSubmissionEligibilityDelay;

    /// @notice Reward for submitting DKG result
    uint256 public dkgResultSubmissionReward;

    /// @notice Reward for unlocking the sortition pool if DKG timed out
    uint256 public sortitionPoolUnlockingReward;

    /// @notice Slashing amount for not submitting relay entry
    uint256 public relayEntrySubmissionFailureSlashingAmount;

    /// @notice Slashing amount for submitting malicious DKG result
    uint256 public maliciousDkgResultSlashingAmount;

    /// @notice Updates the values of relay entry parameters
    /// @dev Can be called only by the contract owner, which should be the
    ///      random beacon governance contract
    /// @param _relayRequestFee New relay request fee
    /// @param _relayEntrySubmissionEligibilityDelay New relay entry submission
    ///        eligibility delay
    /// @param _relayEntryHardTimeout New relay entry hard timeout
    /// @param _callbackGasLimit New callback gas limit
    function updateRelayEntryParameters(
        uint256 _relayRequestFee,
        uint256 _relayEntrySubmissionEligibilityDelay,
        uint256 _relayEntryHardTimeout,
        uint256 _callbackGasLimit
    ) external onlyOwner {
        relayRequestFee = _relayRequestFee;
        relayEntrySubmissionEligibilityDelay = _relayEntrySubmissionEligibilityDelay;
        relayEntryHardTimeout = _relayEntryHardTimeout;
        callbackGasLimit = _callbackGasLimit;
    }

    /// @notice Updates the values of group creation parameters
    /// @dev Can be called only by the contract owner, which should be the
    ///      random beacon governance contract
    /// @param _groupCreationFrequency New group creation frequency
    /// @param _groupLifetime New group lifetime
    /// @param _dkgResultChallengePeriodLength New DKG result challenge period
    ///        length
    /// @param _dkgSubmissionEligibilityDelay New DKG submission eligibility
    ///        delay
    function updateGroupCreationParameters(
        uint256 _groupCreationFrequency,
        uint256 _groupLifetime,
        uint256 _dkgResultChallengePeriodLength,
        uint256 _dkgSubmissionEligibilityDelay
    ) external onlyOwner {
        groupCreationFrequency = _groupCreationFrequency;
        groupLifetime = _groupLifetime;
        dkgResultChallengePeriodLength = _dkgResultChallengePeriodLength;
        dkgSubmissionEligibilityDelay = _dkgSubmissionEligibilityDelay;
    }

    /// @notice Updates the values of reward parameters
    /// @dev Can be called only by the contract owner, which should be the
    ///      random beacon governance contract
    /// @param _dkgResultSubmissionReward New DKG result submission reward
    /// @param _sortitionPoolUnlockingReward New sortition pool unlocking reward
    function updateRewardParameters(
        uint256 _dkgResultSubmissionReward,
        uint256 _sortitionPoolUnlockingReward
    ) external onlyOwner {
        dkgResultSubmissionReward = _dkgResultSubmissionReward;
        sortitionPoolUnlockingReward = _sortitionPoolUnlockingReward;
    }

    /// @notice Updates the values of slashing parameters
    /// @dev Can be called only by the contract owner, which should be the
    ///      random beacon governance contract
    /// @param _relayEntrySubmissionFailureSlashingAmount New relay entry
    ///        submission failure amount
    /// @param _maliciousDkgResultSlashingAmount New malicious DKG result
    ///        slashing amount
    function updateSlashingParameters(
        uint256 _relayEntrySubmissionFailureSlashingAmount,
        uint256 _maliciousDkgResultSlashingAmount
    ) external onlyOwner {
        relayEntrySubmissionFailureSlashingAmount = _relayEntrySubmissionFailureSlashingAmount;
        maliciousDkgResultSlashingAmount = _maliciousDkgResultSlashingAmount;
    }
}
