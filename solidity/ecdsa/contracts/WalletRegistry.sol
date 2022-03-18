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

import "./api/IWalletRegistry.sol";
import "./api/IWalletOwner.sol";
import "./libraries/EcdsaAuthorization.sol";
import "./libraries/EcdsaDkg.sol";
import "./libraries/Wallets.sol";
import "./EcdsaDkgValidator.sol";

import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "@keep-network/random-beacon/contracts/api/IRandomBeacon.sol";
import "@keep-network/random-beacon/contracts/api/IRandomBeaconConsumer.sol";

import "@threshold-network/solidity-contracts/contracts/staking/IApplication.sol";
import "@threshold-network/solidity-contracts/contracts/staking/IStaking.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract WalletRegistry is
    IWalletRegistry,
    IRandomBeaconConsumer,
    IApplication,
    Ownable
{
    using EcdsaAuthorization for EcdsaAuthorization.Data;
    using EcdsaDkg for EcdsaDkg.Data;
    using Wallets for Wallets.Data;

    // Libraries data storages
    EcdsaAuthorization.Data internal authorization;
    EcdsaDkg.Data internal dkg;
    Wallets.Data internal wallets;

    // Address that is set as owner of all wallets. Only this address can request
    // new wallets creation and manage their state.
    IWalletOwner public walletOwner;

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
    IStaking public immutable staking;
    IRandomBeacon public randomBeacon;

    // Events
    event DkgStarted(uint256 indexed seed);

    event DkgResultSubmitted(
        bytes32 indexed resultHash,
        uint256 indexed seed,
        EcdsaDkg.Result result
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

    event AuthorizationParametersUpdated(
        uint96 minimumAuthorization,
        uint64 authorizationDecreaseDelay
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

    event RandomBeaconUpgraded(address randomBeacon);

    event WalletOwnerUpdated(address walletOwner);

    event OperatorRegistered(
        address indexed stakingProvider,
        address indexed operator
    );

    event AuthorizationIncreased(
        address indexed stakingProvider,
        address indexed operator,
        uint96 fromAmount,
        uint96 toAmount
    );

    event AuthorizationDecreaseRequested(
        address indexed stakingProvider,
        address indexed operator,
        uint96 fromAmount,
        uint96 toAmount,
        uint64 decreasingAt
    );

    event AuthorizationDecreaseApproved(address indexed stakingProvider);

    event InvoluntaryAuthorizationDecreaseFailed(
        address indexed stakingProvider,
        address indexed operator,
        uint96 fromAmount,
        uint96 toAmount
    );

    event OperatorJoinedSortitionPool(
        address indexed stakingProvider,
        address indexed operator
    );

    event OperatorStatusUpdated(
        address indexed stakingProvider,
        address indexed operator
    );

    modifier onlyStakingContract() {
        require(
            msg.sender == address(staking),
            "Caller is not the staking contract"
        );
        _;
    }

    /// @notice Reverts if called not by the Wallet Owner.
    modifier onlyWalletOwner() {
        require(
            msg.sender == address(walletOwner),
            "Caller is not the Wallet Owner"
        );
        _;
    }

    constructor(
        SortitionPool _sortitionPool,
        IStaking _staking,
        EcdsaDkgValidator _ecdsaDkgValidator,
        IRandomBeacon _randomBeacon
    ) {
        sortitionPool = _sortitionPool;
        staking = _staking;
        randomBeacon = _randomBeacon;

        // TODO: Implement governance for the parameters
        // TODO: revisit all initial values

        // slither-disable-next-line too-many-digits
        authorization.setMinimumAuthorization(400000e18); // 400k T
        authorization.setAuthorizationDecreaseDelay(5184000); // 60 days
        maliciousDkgResultSlashingAmount = 50000e18;
        maliciousDkgResultNotificationRewardMultiplier = 100;

        dkg.init(_sortitionPool, _ecdsaDkgValidator);
        dkg.setSeedTimeout(1440); // ~6h assuming 15s block time // TODO: Verify value
        dkg.setResultChallengePeriodLength(11520); // ~48h assuming 15s block time
        dkg.setResultSubmissionTimeout(100 * 20); // TODO: Verify value
        dkg.setSubmitterPrecedencePeriodLength(20); // TODO: Verify value
    }

    /// @notice Used by staking provider to set operator address that will
    ///         operate ECDSA node. The given staking provider can set operator
    ///         address only one time. The operator address can not be changed
    ///         and must be unique. Reverts if the operator is already set for
    ///         the staking provider or if the operator address is already in
    ///         use. Reverts if there is a pending authorization decrease for
    ///         the staking provider.
    function registerOperator(address operator) external {
        authorization.registerOperator(operator);
    }

    /// @notice Lets the operator join the sortition pool. The operator address
    ///         must be known - before calling this function, it has to be
    ///         appointed by the staking provider by calling `registerOperator`.
    ///         Also, the operator must have the minimum authorization required
    ///         by ECDSA. Function reverts if there is no minimum stake
    ///         authorized or if the operator is not known. If there was an
    ///         authorization decrease requested, it is activated by starting
    ///         the authorization decrease delay.
    function joinSortitionPool() external {
        authorization.joinSortitionPool(staking, sortitionPool);
    }

    /// @notice Updates status of the operator in the sortition pool. If there
    ///         was an authorization decrease requested, it is activated by
    ///         starting the authorization decrease delay.
    ///         Function reverts if the operator is not known.
    function updateOperatorStatus(address operator) external {
        authorization.updateOperatorStatus(staking, sortitionPool, operator);
    }

    /// @notice Used by T staking contract to inform the application that the
    ///         authorized stake amount for the given staking provider increased.
    ///
    ///         Reverts if the authorization amount is below the minimum.
    ///
    ///         The function is not updating the sortition pool. Sortition pool
    ///         state needs to be updated by the operator with a call to
    ///         `joinSortitionPool` or `updateOperatorStatus`.
    ///
    /// @dev Can only be called by T staking contract.
    function authorizationIncreased(
        address stakingProvider,
        uint96 fromAmount,
        uint96 toAmount
    ) external onlyStakingContract {
        authorization.authorizationIncreased(
            stakingProvider,
            fromAmount,
            toAmount
        );
    }

    /// @notice Used by T staking contract to inform the application that the
    ///         authorization decrease for the given staking provider has been
    ///         requested.
    ///
    ///         Reverts if the amount after deauthorization would be non-zero
    ///         and lower than the minimum authorization.
    ///
    ///         If the operator is not known (`registerOperator` was not called)
    ///         it lets to `approveAuthorizationDecrease` immediatelly. If the
    ///         operator is known (`registerOperator` was called), the operator
    ///         needs to update state of the sortition pool with a call to
    ///         `joinSortitionPool` or `updateOperatorStatus`. After the
    ///         sortition pool state is in sync, authorization decrease delay
    ///         starts.
    ///
    ///         After authorization decrease delay passes, authorization
    ///         decrease request needs to be approved with a call to
    ///         `approveAuthorizationDecrease` function.
    ///
    ///         If there is a pending authorization decrease request, it is
    ///         overwritten.
    ///
    /// @dev Should only be callable by T staking contract.
    function authorizationDecreaseRequested(
        address stakingProvider,
        uint96 fromAmount,
        uint96 toAmount
    ) external onlyStakingContract {
        authorization.authorizationDecreaseRequested(
            stakingProvider,
            fromAmount,
            toAmount
        );
    }

    /// @notice Approves the previously registered authorization decrease
    ///         request. Reverts if authorization decrease delay have not passed
    ///         yet or if the authorization decrease was not requested for the
    ///         given staking provider.
    function approveAuthorizationDecrease(address stakingProvider) external {
        authorization.approveAuthorizationDecrease(staking, stakingProvider);
    }

    /// @notice Used by T staking contract to inform the application the
    ///         authorization has been decreased for the given staking provider
    ///         involuntarily, as a result of slashing.
    ///
    ///         If the operator is not known (`registerOperator` was not called)
    ///         the function does nothing. The operator was never in a sortition
    ///         pool so there is nothing to update.
    ///
    ///         If the operator is known, sortition pool is unlocked, and the
    ///         operator is in the sortition pool, the sortition pool state is
    ///         updated. If the sortition pool is locked, update needs to be
    ///         postponed. Every other staker is incentivized to call
    ///         `updateOperatorStatus` for the problematic operator to increase
    ///         their own rewards in the pool.
    function involuntaryAuthorizationDecrease(
        address stakingProvider,
        uint96 fromAmount,
        uint96 toAmount
    ) external onlyStakingContract {
        authorization.involuntaryAuthorizationDecrease(
            staking,
            sortitionPool,
            stakingProvider,
            fromAmount,
            toAmount
        );
    }

    /// @notice Updates address of the Random Beacon.
    /// @dev Can be called only by the contract owner, which should be the
    ///      wallet registry governance contract. The caller is responsible for
    ///      validating parameters.
    /// @param _randomBeacon Random Beacon address.
    function upgradeRandomBeacon(IRandomBeacon _randomBeacon)
        external
        onlyOwner
    {
        randomBeacon = _randomBeacon;
        emit RandomBeaconUpgraded(address(_randomBeacon));
    }

    /// @notice Updates the values of authorization parameters.
    /// @dev Can be called only by the contract owner, which should be the
    ///      wallet registry governance contract. The caller is responsible for
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

        // slither-disable-next-line reentrancy-events
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
    ///      validating parameters. The wallet owner has to implement `IWalletOwner`
    ///      interface.
    /// @param _walletOwner New wallet owner address.
    function updateWalletOwner(IWalletOwner _walletOwner) external onlyOwner {
        require(
            address(_walletOwner) != address(0),
            "Wallet owner address cannot be zero"
        );

        walletOwner = _walletOwner;
        emit WalletOwnerUpdated(address(_walletOwner));
    }

    /// @notice Requests a new wallet creation.
    /// @dev Can be called only by the owner of wallets.
    ///      It locks the DKG and request a new relay entry. It expects
    ///      that the DKG process will be started once a new relay entry
    ///      gets generated.
    function requestNewWallet() external onlyWalletOwner {
        dkg.lockState();

        randomBeacon.requestRelayEntry(this);
    }

    /// @notice A callback that is executed once a new relay entry gets
    ///         generated. It starts the DKG process.
    /// @dev Can be called only by the random beacon contract.
    /// @param relayEntry Relay entry.
    function __beaconCallback(uint256 relayEntry, uint256) external {
        require(
            msg.sender == address(randomBeacon),
            "Caller is not the Random Beacon"
        );

        dkg.start(relayEntry);
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
    function submitDkgResult(EcdsaDkg.Result calldata dkgResult) external {
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
    function approveDkgResult(EcdsaDkg.Result calldata dkgResult) external {
        uint32[] memory misbehavedMembers = dkg.approveResult(dkgResult);

        (bytes32 walletID, bytes32 publicKeyX, bytes32 publicKeyY) = wallets
            .addWallet(dkgResult.membersHash, dkgResult.groupPubKey);

        emit WalletCreated(walletID, keccak256(abi.encode(dkgResult)));

        // TODO: Disable rewards for misbehavedMembers.
        //slither-disable-next-line redundant-statements
        misbehavedMembers;

        walletOwner.__ecdsaWalletCreatedCallback(
            walletID,
            publicKeyX,
            publicKeyY
        );

        dkg.complete();
    }

    /// @notice Challenges DKG result. If the submitted result is proved to be
    ///         invalid it reverts the DKG back to the result submission phase.
    /// @param dkgResult Result to challenge. Must match the submitted result
    ///        stored during `submitDkgResult`.
    function challengeDkgResult(EcdsaDkg.Result calldata dkgResult) external {
        (
            bytes32 maliciousDkgResultHash,
            uint32 maliciousDkgResultSubmitterId
        ) = dkg.challengeResult(dkgResult);

        address maliciousDkgResultSubmitterAddress = sortitionPool
            .getIDOperator(maliciousDkgResultSubmitterId);

        address[] memory operatorWrapper = new address[](1);
        operatorWrapper[0] = operatorToStakingProvider(
            maliciousDkgResultSubmitterAddress
        );

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
    function isDkgResultValid(EcdsaDkg.Result calldata result)
        external
        view
        returns (bool, string memory)
    {
        return dkg.isResultValid(result);
    }

    /// @notice Check current wallet creation state.
    function getWalletCreationState() external view returns (EcdsaDkg.State) {
        return dkg.currentState();
    }

    /// @notice Checks if awaiting seed timed out.
    /// @return True if awaiting seed timed out, false otherwise.
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

    function getWallet(bytes32 walletID)
        external
        view
        returns (Wallets.Wallet memory)
    {
        return wallets.registry[walletID];
    }

    /// @notice Gets public key of a wallet with a given wallet ID.
    ///         The public key is returned in an uncompressed format as a 64-byte
    ///         concatenation of X and Y coordinates.
    /// @param walletID ID of the wallet.
    /// @return Uncompressed public key of the wallet.
    function getWalletPublicKey(bytes32 walletID)
        external
        view
        returns (bytes memory)
    {
        return wallets.getWalletPublicKey(walletID);
    }

    /// @notice Checks if a wallet with the given ID is registered.
    /// @param walletID Wallet's ID.
    /// @return True if wallet is registered, false otherwise.
    function isWalletRegistered(bytes32 walletID) external view returns (bool) {
        return wallets.isWalletRegistered(walletID);
    }

    // TODO: Add function to close the Wallet so the members are notified that
    // they no longer need to track the wallet.

    /// @notice Retrieves dkg parameters that were set in DKG library.
    function dkgParameters()
        external
        view
        returns (EcdsaDkg.Parameters memory)
    {
        return dkg.parameters;
    }

    /// @notice The minimum authorization amount required so that operator can
    ///         participate in ECDSA Wallet operations.
    function minimumAuthorization() external view returns (uint96) {
        return authorization.parameters.minimumAuthorization;
    }

    /// @notice Delay in seconds that needs to pass between the time
    ///         authorization decrease is requested and the time that request
    ///         can get approved.
    function authorizationDecreaseDelay() external view returns (uint64) {
        return authorization.parameters.authorizationDecreaseDelay;
    }

    /// @notice Returns the current value of the staking provider's eligible
    ///         stake. Eligible stake is defined as the currently authorized
    ///         stake minus the pending authorization decrease. Eligible stake
    ///         is what is used for operator's weight in the sortition pool.
    function eligibleStake(address stakingProvider)
        external
        view
        returns (uint96)
    {
        return authorization.eligibleStake(staking, stakingProvider);
    }

    /// @notice Returns the remaining time in seconds that needs to pass before
    ///         the requested authorization decrease can be approved.
    ///         If the sortition pool state was not updated yet by the operator
    ///         after requesting the authorization decrease, returns
    ///         `type(uint64).max`.
    function remainingAuthorizationDecreaseDelay(address stakingProvider)
        external
        view
        returns (uint64)
    {
        return
            authorization.remainingAuthorizationDecreaseDelay(stakingProvider);
    }

    /// @notice Returns operator registered for the given staking provider.
    function stakingProviderToOperator(address stakingProvider)
        external
        view
        returns (address)
    {
        return authorization.stakingProviderToOperator[stakingProvider];
    }

    /// @notice Returns staking provider of the given operator.
    function operatorToStakingProvider(address operator)
        public
        view
        returns (address)
    {
        return authorization.operatorToStakingProvider[operator];
    }

    /// @notice Checks if the operator's authorized stake is in sync with
    ///         operator's weight in the sortition pool.
    ///         If the operator is not in the sortition pool and their
    ///         authorized stake is non-zero, function returns false.
    function isOperatorUpToDate(address operator) external view returns (bool) {
        return
            authorization.isOperatorUpToDate(staking, sortitionPool, operator);
    }

    /// @notice Selects a new group of operators based on the provided seed.
    ///         At least one operator has to be registered in the pool,
    ///         otherwise the function fails reverting the transaction.
    /// @param seed Number used to select operators to the group.
    /// @return IDs of selected group members.
    function selectGroup(bytes32 seed) external view returns (uint32[] memory) {
        // TODO: Read seed from EcdsaDkg
        return sortitionPool.selectGroup(EcdsaDkg.groupSize, seed);
    }
}
