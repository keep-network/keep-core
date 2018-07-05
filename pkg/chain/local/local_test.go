package local

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
)

func TestLocalSubmitRelayEntry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chainHandle := Connect(10, 4).ThresholdRelay()
	relayEntryPromise := chainHandle.SubmitRelayEntry(
		&event.Entry{
			RequestID: big.NewInt(int64(19)),
			GroupID:   big.NewInt(int64(1)),
		},
	)

	done := make(chan *event.Entry)
	relayEntryPromise.OnSuccess(func(entry *event.Entry) {
		done <- entry
	}).OnFailure(func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	})

	select {
	case entry := <-done:
		expected := int64(19)
		if entry.RequestID.Int64() != expected {
			t.Fatalf(
				"expected [%v], got [%v]",
				expected,
				entry.RequestID.Int64(),
			)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}

}

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

			c := Connect(10, 4)
			countWait, err := c.BlockCounter()
			if err != nil {
				t.Fatalf("failed to set up block counter: [%v]", err)
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
