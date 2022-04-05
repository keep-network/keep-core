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

import "./RandomBeacon.sol";
import "./libraries/GovernanceBeaconParams.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @title Keep Random Beacon Governance
/// @notice Owns the `RandomBeacon` contract and is responsible for updating its
///         governable parameters in respect to governance delay individual
///         for each parameter.
contract RandomBeaconGovernance is Ownable {
    using GovernanceBeaconParams for GovernanceBeaconParams.Data;

    GovernanceBeaconParams.Data internal governanceBeaconParams;

    RandomBeacon public randomBeacon;

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

    constructor(RandomBeacon _randomBeacon, uint256 _governanceDelay) {
        governanceBeaconParams.init(_governanceDelay);

        randomBeacon = _randomBeacon;
    }

    /// @notice Begins the governance delay update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newGovernanceDelay New governance delay
    function beginGovernanceDelayUpdate(uint256 _newGovernanceDelay)
        external
        onlyOwner
    {
        governanceBeaconParams.beginGovernanceDelayUpdate(_newGovernanceDelay);
    }

    function getGovernanceDelay() external view returns (uint256) {
        return governanceBeaconParams.getGovernanceDelay();
    }

    /// @notice Finalizes the governance delay update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeGovernanceDelayUpdate() external onlyOwner {
        governanceBeaconParams.finalizeGovernanceDelayUpdate();
    }

    /// @notice Begins the random beacon ownership transfer process.
    /// @dev Can be called only by the contract owner.
    function beginRandomBeaconOwnershipTransfer(address _newRandomBeaconOwner)
        external
        onlyOwner
    {
        governanceBeaconParams.beginRandomBeaconOwnershipTransfer(
            _newRandomBeaconOwner
        );
    }

    /// @notice Finalizes the random beacon ownership transfer process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRandomBeaconOwnershipTransfer() external onlyOwner {
        randomBeacon.transferOwnership(
            governanceBeaconParams.getNewRandomBeaconOwner()
        );

        governanceBeaconParams.finalizeRandomBeaconOwnershipTransfer();
    }

    /// @notice Begins the relay request fee update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayRequestFee New relay request fee
    function beginRelayRequestFeeUpdate(uint256 _newRelayRequestFee)
        external
        onlyOwner
    {
        governanceBeaconParams.beginRelayRequestFeeUpdate(_newRelayRequestFee);
    }

    /// @notice Finalizes the relay request fee update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayRequestFeeUpdate() external onlyOwner {
        randomBeacon.updateRelayEntryParameters(
            governanceBeaconParams.getNewRelayRequestFee(),
            randomBeacon.relayEntrySoftTimeout(),
            randomBeacon.relayEntryHardTimeout(),
            randomBeacon.callbackGasLimit()
        );

        governanceBeaconParams.finalizeRelayRequestFeeUpdate();
    }

    /// @notice Begins the relay entry soft timeout update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntrySoftTimeout New relay entry submission timeout in blocks
    function beginRelayEntrySoftTimeoutUpdate(uint256 _newRelayEntrySoftTimeout)
        external
        onlyOwner
    {
        governanceBeaconParams.beginRelayEntrySoftTimeoutUpdate(
            _newRelayEntrySoftTimeout
        );
    }

    /// @notice Finalizes the relay entry soft timeout update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayEntrySoftTimeoutUpdate() external onlyOwner {
        randomBeacon.updateRelayEntryParameters(
            randomBeacon.relayRequestFee(),
            governanceBeaconParams.getNewRelayEntrySoftTimeout(),
            randomBeacon.relayEntryHardTimeout(),
            randomBeacon.callbackGasLimit()
        );

        governanceBeaconParams.finalizeRelayEntrySoftTimeoutUpdate();
    }

    /// @notice Begins the relay entry hard timeout update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntryHardTimeout New relay entry hard timeout in blocks
    function beginRelayEntryHardTimeoutUpdate(uint256 _newRelayEntryHardTimeout)
        external
        onlyOwner
    {
        governanceBeaconParams.beginRelayEntryHardTimeoutUpdate(
            _newRelayEntryHardTimeout
        );
    }

    /// @notice Finalizes the relay entry hard timeout update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayEntryHardTimeoutUpdate() external onlyOwner {
        randomBeacon.updateRelayEntryParameters(
            randomBeacon.relayRequestFee(),
            randomBeacon.relayEntrySoftTimeout(),
            governanceBeaconParams.getNewRelayEntryHardTimeout(),
            randomBeacon.callbackGasLimit()
        );
        governanceBeaconParams.finalizeRelayEntryHardTimeoutUpdate();
    }

    /// @notice Begins the callback gas limit update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newCallbackGasLimit New callback gas limit
    function beginCallbackGasLimitUpdate(uint256 _newCallbackGasLimit)
        external
        onlyOwner
    {
        governanceBeaconParams.beginCallbackGasLimitUpdate(
            _newCallbackGasLimit
        );
    }

    /// @notice Finalizes the callback gas limit update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeCallbackGasLimitUpdate() external onlyOwner {
        randomBeacon.updateRelayEntryParameters(
            randomBeacon.relayRequestFee(),
            randomBeacon.relayEntrySoftTimeout(),
            randomBeacon.relayEntryHardTimeout(),
            governanceBeaconParams.getNewCallbackGasLimit()
        );

        governanceBeaconParams.finalizeCallbackGasLimitUpdate();
    }

    /// @notice Begins the group creation frequency update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newGroupCreationFrequency New group creation frequency
    function beginGroupCreationFrequencyUpdate(
        uint256 _newGroupCreationFrequency
    ) external onlyOwner {
        governanceBeaconParams.beginGroupCreationFrequencyUpdate(
            _newGroupCreationFrequency
        );
    }

    /// @notice Finalizes the group creation frequency update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeGroupCreationFrequencyUpdate() external onlyOwner {
        randomBeacon.updateGroupCreationParameters(
            governanceBeaconParams.getNewGroupCreationFrequency(),
            randomBeacon.groupLifetime()
        );

        governanceBeaconParams.finalizeGroupCreationFrequencyUpdate();
    }

    /// @notice Begins the group lifetime update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newGroupLifetime New group lifetime in blocks
    function beginGroupLifetimeUpdate(uint256 _newGroupLifetime)
        external
        onlyOwner
    {
        governanceBeaconParams.beginGroupLifetimeUpdate(_newGroupLifetime);
    }

    /// @notice Finalizes the group creation frequency update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeGroupLifetimeUpdate() external onlyOwner {
        randomBeacon.updateGroupCreationParameters(
            randomBeacon.groupCreationFrequency(),
            governanceBeaconParams.getNewGroupLifetime()
        );

        governanceBeaconParams.finalizeGroupLifetimeUpdate();
    }

    /// @notice Begins the DKG result challenge period length update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgResultChallengePeriodLength New DKG result challenge
    ///        period length in blocks
    function beginDkgResultChallengePeriodLengthUpdate(
        uint256 _newDkgResultChallengePeriodLength
    ) external onlyOwner {
        governanceBeaconParams.beginDkgResultChallengePeriodLengthUpdate(
            _newDkgResultChallengePeriodLength
        );
    }

    /// @notice Finalizes the DKG result challenge period length update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgResultChallengePeriodLengthUpdate() external onlyOwner {
        randomBeacon.updateDkgParameters(
            governanceBeaconParams.getNewDkgResultChallengePeriodLength(),
            randomBeacon.dkgResultSubmissionTimeout(),
            randomBeacon.dkgSubmitterPrecedencePeriodLength()
        );

        governanceBeaconParams.finalizeDkgResultChallengePeriodLengthUpdate();
    }

    /// @notice Begins the DKG result submission timeout update
    ///         process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgResultSubmissionTimeout New DKG result submission
    ///        timeout in blocks
    function beginDkgResultSubmissionTimeoutUpdate(
        uint256 _newDkgResultSubmissionTimeout
    ) external onlyOwner {
        governanceBeaconParams.beginDkgResultSubmissionTimeoutUpdate(
            _newDkgResultSubmissionTimeout
        );
    }

    /// @notice Finalizes the DKG result submission timeout update
    ///         process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgResultSubmissionTimeoutUpdate() external onlyOwner {
        randomBeacon.updateDkgParameters(
            randomBeacon.dkgResultChallengePeriodLength(),
            governanceBeaconParams.getNewDkgResultSubmissionTimeout(),
            randomBeacon.dkgSubmitterPrecedencePeriodLength()
        );
        governanceBeaconParams.finalizeDkgResultSubmissionTimeoutUpdate();
    }

    /// @notice Begins the DKG submitter precedence period length.
    /// @dev Can be called only by the contract owner.
    /// @param _newSubmitterPrecedencePeriodLength New DKG submitter precedence
    ///        period length in blocks
    function beginDkgSubmitterPrecedencePeriodLengthUpdate(
        uint256 _newSubmitterPrecedencePeriodLength
    ) external onlyOwner {
        governanceBeaconParams.beginDkgSubmitterPrecedencePeriodLengthUpdate(
            _newSubmitterPrecedencePeriodLength
        );
    }

    /// @notice Finalizes the DKG submitter precedence period length.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
        external
        onlyOwner
    {
        randomBeacon.updateDkgParameters(
            randomBeacon.dkgResultChallengePeriodLength(),
            randomBeacon.dkgResultSubmissionTimeout(),
            governanceBeaconParams.getNewDkgSubmitterPrecedencePeriodLength()
        );
        governanceBeaconParams
            .finalizeDkgSubmitterPrecedencePeriodLengthUpdate();
    }

    /// @notice Begins the DKG result submission reward update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgResultSubmissionReward New DKG result submission reward
    function beginDkgResultSubmissionRewardUpdate(
        uint256 _newDkgResultSubmissionReward
    ) external onlyOwner {
        governanceBeaconParams.beginDkgResultSubmissionRewardUpdate(
            _newDkgResultSubmissionReward
        );
    }

    /// @notice Finalizes the DKG result submission reward update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgResultSubmissionRewardUpdate() external onlyOwner {
        randomBeacon.updateRewardParameters(
            governanceBeaconParams.getNewDkgResultSubmissionReward(),
            randomBeacon.sortitionPoolUnlockingReward(),
            randomBeacon.ineligibleOperatorNotifierReward(),
            randomBeacon.sortitionPoolRewardsBanDuration(),
            randomBeacon.relayEntryTimeoutNotificationRewardMultiplier(),
            randomBeacon.unauthorizedSigningNotificationRewardMultiplier(),
            randomBeacon.dkgMaliciousResultNotificationRewardMultiplier()
        );

        governanceBeaconParams.finalizeDkgResultSubmissionRewardUpdate();
    }

    /// @notice Begins the sortition pool unlocking reward update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newSortitionPoolUnlockingReward New sortition pool unlocking reward
    function beginSortitionPoolUnlockingRewardUpdate(
        uint256 _newSortitionPoolUnlockingReward
    ) external onlyOwner {
        governanceBeaconParams.beginSortitionPoolUnlockingRewardUpdate(
            _newSortitionPoolUnlockingReward
        );
    }

    /// @notice Finalizes the sortition pool unlocking reward update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeSortitionPoolUnlockingRewardUpdate() external onlyOwner {
        randomBeacon.updateRewardParameters(
            randomBeacon.dkgResultSubmissionReward(),
            governanceBeaconParams.getNewSortitionPoolUnlockingReward(),
            randomBeacon.ineligibleOperatorNotifierReward(),
            randomBeacon.sortitionPoolRewardsBanDuration(),
            randomBeacon.relayEntryTimeoutNotificationRewardMultiplier(),
            randomBeacon.unauthorizedSigningNotificationRewardMultiplier(),
            randomBeacon.dkgMaliciousResultNotificationRewardMultiplier()
        );

        governanceBeaconParams.finalizeSortitionPoolUnlockingRewardUpdate();
    }

    /// @notice Begins the ineligible operator notifier reward update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newIneligibleOperatorNotifierReward New ineligible operator
    ///        notifier reward.
    function beginIneligibleOperatorNotifierRewardUpdate(
        uint256 _newIneligibleOperatorNotifierReward
    ) external onlyOwner {
        governanceBeaconParams.beginIneligibleOperatorNotifierRewardUpdate(
            _newIneligibleOperatorNotifierReward
        );
    }

    /// @notice Finalizes the ineligible operator notifier reward update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeIneligibleOperatorNotifierRewardUpdate()
        external
        onlyOwner
    {
        randomBeacon.updateRewardParameters(
            randomBeacon.dkgResultSubmissionReward(),
            randomBeacon.sortitionPoolUnlockingReward(),
            governanceBeaconParams.getNewIneligibleOperatorNotifierReward(),
            randomBeacon.sortitionPoolRewardsBanDuration(),
            randomBeacon.relayEntryTimeoutNotificationRewardMultiplier(),
            randomBeacon.unauthorizedSigningNotificationRewardMultiplier(),
            randomBeacon.dkgMaliciousResultNotificationRewardMultiplier()
        );

        governanceBeaconParams.finalizeIneligibleOperatorNotifierRewardUpdate();
    }

    /// @notice Begins the sortition pool rewards ban duration update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newSortitionPoolRewardsBanDuration New sortition pool rewards
    ///        ban duration.
    function beginSortitionPoolRewardsBanDurationUpdate(
        uint256 _newSortitionPoolRewardsBanDuration
    ) external onlyOwner {
        governanceBeaconParams.beginSortitionPoolRewardsBanDurationUpdate(
            _newSortitionPoolRewardsBanDuration
        );
    }

    /// @notice Finalizes the sortition pool rewards ban duration update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeSortitionPoolRewardsBanDurationUpdate()
        external
        onlyOwner
    {
        randomBeacon.updateRewardParameters(
            randomBeacon.dkgResultSubmissionReward(),
            randomBeacon.sortitionPoolUnlockingReward(),
            randomBeacon.ineligibleOperatorNotifierReward(),
            governanceBeaconParams.getNewSortitionPoolRewardsBanDuration(),
            randomBeacon.relayEntryTimeoutNotificationRewardMultiplier(),
            randomBeacon.unauthorizedSigningNotificationRewardMultiplier(),
            randomBeacon.dkgMaliciousResultNotificationRewardMultiplier()
        );

        governanceBeaconParams.finalizeSortitionPoolRewardsBanDurationUpdate();
    }

    /// @notice Begins the unauthorized signing notification reward multiplier
    ///         update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newUnauthorizedSigningNotificationRewardMultiplier New unauthorized
    ///         signing notification reward multiplier.
    function beginUnauthorizedSigningNotificationRewardMultiplierUpdate(
        uint256 _newUnauthorizedSigningNotificationRewardMultiplier
    ) external onlyOwner {
        governanceBeaconParams
            .beginUnauthorizedSigningNotificationRewardMultiplierUpdate(
                _newUnauthorizedSigningNotificationRewardMultiplier
            );
    }

    /// @notice Finalizes the unauthorized signing notification reward
    ///         multiplier update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeUnauthorizedSigningNotificationRewardMultiplierUpdate()
        external
        onlyOwner
    {
        randomBeacon.updateRewardParameters(
            randomBeacon.dkgResultSubmissionReward(),
            randomBeacon.sortitionPoolUnlockingReward(),
            randomBeacon.ineligibleOperatorNotifierReward(),
            randomBeacon.sortitionPoolRewardsBanDuration(),
            randomBeacon.relayEntryTimeoutNotificationRewardMultiplier(),
            governanceBeaconParams
                .getNewUnauthorizedSigningNotificationRewardMultiplier(),
            randomBeacon.dkgMaliciousResultNotificationRewardMultiplier()
        );

        governanceBeaconParams
            .finalizeUnauthorizedSigningNotificationRewardMultiplierUpdate();
    }

    /// @notice Begins the relay entry timeout notification reward multiplier
    ///         update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntryTimeoutNotificationRewardMultiplier New relay
    ///        entry timeout notification reward multiplier.
    function beginRelayEntryTimeoutNotificationRewardMultiplierUpdate(
        uint256 _newRelayEntryTimeoutNotificationRewardMultiplier
    ) external onlyOwner {
        governanceBeaconParams
            .beginRelayEntryTimeoutNotificationRewardMultiplierUpdate(
                _newRelayEntryTimeoutNotificationRewardMultiplier
            );
    }

    /// @notice Finalizes the relay entry timeout notification reward
    ///         multiplier update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayEntryTimeoutNotificationRewardMultiplierUpdate()
        external
        onlyOwner
    {
        randomBeacon.updateRewardParameters(
            randomBeacon.dkgResultSubmissionReward(),
            randomBeacon.sortitionPoolUnlockingReward(),
            randomBeacon.ineligibleOperatorNotifierReward(),
            randomBeacon.sortitionPoolRewardsBanDuration(),
            governanceBeaconParams
                .getNewRelayEntryTimeoutNotificationRewardMultiplier(),
            randomBeacon.unauthorizedSigningNotificationRewardMultiplier(),
            randomBeacon.dkgMaliciousResultNotificationRewardMultiplier()
        );

        governanceBeaconParams
            .finalizeRelayEntryTimeoutNotificationRewardMultiplierUpdate();
    }

    // Tu zaczynaj
    // ok
    /// @notice Begins the DKG malicious result notification reward multiplier
    ///         update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgMaliciousResultNotificationRewardMultiplier New DKG
    ///        malicious result notification reward multiplier.
    function beginDkgMaliciousResultNotificationRewardMultiplierUpdate(
        uint256 _newDkgMaliciousResultNotificationRewardMultiplier
    ) external onlyOwner {
        governanceBeaconParams
            .beginDkgMaliciousResultNotificationRewardMultiplierUpdate(
                _newDkgMaliciousResultNotificationRewardMultiplier
            );
    }

    // ok
    /// @notice Finalizes the DKG malicious result notification reward
    ///         multiplier update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgMaliciousResultNotificationRewardMultiplierUpdate()
        external
        onlyOwner
    {
        randomBeacon.updateRewardParameters(
            randomBeacon.dkgResultSubmissionReward(),
            randomBeacon.sortitionPoolUnlockingReward(),
            randomBeacon.ineligibleOperatorNotifierReward(),
            randomBeacon.sortitionPoolRewardsBanDuration(),
            randomBeacon.relayEntryTimeoutNotificationRewardMultiplier(),
            randomBeacon.unauthorizedSigningNotificationRewardMultiplier(),
            governanceBeaconParams
                .getNewDkgMaliciousResultNotificationRewardMultiplier()
        );

        governanceBeaconParams
            .finalizeDkgMaliciousResultNotificationRewardMultiplierUpdate();
    }

    /// @notice Begins the relay entry submission failure slashing amount update
    ///         process.
    /// @dev Can be called only by the contract owner.
    /// @param _newRelayEntrySubmissionFailureSlashingAmount New relay entry
    ///        submission failure slashing amount
    function beginRelayEntrySubmissionFailureSlashingAmountUpdate(
        uint96 _newRelayEntrySubmissionFailureSlashingAmount
    ) external onlyOwner {
        governanceBeaconParams
            .beginRelayEntrySubmissionFailureSlashingAmountUpdate(
                _newRelayEntrySubmissionFailureSlashingAmount
            );
    }

    /// @notice Finalizes the relay entry submission failure slashing amount
    ///         update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayEntrySubmissionFailureSlashingAmountUpdate()
        external
        onlyOwner
    {
        randomBeacon.updateSlashingParameters(
            governanceBeaconParams
                .getNewRelayEntrySubmissionFailureSlashingAmount(),
            randomBeacon.maliciousDkgResultSlashingAmount(),
            randomBeacon.unauthorizedSigningSlashingAmount()
        );

        governanceBeaconParams
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate();
    }

    /// @notice Begins the malicious DKG result slashing amount update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newMaliciousDkgResultSlashingAmount New malicious DKG result
    ///        slashing amount
    function beginMaliciousDkgResultSlashingAmountUpdate(
        uint256 _newMaliciousDkgResultSlashingAmount
    ) external onlyOwner {
        governanceBeaconParams.beginMaliciousDkgResultSlashingAmountUpdate(
            _newMaliciousDkgResultSlashingAmount
        );
    }

    /// @notice Finalizes the malicious DKG result slashing amount update
    ///         process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeMaliciousDkgResultSlashingAmountUpdate()
        external
        onlyOwner
    {
        randomBeacon.updateSlashingParameters(
            randomBeacon.relayEntrySubmissionFailureSlashingAmount(),
            governanceBeaconParams.getNewMaliciousDkgResultSlashingAmount(),
            randomBeacon.unauthorizedSigningSlashingAmount()
        );

        governanceBeaconParams.finalizeMaliciousDkgResultSlashingAmountUpdate();
    }

    /// @notice Begins the unauthorized signing slashing amount update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newUnauthorizedSigningSlashingAmount New unauthorized signing
    ///        slashing amount
    function beginUnauthorizedSigningSlashingAmountUpdate(
        uint256 _newUnauthorizedSigningSlashingAmount
    ) external onlyOwner {
        governanceBeaconParams.beginUnauthorizedSigningSlashingAmountUpdate(
            _newUnauthorizedSigningSlashingAmount
        );
    }

    /// @notice Finalizes the unauthorized signing slashing amount update
    ///         process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeUnauthorizedSigningSlashingAmountUpdate()
        external
        onlyOwner
    {
        randomBeacon.updateSlashingParameters(
            randomBeacon.relayEntrySubmissionFailureSlashingAmount(),
            randomBeacon.maliciousDkgResultSlashingAmount(),
            governanceBeaconParams.getNewUnauthorizedSigningSlashingAmount()
        );

        governanceBeaconParams.finalizeUnauthorizedSigningSlashingAmountUpdate();
    }

    /// @notice Begins the minimum authorization amount update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newMinimumAuthorization New minimum authorization amount.
    function beginMinimumAuthorizationUpdate(uint96 _newMinimumAuthorization)
        external
        onlyOwner
    {
        governanceBeaconParams.beginMinimumAuthorizationUpdate(
            _newMinimumAuthorization
        );
    }

    /// @notice Finalizes the minimum authorization amount update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeMinimumAuthorizationUpdate() external onlyOwner {
        randomBeacon.updateAuthorizationParameters(
            governanceBeaconParams.getNewMinimumAuthorization(),
            randomBeacon.authorizationDecreaseDelay()
        );
        governanceBeaconParams.finalizeMinimumAuthorizationUpdate();
    }

    /// @notice Begins the authorization decrease delay update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newAuthorizationDecreaseDelay New authorization decrease delay
    function beginAuthorizationDecreaseDelayUpdate(
        uint64 _newAuthorizationDecreaseDelay
    ) external onlyOwner {
        governanceBeaconParams.beginAuthorizationDecreaseDelayUpdate(
            _newAuthorizationDecreaseDelay
        );
    }

    /// @notice Finalizes the authorization decrease delay update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeAuthorizationDecreaseDelayUpdate() external onlyOwner {
        randomBeacon.updateAuthorizationParameters(
            randomBeacon.minimumAuthorization(),
            governanceBeaconParams.getNewAuthorizationDecreaseDelay()
        );

        governanceBeaconParams.finalizeAuthorizationDecreaseDelayUpdate();
    }

    /// @notice Withdraws rewards belonging to operators marked as ineligible
    ///         for sortition pool rewards.
    /// @dev Can be called only by the contract owner.
    /// @param recipient Recipient of withdrawn rewards.
    function withdrawIneligibleRewards(address recipient) external onlyOwner {
        randomBeacon.withdrawIneligibleRewards(recipient);
    }

    /// @notice Get the time remaining until the governance delay can be updated.
    /// @return Remaining time in seconds.
    function getRemainingGovernanceDelayUpdateTime()
        external
        view
        returns (uint256)
    {
        return governanceBeaconParams.getRemainingGovernanceDelayUpdateTime();
    }

    /// @notice Get the time remaining until the random beacon ownership can
    ///         be transferred.
    /// @return Remaining time in seconds.
    function getRemainingRandomBeaconOwnershipTransferDelayTime()
        external
        view
        returns (uint256)
    {
        return
            governanceBeaconParams
                .getRemainingRandomBeaconOwnershipTransferDelayTime();
    }

    /// @notice Get the time remaining until the relay request fee can be
    ///         updated.
    /// @return Remaining time in seconds.
    function getRemainingRelayRequestFeeUpdateTime()
        external
        view
        returns (uint256)
    {
        return governanceBeaconParams.getRemainingRelayRequestFeeUpdateTime();
    }

    /// @notice Get the time remaining until the relay entry submission soft
    ///         timeout can be updated.
    /// @return Remaining time in seconds.
    function getRemainingRelayEntrySoftTimeoutUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            governanceBeaconParams
                .getRemainingRelayEntrySoftTimeoutUpdateTime();
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
            governanceBeaconParams
                .getRemainingRelayEntryHardTimeoutUpdateTime();
    }

    /// @notice Get the time remaining until the callback gas limit can be
    ///         updated.
    /// @return Remaining time in seconds.
    function getRemainingCallbackGasLimitUpdateTime()
        external
        view
        returns (uint256)
    {
        return governanceBeaconParams.getRemainingCallbackGasLimitUpdateTime();
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
            governanceBeaconParams
                .getRemainingGroupCreationFrequencyUpdateTime();
    }

    /// @notice Get the time remaining until the group lifetime can be updated.
    /// @return Remaining time in seconds.
    function getRemainingGroupLifetimeUpdateTime()
        external
        view
        returns (uint256)
    {
        return governanceBeaconParams.getRemainingGroupLifetimeUpdateTime();
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
            governanceBeaconParams
                .getRemainingDkgResultChallengePeriodLengthUpdateTime();
    }

    /// @notice Get the time remaining until the DKG result submission timeout
    ///         can be updated.
    /// @return Remaining time in seconds.
    function getRemainingDkgResultSubmissionTimeoutUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            governanceBeaconParams
                .getRemainingDkgResultSubmissionTimeoutUpdateTime();
    }

    /// @notice Get the time remaining until the wallet owner can be updated.
    /// @return Remaining time in seconds.
    function getRemainingDkgSubmitterPrecedencePeriodLengthUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            governanceBeaconParams
                .getRemainingDkgSubmitterPrecedencePeriodLengthUpdateTime();
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
            governanceBeaconParams
                .getRemainingDkgResultSubmissionRewardUpdateTime();
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
            governanceBeaconParams
                .getRemainingSortitionPoolUnlockingRewardUpdateTime();
    }

    /// @notice Get the time remaining until the ineligible operator notifier
    ///         reward can be updated.
    /// @return Remaining time in seconds.
    function getRemainingIneligibleOperatorNotifierRewardUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            governanceBeaconParams
                .getRemainingIneligibleOperatorNotifierRewardUpdateTime();
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
            governanceBeaconParams
                .getRemainingSortitionPoolRewardsBanDurationUpdateTime();
    }

    /// @notice Get the time remaining until the unauthorized signing
    ///         notification reward multiplier duration can be updated.
    /// @return Remaining time in seconds.
    function getRemainingUnauthorizedSigningNotificationRewardMultiplierUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            governanceBeaconParams
                .getRemainingUnauthorizedSigningNotificationRewardMultiplierUpdateTime();
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
            governanceBeaconParams
                .getRemainingRelayEntryTimeoutNotificationRewardMultiplierUpdateTime();
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
            governanceBeaconParams
                .getRemainingDkgMaliciousResultNotificationRewardMultiplierUpdateTime();
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
            governanceBeaconParams
                .getRemainingRelayEntrySubmissionFailureSlashingAmountUpdateTime();
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
            governanceBeaconParams
                .getRemainingMaliciousDkgResultSlashingAmountUpdateTime();
    }

    /// @notice Get the time remaining until the unauthorized signing
    ///         slashing amount can be updated.
    /// @return Remaining time in seconds.
    function getRemainingUnauthorizedSigningSlashingAmountUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            governanceBeaconParams
                .getRemainingUnauthorizedSigningSlashingAmountUpdateTime();
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
            governanceBeaconParams.getRemainingMimimumAuthorizationUpdateTime();
    }

    function getRemainingAuthorizationDecreaseDelayUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            governanceBeaconParams
                .getRemainingAuthorizationDecreaseDelayUpdateTime();
    }
}
