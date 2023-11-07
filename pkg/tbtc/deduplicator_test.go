package tbtc

import (
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/bitcoin"

	"github.com/keep-network/keep-common/pkg/cache"
)

const testDKGSeedCachePeriod = 1 * time.Second
const testDKGResultHashCachePeriod = 1 * time.Second
const testHeartbeatRequestCachePeriod = 1 * time.Second
const testDepositSweepProposalCachePeriod = 1 * time.Second
const testRedemptionProposalCachePeriod = 1 * time.Second

func TestNotifyDKGStarted(t *testing.T) {
	deduplicator := deduplicator{
		dkgSeedCache: cache.NewTimeCache(testDKGSeedCachePeriod),
	}

	seed1 := big.NewInt(100)
	seed2 := big.NewInt(200)

	// Add the first seed.
	canJoinDKG := deduplicator.notifyDKGStarted(seed1)
	if !canJoinDKG {
		t.Fatal("should be allowed to join DKG")
	}

	// Add the second seed.
	canJoinDKG = deduplicator.notifyDKGStarted(seed2)
	if !canJoinDKG {
		t.Fatal("should be allowed to join DKG")
	}

	// Add the first seed before caching period elapses.
	canJoinDKG = deduplicator.notifyDKGStarted(seed1)
	if canJoinDKG {
		t.Fatal("should not be allowed to join DKG")
	}

	// Wait until caching period elapses.
	time.Sleep(testDKGSeedCachePeriod)

	// Add the first seed again.
	canJoinDKG = deduplicator.notifyDKGStarted(seed1)
	if !canJoinDKG {
		t.Fatal("should be allowed to join DKG")
	}
}

func TestNotifyDKGResultSubmitted(t *testing.T) {
	deduplicator := deduplicator{
		dkgResultHashCache: cache.NewTimeCache(testDKGResultHashCachePeriod),
	}

	hash1Bytes, err := hex.DecodeString("92327ddff69a2b8c7ae787c5d590a2f14586089e6339e942d56e82aa42052cd9")
	if err != nil {
		t.Fatal(err)
	}
	var hash1 [32]byte
	copy(hash1[:], hash1Bytes)

	hash2Bytes, err := hex.DecodeString("23c0062913c4614bdff07f94475ceb4c585df53f71611776c3521ed8f8785913")
	if err != nil {
		t.Fatal(err)
	}
	var hash2 [32]byte
	copy(hash2[:], hash2Bytes)

	// Add the original parameters.
	canProcess := deduplicator.notifyDKGResultSubmitted(big.NewInt(100), hash1, 500)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add with different seed.
	canProcess = deduplicator.notifyDKGResultSubmitted(big.NewInt(101), hash1, 500)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add with different result hash.
	canProcess = deduplicator.notifyDKGResultSubmitted(big.NewInt(100), hash2, 500)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add with different result block.
	canProcess = deduplicator.notifyDKGResultSubmitted(big.NewInt(100), hash1, 501)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add with all different parameters.
	canProcess = deduplicator.notifyDKGResultSubmitted(big.NewInt(101), hash2, 501)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add the original parameters before caching period elapses.
	canProcess = deduplicator.notifyDKGResultSubmitted(big.NewInt(100), hash1, 500)
	if canProcess {
		t.Fatal("should not be allowed to process")
	}

	// Wait until caching period elapses.
	time.Sleep(testDKGResultHashCachePeriod)

	// Add the original parameters again.
	canProcess = deduplicator.notifyDKGResultSubmitted(big.NewInt(100), hash1, 500)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}
}

func TestNotifyHeartbeatRequestSubmitted(t *testing.T) {
	deduplicator := deduplicator{
		heartbeatRequestCache: cache.NewTimeCache(testHeartbeatRequestCachePeriod),
	}

	walletPublicKeyHash1 := [20]byte{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5}
	walletPublicKeyHash2 := [20]byte{2, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 6}
	heartbeatMessage1 := []byte{1, 2}
	heartbeatMessage2 := []byte{3, 4}

	// Original heartbeat request
	canProcess := deduplicator.notifyHeartbeatRequestSubmitted(walletPublicKeyHash1, heartbeatMessage1)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// With another message
	canProcess = deduplicator.notifyHeartbeatRequestSubmitted(walletPublicKeyHash1, heartbeatMessage2)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// With another wallet
	canProcess = deduplicator.notifyHeartbeatRequestSubmitted(walletPublicKeyHash2, heartbeatMessage1)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// With the same wallet and message
	canProcess = deduplicator.notifyHeartbeatRequestSubmitted(walletPublicKeyHash1, heartbeatMessage1)
	if canProcess {
		t.Fatal("should not be allowed to process")
	}

	// Wait until caching period elapses.
	time.Sleep(testHeartbeatRequestCachePeriod)

	// With the same wallet and message again
	canProcess = deduplicator.notifyHeartbeatRequestSubmitted(walletPublicKeyHash1, heartbeatMessage1)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}
}

