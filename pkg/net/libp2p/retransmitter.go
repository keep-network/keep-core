package libp2p

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
)

// retransmitter is a message retransmission strategy for libp2p broadcast
// channel retransmitting message for the lifetime of the context with the
// given interval between retransmissions.
//
// libp2p pubsub used internally by the broadcast channel does not guarantee
// message delivery. To improve the delivery rate, each message can be
// retransmitted several times.
type retransmitter struct {
	interval time.Duration
	cache    *timeCache
}

func newRetransmitter(intervalMilliseconds int) *retransmitter {
	interval := time.Duration(intervalMilliseconds) * time.Millisecond
	cacheLifetime := 2 * time.Minute // TODO: can be invalidated with context

	cache := newTimeCache(cacheLifetime)

	return &retransmitter{
		interval: interval,
		cache:    cache,
	}
}

// scheduleRetransmission takes the provided message and retransmits it
// with the given interval for the entire lifetime of the provided context.
// For each retransmission, send function is called with a copy of the original
// message and message retransmission counter set to the appropriate value.
func (r *retransmitter) scheduleRetransmissions(
	ctx context.Context,
	message *pb.NetworkMessage,
	send func(*pb.NetworkMessage) error,
) {
	go func() {
		retransmission := uint32(0)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(r.interval)
				retransmission++
				messageCopy := *message
				messageCopy.Retransmission = retransmission

				if err := send(&messageCopy); err != nil {
					logger.Errorf(
						"could not retransmit message: [%v]",
						err,
					)
				}

			}
		}
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
