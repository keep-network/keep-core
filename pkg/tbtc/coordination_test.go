package tbtc

import (
	"context"
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

func TestCoordinationWindow_IsAfter(t *testing.T) {
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
