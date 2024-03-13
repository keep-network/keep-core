package tbtc

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/subscription"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
)

const localChainOperatorID = chain.OperatorID(1)

type localChain struct {
	dkgResultSubmissionHandlersMutex sync.Mutex
	dkgResultSubmissionHandlers      map[int]func(submission *DKGResultSubmittedEvent)

	dkgResultApprovalHandlersMutex sync.Mutex
	dkgResultApprovalHandlers      map[int]func(submission *DKGResultApprovedEvent)

	dkgResultApprovalGuard func() bool

	dkgResultChallengeHandlersMutex sync.Mutex
	dkgResultChallengeHandlers      map[int]func(submission *DKGResultChallengedEvent)

	dkgMutex       sync.Mutex
	dkgState       DKGState
	dkgResult      *DKGChainResult
	dkgResultValid bool

	walletsMutex sync.Mutex
	wallets      map[[20]byte]*WalletChainData

	blocksByTimestampMutex sync.Mutex
	blocksByTimestamp      map[uint64]uint64

	blocksHashesByNumberMutex sync.Mutex
	blocksHashesByNumber      map[uint64][32]byte

	pastDepositRevealedEventsMutex sync.Mutex
	pastDepositRevealedEvents      map[[32]byte][]*DepositRevealedEvent

	depositSweepProposalValidationsMutex sync.Mutex
	depositSweepProposalValidations      map[[32]byte]bool

	pendingRedemptionRequestsMutex sync.Mutex
	pendingRedemptionRequests      map[[32]byte]*RedemptionRequest

	redemptionProposalValidationsMutex sync.Mutex
	redemptionProposalValidations      map[[32]byte]bool

	movingFundsProposalValidationsMutex sync.Mutex
	movingFundsProposalValidations      map[[32]byte]bool

	heartbeatProposalValidationsMutex sync.Mutex
	heartbeatProposalValidations      map[[16]byte]bool

	depositRequestsMutex sync.Mutex
	depositRequests      map[[32]byte]*DepositChainRequest

	blockCounter       chain.BlockCounter
	operatorPrivateKey *operator.PrivateKey
}

func (lc *localChain) BlockCounter() (chain.BlockCounter, error) {
	return lc.blockCounter, nil
}

func (lc *localChain) Signing() chain.Signing {
	return local_v1.NewSigner(lc.operatorPrivateKey)
}

func (lc *localChain) OperatorKeyPair() (
	*operator.PrivateKey,
	*operator.PublicKey,
	error,
) {
	return lc.operatorPrivateKey, &lc.operatorPrivateKey.PublicKey, nil
}

func (lc *localChain) GetBlockNumberByTimestamp(timestamp uint64) (
	uint64,
	error,
) {
	lc.blocksByTimestampMutex.Lock()
	defer lc.blocksByTimestampMutex.Unlock()

	block, ok := lc.blocksByTimestamp[timestamp]
	if !ok {
		return 0, fmt.Errorf("block not found")
	}

	return block, nil
}

//lint:ignore U1000 This function can be useful for future.
func (lc *localChain) setBlockNumberByTimestamp(timestamp uint64, block uint64) {
	lc.blocksByTimestampMutex.Lock()
	defer lc.blocksByTimestampMutex.Unlock()

	lc.blocksByTimestamp[timestamp] = block
}

func (lc *localChain) GetBlockHashByNumber(blockNumber uint64) (
	[32]byte,
	error,
) {
	lc.blocksHashesByNumberMutex.Lock()
	defer lc.blocksHashesByNumberMutex.Unlock()

	blockHash, ok := lc.blocksHashesByNumber[blockNumber]
	if !ok {
		return [32]byte{}, fmt.Errorf("block not found")
	}

	return blockHash, nil
}

func (lc *localChain) setBlockHashByNumber(
	blockNumber uint64,
	blockHashString string,
) {
	lc.blocksHashesByNumberMutex.Lock()
	defer lc.blocksHashesByNumberMutex.Unlock()

	blockHashBytes, err := hex.DecodeString(blockHashString)
	if err != nil {
		panic(err)
	}

	var blockHash [32]byte
	copy(blockHash[:], blockHashBytes)

	lc.blocksHashesByNumber[blockNumber] = blockHash
}

