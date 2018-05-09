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
			expectation:  time.Duration(100000000),
			errorMessage: "Failed to wait for a single block",
		},
		"waited for a longer time": {
			blockWait:    2,
			expectation:  time.Duration(200000000),
			errorMessage: "Failed to wait for 2 blocks",
		},
		"doesn't wait if 0 blocks": {
			blockWait:    0,
			expectation:  time.Duration(1000),
			errorMessage: "Failed for a 0 block wait",
		},
		"invalid value": {
			blockWait:    -1,
			expectation:  time.Duration(0),
			errorMessage: "Waiting for a time when it should have errored",
		},
	}

	countWait := BlockCounter()

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			start := time.Now().UTC()
			countWait.WaitForBlocks(test.blockWait)
			end := time.Now().UTC()

			elapsed := end.Sub(start)
			if elapsed < test.expectation {
				t.Error(test.errorMessage)
			}
		})
	}
}
