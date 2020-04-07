package chain

import (
	"math/big"
)

// Ticket represents group selection ticket as seen on-chain.
type Ticket struct {
	Value [8]byte // W_k
	Proof *TicketProof
}

// IntValue returns ticket value as a big integer.
func (t *Ticket) IntValue() *big.Int {
	return new(big.Int).SetBytes(t.Value[:])
}

// TicketProof represents group selection ticket proof as seen on-chain.
type TicketProof struct {
	StakerValue        *big.Int
	VirtualStakerIndex *big.Int
}