func (lc *localChain) OperatorToStakingProvider() (chain.Address, bool, error) {
	panic("unsupported")
}

func (lc *localChain) EligibleStake(stakingProvider chain.Address) (*big.Int, error) {
	panic("unsupported")
}

func (lc *localChain) IsPoolLocked() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) IsOperatorInPool() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) IsOperatorUpToDate() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) JoinSortitionPool() error {
	panic("unsupported")
}

func (lc *localChain) UpdateOperatorStatus() error {
	panic("unsupported")
}

func (lc *localChain) IsEligibleForRewards() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) CanRestoreRewardEligibility() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) RestoreRewardEligibility() error {
	panic("unsupported")
}

func (lc *localChain) IsChaosnetActive() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) IsBetaOperator() (bool, error) {
	panic("unsupported")
}

func (lc *localChain) GetMovedFundsSweepRequest(
	movingFundsTxHash bitcoin.Hash,
	movingFundsTxOutpointIndex uint32,
) (*MovedFundsSweepRequest, bool, error) {
	panic("unsupported")
}

func (lc *localChain) GetOperatorID(
	operatorAddress chain.Address,
) (chain.OperatorID, error) {
	thisOperatorAddress, err := lc.operatorAddress()
	if err != nil {
		return 0, err
	}

	if thisOperatorAddress != operatorAddress {
		return 0, fmt.Errorf("local chain allows for one operator only")
	}

	return localChainOperatorID, nil
}

func (lc *localChain) SelectGroup() (*GroupSelectionResult, error) {
	panic("not implemented")
}

func (lc *localChain) OnDKGStarted(
	handler func(event *DKGStartedEvent),
) subscription.EventSubscription {
	panic("unsupported")
}

func (lc *localChain) PastDKGStartedEvents(
	filter *DKGStartedEventFilter,
) ([]*DKGStartedEvent, error) {
	panic("unsupported")
}

func (lc *localChain) OnDKGResultSubmitted(
	handler func(event *DKGResultSubmittedEvent),
) subscription.EventSubscription {
	lc.dkgResultSubmissionHandlersMutex.Lock()
	defer lc.dkgResultSubmissionHandlersMutex.Unlock()

	handlerID := generateHandlerID()
	lc.dkgResultSubmissionHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		lc.dkgResultSubmissionHandlersMutex.Lock()
		defer lc.dkgResultSubmissionHandlersMutex.Unlock()

		delete(lc.dkgResultSubmissionHandlers, handlerID)
	})
}

func (lc *localChain) OnDKGResultChallenged(
	handler func(event *DKGResultChallengedEvent),
) subscription.EventSubscription {
	lc.dkgResultChallengeHandlersMutex.Lock()
	defer lc.dkgResultChallengeHandlersMutex.Unlock()

	handlerID := generateHandlerID()
	lc.dkgResultChallengeHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		lc.dkgResultChallengeHandlersMutex.Lock()
		defer lc.dkgResultChallengeHandlersMutex.Unlock()

		delete(lc.dkgResultChallengeHandlers, handlerID)
	})
}

func (lc *localChain) OnDKGResultApproved(
	handler func(event *DKGResultApprovedEvent),
) subscription.EventSubscription {
	lc.dkgResultApprovalHandlersMutex.Lock()
	defer lc.dkgResultApprovalHandlersMutex.Unlock()

	handlerID := generateHandlerID()
	lc.dkgResultApprovalHandlers[handlerID] = handler

	return subscription.NewEventSubscription(func() {
		lc.dkgResultApprovalHandlersMutex.Lock()
		defer lc.dkgResultApprovalHandlersMutex.Unlock()

		delete(lc.dkgResultApprovalHandlers, handlerID)
	})
}

func (lc *localChain) startDKG() error {
	lc.dkgMutex.Lock()
	defer lc.dkgMutex.Unlock()

	if lc.dkgState != Idle {
		return fmt.Errorf("DKG not idle")
	}

	lc.dkgState = AwaitingResult

	return nil
}

