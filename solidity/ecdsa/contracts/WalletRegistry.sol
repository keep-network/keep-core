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

contract WalletRegistry is Ownable {
    using DKG for DKG.Data;
    using Wallets for Wallets.Data;

    // Libraries data storages
    DKG.Data internal dkg;
    Wallets.Data internal wallets;

    // Address that is set as owner of all wallets. Only this address can request
    // new wallets creation, manage their state or request signatures from wallets.
    address public walletOwner;

    uint256 public constant relayEntry = 12345; // TODO: get value from Random Beacon

    /// @notice Slashing amount for supporting malicious DKG result. Every
    ///         DKG result submitted can be challenged for the time of
    ///         `dkgResultChallengePeriodLength`. If the DKG result submitted
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
    IERC20 public immutable tToken;
    /// TODO: Add a dependency to `threshold-network/solidity-contracts` and use
    /// IStaking interface from there.
    IWalletStaking public immutable staking;

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
        bytes32 indexed walletID,
        bytes32 indexed dkgResultHash
    );

    event SignatureRequested(bytes32 indexed walletID, bytes32 indexed digest);

    event SignatureSubmitted(
        bytes32 indexed walletID,
        bytes32 indexed digest,
        Wallets.Signature signature
    );

    event DkgMaliciousResultSlashed(
        bytes32 indexed resultHash,
        uint256 slashingAmount,
        address maliciousSubmitter
    );

    constructor(
        SortitionPool _sortitionPool,
        IERC20 _tToken,
        IWalletStaking _staking,
        DKGValidator _dkgValidator,
        address _walletOwner
    ) {
        sortitionPool = _sortitionPool;
        tToken = _tToken;
        staking = _staking;
        walletOwner = _walletOwner;

        // TODO: Implement governance for the parameters
        // TODO: revisit all initial values

        maliciousDkgResultSlashingAmount = 50000e18;
        maliciousDkgResultNotificationRewardMultiplier = 100;

        dkg.init(_sortitionPool, _dkgValidator);
        dkg.setResultChallengePeriodLength(11520); // ~48h assuming 15s block time
        dkg.setResultSubmissionEligibilityDelay(20);
    }

    // TODO: Update to governable params
    function updateDkgParams(
        uint256 newResultChallengePeriodLength,
        uint256 newResultSubmissionEligibilityDelay
    ) external {
        dkg.setResultChallengePeriodLength(newResultChallengePeriodLength);
        dkg.setResultSubmissionEligibilityDelay(
            newResultSubmissionEligibilityDelay
        );
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
            staking.authorizedStake(operator, address(this)) // FIXME: authorizedStake expects `stakingProvider` instead of `operator`
        );
    }

    /// @notice Updates the sortition pool status of the caller.
    function updateOperatorStatus() external {
        sortitionPool.updateOperatorStatus(
            msg.sender,
            staking.authorizedStake(msg.sender, address(this)) // FIXME: authorizedStake expects `stakingProvider` instead of `msg.sender`
        );
    }

    /// @notice Requests a new wallet creation.
    /// @dev Can be called only by the owner of wallets.
    ///      It starts a new DKG process.
    function requestNewWallet() external {
        require(msg.sender == walletOwner, "Caller is not the Wallet Owner");

        dkg.lockState();
        dkg.start(
            uint256(keccak256(abi.encodePacked(relayEntry, block.number)))
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

    /// @notice Notifies about DKG timeout. Pays the sortition pool unlocking
    ///         reward to the notifier.
    function notifyDkgTimeout() external {
        dkg.notifyTimeout();

        // TODO: Implement transferDkgRewards
        // transferDkgRewards(msg.sender, sortitionPoolUnlockingReward);

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
    ///         A new wallet based on the DKG result details.
    /// @param dkgResult Result to approve. Must match the submitted result
    ///        stored during `submitDkgResult`.
    function approveDkgResult(DKG.Result calldata dkgResult) external {
        uint32[] memory misbehavedMembers = dkg.approveResult(dkgResult);

        // TODO: Revisit at the end of ECDSA Wallets implementation to see if
        // usage was clarified and we can simplify it, as walletID is assumed to
        // be the same as groupPubKeyHash.
        bytes32 walletID = wallets.addWallet(
            dkgResult.membersHash,
            keccak256(dkgResult.groupPubKey)
        );

        emit WalletCreated(walletID, keccak256(abi.encode(dkgResult)));

        // TODO: Transfer DKG rewards and disable rewards for misbehavedMembers.
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

        emit DkgMaliciousResultSlashed(
            maliciousDkgResultHash,
            maliciousDkgResultSlashingAmount,
            maliciousDkgResultSubmitterAddress
        );

        address[] memory operatorWrapper = new address[](1);
        operatorWrapper[0] = maliciousDkgResultSubmitterAddress;

        staking.seize(
            maliciousDkgResultSlashingAmount,
            maliciousDkgResultNotificationRewardMultiplier,
            msg.sender,
            operatorWrapper
        );
    }

    /// @notice Check current wallet creation state.
    function getWalletCreationState() external view returns (DKG.State) {
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

    /// @notice Requests a new signature.
    /// @param publicKeyHash ID of a wallet that should calculate a signature.
    /// @param digest Digest to sign.
    function requestSignature(bytes32 publicKeyHash, bytes32 digest) external {
        require(msg.sender == walletOwner, "Caller is not the Wallet Owner");

        wallets.requestSignature(publicKeyHash, digest);
    }

    /// @notice Submits a calculated signature for the digest that is currently
    ///         under signing.
    /// @param walletID ID of a wallet that should calculate a signature.
    /// @param signature Calculated signature.
    function submitSignature(
        bytes32 walletID,
        Wallets.Signature calldata signature
    ) external {
        wallets.submitSignature(walletID, signature);
    }

    function getWallet(bytes32 walletID)
        external
        view
        returns (Wallets.Wallet memory)
    {
        return wallets.registry[walletID];
    }
}
