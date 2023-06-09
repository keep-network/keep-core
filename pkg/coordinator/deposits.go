package coordinator

import (
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/tbtc"
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

// ListDeposits gets deposits from the chain.
func ListDeposits(
	tbtcChain tbtc.Chain,
	btcChain bitcoin.Chain,
	walletPublicKeyHashString string,
	hideSwept bool,
	sortByAmount bool,
	head int,
	tail int,
) error {
	deposits, err := getDeposits(
		tbtcChain,
		btcChain,
		walletPublicKeyHashString,
		sortByAmount,
		head,
		tail,
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

	// Filter
	if hideSwept {
		deposits = removeSwept(deposits)
	}

	// Print
	if err := printTable(deposits); err != nil {
		return fmt.Errorf("failed to print deposits table: %v", err)
	}

	return nil
}

func removeSwept(deposits []depositEntry) []depositEntry {
	result := []depositEntry{}
	for _, deposit := range deposits {
		if deposit.isSwept {
			continue
		}
		result = append(result, deposit)
	}
	return result
}

func getDeposits(
	tbtcChain tbtc.Chain,
	btcChain bitcoin.Chain,
	walletPublicKeyHashString string,
	sortByAmount bool,
	head int,
	tail int,
) ([]depositEntry, error) {
	logger.Infof("reading deposits from chain...")

	filter := &tbtc.DepositRevealedEventFilter{}
	if len(walletPublicKeyHashString) > 0 {
		walletPublicKeyHash, err := hexToWalletPublicKeyHash(walletPublicKeyHashString)
		if err != nil {
			return []depositEntry{}, fmt.Errorf("failed to extract wallet public key hash: %v", err)
		}

		filter.WalletPublicKeyHash = [][20]byte{walletPublicKeyHash}
	}

	allDepositRevealedEvents, err := tbtcChain.PastDepositRevealedEvents(filter)
	if err != nil {
		return []depositEntry{}, fmt.Errorf(
			"failed to get past deposit revealed events: [%w]",
			err,
		)
	}

	logger.Infof("found %d DepositRevealed events", len(allDepositRevealedEvents))

	// Order
	sort.SliceStable(allDepositRevealedEvents, func(i, j int) bool {
		return allDepositRevealedEvents[i].BlockNumber > allDepositRevealedEvents[j].BlockNumber
	})

	if sortByAmount {
		sort.SliceStable(allDepositRevealedEvents, func(i, j int) bool {
			return allDepositRevealedEvents[i].Amount < allDepositRevealedEvents[j].Amount
		})
	}

	// Filter
	depositRevealedEvents := []*tbtc.DepositRevealedEvent{}

	if len(allDepositRevealedEvents) > head+tail && (head > 0 || tail > 0) {
		// Head
		depositRevealedEvents = append(
			depositRevealedEvents,
			allDepositRevealedEvents[:head]...,
		)
		// Tail
		depositRevealedEvents = append(
			depositRevealedEvents,
			allDepositRevealedEvents[len(allDepositRevealedEvents)-tail:]...,
		)
	} else {
		depositRevealedEvents = append(
			depositRevealedEvents,
			allDepositRevealedEvents...,
		)
	}

	result := make([]depositEntry, len(depositRevealedEvents))
	for i, event := range depositRevealedEvents {
		logger.Debugf("getting details of deposit %d/%d", i+1, len(depositRevealedEvents))

		depositKey := tbtcChain.BuildDepositKey(event.FundingTxHash, event.FundingOutputIndex)

		depositRequest, err := tbtcChain.GetDepositRequest(event.FundingTxHash, event.FundingOutputIndex)
		if err != nil {
			return result, fmt.Errorf(
				"failed to get deposit request: [%w]",
				err,
			)
		}

		confirmations, err := btcChain.GetTransactionConfirmations(event.FundingTxHash)
		if err != nil {
			logger.Errorf(
				"failed to get bitcoin transaction confirmations: [%v]",
				err,
			)
		}

		result[i] = depositEntry{
			walletPublicKeyHash:           event.WalletPublicKeyHash,
			depositKey:                    hexutils.Encode(depositKey.Bytes()),
			revealBlock:                   event.BlockNumber,
			isSwept:                       depositRequest.SweptAt.Unix() != 0,
			fundingTransactionHash:        event.FundingTxHash,
			fundingTransactionOutputIndex: event.FundingOutputIndex,
			amountBtc:                     convertSatToBtc(float64(depositRequest.Amount)),
			confirmations:                 confirmations,
		}
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
