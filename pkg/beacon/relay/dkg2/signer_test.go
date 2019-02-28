package dkg2

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/bls"
)

func TestSignAndComplete(t *testing.T) {
	threshold := 6

	// Obtained by running `TestFullStateTransitions` and outputting shares.

	// groupPublicKey: bn256.G2((0b9a57387a8e2c5f588ddece242ab9aa005eea7a2edd93c2890d6e6b6a4a379c, 23023e11f05d3c4694f5fb22960fa6e9448cd7d959695f2ca11b8fe737dc03b3), (17c10becb92975a8a9d0fd5bb4d69c964c69cfae94694db45ed4b14d9ba8299b, 1b9cfe419acfb244d95bd71b2c8453859e75437249c6ab085a4a5c8c920f0690))

	privateKeySharesSlice := []string{
		"11317689322678074320253490110336291778019389544269451057844495958842587881889",
		"3321938026855708852034359825983988567719999295358659054897559573286783870881",
		"16637974750087405446600293454316416930247455842507609024774997421245084242640",
		"8811658855469363004088702363302325849186372736438301169002749756782287688357",
		"5216737455961088866926156556557357315061484853999540690387272503653832520818",
		"19745272718386738224116191829677947371012382482669271992696752605818677063650",
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
	fmt.Println(groupPublicKey)

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

	if !bls.Verify(groupPublicKey, message, signature) {
		t.Fatal("Failed to verify recovered signature")
	}
}

// 3 tests - 1 where we have 0 that work, 1 where we have under t, one where we have t
