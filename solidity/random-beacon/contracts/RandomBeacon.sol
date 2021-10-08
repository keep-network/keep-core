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

/// @title Keep Random Beacon
/// @notice Keep Random Beacon contract. It lets anyone request a new
///         relay entry and validates the new relay entry provided by the
///         network. This contract is in charge of all Random Beacon maintenance
///         activities such as group lifecycle or slashing.
/// @dev Should be owned by the governance contract controlling Random Beacon
///      parameters.
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

    /// @notice Relay request fee in T. This fee needs to be provided by the
    ///         account or contract requesting for a new relay entry.
    uint256 public relayRequestFee;

    /// @notice The number of blocks it takes for a group member to become
    ///         eligible to submit the relay entry. At first, there is only one
    ///         member in the group eligible to submit the relay entry. Then,
    ///         after `relayEntrySubmissionEligibilityDelay` blocks, another
    ///         group member becomes eligible so that there are two group
    ///         members eligible to submit the relay entry at that moment. After
    ///         another `relayEntrySubmissionEligibilityDelay` blocks, yet one
    ///         group member becomes eligible so that there are three group
    ///         members eligible to submit the relay entry at that moment. This
    ///         continues until all group members are eligible to submit the
    ///         relay entry or until the relay entry is submitted. If all
    ///         members became eligible to submit the relay entry and one more
    ///         `relayEntrySubmissionEligibilityDelay` passed without the relay
    ///         entry submitted, the group reaches soft timeout for submitting
    ///         the relay entry and the slashing starts.
    uint256 public relayEntrySubmissionEligibilityDelay;

    /// @notice Hard timeout in blocks for a group to submit the relay entry.
    ///         After all group members became eligible to submit the relay
    ///         entry and one more `relayEntrySubmissionEligibilityDelay` blocks
    ///         passed without relay entry submitted, all group members start
    ///         getting slashed. The slashing amount increases linearly until
    ///         the group submits the relay entry or until
    ///         `relayEntryHardTimeout` is reached. When the hard timeout is
    ///         reached, each group member will get slashed for
    ///         `relayEntrySubmissionFailureSlashingAmount`.
    uint256 public relayEntryHardTimeout;

    /// @notice Relay entry callback gas limit. This is the gas limit with which
    ///         callback function provided in the relay request transaction is
    ///         executed. The callback is executed with a new relay entry value
    ///         in the same transaction the relay entry is submitted.
    uint256 public callbackGasLimit;

    /// @notice The frequency of new group creation. Groups are created with
    ///         a fixed frequency of relay requests.
    uint256 public groupCreationFrequency;

    /// @notice Group lifetime in seconds. When a group reached its lifetime, it
    ///         is no longer selected for new relay requests but may still be
    ///         responsible for submitting relay entry if relay request assigned
    ///         to that group is still pending.
    uint256 public groupLifetime;

    /// @notice The number of blocks for which a DKG result can be challenged.
    ///         Anyone can challenge DKG result for a certain number of blocks
    ///         before the result is fully accepted and the group registered in
    ///         the pool of active groups. If the challenge gets accepted, all
    ///         operators who signed the malicious result get slashed for
    ///         `maliciousDkgResultSlashingAmount` and the notifier gets
    ///         rewarded.
    uint256 public dkgResultChallengePeriodLength;

    /// @notice The number of blocks it takes for a group member to become
    ///         eligible to submit the DKG result. At first, there is only one
    ///         member in the group eligible to submit the DKG result. Then,
    ///         after `dkgResultSubmissionEligibilityDelay` blocks, another
    ///         group member becomes eligible so that there are two group
    ///         members eligible to submit the DKG result at that moment. After
    ///         another `dkgResultSubmissionEligibilityDelay` blocks, yet one
    ///         group member becomes eligible to submit the DKG result so that
    ///         there are three group members eligible to submit the DKG result
    ///         at that moment. This continues until all group members are
    ///         eligible to submit the DKG result or until the DKG result is
    ///         submitted. If all members became eligible to submit the DKG
    ///         result and one more `dkgResultSubmissionEligibilityDelay` passed
    ///         without the DKG result submitted, DKG is considered as timed out
    ///         and no DKG result for this group creation can be submitted
    ///         anymore.
    uint256 public dkgResultSubmissionEligibilityDelay;

    /// @notice Reward in T for submitting DKG result. The reward is paid to
    ///         a submitter of a valid DKG result when the DKG result challenge
    ///         period ends.
    uint256 public dkgResultSubmissionReward;

    /// @notice Reward in T for unlocking the sortition pool if DKG timed out.
    ///         When DKG result submission timed out, sortition pool is still
    ///         locked and someone needs to unlock it. Anyone can do it and earn
    ///         `sortitionPoolUnlockingReward`.
    uint256 public sortitionPoolUnlockingReward;

    /// @notice Slashing amount for not submitting relay entry. When
    ///         relay entry hard timeout is reached without the relay entry
    ///         submitted, each group member gets slashed for
    ///         `relayEntrySubmissionFailureSlashingAmount`. If the relay entry
    ///         gets submitted after the soft timeout (see
    ///         `relayEntrySubmissionEligibilityDelay` documentation), but
    ///         before the hard timeout, each group member gets slashed
    ///         proportionally to `relayEntrySubmissionFailureSlashingAmount`
    ///         and the time passed since the soft deadline.
    uint256 public relayEntrySubmissionFailureSlashingAmount;

    /// @notice Slashing amount for supporting malicious DKG result. Every
    ///         DKG result submitted can be challenged for the time of
    ///         `dkgResultChallengePeriodLength`. If the DKG result submitted
    ///         is challenged and proven to be malicious, each operator who
    ///         signed the malicious result is slashed for
    ///         `maliciousDkgResultSlashingAmount`.
    uint256 public maliciousDkgResultSlashingAmount;

    // Libraries data storages

    // TODO: Can we really make it public along with the library functions?
    DKG.Data public dkg;

    event RelayEntryParametersUpdated(
        uint256 relayRequestFee,
        uint256 relayEntrySubmissionEligibilityDelay,
        uint256 relayEntryHardTimeout,
        uint256 callbackGasLimit
    );

    event GroupCreationParametersUpdated(
        uint256 groupCreationFrequency,
        uint256 groupLifetime,
        uint256 dkgResultChallengePeriodLength,
        uint256 dkgResultSubmissionEligibilityDelay
    );

    event RewardParametersUpdated(
        uint256 dkgResultSubmissionReward,
        uint256 sortitionPoolUnlockingReward
    );

    event SlashingParametersUpdated(
        uint256 relayEntrySubmissionFailureSlashingAmount,
        uint256 maliciousDkgResultSlashingAmount
    );

    event DkgStarted(
        uint256 seed,
        uint256 groupSize,
        uint256 dkgResultSubmissionEligibilityDelay
    ); // TODO: Add all other needed paramters

    event DkgResultSubmitted(
        bytes32 indexed seed,
        uint256 index,
        uint256 indexed submitterMemberIndex,
        bytes indexed groupPubKey,
        bytes misbehaved,
        bytes signatures,
        uint256[] signingMembersIndexes,
        address[] members
    ); // TODO: We could add submitter member address or ID

    event DkgResultApproved(
        bytes indexed groupPubKey,
        address indexed submitter
    );

    event DkgTimedOut(uint256 indexed seed);

    event DkgResultChallenged(
        bytes indexed groupPubKey,
        address indexed submitter
    ); // TODO: We could add submitter member address or ID

    // FIXME: This is just for tests
    uint256 internal currentRelayEntry = 420;

    /// @dev Assigns initial values to parameters to make the beacon work
    ///      safely. These parameters are just proposed defaults and they might
    ///      be updated with `update*` functions after the contract deployment
    ///      and before transferring the ownership to the governance contract.
    constructor(
        uint256 groupSize,
        uint256 signatureThreshold,
        uint256 timeDkg
    ) {
        relayRequestFee = 0;
        relayEntrySubmissionEligibilityDelay = 10;
        relayEntryHardTimeout = 5760; // ~24h assuming 15s block time
        callbackGasLimit = 200000;
        groupCreationFrequency = 10;
        groupLifetime = 2 weeks;
        dkgResultChallengePeriodLength = 1440; // ~6h assuming 15s block time
        dkgResultSubmissionEligibilityDelay = 10;
        dkgResultSubmissionReward = 0;
        sortitionPoolUnlockingReward = 0;
        relayEntrySubmissionFailureSlashingAmount = 1000e18;
        maliciousDkgResultSlashingAmount = 50000e18;

        GROUP_SIZE = groupSize;
        SIGNATURE_THRESHOLD = signatureThreshold;
        TIME_DKG = timeDkg;
    }

    /// @notice Updates the values of relay entry parameters.
    /// @dev Can be called only by the contract owner, which should be the
    ///      random beacon governance contract. The caller is responsible for
    ///      validating parameters.
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

    /// @notice Updates the values of group creation parameters.
    /// @dev Can be called only by the contract owner, which should be the
    ///      random beacon governance contract. The caller is responsible for
    ///      validating parameters.
    /// @param _groupCreationFrequency New group creation frequency
    /// @param _groupLifetime New group lifetime
    /// @param _dkgResultChallengePeriodLength New DKG result challenge period
    ///        length
    /// @param _dkgResultSubmissionEligibilityDelay New DKG result submission
    ///        eligibility delay
    function updateGroupCreationParameters(
        uint256 _groupCreationFrequency,
        uint256 _groupLifetime,
        uint256 _dkgResultChallengePeriodLength,
        uint256 _dkgResultSubmissionEligibilityDelay
    ) external onlyOwner {
        groupCreationFrequency = _groupCreationFrequency;
        groupLifetime = _groupLifetime;
        dkgResultChallengePeriodLength = _dkgResultChallengePeriodLength;
        dkgResultSubmissionEligibilityDelay = _dkgResultSubmissionEligibilityDelay;
        emit GroupCreationParametersUpdated(
            groupCreationFrequency,
            groupLifetime,
            dkgResultChallengePeriodLength,
            dkgResultSubmissionEligibilityDelay
        );
    }

    /// @notice Updates the values of reward parameters.
    /// @dev Can be called only by the contract owner, which should be the
    ///      random beacon governance contract. The caller is responsible for
    ///      validating parameters.
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

    /// @notice Updates the values of slashing parameters.
    /// @dev Can be called only by the contract owner, which should be the
    ///      random beacon governance contract. The caller is responsible for
    ///      validating parameters.
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
            dkgResultSubmissionEligibilityDelay,
            dkgResultChallengePeriodLength,
            TIME_DKG
        );

        emit DkgStarted(seed, GROUP_SIZE, dkgResultSubmissionEligibilityDelay);
    }

    function genesis() external {
        // require(groups.groupCount == 0, "not awaiting genesis");

        createGroup(GENESIS_SEED);
    }

    function completeGroupCreation() internal {
        dkg.finish();

        // New groups should be created with a fixed frequency of relay requests
        // TODO: Consider each group a separate contract instance deployed with proxy?
    }

    function notifyDkgTimeout() external {
        dkg.notifyTimeout();

        emit DkgTimedOut(dkg.seed);
    }

    function isDkgInProgress() external view returns (bool) {
        return dkg.isInProgress();
    }

    /// @notice Submits result of DKG protocol. It is on-chain part of phase 14 of
    /// the protocol.
    /// @param dkgResult DKG result.
    function submitDkgResult(DKG.DkgResult calldata dkgResult) external {
        // TODO: Consider adding nonReentrant?

        // validate DKG result
        uint256 resultIndex = dkg.submitDkgResult(dkgResult);

        groups.addGroup(dkgResult.groupPubKey);

        emit DkgResultSubmitted(
            dkg.seed,
            resultIndex,
            dkgResult.submitterMemberIndex,
            dkgResult.groupPubKey,
            dkgResult.misbehaved,
            dkgResult.signatures,
            dkgResult.signingMembersIndexes,
            dkgResult.members
        );
    }

    function challengeDkgResult(
        uint256 resultIndex,
        DKG.DkgResult calldata dkgResult
    ) external {
        dkg.challengeResult(resultIndex, dkgResult);

        // TODO: Slash submitter

        emit DkgResultChallenged(
            dkgResult.groupPubKey,
            dkgResult.members[dkgResult.submitterMemberIndex] // TODO: Double check if this should be `submitterMemberIndex` or `submitterMemberIndex + 1`?
        );
    }

    // Once the challenge period passes, anyone can unlock the sortition pool and mark the DKG result as accepted.
    function approveDkgResult(
        uint256 resultIndex,
        DKG.DkgResult calldata dkgResult
    ) external {
        dkg.acceptResult(resultIndex, dkgResult);

        groups.activateGroup(dkgResult.groupPubKey);
        // groups.setGroupMembers(groupPubKey, members, misbehaved);

        emit DkgResultApproved(
            dkgResult.groupPubKey,
            dkgResult.members[dkgResult.submitterMemberIndex] // TODO: Double check if this should be `submitterMemberIndex` or `submitterMemberIndex + 1`?
        );

        // TODO: Unlock sortition pool
    }

    // function requestRelayEntry() external {
    //     if RELAY_ENTRY_COUNT >= groupCreationFrequency {
    //         createGroup();
    //     }
    // }
}
