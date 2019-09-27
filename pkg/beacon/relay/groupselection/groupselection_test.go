package groupselection

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/internal/byteutils"
)

func TestGenerateTickets(t *testing.T) {
	minimumStake := big.NewInt(1)
	availableStake := big.NewInt(10)
	virtualStakers := availableStake.Int64() / minimumStake.Int64()

	stakingPublicKey, err := newTestPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	stakingPublicKeyECDSA := stakingPublicKey.ToECDSA()
	stakingAddress := crypto.PubkeyToAddress(*stakingPublicKeyECDSA)
	previousBeaconOutput := []byte("test beacon output")

	tickets, err := GenerateTickets(
		previousBeaconOutput,
		stakingAddress.Bytes(),
		availableStake,
		minimumStake,
	)
	if err != nil {
		t.Fatal(err)
	}

	// We should have 10 tickets
	if len(tickets) != int(virtualStakers) {
		t.Fatalf(
			"expected [%d] tickets, received [%d] tickets",
			virtualStakers,
			len(tickets),
		)
	}

	for i, ticket := range tickets {
		expectedIndex := int64(i + 1)
		// Tickets should be sorted in ascending order
		if expectedIndex != ticket.Proof.VirtualStakerIndex.Int64() {
			t.Fatalf("Got index [%d], want index [%d]",
				ticket.Proof.VirtualStakerIndex,
				expectedIndex,
			)
		}

		if ticket.Proof.VirtualStakerIndex == big.NewInt(0) {
			t.Fatal("Virutal stakers should be 1-indexed, not 0-indexed")
		}
	}

}

func TestValidateProofs(t *testing.T) {
	beaconOutput := []byte("test beacon output")
	beaconOutputPadded, _ := byteutils.LeftPadTo32Bytes(beaconOutput)

	stakingPublicKey, err := newTestPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	stakingPublicKeyECDSA := stakingPublicKey.ToECDSA()
	stakingAddress := crypto.PubkeyToAddress(*stakingPublicKeyECDSA)
	stakerValuePadded, _ := byteutils.LeftPadTo32Bytes(stakingAddress.Bytes())

	minimumStake := big.NewInt(1)
	availableStake := big.NewInt(1)
	virtualStakers := big.NewInt(0).Quo(availableStake, minimumStake) // 1
	virtualStakerIndexPadded, _ := byteutils.LeftPadTo32Bytes(virtualStakers.Bytes())

	var valueBytes []byte

	valueBytes = append(valueBytes, beaconOutputPadded...) // V_i
	valueBytes = append(valueBytes, stakerValuePadded...)  // Q_j
	// only 1 virtual staker, which corresponds to the index, vs
	valueBytes = append(valueBytes, virtualStakerIndexPadded...)

	expectedValue := crypto.Keccak256(valueBytes[:])

	tickets, err := GenerateTickets(
		beaconOutput,
		stakingAddress.Bytes(),
		availableStake,
		minimumStake,
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
		expectedValue,
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
