package groupselection

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec"
)

func TestGenerateTickets(t *testing.T) {
	minimumStake := big.NewInt(1)
	availableStake := big.NewInt(10)
	virtualStakers := availableStake.Int64() / minimumStake.Int64()

	stakingPublicKey, err := newTestPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	previousBeaconOutput := []byte("test beacon output")

	tickets, err := GenerateTickets(
		minimumStake,
		availableStake,
		stakingPublicKey.SerializeCompressed(),
		previousBeaconOutput,
	)
	if err != nil {
		t.Fatal(err)
	}

	// we should have 10 tickets
	if len(tickets) != int(virtualStakers) {
		t.Fatalf(
			"expected [%d] tickets, received [%d] tickets",
			virtualStakers,
			len(tickets),
		)
	}

	for _, ticket := range tickets {
		if ticket.Proof.VirtualStakerIndex == big.NewInt(0) {
			t.Fatal("Virutal stakers should be 1-indexed, not 0-indexed")
		}
	}
}

func TestValidateProofs(t *testing.T) {
	minimumStake := big.NewInt(1)
	availableStake := big.NewInt(1)
	virtualStakers := big.NewInt(0).Quo(availableStake, minimumStake)

	stakingPublicKey, err := newTestPublicKey()
	if err != nil {
		t.Fatal(err)
	}

	beaconOutput := []byte("test beacon output")

	// hash(proof) == expected value?
	var valueBytes []byte
	valueBytes = append(valueBytes, beaconOutput...)
	valueBytes = append(valueBytes, stakingPublicKey.SerializeCompressed()...)
	valueBytes = append(valueBytes, virtualStakers.Bytes()...)

	expectedValue := SHAValue(sha256.Sum256(valueBytes[:]))

	tickets, err := GenerateTickets(
		minimumStake,
		availableStake,
		stakingPublicKey.SerializeCompressed(),
		beaconOutput,
	)
	if err != nil {
		t.Fatal(err)
	}

	// we should have virtualStaker number of tickets
	if len(tickets) != int(virtualStakers.Int64()) {
		t.Fatalf(
			"expected [%d] tickets, received [%d] tickets",
			virtualStakers,
			len(tickets),
		)
	}

	if bytes.Compare(
		tickets[0].Value.Bytes(),
		expectedValue.Bytes(),
	) != 0 {
		t.Fatalf(
			"hashed value (%v) doesn't match ticket value (%v)",
			tickets[0].Value,
			expectedValue,
		)
	}
}

func newTestPublicKey() (*btcec.PublicKey, error) {
	ecdsaPrivateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate new ephemeral keypair [%v]",
			err,
		)
	}

	return ecdsaPrivateKey.PubKey(), nil
}
