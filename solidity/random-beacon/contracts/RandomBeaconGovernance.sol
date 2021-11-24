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

import "./RandomBeacon.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @title Keep Random Beacon Governance
/// @notice Owns the `RandomBeacon` contract and is responsible for updating its
///         governable parameters in respect to governance delay individual
///         for each parameter.
contract RandomBeaconGovernance is Ownable {
    uint256 public newRelayRequestFee;
    uint256 public relayRequestFeeChangeInitiated;

    uint256 public newRelayEntrySubmissionEligibilityDelay;
    uint256 public relayEntrySubmissionEligibilityDelayChangeInitiated;

    uint256 public newRelayEntryHardTimeout;
    uint256 public relayEntryHardTimeoutChangeInitiated;

    uint256 public newCallbackGasLimit;
    uint256 public callbackGasLimitChangeInitiated;

    uint256 public newGroupCreationFrequency;
    uint256 public groupCreationFrequencyChangeInitiated;

    uint256 public newGroupLifetime;
    uint256 public groupLifetimeChangeInitiated;

    uint256 public newDkgResultChallengePeriodLength;
    uint256 public dkgResultChallengePeriodLengthChangeInitiated;

    uint256 public newDkgResultSubmissionEligibilityDelay;
    uint256 public dkgResultSubmissionEligibilityDelayChangeInitiated;

    uint256 public newDkgResultSubmissionReward;
    uint256 public dkgResultSubmissionRewardChangeInitiated;

    uint256 public newSortitionPoolUnlockingReward;
    uint256 public sortitionPoolUnlockingRewardChangeInitiated;

    uint256 public newRelayEntrySubmissionFailureSlashingAmount;
    uint256 public relayEntrySubmissionFailureSlashingAmountChangeInitiated;

    uint256 public newMaliciousDkgResultSlashingAmount;
    uint256 public maliciousDkgResultSlashingAmountChangeInitiated;

    uint256 public newSortitionPoolRewardsBanDuration;
    uint256 public sortitionPoolRewardsBanDurationChangeInitiated;

    uint256 public newRelayEntryTimeoutNotificationRewardMultiplier;
    uint256 public relayEntryTimeoutNotificationRewardMultiplierChangeInitiated;

    uint96 public newMinimumAuthorization;
    uint256 public minimumAuthorizationChangeInitiated;

    uint64 public newAuthorizationDecreaseDelay;
    uint256 public authorizationDecreaseDelayChangeInitiated;

    uint256 public newDkgMaliciousResultNotificationRewardMultiplier;
    uint256
        public dkgMaliciousResultNotificationRewardMultiplierChangeInitiated;

    RandomBeacon public randomBeacon;

    // Long governance delay used for critical parameters giving a chance for
    // stakers to opt out before the change is finalized in case they do not
    // agree with that change. The maximum group lifetime must not be longer
    // than this delay.
    //
    // The full list of parameters protected by this delay:
    // - relay entry hard timeout
    // - callback gas limit
    // - group lifetime
    // - relay entry submission failure slashing amount
    // - minimum authorization
    // - authorization decrease delay
    uint256 internal constant CRITICAL_PARAMETER_GOVERNANCE_DELAY = 2 weeks;

    // Short governance delay for non-critical parameters. Honest stakers should
    // not be severely affected by any change of these parameters.
    //
    // The full list of parameters protected by this delay:
    // - relay request fee
    // - group creation frequency
    // - relay entry submission eligibility delay
    // - DKG result challenge period length
    // - DKG result submission eligibility delay
    // - DKG result submission reward
    // - sortition pool rewards ban duration
    // - malicious DKG result slashing amount
    // - sortition pool unlocking reward
    // - relay entry timeout notification reward multiplier
    // - DKG malicious result notification reward multiplier
    uint256 internal constant STANDARD_PARAMETER_GOVERNANCE_DELAY = 12 hours;

    event RelayRequestFeeUpdateStarted(
        uint256 relayRequestFee,
        uint256 timestamp
    );
    event RelayRequestFeeUpdated(uint256 relayRequestFee);

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

    event CallbackGasLimitUpdateStarted(
        uint256 callbackGasLimit,
        uint256 timestamp
    );
    event CallbackGasLimitUpdated(uint256 callbackGasLimit);

    event GroupCreationFrequencyUpdateStarted(
        uint256 groupCreationFrequency,
        uint256 timestamp
    );
    event GroupCreationFrequencyUpdated(uint256 groupCreationFrequency);

    event GroupLifetimeUpdateStarted(uint256 groupLifetime, uint256 timestamp);
    event GroupLifetimeUpdated(uint256 groupLifetime);

    event DkgResultChallengePeriodLengthUpdateStarted(
        uint256 dkgResultChallengePeriodLength,
        uint256 timestamp
    );
    event DkgResultChallengePeriodLengthUpdated(
        uint256 dkgResultChallengePeriodLength
    );

    event DkgResultSubmissionEligibilityDelayUpdateStarted(
        uint256 dkgResultSubmissionEligibilityDelay,
        uint256 timestamp
    );
    event DkgResultSubmissionEligibilityDelayUpdated(
        uint256 dkgResultSubmissionEligibilityDelay
    );

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

    event SortitionPoolRewardsBanDurationUpdateStarted(
        uint256 sortitionPoolRewardsBanDuration,
        uint256 timestamp
    );
    event SortitionPoolRewardsBanDurationUpdated(
        uint256 sortitionPoolRewardsBanDuration
    );

    event RelayEntryTimeoutNotificationRewardMultiplierUpdateStarted(
        uint256 relayEntryTimeoutNotificationRewardMultiplier,
        uint256 timestamp
    );
    event RelayEntryTimeoutNotificationRewardMultiplierUpdated(
        uint256 relayEntryTimeoutNotificationRewardMultiplier
    );

    event MinimumAuthorizationUpdateStarted(
        uint96 minimumAuthorization,
        uint256 timestamp
    );
    event MinimumAuthorizationUpdated(uint96 minimumAuthorization);

    event AuthorizationDecreaseDelayUpdateStarted(
        uint64 authorizationDecreaseDelay,
        uint256 timestamp
    );
    event AuthorizationDecreaseDelayUpdated(uint64 authorizationDecreaseDelay);

    event DkgMaliciousResultNotificationRewardMultiplierUpdateStarted(
        uint256 dkgMaliciousResultNotificationRewardMultiplier,
        uint256 timestamp
    );
    event DkgMaliciousResultNotificationRewardMultiplierUpdated(
        uint256 dkgMaliciousResultNotificationRewardMultiplier
    );

    /// @notice Reverts if called before the governance delay elapses.
    /// @param changeInitiatedTimestamp Timestamp indicating the beginning
    ///        of the change.
    modifier onlyAfterGovernanceDelay(
        uint256 changeInitiatedTimestamp,
        uint256 delay
    ) {
        /* solhint-disable not-rely-on-time */
        require(changeInitiatedTimestamp > 0, "Change not initiated");
        require(
            block.timestamp - changeInitiatedTimestamp >= delay,
            "Governance delay has not elapsed"
        );
        _;
        /* solhint-enable not-rely-on-time */
    }

    constructor(RandomBeacon _randomBeacon) {
        randomBeacon = _randomBeacon;
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
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayRequestFeeUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            relayRequestFeeChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit RelayRequestFeeUpdated(newRelayRequestFee);
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateRelayEntryParameters(
            newRelayRequestFee,
            randomBeacon.relayEntrySubmissionEligibilityDelay(),
            randomBeacon.relayEntryHardTimeout(),
            randomBeacon.callbackGasLimit()
        );
        relayRequestFeeChangeInitiated = 0;
        newRelayRequestFee = 0;
    }

    /// @notice Begins the relay entry submission eligibility delay update
    ///         process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntrySubmissionEligibilityDelay New relay entry
    ///        submission eligibility delay in blocks
    function beginRelayEntrySubmissionEligibilityDelayUpdate(
        uint256 _newRelayEntrySubmissionEligibilityDelay
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        require(
            _newRelayEntrySubmissionEligibilityDelay > 0,
            "Relay entry submission eligibility delay must be > 0"
        );
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
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            relayEntrySubmissionEligibilityDelayChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit RelayEntrySubmissionEligibilityDelayUpdated(
            newRelayEntrySubmissionEligibilityDelay
        );
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateRelayEntryParameters(
            randomBeacon.relayRequestFee(),
            newRelayEntrySubmissionEligibilityDelay,
            randomBeacon.relayEntryHardTimeout(),
            randomBeacon.callbackGasLimit()
        );
        relayEntrySubmissionEligibilityDelayChangeInitiated = 0;
        newRelayEntrySubmissionEligibilityDelay = 0;
    }

    /// @notice Begins the relay entry hard timeout update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntryHardTimeout New relay entry hard timeout in blocks
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
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayEntryHardTimeoutUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            relayEntryHardTimeoutChangeInitiated,
            CRITICAL_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit RelayEntryHardTimeoutUpdated(newRelayEntryHardTimeout);
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateRelayEntryParameters(
            randomBeacon.relayRequestFee(),
            randomBeacon.relayEntrySubmissionEligibilityDelay(),
            newRelayEntryHardTimeout,
            randomBeacon.callbackGasLimit()
        );
        relayEntryHardTimeoutChangeInitiated = 0;
        newRelayEntryHardTimeout = 0;
    }

    /// @notice Begins the callback gas limit update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newCallbackGasLimit New callback gas limit
    function beginCallbackGasLimitUpdate(uint256 _newCallbackGasLimit)
        external
        onlyOwner
    {
        /* solhint-disable not-rely-on-time */
        // slither-disable-next-line too-many-digits
        require(
            _newCallbackGasLimit > 0 && _newCallbackGasLimit <= 1e6,
            "Callback gas limit must be > 0 and <= 1000000"
        );
        newCallbackGasLimit = _newCallbackGasLimit;
        callbackGasLimitChangeInitiated = block.timestamp;
        emit CallbackGasLimitUpdateStarted(
            _newCallbackGasLimit,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the callback gas limit update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeCallbackGasLimitUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            callbackGasLimitChangeInitiated,
            CRITICAL_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit CallbackGasLimitUpdated(newCallbackGasLimit);
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateRelayEntryParameters(
            randomBeacon.relayRequestFee(),
            randomBeacon.relayEntrySubmissionEligibilityDelay(),
            randomBeacon.relayEntryHardTimeout(),
            newCallbackGasLimit
        );
        callbackGasLimitChangeInitiated = 0;
        newCallbackGasLimit = 0;
    }

    /// @notice Begins the group creation frequency update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newGroupCreationFrequency New group creation frequency
    function beginGroupCreationFrequencyUpdate(
        uint256 _newGroupCreationFrequency
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        require(
            _newGroupCreationFrequency > 0,
            "Group creation frequency must be > 0"
        );
        newGroupCreationFrequency = _newGroupCreationFrequency;
        groupCreationFrequencyChangeInitiated = block.timestamp;
        emit GroupCreationFrequencyUpdateStarted(
            _newGroupCreationFrequency,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the group creation frequency update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeGroupCreationFrequencyUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            groupCreationFrequencyChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit GroupCreationFrequencyUpdated(newGroupCreationFrequency);
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateGroupCreationParameters(
            newGroupCreationFrequency,
            randomBeacon.groupLifetime()
        );
        groupCreationFrequencyChangeInitiated = 0;
        newGroupCreationFrequency = 0;
    }

    /// @notice Begins the group lifetime update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newGroupLifetime New group lifetime in blocks
    function beginGroupLifetimeUpdate(uint256 _newGroupLifetime)
        external
        onlyOwner
    {
        /* solhint-disable not-rely-on-time */
        require(
            _newGroupLifetime >= 1 days && _newGroupLifetime <= 2 weeks,
            "Group lifetime must be >= 1 day and <= 2 weeks"
        );
        newGroupLifetime = _newGroupLifetime;
        groupLifetimeChangeInitiated = block.timestamp;
        emit GroupLifetimeUpdateStarted(_newGroupLifetime, block.timestamp);
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the group creation frequency update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeGroupLifetimeUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            groupLifetimeChangeInitiated,
            CRITICAL_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit GroupLifetimeUpdated(newGroupLifetime);
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateGroupCreationParameters(
            randomBeacon.groupCreationFrequency(),
            newGroupLifetime
        );
        groupLifetimeChangeInitiated = 0;
        newGroupLifetime = 0;
    }

    /// @notice Begins the DKG result challenge period length update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgResultChallengePeriodLength New DKG result challenge
    ///        period length in blocks
    function beginDkgResultChallengePeriodLengthUpdate(
        uint256 _newDkgResultChallengePeriodLength
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        require(
            _newDkgResultChallengePeriodLength >= 10,
            "DKG result challenge period length must be >= 10"
        );
        newDkgResultChallengePeriodLength = _newDkgResultChallengePeriodLength;
        dkgResultChallengePeriodLengthChangeInitiated = block.timestamp;
        emit DkgResultChallengePeriodLengthUpdateStarted(
            _newDkgResultChallengePeriodLength,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG result challenge period length update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgResultChallengePeriodLengthUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            dkgResultChallengePeriodLengthChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit DkgResultChallengePeriodLengthUpdated(
            newDkgResultChallengePeriodLength
        );
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateDkgParameters(
            newDkgResultChallengePeriodLength,
            randomBeacon.dkgResultSubmissionEligibilityDelay()
        );
        dkgResultChallengePeriodLengthChangeInitiated = 0;
        newDkgResultChallengePeriodLength = 0;
    }

    /// @notice Begins the DKG result submission eligibility delay update
    ///         process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgResultSubmissionEligibilityDelay New DKG result submission
    ///        eligibility delay in blocks
    function beginDkgResultSubmissionEligibilityDelayUpdate(
        uint256 _newDkgResultSubmissionEligibilityDelay
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        require(
            _newDkgResultSubmissionEligibilityDelay > 0,
            "DKG result submission eligibility delay must be > 0"
        );
        newDkgResultSubmissionEligibilityDelay = _newDkgResultSubmissionEligibilityDelay;
        dkgResultSubmissionEligibilityDelayChangeInitiated = block.timestamp;
        emit DkgResultSubmissionEligibilityDelayUpdateStarted(
            _newDkgResultSubmissionEligibilityDelay,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG result submission eligibility delay update
    ///         process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgResultSubmissionEligibilityDelayUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            dkgResultSubmissionEligibilityDelayChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit DkgResultSubmissionEligibilityDelayUpdated(
            newDkgResultSubmissionEligibilityDelay
        );
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateDkgParameters(
            randomBeacon.dkgResultChallengePeriodLength(),
            newDkgResultSubmissionEligibilityDelay
        );
        dkgResultSubmissionEligibilityDelayChangeInitiated = 0;
        newDkgResultSubmissionEligibilityDelay = 0;
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
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgResultSubmissionRewardUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            dkgResultSubmissionRewardChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit DkgResultSubmissionRewardUpdated(newDkgResultSubmissionReward);
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateRewardParameters(
            newDkgResultSubmissionReward,
            randomBeacon.sortitionPoolUnlockingReward(),
            randomBeacon.sortitionPoolRewardsBanDuration(),
            randomBeacon.relayEntryTimeoutNotificationRewardMultiplier(),
            randomBeacon.dkgMaliciousResultNotificationRewardMultiplier()
        );
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
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeSortitionPoolUnlockingRewardUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            sortitionPoolUnlockingRewardChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit SortitionPoolUnlockingRewardUpdated(
            newSortitionPoolUnlockingReward
        );
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateRewardParameters(
            randomBeacon.dkgResultSubmissionReward(),
            newSortitionPoolUnlockingReward,
            randomBeacon.sortitionPoolRewardsBanDuration(),
            randomBeacon.relayEntryTimeoutNotificationRewardMultiplier(),
            randomBeacon.dkgMaliciousResultNotificationRewardMultiplier()
        );
        sortitionPoolUnlockingRewardChangeInitiated = 0;
        newSortitionPoolUnlockingReward = 0;
    }

    /// @notice Begins the sortition pool rewards ban duration update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newSortitionPoolRewardsBanDuration New sortition pool rewards
    ///        ban duration.
    function beginSortitionPoolRewardsBanDurationUpdate(
        uint256 _newSortitionPoolRewardsBanDuration
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        newSortitionPoolRewardsBanDuration = _newSortitionPoolRewardsBanDuration;
        sortitionPoolRewardsBanDurationChangeInitiated = block.timestamp;
        emit SortitionPoolRewardsBanDurationUpdateStarted(
            _newSortitionPoolRewardsBanDuration,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the sortition pool rewards ban duration update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeSortitionPoolRewardsBanDurationUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            sortitionPoolRewardsBanDurationChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit SortitionPoolRewardsBanDurationUpdated(
            newSortitionPoolRewardsBanDuration
        );
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateRewardParameters(
            randomBeacon.dkgResultSubmissionReward(),
            randomBeacon.sortitionPoolUnlockingReward(),
            newSortitionPoolRewardsBanDuration,
            randomBeacon.relayEntryTimeoutNotificationRewardMultiplier(),
            randomBeacon.dkgMaliciousResultNotificationRewardMultiplier()
        );
        sortitionPoolRewardsBanDurationChangeInitiated = 0;
        newSortitionPoolRewardsBanDuration = 0;
    }

    /// @notice Begins the relay entry timeout notification reward multiplier
    ///         update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntryTimeoutNotificationRewardMultiplier New relay
    ///        entry timeout notification reward multiplier.
    function beginRelayEntryTimeoutNotificationRewardMultiplierUpdate(
        uint256 _newRelayEntryTimeoutNotificationRewardMultiplier
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        require(
            _newRelayEntryTimeoutNotificationRewardMultiplier <= 100,
            "Maximum value is 100"
        );

        newRelayEntryTimeoutNotificationRewardMultiplier = _newRelayEntryTimeoutNotificationRewardMultiplier;
        relayEntryTimeoutNotificationRewardMultiplierChangeInitiated = block
            .timestamp;
        emit RelayEntryTimeoutNotificationRewardMultiplierUpdateStarted(
            _newRelayEntryTimeoutNotificationRewardMultiplier,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the relay entry timeout notification reward
    ///         multiplier update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayEntryTimeoutNotificationRewardMultiplierUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            relayEntryTimeoutNotificationRewardMultiplierChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit RelayEntryTimeoutNotificationRewardMultiplierUpdated(
            newRelayEntryTimeoutNotificationRewardMultiplier
        );
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateRewardParameters(
            randomBeacon.dkgResultSubmissionReward(),
            randomBeacon.sortitionPoolUnlockingReward(),
            randomBeacon.sortitionPoolRewardsBanDuration(),
            newRelayEntryTimeoutNotificationRewardMultiplier,
            randomBeacon.dkgMaliciousResultNotificationRewardMultiplier()
        );
        relayEntryTimeoutNotificationRewardMultiplierChangeInitiated = 0;
        newRelayEntryTimeoutNotificationRewardMultiplier = 0;
    }

    /// @notice Begins the DKG malicious result notification reward multiplier
    ///         update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgMaliciousResultNotificationRewardMultiplier New DKG
    ///        malicious result notification reward multiplier.
    function beginDkgMaliciousResultNotificationRewardMultiplierUpdate(
        uint256 _newDkgMaliciousResultNotificationRewardMultiplier
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        require(
            _newDkgMaliciousResultNotificationRewardMultiplier <= 100,
            "Maximum value is 100"
        );

        newDkgMaliciousResultNotificationRewardMultiplier = _newDkgMaliciousResultNotificationRewardMultiplier;
        dkgMaliciousResultNotificationRewardMultiplierChangeInitiated = block
            .timestamp;
        emit DkgMaliciousResultNotificationRewardMultiplierUpdateStarted(
            _newDkgMaliciousResultNotificationRewardMultiplier,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG malicious result notification reward
    ///         multiplier update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgMaliciousResultNotificationRewardMultiplierUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            dkgMaliciousResultNotificationRewardMultiplierChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit DkgMaliciousResultNotificationRewardMultiplierUpdated(
            newDkgMaliciousResultNotificationRewardMultiplier
        );
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateRewardParameters(
            randomBeacon.dkgResultSubmissionReward(),
            randomBeacon.sortitionPoolUnlockingReward(),
            randomBeacon.sortitionPoolRewardsBanDuration(),
            randomBeacon.relayEntryTimeoutNotificationRewardMultiplier(),
            newDkgMaliciousResultNotificationRewardMultiplier
        );
        dkgMaliciousResultNotificationRewardMultiplierChangeInitiated = 0;
        newDkgMaliciousResultNotificationRewardMultiplier = 0;
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
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            relayEntrySubmissionFailureSlashingAmountChangeInitiated,
            CRITICAL_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit RelayEntrySubmissionFailureSlashingAmountUpdated(
            newRelayEntrySubmissionFailureSlashingAmount
        );
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateSlashingParameters(
            newRelayEntrySubmissionFailureSlashingAmount,
            randomBeacon.maliciousDkgResultSlashingAmount()
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
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeMaliciousDkgResultSlashingAmountUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            maliciousDkgResultSlashingAmountChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit MaliciousDkgResultSlashingAmountUpdated(
            newMaliciousDkgResultSlashingAmount
        );
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateSlashingParameters(
            randomBeacon.relayEntrySubmissionFailureSlashingAmount(),
            newMaliciousDkgResultSlashingAmount
        );
        maliciousDkgResultSlashingAmountChangeInitiated = 0;
        newMaliciousDkgResultSlashingAmount = 0;
    }

    /// @notice Begins the minimum authorization amount update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newMinimumAuthorization New minimum authorization amount.
    function beginMinimumAuthorizationUpdate(uint96 _newMinimumAuthorization)
        external
        onlyOwner
    {
        /* solhint-disable not-rely-on-time */
        newMinimumAuthorization = _newMinimumAuthorization;
        minimumAuthorizationChangeInitiated = block.timestamp;
        emit MinimumAuthorizationUpdateStarted(
            _newMinimumAuthorization,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the minimum authorization amount update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeMinimumAuthorizationUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            minimumAuthorizationChangeInitiated,
            CRITICAL_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit MinimumAuthorizationUpdated(newMinimumAuthorization);
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateAuthorizationParameters(
            newMinimumAuthorization,
            randomBeacon.authorizationDecreaseDelay()
        );
        minimumAuthorizationChangeInitiated = 0;
        newMinimumAuthorization = 0;
    }

    /// @notice Begins the authorization decrease delay update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newAuthorizationDecreaseDelay New authorization decrease delay
    function beginAuthorizationDecreaseDelayUpdate(
        uint64 _newAuthorizationDecreaseDelay
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        newAuthorizationDecreaseDelay = _newAuthorizationDecreaseDelay;
        authorizationDecreaseDelayChangeInitiated = block.timestamp;
        emit AuthorizationDecreaseDelayUpdateStarted(
            _newAuthorizationDecreaseDelay,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the authorization decrease delay update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeAuthorizationDecreaseDelayUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            authorizationDecreaseDelayChangeInitiated,
            CRITICAL_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit AuthorizationDecreaseDelayUpdated(newAuthorizationDecreaseDelay);
        // slither-disable-next-line reentrancy-no-eth
        randomBeacon.updateAuthorizationParameters(
            randomBeacon.minimumAuthorization(),
            newAuthorizationDecreaseDelay
        );
        authorizationDecreaseDelayChangeInitiated = 0;
        newAuthorizationDecreaseDelay = 0;
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
            getRemainingChangeTime(
                relayRequestFeeChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
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
            getRemainingChangeTime(
                relayEntrySubmissionEligibilityDelayChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
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
            getRemainingChangeTime(
                relayEntryHardTimeoutChangeInitiated,
                CRITICAL_PARAMETER_GOVERNANCE_DELAY
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
            getRemainingChangeTime(
                callbackGasLimitChangeInitiated,
                CRITICAL_PARAMETER_GOVERNANCE_DELAY
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
            getRemainingChangeTime(
                groupCreationFrequencyChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
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
            getRemainingChangeTime(
                groupLifetimeChangeInitiated,
                CRITICAL_PARAMETER_GOVERNANCE_DELAY
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
            getRemainingChangeTime(
                dkgResultChallengePeriodLengthChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the DKG result submission
    ///         eligibility delay can be updated.
    /// @return Remaining time in seconds.
    function getRemainingDkgResultSubmissionEligibilityDelayUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                dkgResultSubmissionEligibilityDelayChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
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
            getRemainingChangeTime(
                dkgResultSubmissionRewardChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
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
            getRemainingChangeTime(
                sortitionPoolUnlockingRewardChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
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
            getRemainingChangeTime(
                relayEntrySubmissionFailureSlashingAmountChangeInitiated,
                CRITICAL_PARAMETER_GOVERNANCE_DELAY
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
            getRemainingChangeTime(
                maliciousDkgResultSlashingAmountChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the minimum authorization amount
    ///         can be updated.
    /// @return Remaining time in seconds.
    function getRemainingMimimumAuthorizationUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                minimumAuthorizationChangeInitiated,
                CRITICAL_PARAMETER_GOVERNANCE_DELAY
            );
    }

    function getRemainingAuthorizationDecreaseDelayUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                authorizationDecreaseDelayChangeInitiated,
                CRITICAL_PARAMETER_GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the sortition pool rewards ban
    ///         duration can be updated.
    /// @return Remaining time in seconds.
    function getRemainingSortitionPoolRewardsBanDurationUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                sortitionPoolRewardsBanDurationChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the relay entry timeout
    ///         notification reward multiplier duration can be updated.
    /// @return Remaining time in seconds.
    function getRemainingRelayEntryTimeoutNotificationRewardMultiplierUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                relayEntryTimeoutNotificationRewardMultiplierChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the DKG malicious result
    ///         notification reward multiplier duration can be updated.
    /// @return Remaining time in seconds.
    function getRemainingDkgMaliciousResultNotificationRewardMultiplierUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                dkgMaliciousResultNotificationRewardMultiplierChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
            );
    }

    /// @notice Gets the time remaining until the governable parameter update
    ///         can be committed.
    /// @param changeTimestamp Timestamp indicating the beginning of the change.
    /// @param delay Governance delay.
    /// @return Remaining time in seconds.
    function getRemainingChangeTime(uint256 changeTimestamp, uint256 delay)
        internal
        view
        returns (uint256)
    {
        require(changeTimestamp > 0, "Change not initiated");
        /* solhint-disable-next-line not-rely-on-time */
        uint256 elapsed = block.timestamp - changeTimestamp;
        if (elapsed >= delay) {
            return 0;
        } else {
            return delay - elapsed;
        }
    }
}
