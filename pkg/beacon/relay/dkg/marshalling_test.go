package dkg

import (
	"reflect"
	"testing"

	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
)

func TestThresholdSignerRoundtrip(t *testing.T) {
	thresholdSigner := &ThresholdSigner{
		memberIndex:          group.MemberIndex(2),
		groupPublicKey:       new(bn256.G2).ScalarBaseMult(big.NewInt(10)),
		groupPrivateKeyShare: big.NewInt(1),
	}

	unmarshaled := &ThresholdSigner{}

	err := pbutils.RoundTrip(thresholdSigner, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(thresholdSigner, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled threshold signer")
	}
}
