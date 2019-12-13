package libp2p

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	cache    *timeCache
}

func newRetransmitter(cycles int, intervalMilliseconds int) *retransmitter {
	interval := time.Duration(intervalMilliseconds) * time.Millisecond
	retransmissionDuration := time.Duration(cycles) * interval
	cacheLifetime := 2 * time.Minute

	cache := newTimeCache(retransmissionDuration + cacheLifetime)

	return &retransmitter{
		cycles:   uint32(cycles),
		interval: interval,
		cache:    cache,
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

// receive takes the received message and calls the provided onFirstTimeReceived
// function only if the message was not received before. The message can be
// the original one or a retransmission. To decide whether the given message
// was received before, retransmitter evaluates a fingerprint of the message
// which includes all the fields but the retransmission counter.
func (r *retransmitter) receive(
	message *pb.NetworkMessage,
	onFirstTimeReceived func() error,
) error {
	fingerprint, err := calculateFingerprint(message)
	if err != nil {
		return fmt.Errorf("could not calculate message fingerprint: [%v]", err)
	}

	if r.cache.has(fingerprint) && message.Retransmission != 0 {
		return nil
	}

	r.cache.add(fingerprint)
	return onFirstTimeReceived()
}

func calculateFingerprint(message *pb.NetworkMessage) (string, error) {
	// Reset retransmission counter to 0. We do not want the retransmission
	// counter value to change the message fingerprint.
	messageCopy := *message
	messageCopy.Retransmission = 0

	bytes, err := messageCopy.Marshal()
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(bytes)
	return hex.EncodeToString(sum[:]), nil
}
