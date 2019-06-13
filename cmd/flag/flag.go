package flag

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

// TransactionHash is a container for a transaction hash that can be set from
// and converted to a string. It is compatible with the interface used by
// urfave/cli to specify flag types.
type TransactionHash struct {
	Hash *common.Hash
}

// Set sets the transaction hash flag's value from the given string.
func (tg *TransactionHash) Set(value string) error {
	if len(value) < 2 || value[0:2] != "0x" {
		return fmt.Errorf(
			"[%v] must be a hex string starting with 0x to be a valid transaction hash",
			value,
		)
	}

	bytes := common.FromHex(value)
	if len(bytes) != common.HashLength {
		return fmt.Errorf(
			"[%v] has [%v] bytes, must be a hex string of exactly 32 bytes to be a valid transaction hash",
			value,
			len(bytes),
		)
	}

	hash := common.BytesToHash(bytes)
	tg.Hash = &hash

	return nil
}

// String returns a string representation of the transaction hash flag.
func (tg *TransactionHash) String() string {
	if tg.Hash == nil {
		return "unset"
	}

	return tg.Hash.Hex()
}

// Uint256 turns *big.Int into a flag.Value.
type Uint256 struct {
	Uint *big.Int
}

// Set sets the Uint256 flag's value from the given string.
func (u256 *Uint256) Set(s string) error {
	// TODO It would be really nice to give more guidance here, e.g. the number
	// TODO is too big vs simply invalid.
	int, ok := math.ParseBig256(s)
	if !ok || len(s) == 0 || int.Sign() == -1 {
		return fmt.Errorf(
			"[%v] must be a positive 256-bit or smaller hex or decimal string",
			s,
		)
	}

	u256.Uint = int
	return nil
}

// String returns a string representation of the Uint256 flag as hex.
func (u256 *Uint256) String() string {
	if u256.Uint == nil {
		return "unset"
	}

	return "0x" + u256.Uint.Text(16)
}
