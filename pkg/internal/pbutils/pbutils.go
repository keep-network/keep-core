// Package pbutils provides helper utilities for working with protobuf objects.
// These utilities are mostly aimed at testing.
package pbutils

import (
	"github.com/gogo/protobuf/proto"
	fuzz "github.com/google/gofuzz"
)

// RoundTrip takes a marshaler and unmarshaler, marshals the marshaler, and then
// unmarshals the result into the unmarshaler. If either procedure errors out,
// it returns an error; otherwise it returns nil and the unmarshaler is left
// with the results of the round-trip.
//
// This is a utility meant to facilitate tests that verify round-trip marshaling
// of objects with custom protobuf marshaling.
func RoundTrip(
	marshaler proto.Marshaler,
	unmarshaler proto.Unmarshaler,
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
func FuzzUnmarshaler(unmarshaler proto.Unmarshaler) {
	for i := 0; i < 100; i++ {
		var messageBytes []byte

		f := fuzz.New().NilChance(0.01).NumElements(0, 512)
		f.Fuzz(&messageBytes)

		_ = unmarshaler.Unmarshal(messageBytes)
	}
}
