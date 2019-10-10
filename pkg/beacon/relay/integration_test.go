// Package relay_test contains integration tests for the whole random beacon
// roundtrip including DKG and threshold signing.
package relay_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/entry"
	"github.com/keep-network/keep-core/pkg/bls"

	"github.com/keep-network/keep-core/pkg/altbn128"

	"github.com/keep-network/keep-core/pkg/internal/dkgtest"
	"github.com/keep-network/keep-core/pkg/internal/entrytest"
	"github.com/keep-network/keep-core/pkg/net"
)

const groupSize = 10
const honestThreshold = 6

var previousEntry, _ = new(big.Int).SetString("132847218974128941824981812", 10)
var seed, _ = new(big.Int).SetString("123789127389127398172398123", 10)

func TestExecute_Bug(t *testing.T) {
	publicKey, err := hex.DecodeString("1bc42f2964cd7e8f5bfa14a8afd35643377d9f8584b9612cad8efb1d42e38c76242a72937cbc434c231e40fa8b48374985415cf95a1a17a4b86cd9fad722d7ed")
	if err != nil {
		t.Fatal(err)
	}

	groupPublicKey, err := altbn128.DecompressToG2(publicKey)
	if err != nil {
		t.Fatal(err)
	}

	share1, _ := new(big.Int).SetString("18741891083689489722858772803753298790305147842591142472576932259609513235297", 10)
	share2, _ := new(big.Int).SetString("10972812093949641174507199649862632899818536049581886008324019953839962045102", 10)
	share3, _ := new(big.Int).SetString("3553471460086587272012285416046512861599745883229014640316048150407697061205", 10)
	share4, _ := new(big.Int).SetString("18372112053939603237620435847562213764197141743948562712251221035888526779223", 10)
	share5, _ := new(big.Int).SetString("11652248131830138626838839453895185430513994830908461536733130237130834207922", 10)

	signer1 := dkg.NewThresholdSigner(1, groupPublicKey, share1)
	signer2 := dkg.NewThresholdSigner(2, groupPublicKey, share2)
	signer3 := dkg.NewThresholdSigner(3, groupPublicKey, share3)
	signer4 := dkg.NewThresholdSigner(4, groupPublicKey, share4)
	signer5 := dkg.NewThresholdSigner(5, groupPublicKey, share5)

	signers := []*dkg.ThresholdSigner{
		signer1, signer2, signer3, signer4, signer5,
	}

	threshold := 3 // threshold from operator contract
	previousEntry, _ := new(big.Int).SetString("18550767884034693229942683654160955525259796233210482098905184301929613300719", 10)
	seed, _ := new(big.Int).SetString("3169316253148515269823481756057378585363364814898551003472951365150675163", 10)

	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		return msg
	}
	signingResult, err := entrytest.RunTest(
		signers,
		threshold,
		interceptor,
		previousEntry,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	signature, err := getSignature(signingResult)
	if err != nil {
		t.Fatal(err)
	}

	entryToSign := entry.CombineToSign(previousEntry, seed)
	if !bls.Verify(groupPublicKey, entryToSign, signature) {
		t.Errorf("threshold signature failed BLS verification")
	}
}

// Success: all members of the signing group participate in signing.
func TestExecute_AllMembersSigning(t *testing.T) {
	t.Parallel()

	signingMembersCount := groupSize
	dkgResult, signingResult := runTest(
		t,
		groupSize,
		honestThreshold,
		signingMembersCount,
	)

	dkgtest.AssertDkgResultPublished(t, dkgResult)
	dkgtest.AssertSamePublicKey(t, dkgResult)
	entrytest.AssertEntryPublished(t, signingResult)
	entrytest.AssertNoSignerFailures(t, signingResult)

	groupPublicKey, err := getFirstGroupPublicKey(dkgResult)
	if err != nil {
		t.Fatal(err)
	}

	signature, err := getSignature(signingResult)
	if err != nil {
		t.Fatal(err)
	}

	entryToSign := entry.CombineToSign(previousEntry, seed)
	if !bls.Verify(groupPublicKey, entryToSign, signature) {
		t.Errorf("threshold signature failed BLS verification")
	}
}

// Success: honest threshold of the signing group members participate in
// signing.
func TestExecute_HonestThresholdMembersSigning(t *testing.T) {
	t.Parallel()

	signingMembersCount := honestThreshold
	dkgResult, signingResult := runTest(
		t,
		groupSize,
		honestThreshold,
		signingMembersCount,
	)

	dkgtest.AssertDkgResultPublished(t, dkgResult)
	dkgtest.AssertSamePublicKey(t, dkgResult)
	entrytest.AssertEntryPublished(t, signingResult)
	entrytest.AssertNoSignerFailures(t, signingResult)

	groupPublicKey, err := getFirstGroupPublicKey(dkgResult)
	if err != nil {
		t.Fatal(err)
	}

	signature, err := getSignature(signingResult)
	if err != nil {
		t.Fatal(err)
	}

	entryToSign := entry.CombineToSign(previousEntry, seed)
	if !bls.Verify(groupPublicKey, entryToSign, signature) {
		t.Errorf("threshold signature failed BLS verification")
	}
}

// Failure: Less than honest threshold signing group members participate in
// signing.
func TestExecute_LessThanHonestThresholdMembersSigning(t *testing.T) {
	t.Parallel()

	signingMembersCount := honestThreshold - 1
	dkgResult, signingResult := runTest(
		t,
		groupSize,
		honestThreshold,
		signingMembersCount,
	)

	dkgtest.AssertDkgResultPublished(t, dkgResult)
	dkgtest.AssertSamePublicKey(t, dkgResult)
	entrytest.AssertEntryNotPublished(t, signingResult)
	entrytest.AssertSignerFailuresCount(t, signingResult, signingMembersCount)
}

func runTest(t *testing.T, groupSize, honestThreshold, honestSignersCount int) (
	*dkgtest.Result,
	*entrytest.Result,
) {
	interceptor := func(msg net.TaggedMarshaler) net.TaggedMarshaler {
		return msg
	}

	dkgSeed := dkgtest.RandomSeed(t)
	dkgResult, err := dkgtest.RunTest(groupSize, honestThreshold, dkgSeed, interceptor)
	if err != nil {
		t.Fatal(err)
	}

	signers := dkgResult.GetSigners()[0:honestSignersCount]

	signingResult, err := entrytest.RunTest(
		signers,
		honestThreshold,
		interceptor,
		previousEntry,
		seed,
	)
	if err != nil {
		t.Fatal(err)
	}

	return dkgResult, signingResult
}

func getSignature(result *entrytest.Result) (*bn256.G1, error) {
	entry := result.EntryValue()
	if entry == nil {
		return nil, fmt.Errorf("no new entry")
	}

	return altbn128.DecompressToG1(entry.Bytes())
}

func getFirstGroupPublicKey(result *dkgtest.Result) (*bn256.G2, error) {
	signers := result.GetSigners()
	if len(signers) == 0 {
		return nil, fmt.Errorf("no signers in result")
	}

	return altbn128.DecompressToG2(signers[0].GroupPublicKeyBytes())
}
