package coordinator

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

type depositEntry struct {
	walletPublicKeyHash [20]byte

	depositKey  string
	revealBlock uint64
	isSwept     bool

	fundingTxHash      bitcoin.Hash
	fundingOutputIndex uint32
	amountBtc          float64
	confirmations      uint
}

// ListDeposits gets deposits from the chain and prints them to standard output.
func ListDeposits(
	tbtcChain tbtc.Chain,
	btcChain bitcoin.Chain,
	walletPublicKeyHash [20]byte,
	head int,
	skipSwept bool,
) error {
	deposits, err := getDeposits(
		tbtcChain,
		btcChain,
		walletPublicKeyHash,
		head,
		skipSwept,
		false,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to get deposits: [%w]",
			err,
		)
	}

	if len(deposits) == 0 {
		return fmt.Errorf("no deposits found")
	}

	// Print
	if err := printTable(deposits); err != nil {
		return fmt.Errorf("failed to print deposits table: %v", err)
	}

	return nil
}

// FindDepositsToSweep finds deposits that can be swept.
// If a wallet public key hash is provided, it will find unswept deposits for the
// given wallet. If a wallet public key hash is nil, it will check all wallets
// starting from the oldest one to find a first wallet containing unswept deposits
// and return those deposits.
// maxNumberOfDeposits is used as a ceiling for the number of deposits in the
// result. If number of discovered deposits meets the maxNumberOfDeposits the
// function will stop fetching more deposits.
// This function will return a wallet public key hash and a list of deposits from
// the wallet that can be swept.
// Deposits with insufficient number of funding transaction confirmations will
// not be taken into consideration for sweeping.
// The result will not mix deposits for different wallets.
// TODO: Cache immutable data
func FindDepositsToSweep(
	tbtcChain tbtc.Chain,
	btcChain bitcoin.Chain,
	walletPublicKeyHash [20]byte,
	maxNumberOfDeposits uint16,
) ([20]byte, []*DepositSweepDetails, error) {
	logger.Infof("deposit sweep max size: %d", maxNumberOfDeposits)

	getDepositsToSweepFromWallet := func(walletToSweep [20]byte) ([]depositEntry, error) {
		unsweptDeposits, err := getDeposits(
			tbtcChain,
			btcChain,
			walletToSweep,
			int(maxNumberOfDeposits),
			true,
			true,
		)
		if err != nil {
			return nil,
				fmt.Errorf(
					"failed to get deposits for [%s] wallet: [%w]",
					walletToSweep,
					err,
				)
		}
		return unsweptDeposits, nil
	}

	depositsToSweep := make([]depositEntry, 0, maxNumberOfDeposits)
	// If walletPublicKeyHash is not provided we need to find a wallet that has
	// unswept deposits.
	if walletPublicKeyHash == [20]byte{} {
		walletRegisteredEvents, err := tbtcChain.PastNewWalletRegisteredEvents(nil)
		if err != nil {
			return [20]byte{}, nil, fmt.Errorf("failed to get registered wallets: [%w]", err)
		}

		// Take the oldest first
		sort.SliceStable(walletRegisteredEvents, func(i, j int) bool {
			return walletRegisteredEvents[i].BlockNumber < walletRegisteredEvents[j].BlockNumber
		})

		sweepingWallets := walletRegisteredEvents
		// Only two the most recently created wallets are sweeping.
		if len(walletRegisteredEvents) >= 2 {
			sweepingWallets = walletRegisteredEvents[len(walletRegisteredEvents)-2:]
		}

		for _, registeredWallet := range sweepingWallets {
			logger.Infof(
				"fetching deposits from wallet [%s]...",
				hexutils.Encode(registeredWallet.WalletPublicKeyHash[:]),
			)

			unsweptDeposits, err := getDepositsToSweepFromWallet(
				registeredWallet.WalletPublicKeyHash,
			)
			if err != nil {
				return [20]byte{}, nil, err
			}

			// Check if there are any unswept deposits in this wallet. If so
			// sweep this wallet and don't check the other wallet.
			if len(unsweptDeposits) > 0 {
				walletPublicKeyHash = registeredWallet.WalletPublicKeyHash
				depositsToSweep = unsweptDeposits
				break
			}
		}
	} else {
		logger.Infof(
			"fetching deposits from wallet [%s]...",
			hexutils.Encode(walletPublicKeyHash[:]),
		)
		unsweptDeposits, err := getDepositsToSweepFromWallet(
			walletPublicKeyHash,
		)
		if err != nil {
			return [20]byte{}, nil, err
		}
		depositsToSweep = unsweptDeposits
	}

	if len(depositsToSweep) == 0 {
		return [20]byte{}, nil, nil
	}

	logger.Infof(
		"found [%d] deposits to sweep for wallet [%s]",
		len(depositsToSweep),
		hexutils.Encode(walletPublicKeyHash[:]),
	)

	for i, deposit := range depositsToSweep {
		logger.Debugf(
			"deposit [%d/%d] - %s",
			i+1,
			len(depositsToSweep),
			fmt.Sprintf(
				"depositKey: [%s], reveal block: [%d], funding transaction: [%s], output index: [%d]",
				deposit.depositKey,
				deposit.revealBlock,
				deposit.fundingTxHash.Hex(bitcoin.ReversedByteOrder),
				deposit.fundingOutputIndex,
			))
	}

	result := make([]*DepositSweepDetails, len(depositsToSweep))
	for i, deposit := range depositsToSweep {
		result[i] = &DepositSweepDetails{
			FundingTxHash:      deposit.fundingTxHash,
			FundingOutputIndex: deposit.fundingOutputIndex,
			RevealBlock:        deposit.revealBlock,
		}
	}

	return walletPublicKeyHash, result, nil
}

