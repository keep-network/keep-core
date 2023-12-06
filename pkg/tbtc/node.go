package tbtc

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	chainpkg "github.com/keep-network/keep-core/pkg/chain"

	"go.uber.org/zap"

	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/announcer"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing"
)

const (
	// signingAttemptsLimit determines the maximum number of signing attempts
	// that can be performed for the given message being subject of signing.
	//
	// The value of `5` should be enough to produce the signature even with
	// `2` malicious members in a signing group of `100` members. To produce
	// the signature, `51` members must be selected out of the honest `98`.
	// The probability of successful signing in that case is:
	// `P = (98 choose 51) / (100 choose 51) = ~0.24` which means we need
	// `5` attempts on the worst case.
	//
	// A greater limit does not necessarily make sense. Presence of more than
	// `2` malicious members in the signing group has a very small probability.
	// Moreover, the signature must be produced in the reasonable time.
	// That being said, the value `5` seems to be reasonable trade-off.
	signingAttemptsLimit = 5
)

// TODO: Unit tests for `node.go`.

// node represents the current state of an ECDSA node.
type node struct {
	groupParameters *GroupParameters

	chain          Chain
	btcChain       bitcoin.Chain
	netProvider    net.Provider
	walletRegistry *walletRegistry

	// walletDispatcher ensures only one action is executed by a wallet at
	// a time. All possible activities of a created wallet must be represented
	// by appropriate actions dispatched through this component.
	walletDispatcher *walletDispatcher

	// protocolLatch makes sure no expensive number generator operations are
	// running when signing or generating a wallet key are executed. The
	// protocolLatch is used by dkgExecutor and signingExecutor.
	protocolLatch *generator.ProtocolLatch

	// dkgExecutor encapsulates the logic of distributed key generation.
	//
	// dkgExecutor MUST NOT be used outside this struct.
	dkgExecutor *dkgExecutor

	signingExecutorsMutex sync.Mutex
	// signingExecutors is the cache holding signing executors for specific wallets.
	// The cache key is the uncompressed public key (with 04 prefix) of the wallet.
	// signingExecutor encapsulates the generic logic of signing messages.
	//
	// signingExecutors MUST NOT be used outside this struct. Please use
	// wallet actions and walletDispatcher to execute an action on an existing
	// wallet.
	signingExecutors map[string]*signingExecutor

	operatorIDsMutex sync.Mutex
	// operatorIDsCache is the cache mapping operator addresses to operator IDs.
	// The cache holds IDs of singing group operators of this node.
	operatorIDsCache map[chainpkg.Address]chainpkg.OperatorID

	coordinationExecutorsMutex sync.Mutex
	// coordinationExecutors is the cache holding coordination executors for
	// specific wallets. The cache key is the uncompressed public key
	// (with 04 prefix) of the wallet. The coordinationExecutor encapsulates the
	// logic of the wallet coordination procedure.
	//
	// coordinationExecutors MUST NOT be used outside this struct.
	coordinationExecutors map[string]*coordinationExecutor

	// proposalGenerator is the implementation of the coordination proposal
	// generator used by the node.
	proposalGenerator CoordinationProposalGenerator
}

func newNode(
	groupParameters *GroupParameters,
	chain Chain,
	btcChain bitcoin.Chain,
	netProvider net.Provider,
	keyStorePersistance persistence.ProtectedHandle,
	workPersistence persistence.BasicHandle,
	scheduler *generator.Scheduler,
	proposalGenerator CoordinationProposalGenerator,
	config Config,
) (*node, error) {
	walletRegistry := newWalletRegistry(keyStorePersistance)

	latch := generator.NewProtocolLatch()
	scheduler.RegisterProtocol(latch)

	node := &node{
		groupParameters:       groupParameters,
		chain:                 chain,
		btcChain:              btcChain,
		netProvider:           netProvider,
		walletRegistry:        walletRegistry,
		walletDispatcher:      newWalletDispatcher(),
		protocolLatch:         latch,
		signingExecutors:      make(map[string]*signingExecutor),
		operatorIDsCache:      make(map[chainpkg.Address]chainpkg.OperatorID),
		coordinationExecutors: make(map[string]*coordinationExecutor),
		proposalGenerator:     proposalGenerator,
	}

	// Only the operator address is known at this point and can be pre-fetched.
	// The operator ID must be determined later as the operator may not be in
	// the sortition pool yet.
	operatorAddress, err := node.operatorAddress()
	if err != nil {
		return nil, fmt.Errorf("cannot get node's operator adress: [%v]", err)
	}

	// TODO: This chicken and egg problem should be solved when
	// waitForBlockHeight becomes a part of BlockHeightWaiter interface.
	node.dkgExecutor = newDkgExecutor(
		node.groupParameters,
		node.operatorID,
		operatorAddress,
		chain,
		netProvider,
		walletRegistry,
		latch,
		config,
		workPersistence,
		scheduler,
		node.waitForBlockHeight,
	)

	return node, nil
}

