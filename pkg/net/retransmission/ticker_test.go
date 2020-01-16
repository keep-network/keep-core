package retransmission

import (
	"context"
	"testing"
	"time"
)

func TestOnTick(t *testing.T) {
	ticks := make(chan uint64)
	ticker := NewTicker(ticks)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tickCount := 0
	ticker.onTick(ctx, func() { tickCount++ })

	ticks <- 1
	ticks <- 2
	time.Sleep(10 * time.Millisecond)

	if tickCount != 2 {
		t.Errorf("expected [2] executions of handler, had [%v]", tickCount)
	}
}

func TestOnTickTimeTicker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 105*time.Millisecond)
	defer cancel()

	ticker := NewTimeTicker(ctx, 10*time.Millisecond)

	tickCount := 0
	ticker.onTick(ctx, func() { tickCount++ })

	<-ctx.Done()

	if tickCount != 10 {
		t.Errorf("expected [10] executions of handler, had [%v]", tickCount)
	}
}

func TestUnregisterHandler(t *testing.T) {
	ticks := make(chan uint64)
	ticker := NewTicker(ticks)

	ctx1, cancel1 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel1()
	ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel2()

	tickCount1 := 0
	ticker.onTick(ctx1, func() { tickCount1++ })

	tickCount2 := 0
	ticker.onTick(ctx2, func() { tickCount2++ })

	ticks <- 1
	ticks <- 2
	<-ctx1.Done()
	ticks <- 3
	<-ctx2.Done()
	ticks <- 4
	time.Sleep(10 * time.Millisecond)

	if tickCount1 != 2 {
		t.Errorf("expected [2] executions of the first handler, had [%v]", tickCount1)
	}
	if tickCount2 != 3 {
		t.Errorf("expected [3] executions of the second handler, had [%v]", tickCount2)
	}
}

func TestCloseTicker(t *testing.T) {
	ticks := make(chan uint64)
	ticker := NewTicker(ticks)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker.onTick(ctx, func() {})

	close(ticks)
	time.Sleep(10 * time.Millisecond)

	if len(ticker.handlers) != 0 {
		t.Errorf(
			"all handlers should be unregistered, still has [%v]",
			len(ticker.handlers),
		)
	}
}

func TestCloseTimeTicker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 105*time.Millisecond)
	defer cancel()

	ticker := NewTimeTicker(ctx, 10*time.Millisecond)

	ticker.onTick(ctx, func() {})

	<-ctx.Done()

	time.Sleep(10 * time.Millisecond)

	if len(ticker.handlers) != 0 {
		t.Errorf(
			"all handlers should be unregistered, still has [%v]",
			len(ticker.handlers),
		)
	}
}
