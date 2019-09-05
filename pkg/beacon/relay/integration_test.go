package relay_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/entry"

	"github.com/keep-network/keep-core/pkg/altbn128"
	"github.com/keep-network/keep-core/pkg/bls"

	"github.com/keep-network/keep-core/pkg/internal/dkgtest"
	"github.com/keep-network/keep-core/pkg/internal/entrytest"
	"github.com/keep-network/keep-core/pkg/net"
)

func TestExecute_HappyPath(t *testing.T) {
	t.Parallel()

	groupSize := 5
	threshold := 3

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		return msg
	}

	dkgResult, err := dkgtest.RunTest(groupSize, threshold-1, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	groupPublicKey, err := altbn128.DecompressToG2(dkgResult.Signers[0].GroupPublicKeyBytes())
	if err != nil {
		t.Fatal(err)
	}

	previousEntry := big.NewInt(rand.Int63())
	seed := big.NewInt(rand.Int63())

	signingResult, err := entrytest.RunTest(
		dkgResult.Signers,
		threshold,
		interceptor,
		previousEntry,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	entryToSign := entry.CombineToSign(previousEntry, seed)
	signature, err := altbn128.DecompressToG1(signingResult.Entry.Value.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Signature I will check: [%v]\n", signature.String())

	fmt.Printf("[%v]\n", signingResult.Entry)
	fmt.Printf("[%v]\n", bls.Verify(groupPublicKey, entryToSign, signature))
}