// operatorAddress returns the node's operator address.
func (n *node) operatorAddress() (chainpkg.Address, error) {
	_, operatorPublicKey, err := n.chain.OperatorKeyPair()
	if err != nil {
		return "", fmt.Errorf("failed to get operator public key: [%v]", err)
	}

	operatorAddress, err := n.chain.Signing().PublicKeyToAddress(operatorPublicKey)
	if err != nil {
		return "", fmt.Errorf(
			"failed to convert operator public key to address: [%v]",
			err,
		)
	}

	return operatorAddress, nil
}

// operatorAddress returns the node's operator ID.
func (n *node) operatorID() (chainpkg.OperatorID, error) {
	operatorAddress, err := n.operatorAddress()
	if err != nil {
		return 0, fmt.Errorf("failed to get operator address: [%v]", err)
	}

	operatorID, err := n.chain.GetOperatorID(operatorAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to get operator ID: [%v]", err)
	}

	return operatorID, nil
}

// joinDKGIfEligible takes a seed value and undergoes the process of the
// distributed key generation if this node's operator proves to be eligible for
// the group generated by that seed. This is an interactive on-chain process,
// and joinDKGIfEligible can block for an extended period of time while it
// completes the on-chain operation. The execution can be delayed by an
// arbitrary number of blocks using the delayBlocks argument. This allows
// confirming the state on-chain - e.g. wait for the required number of
// confirming blocks - before executing the off-chain action.
func (n *node) joinDKGIfEligible(
	seed *big.Int,
	startBlock uint64,
	delayBlocks uint64,
) {
	n.dkgExecutor.executeDkgIfEligible(seed, startBlock, delayBlocks)
}

// validateDKG performs the submitted DKG result validation process.
// If the result is not valid, this function submits an on-chain result
// challenge. If the result is valid and the given node was involved in the DKG,
// this function schedules an on-chain approve that is submitted once the
// challenge period elapses.
func (n *node) validateDKG(
	seed *big.Int,
	submissionBlock uint64,
	result *DKGChainResult,
	resultHash [32]byte,
) {
	n.dkgExecutor.executeDkgValidation(seed, submissionBlock, result, resultHash)
}

// getSigningExecutor gets the signing executor responsible for executing
// signing related to a specific wallet whose part is controlled by this node.
// The second boolean return value indicates whether the node controls at least
// one signer for the given wallet.
func (n *node) getSigningExecutor(
	walletPublicKey *ecdsa.PublicKey,
) (*signingExecutor, bool, error) {
	n.signingExecutorsMutex.Lock()
	defer n.signingExecutorsMutex.Unlock()

	walletPublicKeyBytes, err := marshalPublicKey(walletPublicKey)
	if err != nil {
		return nil, false, fmt.Errorf("cannot marshal wallet public key: [%v]", err)
	}

	executorKey := hex.EncodeToString(walletPublicKeyBytes)

	if executor, exists := n.signingExecutors[executorKey]; exists {
		return executor, true, nil
	}

	executorLogger := logger.With(
		zap.String("wallet", fmt.Sprintf("0x%x", walletPublicKeyBytes)),
	)

	signers := n.walletRegistry.getSigners(walletPublicKey)
	if len(signers) == 0 {
		// This is not an error because the node simply does not control
		// the given wallet.
		return nil, false, nil
	}

	// All signers belong to one wallet. Take that wallet from the
	// first signer.
	wallet := signers[0].wallet

	channelName := fmt.Sprintf(
		"%s-%s",
		ProtocolName,
		hex.EncodeToString(walletPublicKeyBytes),
	)

	broadcastChannel, err := n.netProvider.BroadcastChannelFor(channelName)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get broadcast channel: [%v]", err)
	}

	signing.RegisterUnmarshallers(broadcastChannel)
	announcer.RegisterUnmarshaller(broadcastChannel)
	broadcastChannel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &signingDoneMessage{}
	})

	membershipValidator := group.NewMembershipValidator(
		executorLogger,
		wallet.signingGroupOperators,
		n.chain.Signing(),
	)

	err = broadcastChannel.SetFilter(membershipValidator.IsInGroup)
	if err != nil {
		return nil, false, fmt.Errorf(
			"could not set filter for channel [%v]: [%v]",
			broadcastChannel.Name(),
			err,
		)
	}

	executorLogger.Infof(
		"signing executor created; controlling [%v] signers",
		len(signers),
	)

	blockCounter, err := n.chain.BlockCounter()
	if err != nil {
		return nil, false, fmt.Errorf(
			"could not get block counter: [%v]",
			err,
		)
	}

	executor := newSigningExecutor(
		signers,
		broadcastChannel,
		membershipValidator,
		n.groupParameters,
		n.protocolLatch,
		blockCounter.CurrentBlock,
		n.waitForBlockHeight,
		signingAttemptsLimit,
	)

	n.signingExecutors[executorKey] = executor

	return executor, true, nil
}

