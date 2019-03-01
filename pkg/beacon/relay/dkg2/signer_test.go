package dkg2

import (
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/bls"
)

func TestSignAndComplete(t *testing.T) {
	threshold := 6

	// Obtained by running `TestFullStateTransitions` and outputting shares.
	privateKeySharesSlice := []string{
		"7965280207209549879164292761852524476109477664957641865927295346590476704711",
		"6106610144639464785158072029008498287824734372346580964957055618768317731307",
		"1440342545552619026306193227443878557590393893287770611896928291936003327327",
		"4784636576436291513911847371324649154255552173376284038648549444705786002507",
		"11155341033982526633651061739311639943498023031242960281134125600620928234282",
		"476414470139165825449834970450656200990319293028518982997364506493673497757",
	}

	// First get SecretKeyShares from slice of privateKeyShares
	var publicKeyShares []*bls.PublicKeyShare
	for i, privateKeyShareString := range privateKeySharesSlice {
		privateKeyShare, _ := new(big.Int).SetString(privateKeyShareString, 10)
		publicKeyShare := (&bls.SecretKeyShare{
			I: i + 1,
			V: privateKeyShare,
		}).PublicKeyShare()
		publicKeyShares = append(publicKeyShares, publicKeyShare)
	}
	// Build up the group public key
	groupPublicKey, err := bls.RecoverPublicKey(publicKeyShares, threshold)
	if err != nil {
		t.Fatal(err)
	}

	var signers []*ThresholdSigner
	for i, privateKeyShare := range privateKeySharesSlice {
		share, _ := new(big.Int).SetString(privateKeyShare, 10)
		signers = append(signers, &ThresholdSigner{
			memberID:             gjkr.MemberID(i + 1),
			groupPublicKey:       groupPublicKey,
			groupPrivateKeyShare: share,
		})
	}

	message := []byte("hello world")
	shares := make([]*bls.SignatureShare, 0)
	for _, signer := range signers {
		shares = append(shares,
			&bls.SignatureShare{
				I: int(signer.MemberID()),
				V: signer.CalculateSignatureShare(message),
			},
		)
	}

	signature, err := signers[0].CompleteSignature(shares, threshold)
	if err != nil {
		t.Fatal(err)
	}

	if !bls.Verify(groupPublicKey, message, signature) {
		t.Fatal("Failed to verify recovered signature")
	}
}

// 3 tests - 1 where we have 0 that work, 1 where we have under t, one where we have t