func TestNotifyDepositSweepProposalSubmitted(t *testing.T) {
	deduplicator := deduplicator{
		depositSweepProposalCache: cache.NewTimeCache(
			testDepositSweepProposalCachePeriod,
		),
	}

	newHash := func(t *testing.T, value string) bitcoin.Hash {
		hash, err := bitcoin.NewHashFromString(value, bitcoin.InternalByteOrder)
		if err != nil {
			t.Fatal(err)
		}

		return hash
	}

	// Original proposal.
	proposal := &DepositSweepProposal{
		WalletPublicKeyHash: [20]byte{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
		DepositsKeys: []struct {
			FundingTxHash      bitcoin.Hash
			FundingOutputIndex uint32
		}{
			{newHash(t, "74d0e353cdba99a6c17ce2cfeab62a26c09b5eb756eccdcfb83dbc12e67b18bc"), 4},
			{newHash(t, "f8eaf242a55ea15e602f9f990e33f67f99dfbe25d1802bbde63cc1caabf99668"), 0},
		},
		SweepTxFee: big.NewInt(1000),
	}

	// Proposal with different wallet.
	proposalDiffWallet := &DepositSweepProposal{
		WalletPublicKeyHash: [20]byte{2, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
		DepositsKeys: []struct {
			FundingTxHash      bitcoin.Hash
			FundingOutputIndex uint32
		}{
			{newHash(t, "74d0e353cdba99a6c17ce2cfeab62a26c09b5eb756eccdcfb83dbc12e67b18bc"), 4},
			{newHash(t, "f8eaf242a55ea15e602f9f990e33f67f99dfbe25d1802bbde63cc1caabf99668"), 0},
		},
		SweepTxFee: big.NewInt(1000),
	}

	// Proposal with different deposits.
	proposalDiffDeposits := &DepositSweepProposal{
		WalletPublicKeyHash: [20]byte{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
		DepositsKeys: []struct {
			FundingTxHash      bitcoin.Hash
			FundingOutputIndex uint32
		}{
			{newHash(t, "84d0e353cdba99a6c17ce2cfeab62a26c09b5eb756eccdcfb83dbc12e67b18bc"), 4},
			{newHash(t, "f8eaf242a55ea15e602f9f990e33f67f99dfbe25d1802bbde63cc1caabf99668"), 0},
		},
		SweepTxFee: big.NewInt(1000),
	}

	// Proposal with same deposits but in different order.
	proposalDiffDepositsOrder := &DepositSweepProposal{
		WalletPublicKeyHash: [20]byte{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
		DepositsKeys: []struct {
			FundingTxHash      bitcoin.Hash
			FundingOutputIndex uint32
		}{
			{newHash(t, "f8eaf242a55ea15e602f9f990e33f67f99dfbe25d1802bbde63cc1caabf99668"), 0},
			{newHash(t, "74d0e353cdba99a6c17ce2cfeab62a26c09b5eb756eccdcfb83dbc12e67b18bc"), 4},
		},
		SweepTxFee: big.NewInt(1000),
	}

	// Proposal with different sweep tx fee.
	proposalDiffSweepTxFee := &DepositSweepProposal{
		WalletPublicKeyHash: [20]byte{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
		DepositsKeys: []struct {
			FundingTxHash      bitcoin.Hash
			FundingOutputIndex uint32
		}{
			{newHash(t, "74d0e353cdba99a6c17ce2cfeab62a26c09b5eb756eccdcfb83dbc12e67b18bc"), 4},
			{newHash(t, "f8eaf242a55ea15e602f9f990e33f67f99dfbe25d1802bbde63cc1caabf99668"), 0},
		},
		SweepTxFee: big.NewInt(1001),
	}

	// Add the original proposal.
	canProcess := deduplicator.notifyDepositSweepProposalSubmitted(proposal)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add the original proposal before caching period elapses.
	canProcess = deduplicator.notifyDepositSweepProposalSubmitted(proposal)
	if canProcess {
		t.Fatal("should not be allowed to process")
	}

	// Add the proposal with different wallet.
	canProcess = deduplicator.notifyDepositSweepProposalSubmitted(proposalDiffWallet)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add the proposal with different deposits.
	canProcess = deduplicator.notifyDepositSweepProposalSubmitted(proposalDiffDeposits)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add the proposal with different deposits order.
	canProcess = deduplicator.notifyDepositSweepProposalSubmitted(proposalDiffDepositsOrder)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add the proposal with different sweep tx fee.
	canProcess = deduplicator.notifyDepositSweepProposalSubmitted(proposalDiffSweepTxFee)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Wait until caching period elapses.
	time.Sleep(testDepositSweepProposalCachePeriod)

	// Add the original proposal again.
	canProcess = deduplicator.notifyDepositSweepProposalSubmitted(proposal)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}
}

func TestNotifyRedemptionProposalSubmitted(t *testing.T) {
	deduplicator := deduplicator{
		redemptionProposalCache: cache.NewTimeCache(
			testRedemptionProposalCachePeriod,
		),
	}

	newScript := func(t *testing.T, value string) bitcoin.Script {
		hash, err := hex.DecodeString(value)
		if err != nil {
			t.Fatal(err)
		}

		return hash
	}

	// Original proposal.
	proposal := &RedemptionProposal{
		WalletPublicKeyHash: [20]byte{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
		RedeemersOutputScripts: []bitcoin.Script{
			newScript(t, "74d0e353cdba99a6c17ce2cfeab62a26c09b5eb756eccdcfb83dbc12e67b18bc"),
			newScript(t, "f8eaf242a55ea15e602f9f990e33f67f99dfbe25d1802bbde63cc1caabf99668"),
		},
		RedemptionTxFee: big.NewInt(1000),
	}

	// Proposal with different wallet.
	proposalDiffWallet := &RedemptionProposal{
		WalletPublicKeyHash: [20]byte{2, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
		RedeemersOutputScripts: []bitcoin.Script{
			newScript(t, "74d0e353cdba99a6c17ce2cfeab62a26c09b5eb756eccdcfb83dbc12e67b18bc"),
			newScript(t, "f8eaf242a55ea15e602f9f990e33f67f99dfbe25d1802bbde63cc1caabf99668"),
		},
		RedemptionTxFee: big.NewInt(1000),
	}

	// Proposal with different redeemer scripts.
	proposalDiffScripts := &RedemptionProposal{
		WalletPublicKeyHash: [20]byte{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
		RedeemersOutputScripts: []bitcoin.Script{
			newScript(t, "84d0e353cdba99a6c17ce2cfeab62a26c09b5eb756eccdcfb83dbc12e67b18bc"),
			newScript(t, "f8eaf242a55ea15e602f9f990e33f67f99dfbe25d1802bbde63cc1caabf99668"),
		},
		RedemptionTxFee: big.NewInt(1000),
	}

	// Proposal with same redeemer scripts but in different order.
	proposalDiffScriptsOrder := &RedemptionProposal{
		WalletPublicKeyHash: [20]byte{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
		RedeemersOutputScripts: []bitcoin.Script{
			newScript(t, "f8eaf242a55ea15e602f9f990e33f67f99dfbe25d1802bbde63cc1caabf99668"),
			newScript(t, "74d0e353cdba99a6c17ce2cfeab62a26c09b5eb756eccdcfb83dbc12e67b18bc"),
		},
		RedemptionTxFee: big.NewInt(1000),
	}

	// Proposal with different redemption tx fee.
	proposalDiffRedemptionTxFee := &RedemptionProposal{
		WalletPublicKeyHash: [20]byte{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
		RedeemersOutputScripts: []bitcoin.Script{
			newScript(t, "74d0e353cdba99a6c17ce2cfeab62a26c09b5eb756eccdcfb83dbc12e67b18bc"),
			newScript(t, "f8eaf242a55ea15e602f9f990e33f67f99dfbe25d1802bbde63cc1caabf99668"),
		},
		RedemptionTxFee: big.NewInt(1001),
	}

	// Add the original proposal.
	canProcess := deduplicator.notifyRedemptionProposalSubmitted(proposal)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add the original proposal before caching period elapses.
	canProcess = deduplicator.notifyRedemptionProposalSubmitted(proposal)
	if canProcess {
		t.Fatal("should not be allowed to process")
	}

	// Add the proposal with different wallet.
	canProcess = deduplicator.notifyRedemptionProposalSubmitted(proposalDiffWallet)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add the proposal with different redeemer scripts.
	canProcess = deduplicator.notifyRedemptionProposalSubmitted(proposalDiffScripts)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add the proposal with different redeemer scripts order.
	canProcess = deduplicator.notifyRedemptionProposalSubmitted(proposalDiffScriptsOrder)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Add the proposal with different redemption tx fee.
	canProcess = deduplicator.notifyRedemptionProposalSubmitted(proposalDiffRedemptionTxFee)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}

	// Wait until caching period elapses.
	time.Sleep(testRedemptionProposalCachePeriod)

	// Add the original proposal again.
	canProcess = deduplicator.notifyRedemptionProposalSubmitted(proposal)
	if !canProcess {
		t.Fatal("should be allowed to process")
	}
}
