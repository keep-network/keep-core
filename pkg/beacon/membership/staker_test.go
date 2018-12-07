package membership

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec"
)

func TestGenerateTickets(t *testing.T) {
	staker, err := newTestStaker(10)
	if err != nil {
		t.Fatal(err)
	}
	previousBeaconOutput := []byte("test beacon output")

	tickets, err := staker.GenerateTickets(previousBeaconOutput)
	if err != nil {
		t.Fatal(err)
	}

	// we should have 10 tickets
	if len(tickets) != 10 {
		t.Fatal("bad things in paradise")
	}

	for _, ticket := range tickets {
		if ticket.VirtualStakerIndex == 0 {
			t.Fatal("Virutal stakers should be 1-indexed, not 0-indexed")
		}
	}
}

func TestValidateProofs(t *testing.T) {
	staker, err := newTestStaker(1)
	if err != nil {
		t.Fatal(err)
	}

	beaconOutput := []byte("test beacon output")

	// hash(proof) == value?
	var proofBytes []byte
	proofBytes = append(proofBytes, beaconOutput...)
	proofBytes = append(proofBytes, staker.PubKey.SerializeCompressed()...)

	virtualStakerBytes := make([]byte, 64)
	binary.LittleEndian.PutUint64(virtualStakerBytes, staker.VirtualStakers)
	proofBytes = append(proofBytes, virtualStakerBytes...)

	expectedValue := sha256.Sum256(proofBytes[:])

	tickets, err := staker.GenerateTickets(beaconOutput)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(
		toByteSlice(tickets[0].Value),
		toByteSlice(expectedValue),
	) != 0 {
		t.Fatalf(
			"hashed value (%v) doesn't match ticket value (%v)",
			tickets[0].Value,
			expectedValue,
		)
	}

}

func newTestStaker(virtualStakers int) (*Staker, error) {
	ecdsaPrivateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate new ephemeral keypair [%v]",
			err,
		)
	}

	return NewStaker(ecdsaPrivateKey.PubKey(), uint64(virtualStakers)), nil
}
