package tbtc

import (
	"context"
	"encoding/hex"
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
