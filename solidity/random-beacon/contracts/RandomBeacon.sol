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

pragma solidity ^0.8.9;

import "./libraries/Authorization.sol";
import "./libraries/DKG.sol";
import "./libraries/Groups.sol";
import "./libraries/Relay.sol";
import "./libraries/Groups.sol";
import "./libraries/Callback.sol";
import "./libraries/Heartbeat.sol";
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

    function eligibleStake(address operator, address operatorContract)
        external
        view
        returns (uint256);
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

    /// @notice Reward in T for submitting DKG result. The reward is paid to
    ///         a submitter of a valid DKG result when the DKG result challenge
    ///         period ends.
    uint256 public dkgResultSubmissionReward;

    /// @notice Reward in T for unlocking the sortition pool if DKG timed out.
    ///         When DKG result submission timed out, sortition pool is still
    ///         locked and someone needs to unlock it. Anyone can do it and earn
    ///         `sortitionPoolUnlockingReward`.
    uint256 public sortitionPoolUnlockingReward;

    /// @notice Reward in T for notifying the operator is ineligible.
    uint256 public ineligibleOperatorNotifierReward;

    /// @notice Slashing amount for supporting malicious DKG result. Every
    ///         DKG result submitted can be challenged for the time of
    ///         `dkgResultChallengePeriodLength`. If the DKG result submitted
    ///         is challenged and proven to be malicious, each operator who
    ///         signed the malicious result is slashed for
    ///         `maliciousDkgResultSlashingAmount`.
    uint256 public maliciousDkgResultSlashingAmount;

    /// @notice Slashing amount when an unauthorized signing has been proved,
    ///         which means the private key has been leaked and all the group
    ///         members should be punished.
    uint256 public unauthorizedSigningSlashingAmount;

    /// @notice Duration of the sortition pool rewards ban imposed on operators
    ///         who missed their turn for relay entry or DKG result submission.
    uint256 public sortitionPoolRewardsBanDuration;

    /// @notice Percentage of the staking contract malicious behavior
    ///         notification reward which will be transferred to the notifier
    ///         reporting about relay entry timeout. Notifiers are rewarded
    ///         from a notifiers treasury pool. For example, if
    ///         notification reward is 1000 and the value of the multiplier is
    ///         5, the notifier will receive: 5% of 1000 = 50 per each
    ///         operator affected.
    uint256 public relayEntryTimeoutNotificationRewardMultiplier;

    /// @notice Percentage of the staking contract malicious behavior
    ///         notification reward which will be transferred to the notifier
    ///         reporting about unauthorized signing. Notifiers are rewarded
    ///         from a notifiers treasury pool. For example, if a
    ///         notification reward is 1000 and the value of the multiplier is
    ///         5, the notifier will receive: 5% of 1000 = 50 per each
    ///         operator affected.
    uint256 public unauthorizedSigningNotificationRewardMultiplier;

    /// @notice Percentage of the staking contract malicious behavior
    ///         notification reward which will be transferred to the notifier
    ///         reporting about a malicious DKG result. Notifiers are rewarded
    ///         from a notifiers treasury pool. For example, if
    ///         notification reward is 1000 and the value of the multiplier is
    ///         5, the notifier will receive: 5% of 1000 = 50 per each
    ///         operator affected.
    uint256 public dkgMaliciousResultNotificationRewardMultiplier;

    // Other parameters

    /// @notice Stores current failed heartbeat nonce for given group.
    ///         Each claim is made with an unique nonce which protects
    ///         against claim replay.
    mapping(uint64 => uint256) public failedHeartbeatNonce; // groupId -> nonce

    // External dependencies

    SortitionPool public sortitionPool;
    IERC20 public tToken;
    IRandomBeaconStaking public staking;

    // Token bookkeeping

    /// @notice Rewards pool for DKG actions. This pool is funded from relay
    ///         request fees and external donates. Funds are used to cover
    ///         rewards for actors approving DKG result or notifying about
    ///         DKG timeout.
    uint256 public dkgRewardsPool;
    /// @notice Rewards pool for heartbeat notifiers. This pool is funded from
    ///         external donates. Funds are used to cover rewards for actors
    ///         notifying about failed heartbeats.
    uint256 public heartbeatNotifierRewardsPool;

    // Libraries data storages

    Authorization.Data internal authorization;
    DKG.Data internal dkg;
    Groups.Data internal groups;
    Relay.Data internal relay;
    Callback.Data internal callback;

    // Events

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
        uint256 ineligibleOperatorNotifierReward,
        uint256 sortitionPoolRewardsBanDuration,
        uint256 relayEntryTimeoutNotificationRewardMultiplier,
        uint256 unauthorizedSigningNotificationRewardMultiplier,
        uint256 dkgMaliciousResultNotificationRewardMultiplier
    );

    event SlashingParametersUpdated(
        uint256 relayEntrySubmissionFailureSlashingAmount,
        uint256 maliciousDkgResultSlashingAmount,
        uint256 unauthorizedSigningSlashingAmount
    );

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
        address indexed challenger,
        string reason
    );

    event DkgMaliciousResultSlashed(
        bytes32 indexed resultHash,
        uint256 slashingAmount,
        address maliciousSubmitter
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

    event RelayEntrySubmitted(
        uint256 indexed requestId,
        address submitter,
        bytes entry
    );

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

    event UnauthorizedSigningSlashed(
        uint64 indexed groupId,
        uint256 unauthorizedSigningSlashingAmount,
        address[] groupMembers
    );

    event CallbackFailed(uint256 entry, uint256 entrySubmittedBlock);

    event HeartbeatFailed(uint64 groupId, uint256 nonce, address notifier);

    /// @dev Assigns initial values to parameters to make the beacon work
    ///      safely. These parameters are just proposed defaults and they might
    ///      be updated with `update*` functions after the contract deployment
    ///      and before transferring the ownership to the governance contract.
    constructor(
        SortitionPool _sortitionPool,
        IERC20 _tToken,
        IRandomBeaconStaking _staking,
        DKGValidator _dkgValidator
    ) {
        sortitionPool = _sortitionPool;
        tToken = _tToken;
        staking = _staking;

        // TODO: revisit all initial values
        callbackGasLimit = 50000;
        groupCreationFrequency = 5;

        dkgResultSubmissionReward = 1000e18;
        sortitionPoolUnlockingReward = 100e18;
        ineligibleOperatorNotifierReward = 0;
        maliciousDkgResultSlashingAmount = 50000e18;
        unauthorizedSigningSlashingAmount = 100e3 * 1e18;
        sortitionPoolRewardsBanDuration = 2 weeks;
        relayEntryTimeoutNotificationRewardMultiplier = 40;
        unauthorizedSigningNotificationRewardMultiplier = 50;
        dkgMaliciousResultNotificationRewardMultiplier = 100;
        // slither-disable-next-line too-many-digits
        authorization.setMinimumAuthorization(100000 * 1e18);

        dkg.init(_sortitionPool, _dkgValidator);
        dkg.setResultChallengePeriodLength(11520); // ~48h assuming 15s block time
        dkg.setResultSubmissionEligibilityDelay(20);

        relay.initSeedEntry();
        relay.setRelayRequestFee(200e18);
        relay.setRelayEntrySubmissionEligibilityDelay(20);
        relay.setRelayEntryHardTimeout(5760); // ~24h assuming 15s block time
        relay.setRelayEntrySubmissionFailureSlashingAmount(1000e18);

        groups.setGroupLifetime(403200); // ~10 weeks assuming 15s block time
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
    /// @param _groupLifetime New group lifetime in blocks
    function updateGroupCreationParameters(
        uint256 _groupCreationFrequency,
        uint256 _groupLifetime
    ) external onlyOwner {
        groupCreationFrequency = _groupCreationFrequency;

        groups.setGroupLifetime(_groupLifetime);

        emit GroupCreationParametersUpdated(
            groupCreationFrequency,
            _groupLifetime
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
        dkg.setResultChallengePeriodLength(_dkgResultChallengePeriodLength);
        dkg.setResultSubmissionEligibilityDelay(
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
    /// @param _ineligibleOperatorNotifierReward New value of the ineligible
    ///        operator notifier reward.
    /// @param _sortitionPoolRewardsBanDuration New sortition pool rewards
    ///        ban duration in seconds.
    /// @param _relayEntryTimeoutNotificationRewardMultiplier New value of the
    ///        relay entry timeout notification reward multiplier.
    /// @param _unauthorizedSigningNotificationRewardMultiplier New value of the
    ///        unauthorized signing notification reward multiplier.
    /// @param _dkgMaliciousResultNotificationRewardMultiplier New value of the
    ///        DKG malicious result notification reward multiplier.
    function updateRewardParameters(
        uint256 _dkgResultSubmissionReward,
        uint256 _sortitionPoolUnlockingReward,
        uint256 _ineligibleOperatorNotifierReward,
        uint256 _sortitionPoolRewardsBanDuration,
        uint256 _relayEntryTimeoutNotificationRewardMultiplier,
        uint256 _unauthorizedSigningNotificationRewardMultiplier,
        uint256 _dkgMaliciousResultNotificationRewardMultiplier
    ) external onlyOwner {
        dkgResultSubmissionReward = _dkgResultSubmissionReward;
        sortitionPoolUnlockingReward = _sortitionPoolUnlockingReward;
        ineligibleOperatorNotifierReward = _ineligibleOperatorNotifierReward;
        sortitionPoolRewardsBanDuration = _sortitionPoolRewardsBanDuration;
        relayEntryTimeoutNotificationRewardMultiplier = _relayEntryTimeoutNotificationRewardMultiplier;
        unauthorizedSigningNotificationRewardMultiplier = _unauthorizedSigningNotificationRewardMultiplier;
        dkgMaliciousResultNotificationRewardMultiplier = _dkgMaliciousResultNotificationRewardMultiplier;
        emit RewardParametersUpdated(
            dkgResultSubmissionReward,
            sortitionPoolUnlockingReward,
            ineligibleOperatorNotifierReward,
            sortitionPoolRewardsBanDuration,
            relayEntryTimeoutNotificationRewardMultiplier,
            unauthorizedSigningNotificationRewardMultiplier,
            dkgMaliciousResultNotificationRewardMultiplier
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
    /// @param _unauthorizedSigningSlashingAmount New unauthorized signing
    ///        slashing amount
    function updateSlashingParameters(
        uint256 _relayEntrySubmissionFailureSlashingAmount,
        uint256 _maliciousDkgResultSlashingAmount,
        uint256 _unauthorizedSigningSlashingAmount
    ) external onlyOwner {
        relay.setRelayEntrySubmissionFailureSlashingAmount(
            _relayEntrySubmissionFailureSlashingAmount
        );
        maliciousDkgResultSlashingAmount = _maliciousDkgResultSlashingAmount;
        unauthorizedSigningSlashingAmount = _unauthorizedSigningSlashingAmount;
        emit SlashingParametersUpdated(
            _relayEntrySubmissionFailureSlashingAmount,
            maliciousDkgResultSlashingAmount,
            unauthorizedSigningSlashingAmount
        );
    }

    /// @notice Withdraws rewards belonging to operators marked as ineligible
    ///         for sortition pool rewards.
    /// @dev Can be called only by the contract owner, which should be the
    ///      random beacon governance contract.
    /// @param recipient Recipient of withdrawn rewards.
    function withdrawIneligibleRewards(address recipient) external onlyOwner {
        sortitionPool.withdrawIneligible(recipient);
    }

    /// @notice Registers the caller in the sortition pool.
    function registerOperator() external {
        address operator = msg.sender;

        require(
            !sortitionPool.isOperatorInPool(operator),
            "Operator is already registered"
        );

        sortitionPool.insertOperator(
            operator,
            staking.eligibleStake(operator, address(this))
        );
    }

    /// @notice Updates the sortition pool status of the caller.
    function updateOperatorStatus() external {
        sortitionPool.updateOperatorStatus(
            msg.sender,
            staking.eligibleStake(msg.sender, address(this))
        );
    }

    /// @notice Triggers group selection if there are no active groups.
    function genesis() external {
        require(groups.numberOfActiveGroups() == 0, "Not awaiting genesis");

        dkg.lockState();
        dkg.start(
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
        dkg.submitResult(dkgResult);

        groups.addCandidateGroup(
            dkgResult.groupPubKey,
            dkgResult.members,
            dkgResult.misbehavedMembersIndices
        );
    }

    /// @notice Notifies about DKG timeout. Pays the sortition pool unlocking
    ///         reward to the notifier.
    function notifyDkgTimeout() external {
        dkg.notifyTimeout();

        transferDkgRewards(msg.sender, sortitionPoolUnlockingReward);

        dkg.complete();
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
        uint32[] memory misbehavedMembers = dkg.approveResult(dkgResult);

        transferDkgRewards(msg.sender, dkgResultSubmissionReward);

        if (misbehavedMembers.length > 0) {
            sortitionPool.setRewardIneligibility(
                misbehavedMembers,
                // solhint-disable-next-line not-rely-on-time
                block.timestamp + sortitionPoolRewardsBanDuration
            );
        }

        groups.activateCandidateGroup();
        dkg.complete();
    }

    /// @notice Challenges DKG result. If the submitted result is proved to be
    ///         invalid it reverts the DKG back to the result submission phase.
    ///         It removes a candidate group that was previously registered with
    ///         the DKG result submission.
    /// @param dkgResult Result to challenge. Must match the submitted result
    ///        stored during `submitDkgResult`.
    function challengeDkgResult(DKG.Result calldata dkgResult) external {
        (bytes32 maliciousResultHash, uint32 maliciousSubmitter) = dkg
            .challengeResult(dkgResult);

        uint256 slashingAmount = maliciousDkgResultSlashingAmount;
        address maliciousSubmitterAddresses = sortitionPool.getIDOperator(
            maliciousSubmitter
        );

        groups.popCandidateGroup();

        emit DkgMaliciousResultSlashed(
            maliciousResultHash,
            slashingAmount,
            maliciousSubmitterAddresses
        );

        address[] memory operatorWrapper = new address[](1);
        operatorWrapper[0] = maliciousSubmitterAddresses;
        staking.seize(
            slashingAmount,
            dkgMaliciousResultNotificationRewardMultiplier,
            msg.sender,
            operatorWrapper
        );
    }

    /// @notice Check current group creation state.
    function getGroupCreationState() external view returns (DKG.State) {
        return dkg.currentState();
    }

    /// @notice Checks if DKG timed out. The DKG timeout period includes time required
    ///         for off-chain protocol execution and time for the result publication
    ///         for all group members. After this time result cannot be submitted
    ///         and DKG can be notified about the timeout.
    /// @return True if DKG timed out, false otherwise.
    function hasDkgTimedOut() external view returns (bool) {
        return dkg.hasDkgTimedOut();
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
            uint256(keccak256(AltBn128.g1Marshal(relay.previousEntry)))
        );

        relay.requestEntry(groupId);

        fundDkgRewardsPool(msg.sender, relay.relayRequestFee);

        callback.setCallbackContract(callbackContract);

        // If the current request should trigger group creation we need to lock
        // DKG state (lock sortition pool) to prevent operators from changing
        // its state before relay entry is known. That entry will be used as a
        // group selection seed.
        if (
            relay.requestCount % groupCreationFrequency == 0 &&
            dkg.currentState() == DKG.State.IDLE
        ) {
            dkg.lockState();
        }
    }

    /// @notice Creates a new relay entry.
    /// @param entry Group BLS signature over the previous entry.
    /// @param groupMembers Group member ids that participated in dkg (excluding IA/DQ).
    function submitRelayEntry(
        bytes calldata entry,
        uint32[] calldata groupMembers
    ) external {
        uint256 currentRequestId = relay.currentRequestID;

        Groups.Group storage group = groups.getGroup(
            relay.currentRequestGroupID
        );

        require(
            group.membersHash == keccak256(abi.encode(groupMembers)),
            "Invalid group members"
        );

        uint256 slashingAmount = relay.submitEntry(entry, group.groupPubKey);

        if (slashingAmount > 0) {
            address[] memory groupMembersIds = sortitionPool.getIDOperators(
                groupMembers
            );

            try staking.slash(slashingAmount, groupMembersIds) {
                // slither-disable-next-line reentrancy-events
                emit RelayEntryDelaySlashed(
                    currentRequestId,
                    slashingAmount,
                    groupMembersIds
                );
            } catch {
                // Should never happen but we want to ensure a non-critical path
                // failure from an external contract does not stop group members
                // from submitting a valid relay entry.
                emit RelayEntryDelaySlashingFailed(
                    currentRequestId,
                    slashingAmount,
                    groupMembersIds
                );
            }
        }

        // If DKG is awaiting a seed, that means the we should start the actual
        // group creation process.
        if (dkg.currentState() == DKG.State.AWAITING_SEED) {
            dkg.start(uint256(keccak256(entry)));
        }

        callback.executeCallback(uint256(keccak256(entry)), callbackGasLimit);
    }

    /// @notice Reports a relay entry timeout.
    /// @param groupMembers Group member ids that participated in dkg (excluding IA/DQ).
    function reportRelayEntryTimeout(uint32[] calldata groupMembers) external {
        uint64 groupId = relay.currentRequestGroupID;
        Groups.Group storage group = groups.getGroup(groupId);

        require(
            group.membersHash == keccak256(abi.encode(groupMembers)),
            "Invalid group members"
        );

        uint256 slashingAmount = relay
            .relayEntrySubmissionFailureSlashingAmount;
        address[] memory groupMembersIds = sortitionPool.getIDOperators(
            groupMembers
        );

        emit RelayEntryTimeoutSlashed(
            relay.currentRequestID,
            slashingAmount,
            groupMembersIds
        );

        staking.seize(
            slashingAmount,
            relayEntryTimeoutNotificationRewardMultiplier,
            msg.sender,
            groupMembersIds
        );

        groups.terminateGroup(groupId);
        groups.expireOldGroups();

        if (groups.numberOfActiveGroups() > 0) {
            groupId = groups.selectGroup(
                uint256(keccak256(AltBn128.g1Marshal(relay.previousEntry)))
            );
            relay.retryOnEntryTimeout(groupId);
        } else {
            relay.cleanupOnEntryTimeout();

            // If DKG is awaiting a seed, we should notify about its timeout to
            // avoid blocking the future group creation.
            if (dkg.currentState() == DKG.State.AWAITING_SEED) {
                dkg.notifySeedTimedOut();
            }
        }
    }

    /// @notice Reports unauthorized groups signing. Must provide a valid signature
    ///         of the sender's address as a message. Successful signature
    ///         verification means the private key has been leaked and all group
    ///         members should be punished by slashing their tokens. Group has
    ///         to be active or expired. Unauthorized signing cannot be reported
    ///         for a terminated group. In case of reporting unauthorized
    ///         signing for a terminated group, or when the signature is invalid,
    ///         function reverts.
    /// @param signedMsgSender Signature of the sender's address as a message.
    /// @param groupId Group that is being reported for leaking a private key.
    /// @param groupMembers Group member ids that participated in dkg (excluding IA/DQ).
    function reportUnauthorizedSigning(
        bytes memory signedMsgSender,
        uint64 groupId,
        uint32[] calldata groupMembers
    ) external {
        Groups.Group memory group = groups.getGroup(groupId);

        require(
            group.membersHash == keccak256(abi.encode(groupMembers)),
            "Invalid group members"
        );

        require(!group.terminated, "Group cannot be terminated");

        require(
            BLS.verifyBytes(
                group.groupPubKey,
                abi.encodePacked(msg.sender),
                signedMsgSender
            ),
            "Invalid signature"
        );

        groups.terminateGroup(groupId);

        address[] memory groupMembersIds = sortitionPool.getIDOperators(
            groupMembers
        );

        emit UnauthorizedSigningSlashed(
            groupId,
            unauthorizedSigningSlashingAmount,
            groupMembersIds
        );

        staking.seize(
            unauthorizedSigningSlashingAmount,
            unauthorizedSigningNotificationRewardMultiplier,
            msg.sender,
            groupMembersIds
        );
    }

    /// @notice Notifies about a failed group heartbeat. Using this function,
    ///         a majority of the group can decide about punishing specific
    ///         group members who failed to provide a heartbeat. If provided
    ///         claim is proved to be valid and signed by sufficient number
    ///         of group members, operators of members deemed as failed are
    ///         banned for sortition pool rewards for duration specified by
    ///         `sortitionPoolRewardsBanDuration` parameter. The sender of
    ///         the claim must be one of the claim signers. The sender is
    ///         rewarded from `heartbeatNotifierRewardsPool`. Exact reward
    ///         amount is multiplication of operators marked as ineligible
    ///         and `ineligibleOperatorNotifierReward` factor. This function
    ///         can be called only for active and non-terminated groups.
    /// @param claim Failure claim.
    /// @param nonce Current failed heartbeat nonce for given group. Must
    ///        be the same as the stored one.
    /// @param groupMembers Group member ids that participated in dkg (excluding IA/DQ).
    function notifyFailedHeartbeat(
        Heartbeat.FailureClaim calldata claim,
        uint256 nonce,
        uint32[] calldata groupMembers
    ) external {
        uint64 groupId = claim.groupId;

        require(nonce == failedHeartbeatNonce[groupId], "Invalid nonce");

        require(
            groups.isGroupActive(groupId),
            "Group must be active and non-terminated"
        );

        Groups.Group storage group = groups.getGroup(groupId);

        require(
            group.membersHash == keccak256(abi.encode(groupMembers)),
            "Invalid group members"
        );

        uint32[] memory ineligibleOperators = Heartbeat.verifyFailureClaim(
            sortitionPool,
            claim,
            group,
            nonce,
            groupMembers
        );

        failedHeartbeatNonce[groupId]++;

        emit HeartbeatFailed(groupId, nonce, msg.sender);

        sortitionPool.setRewardIneligibility(
            ineligibleOperators,
            // solhint-disable-next-line not-rely-on-time
            block.timestamp + sortitionPoolRewardsBanDuration
        );

        transferHeartbeatNotifierRewards(
            msg.sender,
            ineligibleOperatorNotifierReward * ineligibleOperators.length
        );
    }

    /// @notice Funds the DKG rewards pool.
    /// @param from Address of the funder. The funder must give a sufficient
    ///        allowance for this contract to make a successful call.
    /// @param value Token value transferred by the funder.
    function fundDkgRewardsPool(address from, uint256 value) public {
        dkgRewardsPool += value;
        tToken.safeTransferFrom(from, address(this), value);
    }

    /// @notice Funds the heartbeat notifier rewards pool.
    /// @param from Address of the funder. The funder must give a sufficient
    ///        allowance for this contract to make a successful call.
    /// @param value Token value transferred by the funder.
    function fundHeartbeatNotifierRewardsPool(address from, uint256 value)
        external
    {
        heartbeatNotifierRewardsPool += value;
        tToken.safeTransferFrom(from, address(this), value);
    }

    /// @notice Makes a transfer using DKG rewards pool.
    /// @param to Address of the recipient.
    /// @param value Token value transferred to the recipient.
    function transferDkgRewards(address to, uint256 value) internal {
        uint256 actualValue = Math.min(dkgRewardsPool, value);
        dkgRewardsPool -= actualValue;
        tToken.safeTransfer(to, actualValue);
    }

    /// @notice Makes a transfer using heartbeat notifier rewards pool.
    /// @param to Address of the recipient.
    /// @param value Token value transferred to the recipient.
    function transferHeartbeatNotifierRewards(address to, uint256 value)
        internal
    {
        uint256 actualValue = Math.min(heartbeatNotifierRewardsPool, value);
        heartbeatNotifierRewardsPool -= actualValue;
        tToken.safeTransfer(to, actualValue);
    }

    /// @notice Locks the state of group creation.
    /// @dev This function is meant to be used by test stubs which inherits
    ///      from this contract and needs to lock the DKG state arbitrarily.
    function dkgLockState() internal {
        dkg.lockState();
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

    /// @notice Group lifetime in blocks. When a group reached its lifetime, it
    ///         is no longer selected for new relay requests but may still be
    ///         responsible for submitting relay entry if relay request assigned
    ///         to that group is still pending.
    function groupLifetime() external view returns (uint256) {
        return groups.groupLifetime;
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

    /// @notice Selects a new group of operators based on the provided seed.
    ///         At least one operator has to be registered in the pool,
    ///         otherwise the function fails reverting the transaction.
    /// @param seed Number used to select operators to the group.
    /// @return IDs of selected group members.
    function selectGroup(bytes32 seed) external view returns (uint32[] memory) {
        return sortitionPool.selectGroup(DKG.groupSize, seed);
    }
}
