package relay

import "math/big"

// https://www.wolframalpha.com/input/?i=pi+to+78+digits
const piAsString = "31415926535897932384626433832795028841971693993751058209749445923078164062862"

// Genesis is the seed value for the network. The n digits of pi that fit into a
// *big.Int represent a "nothing up our sleeve" value that all consumers of this
// network can verify.
func GenesisEntryValue() *big.Int {
	value, success := big.NewInt(0).SetString(piAsString, 10)
	if !success {
		panic("failed to parse Pi as string")
	}
	return value
}