// getSigningGroupOperatorID gets the operator ID of the signing group operator
// based on the provided operator address. The operator ID is cached for future
// efficient retrievals.
func (n *node) getSigningGroupOperatorID(
	operatorAddress chainpkg.Address,
) (chainpkg.OperatorID, error) {
	n.operatorIDsMutex.Lock()
	defer n.operatorIDsMutex.Unlock()

	if operatorID, exists := n.operatorIDsCache[operatorAddress]; exists {
		return operatorID, nil
	}

	operatorID, err := n.chain.GetOperatorID(operatorAddress)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get operator ID for operator with address [%s]: [%v]",
			operatorAddress,
			err,
		)
	}

	n.operatorIDsCache[operatorAddress] = operatorID
	return operatorID, nil
}

// getCoordinationExecutor gets the coordination executor responsible for
// executing coordination related to a specific wallet whose part is controlled
// by this node. The second boolean return value indicates whether the node
// controls at least one signer for the given wallet.
func (n *node) getCoordinationExecutor(
	walletPublicKey *ecdsa.PublicKey,
) (*coordinationExecutor, bool, error) {
	n.coordinationExecutorsMutex.Lock()
	defer n.coordinationExecutorsMutex.Unlock()

	walletPublicKeyBytes, err := marshalPublicKey(walletPublicKey)
	if err != nil {
		return nil, false, fmt.Errorf("cannot marshal wallet public key: [%v]", err)
	}

	executorKey := hex.EncodeToString(walletPublicKeyBytes)

	if executor, exists := n.coordinationExecutors[executorKey]; exists {
		return executor, true, nil
	}

	executorLogger := logger.With(
		zap.String("wallet", fmt.Sprintf("0x%x", walletPublicKeyBytes)),
	)

	signers := n.walletRegistry.getSigners(walletPublicKey)
	if len(signers) == 0 {
		// This is not an error because the node simply does not control
		// the given wallet.
		return nil, false, nil
	}

	// All signers belong to one wallet. Take that wallet from the
	// first signer.
	wallet := signers[0].wallet

	channelName := fmt.Sprintf(
		"%s-%s-coordination",
		ProtocolName,
		hex.EncodeToString(walletPublicKeyBytes),
	)

	broadcastChannel, err := n.netProvider.BroadcastChannelFor(channelName)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get broadcast channel: [%v]", err)
	}

	broadcastChannel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &coordinationMessage{}
	})

	membershipValidator := group.NewMembershipValidator(
		executorLogger,
		wallet.signingGroupOperators,
		n.chain.Signing(),
	)

	err = broadcastChannel.SetFilter(membershipValidator.IsInGroup)
	if err != nil {
		return nil, false, fmt.Errorf(
			"could not set filter for channel [%v]: [%v]",
			broadcastChannel.Name(),
			err,
		)
	}

	// The coordination executor does not need access to signers' key material.
	// It is enough to pass only their member indexes.
	membersIndexes := make([]group.MemberIndex, len(signers))
	for i, s := range signers {
		membersIndexes[i] = s.signingGroupMemberIndex
	}

	operatorAddress, err := n.operatorAddress()
	if err != nil {
		return nil, false, fmt.Errorf("failed to get operator address: [%v]", err)
	}

	executor := newCoordinationExecutor(
		n.chain,
		wallet,
		membersIndexes,
		operatorAddress,
		n.proposalGenerator,
		broadcastChannel,
		membershipValidator,
		n.protocolLatch,
		n.waitForBlockHeight,
	)

	n.coordinationExecutors[executorKey] = executor

	executorLogger.Infof(
		"coordination executor created; controlling [%v] signers",
		len(signers),
	)

	return executor, true, nil
}

