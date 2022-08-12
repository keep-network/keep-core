package miner

import (
	"math/big"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

var one = big.NewInt(1)

// TestMineStop tests the situation when two new worker functions are added to
// a miner in a working state. The test ensures the worker functions starts
// doing their work. Then, the miner is stopped and the test ensures the worker
// functions are stopped.
func TestMineStop(t *testing.T) {
	miner := new(Miner)

	number1 := big.NewInt(0)
	number2 := big.NewInt(0)

	miner.Mine(func() {
		number1.Add(number1, one)
	})
	miner.Mine(func() {
		number2.Add(number2, one)
	})

	// give some time for the miner to perform computations
	time.Sleep(10 * time.Millisecond)

	// ensure computations started
	testutils.AssertBigIntNonZero(t, "computation result", number1)
	testutils.AssertBigIntNonZero(t, "computation result", number2)

	// send the stop signal and give the miner some time to stop computations
	miner.Stop()
	time.Sleep(100 * time.Millisecond)

	// at this point, all computations should be stopped, capture the current
	// result
	result1 := new(big.Int).Set(number1)
	result2 := new(big.Int).Set(number2)

	// wait some time and ensure computations stopped
	time.Sleep(20 * time.Millisecond)
	testutils.AssertBigIntsEqual(
		t,
		"computation result after stop signal",
		result1,
		number1,
	)
	testutils.AssertBigIntsEqual(
		t,
		"computation result after stop signal",
		result2,
		number2,
	)
}

// TestMineStopResume tests the situation when two new worker functions are
// added to a miner in a working state. Then, the miner is stopped and after
// some time its work is resumed. The test ensures the worker functions resume
// their work.
func TestMineStopResume(t *testing.T) {
	miner := new(Miner)
	defer miner.Stop()

	number1 := big.NewInt(0)
	number2 := big.NewInt(0)

	miner.Mine(func() {
		number1.Add(number1, one)
	})
	miner.Mine(func() {
		number2.Add(number2, one)
	})

	// send the stop signal and give the miner some time to stop computations
	miner.Stop()
	time.Sleep(100 * time.Millisecond)

	// at this point, all computations should be stopped, capture the current
	// result
	intermediateResult1 := new(big.Int).Set(number1)
	intermediateResult2 := new(big.Int).Set(number2)

	// send the resume signal and give the miner some time to resume computations
	miner.Resume()
	time.Sleep(100 * time.Millisecond)

	// ensure computations have been resumed
	testutils.AssertBigIntsNotEqual(
		t,
		"computation results after resume signal",
		intermediateResult1,
		number1,
	)
	testutils.AssertBigIntsNotEqual(
		t,
		"computation results after resume signal",
		intermediateResult2,
		number2,
	)
}

// TestMineStopResumeStop tests the situation when two new worker functions are
// added to a miner in a working state. Then, the miner is stopped, resumed,
// and stopped again. The test ensures the worker functions are stopped at the
// end of the cycle.
func TestMineStopResumeStop(t *testing.T) {
	miner := new(Miner)

	number1 := big.NewInt(0)
	number2 := big.NewInt(0)

	miner.Mine(func() {
		number1.Add(number1, one)
	})
	miner.Mine(func() {
		number2.Add(number2, one)
	})

	miner.Stop()
	miner.Resume()
	miner.Stop()

	// at this point, all computations should be stopped, capture the current
	// result
	result1 := new(big.Int).Set(number1)
	result2 := new(big.Int).Set(number2)

	// wait some time and ensure computations stopped
	time.Sleep(20 * time.Millisecond)
	testutils.AssertBigIntsEqual(
		t,
		"computation result after stop signal",
		result1,
		number1,
	)
	testutils.AssertBigIntsEqual(
		t,
		"computation result after stop signal",
		result2,
		number2,
	)
}

// TestStopMineResume tests the situation when two new worker functions are
// added to a stopped miner. Then, the miner work is resumed. The test ensures
// the worker functions are not working before the resume and that they are
// working after the resume.
func TestStopMineResume(t *testing.T) {
	miner := new(Miner)
	defer miner.Stop()

	miner.Stop()

	number1 := big.NewInt(0)
	number2 := big.NewInt(0)

	miner.Mine(func() {
		number1.Add(number1, one)
	})
	miner.Mine(func() {
		number2.Add(number2, one)
	})

	// assert computations have not started - the miner is stopped
	testutils.AssertBigIntsEqual(t, "computation result", big.NewInt(0), number1)
	testutils.AssertBigIntsEqual(t, "computation result", big.NewInt(0), number2)

	miner.Resume()
	// give some time for the miner to perform computations
	time.Sleep(10 * time.Millisecond)

	// ensure computations started
	testutils.AssertBigIntNonZero(t, "computation result", number1)
	testutils.AssertBigIntNonZero(t, "computation result", number2)
}
