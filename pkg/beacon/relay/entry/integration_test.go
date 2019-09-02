package entry_test

import (
	"encoding/hex"
	"math/big"
	"math/rand"
	"testing"

	"github.com/keep-network/keep-core/pkg/altbn128"

	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/entry"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/internal/entrytest"
	"github.com/keep-network/keep-core/pkg/net"
)

// all signers can reconstruct signature having all the shares
func TestExecute_HappyPath(t *testing.T) {
	t.Parallel()

	signers := initSigners(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		return msg
	}

	previousEntry := big.NewInt(rand.Int63())
	seed := big.NewInt(rand.Int63())

	result, err := entrytest.RunTest(
		signers,
		threshold,
		interceptor,
		previousEntry,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	entrytest.AssertEntryPublished(t, result)
	entrytest.AssertSignerFailuresCount(t, result, 0)
}

// signer 1 can reconstruct signature - has shares from [1, 3, 4, 5]
// signer 2 can reconstruct signature - has shares from [1, 2, 3, 4, 5]
// signer 3 can reconstruct signature - has shares from [1, 3, 4, 5]
// signer 4 can reconstruct signature - has shares from [1, 3, 4, 5]
// signer 5 can reconstruct signature - has shares from [1, 3, 4, 5]
func TestExecuteIA_signer2(t *testing.T) {
	t.Parallel()

	signers := initSigners(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		signatureShareMessage, ok := msg.(*entry.SignatureShareMessage)
		if ok && signatureShareMessage.SenderID() == group.MemberIndex(2) {
			return nil
		}

		return msg
	}

	previousEntry := big.NewInt(rand.Int63())
	seed := big.NewInt(rand.Int63())

	result, err := entrytest.RunTest(
		signers,
		threshold,
		interceptor,
		previousEntry,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	entrytest.AssertEntryPublished(t, result)
	entrytest.AssertSignerFailuresCount(t, result, 0)
}

// signer 1 can reconstruct signature - has shares from [1, 4, 5]
// signer 2 can reconstruct signature - has shares from [1, 2, 4, 5]
// signer 3 can reconstruct signature - has shares from [1, 3, 4, 5]
// signer 4 can reconstruct signature - has shares from [1, 4, 5]
// signer 5 can reconstruct signature - has shares from [1, 4, 5]
func TestExecuteIA_signers23(t *testing.T) {
	t.Parallel()

	signers := initSigners(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		signatureShareMessage, ok := msg.(*entry.SignatureShareMessage)
		if ok && (signatureShareMessage.SenderID() == group.MemberIndex(2) ||
			signatureShareMessage.SenderID() == group.MemberIndex(3)) {
			return nil
		}

		return msg
	}

	previousEntry := big.NewInt(rand.Int63())
	seed := big.NewInt(rand.Int63())

	result, err := entrytest.RunTest(
		signers,
		threshold,
		interceptor,
		previousEntry,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	entrytest.AssertEntryPublished(t, result)
	entrytest.AssertSignerFailuresCount(t, result, 0)
}

// signer 1 can't reconstruct signature - has shares from [1, 5]
// signer 2 can reconstruct signature - has shares from [1, 2, 5]
// signer 3 can reconstruct signature - has shares from [1, 3, 5]
// signer 4 can reconstruct signature - has shares from [1, 4, 5]
// signer 5 can't reconstruct signature - has shares from [1, 5]
func TestExecuteIA_signers234(t *testing.T) {
	t.Parallel()

	signers := initSigners(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		signatureShareMessage, ok := msg.(*entry.SignatureShareMessage)
		if ok && (signatureShareMessage.SenderID() == group.MemberIndex(2) ||
			signatureShareMessage.SenderID() == group.MemberIndex(3) ||
			signatureShareMessage.SenderID() == group.MemberIndex(4)) {
			return nil
		}

		return msg
	}

	previousEntry := big.NewInt(rand.Int63())
	seed := big.NewInt(rand.Int63())

	result, err := entrytest.RunTest(
		signers,
		threshold,
		interceptor,
		previousEntry,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	entrytest.AssertEntryPublished(t, result)
	entrytest.AssertSignerFailuresCount(t, result, 2)
}

// signer 1 can't reconstruct signature - has shares from [1]
// signer 2 can't reconstruct signature - has shares from [1, 2]
// signer 3 can't reconstruct signature - has shares from [1, 3]
// signer 4 can't reconstruct signature - has shares from [1, 4]
// signer 5 can't reconstruct signature - has shares from [1, 5]
func TestExecuteIA_signers2345(t *testing.T) {
	t.Parallel()

	signers := initSigners(t)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		signatureShareMessage, ok := msg.(*entry.SignatureShareMessage)
		if ok && (signatureShareMessage.SenderID() == group.MemberIndex(2) ||
			signatureShareMessage.SenderID() == group.MemberIndex(3) ||
			signatureShareMessage.SenderID() == group.MemberIndex(4) ||
			signatureShareMessage.SenderID() == group.MemberIndex(5)) {
			return nil
		}

		return msg
	}

	previousEntry := big.NewInt(rand.Int63())
	seed := big.NewInt(rand.Int63())

	result, err := entrytest.RunTest(
		signers,
		threshold,
		interceptor,
		previousEntry,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	entrytest.AssertEntryNotPublished(t, result)
	entrytest.AssertSignerFailuresCount(t, result, groupSize)
}

const groupSize = 5
const threshold = 3

// group of 5 signers created with DKG for threshold = 3
func initSigners(t *testing.T) []*dkg.ThresholdSigner {
	signer1KeyShare, _ := new(big.Int).SetString("19861193134483177941115785550115091961929449607192654978608847685939283615757", 10)
	signer2KeyShare, _ := new(big.Int).SetString("3738797326965009616916531092397687518539287046270322726916185127422400804429", 10)
	signer3KeyShare, _ := new(big.Int).SetString("14915476670403082790011838544775487204401638936796576115323203991337974542037", 10)
	signer4KeyShare, _ := new(big.Int).SetString("12244086356984772336435377139409119114901280060088038611921373199822435295306", 10)
	signer5KeyShare, _ := new(big.Int).SetString("20242210194415278798960033344231036611068078398709436715896774234739639017812", 10)

	groupPublicKeyBytes, err := hex.DecodeString(
		"064514d4fe5de5512c40c5b19ba7995ac86fce10931bdbdde716b3b913acd54710455ad917dce5db6b3764563bba77d055caa656737e983867fdf539c276530f",
	)
	if err != nil {
		t.Fatal(err)
	}

	groupPublicKey, err := altbn128.DecompressToG2(groupPublicKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	signer1 := dkg.NewThresholdSigner(1, groupPublicKey, signer1KeyShare)
	signer2 := dkg.NewThresholdSigner(2, groupPublicKey, signer2KeyShare)
	signer3 := dkg.NewThresholdSigner(3, groupPublicKey, signer3KeyShare)
	signer4 := dkg.NewThresholdSigner(4, groupPublicKey, signer4KeyShare)
	signer5 := dkg.NewThresholdSigner(5, groupPublicKey, signer5KeyShare)

	return []*dkg.ThresholdSigner{signer1, signer2, signer3, signer4, signer5}
}