// handleHeartbeatProposal handles an incoming heartbeat proposal by
// orchestrating and dispatching an appropriate wallet action.
func (n *node) handleHeartbeatProposal(
	wallet wallet,
	proposal *HeartbeatProposal,
	startBlock uint64,
	expiryBlock uint64,
) {
	walletPublicKeyHash := bitcoin.PublicKeyHash(wallet.publicKey)

	signingExecutor, ok, err := n.getSigningExecutor(wallet.publicKey)
	if err != nil {
		logger.Errorf("cannot get signing executor: [%v]", err)
		return
	}
	// This check is actually redundant. We know the node controls some
	// wallet signers as we just got the wallet from the registry using their
	// public key hash. However, we are doing it just in case. The API
	// contract of getSigningExecutor may change one day.
	if !ok {
		logger.Infof(
			"node does not control signers of wallet PKH [0x%x]; "+
				"ignoring the received heartbeat request",
			walletPublicKeyHash,
		)
		return
	}

	walletPublicKeyBytes, err := marshalPublicKey(wallet.publicKey)
	if err != nil {
		logger.Errorf("cannot marshal wallet public key: [%v]", err)
		return
	}

	logger.Infof(
		"node controls signers of wallet PKH [0x%x]; "+
			"plain-text uncompressed public key of that wallet is [0x%x]; "+
			"starting orchestration of the heartbeat action",
		walletPublicKeyHash,
		walletPublicKeyBytes,
	)

	walletActionLogger := logger.With(
		zap.String("wallet", fmt.Sprintf("0x%x", walletPublicKeyBytes)),
		zap.String("action", ActionHeartbeat.String()),
		zap.Uint64("startBlock", startBlock),
		zap.Uint64("expiryBlock", expiryBlock),
	)
	walletActionLogger.Infof("dispatching wallet action")

	action := newHeartbeatAction(
		walletActionLogger,
		wallet,
		signingExecutor,
		proposal.Message[:],
		startBlock,
		expiryBlock,
		n.waitForBlockHeight,
	)

	err = n.walletDispatcher.dispatch(action)
	if err != nil {
		walletActionLogger.Errorf("cannot dispatch wallet action: [%v]", err)
		return
	}

	walletActionLogger.Infof("wallet action dispatched successfully")
}

// handleDepositSweepProposal handles an incoming deposit sweep proposal by
// orchestrating and dispatching an appropriate wallet action.
func (n *node) handleDepositSweepProposal(
	wallet wallet,
	proposal *DepositSweepProposal,
	startBlock uint64,
	expiryBlock uint64,
) {
	signingExecutor, ok, err := n.getSigningExecutor(wallet.publicKey)
	if err != nil {
		logger.Errorf("cannot get signing executor: [%v]", err)
		return
	}
	// This check is actually redundant. We know the node controls some
	// wallet signers as we just got the wallet from the registry using their
	// public key hash. However, we are doing it just in case. The API
	// contract of getSigningExecutor may change one day.
	if !ok {
		logger.Infof(
			"node does not control signers of wallet PKH [0x%x]; "+
				"ignoring the received deposit sweep proposal",
			proposal.WalletPublicKeyHash,
		)
		return
	}

	walletPublicKeyBytes, err := marshalPublicKey(wallet.publicKey)
	if err != nil {
		logger.Errorf("cannot marshal wallet public key: [%v]", err)
		return
	}

	logger.Infof(
		"node controls signers of wallet PKH [0x%x]; "+
			"plain-text uncompressed public key of that wallet is [0x%x]; "+
			"starting orchestration of the deposit sweep action",
		proposal.WalletPublicKeyHash,
		walletPublicKeyBytes,
	)

	walletActionLogger := logger.With(
		zap.String("wallet", fmt.Sprintf("0x%x", walletPublicKeyBytes)),
		zap.String("action", ActionDepositSweep.String()),
		zap.Uint64("startBlock", startBlock),
		zap.Uint64("expiryBlock", expiryBlock),
	)
	walletActionLogger.Infof("dispatching wallet action")

	action := newDepositSweepAction(
		walletActionLogger,
		n.chain,
		n.btcChain,
		wallet,
		signingExecutor,
		proposal,
		startBlock,
		expiryBlock,
		n.waitForBlockHeight,
	)

	err = n.walletDispatcher.dispatch(action)
	if err != nil {
		walletActionLogger.Errorf("cannot dispatch wallet action: [%v]", err)
		return
	}

	walletActionLogger.Infof("wallet action dispatched successfully")
}

