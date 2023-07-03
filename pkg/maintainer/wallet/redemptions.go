package wallet

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/keep-network/keep-core/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

func (wm *walletMaintainer) runRedemptionTask(ctx context.Context) error {
	redemptionMaxSize, err := wm.chain.GetRedemptionMaxSize()
	if err != nil {
		return fmt.Errorf("failed to get redemption max size: [%w]", err)
	}

	walletsPendingRedemptions, err := FindPendingRedemptions(
		wm.chain,
		PendingRedemptionsFilter{
			WalletPublicKeyHashes:   nil,
			WalletsLimit:            wm.config.RedemptionWalletsLimit,
			RedemptionRequestsLimit: redemptionMaxSize,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to find pending redemption requests: [%w]", err)
	}

	if len(walletsPendingRedemptions) == 0 {
		logger.Info("no pending redemption requests")
		return nil
	}

	for walletPublicKeyHash, redeemersOutputScripts := range walletsPendingRedemptions {
		err = wm.runIfWalletUnlocked(
			ctx,
			walletPublicKeyHash,
			tbtc.Redemption,
			func() error {
				return ProposeRedemption(
					wm.chain,
					wm.btcChain,
					walletPublicKeyHash,
					0,
					redeemersOutputScripts,
					false,
				)
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// RedemptionRequest represents a redemption request.
type RedemptionRequest struct {
	WalletPublicKeyHash  [20]byte
	RedemptionKey        string
	RedeemerOutputScript bitcoin.Script
	RequestedAt          time.Time
}

// PendingRedemptionsFilter defines some criteria that are used to filter
// pending redemption requests returned by the FindPendingRedemptions function.
type PendingRedemptionsFilter struct {
	// WalletPublicKeyHashes limits the search space to specific wallets.
	// The nil/empty value of this field means all wallets will be taken into
	// account, starting from the oldest one.
	WalletPublicKeyHashes [][20]byte

	// WalletsLimit limits the total number of wallets, starting from the
	// oldest one. The value of 0 means there is no wallets limit.
	WalletsLimit uint16

	// RedemptionRequestsLimit limits the number of redemptions requests per
	// single wallet. The value of 0 means there is no requests limit per wallet.
	RedemptionRequestsLimit uint16
}

func (prf PendingRedemptionsFilter) String() string {
	wallets := make([]string, len(prf.WalletPublicKeyHashes))
	for i, wallet := range prf.WalletPublicKeyHashes {
		wallets[i] = hex.EncodeToString(wallet[:])
	}

	return fmt.Sprintf(
		"wallets: [%s], wallets limit: [%v], redemption requests limit: [%v]",
		strings.Join(wallets, ", "),
		prf.WalletsLimit,
		prf.RedemptionRequestsLimit,
	)
}

// FindPendingRedemptions finds pending redemptions requests according to
// the provided filter. The returned value is a map, where the key is
// a 20-byte public key hash of a specific wallet and the value is a list
// of pending requests targeting this wallet. It is guaranteed that an existing
// key has always a non-empty slice as value.
func FindPendingRedemptions(
	chain Chain,
	filter PendingRedemptionsFilter,
) (map[[20]byte][]bitcoin.Script, error) {
	logger.Infof(
		"looking for pending redemptions using filter [%s]",
		filter,
	)

	blockCounter, err := chain.BlockCounter()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get block counter: [%w]",
			err,
		)
	}

	currentBlockNumber, err := blockCounter.CurrentBlock()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get current block number: [%w]",
			err,
		)
	}

	redemptionRequestMinAge, err := chain.GetRedemptionRequestMinAge()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get redemption request minimum age: [%w]",
			err,
		)
	}

	_, _, _, _, redemptionTimeout, _, _, err := chain.GetRedemptionParameters()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get redemption parameters: [%w]",
			err,
		)
	}

	getPendingRedemptionsFromWallet := func(
		wallet [20]byte,
	) ([]*RedemptionRequest, error) {
		pendingRedemptions, err := getPendingRedemptions(
			chain,
			wallet,
			currentBlockNumber,
			filter.RedemptionRequestsLimit,
			redemptionTimeout,
			redemptionRequestMinAge,
		)
		if err != nil {
			return nil,
				fmt.Errorf(
					"failed to get pending redemptions for [%s] wallet: [%w]",
					hexutils.Encode(wallet[:]),
					err,
				)
		}
		return pendingRedemptions, nil
	}

	var walletPublicKeyHashes [][20]byte
	if len(filter.WalletPublicKeyHashes) > 0 {
		// Take wallets from the filter.
		walletPublicKeyHashes = filter.WalletPublicKeyHashes
	} else {
		// Take all wallets.
		events, err := chain.PastNewWalletRegisteredEvents(nil)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get new wallet registered wallets: [%w]",
				err,
			)
		}

		// Sort the wallets list from the oldest to the newest.
		sort.SliceStable(events, func(i, j int) bool {
			return events[i].BlockNumber < events[j].BlockNumber
		})

		for _, event := range events {
			walletPublicKeyHashes = append(
				walletPublicKeyHashes,
				event.WalletPublicKeyHash,
			)
		}
	}

	logger.Infof(
		"built an initial list of [%v] wallets that will be checked "+
			"for pending redemption requests",
		len(walletPublicKeyHashes),
	)

	// Apply the wallets number limit if needed.
	if limit := int(filter.WalletsLimit); limit > 0 && len(walletPublicKeyHashes) > limit {
		walletPublicKeyHashes = walletPublicKeyHashes[:limit]

		logger.Infof(
			"limited the initial wallets list to [%v] wallets",
			len(walletPublicKeyHashes),
		)
	}

	result := make(map[[20]byte][]bitcoin.Script)

	for _, walletPublicKeyHash := range walletPublicKeyHashes {
		walletPublicKeyHashHex := hexutils.Encode(walletPublicKeyHash[:])

		logger.Infof(
			"fetching pending redemption requests from wallet [%s]...",
			walletPublicKeyHashHex,
		)

		pendingRedemptions, err := getPendingRedemptionsFromWallet(
			walletPublicKeyHash,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot get pending redemptions for wallet [%s]: [%w]",
				walletPublicKeyHashHex,
				err,
			)
		}

		logger.Infof(
			"found [%d] redemptions for wallet [%s]",
			len(pendingRedemptions),
			walletPublicKeyHashHex,
		)

		for i, pendingRedemption := range pendingRedemptions {
			logger.Infof(
				"redemption [%d/%d] - [%s]",
				i+1,
				len(pendingRedemptions),
				fmt.Sprintf(
					"redemptionKey: [%s], requested at: [%s]",
					pendingRedemption.RedemptionKey,
					pendingRedemption.RequestedAt,
				),
			)

			result[walletPublicKeyHash] = append(
				result[walletPublicKeyHash],
				pendingRedemption.RedeemerOutputScript,
			)
		}
	}

	return result, nil
}

