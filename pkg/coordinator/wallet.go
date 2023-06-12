package coordinator

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/internal/hexutils"
)

// NewWalletPublicKeyHash decodes a wallet public key hash from a string.
func NewWalletPublicKeyHash(str string) ([20]byte, error) {
	var result [20]byte

	walletHex, err := hexutils.Decode(str)
	if err != nil {
		return result, err
	}

	if len(walletHex) != 20 {
		return result, fmt.Errorf("invalid bytes length: [%d], expected: [%d]", len(walletHex), 20)
	}

	copy(result[:], walletHex)

	return result, nil
}
