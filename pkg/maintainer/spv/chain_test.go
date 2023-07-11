package spv

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

type submittedRedemptionProof struct {
	transaction         *bitcoin.Transaction
	proof               *bitcoin.SpvProof
	mainUTXO            bitcoin.UnspentTransactionOutput
	walletPublicKeyHash [20]byte
}

type submittedDepositSweepProof struct {
	transaction *bitcoin.Transaction
	proof       *bitcoin.SpvProof
	mainUTXO    bitcoin.UnspentTransactionOutput
	vault       common.Address
}

type localChain struct {
	mutex sync.Mutex

	blockCounter                            chain.BlockCounter
	pastDepositSweepProposalSubmittedEvents map[[32]byte][]*tbtc.DepositSweepProposalSubmittedEvent
	pastRedemptionProposalSubmittedEvents   map[[32]byte][]*tbtc.RedemptionProposalSubmittedEvent
	wallets                                 map[[20]byte]*tbtc.WalletChainData
	depositRequests                         map[[32]byte]*tbtc.DepositChainRequest
	pendingRedemptionRequests               map[[32]byte]*tbtc.RedemptionRequest
	submittedRedemptionProofs               []*submittedRedemptionProof
	submittedDepositSweepProofs             []*submittedDepositSweepProof

	txProofDifficultyFactor *big.Int
	currentEpoch            uint64
	currentEpochDifficulty  *big.Int
	previousEpochDifficulty *big.Int
}

func newLocalChain() *localChain {
	return &localChain{
		pastDepositSweepProposalSubmittedEvents: make(map[[32]byte][]*tbtc.DepositSweepProposalSubmittedEvent),
		pastRedemptionProposalSubmittedEvents:   make(map[[32]byte][]*tbtc.RedemptionProposalSubmittedEvent),
		wallets:                                 make(map[[20]byte]*tbtc.WalletChainData),
		depositRequests:                         make(map[[32]byte]*tbtc.DepositChainRequest),
		pendingRedemptionRequests:               make(map[[32]byte]*tbtc.RedemptionRequest),
		submittedRedemptionProofs:               make([]*submittedRedemptionProof, 0),
		submittedDepositSweepProofs:             make([]*submittedDepositSweepProof, 0),
	}
}

func (lc *localChain) SubmitDepositSweepProofWithReimbursement(
	transaction *bitcoin.Transaction,
	proof *bitcoin.SpvProof,
	mainUTXO bitcoin.UnspentTransactionOutput,
	vault common.Address,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.submittedDepositSweepProofs = append(
		lc.submittedDepositSweepProofs,
		&submittedDepositSweepProof{
			transaction: transaction,
			proof:       proof,
			mainUTXO:    mainUTXO,
			vault:       vault,
		},
	)

	return nil
}

func (lc *localChain) getSubmittedDepositSweepProofs() []*submittedDepositSweepProof {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	return lc.submittedDepositSweepProofs
}

func (lc *localChain) GetDepositRequest(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) (*tbtc.DepositChainRequest, bool, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	requestKey := buildDepositRequestKey(fundingTxHash, fundingOutputIndex)

	request, ok := lc.depositRequests[requestKey]
	if !ok {
		return nil, false, nil
	}

	return request, true, nil
}

func (lc *localChain) setDepositRequest(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
	depositRequest *tbtc.DepositChainRequest,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	requestKey := buildDepositRequestKey(fundingTxHash, fundingOutputIndex)
	lc.depositRequests[requestKey] = depositRequest
}

func buildDepositRequestKey(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) [32]byte {
	fundingOutputIndexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(fundingOutputIndexBytes, fundingOutputIndex)

	return sha256.Sum256(append(fundingTxHash[:], fundingOutputIndexBytes...))
}

func (lc *localChain) GetWallet(walletPublicKeyHash [20]byte) (
	*tbtc.WalletChainData,
	error,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	walletChainData, ok := lc.wallets[walletPublicKeyHash]
	if !ok {
		return nil, fmt.Errorf("no wallet for given PKH")
	}

	return walletChainData, nil
}

func (lc *localChain) setWallet(
	walletPublicKeyHash [20]byte,
	walletChainData *tbtc.WalletChainData,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.wallets[walletPublicKeyHash] = walletChainData
}

