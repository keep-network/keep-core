// Package blsutils provides helper utilities for working with the
// go-dfinity-crypto library's BLS types. These utilities are mostly aimed at
// testing.
package blsutils

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/dfinity/go-dfinity-crypto/bls"
)

// GenerateID generates a random BLS ID and returns a pointer to it.
func GenerateID() *bls.ID {
	id := bls.ID{}
	idValue := fmt.Sprintf("%v", rand.Int31())
	err := id.SetDecString(idValue)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to generate id from random number %v: [%v]",
			idValue,
			err))
	}

	return &id
}

// GeneratePublicKey generates a random BLS public key and returns a pointer
// to it. Note that the associated secret key is not returned and cannot be
// recovered from the public key.
func GeneratePublicKey() *bls.PublicKey {
	sk := bls.SecretKey{}
	sk.SetByCSPRNG()
	pk := sk.GetPublicKey()

	return pk
}

// GenerateSecretKey generates a random BLS secret key and returns a pointer
// to it.
func GenerateSecretKey() *bls.SecretKey {
	sk := bls.SecretKey{}
	sk.SetByCSPRNG()

	return &sk
}

// AssertIDsEqual takes a testing.T and two ids and reports an error if the two
// ids are not equal.
func AssertIDsEqual(t *testing.T, id *bls.ID, roundTripID *bls.ID) {
	if !id.IsEqual(roundTripID) {
		t.Errorf(
			"IDs are not equal\nexpected: [%s]\nactual:   [%s]",
			id.GetHexString(),
			roundTripID.GetHexString())
	}
}
