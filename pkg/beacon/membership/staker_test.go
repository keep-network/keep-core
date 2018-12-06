package membership

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"testing"

	"github.com/btcsuite/btcd/btcec"
)

func TestGenerateTickets(t *testing.T) {
	ecdsaPrivateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		t.Fatalf(
			"could not generate new ephemeral keypair [%v]",
			err,
		)
	}

	staker := NewStaker(ecdsaPrivateKey.PubKey(), 10)
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
	ecdsaPrivateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		t.Fatalf(
			"could not generate new ephemeral keypair [%v]",
			err,
		)
	}

	staker := NewStaker(ecdsaPrivateKey.PubKey(), 1)
	beaconOutput := []byte("test beacon output")

	// hash(proof) == value?
	var proofBytes []byte
	proofBytes = append(proofBytes, beaconOutput...)
	proofBytes = append(proofBytes, staker.PubKey.SerializeCompressed()...)
	binary.LittleEndian.PutUint64(
		proofBytes,
		staker.VirtualStakers,
	)
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
