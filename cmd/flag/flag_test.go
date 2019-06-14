package flag_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/cmd/flag"
)

func TestTransactionHashSetting(t *testing.T) {
	tests := map[string]struct {
		value        string
		errorMessage string
		hash         common.Hash
	}{
		"valid hash": {
			value:        "0x1234000000000000000000000000000000000000000000000000000000001234",
			errorMessage: "",
			hash:         common.Hash([32]byte{18, 52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 18, 52}),
		},
		"short hash": {
			value:        "0x1234",
			errorMessage: "[0x1234] has [2] bytes, must be a hex string of exactly 32 bytes to be a valid transaction hash",
			hash:         common.Hash{},
		},
		"long hash": {
			value:        "0x123400000000000000000000000000000000000000000000000000000000123400",
			errorMessage: "[0x123400000000000000000000000000000000000000000000000000000000123400] has [33] bytes, must be a hex string of exactly 32 bytes to be a valid transaction hash",
			hash:         common.Hash{},
		},
		"decimal number": {
			value:        "000000000000000000000000000000000000001234",
			errorMessage: "[000000000000000000000000000000000000001234] must be a hex string starting with 0x to be a valid transaction hash",
			hash:         common.Hash{},
		},
		"blank string": {
			value:        "",
			errorMessage: "[] must be a hex string starting with 0x to be a valid transaction hash",
			hash:         common.Hash{},
		},
		"empty hex": {
			value:        "0x",
			errorMessage: "[0x] has [0] bytes, must be a hex string of exactly 32 bytes to be a valid transaction hash",
			hash:         common.Hash{},
		},
		"arbitrary string": {
			value:        "I am a booyan with booyans",
			errorMessage: "[I am a booyan with booyans] must be a hex string starting with 0x to be a valid transaction hash",
			hash:         common.Hash{},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			hash := &flag.TransactionHash{}
			err := hash.Set(test.value)

			message := ""
			if err != nil {
				message = err.Error()
			}

			if message != test.errorMessage {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.errorMessage,
					err,
				)
			}

			if err == nil && *hash.Hash != test.hash {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.hash,
					hash.Hash,
				)
			}
		})
	}
}

func TestTransactionHashStringification(t *testing.T) {
	hash := common.Hash([32]byte{18, 52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 18, 52})

	tests := map[string]struct {
		hash   *common.Hash
		string string
	}{
		"present hash": {
			hash:   &hash,
			string: "0x1234000000000000000000000000000000000000000000000000000000001234",
		},
		"missing hash": {
			hash:   nil,
			string: "unset",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			hash := &flag.TransactionHash{Hash: test.hash}
			string := hash.String()

			if string != test.string {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.string,
					string,
				)
			}
		})
	}
}

func TestUint256Setting(t *testing.T) {
	bigNumber, ok := (&big.Int{}).SetString("8233507321867270975858166353462000283756074959440384344846684898023608160820", 10)
	if !ok {
		t.Fatalf("Couldn't set comparison big number.")
	}

	tests := map[string]struct {
		value        string
		errorMessage string
		uint         *big.Int
	}{
		"256-bit number": {
			value:        "0x1234000000000000000000000000000000000000000000000000000000001234",
			errorMessage: "",
			uint:         bigNumber,
		},
		"<256-bit number": {
			value:        "0x1234",
			errorMessage: "",
			uint:         big.NewInt(int64(4660)),
		},
		">256-bit number": {
			value:        "0x123400000000000000000000000000000000000000000000000000000000123400",
			errorMessage: "[0x123400000000000000000000000000000000000000000000000000000000123400] must be a positive 256-bit or smaller hex or decimal string",
			uint:         nil,
		},
		"decimal number": {
			value:        "000000000000000000000000000000000000001234",
			errorMessage: "",
			uint:         big.NewInt(int64(1234)),
		},
		"negative number": {
			value:        "-0x1234000000000000000000000000000000000000000000000000000000001234",
			errorMessage: "[-0x1234000000000000000000000000000000000000000000000000000000001234] must be a positive 256-bit or smaller hex or decimal string",
			uint:         nil,
		},
		"negative decimal number": {
			value:        "-000000000000000000000000000000000000001234",
			errorMessage: "[-000000000000000000000000000000000000001234] must be a positive 256-bit or smaller hex or decimal string",
			uint:         nil,
		},
		"blank string": {
			value:        "",
			errorMessage: "[] must be a positive 256-bit or smaller hex or decimal string",
			uint:         &big.Int{},
		},
		"empty hex": {
			value:        "0x",
			errorMessage: "[0x] must be a positive 256-bit or smaller hex or decimal string",
			uint:         &big.Int{},
		},
		"arbitrary string": {
			value:        "I am a booyan with booyans",
			errorMessage: "[I am a booyan with booyans] must be a positive 256-bit or smaller hex or decimal string",
			uint:         &big.Int{},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			uintFlag := &flag.Uint256{}
			err := uintFlag.Set(test.value)

			message := ""
			if err != nil {
				message = err.Error()
			}

			if message != test.errorMessage {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.errorMessage,
					err,
				)
			}

			if err == nil && uintFlag.Uint.Cmp(test.uint) != 0 {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.uint.String(),
					uintFlag.Uint.String(),
				)
			}
		})
	}
}

func TestUint256Stringification(t *testing.T) {
	bigNumber, ok := (&big.Int{}).SetString("8233507321867270975858166353462000283756074959440384344846684898023608160820", 10)
	if !ok {
		t.Fatalf("Couldn't set comparison big number.")
	}

	tests := map[string]struct {
		uint   *big.Int
		string string
	}{
		"present uint": {
			uint:   bigNumber,
			string: "0x1234000000000000000000000000000000000000000000000000000000001234",
		},
		"missing uint": {
			uint:   nil,
			string: "unset",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			uint := &flag.Uint256{test.uint}
			string := uint.String()

			if string != test.string {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.string,
					string,
				)
			}
		})
	}
}
