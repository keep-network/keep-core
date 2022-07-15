package ecdsa

// WalletID is a unique identifier of an ECDSA wallet.
type WalletID [32]byte

// Wallet represents an identifiable group of multiple ThresholdSigner
// ready to process signing requests.
type Wallet struct {
	ID WalletID
	// TODO: Implementation.
}

type WalletStorage struct {
	// TODO: Implementation.
}
