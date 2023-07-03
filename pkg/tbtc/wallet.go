package tbtc

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"go.uber.org/zap"
)

// WalletActionType represents actions types that can be performed by a wallet.
type WalletActionType uint8

const (
	ActionNoop WalletActionType = iota
	ActionHeartbeat
	ActionDepositSweep
	ActionRedemption
	ActionMovingFunds
	ActionMovedFundsSweep
)

func (wat WalletActionType) String() string {
	switch wat {
	case ActionNoop:
		return "Noop"
	case ActionHeartbeat:
		return "Heartbeat"
	case ActionDepositSweep:
		return "DepositSweep"
	case ActionRedemption:
		return "Redemption"
	case ActionMovingFunds:
		return "MovingFunds"
	case ActionMovedFundsSweep:
		return "MovedFundsSweep"
	default:
		panic("unknown wallet action type")
	}
}

// walletAction represents an action that can be performed by the wallet.
type walletAction interface {
	// execute carries out the walletAction until completion.
	execute() error

	// wallet returns the wallet the walletAction is bound to.
	wallet() wallet

	// actionType returns the specific type of the walletAction.
	actionType() WalletActionType
}

// WalletState represents the state of a wallet.
type WalletState uint8

const (
	StateUnknown WalletState = iota
	StateLive
	StateMovingFunds
	StateClosing
	StateClosed
	StateTerminated
)

func (ws WalletState) String() string {
	switch ws {
	case StateUnknown:
		return "Unknown"
	case StateLive:
		return "Live"
	case StateMovingFunds:
		return "MovingFunds"
	case StateClosing:
		return "Closing"
	case StateClosed:
		return "Closed"
	case StateTerminated:
		return "Terminated"
	default:
		panic("unknown wallet state")
	}
}

// errWalletBusy is an error returned when the waller cannot execute the
// requested walletAction due to an ongoing work.
var errWalletBusy = fmt.Errorf("wallet is busy")

// walletDispatcher is a component responsible for dispatching wallet actions
// to specific wallets.
type walletDispatcher struct {
	actionsMutex sync.Mutex
	// actions is the mapping holding the currently executed action of the
	// given wallet. The mapping key is the uncompressed public key
	// (with 04 prefix) of the wallet.
	actions map[string]WalletActionType
}

func newWalletDispatcher() *walletDispatcher {
	return &walletDispatcher{
		actions: make(map[string]WalletActionType),
	}
}

// dispatch sends the given walletAction for execution. If the wallet is
// already busy, an errWalletBusy error is returned and the action is ignored.
func (wd *walletDispatcher) dispatch(action walletAction) error {
	wd.actionsMutex.Lock()
	defer wd.actionsMutex.Unlock()

	walletPublicKeyBytes, err := marshalPublicKey(action.wallet().publicKey)
	if err != nil {
		return fmt.Errorf("cannot marshal wallet public key: [%v]", err)
	}

	walletActionLogger := logger.With(
		zap.String("wallet", fmt.Sprintf("0x%x", walletPublicKeyBytes)),
		zap.String("action", action.actionType().String()),
	)

	key := hex.EncodeToString(walletPublicKeyBytes)

	if _, ok := wd.actions[key]; ok {
		return errWalletBusy
	}

	wd.actions[key] = action.actionType()

	go func() {
		defer func() {
			wd.actionsMutex.Lock()
			delete(wd.actions, key)
			wd.actionsMutex.Unlock()
		}()

		walletActionLogger.Infof("starting action execution")

		err := action.execute()
		if err != nil {
			walletActionLogger.Errorf(
				"action execution terminated with error: [%v]",
				err,
			)
			return
		}

		walletActionLogger.Infof("action execution terminated with success")
	}()

	return nil
}

// walletSigningExecutor is an interface meant to decouple the specific
// implementation of the signing executor from the wallet transaction executor.
type walletSigningExecutor interface {
	signBatch(
		ctx context.Context,
		messages []*big.Int,
		startBlock uint64,
	) ([]*tecdsa.Signature, error)
}

// walletTransactionExecutor is a component allowing to sign and broadcast
// wallet Bitcoin transactions.
type walletTransactionExecutor struct {
	btcChain bitcoin.Chain

	executingWallet wallet
	signingExecutor walletSigningExecutor
}

func newWalletTransactionExecutor(
	btcChain bitcoin.Chain,
	executingWallet wallet,
	signingExecutor walletSigningExecutor,
) *walletTransactionExecutor {
	return &walletTransactionExecutor{
		btcChain:        btcChain,
		executingWallet: executingWallet,
		signingExecutor: signingExecutor,
	}
}

