package local_test

import (
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/chain/local"
)

func Test_Chain01(t *testing.T) {
	countWait := local.BlockCounter()
	if testing.Verbose() {
		t.Log("Before Wait")
	}
	start := time.Now()
	countWait.WaitForBlocks(3)
	tm := time.Now()
	elapsed := tm.Sub(start)
	if elapsed < 1400000000 {
		t.Fatalf("Did not wait\n")
	}
	if testing.Verbose() {
		t.Logf("After Wait, %d\n", elapsed)
	}

	start = time.Now()
	countWait.WaitForBlocks(5)
	tm = time.Now()
	elapsed = tm.Sub(start)
	if elapsed < 2400000000 {
		t.Fatalf("Did not wait\n")
	}

	if testing.Verbose() {
		t.Logf("Before test #3 , %d\n", elapsed)
	}
	start = time.Now()
	countWait.WaitForBlocks(0)
	tm = time.Now()
	elapsed = tm.Sub(start)
	if elapsed < 10 {
		t.Fatalf("Did not wait\n")
	}
}
