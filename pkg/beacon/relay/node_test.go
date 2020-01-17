package relay

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

func TestCreateGroupMemberFilter(t *testing.T) {
	signing := &mockSigning{}

	stakersPublicKeys := make([]ecdsa.PublicKey, 5)
	for i := range stakersPublicKeys {
		privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		stakersPublicKeys[i] = privateKey.PublicKey
	}

	stakersAddresses := make([]relaychain.StakerAddress, len(stakersPublicKeys))
	for i := range stakersAddresses {
		stakersAddresses[i] = signing.PublicKeyToAddress(stakersPublicKeys[i])
	}

	// Allow only stakers with index 1 and 2
	filter := createGroupMemberFilter(stakersAddresses[1:3], signing)

	expectedResults := []bool{false, true, true, false, false}
	for i, stakerPublicKey := range stakersPublicKeys {
		actualResult := filter(&stakerPublicKey)

		if expectedResults[i] != actualResult {
			t.Errorf(
				"Unexpected result for staker index [%v]\n"+
					"Expected: %v\nActual:   %v\n",
				i,
				expectedResults[i],
				actualResult,
			)
		}
	}
}

type mockSigning struct{}

func (ms *mockSigning) PublicKey() []byte {
	panic("not implemented")
}

func (ms *mockSigning) Sign(message []byte) ([]byte, error) {
	panic("not implemented")
}

func (ms *mockSigning) Verify(message []byte, signature []byte) (bool, error) {
	panic("not implemented")
}

func (ms *mockSigning) VerifyWithPublicKey(
	message []byte,
	signature []byte,
	publicKey []byte,
) (bool, error) {
	panic("not implemented")
}

func (ms *mockSigning) PublicKeyToAddress(publicKey ecdsa.PublicKey) []byte {
	return elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
}
