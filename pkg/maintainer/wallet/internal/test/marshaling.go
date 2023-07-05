package test

import (
	"encoding/json"
	"fmt"
	walletmtr "github.com/keep-network/keep-core/pkg/maintainer/wallet"
	"math/big"
	"time"

	"github.com/keep-network/keep-core/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

// UnmarshalJSON implements a custom JSON unmarshaling logic to produce a
// proper FindDepositsToSweepTestScenario.
func (dsts *FindDepositsToSweepTestScenario) UnmarshalJSON(data []byte) error {
	type findDepositsToSweepTestScenario struct {
		Title               string
		WalletPublicKeyHash string
		MaxNumberOfDeposits uint16
		Wallets             []struct {
			WalletPublicKeyHash     string
			RegistrationBlockNumber uint64
		}
		Deposits []struct {
			FundingTxHash          string
			FundingOutputIndex     uint32
			FundingTxConfirmations uint
			WalletPublicKeyHash    string
			RevealBlockNumber      uint64
			SweptAt                int64
		}
		ExpectedWalletPublicKeyHash string
		ExpectedUnsweptDeposits     []struct {
			FundingTxHash      string
			FundingOutputIndex uint32
			RevealBlockNumber  uint64
		}
	}

	var unmarshaled findDepositsToSweepTestScenario

	err := json.Unmarshal(data, &unmarshaled)
	if err != nil {
		return err
	}

	// Unmarshal title.
	dsts.Title = unmarshaled.Title

	// Unmarshal max number of deposits.
	if len(unmarshaled.WalletPublicKeyHash) > 0 {
		copy(dsts.WalletPublicKeyHash[:], hexToSlice(unmarshaled.WalletPublicKeyHash))
	}

	dsts.MaxNumberOfDeposits = unmarshaled.MaxNumberOfDeposits

	// Unmarshal wallets.
	for _, uw := range unmarshaled.Wallets {
		w := new(Wallet)

		copy(w.WalletPublicKeyHash[:], hexToSlice(uw.WalletPublicKeyHash))
		w.RegistrationBlockNumber = uw.RegistrationBlockNumber

		dsts.Wallets = append(dsts.Wallets, w)
	}
	// Unmarshal deposits.
	for i, deposit := range unmarshaled.Deposits {
		d := new(Deposit)

		fundingTxHash, err := bitcoin.NewHashFromString(deposit.FundingTxHash, bitcoin.ReversedByteOrder)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal funding transaction hash for deposit [%d/%d]: [%w]",
				i,
				len(unmarshaled.Deposits),
				err,
			)
		}

		copy(d.WalletPublicKeyHash[:], hexToSlice(deposit.WalletPublicKeyHash))

		d.FundingTxHash = fundingTxHash
		d.FundingOutputIndex = deposit.FundingOutputIndex
		d.FundingTxConfirmations = deposit.FundingTxConfirmations
		d.RevealBlockNumber = deposit.RevealBlockNumber
		d.SweptAt = time.Unix(deposit.SweptAt, 0)

		dsts.Deposits = append(dsts.Deposits, d)
	}

	// Unmarshal expected wallet public key hash.
	if len(unmarshaled.ExpectedWalletPublicKeyHash) > 0 {
		copy(dsts.ExpectedWalletPublicKeyHash[:], hexToSlice(unmarshaled.ExpectedWalletPublicKeyHash))
	}

	// Unmarshal expected unswept deposits.
	for i, deposit := range unmarshaled.ExpectedUnsweptDeposits {
		ud := new(walletmtr.DepositReference)

		fundingTxHash, err := bitcoin.NewHashFromString(deposit.FundingTxHash, bitcoin.ReversedByteOrder)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal funding transaction hash for expected unswept deposit [%d/%d]: [%w]",
				i,
				len(unmarshaled.Deposits),
				err,
			)
		}

		ud.FundingTxHash = fundingTxHash
		ud.FundingOutputIndex = deposit.FundingOutputIndex
		ud.RevealBlock = deposit.RevealBlockNumber

		dsts.ExpectedUnsweptDeposits = append(dsts.ExpectedUnsweptDeposits, ud)
	}

	return nil
}

type depositSweepProposal struct {
	WalletPublicKeyHash string
	DepositsKeys        []struct {
		FundingTxHash      string
		FundingOutputIndex uint32
	}
	SweepTxFee           int64
	DepositsRevealBlocks []int64
}

