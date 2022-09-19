// Package pbutils provides helper utilities for working with protobuf objects.
// These utilities are mostly aimed at testing.
package pbutils

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	fuzz "github.com/google/gofuzz"

	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/internal/pb"
)

// RoundTrip takes a marshaler and unmarshaler, marshals the marshaler, and then
// unmarshals the result into the unmarshaler. If either procedure errors out,
// it returns an error; otherwise it returns nil and the unmarshaler is left
// with the results of the round-trip.
//
// This is a utility meant to facilitate tests that verify round-trip marshaling
// of objects with custom protobuf marshaling.
func RoundTrip(
	marshaler pb.Marshaler,
	unmarshaler pb.Unmarshaler,
) error {
	bytes, err := marshaler.Marshal()
	if err != nil {
		return err
	}

	err = unmarshaler.Unmarshal(bytes)
	if err != nil {
		return err
	}

	return nil
}

// FuzzUnmarshaler tests given unmarshaler with random bytes.
func FuzzUnmarshaler(unmarshaler pb.Unmarshaler) {
	for i := 0; i < 100; i++ {
		var messageBytes []byte

		f := fuzz.New().NilChance(0.01).NumElements(0, 512)
		f.Fuzz(&messageBytes)

		_ = unmarshaler.Unmarshal(messageBytes)
	}
}

// FuzzFuncs returns custom fuzzing functions set.
func FuzzFuncs() []interface{} {
	return []interface{}{
		fuzzBigInt(),
		fuzzEphemeralPublicKey(),
		fuzzEphemeralPrivateKey(),
		fuzzG1(),
		fuzzG2(),
	}
}

func fuzzBigInt() func(*big.Int, fuzz.Continue) {
	return func(int *big.Int, c fuzz.Continue) {
		var abs []big.Word

		c.Fuzz(&abs)

		int.SetBits(abs)
	}
}

func fuzzEphemeralPublicKey() func(*ephemeral.PublicKey, fuzz.Continue) {
	return func(key *ephemeral.PublicKey, c fuzz.Continue) {
		var x, y big.Int

		c.Fuzz(&x)
		c.Fuzz(&y)

		key.Curve = btcec.S256()
		key.X = &x
		key.Y = &y
	}
}

func fuzzEphemeralPrivateKey() func(*ephemeral.PrivateKey, fuzz.Continue) {
	return func(key *ephemeral.PrivateKey, c fuzz.Continue) {
		var (
			publicKey ephemeral.PublicKey
			d         big.Int
		)

		c.Fuzz(&publicKey)
		c.Fuzz(&d)

		key.PublicKey = ecdsa.PublicKey(publicKey)
		key.D = &d
	}
}

func fuzzG1() func(*bn256.G1, fuzz.Continue) {
	return func(g1 *bn256.G1, c fuzz.Continue) {
		var k big.Int

		c.Fuzz(&k)

		g1.ScalarBaseMult(&k)
	}
}

func fuzzG2() func(*bn256.G2, fuzz.Continue) {
	return func(g2 *bn256.G2, c fuzz.Continue) {
		var k big.Int

		c.Fuzz(&k)

		// trim k to reasonable number of bytes to prevent long execution
		if len(k.Bytes()) > 64 {
			k = *new(big.Int).SetBytes(k.Bytes()[:64])
		}

		g2.ScalarBaseMult(&k)
	}
}
