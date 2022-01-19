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
import "./Wallet.sol";
import "./DKGValidator.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "@thesis/solidity-contracts/contracts/clone/CloneFactory.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// TODO: Add a dependency to `threshold-network/solidity-contracts` and use
/// IStaking interface from there.
interface IWalletStaking {
    function eligibleStake(address operator, address operatorContract)
        external
        view
        returns (uint256);
}

contract WalletFactory is CloneFactory, Ownable {
    using DKG for DKG.Data;

    // Libraries data storages
    DKG.Data internal dkg;

    // Reference to a wallet contract that is under creation. Once wallet creation
    // completes the value is reset. The reference is set on DKG result submission
    // and is kept until the current wallet creation process completion.
    Wallet public wallet;

    uint256 public constant relayEntry = 12345; // TODO: get value from Random Beacon

    // External dependencies

    SortitionPool public sortitionPool;
    IERC20 public tToken;
    /// TODO: Add a dependency to `threshold-network/solidity-contracts` and use
    /// IStaking interface from there.
    IWalletStaking public staking;
    // Holds the address of the wallet contract which will be used as a master
    // contract for cloning.
    Wallet public immutable masterWallet;

    // Events

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

    event DkgStateLocked();

    event DkgSeedTimedOut();

    event WalletCreated(
        address indexed walletAddress,
        bytes32 indexed dkgResultHash
    );

    event WalletRemoved(address indexed walletAddress);

    event WalletActivated(address indexed walletAddress);

    constructor(
        SortitionPool _sortitionPool,
        IERC20 _tToken,
        IWalletStaking _staking,
        DKGValidator _dkgValidator,
        Wallet _masterWallet
    ) {
        sortitionPool = _sortitionPool;
        tToken = _tToken;
        staking = _staking;
        masterWallet = _masterWallet;

        dkg.init(_sortitionPool, _dkgValidator);
        // TODO: Implement governance for the parameters
        dkg.setResultChallengePeriodLength(11520); // ~48h assuming 15s block time
        dkg.setResultSubmissionEligibilityDelay(20);
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
    /// @dev Can be called only by the owner of the Wallet Factory.
    ///      It starts a new DKG process.
    function requestNewWallet() external onlyOwner {
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

        address clonedWalletAddress = createClone(address(masterWallet));
        require(
            clonedWalletAddress != address(0),
            "Cloned wallet address is 0"
        );

        emit WalletCreated(
            clonedWalletAddress,
            keccak256(abi.encode(dkgResult))
        );

        // We expect `dkg.submitResult` function verifies the current state of
        // the wallet creation process. It is expected that at this point the
        // wallet reference is not set, hence we won't be overwriting any previous
        // value.
        wallet = Wallet(clonedWalletAddress);

        wallet.init(
            address(this),
            hashGroupMembers(
                dkgResult.members,
                dkgResult.misbehavedMembersIndices
            ),
            keccak256(dkgResult.groupPubKey)
        );
    }

    /// @notice Notifies about DKG timeout. Pays the sortition pool unlocking
    ///         reward to the notifier.
    function notifyDkgTimeout() external {
        dkg.notifyTimeout();

        // TODO: Implement transferDkgRewards
        // transferDkgRewards(msg.sender, sortitionPoolUnlockingReward);

        //slither-disable-next-line reentrancy-benign
        dkg.complete();

        delete wallet;
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

        emit WalletActivated(address(wallet));

        //slither-disable-next-line reentrancy-no-eth
        wallet.activate();

        // TODO: Transfer Wallet's ownership to WalletManager

        // TODO: Transfer DKG rewards and disable rewards for misbehavedMembers.
        //slither-disable-next-line redundant-statements
        misbehavedMembers;

        // TODO: Let the Wallet Manager know that a new wallet was created successfully.

        dkg.complete();
        delete wallet;
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

        emit WalletRemoved(address(wallet));

        // TODO: Implement slashing.
        //slither-disable-next-line redundant-statements
        maliciousResultHash;
        //slither-disable-next-line redundant-statements
        maliciousSubmitter;

        delete wallet;
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

    /// @notice Hash group members that actively participated in a group signing
    ///         key generation. This function filters out IA/DQ members before
    ///         hashing.
    /// @param members Group member addresses as outputted by the group selection
    ///        protocol.
    /// @param misbehavedMembersIndices Array of misbehaved (disqualified or
    ///        inactive) group members. Indices reflect positions
    ///        of members in the group as outputted by the group selection
    ///        protocol. Indices must be in ascending order. The order can be verified
    ///        during the DKG challege phase in DKGValidator contract.
    /// @return Group members hash.
    function hashGroupMembers(
        uint32[] calldata members,
        uint8[] calldata misbehavedMembersIndices
    ) private pure returns (bytes32) {
        if (misbehavedMembersIndices.length > 0) {
            // members that generated a group signing key
            uint32[] memory groupMembers = new uint32[](
                members.length - misbehavedMembersIndices.length
            );
            uint256 k = 0; // misbehaved members counter
            uint256 j = 0; // group members counter
            for (uint256 i = 0; i < members.length; i++) {
                // misbehaved member indices start from 1, so we need to -1 on misbehaved
                if (i != misbehavedMembersIndices[k] - 1) {
                    groupMembers[j] = members[i];
                    j++;
                } else if (k < misbehavedMembersIndices.length - 1) {
                    k++;
                }
            }

            return keccak256(abi.encode(groupMembers));
        }

        return keccak256(abi.encode(members));
    }
}
