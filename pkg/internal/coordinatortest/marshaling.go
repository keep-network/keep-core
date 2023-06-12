package coordinatortest

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/coordinator"
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
		walletPublicKeyHash, err := coordinator.NewWalletPublicKeyHash(unmarshaled.WalletPublicKeyHash)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal wallet public key hash: [%w]",
				err,
			)
		}

		dsts.WalletPublicKeyHash = &walletPublicKeyHash
	}

	dsts.MaxNumberOfDeposits = unmarshaled.MaxNumberOfDeposits

	// Unmarshal wallets.
	for i, wallet := range unmarshaled.Wallets {
		w := new(Wallet)
		walletPublicKeyHash, err := coordinator.NewWalletPublicKeyHash(wallet.WalletPublicKeyHash)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal wallet public key hash for deposit [%d/%d]: [%w]",
				i,
				len(unmarshaled.Deposits),
				err,
			)
		}

		w.WalletPublicKeyHash = walletPublicKeyHash
		w.RegistrationBlockNumber = wallet.RegistrationBlockNumber

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

		walletPublicKeyHash, err := coordinator.NewWalletPublicKeyHash(deposit.WalletPublicKeyHash)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal wallet public key hash for deposit [%d/%d]: [%w]",
				i,
				len(unmarshaled.Deposits),
				err,
			)
		}

		d.FundingTxHash = fundingTxHash
		d.FundingOutputIndex = deposit.FundingOutputIndex
		d.FundingTxConfirmations = deposit.FundingTxConfirmations
		d.WalletPublicKeyHash = walletPublicKeyHash
		d.RevealBlockNumber = deposit.RevealBlockNumber
		d.SweptAt = time.Unix(deposit.SweptAt, 0)

		dsts.Deposits = append(dsts.Deposits, d)
	}

	// Unmarshal expected wallet public key has.
	if len(unmarshaled.ExpectedWalletPublicKeyHash) > 0 {
		expectedWalletPublicKeyHash, err := coordinator.NewWalletPublicKeyHash(unmarshaled.ExpectedWalletPublicKeyHash)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal expected wallet public key hash: [%w]",
				err,
			)
		}
		dsts.ExpectedWalletPublicKeyHash = expectedWalletPublicKeyHash
	}

	// Unmarshal expected unswept deposits.
	for i, deposit := range unmarshaled.ExpectedUnsweptDeposits {
		ud := new(coordinator.DepositSweepDetails)

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

	var err error

	if len(dsp.WalletPublicKeyHash) > 0 {
		result.WalletPublicKeyHash, err = coordinator.NewWalletPublicKeyHash(dsp.WalletPublicKeyHash)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to unmarshal wallet public key hash: [%w]",
				err,
			)
		}
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
		walletPublicKeyHash, err := coordinator.NewWalletPublicKeyHash(unmarshaled.WalletPublicKeyHash)
		if err != nil {
			return fmt.Errorf(
				"failed to unmarshal wallet public key hash: [%w]",
				err,
			)
		}

		psts.WalletPublicKeyHash = walletPublicKeyHash
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