func (dsp *depositSweepProposal) convert() (*tbtc.DepositSweepProposal, error) {
	if dsp == nil {
		return nil, nil
	}

	result := &tbtc.DepositSweepProposal{}

	if len(dsp.WalletPublicKeyHash) > 0 {
		copy(result.WalletPublicKeyHash[:], hexToSlice(dsp.WalletPublicKeyHash))
	}

	result.DepositsKeys = make([]struct {
		FundingTxHash      bitcoin.Hash
		FundingOutputIndex uint32
	}, len(dsp.DepositsKeys))
	for i, depositKey := range dsp.DepositsKeys {
		fundingTxHash, err := bitcoin.NewHashFromString(depositKey.FundingTxHash, bitcoin.ReversedByteOrder)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to unmarshal funding transaction hash for deposit [%d/%d]: [%w]",
				i,
				len(dsp.DepositsKeys),
				err,
			)
		}
		result.DepositsKeys[i].FundingTxHash = fundingTxHash
		result.DepositsKeys[i].FundingOutputIndex = depositKey.FundingOutputIndex
	}

	result.DepositsRevealBlocks = make([]*big.Int, len(dsp.DepositsRevealBlocks))
	for i, depositRevealBlock := range dsp.DepositsRevealBlocks {
		result.DepositsRevealBlocks[i] = big.NewInt(depositRevealBlock)
	}

	result.SweepTxFee = big.NewInt(dsp.SweepTxFee)

	return result, nil
}

// UnmarshalJSON implements a custom JSON unmarshaling logic to produce a
// proper ProposeSweepTestScenario.
func (psts *ProposeSweepTestScenario) UnmarshalJSON(data []byte) error {
	type proposeSweepTestScenario struct {
		Title               string
		WalletPublicKeyHash string
		DepositTxMaxFee     uint64
		Deposits            []struct {
			FundingTxHash          string
			FundingOutputIndex     uint32
			RevealBlock            uint64
			FundingTxConfirmations uint
		}
		SweepTxFee                   int64
		EstimateSatPerVByteFee       int64
		ExpectedDepositSweepProposal *depositSweepProposal
		ExpectedErr                  string
	}

	var unmarshaled proposeSweepTestScenario

	err := json.Unmarshal(data, &unmarshaled)
	if err != nil {
		return err
	}

	// Unmarshal title.
	psts.Title = unmarshaled.Title

	// Unmarshal wallet public key hash.
	if len(unmarshaled.WalletPublicKeyHash) > 0 {
		copy(psts.WalletPublicKeyHash[:], hexToSlice(unmarshaled.WalletPublicKeyHash))
	}

	// Unmarshal deposit transaction max fee.
	psts.DepositTxMaxFee = unmarshaled.DepositTxMaxFee

	// Unmarshal deposits.
	for i, deposit := range unmarshaled.Deposits {
		d := new(ProposeSweepDepositsData)

		fundingTxHash, err := bitcoin.NewHashFromString(deposit.FundingTxHash, bitcoin.ReversedByteOrder)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal funding transaction hash for deposit [%d/%d]: [%w]",
				i,
				len(unmarshaled.Deposits),
				err,
			)
		}

		d.FundingTxHash = fundingTxHash
		d.FundingOutputIndex = deposit.FundingOutputIndex
		d.RevealBlock = deposit.RevealBlock
		d.FundingTxConfirmations = deposit.FundingTxConfirmations

		psts.Deposits = append(psts.Deposits, d)
	}

	// Unmarshal sweep transaction fee.
	psts.SweepTxFee = unmarshaled.SweepTxFee

	// Unmarshal estimate sat per vbyte fee.
	psts.EstimateSatPerVByteFee = unmarshaled.EstimateSatPerVByteFee

	// Unmarshal deposit sweep proposal
	psts.ExpectedDepositSweepProposal, err = unmarshaled.ExpectedDepositSweepProposal.convert()
	if err != nil {
		return fmt.Errorf(
			"failed to unmarshal expected deposit sweep proposal: [%w]",
			err,
		)
	}

	// Unmarshal expected error
	if len(unmarshaled.ExpectedErr) > 0 {
		psts.ExpectedErr = fmt.Errorf(unmarshaled.ExpectedErr)
	}

	return nil
}

