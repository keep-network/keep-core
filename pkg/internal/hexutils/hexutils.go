package hexutils

import (
	"encoding/hex"
	"fmt"
)

// Decode decodes a hex string with 0x prefix.
func Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty hex string")
	}
	if has0xPrefix(input) {
		input = input[2:]
	}

	b, err := hex.DecodeString(input)

	if err != nil {
		return nil, fmt.Errorf("failed to decode string [%s]", input)
	}
	return b, err
}

// Encode encodes b as a hex string with 0x prefix.
func Encode(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}
