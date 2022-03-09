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

import {IWalletOwner} from "./api/IWalletOwner.sol";
import {IRandomBeacon} from "@keep-network/random-beacon/contracts/api/IRandomBeacon.sol";

/// @title Wallet Registry Governance
/// @notice Owns the `WalletRegistry` contract and is responsible for updating
///         its governable parameters in respect to the governance delay.
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

    uint256 public governanceDelay;

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
    modifier onlyAfterGovernanceDelay(uint256 changeInitiatedTimestamp) {
        /* solhint-disable not-rely-on-time */
        require(changeInitiatedTimestamp > 0, "Change not initiated");
        require(
            block.timestamp - changeInitiatedTimestamp >= governanceDelay,
            "Governance delay has not elapsed"
        );
        _;
        /* solhint-enable not-rely-on-time */
    }

    constructor(WalletRegistry _walletRegistry, uint256 _governanceDelay) {
        walletRegistry = _walletRegistry;
        governanceDelay = _governanceDelay;
    }

    /// @notice Upgrades the random beacon.
    /// @dev Can be called only by the contract owner.
    /// @param _newRandomBeacon New random beacon address
    function upgradeRandomBeacon(address _newRandomBeacon) external onlyOwner {
        require(
            _newRandomBeacon != address(0),
            "New random beacon address cannot be zero"
        );

        walletRegistry.upgradeRandomBeacon(IRandomBeacon(_newRandomBeacon));
    }

    /// @notice Initializes the Wallet Owner's address.
    /// @dev Can be called only by the contract owner. It can be called only if
    ///      walletOwner has not been set before. It doesn't enforce a governance
    ///      delay for the initial update. Any subsequent updates should be performed
    ///      with beginWalletOwnerUpdate/finalizeWalletOwnerUpdate with respect
    ///      of a governance delay.
    /// @param _walletOwner The Wallet Owner's address
    function initializeWalletOwner(address _walletOwner) external onlyOwner {
        require(
            address(walletRegistry.walletOwner()) == address(0),
            "Wallet Owner already initialized"
        );
        require(
            _walletOwner != address(0),
            "Wallet Owner address cannot be zero"
        );

        walletRegistry.updateWalletOwner(IWalletOwner(_walletOwner));
    }

    /// @notice Begins the wallet owner update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newWalletOwner New wallet owner address
    function beginWalletOwnerUpdate(address _newWalletOwner)
        external
        onlyOwner
    {
        require(
            address(_newWalletOwner) != address(0),
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
        onlyAfterGovernanceDelay(walletOwnerChangeInitiated)
    {
        emit WalletOwnerUpdated(newWalletOwner);
        // slither-disable-next-line reentrancy-no-eth
        walletRegistry.updateWalletOwner(IWalletOwner(newWalletOwner));
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
        onlyAfterGovernanceDelay(minimumAuthorizationChangeInitiated)
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
        onlyAfterGovernanceDelay(authorizationDecreaseDelayChangeInitiated)
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
            maliciousDkgResultSlashingAmountChangeInitiated
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
            maliciousDkgResultNotificationRewardMultiplierChangeInitiated
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
        onlyAfterGovernanceDelay(dkgSeedTimeoutChangeInitiated)
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
        onlyAfterGovernanceDelay(dkgResultChallengePeriodLengthChangeInitiated)
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
        onlyAfterGovernanceDelay(dkgResultSubmissionTimeoutChangeInitiated)
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
            dkgSubmitterPrecedencePeriodLengthChangeInitiated
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
        return getRemainingChangeTime(minimumAuthorizationChangeInitiated);
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
            getRemainingChangeTime(authorizationDecreaseDelayChangeInitiated);
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
                maliciousDkgResultSlashingAmountChangeInitiated
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
                maliciousDkgResultNotificationRewardMultiplierChangeInitiated
            );
    }

    /// @notice Get the time remaining until the DKG seed timeout can be updated.
    /// @return Remaining time in seconds.
    function getRemainingDkgSeedTimeoutUpdateTime()
        external
        view
        returns (uint256)
    {
        return getRemainingChangeTime(dkgSeedTimeoutChangeInitiated);
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
                dkgResultChallengePeriodLengthChangeInitiated
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
            getRemainingChangeTime(dkgResultSubmissionTimeoutChangeInitiated);
    }

    /// @notice Get the time remaining until the wallet owner can be updated.
    /// @return Remaining time in seconds.
    function getRemainingWalletOwnerUpdateTime()
        external
        view
        returns (uint256)
    {
        return getRemainingChangeTime(walletOwnerChangeInitiated);
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
                dkgSubmitterPrecedencePeriodLengthChangeInitiated
            );
    }

    /// @notice Gets the time remaining until the governable parameter update
    ///         can be committed.
    /// @param changeTimestamp Timestamp indicating the beginning of the change.
    /// @return Remaining time in seconds.
    function getRemainingChangeTime(uint256 changeTimestamp)
        internal
        view
        returns (uint256)
    {
        require(changeTimestamp > 0, "Change not initiated");
        /* solhint-disable-next-line not-rely-on-time */
        uint256 elapsed = block.timestamp - changeTimestamp;
        if (elapsed >= governanceDelay) {
            return 0;
        }

        return governanceDelay - elapsed;
    }

    // TODO: Add function to transfer WalletRegistry ownership to another address.

    // TODO: Add function to update governance delay
}
