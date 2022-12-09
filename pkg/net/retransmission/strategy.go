package retransmission

import "github.com/keep-network/keep-core/pkg/net"

// Strategy represents a specific retransmission strategy.
type Strategy interface {
	// Tick asks the strategy to run the provided retransmission routine.
	// The strategy uses their internal state and logic to decide whether to
	// call the retransmission function or not.
	Tick(retransmitFn RetransmitFn) error
}

// WithStrategy is a strategy factory function that returns the requested
// strategy instance.
func WithStrategy(strategy net.RetransmissionStrategy) Strategy {
	switch strategy {
	case net.StandardRetransmissionStrategy:
		return WithStandardStrategy()
	case net.BackoffRetransmissionStrategy:
		return WithBackoffStrategy()
	default:
		panic("retransmission strategy not implemented")
	}
}

// StandardStrategy is the basic retransmission strategy that triggers the
// retransmission routine on every tick.
type StandardStrategy struct{}

// WithStandardStrategy uses the StandardStrategy as the retransmission
// strategy.
func WithStandardStrategy() *StandardStrategy {
	return &StandardStrategy{}
}

// Tick implements the Strategy.Tick function.
func (ss *StandardStrategy) Tick(retransmitFn RetransmitFn) error {
	return retransmitFn()
}

// BackoffStrategy is a retransmission strategy that triggers the retransmission
// routine with an exponentially increasing delay. That is, the delay between
// first and second retransmission is 1 tick, between second and third is 2
// ticks, between third and fourth is 4 ticks and so on. Graphically, the
// schedule looks as follows: R _ R _ _ R _ _ _ _  R _ _ _ _ _ _ _ _ R
type BackoffStrategy struct {
	tickCounter    uint64
	delay          uint64
	retransmitTick uint64
}

// WithBackoffStrategy uses the BackoffStrategy as the retransmission
// strategy.
func WithBackoffStrategy() *BackoffStrategy {
	return &BackoffStrategy{
		tickCounter:    0,
		delay:          1,
		retransmitTick: 1,
	}
}

// Tick implements the Strategy.Tick function.
func (bos *BackoffStrategy) Tick(retransmitFn RetransmitFn) error {
	bos.tickCounter++

	if bos.tickCounter == bos.retransmitTick {
		bos.retransmitTick += bos.delay + 1
		bos.delay *= 2

		return retransmitFn()
	}

	return nil
}
