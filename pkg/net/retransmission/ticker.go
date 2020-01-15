package retransmission

import (
	"context"
	"sync"
	"time"
)

// Ticker controls the frequency of retransmissions.
type Ticker struct {
	ticks         <-chan uint64
	handlersMutex sync.Mutex
	handlers      map[context.Context]func()
}

// NewTicker creates and starts a new Ticker for the provided channel.
// For each item read from the channel, new tick is triggered. All handlers
// are unregistered and ticker is stopped when the provided channel gets closed.
func NewTicker(ticks <-chan uint64) *Ticker {
	ticker := &Ticker{
		ticks:    ticks,
		handlers: make(map[context.Context]func()),
	}

	go ticker.start()

	return ticker
}

// NewTimeTicker is a convenience function allowing to convert time.Ticker to
// retransmission.Ticker. When the provided time.Ticker is stopped, all handlers
// are unregistered and retransmission.Ticker is stopped.
func NewTimeTicker(ticker *time.Ticker) *Ticker {
	ticks := make(chan uint64)
	go func() {
		for tick := range ticker.C {
			ticks <- uint64(tick.Unix())
		}
	}()

	return NewTicker(ticks)
}

func (t *Ticker) start() {
	for range t.ticks {
		t.handlersMutex.Lock()

		for ctx, handler := range t.handlers {
			if ctx.Err() != nil {
				delete(t.handlers, ctx)
				continue
			}

			handler()
		}

		t.handlersMutex.Unlock()
	}

	for ctx := range t.handlers {
		delete(t.handlers, ctx)
	}
}

func (t *Ticker) onTick(ctx context.Context, handler func()) {
	t.handlersMutex.Lock()
	t.handlers[ctx] = handler
	t.handlersMutex.Unlock()
}
