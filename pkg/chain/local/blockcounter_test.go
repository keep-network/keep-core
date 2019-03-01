package local

import (
	"math/big"
	"testing"
	"time"
)

func TestBlockWaiter(t *testing.T) {
	var tests = map[string]struct {
		blockWait        int
		expectedWaitTime time.Duration
	}{
		"does wait for a block": {
			blockWait:        1,
			expectedWaitTime: blockTime,
		},
		"does wait for two blocks": {
			blockWait:        2,
			expectedWaitTime: 2 * blockTime,
		},
		"does wait for three blocks": {
			blockWait:        3,
			expectedWaitTime: 3 * blockTime,
		},
		"does not wait for 0 blocks": {
			blockWait:        0,
			expectedWaitTime: 0,
		},
		"does not wait for negative number of blocks": {
			blockWait:        -1,
			expectedWaitTime: 0,
		},
	}

	for testName, test := range tests {
		test := test
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			c := Connect(10, 4, big.NewInt(200))
			countWait, err := c.BlockCounter()
			if err != nil {
				t.Fatalf("failed to set up block counter: [%v]", err)
			}

			start := time.Now().UTC()
			countWait.WaitForBlocks(test.blockWait)
			end := time.Now().UTC()

			elapsed := end.Sub(start)

			// Block waiter should wait for test.expectedWaitTime minus some
			// margin at minimum; the margin is needed because clock is not
			// always that precise. Setting it to 5ms for this test.
			minMargin := time.Duration(5) * time.Millisecond
			if elapsed < (test.expectedWaitTime - minMargin) {
				t.Errorf(
					"waited less than expected; expected [%v] at min, waited [%v]",
					test.expectedWaitTime+minMargin,
					elapsed,
				)
			}

			// Block waiter should wait for test.expectedWaitTime plus some
			// margin at maximum; the margin is the time needed for the return
			// instructions to execute, setting it to 25ms for this test.
			maxMargin := time.Duration(25) * time.Millisecond
			if elapsed > (test.expectedWaitTime + maxMargin) {
				t.Errorf(
					"waited longer than expected; expected %v at max, waited %v",
					test.expectedWaitTime,
					elapsed,
				)
			}
		})
	}
}

func TestBlockHeightWaiter(t *testing.T) {
	var tests = map[string]struct {
		blockHeight      int
		initialDelay     time.Duration
		expectedWaitTime time.Duration
	}{
		"does not wait for negative block height": {
			blockHeight:      -1,
			expectedWaitTime: 0,
		},
		"returns immediately for genesis block": {
			blockHeight:      0,
			expectedWaitTime: 0,
		},
		"returns immediately for block height already reached": {
			blockHeight:      2,
			initialDelay:     3 * blockTime,
			expectedWaitTime: 0,
		},
		"waits for block height not yet reached": {
			blockHeight:      5,
			initialDelay:     2 * blockTime,
			expectedWaitTime: 3 * blockTime,
		},
	}

	for testName, test := range tests {
		test := test
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			c := Connect(10, 4, big.NewInt(100))

			blockCounter, err := c.BlockCounter()
			if err != nil {
				t.Fatalf("failed to set up block counter: [%v]", err)
			}

			time.Sleep(test.initialDelay)

			start := time.Now().UTC()
			blockCounter.WaitForBlockHeight(test.blockHeight)
			end := time.Now().UTC()

			elapsed := end.Sub(start)

			// Block waiter should wait for test.expectedWaitTime minus some
			// margin at minimum; the margin is needed because clock is not
			// always that precise. Setting it to 5ms for this test.
			minMargin := time.Duration(5) * time.Millisecond
			if elapsed < (test.expectedWaitTime - minMargin) {
				t.Errorf(
					"waited less than expected; expected [%v] at min, waited [%v]",
					test.expectedWaitTime,
					elapsed,
				)
			}

			// Block waiter should wait for test.expectedWaitTime plus some
			// margin at maximum; the margin is the time needed for the return
			// instructions to execute, setting it to 25ms for this test.
			maxMargin := time.Duration(25) * time.Millisecond
			if elapsed > (test.expectedWaitTime + maxMargin) {
				t.Errorf(
					"waited longer than expected; expected %v at max, waited %v",
					test.expectedWaitTime,
					elapsed,
				)
			}
		})
	}
}