func (lc *localChain) AssembleDKGResult(
	submitterMemberIndex group.MemberIndex,
	groupPublicKey *ecdsa.PublicKey,
	operatingMembersIndexes []group.MemberIndex,
	misbehavedMembersIndexes []group.MemberIndex,
	signatures map[group.MemberIndex][]byte,
	groupSelectionResult *GroupSelectionResult,
) (*DKGChainResult, error) {
	groupPublicKeyBytes := elliptic.Marshal(
		groupPublicKey.Curve,
		groupPublicKey.X,
		groupPublicKey.Y,
	)

	signingMembersIndexes := make([]group.MemberIndex, 0)
	signaturesConcatenation := make([]byte, 0)
	for memberIndex, signature := range signatures {
		signingMembersIndexes = append(signingMembersIndexes, memberIndex)
		signaturesConcatenation = append(signaturesConcatenation, signature...)
	}

	operatingOperatorsIDsBytes := make([]byte, 0)
	for _, operatingMemberID := range operatingMembersIndexes {
		operatorIDBytes := make([]byte, 4)
		operatorID := groupSelectionResult.OperatorsIDs[operatingMemberID-1]
		binary.BigEndian.PutUint32(operatorIDBytes, operatorID)

		operatingOperatorsIDsBytes = append(
			operatingOperatorsIDsBytes,
			operatorIDBytes...,
		)
	}

	return &DKGChainResult{
		SubmitterMemberIndex:     submitterMemberIndex,
		GroupPublicKey:           groupPublicKeyBytes,
		MisbehavedMembersIndexes: misbehavedMembersIndexes,
		Signatures:               signaturesConcatenation,
		SigningMembersIndexes:    signingMembersIndexes,
		Members:                  groupSelectionResult.OperatorsIDs,
		MembersHash:              sha3.Sum256(operatingOperatorsIDsBytes),
	}, nil
}

func (lc *localChain) SubmitDKGResult(
	dkgResult *DKGChainResult,
) error {
	lc.dkgResultSubmissionHandlersMutex.Lock()
	defer lc.dkgResultSubmissionHandlersMutex.Unlock()

	lc.dkgMutex.Lock()
	defer lc.dkgMutex.Unlock()

	if lc.dkgState != AwaitingResult {
		return fmt.Errorf("not awaiting DKG result")
	}

	blockNumber, err := lc.blockCounter.CurrentBlock()
	if err != nil {
		return fmt.Errorf("failed to get the current block")
	}

	resultHash := computeDkgChainResultHash(dkgResult)

	for _, handler := range lc.dkgResultSubmissionHandlers {
		handler(&DKGResultSubmittedEvent{
			Seed:        nil,
			ResultHash:  resultHash,
			Result:      dkgResult,
			BlockNumber: blockNumber,
		})
	}

	lc.dkgState = Challenge
	lc.dkgResult = dkgResult

	return nil
}

func (lc *localChain) GetDKGState() (DKGState, error) {
	lc.dkgMutex.Lock()
	defer lc.dkgMutex.Unlock()

	return lc.dkgState, nil
}

func (lc *localChain) CalculateDKGResultSignatureHash(
	groupPublicKey *ecdsa.PublicKey,
	misbehavedMembersIndexes []group.MemberIndex,
	startBlock uint64,
) (dkg.ResultSignatureHash, error) {
	if groupPublicKey == nil {
		return dkg.ResultSignatureHash{}, fmt.Errorf("group public key is nil")
	}

	encoded := fmt.Sprint(
		groupPublicKey,
		misbehavedMembersIndexes,
		startBlock,
	)

	return sha3.Sum256([]byte(encoded)), nil
}

func (lc *localChain) IsDKGResultValid(dkgResult *DKGChainResult) (bool, error) {
	lc.dkgMutex.Lock()
	defer lc.dkgMutex.Unlock()

	return lc.dkgResultValid, nil
}

func (lc *localChain) setDKGResultValidity(
	isValid bool,
) error {
	lc.dkgMutex.Lock()
	defer lc.dkgMutex.Unlock()

	lc.dkgResultValid = isValid

	return nil
}

