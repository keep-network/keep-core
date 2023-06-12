package coordinatortest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/coordinator"
)

// UnmarshalJSON implements a custom JSON unmarshaling logic to produce a
// proper DepositsSweepTestScenario.
func (dsts *DepositsSweepTestScenario) UnmarshalJSON(data []byte) error {
	type depositsSweepTestScenario struct {
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

	var unmarshaled depositsSweepTestScenario

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

		ud.FundingTransactionHash = fundingTxHash
		ud.FundingTransactionOutputIndex = deposit.FundingOutputIndex
		ud.RevealBlock = deposit.RevealBlockNumber

		dsts.ExpectedUnsweptDeposits = append(dsts.ExpectedUnsweptDeposits, ud)
	}

	return nil
}
