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

    Wallet[] public wallets;

    uint256 public relayEntry = 12345; // TODO: get value from Random Beacon

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

    // External dependencies

    SortitionPool public sortitionPool;
    IERC20 public tToken;
    /// TODO: Add a dependency to `threshold-network/solidity-contracts` and use
    /// IStaking interface from there.
    IWalletStaking public staking;
    // Holds the address of the wallet contract which will be used as a master
    // contract for cloning.
    Wallet public immutable masterWallet;

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

    // TODO: Revisit to implement mechanism for a fresh wallet creation.
    function requestNewWallet() external onlyOwner {
        dkg.lockState();
        dkg.start(
            uint256(keccak256(abi.encodePacked(relayEntry, block.number)))
        );
    }

    function submitDkgResult(DKG.Result calldata dkgResult) external {
        dkg.submitResult(dkgResult);

        address clonedWalletAddress = createClone(address(masterWallet));
        require(
            clonedWalletAddress != address(0),
            "Cloned wallet address is 0"
        );

        Wallet wallet = Wallet(clonedWalletAddress);

        wallets.push(wallet);

        wallet.init(
            address(this),
            hashGroupMembers(
                dkgResult.members,
                dkgResult.misbehavedMembersIndices
            ),
            keccak256(dkgResult.groupPubKey)
        );

        emit WalletCreated(address(wallet), keccak256(abi.encode(dkgResult)));
    }

    /// @notice Notifies about DKG timeout. Pays the sortition pool unlocking
    ///         reward to the notifier.
    function notifyDkgTimeout() external {
        dkg.notifyTimeout();

        // TODO: Implement transferDkgRewards
        // transferDkgRewards(msg.sender, sortitionPoolUnlockingReward);

        dkg.complete();
    }

    function approveDkgResult(DKG.Result calldata dkgResult) external {
        uint32[] memory misbehavedMembers = dkg.approveResult(dkgResult);

        Wallet latestWallet = wallets[wallets.length - 1];
        latestWallet.activate();

        // TODO: Transfer Wallet's ownership to WalletManager

        // TODO: Transfer DKG rewards and disable rewards for misbehavedMembers.
        misbehavedMembers;

        emit WalletActivated(address(latestWallet));

        dkg.complete();
    }

    function challengeDkgResult(DKG.Result calldata dkgResult) external {
        Wallet latestWallet = wallets[wallets.length - 1];
        require(
            latestWallet.activationBlockNumber() == 0,
            "The latest registered wallet was already activated"
        );

        (bytes32 maliciousResultHash, uint32 maliciousSubmitter) = dkg
            .challengeResult(dkgResult);

        wallets.pop();
        emit WalletRemoved(address(latestWallet));

        // TODO: Implement slashing.
        maliciousResultHash;
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

    /// @notice Returns registered wallets.
    function getWallets() external view returns (Wallet[] memory) {
        return wallets;
    }

    // TODO: Add timeouts

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
