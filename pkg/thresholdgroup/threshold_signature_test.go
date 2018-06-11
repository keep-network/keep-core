package thresholdgroup

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/dfinity/go-dfinity-crypto/bls"
)

const (
	defaultSigningThreshold = defaultDishonestThreshold + 1
)

func hasNonZeroBytes(bytes []byte) bool {
	hasNonZeroBytes := false
	for _, bte := range bytes {
		if bte != byte(0) {
			hasNonZeroBytes = true
			break
		}
	}

	return hasNonZeroBytes
}

func TestMemberExposesGroupPublicKey(t *testing.T) {
	member, _ := buildMembers("")

	bytes := member.GroupPublicKeyBytes()
	if !hasNonZeroBytes(bytes[:]) {
		t.Errorf(
			"\nexpected: nonzero bytes in [%v]\nactual:   no nonzero bytes found",
			bytes,
		)
	}
}

func TestMembersHaveSameGroupPublicKey(t *testing.T) {
	member, otherMembers := buildMembers("")

	publicKey := member.GroupPublicKeyBytes()
	for _, otherMember := range otherMembers {
		memberKey := otherMember.GroupPublicKeyBytes()
		if !reflect.DeepEqual(publicKey, memberKey) {
			t.Errorf(
				"for id %v\nexpected: [%v]\nactual:   [%v]",
				otherMember.MemberID(),
				publicKey,
				memberKey,
			)
		}
	}
}

func TestMembersProduceSignatureShare(t *testing.T) {
	member, otherMembers := buildMembers("")
	allMembers := append(otherMembers, member)

	message := fmt.Sprintf("%v", rand.Int63())

	allShares := make([][]byte, 0)
	for _, member := range allMembers {
		allShares = append(allShares, member.SignatureShare(message))
	}

	for i, share := range allShares {
		if !hasNonZeroBytes(share) {
			t.Errorf(
				"at index %v\nexpected: nonzero bytes in [%v]\nactual:   no nonzero bytes found",
				i,
				share,
			)
		}
	}
}

func TestMemberProducesSignatureFromShares(t *testing.T) {
	var tests = map[string]struct {
		participatingMembers int
		expectedVerification bool
		expectedError        error
	}{
		"with all members participating": {
			participatingMembers: defaultGroupSize,
			expectedVerification: true,
			expectedError:        nil,
		},
		"with more than a sign threshold participating": {
			participatingMembers: defaultSigningThreshold + 1,
			expectedVerification: true,
			expectedError:        nil,
		},
		"with a sign threshold participating": {
			participatingMembers: defaultSigningThreshold,
			expectedVerification: true,
			expectedError:        nil,
		},
		"with less than a sign threshold participating": {
			participatingMembers: defaultSigningThreshold - 1,
			expectedVerification: false,
			expectedError: fmt.Errorf(
				"%v shares are insufficient for a complete signature; need %v",
				defaultSigningThreshold-1,
				defaultSigningThreshold,
			),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			member, otherMembers := buildMembers("")
			allMembers := append(otherMembers, member)

			message := fmt.Sprintf("%v", rand.Int63())

			shares := make(map[bls.ID][]byte)
			for _, i := range rand.Perm(test.participatingMembers) {
				randomMember := allMembers[i]
				shares[randomMember.BlsID] = randomMember.SignatureShare(message)
			}

			signature, err := member.CompleteSignature(shares)
			if !reflect.DeepEqual(err, test.expectedError) {
				t.Fatalf(
					"\nexpected: %v\nactual:   %v",
					test.expectedError,
					err,
				)
			}
			verification := false
			// Verification *will* segfault for certain invalid signatures.
			if err == nil {
				verification = signature.Verify(member.groupPublicKey, message)
			}

			actualText := "verified signature"
			if !verification {
				actualText = "unverified signature"
			}
			expectedText := "verified signature"
			if !test.expectedVerification {
				expectedText = "unverified signature"
			}

			if verification != test.expectedVerification {
				t.Fatalf(
					"\nexpected: %v\nactual:   %v",
					expectedText,
					actualText,
				)
			}
		})
	}
}

func TestMemberVerifiesSignatureFromShares(t *testing.T) {
	var tests = map[string]struct {
		participatingMembers int
		expectedVerification bool
		expectedError        error
	}{
		"with all members participating": {
			participatingMembers: defaultGroupSize,
			expectedVerification: true,
			expectedError:        nil,
		},
		"with more than a sign threshold participating": {
			participatingMembers: defaultSigningThreshold + 1,
			expectedVerification: true,
			expectedError:        nil,
		},
		"with a sign threshold participating": {
			participatingMembers: defaultSigningThreshold,
			expectedVerification: true,
			expectedError:        nil,
		},
		"with less than a sign threshold participating": {
			participatingMembers: defaultSigningThreshold - 1,
			expectedVerification: false,
			expectedError: fmt.Errorf(
				"%v shares are insufficient for a complete signature; need %v",
				defaultSigningThreshold-1,
				defaultSigningThreshold,
			),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			member, otherMembers := buildMembers("")
			allMembers := append(otherMembers, member)

			message := fmt.Sprintf("%v", rand.Int63())

			shares := make(map[bls.ID][]byte)
			for _, i := range rand.Perm(test.participatingMembers) {
				randomMember := allMembers[i]
				shares[randomMember.BlsID] = randomMember.SignatureShare(message)
			}

			verification, err := member.VerifySignature(shares, message)
			if !reflect.DeepEqual(err, test.expectedError) {
				t.Fatalf(
					"\nexpected: %v\nactual:   %v",
					test.expectedError,
					err,
				)
			}

			actualText := "verified signature"
			if !verification {
				actualText = "unverified signature"
			}
			expectedText := "verified signature"
			if !test.expectedVerification {
				expectedText = "unverified signature"
			}

			if verification != test.expectedVerification {
				t.Fatalf(
					"\nexpected: %v\nactual:   %v",
					expectedText,
					actualText,
				)
			}
		})
	}
}