// ProposeRedemption handles redemption proposal submission.
func ProposeRedemption(
	chain Chain,
	btcChain bitcoin.Chain,
	walletPublicKeyHash [20]byte,
	fee int64,
	redeemersOutputScripts []bitcoin.Script,
	dryRun bool,
) error {
	if len(redeemersOutputScripts) == 0 {
		return fmt.Errorf("redemptions list is empty")
	}

	logger.Infof(
		"starting proposing redemption for wallet [%s]...",
		hex.EncodeToString(walletPublicKeyHash[:]),
	)

	// Estimate fee if it's missing. Do not check the estimated fee against
	// the maximum total and per-request fees allowed by the Bridge. This
	// is done during the on-chain validation of the proposal so there is no
	// need to do it here.
	if fee <= 0 {
		logger.Infof("estimating redemption transaction fee...")

		estimatedFee, err := EstimateRedemptionFee(
			btcChain,
			redeemersOutputScripts,
		)
		if err != nil {
			return fmt.Errorf(
				"cannot estimate redemption transaction fee: [%w]",
				err,
			)
		}

		fee = estimatedFee
	}

	logger.Infof("redemption transaction fee: [%d]", fee)

	logger.Infof("preparing a redemption proposal...")

	proposal := &tbtc.RedemptionProposal{
		WalletPublicKeyHash:    walletPublicKeyHash,
		RedeemersOutputScripts: redeemersOutputScripts,
		RedemptionTxFee:        big.NewInt(fee),
	}

	logger.Infof("validating the redemption proposal...")

	if _, err := tbtc.ValidateRedemptionProposal(
		logger,
		proposal,
		chain,
	); err != nil {
		return fmt.Errorf("failed to verify redemption proposal: %v", err)
	}

	if !dryRun {
		logger.Infof("submitting the redemption proposal...")
		if err := chain.SubmitRedemptionProposalWithReimbursement(proposal); err != nil {
			return fmt.Errorf("failed to submit redemption proposal: %v", err)
		}
	} else {
		logger.Infof("skipping transaction submission in dry-run mode")
	}

	return nil
}

