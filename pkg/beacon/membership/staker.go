package membership

import "github.com/btcsuite/btcd/btcec"

// Staker represents an on-chain identity and staked amount.
type Staker struct {
	PubKey btcec.PublicKey // Q_j
	Weight uint64
}