// UnmarshalJSON implements a custom JSON unmarshaling logic to produce a
// proper FindPendingRedemptionsTestScenario.
func (fprts *FindPendingRedemptionsTestScenario) UnmarshalJSON(data []byte) error {
	type findPendingRedemptionsTestScenario struct {
		Title           string
		ChainParameters struct {
			AverageBlockTime int64
			CurrentBlock     uint64
			RequestTimeout   uint32
			RequestMinAge    uint32
		}
		Filter struct {
			WalletPublicKeyHashes []string
			WalletsLimit          uint16
			RequestsLimit         uint16
			RequestAmountLimit    uint64
		}
		Wallets []struct {
			WalletPublicKeyHash     string
			RegistrationBlockNumber uint64
		}
		PendingRedemptions []struct {
			WalletPublicKeyHash  string
			RedeemerOutputScript string
			RequestedAmount      uint64
			Age                  int64
		}
		ExpectedWalletsPendingRedemptions map[string][]string
	}

	var unmarshaled findPendingRedemptionsTestScenario

	err := json.Unmarshal(data, &unmarshaled)
	if err != nil {
		return err
	}

	fprts.Title = unmarshaled.Title

	fprts.ChainParameters.AverageBlockTime =
		time.Duration(unmarshaled.ChainParameters.AverageBlockTime) * time.Second
	fprts.ChainParameters.CurrentBlock = unmarshaled.ChainParameters.CurrentBlock
	fprts.ChainParameters.RequestTimeout = unmarshaled.ChainParameters.RequestTimeout
	fprts.ChainParameters.RequestMinAge = unmarshaled.ChainParameters.RequestMinAge

	for _, wpkhString := range unmarshaled.Filter.WalletPublicKeyHashes {
		var wpkh [20]byte
		copy(wpkh[:], hexToSlice(wpkhString))

		fprts.Filter.WalletPublicKeyHashes = append(
			fprts.Filter.WalletPublicKeyHashes,
			wpkh,
		)
	}

	fprts.Filter.WalletsLimit = unmarshaled.Filter.WalletsLimit
	fprts.Filter.RequestsLimit = unmarshaled.Filter.RequestsLimit
	fprts.Filter.RequestAmountLimit = unmarshaled.Filter.RequestAmountLimit

	for _, w := range unmarshaled.Wallets {
		var wpkh [20]byte
		copy(wpkh[:], hexToSlice(w.WalletPublicKeyHash))

		fprts.Wallets = append(fprts.Wallets, &Wallet{
			WalletPublicKeyHash:     wpkh,
			RegistrationBlockNumber: w.RegistrationBlockNumber,
		})
	}

	now := time.Now()
	currentBlock := fprts.ChainParameters.CurrentBlock
	averageBlockTime := fprts.ChainParameters.AverageBlockTime

	for _, pr := range unmarshaled.PendingRedemptions {
		var wpkh [20]byte
		copy(wpkh[:], hexToSlice(pr.WalletPublicKeyHash))

		age := time.Duration(pr.Age) * time.Second
		ageBlocks := uint64(age.Milliseconds() / averageBlockTime.Milliseconds())

		requestedAt := now.Add(-age)
		requestBlock := currentBlock - ageBlocks

		fprts.PendingRedemptions = append(
			fprts.PendingRedemptions,
			&RedemptionRequest{
				WalletPublicKeyHash:  wpkh,
				RedeemerOutputScript: hexToSlice(pr.RedeemerOutputScript),
				RequestedAmount:      pr.RequestedAmount,
				RequestedAt:          requestedAt,
				RequestBlock:         requestBlock,
			},
		)
	}

	fprts.ExpectedWalletsPendingRedemptions = make(map[[20]byte][]bitcoin.Script)
	for wpkhString, scripts := range unmarshaled.ExpectedWalletsPendingRedemptions {
		var wpkh [20]byte
		copy(wpkh[:], hexToSlice(wpkhString))

		convertedScripts := make([]bitcoin.Script, len(scripts))
		for i, script := range scripts {
			convertedScripts[i] = hexToSlice(script)
		}

		fprts.ExpectedWalletsPendingRedemptions[wpkh] = convertedScripts
	}

	return nil
}

func hexToSlice(hexString string) []byte {
	if len(hexString) == 0 {
		return []byte{}
	}

	bytes, err := hexutils.Decode(hexString)
	if err != nil {
		panic(err)
	}

	return bytes
}
