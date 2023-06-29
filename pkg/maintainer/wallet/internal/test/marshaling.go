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
		ExpectedDepositSweepProposal depositSweepProposal
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