// signTransaction performs signing of an unsigned Bitcoin transaction
// and returns a signed transaction ready to be broadcasted over the
// Bitcoin network.
func (wte *walletTransactionExecutor) signTransaction(
	signTxLogger log.StandardLogger,
	unsignedTx *bitcoin.TransactionBuilder,
	signingStartBlock uint64,
	signingTimesOutAt time.Time,
) (*bitcoin.Transaction, error) {
	signTxLogger.Infof("computing transaction's sig hashes")

	sigHashes, err := unsignedTx.ComputeSignatureHashes()
	if err != nil {
		return nil, fmt.Errorf(
			"error while computing transaction's sig hashes: [%v]",
			err,
		)
	}

	signTxLogger.Infof("signing transaction's sig hashes")

	signingCtx, cancelSigningCtx := context.WithTimeout(
		context.Background(),
		time.Until(signingTimesOutAt),
	)
	defer cancelSigningCtx()

	signatures, err := wte.signingExecutor.signBatch(
		signingCtx,
		sigHashes,
		signingStartBlock,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error while signing transaction's sig hashes: [%v]",
			err,
		)
	}

	signTxLogger.Infof("applying transaction's signatures")

	containers := make([]*bitcoin.SignatureContainer, len(signatures))
	for i, signature := range signatures {
		containers[i] = &bitcoin.SignatureContainer{
			R:         signature.R,
			S:         signature.S,
			PublicKey: wte.executingWallet.publicKey,
		}
	}

	tx, err := unsignedTx.AddSignatures(containers)
	if err != nil {
		return nil, fmt.Errorf(
			"error while applying transaction's signatures: [%v]",
			err,
		)
	}

	signTxLogger.Infof("transaction created successfully")

	return tx, nil
}

// broadcastTransaction broadcasts a signed Bitcoin transaction until
// the transaction lands in the Bitcoin mempool or the provided timeout
// is hit, whichever comes first.
func (wte *walletTransactionExecutor) broadcastTransaction(
	broadcastTxLogger log.StandardLogger,
	tx *bitcoin.Transaction,
	timeout time.Duration,
	checkDelay time.Duration,
) error {
	txHash := tx.Hash()

	broadcastCtx, cancelBroadcastCtx := context.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancelBroadcastCtx()

	broadcastAttempt := 0

	for {
		select {
		case <-broadcastCtx.Done():
			return fmt.Errorf("broadcast timeout exceeded")
		default:
			broadcastAttempt++

			broadcastTxLogger.Infof(
				"broadcasting transaction on the Bitcoin chain - attempt [%v]",
				broadcastAttempt,
			)

			err := wte.btcChain.BroadcastTransaction(tx)
			if err != nil {
				broadcastTxLogger.Warnf(
					"broadcasting failed: [%v]; transaction could be "+
						"broadcasted by another wallet operators though",
					err,
				)
			} else {
				broadcastTxLogger.Infof("broadcasting completed")
			}

			broadcastTxLogger.Infof(
				"waiting [%v] before checking whether the "+
					"transaction is known on Bitcoin chain",
				checkDelay,
			)

			select {
			case <-time.After(checkDelay):
			case <-broadcastCtx.Done():
				return fmt.Errorf("broadcast timeout exceeded")
			}

			broadcastTxLogger.Infof(
				"checking whether the transaction is known on Bitcoin chain",
			)

			_, err = wte.btcChain.GetTransactionConfirmations(txHash)
			if err != nil {
				broadcastTxLogger.Warnf(
					"cannot say whether the transaction is known "+
						"on Bitcoin chain; check returned an error: [%v]",
					err,
				)
				continue
			}

			broadcastTxLogger.Infof("transaction is known on Bitcoin chain")
			return nil
		}
	}
}

// wallet represents a tBTC wallet. A wallet is one of the basic building
// blocks of the system that takes BTC under custody during the deposit
// process and gives that BTC back during redemptions.
type wallet struct {
	// publicKey is the unique ECDSA public key that identifies the
	// given wallet. This public key is also used to derive contract-specific
	// wallet identifiers (e.g. the Bridge contract identifies the wallet using
	// the SHA-256+RIPEMD-160 hash computed over the compressed ECDSA public key)
	publicKey *ecdsa.PublicKey
	// signingGroupOperators is the list holding operators' addresses that
	// form the whole wallet's signing group. This list may differ from the
	// original list outputted by the sortition protocol as it contains only
	// those signing group members who behaved properly during the DKG
	// protocol so all misbehaved members are not included here.
	// This list's size is always in the range [GroupQuorum, GroupSize].
	//
	// Each item in this list represents the given signing group member (seat)
	// and has a group.MemberIndex that is just the element's list index
	// incremented by one (e.g. element with index 0 has the group.MemberIndex
	// equal to 1 and so on).
	signingGroupOperators []chain.Address
}

