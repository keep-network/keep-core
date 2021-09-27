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

import "./GovernanceUtils.sol";

import "@openzeppelin/contracts/access/Ownable.sol";

/// @title Random Beacon Governable
/// @notice This is a base contract that holds governable parameters for Random
///         Beacon as well as their update logic
contract RandomBeaconGovernable is Ownable {
    /// @notice The time delay that needs to pass between initializing and
    ///         finalizing update of any governable parameter in this contract.
    uint256 public constant GOVERNANCE_DELAY = 7 days;

    /// @notice Relay request fee in T
    uint256 public relayRequestFee;
    uint256 public newRelayRequestFee;
    uint256 public relayRequestFeeChangeInitiated;

    /// @notice Reward for submitting DKG result
    uint256 public dkgResultSubmissionReward;
    uint256 public newDkgResultSubmissionReward;
    uint256 public dkgResultSubmissionRewardChangeInitiated;

    /// @notice Reward for unlocking the sortition pool if DKG timed out
    uint256 public sortitionPoolUnlockingReward;
    uint256 public newSortitionPoolUnlockingReward;
    uint256 public sortitionPoolUnlockingRewardChangeInitiated;

    /// @notice Slashing amount for not submitting relay entry
    uint256 public relayEntrySubmissionFailureSlashingAmount;
    uint256 public newRelayEntrySubmissionFailureSlashingAmount;
    uint256 public relayEntrySubmissionFailureSlashingAmountChangeInitiated;

    /// @notice Slashing amount for submitting malicious DKG result
    uint256 public maliciousDkgResultSlashingAmount;
    uint256 public newMaliciousDkgResultSlashingAmount;
    uint256 public maliciousDkgResultSlashingAmountChangeInitiated;

    /// @notice The number of blocks for which a DKG result can be challenged
    uint256 public dkgResultChallengePeriodLength;
    uint256 public newDkgResultChallengePeriodLength;
    uint256 public dkgResultChallengePeriodLengthChangeInitiated;

    /// @notice The number of blocks for a member to become eligible to submit
    ///         relay entry
    uint256 public relayEntrySubmissionEligibilityDelay;
    uint256 public newRelayEntrySubmissionEligibilityDelay;
    uint256 public relayEntrySubmissionEligibilityDelayChangeInitiated;

    /// @notice The number of blocks for a member to become eligible to submit
    ///         DKG result
    uint256 public dkgSubmissionEligibilityDelay;
    uint256 public newDkgSubmissionEligibilityDelay;
    uint256 public dkgSubmissionEligibilityDelayChangeInitiated;

    /// @notice Hard timeout for a relay entry in blocks
    uint256 public relayEntryHardTimeout;
    uint256 public newRelayEntryHardTimeout;
    uint256 public relayEntryHardTimeoutChangeInitiated;

    /// @notice The frequency of a new group creation
    uint256 public groupCreationFrequency;
    uint256 public newGroupCreationFrequency;
    uint256 public groupCreationFrequencyChangeInitiated;

    /// @notice Group lifetime in seconds
    uint256 public groupLifetime;
    uint256 public newGroupLifetime;
    uint256 public groupLifetimeChangeInitiated;

    /// @notice Callback gas limit
    uint256 public callbackGasLimit;
    uint256 public newCallbackGasLimit;
    uint256 public callbackGasLimitChangeInitiated;

    event RelayRequestFeeUpdateStarted(
        uint256 relayRequestFee,
        uint256 timestamp
    );
    event RelayRequestFeeUpdated(uint256 relayRequestFee);

    event DkgResultSubmissionRewardUpdateStarted(
        uint256 dkgResultSubmissionReward,
        uint256 timestamp
    );
    event DkgResultSubmissionRewardUpdated(uint256 dkgResultSubmissionReward);

    event SortitionPoolUnlockingRewardUpdateStarted(
        uint256 sortitionPoolUnlockingReward,
        uint256 timestamp
    );
    event SortitionPoolUnlockingRewardUpdated(
        uint256 sortitionPoolUnlockingReward
    );

    event RelayEntrySubmissionFailureSlashingAmountUpdateStarted(
        uint256 relayEntrySubmissionFailureSlashingAmount,
        uint256 timestamp
    );
    event RelayEntrySubmissionFailureSlashingAmountUpdated(
        uint256 relayEntrySubmissionFailureSlashingAmount
    );

    event MaliciousDkgResultSlashingAmountUpdateStarted(
        uint256 maliciousDkgResultSlashingAmount,
        uint256 timestamp
    );
    event MaliciousDkgResultSlashingAmountUpdated(
        uint256 maliciousDkgResultSlashingAmount
    );

    event DkgResultChallengePeriodLengthUpdateStarted(
        uint256 dkgResultChallengePeriodLength,
        uint256 timestamp
    );
    event DkgResultChallengePeriodLengthUpdated(
        uint256 dkgResultChallengePeriodLength
    );

    event RelayEntrySubmissionEligibilityDelayUpdateStarted(
        uint256 relayEntrySubmissionEligibilityDelay,
        uint256 timestamp
    );
    event RelayEntrySubmissionEligibilityDelayUpdated(
        uint256 relayEntrySubmissionEligibilityDelay
    );

    event RelayEntryHardTimeoutUpdateStarted(
        uint256 relayEntryHardTimeout,
        uint256 timestamp
    );
    event RelayEntryHardTimeoutUpdated(uint256 relayEntryHardTimeout);

    event GroupCreationFrequencyUpdateStarted(
        uint256 groupCreationFrequency,
        uint256 timestamp
    );
    event GroupCreationFrequencyUpdated(uint256 groupCreationFrequency);

    event GroupLifetimeUpdateStarted(uint256 groupLifetime, uint256 timestamp);
    event GroupLifetimeUpdated(uint256 groupLifetime);

    event CallbackGasLimitUpdateStarted(
        uint256 callbackGasLimit,
        uint256 timestamp
    );
    event CallbackGasLimitUpdated(uint256 callbackGasLimit);

    /// @notice Reverts if called before the governance delay elapses.
    /// @param changeInitiatedTimestamp Timestamp indicating the beginning
    ///        of the change.
    modifier onlyAfterGovernanceDelay(uint256 changeInitiatedTimestamp) {
        /* solhint-disable not-rely-on-time */
        require(changeInitiatedTimestamp > 0, "Change not initiated");
        require(
            block.timestamp - changeInitiatedTimestamp >= GOVERNANCE_DELAY,
            "Governance delay has not elapsed"
        );
        _;
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Begins the relay request fee update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayRequestFee New relay request fee
    function beginRelayRequestFeeUpdate(uint256 _newRelayRequestFee)
        external
        onlyOwner
    {
        /* solhint-disable not-rely-on-time */
        newRelayRequestFee = _newRelayRequestFee;
        relayRequestFeeChangeInitiated = block.timestamp;
        emit RelayRequestFeeUpdateStarted(_newRelayRequestFee, block.timestamp);
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the relay request fee update process.
    /// @dev Can be called only by the contract owner, after the the
    ///      governance delay elapses.
    function finalizeRelayRequestFeeUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(relayRequestFeeChangeInitiated)
    {
        relayRequestFee = newRelayRequestFee;
        emit RelayRequestFeeUpdated(relayRequestFee);
        relayRequestFeeChangeInitiated = 0;
        newRelayRequestFee = 0;
    }

    /// @notice Begins the DKG result submission reward update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgResultSubmissionReward New DKG result submission reward
    function beginDkgResultSubmissionRewardUpdate(
        uint256 _newDkgResultSubmissionReward
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        newDkgResultSubmissionReward = _newDkgResultSubmissionReward;
        dkgResultSubmissionRewardChangeInitiated = block.timestamp;
        emit DkgResultSubmissionRewardUpdateStarted(
            _newDkgResultSubmissionReward,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG result submission reward update process.
    /// @dev Can be called only by the contract owner, after the the
    ///      governance delay elapses.
    function finalizeDkgResultSubmissionRewardUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(dkgResultSubmissionRewardChangeInitiated)
    {
        dkgResultSubmissionReward = newDkgResultSubmissionReward;
        emit DkgResultSubmissionRewardUpdated(dkgResultSubmissionReward);
        dkgResultSubmissionRewardChangeInitiated = 0;
        newDkgResultSubmissionReward = 0;
    }

    /// @notice Begins the sortition pool unlocking reward update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newSortitionPoolUnlockingReward New sortition pool unlocking reward
    function beginSortitionPoolUnlockingRewardUpdate(
        uint256 _newSortitionPoolUnlockingReward
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        newSortitionPoolUnlockingReward = _newSortitionPoolUnlockingReward;
        sortitionPoolUnlockingRewardChangeInitiated = block.timestamp;
        emit SortitionPoolUnlockingRewardUpdateStarted(
            _newSortitionPoolUnlockingReward,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the sortition pool unlocking reward update process.
    /// @dev Can be called only by the contract owner, after the the
    ///      governance delay elapses.
    function finalizeSortitionPoolUnlockingRewardUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(sortitionPoolUnlockingRewardChangeInitiated)
    {
        sortitionPoolUnlockingReward = newSortitionPoolUnlockingReward;
        emit SortitionPoolUnlockingRewardUpdated(sortitionPoolUnlockingReward);
        sortitionPoolUnlockingRewardChangeInitiated = 0;
        newSortitionPoolUnlockingReward = 0;
    }

    /// @notice Begins the relay entry submission failure slashing amount update
    ///         process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntrySubmissionFailureSlashingAmount New relay entry
    ///        submission failure slashing amount
    function beginRelayEntrySubmissionFailureSlashingAmountUpdate(
        uint256 _newRelayEntrySubmissionFailureSlashingAmount
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        newRelayEntrySubmissionFailureSlashingAmount = _newRelayEntrySubmissionFailureSlashingAmount;
        relayEntrySubmissionFailureSlashingAmountChangeInitiated = block
            .timestamp;
        emit RelayEntrySubmissionFailureSlashingAmountUpdateStarted(
            _newRelayEntrySubmissionFailureSlashingAmount,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the relay entry submission failure slashing amount
    ///         update process.
    /// @dev Can be called only by the contract owner, after the the
    ///      governance delay elapses.
    function finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            relayEntrySubmissionFailureSlashingAmountChangeInitiated
        )
    {
        relayEntrySubmissionFailureSlashingAmount = newRelayEntrySubmissionFailureSlashingAmount;
        emit RelayEntrySubmissionFailureSlashingAmountUpdated(
            relayEntrySubmissionFailureSlashingAmount
        );
        relayEntrySubmissionFailureSlashingAmountChangeInitiated = 0;
        newRelayEntrySubmissionFailureSlashingAmount = 0;
    }

    /// @notice Begins the malicious DKG result slashing amount update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newMaliciousDkgResultSlashingAmount New malicious DKG result
    ///        slashing amount
    function beginMaliciousDkgResultSlashingAmountUpdate(
        uint256 _newMaliciousDkgResultSlashingAmount
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        newMaliciousDkgResultSlashingAmount = _newMaliciousDkgResultSlashingAmount;
        maliciousDkgResultSlashingAmountChangeInitiated = block.timestamp;
        emit MaliciousDkgResultSlashingAmountUpdateStarted(
            _newMaliciousDkgResultSlashingAmount,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the malicious DKG result slashing amount update
    ///         process.
    /// @dev Can be called only by the contract owner, after the the
    ///      governance delay elapses.
    function finalizeMaliciousDkgResultSlashingAmountUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            maliciousDkgResultSlashingAmountChangeInitiated
        )
    {
        maliciousDkgResultSlashingAmount = newMaliciousDkgResultSlashingAmount;
        emit MaliciousDkgResultSlashingAmountUpdated(
            maliciousDkgResultSlashingAmount
        );
        maliciousDkgResultSlashingAmountChangeInitiated = 0;
        newMaliciousDkgResultSlashingAmount = 0;
    }

    /// @notice Begins the DKG result challenge period length update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgResultChallengePeriodLength New DKG result challenge
    ///        period length
    function beginDkgResultChallengePeriodLengthUpdate(
        uint256 _newDkgResultChallengePeriodLength
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        newDkgResultChallengePeriodLength = _newDkgResultChallengePeriodLength;
        dkgResultChallengePeriodLengthChangeInitiated = block.timestamp;
        emit DkgResultChallengePeriodLengthUpdateStarted(
            _newDkgResultChallengePeriodLength,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG result challenge period length update process.
    /// @dev Can be called only by the contract owner, after the the
    ///      governance delay elapses.
    function finalizeDkgResultChallengePeriodLengthUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(dkgResultChallengePeriodLengthChangeInitiated)
    {
        dkgResultChallengePeriodLength = newDkgResultChallengePeriodLength;
        emit DkgResultChallengePeriodLengthUpdated(
            dkgResultChallengePeriodLength
        );
        dkgResultChallengePeriodLengthChangeInitiated = 0;
        newDkgResultChallengePeriodLength = 0;
    }

    /// @notice Begins the relay entry submission eligibility delay update
    ///         process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntrySubmissionEligibilityDelay New relay entry
    ///        submission eligibility delay
    function beginRelayEntrySubmissionEligibilityDelayUpdate(
        uint256 _newRelayEntrySubmissionEligibilityDelay
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        newRelayEntrySubmissionEligibilityDelay = _newRelayEntrySubmissionEligibilityDelay;
        relayEntrySubmissionEligibilityDelayChangeInitiated = block.timestamp;
        emit RelayEntrySubmissionEligibilityDelayUpdateStarted(
            _newRelayEntrySubmissionEligibilityDelay,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the relay entry submission eligibility delay update
    ////        process.
    /// @dev Can be called only by the contract owner, after the the
    ///      governance delay elapses.
    function finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            relayEntrySubmissionEligibilityDelayChangeInitiated
        )
    {
        relayEntrySubmissionEligibilityDelay = newRelayEntrySubmissionEligibilityDelay;
        emit RelayEntrySubmissionEligibilityDelayUpdated(
            relayEntrySubmissionEligibilityDelay
        );
        relayEntrySubmissionEligibilityDelayChangeInitiated = 0;
        newRelayEntrySubmissionEligibilityDelay = 0;
    }

    /// @notice Begins the relay entry hard timeout update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntryHardTimeout New relay entry hard timeout
    function beginRelayEntryHardTimeoutUpdate(uint256 _newRelayEntryHardTimeout)
        external
        onlyOwner
    {
        /* solhint-disable not-rely-on-time */
        newRelayEntryHardTimeout = _newRelayEntryHardTimeout;
        relayEntryHardTimeoutChangeInitiated = block.timestamp;
        emit RelayEntryHardTimeoutUpdateStarted(
            _newRelayEntryHardTimeout,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the relay entry hard timeout update process.
    /// @dev Can be called only by the contract owner, after the the
    ///      governance delay elapses.
    function finalizeRelayEntryHardTimeoutUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(relayEntryHardTimeoutChangeInitiated)
    {
        relayEntryHardTimeout = newRelayEntryHardTimeout;
        emit RelayEntryHardTimeoutUpdated(relayEntryHardTimeout);
        relayEntryHardTimeoutChangeInitiated = 0;
        newRelayEntryHardTimeout = 0;
    }

    /// @notice Begins the group creation frequency update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newGroupCreationFrequency New group creation frequency
    function beginGroupCreationFrequencyUpdate(
        uint256 _newGroupCreationFrequency
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        newGroupCreationFrequency = _newGroupCreationFrequency;
        groupCreationFrequencyChangeInitiated = block.timestamp;
        emit GroupCreationFrequencyUpdateStarted(
            _newGroupCreationFrequency,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the group creation frequency update process.
    /// @dev Can be called only by the contract owner, after the the
    ///      governance delay elapses.
    function finalizeGroupCreationFrequencyUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(groupCreationFrequencyChangeInitiated)
    {
        groupCreationFrequency = newGroupCreationFrequency;
        emit GroupCreationFrequencyUpdated(groupCreationFrequency);
        groupCreationFrequencyChangeInitiated = 0;
        newGroupCreationFrequency = 0;
    }

    /// @notice Begins the group lifetime update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newGroupLifetime New group lifetime
    function beginGroupLifetimeUpdate(uint256 _newGroupLifetime)
        external
        onlyOwner
    {
        /* solhint-disable not-rely-on-time */
        newGroupLifetime = _newGroupLifetime;
        groupLifetimeChangeInitiated = block.timestamp;
        emit GroupLifetimeUpdateStarted(_newGroupLifetime, block.timestamp);
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the group creation frequency update process.
    /// @dev Can be called only by the contract owner, after the the
    ///      governance delay elapses.
    function finalizeGroupLifetimeUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(groupLifetimeChangeInitiated)
    {
        groupLifetime = newGroupLifetime;
        emit GroupLifetimeUpdated(groupLifetime);
        groupLifetimeChangeInitiated = 0;
        newGroupLifetime = 0;
    }

    /// @notice Begins the callback gas limit update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newCallbackGasLimit New callback gas limit
    function beginCallbackGasLimitUpdate(uint256 _newCallbackGasLimit)
        external
        onlyOwner
    {
        /* solhint-disable not-rely-on-time */
        newCallbackGasLimit = _newCallbackGasLimit;
        callbackGasLimitChangeInitiated = block.timestamp;
        emit CallbackGasLimitUpdateStarted(
            _newCallbackGasLimit,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the callback gas limit update process.
    /// @dev Can be called only by the contract owner, after the the
    ///      governance delay elapses.
    function finalizeCallbackGasLimitUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(callbackGasLimitChangeInitiated)
    {
        callbackGasLimit = newCallbackGasLimit;
        emit CallbackGasLimitUpdated(callbackGasLimit);
        callbackGasLimitChangeInitiated = 0;
        newCallbackGasLimit = 0;
    }

    /// @notice Get the time remaining until the relay request fee can be
    ///         updated.
    /// @return Remaining time in seconds.
    function getRemainingRelayRequestFeeUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            GovernanceUtils.getRemainingChangeTime(
                relayRequestFeeChangeInitiated,
                GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the DKG result submission reward
    ///         can be updated.
    /// @return Remaining time in seconds.
    function getRemainingDkgResultSubmissionRewardUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            GovernanceUtils.getRemainingChangeTime(
                dkgResultSubmissionRewardChangeInitiated,
                GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the sortition pool unlocking reward
    ///         can be updated.
    /// @return Remaining time in seconds.
    function getRemainingSortitionPoolUnlockingRewardUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            GovernanceUtils.getRemainingChangeTime(
                sortitionPoolUnlockingRewardChangeInitiated,
                GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the relay entry submission failure
    ///         slashing amount can be updated.
    /// @return Remaining time in seconds.
    function getRemainingRelayEntrySubmissionFailureSlashingAmountUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            GovernanceUtils.getRemainingChangeTime(
                relayEntrySubmissionFailureSlashingAmountChangeInitiated,
                GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the malicious DKG result
    ///         slashing amount can be updated.
    /// @return Remaining time in seconds.
    function getRemainingMaliciousDkgResultSlashingAmountUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            GovernanceUtils.getRemainingChangeTime(
                maliciousDkgResultSlashingAmountChangeInitiated,
                GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the DKG result challenge period
    ///         length can be updated.
    /// @return Remaining time in seconds.
    function getRemainingDkgResultChallengePeriodLengthUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            GovernanceUtils.getRemainingChangeTime(
                dkgResultChallengePeriodLengthChangeInitiated,
                GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the relay entry submission
    ///         eligibility delay can be updated.
    /// @return Remaining time in seconds.
    function getRemainingRelayEntrySubmissionEligibilityDelayUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            GovernanceUtils.getRemainingChangeTime(
                relayEntrySubmissionEligibilityDelayChangeInitiated,
                GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the relay entry hard timeout can be
    ///         updated.
    /// @return Remaining time in seconds.
    function getRemainingRelayEntryHardTimeoutUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            GovernanceUtils.getRemainingChangeTime(
                relayEntryHardTimeoutChangeInitiated,
                GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the group creation frequency can be
    ///         updated.
    /// @return Remaining time in seconds.
    function getRemainingGroupCreationFrequencyUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            GovernanceUtils.getRemainingChangeTime(
                groupCreationFrequencyChangeInitiated,
                GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the group lifetime can be updated.
    /// @return Remaining time in seconds.
    function getRemainingGroupLifetimeUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            GovernanceUtils.getRemainingChangeTime(
                groupLifetimeChangeInitiated,
                GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the callback gas limit can be
    ///         updated.
    /// @return Remaining time in seconds.
    function getRemainingCallbackGasLimitUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            GovernanceUtils.getRemainingChangeTime(
                callbackGasLimitChangeInitiated,
                GOVERNANCE_DELAY
            );
    }
}
