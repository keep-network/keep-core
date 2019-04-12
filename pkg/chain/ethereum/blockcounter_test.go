package ethereum

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestWaitForNewMinedBlock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	blockCounter := &ethereumBlockCounter{
		latestBlockHeightSeen: uint64(1),
		waiters:               make(map[uint64][]chan uint64),
		subscriptionChannel:   make(chan block),
	}

	block2Waiter, err := blockCounter.BlockHeightWaiter(2)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		blockCounter.subscriptionChannel <- block{Number: "1"}
		blockCounter.subscriptionChannel <- block{Number: "2"}
	}()

	go blockCounter.receiveBlocks()

	select {
	case block := <-block2Waiter:
		if block != 2 {
			t.Fatalf(
				"unexpected block number\nexpected: [2]\nactual:   [%v]\n",
				block,
			)
		}

	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestExecuteBlockHandlerOnlyOnce(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	blockCounter := &ethereumBlockCounter{
		latestBlockHeightSeen: uint64(1),
		waiters:               make(map[uint64][]chan uint64),
		subscriptionChannel:   make(chan block),
	}

	block2Waiter, err := blockCounter.BlockHeightWaiter(2)
	if err != nil {
		t.Fatal(err)
	}

	block3Waiter, err := blockCounter.BlockHeightWaiter(3)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		blockCounter.subscriptionChannel <- block{Number: "1"}
		blockCounter.subscriptionChannel <- block{Number: "2"}
		blockCounter.subscriptionChannel <- block{Number: "2"}
		blockCounter.subscriptionChannel <- block{Number: "3"}
		blockCounter.subscriptionChannel <- block{Number: "4"}
	}()

	go blockCounter.receiveBlocks()

	var waiter2CallCounter uint64
	var waiter3CallCounter uint64

	for {
		select {
		case <-block2Waiter:
			atomic.AddUint64(&waiter2CallCounter, 1)

		case <-block3Waiter:
			atomic.AddUint64(&waiter3CallCounter, 1)

		case <-ctx.Done():
			if waiter2CallCounter != 1 {
				t.Errorf(
					"handler for block 2 should be called only once; was [%v]",
					waiter2CallCounter,
				)
			}
			if waiter3CallCounter != 1 {
				t.Errorf(
					"handler for block 3 should be called only once; was [%v]",
					waiter3CallCounter,
				)
			}

			return
		}
	}
}