// handleRedemptionProposal handles an incoming redemption proposal by
// orchestrating and dispatching an appropriate wallet action.
func (n *node) handleRedemptionProposal(
	wallet wallet,
	proposal *RedemptionProposal,
	startBlock uint64,
	expiryBlock uint64,
) {
	signingExecutor, ok, err := n.getSigningExecutor(wallet.publicKey)
	if err != nil {
		logger.Errorf("cannot get signing executor: [%v]", err)
		return
	}
	// This check is actually redundant. We know the node controls some
	// wallet signers as we just got the wallet from the registry using their
	// public key hash. However, we are doing it just in case. The API
	// contract of getSigningExecutor may change one day.
	if !ok {
		logger.Infof(
			"node does not control signers of wallet PKH [0x%x]; "+
				"ignoring the received redemption proposal",
			proposal.WalletPublicKeyHash,
		)
		return
	}

	walletPublicKeyBytes, err := marshalPublicKey(wallet.publicKey)
	if err != nil {
		logger.Errorf("cannot marshal wallet public key: [%v]", err)
		return
	}

	logger.Infof(
		"node controls signers of wallet PKH [0x%x]; "+
			"plain-text uncompressed public key of that wallet is [0x%x]; "+
			"starting orchestration of the redemption action",
		proposal.WalletPublicKeyHash,
		walletPublicKeyBytes,
	)

	walletActionLogger := logger.With(
		zap.String("wallet", fmt.Sprintf("0x%x", walletPublicKeyBytes)),
		zap.String("action", ActionRedemption.String()),
		zap.Uint64("startBlock", startBlock),
		zap.Uint64("expiryBlock", expiryBlock),
	)
	walletActionLogger.Infof("dispatching wallet action")

	action := newRedemptionAction(
		walletActionLogger,
		n.chain,
		n.btcChain,
		wallet,
		signingExecutor,
		proposal,
		startBlock,
		expiryBlock,
		n.waitForBlockHeight,
	)

	err = n.walletDispatcher.dispatch(action)
	if err != nil {
		walletActionLogger.Errorf("cannot dispatch wallet action: [%v]", err)
		return
	}

	walletActionLogger.Infof("wallet action dispatched successfully")
}

