package local

import (
	"fmt"
	commonlocal "github.com/keep-network/keep-common/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/operator"
)

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
) ([]byte, error) {
	chainPublicKey, err := operatorPublicKeyToChainPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot convert operator key to chain key: [%v]",
			err,
		)
	}

	return s.Signer.PublicKeyToAddress(*chainPublicKey), nil
}
