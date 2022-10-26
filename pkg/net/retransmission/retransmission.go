// Package retransmission implements a simple retransmission mechanism for
// network messages based on their sequence number. Retransmitting message
// several times for the lifetime of the given phase helps to improve message
// delivery rate for senders and receivers who are not perfectly synced on time.
package retransmission

import (
	"context"
	"fmt"
	"sync"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/net"
)

// ScheduleRetransmissions takes the provided message and retransmits it
// for every new tick received from the provided Ticker for the entire lifetime
// of the Context calling the provided retransmit function. The retransmit
// function has to guarantee that every call from this function sends a message
// with the same sequence number.
//
// The provided context is used to identify the retransmission routine of
// the given message. If there is a need to use a single context to send
// multiple messages, it is recommended to use a separate child context
// for each message. Otherwise, only the last message will be retransmitted.
func ScheduleRetransmissions(
	ctx context.Context,
	logger log.StandardLogger,
	ticker *Ticker,
	retransmit func() error,
) {
	go func() {
		ticker.onTick(ctx, func() {
			go func() {
				if err := retransmit(); err != nil {
					logger.Errorf("could not retransmit message: [%v]", err)
				}
			}()
		})
	}()
}

// WithRetransmissionSupport takes the standard network message handler and
// enhances it with functionality allowing to handle retransmissions.
// The returned handler filters out retransmissions and calls the delegate
// handler only if the received message is not a retransmission or if it is
// a retransmission but it has not been seen by the original handler yet.
// The returned handler is thread-safe.
//
// Retransmissions are identified by sender transport ID and message sequence
// number. Two messages with the same sender ID and sequence number are
// considered the same. Handler can not be reused between channels if sequence
// number of message is local for channel.
func WithRetransmissionSupport(delegate func(m net.Message)) func(m net.Message) {
	mutex := &sync.Mutex{}
	cache := make(map[string]bool)

	return func(message net.Message) {
		messageID := fmt.Sprintf(
			"%v-%v",
			message.TransportSenderID().String(),
			message.Seqno(),
		)

		mutex.Lock()
		_, seen := cache[messageID]
		if !seen {
			cache[messageID] = true
		}
		mutex.Unlock()

		if !seen {
			delegate(message)
		}
	}
}
