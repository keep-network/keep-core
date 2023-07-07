package spv

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"math/big"
	"sync"
)

type localChain struct {
	mutex sync.Mutex

	blockCounter                          chain.BlockCounter
	pastRedemptionProposalSubmittedEvents map[[32]byte][]*tbtc.RedemptionProposalSubmittedEvent
	wallets                               map[[20]byte]*tbtc.WalletChainData
	pendingRedemptionRequests             map[[32]byte]*tbtc.RedemptionRequest
}

func newLocalChain() *localChain {
	return &localChain{
		pastRedemptionProposalSubmittedEvents: make(map[[32]byte][]*tbtc.RedemptionProposalSubmittedEvent),
		wallets:                               make(map[[20]byte]*tbtc.WalletChainData),
		pendingRedemptionRequests:             make(map[[32]byte]*tbtc.RedemptionRequest),
	}
}

func (lc *localChain) SubmitDepositSweepProofWithReimbursement(
	transaction *bitcoin.Transaction,
	proof *bitcoin.SpvProof,
	mainUTXO bitcoin.UnspentTransactionOutput,
	vault common.Address,
) error {
	panic("unsupported")
}

func (lc *localChain) GetDepositRequest(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) (*tbtc.DepositChainRequest, bool, error) {
	panic("unsupported")
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
	panic("unsupported")
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
	panic("unsupported")
}

func (lc *localChain) PastRedemptionProposalSubmittedEvents(
	filter *tbtc.RedemptionProposalSubmittedEventFilter,
) ([]*tbtc.RedemptionProposalSubmittedEvent, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	eventsKey, err := buildPastRedemptionProposalSubmittedEventsKey(filter)
	if err != nil {
		return nil, err
	}

	events, ok := lc.pastRedemptionProposalSubmittedEvents[eventsKey]
	if !ok {
		return nil, fmt.Errorf("no events for given filter")
	}

	return events, nil
}

func (lc *localChain) AddPastRedemptionProposalSubmittedEvent(
	filter *tbtc.RedemptionProposalSubmittedEventFilter,
	event *tbtc.RedemptionProposalSubmittedEvent,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	eventsKey, err := buildPastRedemptionProposalSubmittedEventsKey(filter)
	if err != nil {
		return err
	}

	lc.pastRedemptionProposalSubmittedEvents[eventsKey] = append(
		lc.pastRedemptionProposalSubmittedEvents[eventsKey],
		event,
	)

	return nil
}

func buildPastRedemptionProposalSubmittedEventsKey(
	filter *tbtc.RedemptionProposalSubmittedEventFilter,
) ([32]byte, error) {
	if filter == nil {
		return [32]byte{}, nil
	}

	var buffer bytes.Buffer

	startBlock := make([]byte, 8)
	binary.BigEndian.PutUint64(startBlock, filter.StartBlock)
	buffer.Write(startBlock)

	if filter.EndBlock != nil {
		endBlock := make([]byte, 8)
		binary.BigEndian.PutUint64(startBlock, *filter.EndBlock)
		buffer.Write(endBlock)
	}

	for _, coordinator := range filter.Coordinator {
		coordinatorBytes, err := hex.DecodeString(coordinator.String())
		if err != nil {
			return [32]byte{}, err
		}

		buffer.Write(coordinatorBytes)
	}

	buffer.Write(filter.WalletPublicKeyHash[:])

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
	panic("unsupported")
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
	panic("unsupported")
}

func (lc *localChain) ProofLength() (uint64, error) {
	panic("unsupported")
}

func (lc *localChain) GetCurrentAndPrevEpochDifficulty() (
	*big.Int,
	*big.Int,
	error,
) {
	panic("unsupported")
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
