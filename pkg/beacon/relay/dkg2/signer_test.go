package dkg2

import (
	"fmt"
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/bls"
)

func TestSignAndComplete(t *testing.T) {
	groupPublicKey := new(bn256.G1).ScalarBaseMult(big.NewInt(34))

	signers := []*ThresholdSigner{
		&ThresholdSigner{1, groupPublicKey, big.NewInt(2)},
		&ThresholdSigner{1, groupPublicKey, big.NewInt(12)},
		&ThresholdSigner{1, groupPublicKey, big.NewInt(10)},
		&ThresholdSigner{1, groupPublicKey, big.NewInt(5)},
		&ThresholdSigner{1, groupPublicKey, big.NewInt(4)},
		&ThresholdSigner{1, groupPublicKey, big.NewInt(1)},
	} // 2 + 12 + 10 + 5 + 4 + 1 = 34

	message := []byte("hello world")

	shares := make([]*bn256.G1, len(signers))
	for i, signer := range signers {
		shares[i] = signer.SignatureShare(message)
	}

	complete := signers[0].CompleteSignature(shares)

	fmt.Printf("actual = [%v]\n", complete)
	fmt.Printf("expected = [%v]\n", bls.Sign(big.NewInt(34), message))

	// TODO: use byteutils to compare
}