func (n *node) HandleMovingFundsProposal(sourceWalletPublicKeyHash [20]byte) {
	go func() {
		logger.Info(
			"moving funds proposal initiated for wallet with PKH [0x%x]",
			sourceWalletPublicKeyHash,
		)

		// Make sure the wallet meets the criteria for moving funds proposal.
		sourceWallet, found := n.walletRegistry.getWalletByPublicKeyHash(
			sourceWalletPublicKeyHash,
		)
		if !found {
			logger.Errorf(
				"skipping moving funds proposal for wallet with PKH "+
					"[0x%x] as the node does not control it",
				sourceWalletPublicKeyHash,
			)
			return
		}

		sourceWalletChainData, err := n.chain.GetWallet(sourceWalletPublicKeyHash)
		if err != nil {
			logger.Errorf(
				"failed to get wallet data for source wallet with PKH "+
					"[0x%x]: [%v]",
				sourceWalletPublicKeyHash,
				err,
			)
			return
		}

		walletMainUtxo, err := DetermineWalletMainUtxo(
			sourceWalletPublicKeyHash,
			n.chain,
			n.btcChain,
		)
		if err != nil {
			logger.Errorf(
				"skipping moving funds proposal for wallet with PKH "+
					"[0x%x] due to error determining wallet main UTXO: [%v]",
				sourceWalletPublicKeyHash,
				err,
			)
			return
		}

		walletBalance := int64(0)
		if walletMainUtxo != nil {
			walletBalance = walletMainUtxo.Value
		}

		if sourceWalletChainData.State != StateMovingFunds {
			logger.Infof(
				"ignoring moving funds proposal for wallet with PKH [0x%x]; "+
					"wallet not in MovingFunds state",
				sourceWalletPublicKeyHash,
			)
			return
		}

		if sourceWalletChainData.PendingRedemptionsValue > 0 {
			logger.Infof(
				"ignoring moving funds proposal for wallet with PKH [0x%x]; "+
					"wallet has pending redemptions",
				sourceWalletPublicKeyHash,
			)
			return
		}

		if sourceWalletChainData.PendingMovedFundsSweepRequestsCount > 0 {
			logger.Infof(
				"ignoring moving funds proposal for wallet with PKH [0x%x]; "+
					"wallet has pending moved funds sweep requests",
				sourceWalletPublicKeyHash,
			)
			return
		}

		if sourceWalletChainData.MovingFundsTargetWalletsCommitmentHash != [32]byte{} {
			logger.Infof(
				"ignoring moving funds proposal for wallet with PKH [0x%x]; "+
					"wallet has already submitted commitment",
				sourceWalletPublicKeyHash,
			)
			return
		}

		if walletBalance <= 0 {
			// The wallet's balance cannot be `0`. Since we are dealing with
			// a signed integer we also check it's not negative just in case.
			logger.Infof(
				"ignoring moving funds proposal for wallet with PKH [0x%x]; "+
					"wallet does not have a positive balance",
				sourceWalletPublicKeyHash,
			)
			return
		}

		logger.Infof(
			"proceeding with moving funds commitment for wallet with "+
				"PKH [0x%x]",
			sourceWalletPublicKeyHash,
		)

		// Prepare the list of target wallets.
		liveWalletsCount, err := n.chain.GetLiveWalletsCount()
		if err != nil {
			logger.Errorf("failed to get live wallets count: [%v]", err)
			return
		}

		if liveWalletsCount == 0 {
			logger.Infof(
				"skipping moving funds proposal for wallet with PKH [0x%x] due"+
					"to lack of live wallets",
				sourceWalletPublicKeyHash,
			)
			return
		}

		_, _, _, _, _, walletMaxBtcTransfer, _, err := n.chain.GetWalletParameters()
		if err != nil {
			logger.Errorf("failed to get wallet parameters: [%v]", err)
			return
		}

		if walletMaxBtcTransfer == 0 {
			logger.Errorf("wallet max BTC transfer parameter must be positive")
			return
		}

		ceilingDivide := func(x, y uint64) uint64 {
			// The divisor must be positive, but we do not need to check it as
			// this function will be executed with wallet max BTC transfer as
			// the divisor and we already ensured it is positive.
			return (x + y - 1) / y
		}
		min := func(x, y uint64) uint64 {
			if x < y {
				return x
			}
			return y
		}

		targetWalletsCount := min(
			uint64(liveWalletsCount),
			ceilingDivide(uint64(walletBalance), walletMaxBtcTransfer),
		)

		// Prepare a list of target wallets using the new wallets registration
		// events. Retrieve only the necessary number of live wallets.
		// The iteration is started from the end of the
		events, err := n.chain.PastNewWalletRegisteredEvents(nil)
		if err != nil {
			logger.Errorf(
				"failed to get past new wallet registered events: [%v]",
				err,
			)
			return
		}

		targetWallets := make([][20]byte, 0)

		for i := len(events) - 1; i >= 0; i-- {
			walletPubKeyHash := events[i].WalletPublicKeyHash
			if walletPubKeyHash == sourceWalletPublicKeyHash {
				// Just in case make sure not to include the source wallet
				// itself.
				continue
			}
			wallet, err := n.chain.GetWallet(walletPubKeyHash)
			if err != nil {
				logger.Errorf(
					"failed to get wallet data for wallet with PKH [0x%x]: [%v]",
					walletPubKeyHash,
					err,
				)
				continue
			}
			if wallet.State == StateLive {
				targetWallets = append(targetWallets, walletPubKeyHash)
			}
			if len(targetWallets) == int(targetWalletsCount) {
				// Stop the iteration if enough live wallets have been gathered.
				break
			}
		}

		if len(targetWallets) != int(targetWalletsCount) {
			logger.Errorf(
				"failed to get enough target wallets: required [%v]; "+
					"gathered [%v]",
				targetWalletsCount,
				len(targetWallets),
			)
			return
		}

		// Sort the target wallets according to their numerical representation
		// as the on-chain contract expects.
		sort.Slice(targetWallets, func(i, j int) bool {
			bigIntI := new(big.Int).SetBytes(targetWallets[i][:])
			bigIntJ := new(big.Int).SetBytes(targetWallets[j][:])
			return bigIntI.Cmp(bigIntJ) < 0
		})

		logger.Infof("gathered [%v] target wallets", len(targetWallets))

		walletMemberIDs := make([]uint32, 0)
		for _, operatorAddress := range sourceWallet.signingGroupOperators {
			operatorID, err := n.getSigningGroupOperatorID(operatorAddress)
			if err != nil {
				logger.Errorf(
					"failed to get operator ID belonging to wallet with "+
						"PKH [0x%x]: [%v]",
					sourceWalletPublicKeyHash,
					err,
				)
				return
			}
			walletMemberIDs = append(walletMemberIDs, operatorID)
		}

		latestBlockHeight, err := n.btcChain.GetLatestBlockHeight()
		if err != nil {
			logger.Errorf(
				"failed to get latest Bitcoin block height: [%v]",
				err,
			)
			return
		}

		// Use the latest Bitcoin block height to determine the wallet member
		// index. Increase the result of the modulo operation by one since the
		// wallet member index must be within range [1, len(walletMemberIDs)].
		walletMemberIndex := (int(latestBlockHeight) % len(walletMemberIDs)) + 1

		err = n.chain.SubmitMovingFundsCommitment(
			sourceWalletPublicKeyHash,
			*walletMainUtxo,
			walletMemberIDs,
			uint32(walletMemberIndex),
			targetWallets,
		)
		if err != nil {
			logger.Errorf(
				"failed to submit moving funds commitment for wallet wit PKH "+
					"[0x%x]: [%v]",
				sourceWalletPublicKeyHash,
				err,
			)
			return
		}

		logger.Infof(
			"finished moving funds commitment for wallet with PKH [0x%x]",
			sourceWalletPublicKeyHash,
		)

		// TODO: Add construction of the move funds Bitcoin transaction.
		//       Before proceeding with the Bitcoin transaction, check if the
		//       commitment was successfully submitted.
	}()
}

