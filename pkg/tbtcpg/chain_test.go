package tbtcpg

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

type depositParameters = struct {
	dustThreshold      uint64
	treasuryFeeDivisor uint64
	txMaxFee           uint64
	revealAheadPeriod  uint32
}

type redemptionParameters = struct {
	dustThreshold                   uint64
	treasuryFeeDivisor              uint64
	txMaxFee                        uint64
	txMaxTotalFee                   uint64
	timeout                         uint32
	timeoutSlashingAmount           *big.Int
	timeoutNotifierRewardMultiplier uint32
}

type LocalChain struct {
	mutex sync.Mutex

	depositRequests                 map[[32]byte]*tbtc.DepositChainRequest
	pastDepositRevealedEvents       map[[32]byte][]*tbtc.DepositRevealedEvent
	pastNewWalletRegisteredEvents   map[[32]byte][]*tbtc.NewWalletRegisteredEvent
	depositParameters               depositParameters
	depositSweepProposalValidations map[[32]byte]bool
	redemptionParameters            redemptionParameters
	redemptionRequestMinAge         uint32
	blockCounter                    chain.BlockCounter
	pastRedemptionRequestedEvents   map[[32]byte][]*tbtc.RedemptionRequestedEvent
	averageBlockTime                time.Duration
	pendingRedemptionRequests       map[[32]byte]*tbtc.RedemptionRequest
	redemptionProposalValidations   map[[32]byte]bool
	heartbeatProposalValidations    map[[16]byte]bool
}

func NewLocalChain() *LocalChain {
	return &LocalChain{
		depositRequests:                 make(map[[32]byte]*tbtc.DepositChainRequest),
		pastDepositRevealedEvents:       make(map[[32]byte][]*tbtc.DepositRevealedEvent),
		pastNewWalletRegisteredEvents:   make(map[[32]byte][]*tbtc.NewWalletRegisteredEvent),
		depositSweepProposalValidations: make(map[[32]byte]bool),
		pastRedemptionRequestedEvents:   make(map[[32]byte][]*tbtc.RedemptionRequestedEvent),
		pendingRedemptionRequests:       make(map[[32]byte]*tbtc.RedemptionRequest),
		redemptionProposalValidations:   make(map[[32]byte]bool),
		heartbeatProposalValidations:    make(map[[16]byte]bool),
	}
}

