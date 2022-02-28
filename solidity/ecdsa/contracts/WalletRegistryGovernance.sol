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

import "./WalletRegistry.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @title Wallet Registry Governance
/// @notice Owns the `WalletRegistry` contract and is responsible for updating its
///         governable parameters in respect to governance delay individual
///         for each parameter.
contract WalletRegistryGovernance is Ownable {
    address public newWalletOwner;
    uint256 public walletOwnerChangeInitiated;

    uint96 public newMinimumAuthorization;
    uint256 public minimumAuthorizationChangeInitiated;

    uint64 public newAuthorizationDecreaseDelay;
    uint256 public authorizationDecreaseDelayChangeInitiated;

    uint96 public newMaliciousDkgResultSlashingAmount;
    uint256 public maliciousDkgResultSlashingAmountChangeInitiated;

    uint256 public newMaliciousDkgResultNotificationRewardMultiplier;
    uint256
        public maliciousDkgResultNotificationRewardMultiplierChangeInitiated;

    uint256 public newDkgSeedTimeout;
    uint256 public dkgSeedTimeoutChangeInitiated;

    uint256 public newDkgResultChallengePeriodLength;
    uint256 public dkgResultChallengePeriodLengthChangeInitiated;

    uint256 public newDkgResultSubmissionTimeout;
    uint256 public dkgResultSubmissionTimeoutChangeInitiated;

    uint256 public newSubmitterPrecedencePeriodLength;
    uint256 public dkgSubmitterPrecedencePeriodLengthChangeInitiated;

    WalletRegistry public walletRegistry;

    // Long governance delay used for critical parameters giving a chance for
    // stakers to opt out before the change is finalized in case they do not
    // agree with that change. The maximum group lifetime must not be longer
    // than this delay.
    //
    // The full list of parameters protected by this delay:
    // - wallet owner
    // - minimum authorization
    // - authorization decrease delay
    uint256 internal constant CRITICAL_PARAMETER_GOVERNANCE_DELAY = 2 weeks;

    // Short governance delay for non-critical parameters. Honest stakers should
    // not be severely affected by any change of these parameters.
    //
    // The full list of parameters protected by this delay:
    // - malicious DKG result notification reward multiplier
    // - malicious DKG result slashing amount
    // - DKG seed timeout
    // - DKG result challenge period length
    // - DKG result submission eligibility delay
    // - DKG submitter precedence period length
    uint256 internal constant STANDARD_PARAMETER_GOVERNANCE_DELAY = 12 hours;

    event WalletOwnerUpdateStarted(address walletOwner, uint256 timestamp);
    event WalletOwnerUpdated(address walletOwner);

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

    event MaliciousDkgResultSlashingAmountUpdateStarted(
        uint256 maliciousDkgResultSlashingAmount,
        uint256 timestamp
    );
    event MaliciousDkgResultSlashingAmountUpdated(
        uint256 maliciousDkgResultSlashingAmount
    );

    event MaliciousDkgResultNotificationRewardMultiplierUpdateStarted(
        uint256 maliciousDkgResultNotificationRewardMultiplier,
        uint256 timestamp
    );
    event MaliciousDkgResultNotificationRewardMultiplierUpdated(
        uint256 maliciousDkgResultNotificationRewardMultiplier
    );

    event DkgSeedTimeoutUpdateStarted(
        uint256 dkgSeedTimeout,
        uint256 timestamp
    );
    event DkgSeedTimeoutUpdated(uint256 dkgSeedTimeout);

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

    /// @notice Reverts if called before the governance delay elapses.
    /// @param changeInitiatedTimestamp Timestamp indicating the beginning
    ///        of the change.
    /// @param delay Governance delay
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

    constructor(WalletRegistry _walletRegistry) {
        walletRegistry = _walletRegistry;
    }

    /// @notice Upgrades the random beacon.
    /// @dev Can be called only by the contract owner.
    /// @param _newRandomBeacon New random beacon address
    function upgradeRandomBeacon(address _newRandomBeacon) external onlyOwner {
        require(
            _newRandomBeacon != address(0),
            "New random beacon address cannot be zero"
        );

        walletRegistry.upgradeRandomBeacon(_newRandomBeacon);
    }

    /// @notice Begins the wallet owner update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newWalletOwner New wallet owner address
    function beginWalletOwnerUpdate(address _newWalletOwner)
        external
        onlyOwner
    {
        require(
            _newWalletOwner != address(0),
            "New wallet owner address cannot be zero"
        );
        /* solhint-disable not-rely-on-time */
        newWalletOwner = _newWalletOwner;
        walletOwnerChangeInitiated = block.timestamp;
        emit WalletOwnerUpdateStarted(_newWalletOwner, block.timestamp);
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the wallet owner update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeWalletOwnerUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            walletOwnerChangeInitiated,
            CRITICAL_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit WalletOwnerUpdated(newWalletOwner);
        // slither-disable-next-line reentrancy-no-eth
        walletRegistry.updateWalletOwner(newWalletOwner);
        walletOwnerChangeInitiated = 0;
        newWalletOwner = address(0);
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
        walletRegistry.updateAuthorizationParameters(
            newMinimumAuthorization,
            walletRegistry.authorizationDecreaseDelay()
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
        walletRegistry.updateAuthorizationParameters(
            walletRegistry.minimumAuthorization(),
            newAuthorizationDecreaseDelay
        );
        authorizationDecreaseDelayChangeInitiated = 0;
        newAuthorizationDecreaseDelay = 0;
    }

    /// @notice Begins the malicious DKG result slashing amount update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newMaliciousDkgResultSlashingAmount New malicious DKG result
    ///        slashing amount
    function beginMaliciousDkgResultSlashingAmountUpdate(
        uint96 _newMaliciousDkgResultSlashingAmount
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
        walletRegistry.updateSlashingParameters(
            newMaliciousDkgResultSlashingAmount
        );
        maliciousDkgResultSlashingAmountChangeInitiated = 0;
        newMaliciousDkgResultSlashingAmount = 0;
    }

    /// @notice Begins the DKG malicious result notification reward multiplier
    ///         update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newMaliciousDkgResultNotificationRewardMultiplier New DKG
    ///        malicious result notification reward multiplier.
    function beginMaliciousDkgResultNotificationRewardMultiplierUpdate(
        uint256 _newMaliciousDkgResultNotificationRewardMultiplier
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        require(
            _newMaliciousDkgResultNotificationRewardMultiplier <= 100,
            "Maximum value is 100"
        );

        newMaliciousDkgResultNotificationRewardMultiplier = _newMaliciousDkgResultNotificationRewardMultiplier;
        maliciousDkgResultNotificationRewardMultiplierChangeInitiated = block
            .timestamp;
        emit MaliciousDkgResultNotificationRewardMultiplierUpdateStarted(
            _newMaliciousDkgResultNotificationRewardMultiplier,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG malicious result notification reward
    ///         multiplier update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeMaliciousDkgResultNotificationRewardMultiplierUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            maliciousDkgResultNotificationRewardMultiplierChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit MaliciousDkgResultNotificationRewardMultiplierUpdated(
            newMaliciousDkgResultNotificationRewardMultiplier
        );
        // slither-disable-next-line reentrancy-no-eth
        walletRegistry.updateRewardParameters(
            newMaliciousDkgResultNotificationRewardMultiplier
        );
        maliciousDkgResultNotificationRewardMultiplierChangeInitiated = 0;
        newMaliciousDkgResultNotificationRewardMultiplier = 0;
    }

    /// @notice Begins the DKG seed timeout update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgSeedTimeout New DKG seed timeout in blocks
    function beginDkgSeedTimeoutUpdate(uint256 _newDkgSeedTimeout)
        external
        onlyOwner
    {
        /* solhint-disable not-rely-on-time */
        require(_newDkgSeedTimeout > 0, "DKG seed timeout must be > 0");
        newDkgSeedTimeout = _newDkgSeedTimeout;
        dkgSeedTimeoutChangeInitiated = block.timestamp;
        emit DkgSeedTimeoutUpdateStarted(_newDkgSeedTimeout, block.timestamp);
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG seed timeout update process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgSeedTimeoutUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            dkgSeedTimeoutChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit DkgSeedTimeoutUpdated(newDkgSeedTimeout);
        // slither-disable-next-line reentrancy-no-eth
        walletRegistry.updateDkgParameters(
            newDkgSeedTimeout,
            walletRegistry.dkgParameters().resultChallengePeriodLength,
            walletRegistry.dkgParameters().resultSubmissionTimeout,
            walletRegistry.dkgParameters().submitterPrecedencePeriodLength
        );
        dkgSeedTimeoutChangeInitiated = 0;
        newDkgSeedTimeout = 0;
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
        walletRegistry.updateDkgParameters(
            walletRegistry.dkgParameters().seedTimeout,
            newDkgResultChallengePeriodLength,
            walletRegistry.dkgParameters().resultSubmissionTimeout,
            walletRegistry.dkgParameters().submitterPrecedencePeriodLength
        );
        dkgResultChallengePeriodLengthChangeInitiated = 0;
        newDkgResultChallengePeriodLength = 0;
    }

    /// @notice Begins the DKG result submission eligibility delay update
    ///         process.
    /// @dev Can be called only by the contract owner.
    /// @param _newDkgResultSubmissionTimeout New DKG result submission
    ///        eligibility delay in blocks
    function beginDkgResultSubmissionTimeoutUpdate(
        uint256 _newDkgResultSubmissionTimeout
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        require(
            _newDkgResultSubmissionTimeout > 0,
            "DKG result submission eligibility delay must be > 0"
        );
        newDkgResultSubmissionTimeout = _newDkgResultSubmissionTimeout;
        dkgResultSubmissionTimeoutChangeInitiated = block.timestamp;
        emit DkgResultSubmissionTimeoutUpdateStarted(
            _newDkgResultSubmissionTimeout,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG result submission eligibility delay update
    ///         process.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgResultSubmissionTimeoutUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            dkgResultSubmissionTimeoutChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit DkgResultSubmissionTimeoutUpdated(newDkgResultSubmissionTimeout);
        // slither-disable-next-line reentrancy-no-eth
        walletRegistry.updateDkgParameters(
            walletRegistry.dkgParameters().seedTimeout,
            walletRegistry.dkgParameters().resultChallengePeriodLength,
            newDkgResultSubmissionTimeout,
            walletRegistry.dkgParameters().submitterPrecedencePeriodLength
        );
        dkgResultSubmissionTimeoutChangeInitiated = 0;
        newDkgResultSubmissionTimeout = 0;
    }

    /// @notice Begins the DKG submitter precedence period length.
    /// @dev Can be called only by the contract owner.
    /// @param _newSubmitterPrecedencePeriodLength New DKG submitter precedence
    ///        period length in blocks
    function beginDkgSubmitterPrecedencePeriodLengthUpdate(
        uint256 _newSubmitterPrecedencePeriodLength
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        require(
            _newSubmitterPrecedencePeriodLength > 0,
            "DKG submitter precedence period length must be > 0"
        );
        newSubmitterPrecedencePeriodLength = _newSubmitterPrecedencePeriodLength;
        dkgSubmitterPrecedencePeriodLengthChangeInitiated = block.timestamp;
        emit DkgSubmitterPrecedencePeriodLengthUpdateStarted(
            _newSubmitterPrecedencePeriodLength,
            block.timestamp
        );
        /* solhint-enable not-rely-on-time */
    }

    /// @notice Finalizes the DKG submitter precedence period length.
    /// @dev Can be called only by the contract owner, after the governance
    ///      delay elapses.
    function finalizeDkgSubmitterPrecedencePeriodLengthUpdate()
        external
        onlyOwner
        onlyAfterGovernanceDelay(
            dkgSubmitterPrecedencePeriodLengthChangeInitiated,
            STANDARD_PARAMETER_GOVERNANCE_DELAY
        )
    {
        emit DkgSubmitterPrecedencePeriodLengthUpdated(
            newSubmitterPrecedencePeriodLength
        );
        // slither-disable-next-line reentrancy-no-eth
        walletRegistry.updateDkgParameters(
            walletRegistry.dkgParameters().seedTimeout,
            walletRegistry.dkgParameters().resultChallengePeriodLength,
            walletRegistry.dkgParameters().resultSubmissionTimeout,
            newSubmitterPrecedencePeriodLength
        );
        dkgSubmitterPrecedencePeriodLengthChangeInitiated = 0;
        newSubmitterPrecedencePeriodLength = 0;
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

    /// @notice Get the time remaining until the authorization decrease delay
    ///         can be updated.
    /// @return Remaining time in seconds.
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

    /// @notice Get the time remaining until the DKG malicious result
    ///         notification reward multiplier duration can be updated.
    /// @return Remaining time in seconds.
    function getRemainingMaliciousDkgResultNotificationRewardMultiplierUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                maliciousDkgResultNotificationRewardMultiplierChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the DKG seed timeout can be updated.
    /// @return Remaining time in seconds.
    function getRemainingDkgSeedTimeoutUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                dkgSeedTimeoutChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
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
    function getRemainingDkgResultSubmissionTimeoutUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                dkgResultSubmissionTimeoutChangeInitiated,
                STANDARD_PARAMETER_GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the wallet owner can be updated.
    /// @return Remaining time in seconds.
    function getRemainingWalletOwnerUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                walletOwnerChangeInitiated,
                CRITICAL_PARAMETER_GOVERNANCE_DELAY
            );
    }

    /// @notice Get the time remaining until the wallet owner can be updated.
    /// @return Remaining time in seconds.
    function getRemainingSubmitterPrecedencePeriodLengthUpdateTime()
        external
        view
        returns (uint256)
    {
        return
            getRemainingChangeTime(
                dkgSubmitterPrecedencePeriodLengthChangeInitiated,
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
        }

        return delay - elapsed;
    }
}
