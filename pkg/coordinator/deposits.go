package coordinator

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

var (
	// ErrNoDepositsToSweep throws an error if no deposits that could be swept
	// have been found.
	ErrNoDepositsToSweep = errors.New("no deposits to sweep")
)

type depositEntry struct {
	walletPublicKeyHash WalletPublicKeyHash

	depositKey  string
	revealBlock uint64
	isSwept     bool

	fundingTransactionHash        bitcoin.Hash
	fundingTransactionOutputIndex uint32
	amountBtc                     float64
	confirmations                 uint
}

// ListDeposits gets deposits from the chain and prints them to standard output.
func ListDeposits(
	tbtcChain tbtc.Chain,
	btcChain bitcoin.Chain,
	walletPublicKeyHash *WalletPublicKeyHash,
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

func getDeposits(
	tbtcChain tbtc.Chain,
	btcChain bitcoin.Chain,
	walletPublicKeyHash *WalletPublicKeyHash,
	maxNumberOfDeposits int,
	skipSwept bool,
	skipUnconfirmed bool,
) ([]depositEntry, error) {
	logger.Infof("reading revealed deposits from chain...")

	filter := &tbtc.DepositRevealedEventFilter{}
	if walletPublicKeyHash != nil {
		filter.WalletPublicKeyHash = [][20]byte{*walletPublicKeyHash}
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
				walletPublicKeyHash:           event.WalletPublicKeyHash,
				depositKey:                    hexutils.Encode(depositKey.Bytes()),
				revealBlock:                   event.BlockNumber,
				isSwept:                       isSwept,
				fundingTransactionHash:        event.FundingTxHash,
				fundingTransactionOutputIndex: event.FundingOutputIndex,
				amountBtc:                     convertSatToBtc(float64(depositRequest.Amount)),
				confirmations:                 confirmations,
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
			deposit.walletPublicKeyHash,
			deposit.amountBtc,
			deposit.depositKey,
			fmt.Sprintf(
				"%s:%d:%d",
				deposit.fundingTransactionHash.Hex(bitcoin.ReversedByteOrder),
				deposit.fundingTransactionOutputIndex,
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
