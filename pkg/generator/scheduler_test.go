package generator

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

var one = big.NewInt(1)

// TestComputeStop tests the situation when two new worker functions are added
// to a scheduler in a working state. The test ensures the worker functions
// starts doing their work. Then, the scheduler is stopped and the test ensures
// the worker functions are stopped.
func TestComputeStop(t *testing.T) {
	scheduler := new(Scheduler)

	number1 := big.NewInt(0)
	number2 := big.NewInt(0)

	scheduler.compute(func(context.Context) {
		number1.Add(number1, one)
	})
	scheduler.compute(func(context.Context) {
		number2.Add(number2, one)
	})

	// give some time to perform computations
	time.Sleep(10 * time.Millisecond)

	// ensure computations started
	testutils.AssertBigIntNonZero(t, "computation result", number1)
	testutils.AssertBigIntNonZero(t, "computation result", number2)

	// send the stop signal and give some time to stop computations
	scheduler.stop()
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

// TestComputeStopContext covers the same situation as TestComputeStop except
// that it ensures if the context passed to the function is cancelled when the
// work is stopped.
func TestComputeStopContext(t *testing.T) {
	scheduler := new(Scheduler)

	cancelled1 := false
	cancelled2 := false

	scheduler.compute(func(ctx context.Context) {
		// this simulates a long-running task
		<-ctx.Done()
		cancelled1 = true
	})
	scheduler.compute(func(ctx context.Context) {
		// this simulates a long-running task
		<-ctx.Done()
		cancelled2 = true
	})

	// give some time to perform computations
	time.Sleep(10 * time.Millisecond)

	// send the stop signal and give some time to stop computations
	scheduler.stop()
	time.Sleep(100 * time.Millisecond)

	// ensure context got cancelled
	if !cancelled1 {
		t.Errorf("expected context to be cancelled")
	}
	if !cancelled2 {
		t.Errorf("expected context to be cancelled")
	}
}

// TestComputeStopResume tests the situation when two new worker functions are
// added to a scheduler in a working state. Then, the scheduler is stopped and
// after some time its work is resumed. The test ensures the worker functions
// resume their work.
func TestComputeStopResume(t *testing.T) {
	scheduler := new(Scheduler)
	defer scheduler.stop()

	number1 := big.NewInt(0)
	number2 := big.NewInt(0)

	scheduler.compute(func(context.Context) {
		number1.Add(number1, one)
	})
	scheduler.compute(func(context.Context) {
		number2.Add(number2, one)
	})

	// send the stop signal and give some time to stop computations
	scheduler.stop()
	time.Sleep(100 * time.Millisecond)

	// at this point, all computations should be stopped, capture the current
	// result
	intermediateResult1 := new(big.Int).Set(number1)
	intermediateResult2 := new(big.Int).Set(number2)

	// send the resume signal and give some time to resume computations
	scheduler.resume()
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

// TestComputeStopResumeStop tests the situation when two new worker functions
// are added to a scheduler in a working state. Then, the scheduler is stopped,
// resumed, and stopped again. The test ensures the worker functions are stopped
// at the end of the cycle.
func TestComputeStopResumeStop(t *testing.T) {
	scheduler := new(Scheduler)

	number1 := big.NewInt(0)
	number2 := big.NewInt(0)

	scheduler.compute(func(context.Context) {
		number1.Add(number1, one)
	})
	scheduler.compute(func(context.Context) {
		number2.Add(number2, one)
	})

	scheduler.stop()
	scheduler.resume()
	scheduler.stop()

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

// TestStopComputeResume tests the situation when two new worker functions are
// added to a stopped scheduler. Then, the scheduler work is resumed. The test
// ensures the worker functions are not working before the resume and that they
// are working after the resume.
func TestStopComputeResume(t *testing.T) {
	scheduler := new(Scheduler)
	defer scheduler.stop()

	scheduler.stop()

	number1 := big.NewInt(0)
	number2 := big.NewInt(0)

	scheduler.compute(func(context.Context) {
		number1.Add(number1, one)
	})
	scheduler.compute(func(context.Context) {
		number2.Add(number2, one)
	})

	// assert computations have not started - the scheduler is stopped
	testutils.AssertBigIntsEqual(t, "computation result", big.NewInt(0), number1)
	testutils.AssertBigIntsEqual(t, "computation result", big.NewInt(0), number2)

	scheduler.resume()
	// give some time to perform computations;
	// given the goroutines are started after Resume call, we are giving this
	// test a bit more time than others
	time.Sleep(250 * time.Millisecond)

	// ensure computations started
	testutils.AssertBigIntNonZero(t, "computation result", number1)
	testutils.AssertBigIntNonZero(t, "computation result", number2)
}
