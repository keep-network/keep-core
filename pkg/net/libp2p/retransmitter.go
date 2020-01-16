package libp2p

import (
	"time"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
)

// retransmitter is a message retransmission strategy for libp2p broadcast
// channel retransmitting message for the certain number of cycles and with the
// given interval.
//
// libp2p pubsub used internally by the broadcast channel does not guarantee
// message delivery. To improve the delivery rate, each message can be
// retransmitted a certain number of times.
type retransmitter struct {
	cycles   uint32
	interval time.Duration
}

func newRetransmitter(cycles int, intervalMilliseconds int) *retransmitter {
	interval := time.Duration(intervalMilliseconds) * time.Millisecond

	return &retransmitter{
		cycles:   uint32(cycles),
		interval: interval,
	}
}

// scheduleRetransmission takes the provided message and retransmits it
// according to the configured number of cycles and interval using the
// provided send function. For each retransmission, send function is
// called with a copy of the original message and message retransmission
// counter set to the appropriate value.
func (r *retransmitter) scheduleRetransmission(
	message *pb.NetworkMessage,
	send func(*pb.NetworkMessage) error,
) {
	go func() {
		for i := uint32(1); i <= r.cycles; i++ {
			time.Sleep(r.interval)

			messageCopy := *message
			messageCopy.Retransmission = i

			go func() {
				if err := send(&messageCopy); err != nil {
					logger.Errorf(
						"could not retransmit message: [%v]",
						err,
					)
				}
			}()
		}
		return
	}()
}
