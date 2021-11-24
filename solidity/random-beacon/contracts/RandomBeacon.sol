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
//                           Trust math, not hardware.

pragma solidity ^0.8.6;

import "./libraries/Authorization.sol";
import "./libraries/DKG.sol";
import "./libraries/GasStation.sol";
import "./libraries/Groups.sol";
import "./libraries/Relay.sol";
import "./libraries/Groups.sol";
import "./libraries/Callback.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";

/// @title Staking contract interface
/// @notice This is an interface with just a few function signatures of the
///         Staking contract, which is available at
///         https://github.com/threshold-network/solidity-contracts/blob/main/contracts/staking/IStaking.sol
///
/// TODO: Add a dependency to `threshold-network/solidity-contracts` and use
///       staking interface from there.
interface IRandomBeaconStaking {
    function slash(uint256 amount, address[] memory operators) external;

    function seize(
        uint256 amount,
        uint256 rewardMultiplier,
        address notifier,
        address[] memory operators
    ) external;
}

/// @title Keep Random Beacon
/// @notice Keep Random Beacon contract. It lets anyone request a new
///         relay entry and validates the new relay entry provided by the
///         network. This contract is in charge of all Random Beacon maintenance
///         activities such as group lifecycle or slashing.
/// @dev Should be owned by the governance contract controlling Random Beacon
///      parameters.
contract RandomBeacon is Ownable {
    using SafeERC20 for IERC20;
    using Authorization for Authorization.Data;
    using DKG for DKG.Data;
    using Groups for Groups.Data;
    using Relay for Relay.Data;
    using Callback for Callback.Data;
    using GasStation for GasStation.Data;

    // Constant parameters

    /// @notice Seed value used for the genesis group selection.
    /// https://www.wolframalpha.com/input/?i=pi+to+78+digits
    uint256 public constant genesisSeed =
        31415926535897932384626433832795028841971693993751058209749445923078164062862;

    // Governable parameters

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

    /// @notice Reward in T for submitting DKG result. The reward is paid to
    ///         a submitter of a valid DKG result when the DKG result challenge
    ///         period ends.
    uint256 public dkgResultSubmissionReward;

    /// @notice Reward in T for unlocking the sortition pool if DKG timed out.
    ///         When DKG result submission timed out, sortition pool is still
    ///         locked and someone needs to unlock it. Anyone can do it and earn
    ///         `sortitionPoolUnlockingReward`.
    uint256 public sortitionPoolUnlockingReward;

    /// @notice Slashing amount for supporting malicious DKG result. Every
    ///         DKG result submitted can be challenged for the time of
    ///         `dkgResultChallengePeriodLength`. If the DKG result submitted
    ///         is challenged and proven to be malicious, each operator who
    ///         signed the malicious result is slashed for
    ///         `maliciousDkgResultSlashingAmount`.
    uint256 public maliciousDkgResultSlashingAmount;

    /// @notice Duration of the sortition pool rewards ban imposed on operators
    ///         who missed their turn for relay entry or DKG result submission.
    uint256 public sortitionPoolRewardsBanDuration;

    /// @notice Percentage of the staking contract malicious behavior
    ///         notification reward which will be transferred to the notifier
    ///         reporting about relay entry timeout. Notifiers are rewarded
    ///         from a separate pool funded from slashed tokens. For example, if
    ///         notification reward is 1000 and the value of the multiplier is
    ///         5, the notifier will receive: 5% of 1000 = 50 per each
    ///         operator affected.
    uint256 public relayEntryTimeoutNotificationRewardMultiplier;

    /// @notice Percentage of the staking contract malicious behavior
    ///         notification reward which will be transferred to the notifier
    ///         reporting about a malicious DKG result. Notifiers are rewarded
    ///         from a separate pool funded from slashed tokens. For example, if
    ///         notification reward is 1000 and the value of the multiplier is
    ///         5, the notifier will receive: 5% of 1000 = 50 per each
    ///         operator affected.
    uint256 public dkgMaliciousResultNotificationRewardMultiplier;

    SortitionPool public sortitionPool;
    IERC20 public tToken;
    IRandomBeaconStaking public staking;

    // Libraries data storages
    Authorization.Data internal authorization;
    DKG.Data internal dkg;
    Groups.Data internal groups;
    Relay.Data internal relay;
    Callback.Data internal callback;
    GasStation.Data internal gasStation;

    event AuthorizationParametersUpdated(
        uint96 minimumAuthorization,
        uint64 authorizationDecreaseDelay
    );

    event RelayEntryParametersUpdated(
        uint256 relayRequestFee,
        uint256 relayEntrySubmissionEligibilityDelay,
        uint256 relayEntryHardTimeout,
        uint256 callbackGasLimit
    );

    event DkgParametersUpdated(
        uint256 dkgResultChallengePeriodLength,
        uint256 dkgResultSubmissionEligibilityDelay
    );

    event GroupCreationParametersUpdated(
        uint256 groupCreationFrequency,
        uint256 groupLifetime
    );

    event RewardParametersUpdated(
        uint256 dkgResultSubmissionReward,
        uint256 sortitionPoolUnlockingReward,
        uint256 sortitionPoolRewardsBanDuration,
        uint256 relayEntryTimeoutNotificationRewardMultiplier,
        uint256 dkgMaliciousResultNotificationRewardMultiplier
    );

    event SlashingParametersUpdated(
        uint256 relayEntrySubmissionFailureSlashingAmount,
        uint256 maliciousDkgResultSlashingAmount
    );

    // Events copied from library to workaround issue https://github.com/ethereum/solidity/issues/9765

    event DkgStarted(uint256 indexed seed);

    event DkgResultSubmitted(
        bytes32 indexed resultHash,
        uint256 indexed seed,
        uint256 submitterMemberIndex,
        bytes indexed groupPubKey,
        uint8[] misbehavedMembersIndices,
        bytes signatures,
        uint256[] signingMembersIndices,
        uint32[] members
    );

    event DkgTimedOut();

    event DkgResultApproved(
        bytes32 indexed resultHash,
        address indexed approver
    );

    event DkgResultChallenged(
        bytes32 indexed resultHash,
        address indexed challenger
    );

    event DkgMaliciousResultSlashed(
        bytes32 indexed resultHash,
        uint256 slashingAmount,
        address[] groupMembers
    );

    event DkgStateLocked();

    event DkgSeedTimedOut();

    event CandidateGroupRegistered(bytes indexed groupPubKey);

    event CandidateGroupRemoved(bytes indexed groupPubKey);

    event GroupActivated(uint64 indexed groupId, bytes indexed groupPubKey);

    event RelayEntryRequested(
        uint256 indexed requestId,
        uint64 groupId,
        bytes previousEntry
    );

    event RelayEntrySubmitted(uint256 indexed requestId, bytes entry);

    event RelayEntryTimedOut(
        uint256 indexed requestId,
        uint64 terminatedGroupId
    );

    event RelayEntryDelaySlashed(
        uint256 indexed requestId,
        uint256 slashingAmount,
        address[] groupMembers
    );

    event RelayEntryDelaySlashingFailed(
        uint256 indexed requestId,
        uint256 slashingAmount,
        address[] groupMembers
    );

    event RelayEntryTimeoutSlashed(
        uint256 indexed requestId,
        uint256 slashingAmount,
        address[] groupMembers
    );

    event CallbackFailed(uint256 entry, uint256 entrySubmittedBlock);

    event BanRewardsFailed(uint32[] ids);

    /// @dev Assigns initial values to parameters to make the beacon work
    ///      safely. These parameters are just proposed defaults and they might
    ///      be updated with `update*` functions after the contract deployment
    ///      and before transferring the ownership to the governance contract.
    constructor(
        SortitionPool _sortitionPool,
        IERC20 _tToken,
        IRandomBeaconStaking _staking
    ) {
        sortitionPool = _sortitionPool;
        tToken = _tToken;
        staking = _staking;

        // TODO: revisit all initial values
        callbackGasLimit = 200e3;
        groupCreationFrequency = 10;
        groupLifetime = 2 weeks;
        dkgResultSubmissionReward = 0;
        sortitionPoolUnlockingReward = 0;
        maliciousDkgResultSlashingAmount = 50000e18;
        sortitionPoolRewardsBanDuration = 2 weeks;
        relayEntryTimeoutNotificationRewardMultiplier = 5;
        dkgMaliciousResultNotificationRewardMultiplier = 5;
        // slither-disable-next-line too-many-digits
        authorization.setMinimumAuthorization(100000 * 1e18);

        dkg.setResultChallengePeriodLength(sortitionPool, 1440); // ~6h assuming 15s block time
        dkg.setResultSubmissionEligibilityDelay(sortitionPool, 10);

        relay.initSeedEntry();
        relay.setRelayEntrySubmissionEligibilityDelay(10);
        relay.setRelayEntryHardTimeout(5760); // ~24h assuming 15s block time
        relay.setRelayEntrySubmissionFailureSlashingAmount(1000e18);
    }

    /// @notice Updates the values of authorization parameters.
    /// @dev Can be called only by the contract owner, which should be the
    ///      random beacon governance contract. The caller is responsible for
    ///      validating parameters.
    /// @param _minimumAuthorization New minimum authorization amount
    /// @param _authorizationDecreaseDelay New authorization decrease delay in
    ///        seconds
    function updateAuthorizationParameters(
        uint96 _minimumAuthorization,
        uint64 _authorizationDecreaseDelay
    ) external onlyOwner {
        authorization.setMinimumAuthorization(_minimumAuthorization);
        authorization.setAuthorizationDecreaseDelay(
            _authorizationDecreaseDelay
        );

        emit AuthorizationParametersUpdated(
            _minimumAuthorization,
            _authorizationDecreaseDelay
        );
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
        callbackGasLimit = _callbackGasLimit;

        relay.setRelayRequestFee(_relayRequestFee);
        relay.setRelayEntrySubmissionEligibilityDelay(
            _relayEntrySubmissionEligibilityDelay
        );
        relay.setRelayEntryHardTimeout(_relayEntryHardTimeout);

        emit RelayEntryParametersUpdated(
            _relayRequestFee,
            _relayEntrySubmissionEligibilityDelay,
            _relayEntryHardTimeout,
            callbackGasLimit
        );
    }

    /// @notice Updates the values of group creation parameters.
    /// @dev Can be called only by the contract owner, which should be the
    ///      random beacon governance contract. The caller is responsible for
    ///      validating parameters.
    /// @param _groupCreationFrequency New group creation frequency
    /// @param _groupLifetime New group lifetime
    function updateGroupCreationParameters(
        uint256 _groupCreationFrequency,
        uint256 _groupLifetime
    ) external onlyOwner {
        groupCreationFrequency = _groupCreationFrequency;
        groupLifetime = _groupLifetime;

        emit GroupCreationParametersUpdated(
            groupCreationFrequency,
            groupLifetime
        );
    }

    /// @notice Updates the values of DKG parameters.
    /// @dev Can be called only by the contract owner, which should be the
    ///      random beacon governance contract. The caller is responsible for
    ///      validating parameters.
    /// @param _dkgResultChallengePeriodLength New DKG result challenge period
    ///        length
    /// @param _dkgResultSubmissionEligibilityDelay New DKG result submission
    ///        eligibility delay
    function updateDkgParameters(
        uint256 _dkgResultChallengePeriodLength,
        uint256 _dkgResultSubmissionEligibilityDelay
    ) external onlyOwner {
        dkg.setResultChallengePeriodLength(
            sortitionPool,
            _dkgResultChallengePeriodLength
        );
        dkg.setResultSubmissionEligibilityDelay(
            sortitionPool,
            _dkgResultSubmissionEligibilityDelay
        );

        emit DkgParametersUpdated(
            dkgResultChallengePeriodLength(),
            dkgResultSubmissionEligibilityDelay()
        );
    }

    /// @notice Updates the values of reward parameters.
    /// @dev Can be called only by the contract owner, which should be the
    ///      random beacon governance contract. The caller is responsible for
    ///      validating parameters.
    /// @param _dkgResultSubmissionReward New DKG result submission reward
    /// @param _sortitionPoolUnlockingReward New sortition pool unlocking reward
    /// @param _sortitionPoolRewardsBanDuration New sortition pool rewards
    ///        ban duration in seconds.
    /// @param _relayEntryTimeoutNotificationRewardMultiplier New value of the
    ///        relay entry timeout notification reward multiplier.
    /// @param _dkgMaliciousResultNotificationRewardMultiplier New value of the
    ///        DKG malicious result notification reward multiplier.
    function updateRewardParameters(
        uint256 _dkgResultSubmissionReward,
        uint256 _sortitionPoolUnlockingReward,
        uint256 _sortitionPoolRewardsBanDuration,
        uint256 _relayEntryTimeoutNotificationRewardMultiplier,
        uint256 _dkgMaliciousResultNotificationRewardMultiplier
    ) external onlyOwner {
        dkgResultSubmissionReward = _dkgResultSubmissionReward;
        sortitionPoolUnlockingReward = _sortitionPoolUnlockingReward;
        sortitionPoolRewardsBanDuration = _sortitionPoolRewardsBanDuration;
        relayEntryTimeoutNotificationRewardMultiplier = _relayEntryTimeoutNotificationRewardMultiplier;
        dkgMaliciousResultNotificationRewardMultiplier = _dkgMaliciousResultNotificationRewardMultiplier;
        emit RewardParametersUpdated(
            dkgResultSubmissionReward,
            sortitionPoolUnlockingReward,
            sortitionPoolRewardsBanDuration,
            relayEntryTimeoutNotificationRewardMultiplier,
            dkgMaliciousResultNotificationRewardMultiplier
        );
    }

    /// @notice The number of blocks for which a DKG result can be challenged.
    ///         Anyone can challenge DKG result for a certain number of blocks
    ///         before the result is fully accepted and the group registered in
    ///         the pool of active groups. If the challenge gets accepted, all
    ///         operators who signed the malicious result get slashed for
    ///         `maliciousDkgResultSlashingAmount` and the notifier gets
    ///         rewarded.
    function dkgResultChallengePeriodLength() public view returns (uint256) {
        return dkg.parameters.resultChallengePeriodLength;
    }

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
    function dkgResultSubmissionEligibilityDelay()
        public
        view
        returns (uint256)
    {
        return dkg.parameters.resultSubmissionEligibilityDelay;
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
        relay.setRelayEntrySubmissionFailureSlashingAmount(
            _relayEntrySubmissionFailureSlashingAmount
        );
        maliciousDkgResultSlashingAmount = _maliciousDkgResultSlashingAmount;
        emit SlashingParametersUpdated(
            _relayEntrySubmissionFailureSlashingAmount,
            maliciousDkgResultSlashingAmount
        );
    }

    /// @notice Registers the caller in the sortition pool.
    /// @dev Creates a gas deposit tied to the operator address. The gas
    ///      deposit is released when the operator is banned from sortition
    ///      pool rewards or leaves the pool during status update.
    function registerOperator() external {
        address operator = msg.sender;

        require(
            !sortitionPool.isOperatorInPool(operator),
            "Operator is already registered"
        );

        gasStation.depositGas(operator);
        sortitionPool.insertOperator(operator);
    }

    /// @notice Updates the sortition pool status of the caller.
    function updateOperatorStatus() external {
        sortitionPool.updateOperatorStatus(
            sortitionPool.getOperatorID(msg.sender)
        );

        // If the operator has been removed from the sortition pool during the
        // status update, release its gas deposit.
        if (!sortitionPool.isOperatorInPool(msg.sender)) {
            gasStation.releaseGas(msg.sender);
        }
    }

    /// @notice Checks whether the given operator is eligible to join the
    ///         sortition pool.
    /// @param operator Address of the operator
    function isOperatorEligible(address operator) external view returns (bool) {
        return sortitionPool.isOperatorEligible(operator);
    }

    /// @notice Triggers group selection if there are no active groups.
    function genesis() external {
        require(groups.numberOfActiveGroups() == 0, "not awaiting genesis");

        dkg.lockState(sortitionPool);
        dkg.start(
            sortitionPool,
            uint256(keccak256(abi.encodePacked(genesisSeed, block.number)))
        );
    }

    /// @notice Submits result of DKG protocol. It is on-chain part of phase 14 of
    ///         the protocol. The DKG result consists of result submitting member
    ///         index, calculated group public key, bytes array of misbehaved
    ///         members, concatenation of signatures from group members,
    ///         indices of members corresponding to each signature and
    ///         the list of group members.
    ///         When the result is verified successfully it gets registered and
    ///         waits for an approval. A result can be challenged to verify the
    ///         members list corresponds to the expected set of members determined
    ///         by the sortition pool.
    ///         A candidate group is registered based on the submitted DKG result
    ///         details.
    /// @dev The message to be signed by each member is keccak256 hash of the
    ///      calculated group public key, misbehaved members as bytes and DKG
    ///      start block. The calculated hash should be prefixed with prefixed with
    ///      `\x19Ethereum signed message:\n` before signing, so the message to
    ///      sign is:
    ///      `\x19Ethereum signed message:\n${keccak256(groupPubKey,misbehaved,startBlock)}`
    /// @param dkgResult DKG result.
    function submitDkgResult(DKG.Result calldata dkgResult) external {
        dkg.submitResult(sortitionPool, dkgResult);

        groups.addCandidateGroup(
            dkgResult.groupPubKey,
            dkgResult.members,
            dkgResult.misbehavedMembersIndices
        );
    }

    /// @notice Notifies about DKG timeout. Pays the sortition pool unlocking
    ///         reward to the notifier.
    function notifyDkgTimeout() external {
        dkg.notifyTimeout(sortitionPool);
        // Pay the sortition pool unlocking reward.
        tToken.safeTransfer(msg.sender, sortitionPoolUnlockingReward);
        dkg.complete(sortitionPool);
    }

    /// @notice Approves DKG result. Can be called when the challenge period for
    ///         the submitted result is finished. Considers the submitted result
    ///         as valid, pays reward to the approver, bans misbehaved group
    ///         members from the sortition pool rewards, and completes the group
    ///         creation by activating the candidate group. For the first
    ///         `resultSubmissionEligibilityDelay` blocks after the end of the
    ///         challenge period can be called only by the DKG result submitter.
    ///         After that time, can be called by anyone.
    /// @param dkgResult Result to approve. Must match the submitted result
    ///        stored during `submitDkgResult`.
    function approveDkgResult(DKG.Result calldata dkgResult) external {
        uint32[] memory misbehavedMembers = dkg.approveResult(
            sortitionPool,
            dkgResult
        );

        uint256 maintenancePoolBalance = tToken.balanceOf(address(this));
        uint256 rewardToPay = Math.min(
            maintenancePoolBalance,
            dkgResultSubmissionReward
        );
        tToken.safeTransfer(msg.sender, rewardToPay);

        if (misbehavedMembers.length > 0) {
            banFromRewards(misbehavedMembers, sortitionPoolRewardsBanDuration);
        }

        groups.activateCandidateGroup();
        dkg.complete(sortitionPool);
        // TODO: Check if this function is cheap enough and it will be
        //       profitable for the DKG result submitter to call it.
    }

    /// @notice Challenges DKG result. If the submitted result is proved to be
    ///         invalid it reverts the DKG back to the result submission phase.
    ///         It removes a candidate group that was previously registered with
    ///         the DKG result submission.
    /// @param dkgResult Result to challenge. Must match the submitted result
    ///        stored during `submitDkgResult`.
    function challengeDkgResult(DKG.Result calldata dkgResult) external {
        (bytes32 maliciousResultHash, uint32[] memory maliciousMembers) = dkg
            .challengeResult(sortitionPool, dkgResult);

        uint256 slashingAmount = maliciousDkgResultSlashingAmount;
        address[] memory maliciousMembersAddresses = sortitionPool
            .getIDOperators(maliciousMembers);

        groups.popCandidateGroup();

        emit DkgMaliciousResultSlashed(
            maliciousResultHash,
            slashingAmount,
            maliciousMembersAddresses
        );

        staking.seize(
            slashingAmount,
            dkgMaliciousResultNotificationRewardMultiplier,
            msg.sender,
            maliciousMembersAddresses
        );
    }

    /// @notice Check current group creation state.
    function getGroupCreationState() external view returns (DKG.State) {
        return dkg.currentState(sortitionPool);
    }

    /// @notice Checks if DKG timed out. The DKG timeout period includes time required
    ///         for off-chain protocol execution and time for the result publication
    ///         for all group members. After this time result cannot be submitted
    ///         and DKG can be notified about the timeout.
    /// @return True if DKG timed out, false otherwise.
    function hasDkgTimedOut() external view returns (bool) {
        return dkg.hasDkgTimedOut(sortitionPool);
    }

    function getGroupsRegistry() external view returns (bytes32[] memory) {
        return groups.groupsRegistry;
    }

    function getGroup(uint64 groupId)
        external
        view
        returns (Groups.Group memory)
    {
        return groups.getGroup(groupId);
    }

    function getGroup(bytes memory groupPubKey)
        external
        view
        returns (Groups.Group memory)
    {
        return groups.getGroup(groupPubKey);
    }

    /// @notice Creates a request to generate a new relay entry, which will
    ///         include a random number (by signing the previous entry's
    ///         random number). Requires a request fee denominated in T token.
    /// @param callbackContract Beacon consumer callback contract.
    function requestRelayEntry(IRandomBeaconConsumer callbackContract)
        external
    {
        uint64 groupId = groups.selectGroup(
            uint256(keccak256(relay.previousEntry))
        );

        relay.requestEntry(groupId);

        tToken.safeTransferFrom(
            msg.sender,
            address(this),
            relay.relayRequestFee
        );

        callback.setCallbackContract(callbackContract);

        // If the current request should trigger group creation we need to lock
        // DKG state (lock sortition pool) to prevent operators from changing
        // its state before relay entry is known. That entry will be used as a
        // group selection seed.
        if (
            relay.requestCount % groupCreationFrequency == 0 &&
            dkg.currentState(sortitionPool) == DKG.State.IDLE
        ) {
            dkg.lockState(sortitionPool);
        }
    }

    /// @notice Creates a new relay entry.
    /// @param submitterIndex Index of the entry submitter.
    /// @param entry Group BLS signature over the previous entry.
    function submitRelayEntry(uint256 submitterIndex, bytes calldata entry)
        external
    {
        uint256 currentRequestId = relay.currentRequest.id;

        Groups.Group memory group = groups.getGroup(
            relay.currentRequest.groupId
        );

        (uint32[] memory inactiveMembers, uint256 slashingAmount) = relay
            .submitEntry(sortitionPool, submitterIndex, entry, group);

        if (inactiveMembers.length > 0) {
            banFromRewards(inactiveMembers, sortitionPoolRewardsBanDuration);
        }

        if (slashingAmount > 0) {
            address[] memory groupMembers = sortitionPool.getIDOperators(
                group.members
            );

            try staking.slash(slashingAmount, groupMembers) {
                // slither-disable-next-line reentrancy-events
                emit RelayEntryDelaySlashed(
                    currentRequestId,
                    slashingAmount,
                    groupMembers
                );
            } catch {
                // Should never happen but we want to ensure a non-critical path
                // failure from an external contract does not stop group members
                // from submitting a valid relay entry.
                emit RelayEntryDelaySlashingFailed(
                    currentRequestId,
                    slashingAmount,
                    groupMembers
                );
            }
        }

        // If DKG is awaiting a seed, that means the we should start the actual
        // group creation process.
        if (dkg.currentState(sortitionPool) == DKG.State.AWAITING_SEED) {
            dkg.start(sortitionPool, uint256(keccak256(entry)));
        }

        callback.executeCallback(uint256(keccak256(entry)), callbackGasLimit);
    }

    /// @notice Reports a relay entry timeout.
    function reportRelayEntryTimeout() external {
        uint64 groupId = relay.currentRequest.groupId;
        uint256 slashingAmount = relay
            .relayEntrySubmissionFailureSlashingAmount;
        address[] memory groupMembers = sortitionPool.getIDOperators(
            groups.getGroup(groupId).members
        );

        emit RelayEntryTimeoutSlashed(
            relay.currentRequest.id,
            slashingAmount,
            groupMembers
        );

        staking.seize(
            slashingAmount,
            relayEntryTimeoutNotificationRewardMultiplier,
            msg.sender,
            groupMembers
        );

        // TODO: Once implemented, terminate group using `groupId`.

        if (groups.numberOfActiveGroups() > 0) {
            groupId = groups.selectGroup(
                uint256(keccak256(relay.previousEntry))
            );
            relay.retryOnEntryTimeout(groupId);
        } else {
            relay.cleanupOnEntryTimeout();

            // If DKG is awaiting a seed, we should notify about its timeout to
            // avoid blocking the future group creation.
            if (dkg.currentState(sortitionPool) == DKG.State.AWAITING_SEED) {
                dkg.notifySeedTimedOut(sortitionPool);
            }
        }
    }

    /// @notice Ban given operators from sortition pool rewards.
    /// @dev By the way, this function releases gas deposits made by operators
    ///      during their registration. See `registerOperator` function. This
    ///      action makes banning cheaper gas-wise.
    /// @param ids IDs of banned operators.
    /// @param banDuration Duration of the ban period in seconds.
    function banFromRewards(uint32[] memory ids, uint256 banDuration) internal {
        try sortitionPool.banRewards(ids, banDuration) {
            address[] memory operators = sortitionPool.getIDOperators(ids);

            for (uint256 i = 0; i < operators.length; i++) {
                // TODO: Once `banRewards` is implemented on pool side, revisit
                //       gas station design. Current design is problematic
                //       because operators deposit gas upon registration and
                //       deposits are released either during status update
                //       or rewards ban. The first case is natural as operator
                //       leaves the pool but the latter is hard because deposit
                //       is released and operator still stays in the pool.
                gasStation.releaseGas(operators[i]);
            }
        } catch {
            // Should never happen but we want to ensure a non-critical path
            // failure from an external contract does not stop group members
            // from submitting a valid relay entry.
            // slither-disable-next-line reentrancy-events
            emit BanRewardsFailed(ids);
        }
    }

    /// @notice Locks the state of group creation.
    /// @dev This function is meant to be used by test stubs which inherits
    ///      from this contract and needs to lock the DKG state arbitrarily.
    function dkgLockState() internal {
        dkg.lockState(sortitionPool);
    }

    /// @notice The minimum authorization amount required so that operator can
    ///         participate in the random beacon. This amount is required to
    ///         execute slashing for providing a malicious DKG result or when
    ///         a relay entry times out.
    function minimumAuthorization() external view returns (uint96) {
        return authorization.minimumAuthorization;
    }

    /// @notice Delay in seconds that needs to pass between the time
    ///         authorization decrease is requested and the time that request
    ///         gets approved. Protects against free-riders earning rewards and
    ///         not being active in the network.
    function authorizationDecreaseDelay() external view returns (uint64) {
        return authorization.authorizationDecreaseDelay;
    }

    /// @return Flag indicating whether a relay entry request is currently
    ///         in progress.
    function isRelayRequestInProgress() external view returns (bool) {
        return relay.isRequestInProgress();
    }

    /// @return Relay request fee in T. This fee needs to be provided by the
    ///         account or contract requesting for a new relay entry.
    function relayRequestFee() external view returns (uint256) {
        return relay.relayRequestFee;
    }

    /// @return The number of blocks it takes for a group member to become
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
    function relayEntrySubmissionEligibilityDelay()
        external
        view
        returns (uint256)
    {
        return relay.relayEntrySubmissionEligibilityDelay;
    }

    /// @return Hard timeout in blocks for a group to submit the relay entry.
    ///         After all group members became eligible to submit the relay
    ///         entry and one more `relayEntrySubmissionEligibilityDelay` blocks
    ///         passed without relay entry submitted, all group members start
    ///         getting slashed. The slashing amount increases linearly until
    ///         the group submits the relay entry or until
    ///         `relayEntryHardTimeout` is reached. When the hard timeout is
    ///         reached, each group member will get slashed for
    ///         `relayEntrySubmissionFailureSlashingAmount`.
    function relayEntryHardTimeout() external view returns (uint256) {
        return relay.relayEntryHardTimeout;
    }

    /// @notice Slashing amount for not submitting relay entry. When
    ///         relay entry hard timeout is reached without the relay entry
    ///         submitted, each group member gets slashed for
    ///         `relayEntrySubmissionFailureSlashingAmount`. If the relay entry
    ///         gets submitted after the soft timeout (see
    ///         `relayEntrySubmissionEligibilityDelay` documentation), but
    ///         before the hard timeout, each group member gets slashed
    ///         proportionally to `relayEntrySubmissionFailureSlashingAmount`
    ///         and the time passed since the soft deadline.
    function relayEntrySubmissionFailureSlashingAmount()
        external
        view
        returns (uint256)
    {
        return relay.relayEntrySubmissionFailureSlashingAmount;
    }

    /// @notice Selects a new group of operators based on the provided seed.
    ///         At least one operator has to be registered in the pool,
    ///         otherwise the function fails reverting the transaction.
    /// @param seed Number used to select operators to the group.
    /// @return IDs of selected group members.
    function selectGroup(bytes32 seed) external view returns (uint32[] memory) {
        return sortitionPool.selectGroup(DKG.groupSize, seed);
    }
}
