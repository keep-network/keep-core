package groupselection

import (
	"bytes"
	"math/big"

	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/internal/byteutils"
)

// ticket is a message containing a pseudorandomly generated value, W_k, which is
// used to determine whether a given virtual staker is eligible for the group P
// (the lowest N tickets will be chosen) and a proof of the validity of the value
type ticket struct {
	Value SHAValue // W_k

	Proof *Proof // proof_k = Proof(Q_j, vs)
}

// Proof consists of the components needed to construct the ticket's value, and
// also acts as evidence for an accusing challenge against the ticket's value.
type Proof struct {
	StakerValue        []byte   // Q_j, a staker-specific value
	VirtualStakerIndex *big.Int // vs
}

// newTicket calculates a ticket Value (SHAValue), and returns the ticket with
// the associated Proof.
func newTicket(
	beaconOutput []byte, // V_i
	stakerValue []byte, // Q_j
	virtualStakerIndex *big.Int, // vs
) (*ticket, error) {
	value, err := calculateTicketValue(beaconOutput, stakerValue, virtualStakerIndex)
	if err != nil {
		return nil, fmt.Errorf("ticket value calculation failed [%v]", err)
	}

	return &ticket{
		Value: value,
		Proof: &Proof{
			StakerValue:        stakerValue,
			VirtualStakerIndex: virtualStakerIndex,
		},
	}, nil
}

// IsFromStaker compare ticket staker value against staker address
func (t *ticket) IsFromStaker(stakerAddress []byte) bool {
	return bytes.Compare(t.Proof.StakerValue, stakerAddress) == 0
}

// calculateTicketValue generates a SHAValue from the previous beacon output, the
// staker-specific value, and the virtual staker index.
//
// See Phase 2 of the Group Selection protocol specification.
func calculateTicketValue(
	beaconOutput []byte,
	stakerValue []byte,
	virtualStakerIndex *big.Int,
) (SHAValue, error) {
	var combinedValue []byte
	var keccak256Hash SHAValue

	beaconOutputPadded, err := byteutils.LeftPadTo32Bytes(beaconOutput)
	if err != nil {
		return keccak256Hash, fmt.Errorf("cannot pad a becon output, [%v]", err)
	}

	stakerValuePadded, err := byteutils.LeftPadTo32Bytes(stakerValue)
	if err != nil {
		return keccak256Hash, fmt.Errorf("cannot pad a staker value, [%v]", err)
	}

	virtualStakerIndexPadded, err := byteutils.LeftPadTo32Bytes(virtualStakerIndex.Bytes())
	if err != nil {
		return keccak256Hash, fmt.Errorf("cannot pad a virtual staker index, [%v]", err)
	}

	combinedValue = append(combinedValue, beaconOutputPadded...)
	combinedValue = append(combinedValue, stakerValuePadded...)
	combinedValue = append(combinedValue, virtualStakerIndexPadded...)

	copy(keccak256Hash[:], crypto.Keccak256(combinedValue[:]))

	return SHAValue(keccak256Hash), nil

}

// tickets implements sort.Interface
type tickets []*ticket

// Len is the sort.Interface requirement for Tickets
func (ts tickets) Len() int {
	return len(ts)
}

// Swap is the sort.Interface requirement for Tickets
func (ts tickets) Swap(i, j int) {
	ts[i].Proof.VirtualStakerIndex, ts[j].Proof.VirtualStakerIndex =
		ts[j].Proof.VirtualStakerIndex, ts[i].Proof.VirtualStakerIndex
}

// Less is the sort.Interface requirement for Tickets
func (ts tickets) Less(i, j int) bool {
	iVirtualStakeIndex := ts[i].Proof.VirtualStakerIndex
	jVirtualStakerIndex := ts[j].Proof.VirtualStakerIndex

	switch iVirtualStakeIndex.Cmp(jVirtualStakerIndex) {
	case -1:
		return true
	case 1:
		return false
	}

	return true
}
