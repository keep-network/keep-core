// Package pbutils provides helper utilities for working with protobuf objects.
// These utilities are mostly aimed at testing.
package pbutils

import "github.com/gogo/protobuf/proto"

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
