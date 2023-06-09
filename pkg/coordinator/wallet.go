package coordinator

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/internal/hexutils"
)

// TODO: Consider moving it to the tbtc package.
// WalletPublicKeyHash is a type representing a public key hash of the wallet.
type WalletPublicKeyHash [20]byte

// NewWalletPublicKeyHash decodes a wallet public key hash from a string.
func NewWalletPublicKeyHash(str string) (WalletPublicKeyHash, error) {
	var result WalletPublicKeyHash

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

func (w WalletPublicKeyHash) String() string {
	return hexutils.Encode(w[:])
}
