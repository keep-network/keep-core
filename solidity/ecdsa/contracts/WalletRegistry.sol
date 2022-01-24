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
    function eligibleStake(address operator, address operatorContract)
        external
        view
        returns (uint256);
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
        uint256 indexed walletID,
        bytes32 indexed dkgResultHash
    );

    event SignatureRequested(uint256 indexed walletID, bytes32 indexed digest);

    event SignatureSubmitted(
        uint256 indexed walletID,
        bytes32 indexed digest,
        Wallets.Signature signature
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

        dkg.init(_sortitionPool, _dkgValidator);
        // TODO: Implement governance for the parameters
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
    ///         A candidate wallet is registered based on the submitted DKG result
    ///         details.
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
    /// @dev It transfers Wallet's ownership to the Wallet Factory owner.
    /// @param dkgResult Result to approve. Must match the submitted result
    ///        stored during `submitDkgResult`.
    function approveDkgResult(DKG.Result calldata dkgResult) external {
        uint32[] memory misbehavedMembers = dkg.approveResult(dkgResult);

        uint256 walletID = wallets.addWallet(
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
    ///         It removes a candidate group that was previously registered with
    ///         the DKG result submission.
    /// @param dkgResult Result to challenge. Must match the submitted result
    ///        stored during `submitDkgResult`.
    function challengeDkgResult(DKG.Result calldata dkgResult) external {
        (bytes32 maliciousResultHash, uint32 maliciousSubmitter) = dkg
            .challengeResult(dkgResult);

        // TODO: Implement slashing.
        //slither-disable-next-line redundant-statements
        maliciousResultHash;
        //slither-disable-next-line redundant-statements
        maliciousSubmitter;
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
    /// @param walletID ID of a wallet that should calculate a signature.
    /// @param digest Digest to sign.
    function requestSignature(uint256 walletID, bytes32 digest) external {
        require(msg.sender == walletOwner, "Caller is not the Wallet Owner");

        wallets.requestSignature(walletID, digest);
    }

    /// @notice Submits a calculated signature for the digest that is currently
    ///         under signing.
    /// @param walletID ID of a wallet that should calculate a signature.
    /// @param publicKey Group public key.
    /// @param signature Calculated signature.
    function submitSignature(
        uint256 walletID,
        bytes calldata publicKey,
        Wallets.Signature calldata signature
    ) external {
        wallets.submitSignature(walletID, publicKey, signature);
    }

    function getWallet(uint256 walletID)
        external
        view
        returns (Wallets.Wallet memory)
    {
        return wallets.registry[walletID];
    }

    function getWalletCounter() external view returns (uint256) {
        return wallets.walletCounter;
    }
}
