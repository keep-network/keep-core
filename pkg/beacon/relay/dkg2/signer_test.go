package dkg2

import (
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/bls"
)

func TestSignAndComplete(t *testing.T) {
	var message = []byte("hello world")

	// Obtained by running `TestFullStateTransitions` and outputting shares.
	// MemberIDs are 1-indexed.
	privateKeySharesMap := map[int]string{
		1: "7965280207209549879164292761852524476109477664957641865927295346590476704711",
		2: "6106610144639464785158072029008498287824734372346580964957055618768317731307",
		3: "1440342545552619026306193227443878557590393893287770611896928291936003327327",
		4: "4784636576436291513911847371324649154255552173376284038648549444705786002507",
		5: "11155341033982526633651061739311639943498023031242960281134125600620928234282",
		6: "476414470139165825449834970450656200990319293028518982997364506493673497757",
	}

	var tests = map[string]struct {
		threshold              int
		numberPrivateKeyShares int
		expectedError          string
	}{
		"success: all members sign the message": {
			threshold:              6,
			numberPrivateKeyShares: 6,
			expectedError:          "",
		},
		"success: at least t members sign the message": {
			threshold:              3,
			numberPrivateKeyShares: 4,
			expectedError:          "",
		},
		"failure: less than t members sign a message": {
			threshold:              4,
			numberPrivateKeyShares: 3,
			expectedError:          "not enough shares to reconstruct public key",
		},
	}

	for _, test := range tests {
		privateKeyShares := make(map[int]string)
		for memberID, share := range privateKeySharesMap {
			if len(privateKeyShares) == test.numberPrivateKeyShares {
				break
			}
			privateKeyShares[memberID] = share
		}
		// First get SecretKeyShares from slice of privateKeyShares
		var publicKeyShares []*bls.PublicKeyShare
		for memberID, privateKeyShareString := range privateKeyShares {
			privateKeyShare, _ := new(big.Int).SetString(privateKeyShareString, 10)
			publicKeyShare := (&bls.SecretKeyShare{
				I: memberID,
				V: privateKeyShare,
			}).PublicKeyShare()
			publicKeyShares = append(publicKeyShares, publicKeyShare)
		}
		// Build up the group public key
		groupPublicKey, err := bls.RecoverPublicKey(publicKeyShares, test.threshold)
		if err != nil {
			if err.Error() != test.expectedError {
				t.Errorf(
					"\nexpected: %v\nactual:   %v",
					test.expectedError,
					err,
				)
			}
			// exit the test as we errored correctly
			continue
		}

		var signers []*ThresholdSigner
		for memberID, privateKeyShare := range privateKeyShares {
			share, _ := new(big.Int).SetString(privateKeyShare, 10)
			signers = append(signers, &ThresholdSigner{
				memberID:             gjkr.MemberID(memberID),
				groupPublicKey:       groupPublicKey,
				groupPrivateKeyShare: share,
			})
		}

		shares := make([]*bls.SignatureShare, 0)
		for _, signer := range signers {
			shares = append(shares,
				&bls.SignatureShare{
					I: int(signer.MemberID()),
					V: signer.CalculateSignatureShare(message),
				},
			)
		}

		signature, err := signers[0].CompleteSignature(shares, test.threshold)
		if err != nil {
			t.Fatal(err)
		}

		if !bls.Verify(groupPublicKey, message, signature) {
			t.Fatal("Failed to verify recovered signature")
		}
	}
}
