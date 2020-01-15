package retransmission

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
)

var logger = log.Logger("keep-net-retransmission")

// NetworkMessage enhances net.Message with functions needed to perform message
// retransmission. Specifically, we need to know if the given message is
// a retransmission and know a fingerprint of the message to filter out
// duplicates in the message handler.
type NetworkMessage interface {
	net.Message

	Fingerprint() string
	Retransmission() uint32
}

// ScheduleRetransmissions takes the provided message and retransmits it
// for every new tick received from the provided Ticker for the entire lifetime
// of the provided Context. For each retransmission, send function is
// called with a copy of the original message and message retransmission
// counter set to the appropriate value.
func ScheduleRetransmissions(
	ctx context.Context,
	ticker *Ticker,
	message *pb.NetworkMessage,
	send func(*pb.NetworkMessage) error,
) {
	retransmission := uint32(0)
	ticker.onTick(ctx, func() {
		retransmission++

		messageCopy := *message
		messageCopy.Retransmission = retransmission

		go func() {
			if err := send(&messageCopy); err != nil {
				logger.Errorf(
					"could not retransmit message: [%v]",
					err,
				)
			}
		}()
	})
}

// WithRetransmissionSupport takes the standard network message handler and
// enhances it with functionality allowing to handle retransmissions.
// The returned handler filters out retransmissions and calls the delegate
// handler only if the received message is not a retransmission or if it is
// a retransmission but it has not been seen by the original handler yet.
// The returned handler is thread-safe.
func WithRetransmissionSupport(delegate func(m net.Message)) func(m NetworkMessage) {
	mutex := &sync.Mutex{}
	cache := make(map[string]bool)

	return func(message NetworkMessage) {
		fingerprint := message.Fingerprint()

		mutex.Lock()
		_, seenFingerprint := cache[fingerprint]
		if !seenFingerprint {
			cache[fingerprint] = true
		}
		mutex.Unlock()

		isRetransmission := message.Retransmission() != 0

		if seenFingerprint && isRetransmission {
			return
		}

		delegate(message)
	}
}

// CalculateFingerprint takes the serialized network message and its sender
// and produces a unique fingerprint for that message. Fingerprint is a part
// of NetworkMessage and is used by message handler to filter out
// already known retransmissions.
func CalculateFingerprint(
	sender net.TransportIdentifier,
	payload []byte,
) string {
	sha := sha256.New()
	sha.Write([]byte(sender.String()))
	sha.Write(payload)
	return hex.EncodeToString(sha.Sum(nil)[:])
}

// NewNetworkMessage accepts an ordinary net.Message as well as retransmission
// information and produces NetworkMessage instance that can be accepted
// by retransmission handler.
func NewNetworkMessage(
	original net.Message,
	fingerprint string,
	retransmission uint32,
) NetworkMessage {
	return &networkMessage{
		original:       original,
		fingerprint:    fingerprint,
		retransmission: retransmission,
	}
}

type networkMessage struct {
	original       net.Message
	fingerprint    string
	retransmission uint32
}

func (nm *networkMessage) TransportSenderID() net.TransportIdentifier {
	return nm.original.TransportSenderID()
}

func (nm *networkMessage) Payload() interface{} {
	return nm.original.Payload()
}

func (nm *networkMessage) Type() string {
	return nm.original.Type()
}

func (nm *networkMessage) SenderPublicKey() []byte {
	return nm.original.SenderPublicKey()
}

func (nm *networkMessage) Fingerprint() string {
	return nm.fingerprint
}

func (nm *networkMessage) Retransmission() uint32 {
	return nm.retransmission
}
