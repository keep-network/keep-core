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
//

pragma solidity ^0.8.9;

/// @title Governance Beacon Params
/// @notice Library is used by the `RandomBeaconGovernance` contract and is
///         responsible for storing and updating governence parameters.
library GovernanceBeaconParams {
    struct Data {
        uint256 governanceDelay;
        uint256 newDkgResultSubmissionReward;
        uint256 dkgResultSubmissionRewardChangeInitiated;
        uint256 newSortitionPoolUnlockingReward;
        uint256 sortitionPoolUnlockingRewardChangeInitiated;
        uint256 newIneligibleOperatorNotifierReward;
        uint256 ineligibleOperatorNotifierRewardChangeInitiated;
        uint256 newSortitionPoolRewardsBanDuration;
        uint256 sortitionPoolRewardsBanDurationChangeInitiated;
        uint256 newUnauthorizedSigningNotificationRewardMultiplier;
        uint256 unauthorizedSigningNotificationRewardMultiplierChangeInitiated;
        uint256 newRelayEntryTimeoutNotificationRewardMultiplier;
        uint256 relayEntryTimeoutNotificationRewardMultiplierChangeInitiated;
        uint256 newDkgMaliciousResultNotificationRewardMultiplier;
        uint256 dkgMaliciousResultNotificationRewardMultiplierChangeInitiated;
        uint96 newRelayEntrySubmissionFailureSlashingAmount;
        uint256 relayEntrySubmissionFailureSlashingAmountChangeInitiated;
        uint256 newMaliciousDkgResultSlashingAmount;
        uint256 maliciousDkgResultSlashingAmountChangeInitiated;
        uint256 newUnauthorizedSigningSlashingAmount;
        uint256 unauthorizedSigningSlashingAmountChangeInitiated;
        uint96 newMinimumAuthorization;
        uint256 minimumAuthorizationChangeInitiated;
        uint64 newAuthorizationDecreaseDelay;
        uint256 authorizationDecreaseDelayChangeInitiated;
        uint256 newRelayRequestFee;
        uint256 relayRequestFeeChangeInitiated;
        uint256 newRelayEntrySoftTimeout;
        uint256 relayEntrySoftTimeoutChangeInitiated;
        uint256 newRelayEntryHardTimeout;
        uint256 relayEntryHardTimeoutChangeInitiated;
        uint256 newCallbackGasLimit;
        uint256 callbackGasLimitChangeInitiated;
        uint256 newGroupCreationFrequency;
        uint256 groupCreationFrequencyChangeInitiated;
        uint256 newGroupLifetime;
        uint256 groupLifetimeChangeInitiated;
        uint256 newDkgResultChallengePeriodLength;
        uint256 dkgResultChallengePeriodLengthChangeInitiated;
        uint256 newDkgResultSubmissionTimeout;
        uint256 dkgResultSubmissionTimeoutChangeInitiated;
        uint256 newSubmitterPrecedencePeriodLength;
        uint256 dkgSubmitterPrecedencePeriodLengthChangeInitiated;
        uint256 newGovernanceDelay;
        uint256 governanceDelayChangeInitiated;
        address newRandomBeaconOwner;
        uint256 randomBeaconOwnershipTransferInitiated;
    }

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

    event IneligibleOperatorNotifierRewardUpdateStarted(
        uint256 ineligibleOperatorNotifierReward,
        uint256 timestamp
    );
    event IneligibleOperatorNotifierRewardUpdated(
        uint256 ineligibleOperatorNotifierReward
    );

    event SortitionPoolRewardsBanDurationUpdateStarted(
        uint256 sortitionPoolRewardsBanDuration,
        uint256 timestamp
    );
    event SortitionPoolRewardsBanDurationUpdated(
        uint256 sortitionPoolRewardsBanDuration
    );

    event UnauthorizedSigningNotificationRewardMultiplierUpdateStarted(
        uint256 unauthorizedSigningTimeoutNotificationRewardMultiplier,
        uint256 timestamp
    );
    event UnauthorizedSigningNotificationRewardMultiplierUpdated(
        uint256 unauthorizedSigningTimeoutNotificationRewardMultiplier
    );

    event RelayEntryTimeoutNotificationRewardMultiplierUpdateStarted(
        uint256 relayEntryTimeoutNotificationRewardMultiplier,
        uint256 timestamp
    );
    event RelayEntryTimeoutNotificationRewardMultiplierUpdated(
        uint256 relayEntryTimeoutNotificationRewardMultiplier
    );

    event DkgMaliciousResultNotificationRewardMultiplierUpdateStarted(
        uint256 dkgMaliciousResultNotificationRewardMultiplier,
        uint256 timestamp
    );
    event DkgMaliciousResultNotificationRewardMultiplierUpdated(
        uint256 dkgMaliciousResultNotificationRewardMultiplier
    );

    event RelayEntrySubmissionFailureSlashingAmountUpdateStarted(
        uint96 relayEntrySubmissionFailureSlashingAmount,
        uint256 timestamp
    );
    event RelayEntrySubmissionFailureSlashingAmountUpdated(
        uint96 relayEntrySubmissionFailureSlashingAmount
    );

    event MaliciousDkgResultSlashingAmountUpdateStarted(
        uint256 maliciousDkgResultSlashingAmount,
        uint256 timestamp
    );
    event MaliciousDkgResultSlashingAmountUpdated(
        uint256 maliciousDkgResultSlashingAmount
    );

    event UnauthorizedSigningSlashingAmountUpdateStarted(
        uint256 unauthorizedSigningSlashingAmount,
        uint256 timestamp
    );
    event UnauthorizedSigningSlashingAmountUpdated(
        uint256 unauthorizedSigningSlashingAmount
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

    event RelayRequestFeeUpdateStarted(
        uint256 relayRequestFee,
        uint256 timestamp
    );
    event RelayRequestFeeUpdated(uint256 relayRequestFee);

    event RelayEntrySoftTimeoutUpdateStarted(
        uint256 relayEntrySoftTimeout,
        uint256 timestamp
    );
    event RelayEntrySoftTimeoutUpdated(uint256 relayEntrySoftTimeout);

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

    event DkgResultSubmissionTimeoutUpdateStarted(
        uint256 dkgResultSubmissionTimeout,
        uint256 timestamp
    );
    event DkgResultSubmissionTimeoutUpdated(uint256 dkgResultSubmissionTimeout);

    event DkgSubmitterPrecedencePeriodLengthUpdateStarted(
        uint256 submitterPrecedencePeriodLength,
        uint256 timestamp
    );
    event DkgSubmitterPrecedencePeriodLengthUpdated(
        uint256 submitterPrecedencePeriodLength
    );

    event GovernanceDelayUpdateStarted(
        uint256 governanceDelay,
        uint256 timestamp
    );
    event GovernanceDelayUpdated(uint256 governanceDelay);

    event RandomBeaconOwnershipTransferStarted(
        address newRandomBeaconOwner,
        uint256 timestamp
    );
    event RandomBeaconOwnershipTransferred(address newRandomBeaconOwner);

    /// @notice Reverts if called before the governance delay elapses.
    /// @param changeInitiatedTimestamp Timestamp indicating the beginning
    ///        of the change.
    modifier onlyAfterGovernanceDelay(
        Data storage self,
        uint256 changeInitiatedTimestamp
    ) {
        /* solhint-disable not-rely-on-time */
        require(changeInitiatedTimestamp > 0, "Change not initiated");
        require(
            block.timestamp - changeInitiatedTimestamp >= self.governanceDelay,
            "Governance delay has not elapsed"
        );
        _;
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Inits governance delay param.
    /// @param _governanceDelay Governance delay
    function init(Data storage self, uint256 _governanceDelay) internal {
        self.governanceDelay = _governanceDelay;
    }

    /// @notice Begins the DKG result submission reward update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgResultSubmissionReward New DKG result submission reward
    function beginDkgResultSubmissionRewardUpdate(
        Data storage self,
        uint256 _newDkgResultSubmissionReward
    ) external {
        /* solhint-disable not-rely-on-time */
        self.newDkgResultSubmissionReward = _newDkgResultSubmissionReward;
        self.dkgResultSubmissionRewardChangeInitiated = block.timestamp;
        emit DkgResultSubmissionRewardUpdateStarted(
            _newDkgResultSubmissionReward,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG result submission reward update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgResultSubmissionRewardUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.dkgResultSubmissionRewardChangeInitiated
        )
    {
        emit DkgResultSubmissionRewardUpdated(
            self.newDkgResultSubmissionReward
        );
        self.dkgResultSubmissionRewardChangeInitiated = 0;
        self.newDkgResultSubmissionReward = 0;
    }

    /// @notice Begins the sortition pool unlocking reward update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newSortitionPoolUnlockingReward New sortition pool unlocking reward
    function beginSortitionPoolUnlockingRewardUpdate(
        Data storage self,
        uint256 _newSortitionPoolUnlockingReward
    ) external {
        /* solhint-disable not-rely-on-time */
        self.newSortitionPoolUnlockingReward = _newSortitionPoolUnlockingReward;
        self.sortitionPoolUnlockingRewardChangeInitiated = block.timestamp;
        emit SortitionPoolUnlockingRewardUpdateStarted(
            _newSortitionPoolUnlockingReward,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the sortition pool unlocking reward update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeSortitionPoolUnlockingRewardUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.sortitionPoolUnlockingRewardChangeInitiated
        )
    {
        emit SortitionPoolUnlockingRewardUpdated(
            self.newSortitionPoolUnlockingReward
        );
        self.sortitionPoolUnlockingRewardChangeInitiated = 0;
        self.newSortitionPoolUnlockingReward = 0;
    }

    /// @notice Begins the ineligible operator notifier reward update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newIneligibleOperatorNotifierReward New ineligible operator
    ///        notifier reward.
    function beginIneligibleOperatorNotifierRewardUpdate(
        Data storage self,
        uint256 _newIneligibleOperatorNotifierReward
    ) external {
        /* solhint-disable not-rely-on-time */
        self
            .newIneligibleOperatorNotifierReward = _newIneligibleOperatorNotifierReward;
        self.ineligibleOperatorNotifierRewardChangeInitiated = block.timestamp;
        emit IneligibleOperatorNotifierRewardUpdateStarted(
            _newIneligibleOperatorNotifierReward,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the ineligible operator notifier reward update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeIneligibleOperatorNotifierRewardUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.ineligibleOperatorNotifierRewardChangeInitiated
        )
    {
        emit IneligibleOperatorNotifierRewardUpdated(
            self.newIneligibleOperatorNotifierReward
        );
        self.ineligibleOperatorNotifierRewardChangeInitiated = 0;
        self.newIneligibleOperatorNotifierReward = 0;
    }

    /// @notice Begins the sortition pool rewards ban duration update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newSortitionPoolRewardsBanDuration New sortition pool rewards
    ///        ban duration.
    function beginSortitionPoolRewardsBanDurationUpdate(
        Data storage self,
        uint256 _newSortitionPoolRewardsBanDuration
    ) external {
        /* solhint-disable not-rely-on-time */
        self
            .newSortitionPoolRewardsBanDuration = _newSortitionPoolRewardsBanDuration;
        self.sortitionPoolRewardsBanDurationChangeInitiated = block.timestamp;
        emit SortitionPoolRewardsBanDurationUpdateStarted(
            _newSortitionPoolRewardsBanDuration,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the sortition pool rewards ban duration update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeSortitionPoolRewardsBanDurationUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.sortitionPoolRewardsBanDurationChangeInitiated
        )
    {
        emit SortitionPoolRewardsBanDurationUpdated(
            self.newSortitionPoolRewardsBanDuration
        );
        self.sortitionPoolRewardsBanDurationChangeInitiated = 0;
        self.newSortitionPoolRewardsBanDuration = 0;
    }

    // ok
    /// @notice Begins the unauthorized signing notification reward multiplier
    ///         update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newUnauthorizedSigningNotificationRewardMultiplier New unauthorized
    ///         signing notification reward multiplier.
    function beginUnauthorizedSigningNotificationRewardMultiplierUpdate(
        Data storage self,
        uint256 _newUnauthorizedSigningNotificationRewardMultiplier
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newUnauthorizedSigningNotificationRewardMultiplier <= 100,
            "Maximum value is 100"
        );

        self
            .newUnauthorizedSigningNotificationRewardMultiplier = _newUnauthorizedSigningNotificationRewardMultiplier;
        self
            .unauthorizedSigningNotificationRewardMultiplierChangeInitiated = block
            .timestamp;
        emit UnauthorizedSigningNotificationRewardMultiplierUpdateStarted(
            _newUnauthorizedSigningNotificationRewardMultiplier,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the unauthorized signing notification reward
    ///         multiplier update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeUnauthorizedSigningNotificationRewardMultiplierUpdate(
        Data storage self
    )
        external
        onlyAfterGovernanceDelay(
            self,
            self.unauthorizedSigningNotificationRewardMultiplierChangeInitiated
        )
    {
        emit UnauthorizedSigningNotificationRewardMultiplierUpdated(
            self.newUnauthorizedSigningNotificationRewardMultiplier
        );
        self.unauthorizedSigningNotificationRewardMultiplierChangeInitiated = 0;
        self.newUnauthorizedSigningNotificationRewardMultiplier = 0;
    }

    /// @notice Begins the relay entry timeout notification reward multiplier
    ///         update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntryTimeoutNotificationRewardMultiplier New relay
    ///        entry timeout notification reward multiplier.
    function beginRelayEntryTimeoutNotificationRewardMultiplierUpdate(
        Data storage self,
        uint256 _newRelayEntryTimeoutNotificationRewardMultiplier
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newRelayEntryTimeoutNotificationRewardMultiplier <= 100,
            "Maximum value is 100"
        );

        self
            .newRelayEntryTimeoutNotificationRewardMultiplier = _newRelayEntryTimeoutNotificationRewardMultiplier;
        self
            .relayEntryTimeoutNotificationRewardMultiplierChangeInitiated = block
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
    function finalizeRelayEntryTimeoutNotificationRewardMultiplierUpdate(
        Data storage self
    )
        external
        onlyAfterGovernanceDelay(
            self,
            self.relayEntryTimeoutNotificationRewardMultiplierChangeInitiated
        )
    {
        emit RelayEntryTimeoutNotificationRewardMultiplierUpdated(
            self.newRelayEntryTimeoutNotificationRewardMultiplier
        );
        self.relayEntryTimeoutNotificationRewardMultiplierChangeInitiated = 0;
        self.newRelayEntryTimeoutNotificationRewardMultiplier = 0;
    }

    /// @notice Begins the DKG malicious result notification reward multiplier
    ///         update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgMaliciousResultNotificationRewardMultiplier New DKG
    ///        malicious result notification reward multiplier.
    function beginDkgMaliciousResultNotificationRewardMultiplierUpdate(
        Data storage self,
        uint256 _newDkgMaliciousResultNotificationRewardMultiplier
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newDkgMaliciousResultNotificationRewardMultiplier <= 100,
            "Maximum value is 100"
        );

        self
            .newDkgMaliciousResultNotificationRewardMultiplier = _newDkgMaliciousResultNotificationRewardMultiplier;
        self
            .dkgMaliciousResultNotificationRewardMultiplierChangeInitiated = block
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
    function finalizeDkgMaliciousResultNotificationRewardMultiplierUpdate(
        Data storage self
    )
        external
        onlyAfterGovernanceDelay(
            self,
            self.dkgMaliciousResultNotificationRewardMultiplierChangeInitiated
        )
    {
        emit DkgMaliciousResultNotificationRewardMultiplierUpdated(
            self.newDkgMaliciousResultNotificationRewardMultiplier
        );
        self.dkgMaliciousResultNotificationRewardMultiplierChangeInitiated = 0;
        self.newDkgMaliciousResultNotificationRewardMultiplier = 0;
    }

    /// @notice Begins the relay entry submission failure slashing amount update
    ///         process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntrySubmissionFailureSlashingAmount New relay entry
    ///        submission failure slashing amount
    function beginRelayEntrySubmissionFailureSlashingAmountUpdate(
        Data storage self,
        uint96 _newRelayEntrySubmissionFailureSlashingAmount
    ) external {
        /* solhint-disable not-rely-on-time */
        self
            .newRelayEntrySubmissionFailureSlashingAmount = _newRelayEntrySubmissionFailureSlashingAmount;
        self.relayEntrySubmissionFailureSlashingAmountChangeInitiated = block
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
    function finalizeRelayEntrySubmissionFailureSlashingAmountUpdate(
        Data storage self
    )
        external
        onlyAfterGovernanceDelay(
            self,
            self.relayEntrySubmissionFailureSlashingAmountChangeInitiated
        )
    {
        emit RelayEntrySubmissionFailureSlashingAmountUpdated(
            self.newRelayEntrySubmissionFailureSlashingAmount
        );
        self.relayEntrySubmissionFailureSlashingAmountChangeInitiated = 0;
        self.newRelayEntrySubmissionFailureSlashingAmount = 0;
    }

    /// @notice Begins the malicious DKG result slashing amount update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newMaliciousDkgResultSlashingAmount New malicious DKG result
    ///        slashing amount
    function beginMaliciousDkgResultSlashingAmountUpdate(
        Data storage self,
        uint256 _newMaliciousDkgResultSlashingAmount
    ) external {
        /* solhint-disable not-rely-on-time */
        self
            .newMaliciousDkgResultSlashingAmount = _newMaliciousDkgResultSlashingAmount;
        self.maliciousDkgResultSlashingAmountChangeInitiated = block.timestamp;
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
    function finalizeMaliciousDkgResultSlashingAmountUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.maliciousDkgResultSlashingAmountChangeInitiated
        )
    {
        emit MaliciousDkgResultSlashingAmountUpdated(
            self.newMaliciousDkgResultSlashingAmount
        );
        self.maliciousDkgResultSlashingAmountChangeInitiated = 0;
        self.newMaliciousDkgResultSlashingAmount = 0;
    }

    /// @notice Begins the unauthorized signing slashing amount update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newUnauthorizedSigningSlashingAmount New unauthorized signing
    ///        slashing amount
    function beginUnauthorizedSigningSlashingAmountUpdate(
        Data storage self,
        uint256 _newUnauthorizedSigningSlashingAmount
    ) external {
        /* solhint-disable not-rely-on-time */
        self
            .newUnauthorizedSigningSlashingAmount = _newUnauthorizedSigningSlashingAmount;
        self.unauthorizedSigningSlashingAmountChangeInitiated = block.timestamp;
        emit UnauthorizedSigningSlashingAmountUpdateStarted(
            _newUnauthorizedSigningSlashingAmount,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the unauthorized signing slashing amount update
    ///         process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeUnauthorizedSigningSlashingAmountUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.unauthorizedSigningSlashingAmountChangeInitiated
        )
    {
        emit UnauthorizedSigningSlashingAmountUpdated(
            self.newUnauthorizedSigningSlashingAmount
        );
        self.unauthorizedSigningSlashingAmountChangeInitiated = 0;
        self.newUnauthorizedSigningSlashingAmount = 0;
    }

    /// @notice Begins the minimum authorization amount update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newMinimumAuthorization New minimum authorization amount.
    function beginMinimumAuthorizationUpdate(
        Data storage self,
        uint96 _newMinimumAuthorization
    ) external {
        /* solhint-disable not-rely-on-time */
        self.newMinimumAuthorization = _newMinimumAuthorization;
        self.minimumAuthorizationChangeInitiated = block.timestamp;
        emit MinimumAuthorizationUpdateStarted(
            _newMinimumAuthorization,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the minimum authorization amount update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeMinimumAuthorizationUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(self, self.minimumAuthorizationChangeInitiated)
    {
        emit MinimumAuthorizationUpdated(self.newMinimumAuthorization);
        self.minimumAuthorizationChangeInitiated = 0;
        self.newMinimumAuthorization = 0;
    }

    /// @notice Begins the authorization decrease delay update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newAuthorizationDecreaseDelay New authorization decrease delay
    function beginAuthorizationDecreaseDelayUpdate(
        Data storage self,
        uint64 _newAuthorizationDecreaseDelay
    ) external {
        /* solhint-disable not-rely-on-time */
        self.newAuthorizationDecreaseDelay = _newAuthorizationDecreaseDelay;
        self.authorizationDecreaseDelayChangeInitiated = block.timestamp;
        emit AuthorizationDecreaseDelayUpdateStarted(
            _newAuthorizationDecreaseDelay,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the authorization decrease delay update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeAuthorizationDecreaseDelayUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.authorizationDecreaseDelayChangeInitiated
        )
    {
        emit AuthorizationDecreaseDelayUpdated(
            self.newAuthorizationDecreaseDelay
        );
        self.authorizationDecreaseDelayChangeInitiated = 0;
        self.newAuthorizationDecreaseDelay = 0;
    }

    /// @notice Begins the relay request fee update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayRequestFee New relay request fee
    function beginRelayRequestFeeUpdate(
        Data storage self,
        uint256 _newRelayRequestFee
    ) external {
        /* solhint-disable not-rely-on-time */
        self.newRelayRequestFee = _newRelayRequestFee;
        self.relayRequestFeeChangeInitiated = block.timestamp;
        emit RelayRequestFeeUpdateStarted(_newRelayRequestFee, block.timestamp);
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the relay request fee update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayRequestFeeUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(self, self.relayRequestFeeChangeInitiated)
    {
        emit RelayRequestFeeUpdated(self.newRelayRequestFee);
        self.relayRequestFeeChangeInitiated = 0;
        self.newRelayRequestFee = 0;
    }

    /// @notice Begins the relay entry soft timeout update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntrySoftTimeout New relay entry submission timeout in blocks
    function beginRelayEntrySoftTimeoutUpdate(
        Data storage self,
        uint256 _newRelayEntrySoftTimeout
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newRelayEntrySoftTimeout > 0,
            "Relay entry soft timeout must be > 0"
        );
        self.newRelayEntrySoftTimeout = _newRelayEntrySoftTimeout;
        self.relayEntrySoftTimeoutChangeInitiated = block.timestamp;
        emit RelayEntrySoftTimeoutUpdateStarted(
            _newRelayEntrySoftTimeout,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the relay entry soft timeout update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayEntrySoftTimeoutUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.relayEntrySoftTimeoutChangeInitiated
        )
    {
        emit RelayEntrySoftTimeoutUpdated(self.newRelayEntrySoftTimeout);
        self.relayEntrySoftTimeoutChangeInitiated = 0;
        self.newRelayEntrySoftTimeout = 0;
    }

    /// @notice Begins the relay entry hard timeout update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntryHardTimeout New relay entry hard timeout in blocks
    function beginRelayEntryHardTimeoutUpdate(
        Data storage self,
        uint256 _newRelayEntryHardTimeout
    ) external {
        /* solhint-disable not-rely-on-time */
        self.newRelayEntryHardTimeout = _newRelayEntryHardTimeout;
        self.relayEntryHardTimeoutChangeInitiated = block.timestamp;
        emit RelayEntryHardTimeoutUpdateStarted(
            _newRelayEntryHardTimeout,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the relay entry hard timeout update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayEntryHardTimeoutUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.relayEntryHardTimeoutChangeInitiated
        )
    {
        emit RelayEntryHardTimeoutUpdated(self.newRelayEntryHardTimeout);
        self.relayEntryHardTimeoutChangeInitiated = 0;
        self.newRelayEntryHardTimeout = 0;
    }

    /// @notice Begins the callback gas limit update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newCallbackGasLimit New callback gas limit
    function beginCallbackGasLimitUpdate(
        Data storage self,
        uint256 _newCallbackGasLimit
    ) external {
        /* solhint-disable not-rely-on-time */
        // slither-disable-next-line too-many-digits
        require(
            _newCallbackGasLimit > 0 && _newCallbackGasLimit <= 1e6,
            "Callback gas limit must be > 0 and <= 1000000"
        );
        self.newCallbackGasLimit = _newCallbackGasLimit;
        self.callbackGasLimitChangeInitiated = block.timestamp;
        emit CallbackGasLimitUpdateStarted(
            _newCallbackGasLimit,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the callback gas limit update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeCallbackGasLimitUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(self, self.callbackGasLimitChangeInitiated)
    {
        emit CallbackGasLimitUpdated(self.newCallbackGasLimit);
        self.callbackGasLimitChangeInitiated = 0;
        self.newCallbackGasLimit = 0;
    }

    /// @notice Begins the group creation frequency update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newGroupCreationFrequency New group creation frequency
    function beginGroupCreationFrequencyUpdate(
        Data storage self,
        uint256 _newGroupCreationFrequency
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newGroupCreationFrequency > 0,
            "Group creation frequency must be > 0"
        );
        self.newGroupCreationFrequency = _newGroupCreationFrequency;
        self.groupCreationFrequencyChangeInitiated = block.timestamp;
        emit GroupCreationFrequencyUpdateStarted(
            _newGroupCreationFrequency,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the group creation frequency update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeGroupCreationFrequencyUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.groupCreationFrequencyChangeInitiated
        )
    {
        emit GroupCreationFrequencyUpdated(self.newGroupCreationFrequency);
        self.groupCreationFrequencyChangeInitiated = 0;
        self.newGroupCreationFrequency = 0;
    }

    /// @notice Begins the group lifetime update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newGroupLifetime New group lifetime in blocks
    function beginGroupLifetimeUpdate(
        Data storage self,
        uint256 _newGroupLifetime
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newGroupLifetime >= 1 days && _newGroupLifetime <= 2 weeks,
            "Group lifetime must be >= 1 day and <= 2 weeks"
        );
        self.newGroupLifetime = _newGroupLifetime;
        self.groupLifetimeChangeInitiated = block.timestamp;
        emit GroupLifetimeUpdateStarted(_newGroupLifetime, block.timestamp);
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the group creation frequency update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeGroupLifetimeUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(self, self.groupLifetimeChangeInitiated)
    {
        emit GroupLifetimeUpdated(self.newGroupLifetime);
        self.groupLifetimeChangeInitiated = 0;
        self.newGroupLifetime = 0;
    }

    /// @notice Begins the DKG result challenge period length update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgResultChallengePeriodLength New DKG result challenge
    ///        period length in blocks
    function beginDkgResultChallengePeriodLengthUpdate(
        Data storage self,
        uint256 _newDkgResultChallengePeriodLength
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newDkgResultChallengePeriodLength >= 10,
            "DKG result challenge period length must be >= 10"
        );
        self
            .newDkgResultChallengePeriodLength = _newDkgResultChallengePeriodLength;
        self.dkgResultChallengePeriodLengthChangeInitiated = block.timestamp;
        emit DkgResultChallengePeriodLengthUpdateStarted(
            _newDkgResultChallengePeriodLength,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG result challenge period length update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgResultChallengePeriodLengthUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.dkgResultChallengePeriodLengthChangeInitiated
        )
    {
        emit DkgResultChallengePeriodLengthUpdated(
            self.newDkgResultChallengePeriodLength
        );
        self.dkgResultChallengePeriodLengthChangeInitiated = 0;
        self.newDkgResultChallengePeriodLength = 0;
    }

    /// @notice Begins the DKG result submission timeout update
    ///         process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgResultSubmissionTimeout New DKG result submission
    ///        timeout in blocks
    function beginDkgResultSubmissionTimeoutUpdate(
        Data storage self,
        uint256 _newDkgResultSubmissionTimeout
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newDkgResultSubmissionTimeout > 0,
            "DKG result submission timeout must be > 0"
        );
        self.newDkgResultSubmissionTimeout = _newDkgResultSubmissionTimeout;
        self.dkgResultSubmissionTimeoutChangeInitiated = block.timestamp;
        emit DkgResultSubmissionTimeoutUpdateStarted(
            _newDkgResultSubmissionTimeout,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG result submission timeout update
    ///         process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgResultSubmissionTimeoutUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.dkgResultSubmissionTimeoutChangeInitiated
        )
    {
        emit DkgResultSubmissionTimeoutUpdated(
            self.newDkgResultSubmissionTimeout
        );
        self.dkgResultSubmissionTimeoutChangeInitiated = 0;
        self.newDkgResultSubmissionTimeout = 0;
    }

    /// @notice Begins the DKG submitter precedence period length.
    /// @dev Can be called only by the contract owner.
    /// @param _newSubmitterPrecedencePeriodLength New DKG submitter precedence
    ///        period length in blocks
    function beginDkgSubmitterPrecedencePeriodLengthUpdate(
        Data storage self,
        uint256 _newSubmitterPrecedencePeriodLength
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newSubmitterPrecedencePeriodLength > 0,
            "DKG submitter precedence period length must be > 0"
        );
        self
            .newSubmitterPrecedencePeriodLength = _newSubmitterPrecedencePeriodLength;
        self.dkgSubmitterPrecedencePeriodLengthChangeInitiated = block
            .timestamp;
        emit DkgSubmitterPrecedencePeriodLengthUpdateStarted(
            _newSubmitterPrecedencePeriodLength,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG submitter precedence period length.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgSubmitterPrecedencePeriodLengthUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.dkgSubmitterPrecedencePeriodLengthChangeInitiated
        )
    {
        emit DkgSubmitterPrecedencePeriodLengthUpdated(
            self.newSubmitterPrecedencePeriodLength
        );
        self.dkgSubmitterPrecedencePeriodLengthChangeInitiated = 0;
        self.newSubmitterPrecedencePeriodLength = 0;
    }

    /// @notice Begins the governance delay update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newGovernanceDelay New governance delay
    function beginGovernanceDelayUpdate(
        Data storage self,
        uint256 _newGovernanceDelay
    ) external {
        self.newGovernanceDelay = _newGovernanceDelay;
        /* solhint-disable not-rely-on-time */
        self.governanceDelayChangeInitiated = block.timestamp;
        emit GovernanceDelayUpdateStarted(_newGovernanceDelay, block.timestamp);
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the governance delay update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeGovernanceDelayUpdate(Data storage self)
        external
        onlyAfterGovernanceDelay(self, self.governanceDelayChangeInitiated)
    {
        emit GovernanceDelayUpdated(self.newGovernanceDelay);
        self.governanceDelay = self.newGovernanceDelay;
        self.governanceDelayChangeInitiated = 0;
        self.newGovernanceDelay = 0;
    }

    /// @notice Begins the random beacon ownership transfer process.
    /// @dev Can be called only by the contract owner.
    function beginRandomBeaconOwnershipTransfer(
        Data storage self,
        address _newRandomBeaconOwner
    ) external {
        require(
            address(_newRandomBeaconOwner) != address(0),
            "New random beacon owner address cannot be zero"
        );
        self.newRandomBeaconOwner = _newRandomBeaconOwner;
        /* solhint-disable not-rely-on-time */
        self.randomBeaconOwnershipTransferInitiated = block.timestamp;
        emit RandomBeaconOwnershipTransferStarted(
            _newRandomBeaconOwner,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the random beacon ownership transfer process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRandomBeaconOwnershipTransfer(Data storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            self.randomBeaconOwnershipTransferInitiated
        )
    {
        emit RandomBeaconOwnershipTransferred(self.newRandomBeaconOwner);
        self.randomBeaconOwnershipTransferInitiated = 0;
        self.newRandomBeaconOwner = address(0);
    }

    /// @notice Get the time remaining until the governance delay can be updated.
    /// @return Remaining time in seconds.
    function getRemainingGovernanceDelayUpdateTime(Data storage self)
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(self, self.governanceDelayChangeInitiated);
    }

    /// @notice Get the time remaining until the DKG result submission reward
    ///         can be updated.
    /// @return Remaining time in seconds.
    function getRemainingDkgResultSubmissionRewardUpdateTime(Data storage self)
        internal
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                self,
                self.dkgResultSubmissionRewardChangeInitiated
            );
    }

    /// @notice Get the time remaining until the sortition pool unlocking reward
    ///         can be updated.
    /// @return Remaining time in seconds.
    function getRemainingSortitionPoolUnlockingRewardUpdateTime(
        Data storage self
    ) internal view returns (uint256) {
        return
            getRemainingChangeTime(
                self,
                self.sortitionPoolUnlockingRewardChangeInitiated
            );
    }

    /// @notice Get the time remaining until the ineligible operator notifier
    ///         reward can be updated.
    /// @return Remaining time in seconds.
    function getRemainingIneligibleOperatorNotifierRewardUpdateTime(
        Data storage self
    ) internal view returns (uint256) {
        return
            getRemainingChangeTime(
                self,
                self.ineligibleOperatorNotifierRewardChangeInitiated
            );
    }

    /// @notice Get the time remaining until the sortition pool rewards ban
    ///         duration can be updated.
    /// @return Remaining time in seconds.
    function getRemainingSortitionPoolRewardsBanDurationUpdateTime(
        Data storage self
    ) internal view returns (uint256) {
        return
            getRemainingChangeTime(
                self,
                self.sortitionPoolRewardsBanDurationChangeInitiated
            );
    }

    /// @notice Get the time remaining until the unauthorized signing
    ///         notification reward multiplier duration can be updated.
    /// @return Remaining time in seconds.
    function getRemainingUnauthorizedSigningNotificationRewardMultiplierUpdateTime(
        Data storage self
    ) external view returns (uint256) {
        return
            getRemainingChangeTime(
                self,
                self
                    .unauthorizedSigningNotificationRewardMultiplierChangeInitiated
            );
    }

    /// @notice Get the time remaining until the relay entry timeout
    ///         notification reward multiplier duration can be updated.
    /// @return Remaining time in seconds.
    function getRemainingRelayEntryTimeoutNotificationRewardMultiplierUpdateTime(
        Data storage self
    ) external view returns (uint256) {
        return
            getRemainingChangeTime(
                self,
                self
                    .relayEntryTimeoutNotificationRewardMultiplierChangeInitiated
            );
    }

    /// @notice Get the time remaining until the DKG malicious result
    ///         notification reward multiplier duration can be updated.
    /// @return Remaining time in seconds.
    function getRemainingDkgMaliciousResultNotificationRewardMultiplierUpdateTime(
        Data storage self
    ) external view returns (uint256) {
        return
            getRemainingChangeTime(
                self,
                self
                    .dkgMaliciousResultNotificationRewardMultiplierChangeInitiated
            );
    }

    /// @notice Get the time remaining until the relay entry submission failure
    ///         slashing amount can be updated.
    /// @return Remaining time in seconds.
    function getRemainingRelayEntrySubmissionFailureSlashingAmountUpdateTime(
        Data storage self
    ) external view returns (uint256) {
        return
            getRemainingChangeTime(
                self,
                self.relayEntrySubmissionFailureSlashingAmountChangeInitiated
            );
    }

    /// @notice Get the time remaining until the malicious DKG result
    ///         slashing amount can be updated.
    /// @return Remaining time in seconds.
    function getRemainingMaliciousDkgResultSlashingAmountUpdateTime(
        Data storage self
    ) external view returns (uint256) {
        return
            getRemainingChangeTime(
                self,
                self.maliciousDkgResultSlashingAmountChangeInitiated
            );
    }

    /// @notice Get the time remaining until the unauthorized signing
    ///         slashing amount can be updated.
    /// @return Remaining time in seconds.
    function getRemainingUnauthorizedSigningSlashingAmountUpdateTime(
        Data storage self
    ) external view returns (uint256) {
        return
            getRemainingChangeTime(
                self,
                self.unauthorizedSigningSlashingAmountChangeInitiated
            );
    }

    /// @notice Get the time remaining until the minimum authorization amount
    ///         can be updated.
    /// @return Remaining time in seconds.
    function getRemainingMimimumAuthorizationUpdateTime(Data storage self)
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                self,
                self.minimumAuthorizationChangeInitiated
            );
    }

    /// @notice Get the time remaining until the authorization decrease delay
    ///         can be updated.
    /// @return Remaining time in seconds.
    function getRemainingAuthorizationDecreaseDelayUpdateTime(Data storage self)
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                self,
                self.authorizationDecreaseDelayChangeInitiated
            );
    }

    /// @notice Get the time remaining until the relay request fee can be
    ///         updated.
    /// @return Remaining time in seconds.
    function getRemainingRelayRequestFeeUpdateTime(Data storage self)
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(self, self.relayRequestFeeChangeInitiated);
    }

    /// @notice Get the time remaining until the relay entry submission soft
    ///         timeout can be updated.
    /// @return Remaining time in seconds.
    function getRemainingRelayEntrySoftTimeoutUpdateTime(Data storage self)
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                self,
                self.relayEntrySoftTimeoutChangeInitiated
            );
    }

    /// @notice Get the time remaining until the relay entry hard timeout can be
    ///         updated.
    /// @return Remaining time in seconds.
    function getRemainingRelayEntryHardTimeoutUpdateTime(Data storage self)
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                self,
                self.relayEntryHardTimeoutChangeInitiated
            );
    }

    /// @notice Get the time remaining until the callback gas limit can be
    ///         updated.
    /// @return Remaining time in seconds.
    function getRemainingCallbackGasLimitUpdateTime(Data storage self)
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(self, self.callbackGasLimitChangeInitiated);
    }

    /// @notice Get the time remaining until the group creation frequency can be
    ///         updated.
    /// @return Remaining time in seconds.
    function getRemainingGroupCreationFrequencyUpdateTime(Data storage self)
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                self,
                self.groupCreationFrequencyChangeInitiated
            );
    }

    /// @notice Get the time remaining until the group lifetime can be updated.
    /// @return Remaining time in seconds.
    function getRemainingGroupLifetimeUpdateTime(Data storage self)
        external
        view
        returns (uint256)
    {
        return getRemainingChangeTime(self, self.groupLifetimeChangeInitiated);
    }

    /// @notice Get the time remaining until the DKG result challenge period
    ///         length can be updated.
    /// @return Remaining time in seconds.
    function getRemainingDkgResultChallengePeriodLengthUpdateTime(
        Data storage self
    ) external view returns (uint256) {
        return
            getRemainingChangeTime(
                self,
                self.dkgResultChallengePeriodLengthChangeInitiated
            );
    }

    /// @notice Get the time remaining until the DKG result submission timeout
    ///         can be updated.
    /// @return Remaining time in seconds.
    function getRemainingDkgResultSubmissionTimeoutUpdateTime(Data storage self)
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                self,
                self.dkgResultSubmissionTimeoutChangeInitiated
            );
    }

    /// @notice Get the time remaining until the wallet owner can be updated.
    /// @return Remaining time in seconds.
    function getRemainingDkgSubmitterPrecedencePeriodLengthUpdateTime(
        Data storage self
    ) external view returns (uint256) {
        return
            getRemainingChangeTime(
                self,
                self.dkgSubmitterPrecedencePeriodLengthChangeInitiated
            );
    }

    /// @notice Get the time remaining until the random beacon ownership can
    ///         be transferred.
    /// @return Remaining time in seconds.
    function getRemainingRandomBeaconOwnershipTransferDelayTime(
        Data storage self
    ) external view returns (uint256) {
        return
            getRemainingChangeTime(
                self,
                self.randomBeaconOwnershipTransferInitiated
            );
    }

    /// @notice Gets the new random beacon owner.
    function getNewRandomBeaconOwner(Data storage self)
        internal
        view
        returns (address)
    {
        return self.newRandomBeaconOwner;
    }

    /// @notice Gets the governance delay
    function getGovernanceDelay(Data storage self)
        external
        view
        returns (uint256)
    {
        return self.governanceDelay;
    }

    /// @notice Gets the new dkg result challenge period length
    function getNewDkgResultChallengePeriodLength(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newDkgResultChallengePeriodLength;
    }

    /// @notice Gets the new dkg result submission timeout
    function getNewDkgResultSubmissionTimeout(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newDkgResultSubmissionTimeout;
    }

    /// @notice Gets the new dkg submitter precedence period length
    function getNewDkgSubmitterPrecedencePeriodLength(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newSubmitterPrecedencePeriodLength;
    }

    /// @notice Gets the new group creation frequency
    function getNewGroupCreationFrequency(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newGroupCreationFrequency;
    }

    /// @notice Gets the new group lifetime
    function getNewGroupLifetime(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newGroupLifetime;
    }

    /// @notice Gets the new dkg result submission reward
    function getNewDkgResultSubmissionReward(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newDkgResultSubmissionReward;
    }

    /// @notice Gets the new sortition pool unlocking reward
    function getNewSortitionPoolUnlockingReward(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newSortitionPoolUnlockingReward;
    }

    /// @notice Gets the new ineligible operator notifier reward
    function getNewIneligibleOperatorNotifierReward(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newIneligibleOperatorNotifierReward;
    }

    /// @notice Gets the new sortition pool rewards ban duration
    function getNewSortitionPoolRewardsBanDuration(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newSortitionPoolRewardsBanDuration;
    }

    /// @notice Gets the new unaothorized signing notification reward multiplier
    function getNewUnauthorizedSigningNotificationRewardMultiplier(
        Data storage self
    ) internal view returns (uint256) {
        return self.newUnauthorizedSigningNotificationRewardMultiplier;
    }

    /// @notice Gets the new relay entry timeout notification rewards multiplier
    function getNewRelayEntryTimeoutNotificationRewardMultiplier(
        Data storage self
    ) internal view returns (uint256) {
        return self.newRelayEntryTimeoutNotificationRewardMultiplier;
    }

    /// @notice Gets the new dkg malicious notification reward multiplier
    function getNewDkgMaliciousResultNotificationRewardMultiplier(
        Data storage self
    ) internal view returns (uint256) {
        return self.newDkgMaliciousResultNotificationRewardMultiplier;
    }

    /// @notice Gets the new relay entry submission failure slashing amount
    function getNewRelayEntrySubmissionFailureSlashingAmount(Data storage self)
        internal
        view
        returns (uint96)
    {
        return self.newRelayEntrySubmissionFailureSlashingAmount;
    }

    /// @notice Gets the new malicious dkg result slashing amount
    function getNewMaliciousDkgResultSlashingAmount(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newMaliciousDkgResultSlashingAmount;
    }

    /// @notice Gets the new unauthorized signing slashing amount
    function getNewUnauthorizedSigningSlashingAmount(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newUnauthorizedSigningSlashingAmount;
    }

    /// @notice Gets the new minimum authorization
    function getNewMinimumAuthorization(Data storage self)
        internal
        view
        returns (uint96)
    {
        return self.newMinimumAuthorization;
    }

    /// @notice Gets the new authorization decrease delay
    function getNewAuthorizationDecreaseDelay(Data storage self)
        internal
        view
        returns (uint64)
    {
        return self.newAuthorizationDecreaseDelay;
    }

    /// @notice Gets the new relay request fee
    function getNewRelayRequestFee(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newRelayRequestFee;
    }

    /// @notice Gets the new relay entry soft timeout
    function getNewRelayEntrySoftTimeout(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newRelayEntrySoftTimeout;
    }

    /// @notice Gets the new relay entry hard timeout
    function getNewRelayEntryHardTimeout(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newRelayEntryHardTimeout;
    }

    /// @notice Gets the new callback gas limit
    function getNewCallbackGasLimit(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newCallbackGasLimit;
    }

    /// @notice Gets the time remaining until the governable parameter update
    ///         can be committed.
    /// @param changeTimestamp Timestamp indicating the beginning of the change.
    /// @return Remaining time in seconds.
    function getRemainingChangeTime(Data storage self, uint256 changeTimestamp)
        internal
        view
        returns (uint256)
    {
        require(changeTimestamp > 0, "Change not initiated");
        /* solhint-disable-next-line not-rely-on-time */
        uint256 elapsed = block.timestamp - changeTimestamp;
        if (elapsed >= self.governanceDelay) {
            return 0;
        }

        return self.governanceDelay - elapsed;
    }
}
