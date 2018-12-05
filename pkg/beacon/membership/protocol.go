// Package membership contains code that implements the Random Beacon Group
// Selection protocol as described in
// http://docs.keep.network/the-beaconness-of-keep/random-beacon.pdf
package membership

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// Ticket is a message containing a pseudorandomly generated value, W_k, which is
// used to determine whether a given virtual staker is eligible for the group P
// (the lowest N tickets will be chosen) and a proof of the validity of the value
type Ticket struct {
	Value [sha256.Size]byte

	// Proof
	PreviousBeaconValue []byte
	StakerValue         []byte // Staker-specific value, Q_j
	VirtualStakerIndex  uint64
}

// calculateTicket generates a Ticket from the previous beacon output, the
// staker's ECDSA public key, and the virtual staker index. This function is
// intended to be called in a loop, ranging over the list of virtual stakers.

// See Phase 1 of the Group Selection protocol specification.
func (s *Staker) calculateTicket(
	beaconOutput []byte,
	virtualStakerIndex uint64,
) (*Ticket, error) {
	if virtualStakerIndex > s.Weight || virtualStakerIndex < 1 {
		return nil, fmt.Errorf(
			"virtualStakerIndex not in range [1, %d]",
			s.Weight,
		)
	}
	var combinedProof []byte

	combinedProof = append(combinedProof, beaconOutput...)
	combinedProof = append(combinedProof, s.PubKey.SerializeCompressed()...)
	binary.LittleEndian.PutUint64(
		combinedProof[len(combinedProof)-1:],
		virtualStakerIndex,
	)
	value := sha256.Sum256(combinedProof[:])

	return &Ticket{
		Value:               value,
		PreviousBeaconValue: beaconOutput,
		StakerValue:         s.PubKey.SerializeCompressed(),
		VirtualStakerIndex:  virtualStakerIndex,
	}, nil
}
