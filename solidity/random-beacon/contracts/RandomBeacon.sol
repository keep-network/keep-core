// SPDX-License-Identifier: MIT
/*
▓▓▌ ▓▓ ▐▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▄
▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
  ▓▓▓▓▓▓    ▓▓▓▓▓▓▓▀    ▐▓▓▓▓▓▓    ▐▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
  ▓▓▓▓▓▓▄▄▓▓▓▓▓▓▓▀      ▐▓▓▓▓▓▓▄▄▄▄         ▓▓▓▓▓▓▄▄▄▄         ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
  ▓▓▓▓▓▓▓▓▓▓▓▓▓▀        ▐▓▓▓▓▓▓▓▓▓▓         ▓▓▓▓▓▓▓▓▓▓▌        ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
  ▓▓▓▓▓▓▀▀▓▓▓▓▓▓▄       ▐▓▓▓▓▓▓▀▀▀▀         ▓▓▓▓▓▓▀▀▀▀         ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▀
  ▓▓▓▓▓▓   ▀▓▓▓▓▓▓▄     ▐▓▓▓▓▓▓     ▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌
▓▓▓▓▓▓▓▓▓▓ █▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓
▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓

                           Trust math, not hardware.
*/
pragma solidity ^0.8.6;

import "./libraries/DKG.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @title Random beacon
/// @notice Random beacon contract which represents the on-chain part of the
///         random beacon functionality
/// @dev Should be owned by the random beacon governance contract
contract RandomBeacon is Ownable {
    using DKG for DKG.Data;

    // Constant parameters

    /// @notice Seed value used for the genesis group selection.
    /// https://www.wolframalpha.com/input/?i=pi+to+78+digits
    uint256 public constant GENESIS_SEED =
        31415926535897932384626433832795028841971693993751058209749445923078164062862;

    /// @notice Size of a group in the threshold relay.
    uint256 public immutable GROUP_SIZE;

    /// @notice The minimum number of signatures required to support DKG result.
    uint256 public immutable SIGNATURE_THRESHOLD;

    /// @notice Time in blocks after which DKG result is complete and ready to be
    // published by clients.
    uint256 public immutable TIME_DKG;

    // Governable parameters

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
    // FIXME: we set it with a value just for tests, the value should be initialized properly in another PR.
    uint256 public dkgSubmissionEligibilityDelay = 10;
    /// @notice Reward for submitting DKG result
    uint256 public dkgResultSubmissionReward;

    /// @notice Reward for unlocking the sortition pool if DKG timed out
    uint256 public sortitionPoolUnlockingReward;

    /// @notice Slashing amount for not submitting relay entry
    uint256 public relayEntrySubmissionFailureSlashingAmount;

    /// @notice Slashing amount for submitting malicious DKG result
    uint256 public maliciousDkgResultSlashingAmount;

    // Libraries data storages

    // TODO: Can we really make it public along with the library functions?
    DKG.Data public dkg;

    error NotAwaitingGenesis(uint256 groupCount);

    event RelayEntryParametersUpdated(
        uint256 relayRequestFee,
        uint256 relayEntrySubmissionEligibilityDelay,
        uint256 RelayEntryHardTimeout,
        uint256 callbackGasLimit
    );

    event GroupCreationParametersUpdated(
        uint256 groupCreationFrequency,
        uint256 groupLifetime,
        uint256 dkgResultChallengePeriodLength,
        uint256 dkgSubmissionEligibilityDelay
    );

    event RewardParametersUpdated(
        uint256 dkgResultSubmissionReward,
        uint256 sortitionPoolUnlockingReward
    );

    event SlashingParametersUpdated(
        uint256 relayEntrySubmissionFailureSlashingAmount,
        uint256 maliciousDkgResultSlashingAmount
    );

    // Events copied from library to workaround issue https://github.com/ethereum/solidity/issues/9765
    event DkgStarted(
        uint256 seed,
        uint256 groupSize,
        uint256 dkgSubmissionEligibilityDelay
    );
    event DkgTimedOut(uint256 seed);
    event DkgCompleted(uint256 seed);

    // FIXME: This is just for tests
    uint256 internal currentRelayEntry = 420;

    constructor(
        uint256 groupSize,
        uint256 signatureThreshold,
        uint256 timeDkg
    ) {
        GROUP_SIZE = groupSize;
        SIGNATURE_THRESHOLD = signatureThreshold;
        TIME_DKG = timeDkg;
    }

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
        emit RelayEntryParametersUpdated(
            relayRequestFee,
            relayEntrySubmissionEligibilityDelay,
            relayEntryHardTimeout,
            callbackGasLimit
        );
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
        emit GroupCreationParametersUpdated(
            groupCreationFrequency,
            groupLifetime,
            dkgResultChallengePeriodLength,
            dkgSubmissionEligibilityDelay
        );
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
        emit RewardParametersUpdated(
            dkgResultSubmissionReward,
            sortitionPoolUnlockingReward
        );
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
        emit SlashingParametersUpdated(
            relayEntrySubmissionFailureSlashingAmount,
            maliciousDkgResultSlashingAmount
        );
    }

    function createGroup(uint256 seed) internal {
        // Sortition performed off-chain

        dkg.start(
            seed,
            GROUP_SIZE,
            SIGNATURE_THRESHOLD,
            dkgSubmissionEligibilityDelay,
            TIME_DKG
        );
    }

    function genesis() external {
        // if (groups.groupCount > 0)
        //     revert NotAwaitingGenesis(dkg.groupCount, dkg.currentState);

        createGroup(GENESIS_SEED);
    }

    function completeGroupCreation() internal {
        dkg.finish();

        // New groups should be created with a fixed frequency of relay requests
        // TODO: Consider each group a separate contract instance deployed with proxy?
    }

    function notifyDkgTimeout() external {
        dkg.notifyTimeout();
    }

    function isDkgInProgress() external view returns (bool) {
        return dkg.isInProgress();
    }

    // params:
    // - dkg result
    // - group members (for verification)
    function submitDkgResult() external {
        // TODO: Consider adding nonReentrant?

        // validate DKG result
        // dkgResultVerification.verify(
        //     submitterMemberIndex,
        //     groupPubKey,
        //     misbehaved,
        //     signatures,
        //     signingMembersIndexes,
        //     members,
        //     groupSelection.ticketSubmissionStartBlock +
        //         groupSelection.ticketSubmissionTimeout
        // );

        // check member eligibility to submit result to submit result, w odpowiednim przedziale blokow gosc z tym ID

        // if enough results
        completeGroupCreation();
    }

    // function requestRelayEntry() external {
    //     if RELAY_ENTRY_COUNT >= groupCreationFrequency {
    //         createGroup();
    //     }
    // }
}
