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

/// TODO: Add a dependency to `threshold-network/solidity-contracts` and use
/// IStaking interface from there.
interface IWalletStaking {
    function eligibleStake(address operator, address operatorContract)
        external
        view
        returns (uint256);
}

contract WalletFactory is CloneFactory {
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

    event WalletCreated(address walletAddress);

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
    // Consider who the caller will be: anyone or onlyOwner?
    function createNewWallet() external {
        dkg.lockState();
        dkg.start(
            uint256(keccak256(abi.encodePacked(relayEntry, block.number)))
        );
    }

    function submitDkgResult(DKG.Result calldata dkgResult) external {
        dkg.submitResult(dkgResult);

        // FIXME: We use all members now. We should filter out DQ/IA members and
        // have just the active members that will be part of the wallet.
        // See https://github.com/keep-network/keep-core/pull/2768
        uint32[] memory walletMembers = dkgResult.members;

        address clonedWalletAddress = createClone(address(masterWallet));
        require(
            clonedWalletAddress != address(0),
            "Cloned wallet address is 0"
        );

        Wallet wallet = Wallet(clonedWalletAddress);

        wallets.push(wallet);

        wallet.init(walletMembers);

        emit WalletCreated(address(wallet));
    }

    function approveDkgResult(DKG.Result calldata dkgResult) external {
        uint32[] memory misbehavedMembers = dkg.approveResult(dkgResult);

        // TODO: Transfer DKG rewards and disable rewards for misbehavedMembers.
        misbehavedMembers;

        Wallet latestWallet = wallets[wallets.length - 1];

        latestWallet.activate();
        dkg.complete();
    }

    function challengeDkgResult(DKG.Result calldata dkgResult) external {
        require(
            wallets[wallets.length - 1].activationBlockNumber() == 0,
            "The latest registered wallet was already activated"
        );

        (bytes32 maliciousResultHash, uint32 maliciousSubmitter) = dkg
            .challengeResult(dkgResult);

        wallets.pop();

        // TODO: Implement slashing.
        maliciousResultHash;
        maliciousSubmitter;
    }

    // TODO: Add timeouts
}