func (lc *localChain) ComputeMainUtxoHash(
	mainUtxo *bitcoin.UnspentTransactionOutput,
) [32]byte {
	var buffer bytes.Buffer

	buffer.Write(mainUtxo.Outpoint.TransactionHash[:])

	outputIndex := make([]byte, 4)
	binary.BigEndian.PutUint32(outputIndex, mainUtxo.Outpoint.OutputIndex)
	buffer.Write(outputIndex)

	value := make([]byte, 8)
	binary.BigEndian.PutUint64(value, uint64(mainUtxo.Value))
	buffer.Write(value)

	return sha256.Sum256(buffer.Bytes())
}

func (lc *localChain) TxProofDifficultyFactor() (*big.Int, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	if lc.txProofDifficultyFactor == nil {
		return nil, fmt.Errorf("transaction proof difficulty factor not set")
	}

	return lc.txProofDifficultyFactor, nil
}

func (lc *localChain) BlockCounter() (chain.BlockCounter, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	return lc.blockCounter, nil
}

func (lc *localChain) setBlockCounter(blockCounter chain.BlockCounter) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.blockCounter = blockCounter
}

func (lc *localChain) PastDepositSweepProposalSubmittedEvents(
	filter *tbtc.DepositSweepProposalSubmittedEventFilter,
) (
	[]*tbtc.DepositSweepProposalSubmittedEvent,
	error,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	var eventsKey [32]byte
	var err error

	if filter != nil {
		eventsKey, err = buildPastProposalSubmittedEventsKey(
			filter.StartBlock,
			filter.EndBlock,
			filter.Coordinator,
			filter.WalletPublicKeyHash,
		)
		if err != nil {
			return nil, err
		}
	}

	events, ok := lc.pastDepositSweepProposalSubmittedEvents[eventsKey]
	if !ok {
		return nil, fmt.Errorf("no events for given filter")
	}

	return events, nil
}

func (lc *localChain) PastRedemptionProposalSubmittedEvents(
	filter *tbtc.RedemptionProposalSubmittedEventFilter,
) ([]*tbtc.RedemptionProposalSubmittedEvent, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	var eventsKey [32]byte
	var err error

	if filter != nil {
		eventsKey, err = buildPastProposalSubmittedEventsKey(
			filter.StartBlock,
			filter.EndBlock,
			filter.Coordinator,
			filter.WalletPublicKeyHash,
		)
		if err != nil {
			return nil, err
		}
	}

	events, ok := lc.pastRedemptionProposalSubmittedEvents[eventsKey]
	if !ok {
		return nil, fmt.Errorf("no events for given filter")
	}

	return events, nil
}

func (lc *localChain) AddPastDepositSweepProposalSubmittedEvent(
	filter *tbtc.DepositSweepProposalSubmittedEventFilter,
	event *tbtc.DepositSweepProposalSubmittedEvent,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	var eventsKey [32]byte
	var err error

	if filter != nil {
		eventsKey, err = buildPastProposalSubmittedEventsKey(
			filter.StartBlock,
			filter.EndBlock,
			filter.Coordinator,
			filter.WalletPublicKeyHash,
		)
		if err != nil {
			return err
		}
	}

	lc.pastDepositSweepProposalSubmittedEvents[eventsKey] = append(
		lc.pastDepositSweepProposalSubmittedEvents[eventsKey],
		event,
	)

	return nil
}

func (lc *localChain) AddPastRedemptionProposalSubmittedEvent(
	filter *tbtc.RedemptionProposalSubmittedEventFilter,
	event *tbtc.RedemptionProposalSubmittedEvent,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	var eventsKey [32]byte
	var err error

	if filter != nil {
		eventsKey, err = buildPastProposalSubmittedEventsKey(
			filter.StartBlock,
			filter.EndBlock,
			filter.Coordinator,
			filter.WalletPublicKeyHash,
		)
		if err != nil {
			return err
		}
	}

	lc.pastRedemptionProposalSubmittedEvents[eventsKey] = append(
		lc.pastRedemptionProposalSubmittedEvents[eventsKey],
		event,
	)

	return nil
}

func buildPastProposalSubmittedEventsKey(
	startBlock uint64,
	endBlock *uint64,
	coordinators []chain.Address,
	walletPublicKeyHash [20]byte,
) ([32]byte, error) {
	var buffer bytes.Buffer

	startBlockBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(startBlockBytes, startBlock)
	buffer.Write(startBlockBytes)

	if endBlock != nil {
		endBlockBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(endBlockBytes, *endBlock)
		buffer.Write(endBlockBytes)
	}

	for _, coordinator := range coordinators {
		coordinatorBytes, err := hex.DecodeString(coordinator.String())
		if err != nil {
			return [32]byte{}, err
		}

		buffer.Write(coordinatorBytes)
	}

	buffer.Write(walletPublicKeyHash[:])

	return sha256.Sum256(buffer.Bytes()), nil
}

