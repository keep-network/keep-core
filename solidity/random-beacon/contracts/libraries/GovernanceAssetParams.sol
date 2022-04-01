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

/// TODO: add desc
library GovernanceAssetParams {
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
        uint256 newRelayEntrySubmissionFailureSlashingAmount;
        uint256 relayEntrySubmissionFailureSlashingAmountChangeInitiated;
        uint256 newMaliciousDkgResultSlashingAmount;
        uint256 maliciousDkgResultSlashingAmountChangeInitiated;
        uint256 newUnauthorizedSigningSlashingAmount;
        uint256 unauthorizedSigningSlashingAmountChangeInitiated;
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

    event UnauthorizedSigningSlashingAmountUpdateStarted(
        uint256 unauthorizedSigningSlashingAmount,
        uint256 timestamp
    );
    event UnauthorizedSigningSlashingAmountUpdated(
        uint256 unauthorizedSigningSlashingAmount
    );

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
        uint256 _newRelayEntrySubmissionFailureSlashingAmount
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

    function getNewDkgResultSubmissionReward(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newDkgResultSubmissionReward;
    }

    function getNewSortitionPoolUnlockingReward(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newSortitionPoolUnlockingReward;
    }

    function getNewIneligibleOperatorNotifierReward(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newIneligibleOperatorNotifierReward;
    }

    function getNewSortitionPoolRewardsBanDuration(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newSortitionPoolRewardsBanDuration;
    }

    function getNewUnauthorizedSigningNotificationRewardMultiplier(
        Data storage self
    ) internal view returns (uint256) {
        return self.newUnauthorizedSigningNotificationRewardMultiplier;
    }

    function getNewRelayEntryTimeoutNotificationRewardMultiplier(
        Data storage self
    ) internal view returns (uint256) {
        return self.newRelayEntryTimeoutNotificationRewardMultiplier;
    }

    function getNewDkgMaliciousResultNotificationRewardMultiplier(
        Data storage self
    ) internal view returns (uint256) {
        return self.newDkgMaliciousResultNotificationRewardMultiplier;
    }

    function getNewRelayEntrySubmissionFailureSlashingAmount(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newRelayEntrySubmissionFailureSlashingAmount;
    }

    function getNewMaliciousDkgResultSlashingAmount(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newMaliciousDkgResultSlashingAmount;
    }

    function getNewUnauthorizedSigningSlashingAmount(Data storage self)
        internal
        view
        returns (uint256)
    {
        return self.newUnauthorizedSigningSlashingAmount;
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
