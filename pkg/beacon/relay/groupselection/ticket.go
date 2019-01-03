// Package groupselection contains code that implements the Random Beacon Group
// Selection protocol as described in
// http://docs.keep.network/the-beaconness-of-keep/random-beacon.pdf
package groupselection

import (
	"bytes"
	"crypto/sha256"
	"math/big"
)

// Ticket is a message containing a pseudorandomly generated value, W_k, which is
// used to determine whether a given virtual staker is eligible for the group P
// (the lowest N tickets will be chosen) and a proof of the validity of the value
type Ticket struct {
	Value SHAValue // W_k

	Proof *Proof // Proof(Q_j, vs)
}

// Proof consists of the components needed to construct the Ticket's value, and
// also acts as evidence for an accusing challenge against the Ticket's value.
type Proof struct {
	StakerValue        []byte   // Q_j, a staker-specific value
	VirtualStakerIndex *big.Int // vs
}

// calculateTicket generates a Ticket from the previous beacon output, the
// staker-specific value, and the virtual staker index.
//
// See Phase 2 of the Group Selection protocol specification.
func calculateTicket(
	beaconOutput []byte,
	stakerValue []byte,
	virtualStakerIndex *big.Int,
) *Ticket {
	var combinedValue []byte
	combinedValue = append(combinedValue, beaconOutput...)
	combinedValue = append(combinedValue, stakerValue...)
	combinedValue = append(combinedValue, virtualStakerIndex.Bytes()...)

	value := SHAValue(sha256.Sum256(combinedValue[:]))

	return &Ticket{
		Value: value,
		Proof: &Proof{
			StakerValue:        stakerValue,
			VirtualStakerIndex: virtualStakerIndex,
		},
	}
}

// tickets implements sort.Interface
type tickets []*Ticket

// Len is the sort.Interface requirement for Tickets
func (ts tickets) Len() int {
	return len(ts)
}

// Swap is the sort.Interface requirement for Tickets
func (ts tickets) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

// Less is the sort.Interface requirement for Tickets
func (ts tickets) Less(i, j int) bool {
	iBytes := ts[i].Value.Bytes()
	jBytes := ts[j].Value.Bytes()

	switch bytes.Compare(iBytes, jBytes) {
	case -1:
		return true
	case 1:
		return false
	}

	return true
}