// coordinationLayerSettings represents settings for the coordination layer.
type coordinationLayerSettings struct {
	// executeCoordinationProcedureFn is a function executing the coordination
	// procedure for the given wallet and coordination window.
	executeCoordinationProcedureFn func(
		node *node,
		window *coordinationWindow,
		walletPublicKey *ecdsa.PublicKey,
	) (*coordinationResult, bool)

	// processCoordinationResultFn is a function processing the given
	// coordination result.
	processCoordinationResultFn func(
		node *node,
		result *coordinationResult,
	)
}

// runCoordinationLayer starts the coordination layer of the node. It is
// responsible for detecting new coordination windows, running coordination
// procedures for all wallets controlled by the node, and processing
// coordination results.
func (n *node) runCoordinationLayer(
	ctx context.Context,
	settings ...*coordinationLayerSettings,
) error {
	// Resolve settings for the coordination layer.
	var cls *coordinationLayerSettings
	switch len(settings) {
	case 1:
		cls = settings[0]
	default:
		cls = &coordinationLayerSettings{
			executeCoordinationProcedureFn: executeCoordinationProcedure,
			processCoordinationResultFn:    processCoordinationResult,
		}
	}

	blockCounter, err := n.chain.BlockCounter()
	if err != nil {
		return fmt.Errorf("cannot get block counter: [%w]", err)
	}

	coordinationResultChan := make(chan *coordinationResult)

	// Prepare a callback function that will be called every time a new
	// coordination window is detected.
	onWindowFn := func(window *coordinationWindow) {
		// Fetch all wallets controlled by the node. It is important to
		// get the wallets every time the window is triggered as the
		// node may have started controlling a new wallet in the meantime.
		walletsPublicKeys := n.walletRegistry.getWalletsPublicKeys()

		for _, currentWalletPublicKey := range walletsPublicKeys {
			// Run an independent coordination procedure for the given wallet
			// in a separate goroutine. The coordination result will be sent
			// to the coordination result channel.
			go func(walletPublicKey *ecdsa.PublicKey) {
				result, ok := cls.executeCoordinationProcedureFn(
					n,
					window,
					walletPublicKey,
				)
				if ok {
					coordinationResultChan <- result
				}
			}(currentWalletPublicKey)
		}
	}

	// Start the coordination windows watcher.
	go watchCoordinationWindows(
		ctx,
		blockCounter.WatchBlocks,
		onWindowFn,
	)

	// Start the coordination result processor.
	go func() {
		for {
			select {
			case result := <-coordinationResultChan:
				go cls.processCoordinationResultFn(n, result)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// executeCoordinationProcedure executes the coordination procedure for the
// given wallet and coordination window.
func executeCoordinationProcedure(
	node *node,
	window *coordinationWindow,
	walletPublicKey *ecdsa.PublicKey,
) (*coordinationResult, bool) {
	walletPublicKeyBytes, err := marshalPublicKey(walletPublicKey)
	if err != nil {
		logger.Errorf("cannot marshal wallet public key: [%v]", err)
		return nil, false
	}

	procedureLogger := logger.With(
		zap.Uint64("coordinationBlock", window.coordinationBlock),
		zap.String("wallet", fmt.Sprintf("0x%x", walletPublicKeyBytes)),
	)

	procedureLogger.Infof("starting coordination procedure")

	executor, ok, err := node.getCoordinationExecutor(walletPublicKey)
	if err != nil {
		procedureLogger.Errorf("cannot get coordination executor: [%v]", err)
		return nil, false
	}
	// This check is actually redundant. We know the node controls some
	// wallet signers as we just got the wallet from the registry.
	// However, we are doing it just in case. The API contract of
	// getWalletsPublicKeys and/or getCoordinationExecutor may change one day.
	if !ok {
		procedureLogger.Infof("node does not control signers of this wallet")
		return nil, false
	}

	result, err := executor.coordinate(window)
	if err != nil {
		procedureLogger.Errorf("coordination procedure failed: [%v]", err)
		return nil, false
	}

	procedureLogger.Infof(
		"coordination procedure finished successfully with result [%s]",
		result,
	)

	return result, true
}

// processCoordinationResult processes the given coordination result.
func processCoordinationResult(node *node, result *coordinationResult) {
	logger.Infof("processing coordination result [%s]", result)

	// TODO: In the future, create coordination faults cache and
	//       record faults from the processed results there.

	startBlock := result.window.endBlock()
	expiryBlock := startBlock + result.proposal.ValidityBlocks()

	switch result.proposal.ActionType() {
	case ActionHeartbeat:
		if proposal, ok := result.proposal.(*HeartbeatProposal); ok {
			node.handleHeartbeatProposal(
				result.wallet,
				proposal,
				startBlock,
				expiryBlock,
			)
		}
	case ActionDepositSweep:
		if proposal, ok := result.proposal.(*DepositSweepProposal); ok {
			node.handleDepositSweepProposal(
				result.wallet,
				proposal,
				startBlock,
				expiryBlock,
			)
		}
	case ActionRedemption:
		if proposal, ok := result.proposal.(*RedemptionProposal); ok {
			node.handleRedemptionProposal(
				result.wallet,
				proposal,
				startBlock,
				expiryBlock,
			)
		}
	// TODO: Uncomment when moving funds support is implemented.
	// case ActionMovingFunds:
	//	 if proposal, ok := result.proposal.(*MovingFundsProposal); ok {
	//	 	 node.handleMovingFundsProposal(
	//	 	 	 result.wallet,
	//	 	 	 proposal,
	//	 	 	 startBlock,
	//	 	 	 expiryBlock,
	//	 	 )
	//	 }
	// case ActionMovedFundsSweep:
	//	 if proposal, ok := result.proposal.(*MovedFundsSweepProposal); ok {
	//	 	 node.handleMovedFundsSweepProposal(
	//	 	 	 result.wallet,
	//	 	 	 proposal,
	//	 	 	 startBlock,
	//	 	 	 expiryBlock,
	//	 	 )
	//	 }
	default:
		logger.Errorf("no handler for coordination result [%s]", result)
	}
}

// waitForBlockFn represents a function blocking the execution until the given
// block height.
type waitForBlockFn func(context.Context, uint64) error

// getCurrentBlockFn represents a function returning the current block height.
type getCurrentBlockFn func() (uint64, error)

// TODO: this should become a part of BlockHeightWaiter interface.
func (n *node) waitForBlockHeight(ctx context.Context, blockHeight uint64) error {
	blockCounter, err := n.chain.BlockCounter()
	if err != nil {
		return err
	}

	wait, err := blockCounter.BlockHeightWaiter(blockHeight)
	if err != nil {
		return err
	}

	select {
	case <-wait:
	case <-ctx.Done():
	}

	return nil
}

// withCancelOnBlock returns a copy of the given ctx that is automatically
// cancelled on the given block or when the parent ctx is done. Note that the
// context can be cancelled earlier if the waitForBlockFn returns an error.
func withCancelOnBlock(
	ctx context.Context,
	block uint64,
	waitForBlockFn waitForBlockFn,
) (context.Context, context.CancelFunc) {
	blockCtx, cancelBlockCtx := context.WithCancel(ctx)

	go func() {
		defer cancelBlockCtx()

		err := waitForBlockFn(ctx, block)
		if err != nil {
			logger.Errorf(
				"failed to wait for block [%v]; "+
					"context cancelled earlier than expected",
				err,
			)
		}
	}()

	return blockCtx, cancelBlockCtx
}