func (lc *localChain) GetPendingRedemptionRequest(
	walletPublicKeyHash [20]byte,
	redeemerOutputScript bitcoin.Script,
) (*tbtc.RedemptionRequest, bool, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	requestKey := buildRedemptionRequestKey(walletPublicKeyHash, redeemerOutputScript)

	request, ok := lc.pendingRedemptionRequests[requestKey]
	if !ok {
		return nil, false, nil
	}

	return request, true, nil
}

func (lc *localChain) setPendingRedemptionRequest(
	walletPublicKeyHash [20]byte,
	request *tbtc.RedemptionRequest,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	requestKey := buildRedemptionRequestKey(
		walletPublicKeyHash,
		request.RedeemerOutputScript,
	)

	lc.pendingRedemptionRequests[requestKey] = request
}

func buildRedemptionRequestKey(
	walletPublicKeyHash [20]byte,
	redeemerOutputScript bitcoin.Script,
) [32]byte {
	return sha256.Sum256(append(walletPublicKeyHash[:], redeemerOutputScript...))
}

func (lc *localChain) SubmitRedemptionProofWithReimbursement(
	transaction *bitcoin.Transaction,
	proof *bitcoin.SpvProof,
	mainUTXO bitcoin.UnspentTransactionOutput,
	walletPublicKeyHash [20]byte,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.submittedRedemptionProofs = append(
		lc.submittedRedemptionProofs,
		&submittedRedemptionProof{
			transaction:         transaction,
			proof:               proof,
			mainUTXO:            mainUTXO,
			walletPublicKeyHash: walletPublicKeyHash,
		},
	)

	return nil
}

func (lc *localChain) getSubmittedRedemptionProofs() []*submittedRedemptionProof {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	return lc.submittedRedemptionProofs
}

func (lc *localChain) Ready() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) IsAuthorized(address chain.Address) (bool, error) {
	panic("unsupported")
}

func (lc *localChain) IsAuthorizedForRefund(address chain.Address) (
	bool,
	error,
) {
	panic("unsupported")
}

func (lc *localChain) Signing() chain.Signing {
	panic("unsupported")
}

func (lc *localChain) Retarget(headers []*bitcoin.BlockHeader) error {
	panic("unsupported")
}

func (lc *localChain) RetargetWithRefund(headers []*bitcoin.BlockHeader) error {
	panic("unsupported")
}

func (lc *localChain) CurrentEpoch() (uint64, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	return lc.currentEpoch, nil
}

func (lc *localChain) ProofLength() (uint64, error) {
	panic("unsupported")
}

func (lc *localChain) GetCurrentAndPrevEpochDifficulty() (
	*big.Int,
	*big.Int,
	error,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	if lc.currentEpochDifficulty == nil || lc.previousEpochDifficulty == nil {
		return nil, nil, fmt.Errorf("epoch difficulties not set")
	}

	return lc.currentEpochDifficulty, lc.previousEpochDifficulty, nil
}

func (lc *localChain) setTxProofDifficultyFactor(
	txProofDifficultyFactor *big.Int,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.txProofDifficultyFactor = txProofDifficultyFactor
}

func (lc *localChain) setCurrentEpoch(currentEpoch uint64) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.currentEpoch = currentEpoch
}

func (lc *localChain) setCurrentAndPrevEpochDifficulty(
	previousEpochDifficulty, currentEpochDifficulty *big.Int,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.currentEpochDifficulty = currentEpochDifficulty
	lc.previousEpochDifficulty = previousEpochDifficulty
}

type mockBlockCounter struct {
	mutex        sync.Mutex
	currentBlock uint64
}

func newMockBlockCounter() *mockBlockCounter {
	return &mockBlockCounter{}
}

func (mbc *mockBlockCounter) WaitForBlockHeight(blockNumber uint64) error {
	panic("unsupported")
}

func (mbc *mockBlockCounter) BlockHeightWaiter(blockNumber uint64) (
	<-chan uint64,
	error,
) {
	panic("unsupported")
}

func (mbc *mockBlockCounter) CurrentBlock() (uint64, error) {
	mbc.mutex.Lock()
	defer mbc.mutex.Unlock()

	return mbc.currentBlock, nil
}

func (mbc *mockBlockCounter) SetCurrentBlock(block uint64) {
	mbc.mutex.Lock()
	defer mbc.mutex.Unlock()

	mbc.currentBlock = block
}

func (mbc *mockBlockCounter) WatchBlocks(ctx context.Context) <-chan uint64 {
	panic("unsupported")
}