func getPendingRedemptions(
	chain Chain,
	walletPublicKeyHash [20]byte,
	currentBlockNumber uint64,
	redemptionRequestsLimit uint16,
	redemptionRequestTimeout uint32,
	redemptionRequestMinAge uint32,
) ([]*RedemptionRequest, error) {
	// We are interested with `RedemptionRequested` events that are not
	// timed out yet. That means there is no sense to look for events that
	// occurred earlier than `now - redemptionRequestTimeout`. However,
	// the event filter expects a block range while `redemptionRequestTimeout`
	// is in seconds. To overcome that problem, we estimate the redemption
	// request timeout in blocks, using the average block time of the host chain.
	// Note that this estimation is not 100% accurate as the actual block time
	// may differ from the assumed one.
	redemptionRequestTimeoutBlocks :=
		uint64(redemptionRequestTimeout) / uint64(chain.AverageBlockTime().Seconds())
	// Then, we set the start block of the filter using the estimated redemption
	// request timeout in blocks. Note that if the actual average block time is
	// lesser than the assumed one, some events being on the edge of the block
	// range may be omitted. To avoid that, we make the block range a little
	// wider by using a constant factor of 1000 blocks.
	filterStartBlock := currentBlockNumber - redemptionRequestTimeoutBlocks - 1000

	filter := &tbtc.RedemptionRequestedEventFilter{
		StartBlock: filterStartBlock,
	}
	if walletPublicKeyHash != [20]byte{} {
		filter.WalletPublicKeyHash = [][20]byte{walletPublicKeyHash}
	}

	events, err := chain.PastRedemptionRequestedEvents(filter)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get past redemption requested events: [%w]",
			err,
		)
	}

	// Take the oldest first.
	sort.SliceStable(
		events, func(i, j int) bool {
			return events[i].BlockNumber < events[j].BlockNumber
		},
	)

	// There may be multiple events targeting the same redemption key
	// (i.e. the same wallet and output script pair). The Bridge contract
	// allows only for one pending request with the given redemption key
	// at the same time. That means we need to deduplicate the events list
	// and take only the latest event for the given redemption key.
	eventsSet := make(map[string]*tbtc.RedemptionRequestedEvent)
	for _, event := range events {
		redemptionKey, err := chain.BuildRedemptionKey(
			event.WalletPublicKeyHash,
			event.RedeemerOutputScript,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to build redemption key: [%v]", err)
		}

		// Events are sorted from the oldest to the newest so there is no
		// need to check for existence. We can just overwrite.
		eventsSet[hexutils.Encode(redemptionKey.Bytes())] = event
	}

	logger.Infof("found [%d] RedemptionRequested events", len(eventsSet))

	logger.Infof("checking pending redemptions details...")

	pendingRedemptions := make([]*RedemptionRequest, 0)

	eventIndex := 0
