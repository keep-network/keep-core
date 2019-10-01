package groupselection

import (
	"math/big"

	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/internal/byteutils"
)

// ticket is a message containing a pseudorandomly generated value, W_k, which is
// used to determine whether a given virtual staker is eligible for the group P
// (the lowest N tickets will be chosen) and a proof of the validity of the value
type ticket struct {
	Value shaValue // W_k

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

// calculateTicketValue generates a shaValue from the previous beacon output, the
// staker-specific value, and the virtual staker index.
//
// See Phase 2 of the Group Selection protocol specification.
func calculateTicketValue(
	beaconOutput []byte,
	stakerValue []byte,
	virtualStakerIndex *big.Int,
) (shaValue, error) {
	var combinedValue []byte
	var keccak256Hash shaValue

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

	return shaValue(keccak256Hash), nil
}
