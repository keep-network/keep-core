package chain

import (
	"math/big"
)

// Ticket represents group selection ticket as seen on-chain.
type Ticket struct {
	Value [8]byte // W_k
	Proof *TicketProof
}

// TicketProof represents group selection ticket proof as seen on-chain.
type TicketProof struct {
	StakerValue        *big.Int
	VirtualStakerIndex *big.Int
}