func getDeposits(
	tbtcChain tbtc.Chain,
	btcChain bitcoin.Chain,
	walletPublicKeyHash [20]byte,
	maxNumberOfDeposits int,
	skipSwept bool,
	skipUnconfirmed bool,
) ([]depositEntry, error) {
	logger.Infof("reading revealed deposits from chain...")

	filter := &tbtc.DepositRevealedEventFilter{}
	if walletPublicKeyHash != [20]byte{} {
		filter.WalletPublicKeyHash = [][20]byte{walletPublicKeyHash}
	}

	depositRevealedEvents, err := tbtcChain.PastDepositRevealedEvents(filter)
	if err != nil {
		return []depositEntry{}, fmt.Errorf(
			"failed to get past deposit revealed events: [%w]",
			err,
		)
	}

	logger.Infof("found %d DepositRevealed events", len(depositRevealedEvents))

	// Take the oldest first
	sort.SliceStable(depositRevealedEvents, func(i, j int) bool {
		return depositRevealedEvents[i].BlockNumber < depositRevealedEvents[j].BlockNumber
	})

	logger.Infof("getting deposits details...")

	resultSliceCapacity := len(depositRevealedEvents)
	if maxNumberOfDeposits > 0 {
		resultSliceCapacity = maxNumberOfDeposits
	}

	result := make([]depositEntry, 0, resultSliceCapacity)
	for i, event := range depositRevealedEvents {
		if len(result) == cap(result) {
			break
		}

		logger.Debugf("getting details of deposit %d/%d", i+1, len(depositRevealedEvents))

		depositKey := tbtcChain.BuildDepositKey(event.FundingTxHash, event.FundingOutputIndex)

		depositRequest, err := tbtcChain.GetDepositRequest(event.FundingTxHash, event.FundingOutputIndex)
		if err != nil {
			return result, fmt.Errorf(
				"failed to get deposit request: [%w]",
				err,
			)
		}

		isSwept := depositRequest.SweptAt.Unix() != 0
		if skipSwept && isSwept {
			logger.Debugf("deposit %d/%d is already swept", i+1, len(depositRevealedEvents))
			continue
		}

		confirmations, err := btcChain.GetTransactionConfirmations(event.FundingTxHash)
		if err != nil {
			logger.Errorf(
				"failed to get bitcoin transaction confirmations: [%v]",
				err,
			)
		}

		if skipUnconfirmed && confirmations < tbtc.DepositSweepRequiredFundingTxConfirmations {
			logger.Debugf(
				"deposit %d/%d funding transaction doesn't have enough confirmations: %d/%d",
				i+1, len(depositRevealedEvents),
				confirmations, tbtc.DepositSweepRequiredFundingTxConfirmations)
			continue
		}

		result = append(
			result,
			depositEntry{
				walletPublicKeyHash: event.WalletPublicKeyHash,
				depositKey:          hexutils.Encode(depositKey.Bytes()),
				revealBlock:         event.BlockNumber,
				isSwept:             isSwept,
				fundingTxHash:       event.FundingTxHash,
				fundingOutputIndex:  event.FundingOutputIndex,
				amountBtc:           convertSatToBtc(float64(depositRequest.Amount)),
				confirmations:       confirmations,
			},
		)
	}

	return result, nil
}

func printTable(deposits []depositEntry) error {
	w := tabwriter.NewWriter(os.Stdout, 2, 4, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "index\twallet\tvalue (BTC)\tdeposit key\trevealed deposit data\tconfirmations\tswept\t\n")

	for i, deposit := range deposits {
		fmt.Fprintf(w, "%d\t%s\t%.5f\t%s\t%s\t%d\t%t\t\n",
			i,
			hexutils.Encode(deposit.walletPublicKeyHash[:]),
			deposit.amountBtc,
			deposit.depositKey,
			fmt.Sprintf(
				"%s:%d:%d",
				deposit.fundingTxHash.Hex(bitcoin.ReversedByteOrder),
				deposit.fundingOutputIndex,
				deposit.revealBlock,
			),
			deposit.confirmations,
			deposit.isSwept,
		)
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush the writer: %v", err)
	}

	return nil
}

func convertSatToBtc(sats float64) float64 {
	return sats / float64(100000000)
}