// groupSize returns the actual size of the wallet's signing group. This
// value may be different from the GroupParameters.GroupSize parameter as some
// candidates may be excluded during distributed key generation.
func (w *wallet) groupSize() int {
	return len(w.signingGroupOperators)
}

// groupDishonestThreshold returns the dishonest threshold for the wallet's
// signing group. The returned value is computed using the wallet's actual
// signing group size for the given honest threshold provided as argument.
func (w *wallet) groupDishonestThreshold(honestThreshold int) int {
	return w.groupSize() - honestThreshold
}

func (w *wallet) String() string {
	publicKey := elliptic.Marshal(
		w.publicKey.Curve,
		w.publicKey.X,
		w.publicKey.Y,
	)

	return fmt.Sprintf(
		"wallet [0x%x] with a signing group of [%v]",
		publicKey,
		len(w.signingGroupOperators),
	)
}

// DetermineWalletMainUtxo determines the plain-text wallet main UTXO
// currently registered in the Bridge on-chain contract. The returned
// main UTXO can be nil if the wallet does not have a main UTXO registered
// in the Bridge at the moment.
//
// WARNING: THIS FUNCTION CANNOT DETERMINE THE MAIN UTXO IF IT COMES FROM A
// BITCOIN TRANSACTION THAT IS NOT ONE OF THE LATEST FIVE TRANSACTIONS
// TARGETING THE GIVEN WALLET PUBLIC KEY HASH. HOWEVER, SUCH A CASE IS
// VERY UNLIKELY.
func DetermineWalletMainUtxo(
	walletPublicKeyHash [20]byte,
	bridgeChain BridgeChain,
	btcChain bitcoin.Chain,
) (*bitcoin.UnspentTransactionOutput, error) {
	walletChainData, err := bridgeChain.GetWallet(walletPublicKeyHash)
	if err != nil {
		return nil, fmt.Errorf("cannot get on-chain data for wallet: [%v]", err)
	}

	// Valid case when the wallet doesn't have a main UTXO registered into
	// the Bridge.
	if walletChainData.MainUtxoHash == [32]byte{} {
		return nil, nil
	}

	// The wallet main UTXO registered in the Bridge almost always comes
	// from the latest BTC transaction made by the wallet. However, there may
	// be cases where the BTC transaction was made but their SPV proof is
	// not yet submitted to the Bridge thus the registered main UTXO points
	// to the second last BTC transaction. In theory, such a gap between
	// the actual latest BTC transaction and the registered main UTXO in
	// the Bridge may be even wider. To cover the worst possible cases, we
	// always take the five latest transactions made by the wallet for
	// consideration.
	transactions, err := btcChain.GetTransactionsForPublicKeyHash(walletPublicKeyHash, 5)
	if err != nil {
		return nil, fmt.Errorf("cannot get transactions history for wallet: [%v]", err)
	}

	walletP2PKH, err := bitcoin.PayToPublicKeyHash(walletPublicKeyHash)
	if err != nil {
		return nil, fmt.Errorf("cannot construct P2PKH for wallet: [%v]", err)
	}
	walletP2WPKH, err := bitcoin.PayToWitnessPublicKeyHash(walletPublicKeyHash)
	if err != nil {
		return nil, fmt.Errorf("cannot construct P2WPKH for wallet: [%v]", err)
	}

	// Start iterating from the latest transaction as the chance it matches
	// the wallet main UTXO is the highest.
	for i := len(transactions) - 1; i >= 0; i-- {
		transaction := transactions[i]

		// Iterate over transaction's outputs and find the one that targets
		// the wallet public key hash.
		for outputIndex, output := range transaction.Outputs {
			script := output.PublicKeyScript
			matchesWallet := bytes.Equal(script, walletP2PKH) ||
				bytes.Equal(script, walletP2WPKH)

			// Once the right output is found, check whether their hash
			// matches the main UTXO hash stored on-chain. If so, this
			// UTXO is the one we are looking for.
			if matchesWallet {
				utxo := &bitcoin.UnspentTransactionOutput{
					Outpoint: &bitcoin.TransactionOutpoint{
						TransactionHash: transaction.Hash(),
						OutputIndex:     uint32(outputIndex),
					},
					Value: output.Value,
				}

				if bridgeChain.ComputeMainUtxoHash(utxo) ==
					walletChainData.MainUtxoHash {
					return utxo, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("main UTXO not found")
}

// EnsureWalletSyncedBetweenChains makes sure all actions taken by the wallet
// on the Bitcoin chain are reflected in the host chain Bridge. This translates
// to two conditions that must be met:
// - The wallet main UTXO registered in the host chain Bridge comes from the
//   latest BTC transaction OR wallet main UTXO is unset and wallet's BTC
//   transaction history is empty. This condition ensures that all expected SPV
//   proofs of confirmed BTC transactions were submitted to the host chain Bridge
//   thus the wallet state held known to the Bridge matches the actual state
//   on the BTC chain.
// - There are no pending BTC transactions in the mempool. This condition
//   ensures the wallet doesn't currently perform any action on the BTC chain.
//   Such a transactions indicate a possible state change in the future
//   but their outcome cannot be determined at this stage so, the wallet
//   should not perform new actions at the moment.
func EnsureWalletSyncedBetweenChains(
	walletPublicKeyHash [20]byte,
	walletMainUtxo *bitcoin.UnspentTransactionOutput,
	btcChain bitcoin.Chain,
) error {
	// Take the recent transactions history for the wallet.
	history, err := btcChain.GetTransactionsForPublicKeyHash(walletPublicKeyHash, 5)
	if err != nil {
		return fmt.Errorf("cannot get transactions history: [%v]", err)
	}

	if walletMainUtxo != nil {
		// If the wallet main UTXO exists, the transaction history must
		// contain at least one item. If it is empty, something went
		// really wrong. This should never happen but check this scenario
		// just in case.
		if len(history) == 0 {
			return fmt.Errorf(
				"wallet main UTXO exists but there are no BTC " +
					"transactions produced by the wallet",
			)
		}

		// The transaction history is not empty for sure. Take the latest BTC
		// transaction from the history.
		latestTransaction := history[len(history)-1]

		// Make sure the wallet main UTXO comes from the latest transaction.
		// That means all expected SPV proofs were submitted to the Bridge.
		// If the wallet main UTXO transaction hash doesn't match the latest
		// transaction, that means the SPV proof for the latest transaction was
		// not submitted to the Bridge yet.
		//
		// Note that it is enough to check that the wallet main UTXO transaction
		// hash matches the latest transaction hash. There is no way the main
		// UTXO changes and the transaction hash stays the same. The Bridge
		// enforces that all wallet transactions form a sequence and refer
		// each other.
		if walletMainUtxo.Outpoint.TransactionHash != latestTransaction.Hash() {
			return fmt.Errorf(
				"wallet main UTXO doesn't come from the latest BTC transaction",
			)
		}
	} else {
		// If the wallet main UTXO doesn't exist, the transaction history must
		// be empty. If it is not, that could mean there is a Bitcoin transaction
		// produced by the wallet whose SPV proof was not submitted to
		// the Bridge yet.
		if len(history) != 0 {
			return fmt.Errorf(
				"wallet main UTXO doesn't exist but there are BTC " +
					"transactions produced by the wallet",
			)
		}
	}

	// Regardless of the main UTXO state, we need to make sure that
	// no pending wallet transactions exist in the mempool. That way,
	// we are handling a plenty of corner cases like transactions races
	// that could potentially lead to fraudulent transactions and funds loss.
	mempool, err := btcChain.GetMempoolForPublicKeyHash(walletPublicKeyHash)
	if err != nil {
		return fmt.Errorf("cannot get mempool: [%v]", err)
	}

	if len(mempool) != 0 {
		return fmt.Errorf("unconfirmed transactions exist in the mempool")
	}

	return nil
}

// signer represents a threshold signer of a tBTC wallet. A signer holds
// a wallet tECDSA private key share and is able to participate in the
// signing process.
type signer struct {
	// wallet points to the tBTC wallet this signer belongs to.
	wallet wallet

	// signingGroupMemberIndex indicates the signer position (seat) in the
	// wallet signing group. Since the final wallet signing group may differ
	// from the original group outputted by the sortition protocol
	// (see wallet.signingGroupOperators documentation for reference), the
	// signingGroupMemberIndex may differ from the member index using
	// during the DKG protocol as well. The value of this index is in the
	// [1, len(wallet.signingGroupOperators)] range.
	signingGroupMemberIndex group.MemberIndex

	// privateKeyShare is the tECDSA private key share required to participate
	// in the signing process.
	privateKeyShare *tecdsa.PrivateKeyShare
}

// newSigner constructs a new instance of the wallet's signer.
func newSigner(
	walletPublicKey *ecdsa.PublicKey,
	walletSigningGroupOperators []chain.Address,
	signingGroupMemberIndex group.MemberIndex,
	privateKeyShare *tecdsa.PrivateKeyShare,
) *signer {
	wallet := wallet{
		publicKey:             walletPublicKey,
		signingGroupOperators: walletSigningGroupOperators,
	}

	return &signer{
		wallet:                  wallet,
		signingGroupMemberIndex: signingGroupMemberIndex,
		privateKeyShare:         privateKeyShare,
	}
}

func (s *signer) String() string {
	return fmt.Sprintf(
		"signer with index [%v] of %s",
		s.signingGroupMemberIndex,
		&s.wallet,
	)
}
