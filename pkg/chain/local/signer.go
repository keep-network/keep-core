package local

import (
	"encoding/hex"
	"fmt"
	commonlocal "github.com/keep-network/keep-common/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/operator"
)

// TODO: Consider moving the local `Signer` out of `keep-common` to this file.
type signer struct {
	*commonlocal.Signer
}

func newSigner(operatorPrivateKey *operator.PrivateKey) *signer {
	chainPrivateKey, _, err := operatorPrivateKeyToChainKeyPair(operatorPrivateKey)
	if err != nil {
		panic(err)
	}

	return &signer{
		commonlocal.NewSigner(chainPrivateKey),
	}
}

func (s *signer) PublicKeyToAddress(
	publicKey *operator.PublicKey,
) (chain.Address, error) {
	chainPublicKey, err := operatorPublicKeyToChainPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf(
			"cannot convert operator key to chain key: [%v]",
			err,
		)
	}

	addressBytes := s.Signer.PublicKeyToAddress(*chainPublicKey)

	return chain.Address(hex.EncodeToString(addressBytes)), nil
}

func (s *signer) PublicKeyBytesToAddress(publicKey []byte) chain.Address {
	addressBytes := s.Signer.PublicKeyBytesToAddress(publicKey)

	return chain.Address(hex.EncodeToString(addressBytes))
}
