package ethereum

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-core/pkg/operator"
)

type signer struct {
	*ethutil.EthereumSigner
}

func newSigner(chainKey *keystore.Key) *signer {
	return &signer{
		ethutil.NewSigner(chainKey.PrivateKey),
	}
}

func (s *signer) PublicKeyToAddress(
	publicKey *operator.PublicKey,
) ([]byte, error) {
	chainPublicKey, err := OperatorPublicKeyToChainPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot convert operator key to chain key: [%v]",
			err,
		)
	}

	return s.EthereumSigner.PublicKeyToAddress(*chainPublicKey), nil
}