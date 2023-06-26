package coordinator_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/coordinator"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

type localTbtcChain struct {
	mutex sync.Mutex

	depositRequests                 map[[32]byte]*tbtc.DepositChainRequest
	pastDepositRevealedEvents       map[[32]byte][]*tbtc.DepositRevealedEvent
	pastNewWalletRegisteredEvents   map[[32]byte][]*coordinator.NewWalletRegisteredEvent
	depositParameters               depositParameters
	depositSweepProposalValidations map[[32]byte]bool
	depositSweepProposals           []*tbtc.DepositSweepProposal
}

type depositParameters = struct {
	dustThreshold      uint64
	treasuryFeeDivisor uint64
	txMaxFee           uint64
	revealAheadPeriod  uint32
}

func newLocalTbtcChain() *localTbtcChain {
	return &localTbtcChain{
		depositRequests:                 make(map[[32]byte]*tbtc.DepositChainRequest),
		pastDepositRevealedEvents:       make(map[[32]byte][]*tbtc.DepositRevealedEvent),
		pastNewWalletRegisteredEvents:   make(map[[32]byte][]*coordinator.NewWalletRegisteredEvent),
		depositSweepProposalValidations: make(map[[32]byte]bool),
	}
}

func (lc *localTbtcChain) PastDepositRevealedEvents(
	filter *tbtc.DepositRevealedEventFilter,
) ([]*tbtc.DepositRevealedEvent, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

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
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

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
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

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
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	requestKey := buildDepositRequestKey(fundingTxHash, fundingOutputIndex)

	lc.depositRequests[requestKey] = request
}

func (lc *localTbtcChain) PastNewWalletRegisteredEvents(
	filter *coordinator.NewWalletRegisteredEventFilter,
) ([]*coordinator.NewWalletRegisteredEvent, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	eventsKey, err := buildPastNewWalletRegisteredEventsKey(filter)
	if err != nil {
		return nil, err
	}

	events, ok := lc.pastNewWalletRegisteredEvents[eventsKey]
	if !ok {
		return nil, fmt.Errorf("no events for given filter")
	}

	return events, nil
}

func (lc *localTbtcChain) addPastNewWalletRegisteredEvent(
	filter *coordinator.NewWalletRegisteredEventFilter,
	event *coordinator.NewWalletRegisteredEvent,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	eventsKey, err := buildPastNewWalletRegisteredEventsKey(filter)
	if err != nil {
		return err
	}

	if _, ok := lc.pastNewWalletRegisteredEvents[eventsKey]; !ok {
		lc.pastNewWalletRegisteredEvents[eventsKey] = []*coordinator.NewWalletRegisteredEvent{}
	}

	lc.pastNewWalletRegisteredEvents[eventsKey] = append(
		lc.pastNewWalletRegisteredEvents[eventsKey],
		event,
	)

	return nil
}

func buildPastNewWalletRegisteredEventsKey(
	filter *coordinator.NewWalletRegisteredEventFilter,
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

func (lc *localTbtcChain) PastRedemptionRequestedEvents(
	filter *tbtc.RedemptionRequestedEventFilter,
) ([]*tbtc.RedemptionRequestedEvent, error) {
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

func (lc *localTbtcChain) BuildRedemptionKey(
	walletPublicKeyHash [20]byte,
	redeemerOutputScript bitcoin.Script,
) (*big.Int, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) GetDepositParameters() (
	dustThreshold uint64,
	treasuryFeeDivisor uint64,
	txMaxFee uint64,
	revealAheadPeriod uint32,
	err error,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	return lc.depositParameters.dustThreshold,
		lc.depositParameters.treasuryFeeDivisor,
		lc.depositParameters.txMaxFee,
		lc.depositParameters.revealAheadPeriod,
		nil
}

func (lc *localTbtcChain) setDepositParameters(
	dustThreshold uint64,
	treasuryFeeDivisor uint64,
	txMaxFee uint64,
	revealAheadPeriod uint32,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.depositParameters = depositParameters{
		dustThreshold:      dustThreshold,
		treasuryFeeDivisor: treasuryFeeDivisor,
		txMaxFee:           txMaxFee,
		revealAheadPeriod:  revealAheadPeriod,
	}
}

func (lc *localTbtcChain) GetRedemptionParameters() (
	dustThreshold uint64,
	treasuryFeeDivisor uint64,
	txMaxFee uint64,
	txMaxTotalFee uint64,
	timeout uint32,
	timeoutSlashingAmount *big.Int,
	timeoutNotifierRewardMultiplier uint32,
	err error,
) {
	panic("unsupported")
}

func (lc *localTbtcChain) ValidateDepositSweepProposal(
	proposal *tbtc.DepositSweepProposal,
	depositsExtraInfo []struct {
		*tbtc.Deposit
		FundingTx *bitcoin.Transaction
	},
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	key, err := buildDepositSweepProposalValidationKey(proposal)
	if err != nil {
		return err
	}

	result, ok := lc.depositSweepProposalValidations[key]
	if !ok {
		return fmt.Errorf("validation result unknown")
	}

	if !result {
		return fmt.Errorf("validation failed")
	}

	return nil
}

func (lc *localTbtcChain) setDepositSweepProposalValidationResult(
	proposal *tbtc.DepositSweepProposal,
	depositsExtraInfo []struct {
		*tbtc.Deposit
		FundingTx *bitcoin.Transaction
	},
	result bool,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	key, err := buildDepositSweepProposalValidationKey(proposal)
	if err != nil {
		return err
	}

	lc.depositSweepProposalValidations[key] = result

	return nil
}

func buildDepositSweepProposalValidationKey(
	proposal *tbtc.DepositSweepProposal,
) ([32]byte, error) {
	var buffer bytes.Buffer

	buffer.Write(proposal.WalletPublicKeyHash[:])

	for _, deposit := range proposal.DepositsKeys {
		buffer.Write(deposit.FundingTxHash[:])

		fundingOutputIndex := make([]byte, 4)
		binary.BigEndian.PutUint32(fundingOutputIndex, deposit.FundingOutputIndex)
		buffer.Write(fundingOutputIndex)
	}

	buffer.Write(proposal.SweepTxFee.Bytes())

	return sha256.Sum256(buffer.Bytes()), nil
}

func (lc *localTbtcChain) SubmitDepositSweepProposalWithReimbursement(
	proposal *tbtc.DepositSweepProposal,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.depositSweepProposals = append(lc.depositSweepProposals, proposal)

	return nil
}

func (lc *localTbtcChain) SubmitRedemptionProposalWithReimbursement(
	proposal *tbtc.RedemptionProposal,
) error {
	panic("unsupported")
}

func (lc *localTbtcChain) GetDepositSweepMaxSize() (uint16, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) ValidateRedemptionProposal(
	proposal *tbtc.RedemptionProposal,
) error {
	panic("unsupported")
}

func (lc *localTbtcChain) GetRedemptionMaxSize() (uint16, error) {
	panic("unsupported")
}

func (lc *localTbtcChain) GetRedemptionRequestMinAge() (uint32, error) {
	panic("unsupported")
}
