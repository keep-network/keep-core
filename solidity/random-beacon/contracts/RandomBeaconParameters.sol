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

/// @title Governable Parameters
/// @notice Contains Random Beacon governable parameters and logic for their
///         update
/// @dev The client contract should take care of authorizations according to
///      their needs.
library GovernableParameters {
    struct Storage {
        uint96 relayRequestFee;
        uint96 newRelayRequestFee;
        uint96 relayEntrySubmissionFailureSlashingAmount;
        uint96 newRelayEntrySubmissionFailureSlashingAmount;
        uint32 relayEntrySubmissionEligibilityDelay;
        uint32 newRelayEntrySubmissionEligibilityDelay;
        uint32 relayEntryHardTimeout;
        uint32 newRelayEntryHardTimeout;
        uint96 dkgResultSubmissionReward;
        uint96 newDkgResultSubmissionReward;
        uint96 maliciousDkgResultSlashingAmount;
        uint96 newMaliciousDkgResultSlashingAmount;
        uint32 dkgSubmissionEligibilityDelay;
        uint32 newDkgSubmissionEligibilityDelay;
        uint32 dkgResultChallengePeriodLength;
        uint32 newDkgResultChallengePeriodLength;
        uint96 sortitionPoolUnlockingReward;
        uint96 newSortitionPoolUnlockingReward;
        uint32 groupCreationFrequency;
        uint32 newGroupCreationFrequency;
        uint32 groupLifetime;
        uint32 newGroupLifetime;
        uint32 callbackGasLimit;
        uint32 newCallbackGasLimit;
        mapping(string => uint256) changeInitiatedTimestamp;
    }

    event RelayRequestFeeUpdateStarted(
        uint256 relayRequestFee,
        uint256 timestamp
    );
    event RelayRequestFeeUpdated(uint256 relayRequestFee);

    event RelayEntrySubmissionFailureSlashingAmountUpdateStarted(
        uint256 relayEntrySubmissionFailureSlashingAmount,
        uint256 timestamp
    );
    event RelayEntrySubmissionFailureSlashingAmountUpdated(
        uint256 relayEntrySubmissionFailureSlashingAmount
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

    event DkgResultSubmissionRewardUpdateStarted(
        uint256 dkgResultSubmissionReward,
        uint256 timestamp
    );
    event DkgResultSubmissionRewardUpdated(uint256 dkgResultSubmissionReward);

    event MaliciousDkgResultSlashingAmountUpdateStarted(
        uint256 maliciousDkgResultSlashingAmount,
        uint256 timestamp
    );
    event MaliciousDkgResultSlashingAmountUpdated(
        uint256 maliciousDkgResultSlashingAmount
    );

    event DkgSubmissionEligibilityDelayUpdateStarted(
        uint256 dkgSubmissionEligibilityDelay,
        uint256 timestamp
    );
    event DkgSubmissionEligibilityDelayUpdated(
        uint256 dkgSubmissionEligibilityDelay
    );

    event DkgResultChallengePeriodLengthUpdateStarted(
        uint256 dkgResultChallengePeriodLength,
        uint256 timestamp
    );
    event DkgResultChallengePeriodLengthUpdated(
        uint256 dkgResultChallengePeriodLength
    );

    event SortitionPoolUnlockingRewardUpdateStarted(
        uint256 sortitionPoolUnlockingReward,
        uint256 timestamp
    );
    event SortitionPoolUnlockingRewardUpdated(
        uint256 sortitionPoolUnlockingReward
    );

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
    /// @param parameter Name of the parameter whose change initiation timestamp
    ///        is will be checked
    /// @param delay Delay that is required to elapse before the parameter value
    ///        can be updated
    modifier onlyAfterGovernanceDelay(
        Storage storage self,
        string memory parameter,
        uint256 delay
    ) {
        /* solhint-disable not-rely-on-time */
        uint256 initiated = self.changeInitiatedTimestamp[parameter];
        require(initiated > 0, "Change not initiated");
        require(
            block.timestamp - initiated >= delay,
            "Governance delay has not elapsed"
        );
        _;
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Begins the relay request fee update process.
    /// @param _newRelayRequestFee New relay request fee
    function beginRelayRequestFeeUpdate(
        Storage storage self,
        uint96 _newRelayRequestFee
    ) external {
        /* solhint-disable not-rely-on-time */
        self.newRelayRequestFee = _newRelayRequestFee;
        self.changeInitiatedTimestamp["relayRequestFee"] = block.timestamp;
        emit RelayRequestFeeUpdateStarted(_newRelayRequestFee, block.timestamp);
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the relay request fee update process.
    /// @dev Can be called only after the governance delay elapses.
    function finalizeRelayRequestFeeUpdate(Storage storage self)
        external
        onlyAfterGovernanceDelay(self, "relayRequestFee", 24 hours)
    {
        self.relayRequestFee = self.newRelayRequestFee;
        emit RelayRequestFeeUpdated(self.relayRequestFee);
        self.changeInitiatedTimestamp["relayRequestFee"] = 0;
        self.newRelayRequestFee = 0;
    }

    /// @notice Begins the relay entry submission failure slashing amount update
    ///         process.
    /// @param _newRelayEntrySubmissionFailureSlashingAmount New relay entry
    ///        submission failure slashing amount
    function beginRelayEntrySubmissionFailureSlashingAmountUpdate(
        Storage storage self,
        uint96 _newRelayEntrySubmissionFailureSlashingAmount
    ) external {
        /* solhint-disable not-rely-on-time */
        self
            .newRelayEntrySubmissionFailureSlashingAmount = _newRelayEntrySubmissionFailureSlashingAmount;
        self.changeInitiatedTimestamp[
            "relayEntrySubmissionFailureSlashingAmount"
        ] = block.timestamp;
        emit RelayEntrySubmissionFailureSlashingAmountUpdateStarted(
            _newRelayEntrySubmissionFailureSlashingAmount,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the relay entry submission failure slashing amount
    ///         update process.
    /// @dev Can be called only after the governance delay elapses.
    function finalizeRelayEntrySubmissionFailureSlashingAmountUpdate(
        Storage storage self
    )
        external
        onlyAfterGovernanceDelay(
            self,
            "relayEntrySubmissionFailureSlashingAmount",
            2 weeks
        )
    {
        self.relayEntrySubmissionFailureSlashingAmount = self
            .newRelayEntrySubmissionFailureSlashingAmount;
        emit RelayEntrySubmissionFailureSlashingAmountUpdated(
            self.relayEntrySubmissionFailureSlashingAmount
        );
        self.changeInitiatedTimestamp[
            "relayEntrySubmissionFailureSlashingAmount"
        ] = 0;
        self.newRelayEntrySubmissionFailureSlashingAmount = 0;
    }

    /// @notice Begins the relay entry submission eligibility delay update
    ///         process.
    /// @param _newRelayEntrySubmissionEligibilityDelay New relay entry
    ///        submission eligibility delay in blocks
    function beginRelayEntrySubmissionEligibilityDelayUpdate(
        Storage storage self,
        uint32 _newRelayEntrySubmissionEligibilityDelay
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newRelayEntrySubmissionEligibilityDelay > 0,
            "Relay entry submission eligibility delay must be greater than 0 blocks"
        );
        self
            .newRelayEntrySubmissionEligibilityDelay = _newRelayEntrySubmissionEligibilityDelay;
        self.changeInitiatedTimestamp[
            "relayEntrySubmissionEligibilityDelay"
        ] = block.timestamp;
        emit RelayEntrySubmissionEligibilityDelayUpdateStarted(
            _newRelayEntrySubmissionEligibilityDelay,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the relay entry submission eligibility delay update
    ////        process.
    /// @dev Can be called only after the governance delay elapses.
    function finalizeRelayEntrySubmissionEligibilityDelayUpdate(
        Storage storage self
    )
        external
        onlyAfterGovernanceDelay(
            self,
            "relayEntrySubmissionEligibilityDelay",
            24 hours
        )
    {
        self.relayEntrySubmissionEligibilityDelay = self
            .newRelayEntrySubmissionEligibilityDelay;
        emit RelayEntrySubmissionEligibilityDelayUpdated(
            self.relayEntrySubmissionEligibilityDelay
        );
        self.changeInitiatedTimestamp[
            "relayEntrySubmissionEligibilityDelay"
        ] = 0;
        self.newRelayEntrySubmissionEligibilityDelay = 0;
    }

    /// @notice Begins the relay entry hard timeout update process.
    /// @param _newRelayEntryHardTimeout New relay entry hard timeout in blocks
    function beginRelayEntryHardTimeoutUpdate(
        Storage storage self,
        uint32 _newRelayEntryHardTimeout
    ) external {
        /* solhint-disable not-rely-on-time */
        self.newRelayEntryHardTimeout = _newRelayEntryHardTimeout;
        self.changeInitiatedTimestamp["relayEntryHardTimeout"] = block
            .timestamp;
        emit RelayEntryHardTimeoutUpdateStarted(
            _newRelayEntryHardTimeout,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the relay entry hard timeout update process.
    /// @dev Can be called only after the governance delay elapses.
    function finalizeRelayEntryHardTimeoutUpdate(Storage storage self)
        external
        onlyAfterGovernanceDelay(self, "relayEntryHardTimeout", 2 weeks)
    {
        self.relayEntryHardTimeout = self.newRelayEntryHardTimeout;
        emit RelayEntryHardTimeoutUpdated(self.relayEntryHardTimeout);
        self.changeInitiatedTimestamp["relayEntryHardTimeout"] = 0;
        self.newRelayEntryHardTimeout = 0;
    }

    /// @notice Begins the DKG result submission reward update process.
    /// @param _newDkgResultSubmissionReward New DKG result submission reward
    function beginDkgResultSubmissionRewardUpdate(
        Storage storage self,
        uint96 _newDkgResultSubmissionReward
    ) external {
        /* solhint-disable not-rely-on-time */
        self.newDkgResultSubmissionReward = _newDkgResultSubmissionReward;
        self.changeInitiatedTimestamp["dkgResultSubmissionReward"] = block
            .timestamp;
        emit DkgResultSubmissionRewardUpdateStarted(
            _newDkgResultSubmissionReward,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG result submission reward update process.
    /// @dev Can be called only after the governance delay elapses.
    function finalizeDkgResultSubmissionRewardUpdate(Storage storage self)
        external
        onlyAfterGovernanceDelay(self, "dkgResultSubmissionReward", 24 hours)
    {
        self.dkgResultSubmissionReward = self.newDkgResultSubmissionReward;
        emit DkgResultSubmissionRewardUpdated(self.dkgResultSubmissionReward);
        self.changeInitiatedTimestamp["dkgResultSubmissionReward"] = 0;
        self.newDkgResultSubmissionReward = 0;
    }

    /// @notice Begins the malicious DKG result slashing amount update process.
    /// @param _newMaliciousDkgResultSlashingAmount New malicious DKG result
    ///        slashing amount
    function beginMaliciousDkgResultSlashingAmountUpdate(
        Storage storage self,
        uint96 _newMaliciousDkgResultSlashingAmount
    ) external {
        /* solhint-disable not-rely-on-time */
        self
            .newMaliciousDkgResultSlashingAmount = _newMaliciousDkgResultSlashingAmount;
        self.changeInitiatedTimestamp[
            "maliciousDkgResultSlashingAmount"
        ] = block.timestamp;
        emit MaliciousDkgResultSlashingAmountUpdateStarted(
            _newMaliciousDkgResultSlashingAmount,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the malicious DKG result slashing amount update
    ///         process.
    /// @dev Can be called only after the governance delay elapses.
    function finalizeMaliciousDkgResultSlashingAmountUpdate(
        Storage storage self
    )
        external
        onlyAfterGovernanceDelay(
            self,
            "maliciousDkgResultSlashingAmount",
            24 hours
        )
    {
        self.maliciousDkgResultSlashingAmount = self
            .newMaliciousDkgResultSlashingAmount;
        emit MaliciousDkgResultSlashingAmountUpdated(
            self.maliciousDkgResultSlashingAmount
        );
        self.changeInitiatedTimestamp["maliciousDkgResultSlashingAmount"] = 0;
        self.newMaliciousDkgResultSlashingAmount = 0;
    }

    /// @notice Begins the DKG submission eligibility delay update process.
    /// @param _newDkgSubmissionEligibilityDelay New DKG submission eligibility
    ///        delay in blocks
    function beginDkgSubmissionEligibilityDelayUpdate(
        Storage storage self,
        uint32 _newDkgSubmissionEligibilityDelay
    ) external {
        require(
            _newDkgSubmissionEligibilityDelay > 0,
            "DKG submission eligibility delay must be greater than 0 blocks"
        );
        /* solhint-disable not-rely-on-time */
        self
            .newDkgSubmissionEligibilityDelay = _newDkgSubmissionEligibilityDelay;
        self.changeInitiatedTimestamp["dkgSubmissionEligibilityDelay"] = block
            .timestamp;
        emit DkgSubmissionEligibilityDelayUpdateStarted(
            _newDkgSubmissionEligibilityDelay,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG submission eligibility delay update process.
    /// @dev Can be called only after the governance delay elapses.
    function finalizeDkgSubmissionEligibilityDelayUpdate(Storage storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            "dkgSubmissionEligibilityDelay",
            24 hours
        )
    {
        self.dkgSubmissionEligibilityDelay = self
            .newDkgSubmissionEligibilityDelay;
        emit DkgSubmissionEligibilityDelayUpdated(
            self.dkgSubmissionEligibilityDelay
        );
        self.changeInitiatedTimestamp["dkgSubmissionEligibilityDelay"] = 0;
        self.newDkgSubmissionEligibilityDelay = 0;
    }

    /// @notice Begins the DKG result challenge period length update process.
    /// @param _newDkgResultChallengePeriodLength New DKG result challenge
    ///        period length in blocks
    function beginDkgResultChallengePeriodLengthUpdate(
        Storage storage self,
        uint32 _newDkgResultChallengePeriodLength
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newDkgResultChallengePeriodLength > 10,
            "DKG result challenge period length must be grater than 10 blocks"
        );
        self
            .newDkgResultChallengePeriodLength = _newDkgResultChallengePeriodLength;
        self.changeInitiatedTimestamp["dkgResultChallengePeriodLength"] = block
            .timestamp;
        emit DkgResultChallengePeriodLengthUpdateStarted(
            _newDkgResultChallengePeriodLength,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG result challenge period length update process.
    /// @dev Can be called only after the governance delay elapses.
    function finalizeDkgResultChallengePeriodLengthUpdate(Storage storage self)
        external
        onlyAfterGovernanceDelay(
            self,
            "dkgResultChallengePeriodLength",
            24 hours
        )
    {
        self.dkgResultChallengePeriodLength = self
            .newDkgResultChallengePeriodLength;
        emit DkgResultChallengePeriodLengthUpdated(
            self.dkgResultChallengePeriodLength
        );
        self.changeInitiatedTimestamp["dkgResultChallengePeriodLength"] = 0;
        self.newDkgResultChallengePeriodLength = 0;
    }

    /// @notice Begins the sortition pool unlocking reward update process.
    /// @param _newSortitionPoolUnlockingReward New sortition pool unlocking reward
    function beginSortitionPoolUnlockingRewardUpdate(
        Storage storage self,
        uint96 _newSortitionPoolUnlockingReward
    ) external {
        /* solhint-disable not-rely-on-time */
        self.newSortitionPoolUnlockingReward = _newSortitionPoolUnlockingReward;
        self.changeInitiatedTimestamp["sortitionPoolUnlockingReward"] = block
            .timestamp;
        emit SortitionPoolUnlockingRewardUpdateStarted(
            _newSortitionPoolUnlockingReward,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the sortition pool unlocking reward update process.
    /// @dev Can be called only after the governance delay elapses.
    function finalizeSortitionPoolUnlockingRewardUpdate(Storage storage self)
        external
        onlyAfterGovernanceDelay(self, "sortitionPoolUnlockingReward", 24 hours)
    {
        self.sortitionPoolUnlockingReward = self
            .newSortitionPoolUnlockingReward;
        emit SortitionPoolUnlockingRewardUpdated(
            self.sortitionPoolUnlockingReward
        );
        self.changeInitiatedTimestamp["sortitionPoolUnlockingReward"] = 0;
        self.newSortitionPoolUnlockingReward = 0;
    }

    /// @notice Begins the group creation frequency update process.
    /// @param _newGroupCreationFrequency New group creation frequency
    function beginGroupCreationFrequencyUpdate(
        Storage storage self,
        uint32 _newGroupCreationFrequency
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newGroupCreationFrequency > 0,
            "Group creation frequency must be grater than zero"
        );
        self.newGroupCreationFrequency = _newGroupCreationFrequency;
        self.changeInitiatedTimestamp["groupCreationFrequency"] = block
            .timestamp;
        emit GroupCreationFrequencyUpdateStarted(
            _newGroupCreationFrequency,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the group creation frequency update process.
    /// @dev Can be called only after the governance delay elapses.
    function finalizeGroupCreationFrequencyUpdate(Storage storage self)
        external
        onlyAfterGovernanceDelay(self, "groupCreationFrequency", 2 weeks)
    {
        self.groupCreationFrequency = self.newGroupCreationFrequency;
        emit GroupCreationFrequencyUpdated(self.groupCreationFrequency);
        self.changeInitiatedTimestamp["groupCreationFrequency"] = 0;
        self.newGroupCreationFrequency = 0;
    }

    /// @notice Begins the group lifetime update process.
    /// @param _newGroupLifetime New group lifetime in seconds
    function beginGroupLifetimeUpdate(
        Storage storage self,
        uint32 _newGroupLifetime
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newGroupLifetime >= 1 days && _newGroupLifetime <= 2 weeks,
            "Group lifetime must be >= 1 day and <= 2 weeks"
        );
        self.newGroupLifetime = _newGroupLifetime;
        self.changeInitiatedTimestamp["groupLifetime"] = block.timestamp;
        emit GroupLifetimeUpdateStarted(_newGroupLifetime, block.timestamp);
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the group creation frequency update process.
    /// @dev Can be called only after the governance delay elapses.
    function finalizeGroupLifetimeUpdate(Storage storage self)
        external
        onlyAfterGovernanceDelay(self, "groupLifetime", 2 weeks)
    {
        self.groupLifetime = self.newGroupLifetime;
        emit GroupLifetimeUpdated(self.groupLifetime);
        self.changeInitiatedTimestamp["groupLifetime"] = 0;
        self.newGroupLifetime = 0;
    }

    /// @notice Begins the callback gas limit update process.
    /// @param _newCallbackGasLimit New callback gas limit
    function beginCallbackGasLimitUpdate(
        Storage storage self,
        uint32 _newCallbackGasLimit
    ) external {
        /* solhint-disable not-rely-on-time */
        require(
            _newCallbackGasLimit > 0 && _newCallbackGasLimit < 1000000,
            "Callback gas limit must be > 0 and < 1000000"
        );
        self.newCallbackGasLimit = _newCallbackGasLimit;
        self.changeInitiatedTimestamp["callbackGasLimit"] = block.timestamp;
        emit CallbackGasLimitUpdateStarted(
            _newCallbackGasLimit,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the callback gas limit update process.
    /// @dev Can be called only after the governance delay elapses.
    function finalizeCallbackGasLimitUpdate(Storage storage self)
        external
        onlyAfterGovernanceDelay(self, "callbackGasLimit", 2 weeks)
    {
        self.callbackGasLimit = self.newCallbackGasLimit;
        emit CallbackGasLimitUpdated(self.callbackGasLimit);
        self.changeInitiatedTimestamp["callbackGasLimit"] = 0;
        self.newCallbackGasLimit = 0;
    }
}

/// @title Random Beacon Parameters
/// @notice This is a base contract that holds governable parameters for Random
///         Beacon as well as their update logic
contract RandomBeaconParameters is Ownable {
    using GovernableParameters for GovernableParameters.Storage;
    GovernableParameters.Storage internal governableParameters;

    event RelayRequestFeeUpdateStarted(
        uint256 relayRequestFee,
        uint256 timestamp
    );
    event RelayRequestFeeUpdated(uint256 relayRequestFee);

    event RelayEntrySubmissionFailureSlashingAmountUpdateStarted(
        uint256 relayEntrySubmissionFailureSlashingAmount,
        uint256 timestamp
    );
    event RelayEntrySubmissionFailureSlashingAmountUpdated(
        uint256 relayEntrySubmissionFailureSlashingAmount
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

    event DkgResultSubmissionRewardUpdateStarted(
        uint256 dkgResultSubmissionReward,
        uint256 timestamp
    );
    event DkgResultSubmissionRewardUpdated(uint256 dkgResultSubmissionReward);

    event MaliciousDkgResultSlashingAmountUpdateStarted(
        uint256 maliciousDkgResultSlashingAmount,
        uint256 timestamp
    );
    event MaliciousDkgResultSlashingAmountUpdated(
        uint256 maliciousDkgResultSlashingAmount
    );

    event DkgSubmissionEligibilityDelayUpdateStarted(
        uint256 dkgSubmissionEligibilityDelay,
        uint256 timestamp
    );
    event DkgSubmissionEligibilityDelayUpdated(
        uint256 dkgSubmissionEligibilityDelay
    );

    event DkgResultChallengePeriodLengthUpdateStarted(
        uint256 dkgResultChallengePeriodLength,
        uint256 timestamp
    );
    event DkgResultChallengePeriodLengthUpdated(
        uint256 dkgResultChallengePeriodLength
    );

    event SortitionPoolUnlockingRewardUpdateStarted(
        uint256 sortitionPoolUnlockingReward,
        uint256 timestamp
    );
    event SortitionPoolUnlockingRewardUpdated(
        uint256 sortitionPoolUnlockingReward
    );

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

    /// @notice Begins the relay request fee update process.
    /// @dev Can be called only by the contract owner.
    /// @param newRelayRequestFee New relay request fee
    function beginRelayRequestFeeUpdate(uint256 newRelayRequestFee)
        external
        onlyOwner
    {
        governableParameters.beginRelayRequestFeeUpdate(
            uint96(newRelayRequestFee)
        );
    }

    /// @notice Finalizes the relay request fee update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayRequestFeeUpdate() external onlyOwner {
        governableParameters.finalizeRelayRequestFeeUpdate();
    }

    /// @notice Begins the relay entry submission failure slashing amount update
    ///         process.
    /// @dev Can be called only by the contract owner.
    /// @param newRelayEntrySubmissionFailureSlashingAmount New relay entry
    ///        submission failure slashing amount
    function beginRelayEntrySubmissionFailureSlashingAmountUpdate(
        uint256 newRelayEntrySubmissionFailureSlashingAmount
    ) external onlyOwner {
        governableParameters
            .beginRelayEntrySubmissionFailureSlashingAmountUpdate(
                uint96(newRelayEntrySubmissionFailureSlashingAmount)
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
        governableParameters
            .finalizeRelayEntrySubmissionFailureSlashingAmountUpdate();
    }

    /// @notice Begins the relay entry submission eligibility delay update
    ///         process.
    /// @dev Can be called only by the contract owner.
    /// @param newRelayEntrySubmissionEligibilityDelay New relay entry
    ///        submission eligibility delay in blocks
    function beginRelayEntrySubmissionEligibilityDelayUpdate(
        uint256 newRelayEntrySubmissionEligibilityDelay
    ) external onlyOwner {
        governableParameters.beginRelayEntrySubmissionEligibilityDelayUpdate(
            uint32(newRelayEntrySubmissionEligibilityDelay)
        );
    }

    /// @notice Finalizes the relay entry submission eligibility delay update
    ////        process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayEntrySubmissionEligibilityDelayUpdate()
        external
        onlyOwner
    {
        governableParameters
            .finalizeRelayEntrySubmissionEligibilityDelayUpdate();
    }

    /// @notice Begins the relay entry hard timeout update process.
    /// @dev Can be called only by the contract owner.
    /// @param newRelayEntryHardTimeout New relay entry hard timeout in blocks
    function beginRelayEntryHardTimeoutUpdate(uint256 newRelayEntryHardTimeout)
        external
        onlyOwner
    {
        governableParameters.beginRelayEntryHardTimeoutUpdate(
            uint32(newRelayEntryHardTimeout)
        );
    }

    /// @notice Finalizes the relay entry hard timeout update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeRelayEntryHardTimeoutUpdate() external onlyOwner {
        governableParameters.finalizeRelayEntryHardTimeoutUpdate();
    }

    /// @notice Begins the DKG result submission reward update process.
    /// @dev Can be called only by the contract owner.
    /// @param newDkgResultSubmissionReward New DKG result submission reward
    function beginDkgResultSubmissionRewardUpdate(
        uint256 newDkgResultSubmissionReward
    ) external onlyOwner {
        governableParameters.beginDkgResultSubmissionRewardUpdate(
            uint96(newDkgResultSubmissionReward)
        );
    }

    /// @notice Finalizes the DKG result submission reward update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgResultSubmissionRewardUpdate() external onlyOwner {
        governableParameters.finalizeDkgResultSubmissionRewardUpdate();
    }

    /// @notice Begins the malicious DKG result slashing amount update process.
    /// @dev Can be called only by the contract owner.
    /// @param newMaliciousDkgResultSlashingAmount New malicious DKG result
    ///        slashing amount
    function beginMaliciousDkgResultSlashingAmountUpdate(
        uint256 newMaliciousDkgResultSlashingAmount
    ) external onlyOwner {
        governableParameters.beginMaliciousDkgResultSlashingAmountUpdate(
            uint96(newMaliciousDkgResultSlashingAmount)
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
        governableParameters.finalizeMaliciousDkgResultSlashingAmountUpdate();
    }

    /// @notice Begins the DKG submission eligibility delay update process.
    /// @dev Can be called only by the contract owner.
    /// @param newDkgSubmissionEligibilityDelay New DKG submission eligibility
    ///        delay in blocks
    function beginDkgSubmissionEligibilityDelayUpdate(
        uint256 newDkgSubmissionEligibilityDelay
    ) external onlyOwner {
        governableParameters.beginDkgSubmissionEligibilityDelayUpdate(
            uint32(newDkgSubmissionEligibilityDelay)
        );
    }

    /// @notice Finalizes the DKG submission eligibility delay update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgSubmissionEligibilityDelayUpdate() external onlyOwner {
        governableParameters.finalizeDkgSubmissionEligibilityDelayUpdate();
    }

    /// @notice Begins the DKG result challenge period length update process.
    /// @dev Can be called only by the contract owner.
    /// @param newDkgResultChallengePeriodLength New DKG result challenge
    ///        period length in blocks
    function beginDkgResultChallengePeriodLengthUpdate(
        uint256 newDkgResultChallengePeriodLength
    ) external onlyOwner {
        governableParameters.beginDkgResultChallengePeriodLengthUpdate(
            uint32(newDkgResultChallengePeriodLength)
        );
    }

    /// @notice Finalizes the DKG result challenge period length update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgResultChallengePeriodLengthUpdate() external onlyOwner {
        governableParameters.finalizeDkgResultChallengePeriodLengthUpdate();
    }

    /// @notice Begins the sortition pool unlocking reward update process.
    /// @dev Can be called only by the contract owner.
    /// @param newSortitionPoolUnlockingReward New sortition pool unlocking reward
    function beginSortitionPoolUnlockingRewardUpdate(
        uint256 newSortitionPoolUnlockingReward
    ) external onlyOwner {
        governableParameters.beginSortitionPoolUnlockingRewardUpdate(
            uint96(newSortitionPoolUnlockingReward)
        );
    }

    /// @notice Finalizes the sortition pool unlocking reward update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeSortitionPoolUnlockingRewardUpdate() external onlyOwner {
        governableParameters.finalizeSortitionPoolUnlockingRewardUpdate();
    }

    /// @notice Begins the group creation frequency update process.
    /// @dev Can be called only by the contract owner.
    /// @param newGroupCreationFrequency New group creation frequency
    function beginGroupCreationFrequencyUpdate(
        uint256 newGroupCreationFrequency
    ) external onlyOwner {
        governableParameters.beginGroupCreationFrequencyUpdate(
            uint32(newGroupCreationFrequency)
        );
    }

    /// @notice Finalizes the group creation frequency update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeGroupCreationFrequencyUpdate() external onlyOwner {
        governableParameters.finalizeGroupCreationFrequencyUpdate();
    }

    /// @notice Begins the group lifetime update process.
    /// @dev Can be called only by the contract owner.
    /// @param newGroupLifetime New group lifetime in seconds
    function beginGroupLifetimeUpdate(uint256 newGroupLifetime)
        external
        onlyOwner
    {
        governableParameters.beginGroupLifetimeUpdate(uint32(newGroupLifetime));
    }

    /// @notice Finalizes the group creation frequency update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeGroupLifetimeUpdate() external onlyOwner {
        governableParameters.finalizeGroupLifetimeUpdate();
    }

    /// @notice Begins the callback gas limit update process.
    /// @dev Can be called only by the contract owner.
    /// @param newCallbackGasLimit New callback gas limit
    function beginCallbackGasLimitUpdate(uint256 newCallbackGasLimit)
        external
        onlyOwner
    {
        governableParameters.beginCallbackGasLimitUpdate(
            uint32(newCallbackGasLimit)
        );
    }

    /// @notice Finalizes the callback gas limit update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeCallbackGasLimitUpdate() external onlyOwner {
        governableParameters.finalizeCallbackGasLimitUpdate();
    }

    /// @notice Returns the relay request fee.
    function relayRequestFee() public view returns (uint256) {
        return uint256(governableParameters.relayRequestFee);
    }

    /// @notice Returns the entry submission failure slashing amount.
    function relayEntrySubmissionFailureSlashingAmount()
        public
        view
        returns (uint256)
    {
        return
            uint256(
                governableParameters.relayEntrySubmissionFailureSlashingAmount
            );
    }

    /// @notice Returns the entry submission eligibility delay in blocks.
    function relayEntrySubmissionEligibilityDelay()
        public
        view
        returns (uint256)
    {
        return
            uint256(governableParameters.relayEntrySubmissionEligibilityDelay);
    }

    /// @notice Returns the relay entry hard timeout in blocks.
    function relayEntryHardTimeout() public view returns (uint256) {
        return uint256(governableParameters.relayEntryHardTimeout);
    }

    /// @notice Returns the DKG result submission reward.
    function dkgResultSubmissionReward() public view returns (uint256) {
        return uint256(governableParameters.dkgResultSubmissionReward);
    }

    /// @notice Returns the malicious DKG result slashing amount.
    function maliciousDkgResultSlashingAmount() public view returns (uint256) {
        return uint256(governableParameters.maliciousDkgResultSlashingAmount);
    }

    /// @notice Returns the DKG submission eligibility delay in blocks.
    function dkgSubmissionEligibilityDelay() public view returns (uint256) {
        return uint256(governableParameters.dkgSubmissionEligibilityDelay);
    }

    /// @notice Returns the DKG result challenge period length in blocks.
    function dkgResultChallengePeriodLength() public view returns (uint256) {
        return uint256(governableParameters.dkgResultChallengePeriodLength);
    }

    /// @notice Returns the sortition pool unlocking reward.
    function sortitionPoolUnlockingReward() public view returns (uint256) {
        return uint256(governableParameters.sortitionPoolUnlockingReward);
    }

    /// @notice Returns the group creation frequency.
    function groupCreationFrequency() public view returns (uint256) {
        return uint256(governableParameters.groupCreationFrequency);
    }

    /// @notice Returns the group lifetime in seconds.
    function groupLifetime() public view returns (uint256) {
        return uint256(governableParameters.groupLifetime);
    }

    /// @notice Returns the callback gas limit.
    function callbackGasLimit() public view returns (uint256) {
        return uint256(governableParameters.callbackGasLimit);
    }
}
