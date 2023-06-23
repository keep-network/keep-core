package coordinator

import (
	"errors"
	"fmt"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"math/big"
	"sort"
	"time"

	"github.com/keep-network/keep-core/internal/hexutils"
)

type redemptionEntry struct {
	walletPublicKeyHash [20]byte

	redemptionKey        string
	redeemerOutputScript bitcoin.Script
	requestedAt          time.Time
}

// FindPendingRedemptions finds pending redemptions requests.
func FindPendingRedemptions(
	chain Chain,
	walletPublicKeyHash [20]byte,
	maxNumberOfDeposits uint16,
) ([20]byte, []bitcoin.Script, error) {
	logger.Infof("redemption max size: %d", maxNumberOfDeposits)

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

	getPendingRedemptionsFromWallet := func(walletToSweep [20]byte) ([]*redemptionEntry, error) {
		pendingRedemptions, err := getPendingRedemptions(
			chain,
			walletToSweep,
			int(maxNumberOfDeposits),
			redemptionTimeout,
			redemptionRequestMinAge,
		)
		if err != nil {
			return nil,
				fmt.Errorf(
					"failed to get deposits for [%s] wallet: [%w]",
					hexutils.Encode(walletToSweep[:]),
					err,
				)
		}
		return pendingRedemptions, nil
	}

	var redemptionsToPropose []*redemptionEntry
	// If walletPublicKeyHash is not provided we need to find a wallet that has
	// pending redemptions.
	if walletPublicKeyHash == [20]byte{} {
		walletRegisteredEvents, err := chain.PastNewWalletRegisteredEvents(nil)
		if err != nil {
			return [20]byte{}, nil, fmt.Errorf("failed to get registered wallets: [%w]", err)
		}

		// Take the oldest first
		sort.SliceStable(walletRegisteredEvents, func(i, j int) bool {
			return walletRegisteredEvents[i].BlockNumber < walletRegisteredEvents[j].BlockNumber
		})

		redeemingWallets := walletRegisteredEvents
		// Only two the most recently created wallets are accepting redemptions.
		if len(walletRegisteredEvents) >= 2 {
			redeemingWallets = walletRegisteredEvents[len(walletRegisteredEvents)-2:]
		}

		for _, registeredWallet := range redeemingWallets {
			logger.Infof(
				"fetching pending redemption requests from wallet [%s]...",
				hexutils.Encode(registeredWallet.WalletPublicKeyHash[:]),
			)

			pendingRedemptions, err := getPendingRedemptionsFromWallet(
				registeredWallet.WalletPublicKeyHash,
			)
			if err != nil {
				return [20]byte{}, nil, err
			}

			// Check if there are any unswept deposits in this wallet. If so
			// sweep this wallet and don't check the other wallet.
			if len(pendingRedemptions) > 0 {
				walletPublicKeyHash = registeredWallet.WalletPublicKeyHash
				redemptionsToPropose = pendingRedemptions
				break
			}
		}
	} else {
		logger.Infof(
			"fetching deposits from wallet [%s]...",
			hexutils.Encode(walletPublicKeyHash[:]),
		)
		unsweptDeposits, err := getPendingRedemptionsFromWallet(
			walletPublicKeyHash,
		)
		if err != nil {
			return [20]byte{}, nil, err
		}
		redemptionsToPropose = unsweptDeposits
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
				redemption.redemptionKey,
				redemption.requestedAt,
			))

		redeemersOutputScripts[i] = redemption.redeemerOutputScript
	}

	return walletPublicKeyHash, redeemersOutputScripts, nil
}

// ProposeRedemption handles redemption proposal submission.
func ProposeRedemption(
	chain Chain,
	walletPublicKeyHash [20]byte,
	fee int64,
	redeemersOutputScripts []bitcoin.Script,
	dryRun bool,
) error {
	if len(redeemersOutputScripts) == 0 {
		return fmt.Errorf("redemptions list is empty")
	}

	// Estimate fee if it's missing.
	if fee <= 0 {
		panic("fee estimation not implemented yet")
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
		logger.Infof("submitting the deposit sweep proposal...")
		if err := chain.SubmitRedemptionProposalWithReimbursement(proposal); err != nil {
			return fmt.Errorf("failed to submit deposit sweep proposal: %v", err)
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
) ([]*redemptionEntry, error) {
	logger.Infof("reading revealed deposits from chain...")

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

	redemptionsRequested, err := chain.PastRedemptionRequestedEvents(filter)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get past redemption requested events: [%w]",
			err,
		)
	}
	logger.Infof("found %d RedemptionRequested events", len(redemptionsRequested))

	// Take the oldest first.
	sort.SliceStable(redemptionsRequested, func(i, j int) bool {
		return redemptionsRequested[i].BlockNumber < redemptionsRequested[j].BlockNumber
	})

	logger.Infof("checking pending redemptions details...")

	pendingRedemptions := make([]*redemptionEntry, 0)

redemptionRequestedLoop:
	for i, event := range redemptionsRequested {
		logger.Debugf("getting pending redemption details %d/%d", i+1, len(redemptionsRequested))

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
					i+1,
					len(redemptionsRequested),
				)

				continue redemptionRequestedLoop
			default:
				return nil, fmt.Errorf(
					"failed to get pending redemption request: [%w]",
					err,
				)
			}
		}

		redemptionKey, err := chain.BuildRedemptionKey(
			event.WalletPublicKeyHash,
			event.RedeemerOutputScript,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to build redemption key: [%v]", err)
		}

		pendingRedemptions = append(pendingRedemptions, &redemptionEntry{
			walletPublicKeyHash:  event.WalletPublicKeyHash,
			redemptionKey:        hexutils.Encode(redemptionKey.Bytes()),
			redeemerOutputScript: event.RedeemerOutputScript,
			requestedAt:          pendingRedemption.RequestedAt,
		})
	}

	resultSliceCapacity := len(redemptionsRequested)
	if maxNumberOfRedemptions > 0 {
		resultSliceCapacity = maxNumberOfRedemptions
	}

	// Sort the pending redemptions.
	sort.SliceStable(pendingRedemptions, func(i, j int) bool {
		return pendingRedemptions[i].requestedAt.Before(pendingRedemptions[j].requestedAt)
	})

	result := make([]*redemptionEntry, 0, resultSliceCapacity)
	for i, pendingRedemption := range pendingRedemptions {
		if len(result) == cap(result) {
			break
		}

		// Check if timeout passed for the redemption request.
		if pendingRedemption.requestedAt.Before(redemptionRequestsRangeStartTimestamp) {
			logger.Infof(
				"redemption request %d/%d has already timed out",
				i+1,
				len(redemptionsRequested),
			)
			continue
		}

		// Check if enough time elapsed since the redemption request.
		if pendingRedemption.requestedAt.After(redemptionRequestsRangeEndTimestamp) {
			logger.Infof(
				"redemption request %d/%d is not old enough",
				i+1,
				len(redemptionsRequested),
			)
			continue
		}

		result = append(result, pendingRedemption)
	}

	return result, nil
}
