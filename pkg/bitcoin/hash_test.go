package bitcoin

import (
	"fmt"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"reflect"
	"testing"
)

func TestHashConversions(t *testing.T) {
	// A hash string in the internal byte order.
	hashString := "5672b911ab0dcc31bb36725de6f4d0c608983da7435443d69ae47e5fc151d909"
	// Same hash string in the reversed byte order.
	reversedHashString := "09d951c15f7ee49ad6435443a73d9808c6d0f4e65d7236bb31cc0dab11b97256"

	// Create the hash using the hash string in the internal byte order.
	hash, err := NewHashFromString(hashString, InternalByteOrder)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertStringsEqual(
		t,
		"hash string in the internal byte order",
		hashString,
		hash.String(InternalByteOrder),
	)

	testutils.AssertStringsEqual(
		t,
		"hash string in the reversed byte order",
		reversedHashString,
		hash.String(ReversedByteOrder),
	)

	// Create the same hash using the hash string in the reversed byte order.
	hashFromReversed, err := NewHashFromString(
		reversedHashString,
		ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	// The internal representation of the hash should be the same regardless
	// how the hash was created.
	testutils.AssertBytesEqual(t, hash[:], hashFromReversed[:])

	// Make sure we have an error if the hash string has a wrong size.
	_, err = NewHashFromString("0x"+hashString, InternalByteOrder)

	expectedErr := fmt.Errorf("wrong hash string size")
	if !reflect.DeepEqual(expectedErr, err) {
		t.Errorf(
			"unexpected error\n"+
				"expected: [%v]\n"+
				"actual:   [%v]",
			expectedErr,
			err,
		)
	}
}

func TestComputeHash(t *testing.T) {
	testutils.AssertStringsEqual(
		t,
		"hash one",
		"a26baf5a9a07d9eb7ba10f43924dcdf3f75f0abf066cd9f0c76f983121302e01",
		ComputeHash([]byte{0x01, 0x02, 0x03, 0x04, 0x05}).String(InternalByteOrder),
	)

	testutils.AssertStringsEqual(
		t,
		"hash two",
		"648f1115f075ad7c3998a0827fb20c61f37308252ad1778a0ae0b9f3268d4b92",
		ComputeHash([]byte{0xAA, 0xBB, 0xCC, 0xDD, 0xEE}).String(InternalByteOrder),
	)

	testutils.AssertStringsEqual(
		t,
		"hash three",
		"1397459010d34bc0a14dd4a5cf9e091cf6b4ebc6a1b21d315135af4139a73b3a",
		ComputeHash([]byte{0x01, 0xAA, 0x02, 0xBB, 0x03}).String(InternalByteOrder),
	)
}
