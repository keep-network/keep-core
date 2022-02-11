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

import "./libraries/DKG.sol";
import "./libraries/Wallets.sol";
import "./DKGValidator.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import {IRandomBeacon} from "@keep-network/random-beacon/contracts/RandomBeacon.sol";
import {IRandomBeaconConsumer} from "@keep-network/random-beacon/contracts/libraries/Callback.sol";

/// TODO: Add a dependency to `threshold-network/solidity-contracts` and use
/// IStaking interface from there.
interface IWalletStaking {
    function authorizedStake(address stakingProvider, address application)
        external
        view
        returns (uint256);

    function seize(
        uint96 amount,
        uint256 rewardMultiplier,
        address notifier,
        address[] memory stakingProviders
    ) external;
}

contract WalletRegistry is IRandomBeaconConsumer, Ownable {
    using DKG for DKG.Data;
    using Wallets for Wallets.Data;

    // Libraries data storages
    DKG.Data internal dkg;
    Wallets.Data internal wallets;

    // Address that is set as owner of all wallets. Only this address can request
    // new wallets creation and manage their state.
    address public walletOwner;

    /// @notice Slashing amount for supporting malicious DKG result. Every
    ///         DKG result submitted can be challenged for the time of DKG's
    ///         `resultChallengePeriodLength` parameter. If the DKG result submitted
    ///         is challenged and proven to be malicious, each operator who
    ///         signed the malicious result is slashed for
    ///         `maliciousDkgResultSlashingAmount`.
    uint96 public maliciousDkgResultSlashingAmount;

    /// @notice Percentage of the staking contract malicious behavior
    ///         notification reward which will be transferred to the notifier
    ///         reporting about a malicious DKG result. Notifiers are rewarded
    ///         from a notifiers treasury pool. For example, if
    ///         notification reward is 1000 and the value of the multiplier is
    ///         5, the notifier will receive: 5% of 1000 = 50 per each
    ///         operator affected.
    uint256 public maliciousDkgResultNotificationRewardMultiplier;

    // External dependencies

    SortitionPool public immutable sortitionPool;
    /// TODO: Add a dependency to `threshold-network/solidity-contracts` and use
    /// IStaking interface from there.
    IWalletStaking public immutable staking;
    IRandomBeacon public randomBeacon;

    // Events
    event DkgStarted(uint256 indexed seed);

    event DkgResultSubmitted(
        bytes32 indexed resultHash,
        uint256 indexed seed,
        DKG.Result result
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

    event DkgStateLocked();

    event DkgSeedTimedOut();

    event WalletCreated(
        bytes32 indexed publicKeyHash,
        bytes32 indexed dkgResultHash
    );

    event DkgMaliciousResultSlashed(
        bytes32 indexed resultHash,
        uint256 slashingAmount,
        address maliciousSubmitter
    );

    event DkgMaliciousResultSlashingFailed(
        bytes32 indexed resultHash,
        uint256 slashingAmount,
        address maliciousSubmitter
    );

    event RewardParametersUpdated(
        uint256 maliciousDkgResultNotificationRewardMultiplier
    );

    event SlashingParametersUpdated(uint256 maliciousDkgResultSlashingAmount);

    event DkgParametersUpdated(
        uint256 seedTimeout,
        uint256 resultChallengePeriodLength,
        uint256 resultSubmissionTimeout,
        uint256 resultSubmitterPrecedencePeriodLength
    );

    event WalletOwnerUpdated(address walletOwner);

    constructor(
        SortitionPool _sortitionPool,
        IWalletStaking _staking,
        DKGValidator _dkgValidator,
        IRandomBeacon _randomBeacon,
        address _walletOwner
    ) {
        sortitionPool = _sortitionPool;
        staking = _staking;
        randomBeacon = _randomBeacon;
        walletOwner = _walletOwner;

        // TODO: Implement governance for the parameters
        // TODO: revisit all initial values

        maliciousDkgResultSlashingAmount = 50000e18;
        maliciousDkgResultNotificationRewardMultiplier = 100;

        dkg.init(_sortitionPool, _dkgValidator);
        dkg.setSeedTimeout(1440); // ~6h assuming 15s block time // TODO: Verify value
        dkg.setResultChallengePeriodLength(11520); // ~48h assuming 15s block time
        dkg.setResultSubmissionTimeout(100 * 20); // TODO: Verify value
        dkg.setSubmitterPrecedencePeriodLength(20); // TODO: Verify value
    }

    // TODO: Update to governable params
    function updateRandomBeacon(IRandomBeacon _newRandomBeacon) external {
        randomBeacon = _newRandomBeacon;
    }

    /// @notice Updates the values of DKG parameters.
    /// @dev Can be called only by the contract owner, which should be the
    ///      wallet registry governance contract. The caller is responsible for
    ///      validating parameters.
    /// @param _seedTimeout New seed timeout.
    /// @param _resultChallengePeriodLength New DKG result challenge period
    ///        length
    /// @param _resultSubmissionTimeout New DKG result submission timeout
    /// @param _submitterPrecedencePeriodLength New submitter precedence period
    ///        length
    function updateDkgParameters(
        uint256 _seedTimeout,
        uint256 _resultChallengePeriodLength,
        uint256 _resultSubmissionTimeout,
        uint256 _submitterPrecedencePeriodLength
    ) external onlyOwner {
        dkg.setSeedTimeout(_seedTimeout);
        dkg.setResultChallengePeriodLength(_resultChallengePeriodLength);
        dkg.setResultSubmissionTimeout(_resultSubmissionTimeout);
        dkg.setSubmitterPrecedencePeriodLength(
            _submitterPrecedencePeriodLength
        );

        emit DkgParametersUpdated(
            _seedTimeout,
            _resultChallengePeriodLength,
            _resultSubmissionTimeout,
            _submitterPrecedencePeriodLength
        );
    }

    /// @notice Updates the values of reward parameters.
    /// @dev Can be called only by the contract owner, which should be the
    ///      wallet registry governance contract. The caller is responsible for
    ///      validating parameters.
    /// @param _maliciousDkgResultNotificationRewardMultiplier New value of the
    ///        DKG malicious result notification reward multiplier.
    function updateRewardParameters(
        uint256 _maliciousDkgResultNotificationRewardMultiplier
    ) external onlyOwner {
        maliciousDkgResultNotificationRewardMultiplier = _maliciousDkgResultNotificationRewardMultiplier;
        emit RewardParametersUpdated(
            maliciousDkgResultNotificationRewardMultiplier
        );
    }

    /// @notice Updates the values of slashing parameters.
    /// @dev Can be called only by the contract owner, which should be the
    ///      wallet registry governance contract. The caller is responsible for
    ///      validating parameters.
    /// @param _maliciousDkgResultSlashingAmount New malicious DKG result
    ///        slashing amount
    function updateSlashingParameters(uint96 _maliciousDkgResultSlashingAmount)
        external
        onlyOwner
    {
        maliciousDkgResultSlashingAmount = _maliciousDkgResultSlashingAmount;
        emit SlashingParametersUpdated(maliciousDkgResultSlashingAmount);
    }

    /// @notice Updates the values of the wallet parameters.
    /// @dev Can be called only by the contract owner, which should be the
    ///      wallet registry governance contract. The caller is responsible for
    ///      validating parameters.
    /// @param _walletOwner New wallet owner address.
    function updateWalletParameters(address _walletOwner) external onlyOwner {
        require(
            _walletOwner != address(0),
            "Wallet owner address cannot be zero"
        );

        walletOwner = _walletOwner;
        emit WalletOwnerUpdated(walletOwner);
    }

    /// @notice Registers the caller in the sortition pool.
    // TODO: Revisit on integration with Token Staking contract.
    function registerOperator() external {
        address operator = msg.sender;

        require(
            !sortitionPool.isOperatorInPool(operator),
            "Operator is already registered"
        );

        sortitionPool.insertOperator(
            operator,
            staking.authorizedStake(operator, address(this)) // FIXME: authorizedStake expects `stakingProvider` instead of `operator`
        );
    }

    /// @notice Updates the sortition pool status of the caller.
    /// @param operator Operator's address.
    // TODO: Revisit on integration with Token Staking contract.
    function updateOperatorStatus(address operator) external {
        sortitionPool.updateOperatorStatus(
            operator,
            staking.authorizedStake(msg.sender, address(this)) // FIXME: authorizedStake expects `stakingProvider` instead of `msg.sender`
        );
    }

    /// @notice Requests a new wallet creation.
    /// @dev Can be called only by the owner of wallets.
    ///      It locks the DKG and request a new random relay entry. It expects
    ///      that the DKG process will be started once a new random relay entry
    ///      gets generated.
    function requestNewWallet() external {
        require(msg.sender == walletOwner, "Caller is not the Wallet Owner");

        dkg.lockState();

        randomBeacon.requestRelayEntry(this);
    }

    /// @notice A callback that is executed once a new random relay entry gets
    ///         generated. It starts the DKG process.
    /// @dev Can be called only by the random beacon contract.
    /// @param randomRelayEntry Random relay entry.
    function __beaconCallback(uint256 randomRelayEntry, uint256) external {
        require(
            msg.sender == address(randomBeacon),
            "Caller is not the Random Beacon"
        );

        dkg.start(randomRelayEntry);
    }

    /// @notice Submits result of DKG protocol.
    ///         The DKG result consists of result submitting member index,
    ///         calculated group public key, bytes array of misbehaved members,
    ///         concatenation of signatures from group members, indices of members
    ///         corresponding to each signature and the list of group members.
    ///         The result is registered optimistically and waits for an approval.
    ///         The result can be challenged when it is believed to be incorrect.
    ///         The challenge verifies the registered result i.a. it checks if members
    ///         list corresponds to the expected set of members determined
    ///         by the sortition pool.
    /// @dev The message to be signed by each member is keccak256 hash of the
    ///      calculated group public key, misbehaved members indices and DKG
    ///      start block. The calculated hash should be prefixed with prefixed with
    ///      `\x19Ethereum signed message:\n` before signing, so the message to
    ///      sign is:
    ///      `\x19Ethereum signed message:\n${keccak256(groupPubKey,misbehavedIndices,startBlock)}`
    /// @param dkgResult DKG result.
    function submitDkgResult(DKG.Result calldata dkgResult) external {
        dkg.submitResult(dkgResult);
    }

    /// @notice Notifies about seed for DKG delivery timeout. It is expected
    ///         that a seed is delivered by the Random Beacon as a relay entry in a
    ///         callback function.
    function notifySeedTimeout() external {
        dkg.notifySeedTimeout();
        dkg.complete();
    }

    /// @notice Notifies about DKG timeout.
    function notifyDkgTimeout() external {
        dkg.notifyDkgTimeout();
        dkg.complete();
    }

    /// @notice Approves DKG result. Can be called when the challenge period for
    ///         the submitted result is finished. Considers the submitted result
    ///         as valid, pays reward to the approver, bans misbehaved group
    ///         members from the sortition pool rewards, and completes the group
    ///         creation by activating the candidate group. For the first
    ///         `resultSubmissionTimeout` blocks after the end of the
    ///         challenge period can be called only by the DKG result submitter.
    ///         After that time, can be called by anyone.
    ///         A new wallet based on the DKG result details.
    /// @param dkgResult Result to approve. Must match the submitted result
    ///        stored during `submitDkgResult`.
    function approveDkgResult(DKG.Result calldata dkgResult) external {
        uint32[] memory misbehavedMembers = dkg.approveResult(dkgResult);

        bytes32 publicKeyHash = keccak256(dkgResult.groupPubKey);

        wallets.addWallet(dkgResult.membersHash, publicKeyHash);

        emit WalletCreated(publicKeyHash, keccak256(abi.encode(dkgResult)));

        // TODO: Disable rewards for misbehavedMembers.
        //slither-disable-next-line redundant-statements
        misbehavedMembers;

        dkg.complete();
    }

    /// @notice Challenges DKG result. If the submitted result is proved to be
    ///         invalid it reverts the DKG back to the result submission phase.
    /// @param dkgResult Result to challenge. Must match the submitted result
    ///        stored during `submitDkgResult`.
    function challengeDkgResult(DKG.Result calldata dkgResult) external {
        (
            bytes32 maliciousDkgResultHash,
            uint32 maliciousDkgResultSubmitterId
        ) = dkg.challengeResult(dkgResult);

        address maliciousDkgResultSubmitterAddress = sortitionPool
            .getIDOperator(maliciousDkgResultSubmitterId);

        address[] memory operatorWrapper = new address[](1);
        operatorWrapper[0] = maliciousDkgResultSubmitterAddress;

        try
            staking.seize(
                maliciousDkgResultSlashingAmount,
                maliciousDkgResultNotificationRewardMultiplier,
                msg.sender,
                operatorWrapper
            )
        {
            // slither-disable-next-line reentrancy-events
            emit DkgMaliciousResultSlashed(
                maliciousDkgResultHash,
                maliciousDkgResultSlashingAmount,
                maliciousDkgResultSubmitterAddress
            );
        } catch {
            // Should never happen but we want to ensure a non-critical path
            // failure from an external contract does not stop the challenge
            // to complete.
            emit DkgMaliciousResultSlashingFailed(
                maliciousDkgResultHash,
                maliciousDkgResultSlashingAmount,
                maliciousDkgResultSubmitterAddress
            );
        }
    }

    /// @notice Checks if DKG result is valid for the current DKG.
    /// @param result DKG result.
    /// @return True if the result is valid. If the result is invalid it returns
    ///         false and an error message.
    function isDkgResultValid(DKG.Result calldata result)
        external
        view
        returns (bool, string memory)
    {
        return dkg.isResultValid(result);
    }

    /// @notice Check current wallet creation state.
    function getWalletCreationState() external view returns (DKG.State) {
        return dkg.currentState();
    }

    /// @notice Checks if seed awaiting timed out.
    /// @return True if seed awaiting timed out, false otherwise.
    function hasSeedTimedOut() external view returns (bool) {
        return dkg.hasSeedTimedOut();
    }

    /// @notice Checks if DKG timed out. The DKG timeout period includes time required
    ///         for off-chain protocol execution and time for the result publication
    ///         for all group members. After this time result cannot be submitted
    ///         and DKG can be notified about the timeout.
    /// @return True if DKG timed out, false otherwise.
    function hasDkgTimedOut() external view returns (bool) {
        return dkg.hasDkgTimedOut();
    }

    function getWallet(bytes32 publicKeyHash)
        external
        view
        returns (Wallets.Wallet memory)
    {
        return wallets.registry[publicKeyHash];
    }

    /// @notice Checks if a wallet with the given public key hash is registered.
    /// @param publicKeyHash Wallet's public key hash.
    /// @return True if wallet is registered, false otherwise.
    function isWalletRegistered(bytes32 publicKeyHash)
        external
        view
        returns (bool)
    {
        return wallets.isWalletRegistered(publicKeyHash);
    }

    /// @notice Retrieves dkg parameters that were set in DKG library.
    function dkgParameters() external view returns (DKG.Parameters memory) {
        return dkg.parameters;
    }
}
