package retransmission

import (
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-net-retransmission")

type retransmittingBroadcastChannel struct {
	delegate net.BroadcastChannel

	retransmissionInterval time.Duration
	retransmissionCycles   int
}

// WithRetransmission decorates the given broadcast channel with the ability
// of automatic message retransmission with desired retransmission interval
// and number of retransmission cycles.
func WithRetransmission(
	delegate net.BroadcastChannel,
	retransmissionInterval time.Duration,
	retransmissionCycles int,
) net.BroadcastChannel {
	return &retransmittingBroadcastChannel{
		delegate:               delegate,
		retransmissionInterval: retransmissionInterval,
		retransmissionCycles:   retransmissionCycles,
	}
}

func (rbc *retransmittingBroadcastChannel) Name() string {
	return rbc.delegate.Name()
}

func (rbc *retransmittingBroadcastChannel) Send(m net.TaggedMarshaler) error {
	go func() {
		for i := 0; i <= rbc.retransmissionCycles; i++ {
			if i != 0 {
				time.Sleep(rbc.retransmissionInterval)
			}
			if err := rbc.delegate.Send(m); err != nil {
				logger.Errorf("Could not send message of type %v: [%v]",
					m.Type(),
					err,
				)
			}
		}
	}()

	return nil
}

func (rbc *retransmittingBroadcastChannel) Recv(
	h net.HandleMessageFunc,
) error {
	return rbc.delegate.Recv(h)
}

func (rbc *retransmittingBroadcastChannel) UnregisterRecv(
	handlerType string,
) error {
	return rbc.delegate.UnregisterRecv(handlerType)
}

func (rbc *retransmittingBroadcastChannel) RegisterUnmarshaler(
	unmarshaler func() net.TaggedUnmarshaler,
) error {
	return rbc.delegate.RegisterUnmarshaler(unmarshaler)
}

func (rbc *retransmittingBroadcastChannel) AddFilter(
	filter net.BroadcastChannelFilter,
) error {
	return rbc.delegate.AddFilter(filter)
}
