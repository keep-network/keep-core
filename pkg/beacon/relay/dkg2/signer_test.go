package dkg2

import (
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/bls"
)

func TestSignAndComplete(t *testing.T) {
	var message = []byte("hello world")

	// Obtained by running `TestFullStateTransitions` and outputting shares
	// and group public key. MemberIDs are 1-indexed.
	privateKeySharesMap := map[int]string{
		1: "+20447821705176695776117400920440893381372259028396365458583014272617533574429",
		2: "+10311498259490277582707403215942210669382166384656845373229012913757750213620",
		3: "+12931471259504366138666739996593106353126621511383680527266384358924878290714",
		4: "+6419497833379686221749005517136305344057260008160836576996924421543109310094",
		5: "+12663820852955513054200605522829082730722446275404347866118837288188251767377",
		6: "+9776197446392571413775134268414163424573815912698180050933918772284497166946",
	}

	groupPublicKeyBytes := []byte{
		16, 225, 37, 168, 24, 49, 229, 90, 189, 2, 116, 144, 153, 193,
		13, 16, 145, 179, 12, 149, 188, 143, 204, 187, 26, 234, 97, 64,
		220, 224, 79, 47, 7, 96, 34, 99, 78, 229, 11, 105, 226, 224,
		190, 36, 93, 101, 69, 59, 77, 214, 30, 38, 28, 32, 14, 119, 222,
		91, 179, 111, 184, 157, 166, 29, 23, 175, 226, 54, 240, 195, 237,
		93, 222, 59, 74, 47, 49, 0, 67, 145, 70, 41, 172, 45, 114, 43, 3,
		125, 247, 77, 208, 176, 240, 31, 240, 231, 20, 114, 77, 45, 177,
		55, 59, 116, 81, 226, 108, 253, 63, 53, 27, 30, 24, 53, 88, 219,
		81, 62, 155, 65, 94, 209, 138, 210, 225, 21, 51, 192,
	}
	groupPublicKey := new(bn256.G2)
	_, err := groupPublicKey.Unmarshal(groupPublicKeyBytes)
	if err != nil {
		t.Fatal(err)
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
			expectedError:          "not enough shares to reconstruct signature",
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

		// First, get SecretKeyShares from slice of privateKeyShares.
		var publicKeyShares []*bls.PublicKeyShare
		for memberID, privateKeyShareString := range privateKeyShares {
			privateKeyShare, _ := new(big.Int).SetString(privateKeyShareString, 10)
			publicKeyShare := (&bls.SecretKeyShare{
				I: memberID,
				V: privateKeyShare,
			}).PublicKeyShare()
			publicKeyShares = append(publicKeyShares, publicKeyShare)
		}

		// Next, build up signers from public key shares.
		var signers []*ThresholdSigner
		for memberID, privateKeyShare := range privateKeyShares {
			share, _ := new(big.Int).SetString(privateKeyShare, 10)
			signers = append(signers, &ThresholdSigner{
				memberID:             gjkr.MemberID(memberID),
				groupPublicKey:       groupPublicKey,
				groupPrivateKeyShare: share,
			})
		}

		// Ensure we get a valid signature share from every signer.
		shares := make([]*bls.SignatureShare, 0)
		for _, signer := range signers {
			shares = append(shares,
				&bls.SignatureShare{
					I: int(signer.MemberID()),
					V: signer.CalculateSignatureShare(message),
				},
			)
		}

		// Attempt to recover a signature from the present shares.
		signature, err := signers[0].CompleteSignature(shares, test.threshold)
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

		// Does the signature match the public key that we have for the group?
		if !bls.Verify(groupPublicKey, message, signature) {
			t.Fatal("Failed to verify recovered signature")
		}
	}
}
