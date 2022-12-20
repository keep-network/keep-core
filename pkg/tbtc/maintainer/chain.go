package maintainer

import "github.com/keep-network/keep-core/pkg/bitcoin"

// BitcoinDifficultyChain is an interface that provides the ability to
// communicate with the Bitcoin difficulty on-chain contract.
type BitcoinDifficultyChain interface {
	// Retarget adds a new epoch to the Bitcoin difficulty relay by providing
	// a proof of the difficulty before and after the retarget.
	Retarget(headers []bitcoin.BlockHeader) error
}

// TODO: Description
type WalletChain interface {
	ActiveWalletPubKeyHash() ([20]byte, error)

	// GetWalletCreationState returns the wallet current wallet creation state
	// in the wallet registry.
	GetWalletCreationState() (DKGState, error)

	WalletParameters() (
		walletCreationPeriod uint32,
		walletCreationMinBtcBalance uint64,
		walletCreationMaxBtcBalance uint64,
		err error,
	)

	GetWalletInfo(walletPubKeyHash [20]byte) (
		publicKeyBytes []byte,
		mainUtxoHash [32]byte,
		createdAt uint32,
		err error,
	)
}

// TODO: Reuse the tbtc.DKGState. It cannot be used now because of a import cycle.
type DKGState int

const (
	Idle DKGState = iota
	AwaitingSeed
	AwaitingResult
	Challenge
)
