package chain

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestAddressesToString(t *testing.T) {
	var tests = map[string]struct {
		addresses      Addresses
		expectedString string
	}{
		"empty addresses": {
			addresses:      Addresses{},
			expectedString: "[]",
		},
		"one address": {
			addresses:      Addresses{"0x123"},
			expectedString: "[1: 0x123]",
		},
		"two addresses": {
			addresses:      Addresses{"0x123", "0x234"},
			expectedString: "[1: 0x123, 2: 0x234]",
		},
		"three addresses": {
			addresses:      Addresses{"0x123", "0x234", "0x345"},
			expectedString: "[1: 0x123, 2: 0x234, 3: 0x345]",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualString := test.addresses.String()
			testutils.AssertStringsEqual(
				t,
				"ToString() result",
				test.expectedString,
				actualString,
			)
		})
	}
}
