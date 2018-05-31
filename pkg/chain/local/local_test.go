package local

import (
	"testing"
	"time"
)

func TestLocalBlockWaiter(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		blockWait    int
		expectation  time.Duration
		errorMessage string
	}{
		"does wait for a block": {
			blockWait:    1,
			expectation:  time.Duration(600) * time.Millisecond,
			errorMessage: "Failed to wait for a single block; expected %s but took %s.",
		},
		"waited for a longer time": {
			blockWait:    2,
			expectation:  time.Duration(600*2) * time.Millisecond,
			errorMessage: "Failed to wait for 2 blocks; expected %s but took %s.",
		},
		"doesn't wait if 0 blocks": {
			blockWait:    0,
			expectation:  time.Duration(20) * time.Microsecond,
			errorMessage: "Failed for a 0 block wait; expected %s but took %s.",
		},
		"invalid value": {
			blockWait:    -1,
			expectation:  time.Duration(20) * time.Microsecond,
			errorMessage: "Waiting for a time when it should have errored; expected %s but took %s.",
		},
	}

	c := Connect()
	countWait := c.BlockCounter()

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			start := time.Now().UTC()
			countWait.WaitForBlocks(test.blockWait)
			end := time.Now().UTC()

			elapsed := end.Sub(start)
			if test.expectation < elapsed {
				t.Errorf(test.errorMessage, test.expectation, elapsed)
			}
		})
	}
}
