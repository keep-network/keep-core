package coordinator_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/subscription"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

type localTbtcChain struct {
	depositRequestsMutex sync.Mutex
	depositRequests      map[[32]byte]*tbtc.DepositChainRequest

	pastDepositRevealedEventsMutex sync.Mutex
	pastDepositRevealedEvents      map[[32]byte][]*tbtc.DepositRevealedEvent

	pastNewWalletRegisteredEventsMutex  sync.Mutex
	pastNewWalletRegisteredEventsEvents map[[32]byte][]*tbtc.NewWalletRegisteredEvent
}

func newLocalTbtcChain() *localTbtcChain {
	return &localTbtcChain{
		depositRequests:                     make(map[[32]byte]*tbtc.DepositChainRequest),
		pastDepositRevealedEvents:           make(map[[32]byte][]*tbtc.DepositRevealedEvent),
		pastNewWalletRegisteredEventsEvents: make(map[[32]byte][]*tbtc.NewWalletRegisteredEvent),
	}
}

func (lc *localTbtcChain) PastDepositRevealedEvents(
	filter *tbtc.DepositRevealedEventFilter,
) ([]*tbtc.DepositRevealedEvent, error) {
	lc.pastDepositRevealedEventsMutex.Lock()
	defer lc.pastDepositRevealedEventsMutex.Unlock()

	eventsKey, err := buildPastDepositRevealedEventsKey(filter)
	if err != nil {
		return nil, err
	}

	events, ok := lc.pastDepositRevealedEvents[eventsKey]
	if !ok {
		return nil, fmt.Errorf("no events for given filter")
	}

	return events, nil
}

func (lc *localTbtcChain) addPastDepositRevealedEvent(
	filter *tbtc.DepositRevealedEventFilter,
	event *tbtc.DepositRevealedEvent,
) error {
	lc.pastDepositRevealedEventsMutex.Lock()
	defer lc.pastDepositRevealedEventsMutex.Unlock()

	eventsKey, err := buildPastDepositRevealedEventsKey(filter)
	if err != nil {
		return err
	}

	lc.pastDepositRevealedEvents[eventsKey] = append(
		lc.pastDepositRevealedEvents[eventsKey],
		event,
	)

	return nil
}

func buildPastDepositRevealedEventsKey(
	filter *tbtc.DepositRevealedEventFilter,
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

	for _, depositor := range filter.Depositor {
		depositorBytes, err := hex.DecodeString(depositor.String())
		if err != nil {
			return [32]byte{}, err
		}

		buffer.Write(depositorBytes)
	}

	for _, walletPublicKeyHash := range filter.WalletPublicKeyHash {
		buffer.Write(walletPublicKeyHash[:])
	}

	return sha256.Sum256(buffer.Bytes()), nil
}

func (lc *localTbtcChain) GetDepositRequest(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) (*tbtc.DepositChainRequest, error) {
	lc.depositRequestsMutex.Lock()
	defer lc.depositRequestsMutex.Unlock()

	requestKey := buildDepositRequestKey(fundingTxHash, fundingOutputIndex)

	request, ok := lc.depositRequests[requestKey]
	if !ok {
		return nil, fmt.Errorf("no request for given key")
	}

	return request, nil
}

func (lc *localTbtcChain) setDepositRequest(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
	request *tbtc.DepositChainRequest,
) {
	lc.depositRequestsMutex.Lock()
	defer lc.depositRequestsMutex.Unlock()

	requestKey := buildDepositRequestKey(fundingTxHash, fundingOutputIndex)

	lc.depositRequests[requestKey] = request
}

func (lc *localTbtcChain) PastNewWalletRegisteredEvents(
	filter *tbtc.NewWalletRegisteredEventFilter,
) ([]*tbtc.NewWalletRegisteredEvent, error) {
	lc.pastNewWalletRegisteredEventsMutex.Lock()
	defer lc.pastNewWalletRegisteredEventsMutex.Unlock()

	eventsKey, err := buildPastNewWalletRegisteredEventsKey(filter)
	if err != nil {
		return nil, err
	}

	events, ok := lc.pastNewWalletRegisteredEventsEvents[eventsKey]
	if !ok {
		return nil, fmt.Errorf("no events for given filter")
	}

	return events, nil
}

