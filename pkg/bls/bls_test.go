package bls

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/altbn128"
)

// Test verifying BLS multi signature.
func TestMultisigBLS(t *testing.T) {

	// Public keys and signatures to aggregate.
	var signatures []*bn256.G1
	var publicKeys []*bn256.G2

	// Generator point of G2 group.
	p2 := new(bn256.G2).ScalarBaseMult(big.NewInt(1))

	// Message to sign.
	msg := altbn128.G1HashToPoint([]byte("Hello!"))

	for i := 0; i < 100; i++ {
		// Get private key.
		k, _, err := bn256.RandomG1(rand.Reader)

		if err != nil {
			t.Errorf("Error generating random point on G1")
		}

		// Get public key.
		pub := new(bn256.G2).ScalarBaseMult(k)
		publicKeys = append(publicKeys, pub)

		// Sign the message.
		sig := new(bn256.G1).ScalarMult(msg, k)
		signatures = append(signatures, sig)
	}

	aggSig := AggregateG1Points(signatures)
	negAggSig := new(bn256.G1).Neg(aggSig)
	aggPub := AggregateG2Points(publicKeys)

	// Perform 2 pairing operations.
	a := []*bn256.G1{negAggSig, msg}
	b := []*bn256.G2{p2, aggPub}
	pairingCheck := bn256.PairingCheck(a, b)

	if !pairingCheck {
		t.Errorf("Error verifying BLS multi signature.")
	}
}