redemptionRequestedLoop:
	for redemptionKey, event := range eventsSet {
		eventIndex++

		logger.Debugf(
			"getting pending redemption details [%d/%d]",
			eventIndex,
			len(eventsSet),
		)

		// Check if there is still a pending redemption for the given redemption
		// requested event.
		pendingRedemption, err := chain.GetPendingRedemptionRequest(
			event.WalletPublicKeyHash,
			event.RedeemerOutputScript,
		)
		if err != nil {
			switch {
			case errors.Is(err, tbtc.ErrPendingRedemptionRequestNotFound):
				logger.Infof(
					"redemption for request [%d/%d] is no longer pending",
					eventIndex,
					len(eventsSet),
				)

				continue redemptionRequestedLoop
			default:
				return nil, fmt.Errorf(
					"failed to get pending redemption request: [%w]",
					err,
				)
			}
		}

		pendingRedemptions = append(
			pendingRedemptions, &RedemptionRequest{
				WalletPublicKeyHash:  event.WalletPublicKeyHash,
				RedemptionKey:        redemptionKey,
				RedeemerOutputScript: event.RedeemerOutputScript,
				RequestedAt:          pendingRedemption.RequestedAt,
			},
		)
	}

	resultSliceCapacity := len(pendingRedemptions)
	if redemptionRequestsLimit > 0 {
		resultSliceCapacity = int(redemptionRequestsLimit)
	}

	// Sort the pending redemptions from oldest to newest.
	sort.SliceStable(
		pendingRedemptions, func(i, j int) bool {
			return pendingRedemptions[i].RequestedAt.Before(pendingRedemptions[j].RequestedAt)
		},
	)

	// Only redemption requests in range:
	// [now - redemptionRequestTimeout, now - redemptionRequestMinAge]
	// should be taken into consideration.
	redemptionRequestsRangeStartTimestamp := time.Now().Add(
		-time.Duration(redemptionRequestTimeout) * time.Second,
	)
	redemptionRequestsRangeEndTimestamp := time.Now().Add(
		-time.Duration(redemptionRequestMinAge) * time.Second,
	)

	result := make([]*RedemptionRequest, 0, resultSliceCapacity)
	for i, pendingRedemption := range pendingRedemptions {
		if len(result) == cap(result) {
			break
		}

		// Check if timeout passed for the redemption request.
		if pendingRedemption.RequestedAt.Before(redemptionRequestsRangeStartTimestamp) {
			logger.Infof(
				"redemption request [%d/%d] has already timed out",
				i+1,
				len(pendingRedemptions),
			)
			continue
		}

		// Check if enough time elapsed since the redemption request.
		if pendingRedemption.RequestedAt.After(redemptionRequestsRangeEndTimestamp) {
			logger.Infof(
				"redemption request [%d/%d] is not old enough",
				i+1,
				len(pendingRedemptions),
			)
			continue
		}

		result = append(result, pendingRedemption)
	}

	return result, nil
}

// EstimateRedemptionFee estimates fee for the redemption transaction that pays
// the provided redeemers output scripts.
func EstimateRedemptionFee(
	btcChain bitcoin.Chain,
	redeemersOutputScripts []bitcoin.Script,
) (int64, error) {
	sizeEstimator := bitcoin.NewTransactionSizeEstimator().
		// 1 P2WPKH main UTXO input.
		AddPublicKeyHashInputs(1, true).
		// 1 P2WPKH change output.
		AddPublicKeyHashOutputs(1, true)

	for _, script := range redeemersOutputScripts {
		switch bitcoin.GetScriptType(script) {
		case bitcoin.P2PKHScript:
			sizeEstimator.AddPublicKeyHashOutputs(1, false)
		case bitcoin.P2WPKHScript:
			sizeEstimator.AddPublicKeyHashOutputs(1, true)
		case bitcoin.P2SHScript:
			sizeEstimator.AddScriptHashOutputs(1, false)
		case bitcoin.P2WSHScript:
			sizeEstimator.AddScriptHashOutputs(1, true)
		default:
			return 0, fmt.Errorf("non-standard redeemer output script type")
		}
	}

	transactionSize, err := sizeEstimator.VirtualSize()
	if err != nil {
		return 0, fmt.Errorf("cannot estimate transaction virtual size: [%v]", err)
	}

	feeEstimator := bitcoin.NewTransactionFeeEstimator(btcChain)

	totalFee, err := feeEstimator.EstimateFee(transactionSize)
	if err != nil {
		return 0, fmt.Errorf("cannot estimate transaction fee: [%v]", err)
	}

	return totalFee, nil
}
