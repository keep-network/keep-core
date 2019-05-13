package registry

import (
	"reflect"
	"testing"

	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
)

func TestDKGResultHashSignatureMessageRoundtrip(t *testing.T) {
	signer := dkg.NewThresholdSigner(
		group.MemberIndex(2),
		new(bn256.G2).ScalarBaseMult(big.NewInt(10)),
		big.NewInt(1),
	)

	membershipMessage := &Membership{
		Signer: signer,
	}

	unmarshaled := &Membership{}

	err := pbutils.RoundTrip(membershipMessage, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(membershipMessage, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}