func (lc *localChain) ChallengeDKGResult(dkgResult *DKGChainResult) error {
	lc.dkgResultChallengeHandlersMutex.Lock()
	defer lc.dkgResultChallengeHandlersMutex.Unlock()

	lc.dkgMutex.Lock()
	defer lc.dkgMutex.Unlock()

	if lc.dkgState != Challenge {
		return fmt.Errorf("not in DKG result challenge period")
	}

	if !reflect.DeepEqual(dkgResult, lc.dkgResult) {
		return fmt.Errorf("result does not match the submitted one")
	}

	if lc.dkgResultValid {
		return fmt.Errorf("submitted result is valid")
	}

	blockNumber, err := lc.blockCounter.CurrentBlock()
	if err != nil {
		return fmt.Errorf("failed to get the current block")
	}

	for _, handler := range lc.dkgResultChallengeHandlers {
		handler(&DKGResultChallengedEvent{
			ResultHash:  computeDkgChainResultHash(dkgResult),
			Challenger:  "",
			Reason:      "",
			BlockNumber: blockNumber,
		})
	}

	lc.dkgState = AwaitingResult
	lc.dkgResult = nil

	return nil
}

func (lc *localChain) ApproveDKGResult(dkgResult *DKGChainResult) error {
	lc.dkgResultApprovalHandlersMutex.Lock()
	defer lc.dkgResultApprovalHandlersMutex.Unlock()

	lc.dkgMutex.Lock()
	defer lc.dkgMutex.Unlock()

	if lc.dkgState != Challenge {
		return fmt.Errorf("not in DKG result challenge period")
	}

	if !reflect.DeepEqual(dkgResult, lc.dkgResult) {
		return fmt.Errorf("result does not match the submitted one")
	}

	if !lc.dkgResultValid {
		return fmt.Errorf("submitted result is invalid")
	}

	if lc.dkgResultApprovalGuard != nil && !lc.dkgResultApprovalGuard() {
		return fmt.Errorf("rejected by guard")
	}

	blockNumber, err := lc.blockCounter.CurrentBlock()
	if err != nil {
		return fmt.Errorf("failed to get the current block")
	}

	for _, handler := range lc.dkgResultApprovalHandlers {
		handler(&DKGResultApprovedEvent{
			ResultHash:  computeDkgChainResultHash(dkgResult),
			Approver:    "",
			BlockNumber: blockNumber,
		})
	}

	lc.dkgState = Idle
	lc.dkgResult = nil

	return nil
}

func (lc *localChain) DKGParameters() (*DKGParameters, error) {
	return &DKGParameters{
		SubmissionTimeoutBlocks:       10,
		ChallengePeriodBlocks:         15,
		ApprovePrecedencePeriodBlocks: 5,
	}, nil
}

