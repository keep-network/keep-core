package chain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Ticket represents group selection ticket as seen on-chain.
type Ticket struct {
	Value *big.Int // W_k
	Proof *TicketProof
}

// ValueBytes returns ticket value as a byte slice.
func (t *Ticket) ValueBytes() []byte {
	return common.LeftPadBytes(t.Value.Bytes(), 8)[:8]
}

// TicketProof represents group selection ticket proof as seen on-chain.
type TicketProof struct {
	StakerValue        *big.Int
	VirtualStakerIndex *big.Int
}