func (lc *localTbtcChain) addPastNewWalletRegisteredEvent(
	filter *tbtc.NewWalletRegisteredEventFilter,
	event *tbtc.NewWalletRegisteredEvent,
) error {
	lc.pastNewWalletRegisteredEventsMutex.Lock()
	defer lc.pastNewWalletRegisteredEventsMutex.Unlock()

	eventsKey, err := buildPastNewWalletRegisteredEventsKey(filter)
	if err != nil {
		return err
	}

	if _, ok := lc.pastNewWalletRegisteredEventsEvents[eventsKey]; !ok {
		lc.pastNewWalletRegisteredEventsEvents[eventsKey] = []*tbtc.NewWalletRegisteredEvent{}
	}

	lc.pastNewWalletRegisteredEventsEvents[eventsKey] = append(
		lc.pastNewWalletRegisteredEventsEvents[eventsKey],
		event,
	)

	return nil
}

func buildPastNewWalletRegisteredEventsKey(
	filter *tbtc.NewWalletRegisteredEventFilter,
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

	for _, ecdsaWalletID := range filter.EcdsaWalletID {
		buffer.Write(ecdsaWalletID[:])
	}

	for _, walletPublicKeyHash := range filter.WalletPublicKeyHash {
		buffer.Write(walletPublicKeyHash[:])
	}

	return sha256.Sum256(buffer.Bytes()), nil
}

func (lc *localTbtcChain) GetWallet(walletPublicKeyHash [20]byte) (*tbtc.WalletChainData, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) ComputeMainUtxoHash(mainUtxo *bitcoin.UnspentTransactionOutput) [32]byte {
	panic("unsupported")
}

func (lc *localTbtcChain) BuildDepositKey(fundingTxHash bitcoin.Hash, fundingOutputIndex uint32) *big.Int {
	depositKeyBytes := buildDepositRequestKey(fundingTxHash, fundingOutputIndex)

	return new(big.Int).SetBytes(depositKeyBytes[:])
}

func buildDepositRequestKey(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) [32]byte {
	fundingOutputIndexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(fundingOutputIndexBytes, fundingOutputIndex)

	return sha256.Sum256(append(fundingTxHash[:], fundingOutputIndexBytes...))
}

func (lc *localTbtcChain) GetDepositParameters() (
	dustThreshold uint64,
	treasuryFeeDivisor uint64,
	txMaxFee uint64,
	revealAheadPeriod uint32,
	err error,
) {
	panic("unsupported")
}

func (lc *localTbtcChain) OnHeartbeatRequestSubmitted(
	handler func(event *tbtc.HeartbeatRequestSubmittedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localTbtcChain) OnDepositSweepProposalSubmitted(
	func(event *tbtc.DepositSweepProposalSubmittedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localTbtcChain) PastDepositSweepProposalSubmittedEvents(
	filter *tbtc.DepositSweepProposalSubmittedEventFilter,
) ([]*tbtc.DepositSweepProposalSubmittedEvent, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) GetWalletLock(
	walletPublicKeyHash [20]byte,
) (time.Time, tbtc.WalletActionType, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) ValidateDepositSweepProposal(
	proposal *tbtc.DepositSweepProposal,
	depositsExtraInfo []struct {
		*tbtc.Deposit
		FundingTx *bitcoin.Transaction
	},
) error {
	panic("unsupported")
}

func (lc *localTbtcChain) SubmitDepositSweepProposalWithReimbursement(
	proposal *tbtc.DepositSweepProposal,
) error {
	panic("unsupported")
}

func (lc *localTbtcChain) GetDepositSweepMaxSize() (uint16, error) {
	panic("unsupported")
}