func (lc *localChain) PastDepositRevealedEvents(
	filter *DepositRevealedEventFilter,
) ([]*DepositRevealedEvent, error) {
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

func (lc *localChain) setPastDepositRevealedEvents(
	filter *DepositRevealedEventFilter,
	events []*DepositRevealedEvent,
) error {
	lc.pastDepositRevealedEventsMutex.Lock()
	defer lc.pastDepositRevealedEventsMutex.Unlock()

	eventsKey, err := buildPastDepositRevealedEventsKey(filter)
	if err != nil {
		return err
	}

	lc.pastDepositRevealedEvents[eventsKey] = events

	return nil
}

func buildPastDepositRevealedEventsKey(
	filter *DepositRevealedEventFilter,
) ([32]byte, error) {
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

func (lc *localChain) GetPendingRedemptionRequest(
	walletPublicKeyHash [20]byte,
	redeemerOutputScript bitcoin.Script,
) (*RedemptionRequest, bool, error) {
	lc.pendingRedemptionRequestsMutex.Lock()
	defer lc.pendingRedemptionRequestsMutex.Unlock()

	requestKey := buildRedemptionRequestKey(walletPublicKeyHash, redeemerOutputScript)

	request, ok := lc.pendingRedemptionRequests[requestKey]
	if !ok {
		return nil, false, nil
	}

	return request, true, nil
}

func (lc *localChain) setPendingRedemptionRequest(
	walletPublicKeyHash [20]byte,
	request *RedemptionRequest,
) {
	lc.pendingRedemptionRequestsMutex.Lock()
	defer lc.pendingRedemptionRequestsMutex.Unlock()

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

func (lc *localChain) GetDepositRequest(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) (*DepositChainRequest, bool, error) {
	lc.depositRequestsMutex.Lock()
	defer lc.depositRequestsMutex.Unlock()

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
	request *DepositChainRequest,
) {
	lc.depositRequestsMutex.Lock()
	defer lc.depositRequestsMutex.Unlock()

	requestKey := buildDepositRequestKey(fundingTxHash, fundingOutputIndex)

	lc.depositRequests[requestKey] = request
}

func buildDepositRequestKey(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) [32]byte {
	buffer := make([]byte, 4)
	binary.BigEndian.PutUint32(buffer[:], fundingOutputIndex)

	return sha256.Sum256(append(fundingTxHash[:], buffer...))
}

func (lc *localChain) GetWallet(walletPublicKeyHash [20]byte) (
	*WalletChainData,
	error,
) {
	lc.walletsMutex.Lock()
	defer lc.walletsMutex.Unlock()

	walletChainData, ok := lc.wallets[walletPublicKeyHash]
	if !ok {
		return nil, fmt.Errorf("no wallet for given PKH")
	}

	return walletChainData, nil
}

func (lc *localChain) setWallet(
	walletPublicKeyHash [20]byte,
	walletChainData *WalletChainData,
) {
	lc.walletsMutex.Lock()
	defer lc.walletsMutex.Unlock()

	lc.wallets[walletPublicKeyHash] = walletChainData
}

func (lc *localChain) ComputeMainUtxoHash(
	mainUtxo *bitcoin.UnspentTransactionOutput,
) [32]byte {
	outputIndexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(outputIndexBytes, mainUtxo.Outpoint.OutputIndex)

	valueBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(valueBytes, uint64(mainUtxo.Value))

	mainUtxoHash := sha256.Sum256(
		append(
			append(
				mainUtxo.Outpoint.TransactionHash[:],
				outputIndexBytes...,
			), valueBytes...,
		),
	)

	return mainUtxoHash
}

func (lc *localChain) ComputeMovingFundsCommitmentHash(targetWallets [][20]byte) [32]byte {
	packedWallets := []byte{}

	for _, wallet := range targetWallets {
		packedWallets = append(packedWallets, wallet[:]...)
		// Each wallet hash must be padded with 12 zero bytes following the
		// actual hash.
		packedWallets = append(packedWallets, make([]byte, 12)...)
	}

	return crypto.Keccak256Hash(packedWallets)
}

func (lc *localChain) operatorAddress() (chain.Address, error) {
	_, operatorPublicKey, err := lc.OperatorKeyPair()
	if err != nil {
		return "", err
	}

	return lc.Signing().PublicKeyToAddress(operatorPublicKey)
}

func (lc *localChain) GetWalletParameters() (
	creationPeriod uint32,
	creationMinBtcBalance uint64,
	creationMaxBtcBalance uint64,
	closureMinBtcBalance uint64,
	maxAge uint32,
	maxBtcTransfer uint64,
	closingPeriod uint32,
	err error,
) {
	panic("unsupported")
}

func (lc *localChain) ValidateDepositSweepProposal(
	walletPublicKeyHash [20]byte,
	proposal *DepositSweepProposal,
	depositsExtraInfo []struct {
		*Deposit
		FundingTx *bitcoin.Transaction
	},
) error {
	lc.depositSweepProposalValidationsMutex.Lock()
	defer lc.depositSweepProposalValidationsMutex.Unlock()

	key, err := buildDepositSweepProposalValidationKey(
		walletPublicKeyHash,
		proposal,
		depositsExtraInfo,
	)
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

func (lc *localChain) setDepositSweepProposalValidationResult(
	walletPublicKeyHash [20]byte,
	proposal *DepositSweepProposal,
	depositsExtraInfo []struct {
		*Deposit
		FundingTx *bitcoin.Transaction
	},
	result bool,
) error {
	lc.depositSweepProposalValidationsMutex.Lock()
	defer lc.depositSweepProposalValidationsMutex.Unlock()

	key, err := buildDepositSweepProposalValidationKey(
		walletPublicKeyHash,
		proposal,
		depositsExtraInfo,
	)
	if err != nil {
		return err
	}

	lc.depositSweepProposalValidations[key] = result

	return nil
}

func buildDepositSweepProposalValidationKey(
	walletPublicKeyHash [20]byte,
	proposal *DepositSweepProposal,
	depositsExtraInfo []struct {
		*Deposit
		FundingTx *bitcoin.Transaction
	},
) ([32]byte, error) {
	var buffer bytes.Buffer

	buffer.Write(walletPublicKeyHash[:])

	for _, deposit := range proposal.DepositsKeys {
		buffer.Write(deposit.FundingTxHash[:])

		fundingOutputIndex := make([]byte, 4)
		binary.BigEndian.PutUint32(fundingOutputIndex, deposit.FundingOutputIndex)
		buffer.Write(fundingOutputIndex)
	}

	buffer.Write(proposal.SweepTxFee.Bytes())

	for _, extra := range depositsExtraInfo {
		depositScript, err := extra.Deposit.Script()
		if err != nil {
			return [32]byte{}, err
		}

		buffer.Write(depositScript)

		fundingTxHash := extra.FundingTx.Hash()
		buffer.Write(fundingTxHash[:])
	}

	return sha256.Sum256(buffer.Bytes()), nil
}

func (lc *localChain) ValidateRedemptionProposal(
	walletPublicKeyHash [20]byte,
	proposal *RedemptionProposal,
) error {
	lc.redemptionProposalValidationsMutex.Lock()
	defer lc.redemptionProposalValidationsMutex.Unlock()

	key, err := buildRedemptionProposalValidationKey(
		walletPublicKeyHash,
		proposal,
	)
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

func (lc *localChain) setRedemptionProposalValidationResult(
	walletPublicKeyHash [20]byte,
	proposal *RedemptionProposal,
	result bool,
) error {
	lc.redemptionProposalValidationsMutex.Lock()
	defer lc.redemptionProposalValidationsMutex.Unlock()

	key, err := buildRedemptionProposalValidationKey(
		walletPublicKeyHash,
		proposal,
	)
	if err != nil {
		return err
	}

	lc.redemptionProposalValidations[key] = result

	return nil
}

func buildRedemptionProposalValidationKey(
	walletPublicKeyHash [20]byte,
	proposal *RedemptionProposal,
) ([32]byte, error) {
	var buffer bytes.Buffer

	buffer.Write(walletPublicKeyHash[:])

	for _, script := range proposal.RedeemersOutputScripts {
		buffer.Write(script)
	}

	buffer.Write(proposal.RedemptionTxFee.Bytes())

	return sha256.Sum256(buffer.Bytes()), nil
}

func (lc *localChain) ValidateHeartbeatProposal(
	walletPublicKeyHash [20]byte,
	proposal *HeartbeatProposal,
) error {
	lc.heartbeatProposalValidationsMutex.Lock()
	defer lc.heartbeatProposalValidationsMutex.Unlock()

	result, ok := lc.heartbeatProposalValidations[proposal.Message]
	if !ok {
		return fmt.Errorf("validation result unknown")
	}

	if !result {
		return fmt.Errorf("validation failed")
	}

	return nil
}

func (lc *localChain) setHeartbeatProposalValidationResult(
	proposal *HeartbeatProposal,
	result bool,
) {
	lc.heartbeatProposalValidationsMutex.Lock()
	defer lc.heartbeatProposalValidationsMutex.Unlock()

	lc.heartbeatProposalValidations[proposal.Message] = result
}

func (lc *localChain) ValidateMovingFundsProposal(
	walletPublicKeyHash [20]byte,
	mainUTXO *bitcoin.UnspentTransactionOutput,
	proposal *MovingFundsProposal,
) error {
	lc.movingFundsProposalValidationsMutex.Lock()
	defer lc.movingFundsProposalValidationsMutex.Unlock()

	key, err := buildMovingFundsProposalValidationKey(
		walletPublicKeyHash,
		mainUTXO,
		proposal,
	)
	if err != nil {
		return err
	}

	result, ok := lc.movingFundsProposalValidations[key]
	if !ok {
		return fmt.Errorf("validation result unknown")
	}

	if !result {
		return fmt.Errorf("validation failed")
	}

	return nil
}

func (lc *localChain) setMovingFundsProposalValidationResult(
	walletPublicKeyHash [20]byte,
	mainUTXO *bitcoin.UnspentTransactionOutput,
	proposal *MovingFundsProposal,
	result bool,
) error {
	lc.movingFundsProposalValidationsMutex.Lock()
	defer lc.movingFundsProposalValidationsMutex.Unlock()

	key, err := buildMovingFundsProposalValidationKey(
		walletPublicKeyHash,
		mainUTXO,
		proposal,
	)
	if err != nil {
		return err
	}

	lc.movingFundsProposalValidations[key] = result

	return nil
}

func buildMovingFundsProposalValidationKey(
	walletPublicKeyHash [20]byte,
	mainUTXO *bitcoin.UnspentTransactionOutput,
	proposal *MovingFundsProposal,
) ([32]byte, error) {
	var buffer bytes.Buffer

	buffer.Write(walletPublicKeyHash[:])

	buffer.Write(mainUTXO.Outpoint.TransactionHash[:])
	binary.Write(&buffer, binary.BigEndian, mainUTXO.Outpoint.OutputIndex)
	binary.Write(&buffer, binary.BigEndian, mainUTXO.Value)

	for _, wallet := range proposal.TargetWallets {
		buffer.Write(wallet[:])
	}

	buffer.Write(proposal.MovingFundsTxFee.Bytes())

	return sha256.Sum256(buffer.Bytes()), nil
}

// Connect sets up the local chain.
func Connect(blockTime ...time.Duration) *localChain {
	operatorPrivateKey, _, err := operator.GenerateKeyPair(local_v1.DefaultCurve)
	if err != nil {
		panic(err)
	}

	return ConnectWithKey(operatorPrivateKey, blockTime...)
}

// ConnectWithKey sets up the local chain using the provided operator private
// key.
func ConnectWithKey(
	operatorPrivateKey *operator.PrivateKey,
	blockTime ...time.Duration,
) *localChain {
	blockCounter, _ := local_v1.BlockCounter(blockTime...)

	localChain := &localChain{
		dkgResultSubmissionHandlers: make(
			map[int]func(submission *DKGResultSubmittedEvent),
		),
		dkgResultApprovalHandlers: make(
			map[int]func(submission *DKGResultApprovedEvent),
		),
		dkgResultChallengeHandlers: make(
			map[int]func(submission *DKGResultChallengedEvent),
		),
		wallets:                         make(map[[20]byte]*WalletChainData),
		blocksByTimestamp:               make(map[uint64]uint64),
		blocksHashesByNumber:            make(map[uint64][32]byte),
		pastDepositRevealedEvents:       make(map[[32]byte][]*DepositRevealedEvent),
		depositSweepProposalValidations: make(map[[32]byte]bool),
		pendingRedemptionRequests:       make(map[[32]byte]*RedemptionRequest),
		redemptionProposalValidations:   make(map[[32]byte]bool),
		movingFundsProposalValidations:  make(map[[32]byte]bool),
		heartbeatProposalValidations:    make(map[[16]byte]bool),
		depositRequests:                 make(map[[32]byte]*DepositChainRequest),
		blockCounter:                    blockCounter,
		operatorPrivateKey:              operatorPrivateKey,
	}

	return localChain
}

func computeDkgChainResultHash(result *DKGChainResult) DKGChainResultHash {
	return sha3.Sum256(result.GroupPublicKey)
}

func generateHandlerID() int {
	// #nosec G404 (insecure random number source (rand))
	// Local chain implementation doesn't require secure randomness.
	return rand.Int()
}
