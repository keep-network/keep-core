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
	groupPublicKey := new(bn256.G2).ScalarBaseMult(groupPrivateKey)

	threshold := 6
	signers := []*ThresholdSigner{
		&ThresholdSigner{1, groupPublicKey, big.NewInt(2)},
		&ThresholdSigner{2, groupPublicKey, big.NewInt(12)},
		&ThresholdSigner{3, groupPublicKey, big.NewInt(10)},
		&ThresholdSigner{4, groupPublicKey, big.NewInt(5)},
		&ThresholdSigner{5, groupPublicKey, big.NewInt(4)},
		&ThresholdSigner{6, groupPublicKey, big.NewInt(1)},
	} // 2 + 12 + 10 + 5 + 4 + 1 = 34

	message := []byte("hello world")

	shares := make([]*bls.SignatureShare, 0)
	for i, signer := range signers {
		shares = append(shares,
			&bls.SignatureShare{
				I: i + 1, // always 1-indexed
				V: signer.CalculateSignatureShare(message),
			},
		)
	}

	signature, err := signers[0].CompleteSignature(shares, threshold)
	if err != nil {
		t.Fatal(err)
	}
	actual := signature.Marshal()
	expected := bls.Sign(groupPrivateKey, message).Marshal()

	testutils.AssertBytesEqual(t, expected, actual)
}

// 3 tests - 1 where we have 0 that work, 1 where we have under t, one where we have t
