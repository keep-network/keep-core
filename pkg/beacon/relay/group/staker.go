package group

import (
	"sort"

	"github.com/btcsuite/btcd/btcec"
)

// Staker represents an on-chain identity and staked amount.
type Staker struct {
	StakeID string
	PubKey  *btcec.PublicKey // Q_j

	// A staker's VirtualStakers is how many minimum-stake stakers a given
	// actual staker could form if they were to blitzpants their stake.
	VirtualStakers uint64
}

// Returns a new Staker with the calculated weight and on-chain public key.
func NewStaker(pubKey *btcec.PublicKey, weight uint64) *Staker {
	return &Staker{
		PubKey:         pubKey,
		VirtualStakers: weight,
	}
}

// GenerateTickets takes the previous beacon output, and calculates the number of
// tickets corresponding to the number of virtual stakers the staker could form.
func (s *Staker) GenerateTickets(beaconOutput []byte) (Tickets, error) {
	var tickets Tickets
	// VirtualStakers are 1-indexed.
	for i := uint64(1); i <= s.VirtualStakers; i++ {
		ticket, err := s.calculateTicket(
			beaconOutput, i,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	sort.Sort(tickets)
	return tickets, nil
}
