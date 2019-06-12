package flag

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// TransactionHash is a container for a transaction hash that can be set from
// and converted to a string. It is compatible with the interface used by
// urfave/cli to specify flag types.
type TransactionHash struct {
	Hash *common.Hash
}

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

func (tg *TransactionHash) String() string {
	if tg.Hash != nil {
		return tg.Hash.Hex()
	} else {
		return "unset"
	}
}
