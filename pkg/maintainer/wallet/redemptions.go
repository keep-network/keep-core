package wallet

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
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

	walletPublicKeyHash, redeemersOutputScripts, err := FindPendingRedemptions(
		wm.chain,
		[20]byte{},
		redemptionMaxSize,
	)
	if err != nil {
		return fmt.Errorf("failed to prepare redemption proposal: [%w]", err)
	}

	if len(redeemersOutputScripts) == 0 {
		logger.Info("no pending redemption requests")
		return nil
	}

	return wm.runIfWalletUnlocked(
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
}

// RedemptionRequest represents a redemption request.
type RedemptionRequest struct {
	WalletPublicKeyHash  [20]byte
	RedemptionKey        string
	RedeemerOutputScript bitcoin.Script
	RequestedAt          time.Time
}

// FindPendingRedemptions finds pending redemptions requests.
func FindPendingRedemptions(
	chain Chain,
	walletPublicKeyHash [20]byte,
	maxNumberOfRedemptions uint16,
) ([20]byte, []bitcoin.Script, error) {
	logger.Infof("redemption max size: %d", maxNumberOfRedemptions)

	redemptionRequestMinAge, err := chain.GetRedemptionRequestMinAge()
	if err != nil {
		return [20]byte{}, nil, fmt.Errorf(
			"failed to get redemption request minimum age: [%w]",
			err,
		)
	}

	_, _, _, _, redemptionTimeout, _, _, err := chain.GetRedemptionParameters()
	if err != nil {
		return [20]byte{}, nil, fmt.Errorf(
			"failed to get redemption parameters: [%w]",
			err,
		)
	}

	getPendingRedemptionsFromWallet := func(wallet [20]byte) ([]*RedemptionRequest, error) {
		pendingRedemptions, err := getPendingRedemptions(
			chain,
			wallet,
			int(maxNumberOfRedemptions),
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

	var redemptionsToPropose []*RedemptionRequest
	// If walletPublicKeyHash is not provided we need to find a wallet that has
	// pending redemptions.
	if walletPublicKeyHash == [20]byte{} {
		events, err := chain.PastNewWalletRegisteredEvents(nil)
		if err != nil {
			return [20]byte{}, nil, fmt.Errorf("failed to get registered wallets: [%w]", err)
		}

		// Take the oldest first
		sort.SliceStable(events, func(i, j int) bool {
			return events[i].BlockNumber < events[j].BlockNumber
		})

		for _, event := range events {
			logger.Infof(
				"fetching pending redemption requests from wallet [%s]...",
				hexutils.Encode(event.WalletPublicKeyHash[:]),
			)

			pendingRedemptions, err := getPendingRedemptionsFromWallet(
				event.WalletPublicKeyHash,
			)
			if err != nil {
				return [20]byte{}, nil, err
			}

			if len(pendingRedemptions) > 0 {
				walletPublicKeyHash = event.WalletPublicKeyHash
				redemptionsToPropose = pendingRedemptions
				break
			}
		}
	} else {
		logger.Infof(
			"fetching pending redemptions from wallet [%s]...",
			hexutils.Encode(walletPublicKeyHash[:]),
		)
		redemptions, err := getPendingRedemptionsFromWallet(
			walletPublicKeyHash,
		)
		if err != nil {
			return [20]byte{}, nil, err
		}
		redemptionsToPropose = redemptions
	}

	if len(redemptionsToPropose) == 0 {
		return [20]byte{}, nil, nil
	}

	logger.Infof(
		"found [%d] redemptions for wallet [%s]",
		len(redemptionsToPropose),
		hexutils.Encode(walletPublicKeyHash[:]),
	)

	redeemersOutputScripts := make([]bitcoin.Script, len(redemptionsToPropose))

	for i, redemption := range redemptionsToPropose {
		logger.Infof(
			"redemption [%d/%d] - %s",
			i+1,
			len(redemptionsToPropose),
			fmt.Sprintf(
				"redemptionKey: [%s], requested at: [%s]",
				redemption.RedemptionKey,
				redemption.RequestedAt,
			))

		redeemersOutputScripts[i] = redemption.RedeemerOutputScript
	}

	return walletPublicKeyHash, redeemersOutputScripts, nil
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
	maxNumberOfRedemptions int,
	redemptionRequestTimeout uint32,
	redemptionRequestMinAge uint32,
) ([]*RedemptionRequest, error) {
	logger.Infof("reading pending redemptions from chain...")

	filter := &tbtc.RedemptionRequestedEventFilter{}
	if walletPublicKeyHash != [20]byte{} {
		filter.WalletPublicKeyHash = [][20]byte{walletPublicKeyHash}
	}

	// Only redemption requests in range:
	// [now - redemptionRequestTimeout, now - redemptionRequestMinAge]
	// should be taken into consideration.
	redemptionRequestsRangeStartTimestamp := time.Now().Add(
		time.Duration(-redemptionRequestTimeout) * time.Second,
	)
	redemptionRequestsRangeEndTimestamp := time.Now().Add(
		time.Duration(-redemptionRequestMinAge) * time.Second,
	)

	// TODO: We should consider narrowing the block range in filter the fetch
	//       only events that are newer than `now - redemptionRequestTimeout`.

	events, err := chain.PastRedemptionRequestedEvents(filter)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get past redemption requested events: [%w]",
			err,
		)
	}

	// Take the oldest first.
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].BlockNumber < events[j].BlockNumber
	})

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

	logger.Infof("found %d redemption requests", len(eventsSet))

	logger.Infof("checking pending redemptions details...")

	pendingRedemptions := make([]*RedemptionRequest, 0)

	eventIndex := 0
redemptionRequestedLoop:
	for redemptionKey, event := range eventsSet {
		eventIndex++

		logger.Debugf(
			"getting pending redemption details %d/%d",
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
				logger.Debugf(
					"redemption for request %d/%d is no longer pending",
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

		pendingRedemptions = append(pendingRedemptions, &RedemptionRequest{
			WalletPublicKeyHash:  event.WalletPublicKeyHash,
			RedemptionKey:        redemptionKey,
			RedeemerOutputScript: event.RedeemerOutputScript,
			RequestedAt:          pendingRedemption.RequestedAt,
		})
	}

	logger.Infof("found %d redemption requests", len(pendingRedemptions))

	resultSliceCapacity := len(pendingRedemptions)
	if maxNumberOfRedemptions > 0 {
		resultSliceCapacity = maxNumberOfRedemptions
	}

	// Sort the pending redemptions.
	sort.SliceStable(pendingRedemptions, func(i, j int) bool {
		return pendingRedemptions[i].RequestedAt.Before(pendingRedemptions[j].RequestedAt)
	})

	result := make([]*RedemptionRequest, 0, resultSliceCapacity)
	for i, pendingRedemption := range pendingRedemptions {
		if len(result) == cap(result) {
			break
		}

		// Check if timeout passed for the redemption request.
		if pendingRedemption.RequestedAt.Before(redemptionRequestsRangeStartTimestamp) {
			logger.Infof(
				"redemption request %d/%d has already timed out",
				i+1,
				len(pendingRedemptions),
			)
			continue
		}

		// Check if enough time elapsed since the redemption request.
		if pendingRedemption.RequestedAt.After(redemptionRequestsRangeEndTimestamp) {
			logger.Infof(
				"redemption request %d/%d is not old enough",
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
