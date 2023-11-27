package tbtc

import (
	"context"
	"encoding/hex"
	"github.com/keep-network/keep-core/pkg/chain"
	"testing"
	"time"

	"github.com/keep-network/keep-core/internal/testutils"
)

func TestCoordinationWindow_ActivePhaseEndBlock(t *testing.T) {
	window := newCoordinationWindow(900)

	testutils.AssertIntsEqual(
		t,
		"active phase end block",
		980,
		int(window.activePhaseEndBlock()),
	)
}

func TestCoordinationWindow_EndBlock(t *testing.T) {
	window := newCoordinationWindow(900)

	testutils.AssertIntsEqual(
		t,
		"end block",
		1000,
		int(window.endBlock()),
	)
}

func TestCoordinationWindow_IsAfterActivePhase(t *testing.T) {
	window := newCoordinationWindow(1800)

	previousWindow := newCoordinationWindow(900)
	sameWindow := newCoordinationWindow(1800)
	nextWindow := newCoordinationWindow(2700)

	testutils.AssertBoolsEqual(
		t,
		"result for nil",
		true,
		window.isAfter(nil),
	)
	testutils.AssertBoolsEqual(
		t,
		"result for previous window",
		true,
		window.isAfter(previousWindow),
	)
	testutils.AssertBoolsEqual(
		t,
		"result for same window",
		false,
		window.isAfter(sameWindow),
	)
	testutils.AssertBoolsEqual(
		t,
		"result for next window",
		false,
		window.isAfter(nextWindow),
	)
}

func TestWatchCoordinationWindows(t *testing.T) {
	watchBlocksFn := func(ctx context.Context) <-chan uint64 {
		blocksChan := make(chan uint64)

		go func() {
			ticker := time.NewTicker(1 * time.Millisecond)
			defer ticker.Stop()

			block := uint64(0)

			for {
				select {
				case <-ticker.C:
					block++
					blocksChan <- block
				case <-ctx.Done():
					return
				}
			}
		}()

		return blocksChan
	}

	receivedWindows := make([]*coordinationWindow, 0)
	onWindowFn := func(window *coordinationWindow) {
		receivedWindows = append(receivedWindows, window)
	}

	ctx, cancelCtx := context.WithTimeout(
		context.Background(),
		2000*time.Millisecond,
	)
	defer cancelCtx()

	go watchCoordinationWindows(ctx, watchBlocksFn, onWindowFn)

	<-ctx.Done()

	testutils.AssertIntsEqual(t, "received windows", 2, len(receivedWindows))
	testutils.AssertIntsEqual(
		t,
		"first window",
		900,
		int(receivedWindows[0].coordinationBlock),
	)
	testutils.AssertIntsEqual(
		t,
		"second window",
		1800,
		int(receivedWindows[1].coordinationBlock),
	)
}

func TestCoordinationExecutor_CoordinationSeed(t *testing.T) {
	window := newCoordinationWindow(900)

	localChain := Connect()

	localChain.setBlockHashByNumber(
		window.coordinationBlock-32,
		"1322996cbcbc38fc924a46f4df5f9064279d3ab43396e58386dac9b87440d64f",
	)

	// Uncompressed public key corresponding to the 20-byte public key hash:
	// aa768412ceed10bd423c025542ca90071f9fb62d.
	publicKeyHex, err := hex.DecodeString(
		"0471e30bca60f6548d7b42582a478ea37ada63b402af7b3ddd57f0c95bb6843175" +
			"aa0d2053a91a050a6797d85c38f2909cb7027f2344a01986aa2f9f8ca7a0c289",
	)
	if err != nil {
		t.Fatal(err)
	}

	coordinatedWallet := wallet{
		// Set only relevant fields.
		publicKey: unmarshalPublicKey(publicKeyHex),
	}

	executor := &coordinationExecutor{
		// Set only relevant fields.
		chain:             localChain,
		coordinatedWallet: coordinatedWallet,
	}

	seed, err := executor.coordinationSeed(window)
	if err != nil {
		t.Fatal(err)
	}

	// Expected seed is sha256(wallet_public_key_hash | safe_block_hash).
	expectedSeed := "e55c779d6d83183409ddc90c6cd5130567f0593349a9c82494b402048ec2d03d"

	testutils.AssertStringsEqual(
		t,
		"coordination seed",
		expectedSeed,
		hex.EncodeToString(seed[:]),
	)
}

func TestCoordinationExecutor_CoordinationLeader(t *testing.T) {
	seedBytes, err := hex.DecodeString(
		"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
	)
	if err != nil {
		t.Fatal(err)
	}

	var seed [32]byte
	copy(seed[:], seedBytes)

	coordinatedWallet := wallet{
		// Set only relevant fields.
		signingGroupOperators: []chain.Address{
			"957ECF59507a6A74b8d98747f07a74De270D3CC3", // member 1
			"5E14c0f27612fbfB7A6FE40b5A6Ec997fA62fc04", // member 2
			"D2662604f8b4540336fBd3c1F48d7e9cdFbD079c", // member 3
			"7CBD87ABC182216A7Aa0E8d19aA21abFA2511383", // member 4
			"FAc73b03884d94a08a5c6c7BB12Ac0b20571F162", // member 5
			"705C76445651530fe0D25eeE287b6164cE2c7216", // member 6
			"7CBD87ABC182216A7Aa0E8d19aA21abFA2511383", // member 7  (same operator as member 4)
			"405ad1f632b49A0617fbdc1fD427aF54BA9Bb3dd", // member 8
			"7CBD87ABC182216A7Aa0E8d19aA21abFA2511383", // member 9  (same operator as member 4)
			"5E14c0f27612fbfB7A6FE40b5A6Ec997fA62fc04", // member 10 (same operator as member 2)
		},
	}

	executor := &coordinationExecutor{
		// Set only relevant fields.
		coordinatedWallet: coordinatedWallet,
	}

	leader := executor.coordinationLeader(seed)

	testutils.AssertStringsEqual(
		t,
		"coordination leader",
		"D2662604f8b4540336fBd3c1F48d7e9cdFbD079c",
		leader.String(),
	)
}
