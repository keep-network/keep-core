package dkg2

import (
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/bls"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestSignAndComplete(t *testing.T) {
	groupPrivateKey := big.NewInt(34)
	groupPublicKey := new(bn256.G1).ScalarBaseMult(groupPrivateKey)

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
		shares[i] = signer.CalculateSignatureShare(message)
	}

	actual := signers[0].CompleteSignature(shares).Marshal()
	expected := bls.Sign(groupPrivateKey, message).Marshal()

	testutils.AssertBytesEqual(t, expected, actual)
}