func (lc *LocalChain) PastDepositRevealedEvents(
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

func (lc *LocalChain) AddPastDepositRevealedEvent(
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

func (lc *LocalChain) GetDepositRequest(
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

func (lc *LocalChain) SetDepositRequest(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
	request *tbtc.DepositChainRequest,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	requestKey := buildDepositRequestKey(fundingTxHash, fundingOutputIndex)

	lc.depositRequests[requestKey] = request
}

func (lc *LocalChain) PastNewWalletRegisteredEvents(
	filter *tbtc.NewWalletRegisteredEventFilter,
) ([]*tbtc.NewWalletRegisteredEvent, error) {
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

func (lc *LocalChain) AddPastNewWalletRegisteredEvent(
	filter *tbtc.NewWalletRegisteredEventFilter,
	event *tbtc.NewWalletRegisteredEvent,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	eventsKey, err := buildPastNewWalletRegisteredEventsKey(filter)
	if err != nil {
		return err
	}

	if _, ok := lc.pastNewWalletRegisteredEvents[eventsKey]; !ok {
		lc.pastNewWalletRegisteredEvents[eventsKey] = []*tbtc.NewWalletRegisteredEvent{}
	}

	lc.pastNewWalletRegisteredEvents[eventsKey] = append(
		lc.pastNewWalletRegisteredEvents[eventsKey],
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

func (lc *LocalChain) PastRedemptionRequestedEvents(
	filter *tbtc.RedemptionRequestedEventFilter,
) ([]*tbtc.RedemptionRequestedEvent, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	eventsKey, err := buildPastRedemptionRequestedEventsKey(filter)
	if err != nil {
		return nil, err
	}

	events, ok := lc.pastRedemptionRequestedEvents[eventsKey]
	if !ok {
		return nil, fmt.Errorf("no events for given filter")
	}

	return events, nil
}

func (lc *LocalChain) AddPastRedemptionRequestedEvent(
	filter *tbtc.RedemptionRequestedEventFilter,
	event *tbtc.RedemptionRequestedEvent,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	eventsKey, err := buildPastRedemptionRequestedEventsKey(filter)
	if err != nil {
		return err
	}

	lc.pastRedemptionRequestedEvents[eventsKey] = append(
		lc.pastRedemptionRequestedEvents[eventsKey],
		event,
	)

	return nil
}

func buildPastRedemptionRequestedEventsKey(
	filter *tbtc.RedemptionRequestedEventFilter,
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

	for _, walletPublicKeyHash := range filter.WalletPublicKeyHash {
		buffer.Write(walletPublicKeyHash[:])
	}

	for _, redeemer := range filter.Redeemer {
		redeemerHex, err := hex.DecodeString(redeemer.String())
		if err != nil {
			return [32]byte{}, err
		}

		buffer.Write(redeemerHex)
	}

	return sha256.Sum256(buffer.Bytes()), nil
}

func (lc *LocalChain) BuildDepositKey(fundingTxHash bitcoin.Hash, fundingOutputIndex uint32) *big.Int {
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

func (lc *LocalChain) BuildRedemptionKey(
	walletPublicKeyHash [20]byte,
	redeemerOutputScript bitcoin.Script,
) (*big.Int, error) {
	redemptionKeyBytes := buildRedemptionRequestKey(
		walletPublicKeyHash,
		redeemerOutputScript,
	)

	return new(big.Int).SetBytes(redemptionKeyBytes[:]), nil
}

func buildRedemptionRequestKey(
	walletPublicKeyHash [20]byte,
	redeemerOutputScript bitcoin.Script,
) [32]byte {
	return sha256.Sum256(append(walletPublicKeyHash[:], redeemerOutputScript...))
}

func (lc *LocalChain) GetDepositParameters() (
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

func (lc *LocalChain) GetPendingRedemptionRequest(
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

func (lc *LocalChain) SetPendingRedemptionRequest(
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

func (lc *LocalChain) SetDepositParameters(
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

func (lc *LocalChain) GetRedemptionParameters() (
	dustThreshold uint64,
	treasuryFeeDivisor uint64,
	txMaxFee uint64,
	txMaxTotalFee uint64,
	timeout uint32,
	timeoutSlashingAmount *big.Int,
	timeoutNotifierRewardMultiplier uint32,
	err error,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	return lc.redemptionParameters.dustThreshold,
		lc.redemptionParameters.treasuryFeeDivisor,
		lc.redemptionParameters.txMaxFee,
		lc.redemptionParameters.txMaxTotalFee,
		lc.redemptionParameters.timeout,
		lc.redemptionParameters.timeoutSlashingAmount,
		lc.redemptionParameters.timeoutNotifierRewardMultiplier,
		nil
}

func (lc *LocalChain) SetRedemptionParameters(
	dustThreshold uint64,
	treasuryFeeDivisor uint64,
	txMaxFee uint64,
	txMaxTotalFee uint64,
	timeout uint32,
	timeoutSlashingAmount *big.Int,
	timeoutNotifierRewardMultiplier uint32,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.redemptionParameters = redemptionParameters{
		dustThreshold:                   dustThreshold,
		treasuryFeeDivisor:              treasuryFeeDivisor,
		txMaxFee:                        txMaxFee,
		txMaxTotalFee:                   txMaxTotalFee,
		timeout:                         timeout,
		timeoutSlashingAmount:           timeoutSlashingAmount,
		timeoutNotifierRewardMultiplier: timeoutNotifierRewardMultiplier,
	}
}

func (lc *LocalChain) ValidateDepositSweepProposal(
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

func (lc *LocalChain) SetDepositSweepProposalValidationResult(
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

func (lc *LocalChain) ValidateRedemptionProposal(
	proposal *tbtc.RedemptionProposal,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	key, err := buildRedemptionProposalValidationKey(proposal)
	if err != nil {
		return err
	}

	result, ok := lc.redemptionProposalValidations[key]
	if !ok {
		return fmt.Errorf("validation result unknown")
	}

	if !result {
		return fmt.Errorf("validation failed")
	}

	return nil
}

func (lc *LocalChain) SetRedemptionProposalValidationResult(
	proposal *tbtc.RedemptionProposal,
	result bool,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	key, err := buildRedemptionProposalValidationKey(proposal)
	if err != nil {
		return err
	}

	lc.redemptionProposalValidations[key] = result

	return nil
}

func (lc *LocalChain) ValidateHeartbeatProposal(
	walletPublicKeyHash [20]byte,
	proposal *tbtc.HeartbeatProposal,
) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	result, ok := lc.heartbeatProposalValidations[proposal.Message]
	if !ok {
		return fmt.Errorf("validation result unknown")
	}

	if !result {
		return fmt.Errorf("validation failed")
	}

	return nil
}

func (lc *LocalChain) SetHeartbeatProposalValidationResult(
	proposal *tbtc.HeartbeatProposal,
	result bool,
) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.heartbeatProposalValidations[proposal.Message] = result
}

func buildRedemptionProposalValidationKey(
	proposal *tbtc.RedemptionProposal,
) ([32]byte, error) {
	var buffer bytes.Buffer

	buffer.Write(proposal.WalletPublicKeyHash[:])

	for _, script := range proposal.RedeemersOutputScripts {
		buffer.Write(script)
	}

	buffer.Write(proposal.RedemptionTxFee.Bytes())

	return sha256.Sum256(buffer.Bytes()), nil
}

func (lc *LocalChain) GetRedemptionMaxSize() (uint16, error) {
	panic("unsupported")
}

func (lc *LocalChain) GetRedemptionRequestMinAge() (uint32, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	return lc.redemptionRequestMinAge, nil
}

func (lc *LocalChain) SetRedemptionRequestMinAge(redemptionRequestMinAge uint32) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.redemptionRequestMinAge = redemptionRequestMinAge
}

func (lc *LocalChain) GetDepositSweepMaxSize() (uint16, error) {
	panic("unsupported")
}

func (lc *LocalChain) BlockCounter() (chain.BlockCounter, error) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	return lc.blockCounter, nil
}

func (lc *LocalChain) SetBlockCounter(blockCounter chain.BlockCounter) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.blockCounter = blockCounter
}

func (lc *LocalChain) AverageBlockTime() time.Duration {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	return lc.averageBlockTime
}

func (lc *LocalChain) SetAverageBlockTime(averageBlockTime time.Duration) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.averageBlockTime = averageBlockTime
}

func (lc *LocalChain) GetWallet(walletPublicKeyHash [20]byte) (
	*tbtc.WalletChainData,
	error,
) {
	panic("unsupported")
}

func (lc *LocalChain) ComputeMainUtxoHash(mainUtxo *bitcoin.UnspentTransactionOutput) [32]byte {
	panic("unsupported")
}

type MockBlockCounter struct {
	mutex        sync.Mutex
	currentBlock uint64
}

func NewMockBlockCounter() *MockBlockCounter {
	return &MockBlockCounter{}
}

func (mbc *MockBlockCounter) WaitForBlockHeight(blockNumber uint64) error {
	panic("unsupported")
}

func (mbc *MockBlockCounter) BlockHeightWaiter(blockNumber uint64) (
	<-chan uint64,
	error,
) {
	panic("unsupported")
}

func (mbc *MockBlockCounter) CurrentBlock() (uint64, error) {
	mbc.mutex.Lock()
	defer mbc.mutex.Unlock()

	return mbc.currentBlock, nil
}

func (mbc *MockBlockCounter) SetCurrentBlock(block uint64) {
	mbc.mutex.Lock()
	defer mbc.mutex.Unlock()

	mbc.currentBlock = block
}

func (mbc *MockBlockCounter) WatchBlocks(ctx context.Context) <-chan uint64 {
	panic("unsupported")
}
