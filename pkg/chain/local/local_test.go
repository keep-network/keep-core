package local

import (
	"testing"
	"time"
)

func TestLocalBlockWaiter(t *testing.T) {

	t.Parallel()

	var tests = map[string]struct {
		wait         int
		want         time.Duration
		errorMessage string
	}{
		"does wait for a block": {
			wait:         1,
			want:         time.Duration(100000000),
			errorMessage: "failed to wait for a single block",
		},
		"waited for a longer time": {
			wait:         2,
			want:         time.Duration(200000000),
			errorMessage: "failed to wait for 2 blocks",
		},
		"doesn't wait if 0 blocks": {
			wait:         0,
			want:         time.Duration(1000),
			errorMessage: "some errror message here...not understanding this case",
		},
		"invalid value": {
			wait:         -1,
			want:         time.Duration(0),
			errorMessage: "Waiting for a time when it should have errored",
		},
	}

	var e chan interface{}

	tim := 15 // Force test to fail if not completed in 15 seconds.
	tick := time.NewTimer(time.Duration(tim) * time.Second)

	go func() {
		select {
		case e <- tick:
			t.Fatal("Test ran too long - it failed")
		}
	}()

	countWait := BlockCounter()

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {

			start := time.Now().UTC()
			countWait.WaitForBlocks(test.wait)
			end := time.Now().UTC()

			elapsed := end.Sub(start)
			if elapsed < test.want {
				t.Error(test.errorMessage)
			}
		})
	}
	tick.Stop()

}
