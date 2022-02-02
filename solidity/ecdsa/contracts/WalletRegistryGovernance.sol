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

    uint256 public newMaliciousDkgResultSlashingAmount;
    uint256 public maliciousDkgResultSlashingAmountChangeInitiated;

    uint256 public newMaliciousDkgResultNotificationRewardMultiplier;
    uint256
        public maliciousDkgResultNotificationRewardMultiplierChangeInitiated;

    uint256 public newDkgResultChallengePeriodLength;
    uint256 public dkgResultChallengePeriodLengthChangeInitiated;

    uint256 public newDkgResultSubmissionEligibilityDelay;
    uint256 public dkgResultSubmissionEligibilityDelayChangeInitiated;

    WalletRegistry public walletRegistry;

    // Long governance delay used for critical parameters giving a chance for
    // stakers to opt out before the change is finalized in case they do not
    // agree with that change. The maximum group lifetime must not be longer
    // than this delay.
    //
    // The full list of parameters protected by this delay:
    // - wallet owner
    uint256 internal constant CRITICAL_PARAMETER_GOVERNANCE_DELAY = 2 weeks;

    // Short governance delay for non-critical parameters. Honest stakers should
    // not be severely affected by any change of these parameters.
    //
    // The full list of parameters protected by this delay:
    // - malicious DKG result notification reward multiplier
    // - malicious DKG result slashing amount
    // - DKG result challenge period length
    // - DKG result submission eligibility delay
    uint256 internal constant STANDARD_PARAMETER_GOVERNANCE_DELAY = 12 hours;

    event WalletOwnerUpdateStarted(
        address walletOwner,
        uint256 timestamp
    );
    event WalletOwnerUpdated(
        address walletOwner
    );

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

    /// @notice Begins the wallet owner update process.
    /// @dev Can be called only by the contract owner.
    /// @param _newWalletOwner New wallet owner address
    function beginWalletOwnerUpdate(
        address _newWalletOwner
    ) external onlyOwner {
        /* solhint-disable not-rely-on-time */
        newWalletOwner = _newWalletOwner;
        walletOwnerChangeInitiated = block.timestamp;
        emit WalletOwnerUpdateStarted(
            _newWalletOwner,
            block.timestamp
        );
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
        emit WalletOwnerUpdated(
            newWalletOwner
        );
        // slither-disable-next-line reentrancy-no-eth
        walletRegistry.updateWalletParameters(
            newWalletOwner
        );
        walletOwnerChangeInitiated = 0;
        newWalletOwner = address(0);
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
            newDkgResultChallengePeriodLength,
            walletRegistry.dkgParameters().resultSubmissionEligibilityDelay
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
        walletRegistry.updateDkgParameters(
            walletRegistry.dkgParameters().resultChallengePeriodLength,
            newDkgResultSubmissionEligibilityDelay
        );
        dkgResultSubmissionEligibilityDelayChangeInitiated = 0;
        newDkgResultSubmissionEligibilityDelay = 0;
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
