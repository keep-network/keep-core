package local

import (
	"testing"
	"time"
)

func TestLocalBlockWaiter(t *testing.T) {
	var tests = map[string]struct {
		blockWait    int
		expectation  time.Duration
		errorMessage string
	}{
		"does wait for a block": {
			blockWait:    1,
			expectation:  time.Duration(525) * time.Millisecond,
			errorMessage: "Failed to wait for a single block; expected %s but took %s.",
		},
		"waited for a longer time": {
			blockWait:    2,
			expectation:  time.Duration(525*2) * time.Millisecond,
			errorMessage: "Failed to wait for 2 blocks; expected %s but took %s.",
		},
		"doesn't wait if 0 blocks": {
			blockWait:    0,
			expectation:  time.Duration(20) * time.Millisecond,
			errorMessage: "Failed for a 0 block wait; expected %s but took %s.",
		},
		"invalid value": {
			blockWait:    -1,
			expectation:  time.Duration(20) * time.Millisecond,
			errorMessage: "Waiting for a time when it should have errored; expected %s but took %s.",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			c := Connect()
			countWait, err := c.BlockCounter()
			if err != nil {
				t.Errorf("failed to setup test")
			}

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
