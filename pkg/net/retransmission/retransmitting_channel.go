package retransmission

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
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
				logger.Errorf(
					"could not send message of type %v: [%v]",
					m.Type(),
					err,
				)
			}
		}
	}()

	return nil
}

func (rbc *retransmittingBroadcastChannel) Recv(
	handleMessageFunc net.HandleMessageFunc,
) error {
	messageCache := make(map[string]bool)
	messageCacheMutex := &sync.Mutex{}

	cachingHandler := func(msg net.Message) error {
		payloadChecksum, err := calculatePayloadChecksum(msg)
		if err != nil {
			return err
		}

		messageCacheMutex.Lock()
		defer messageCacheMutex.Unlock()

		if _, ok := messageCache[payloadChecksum]; !ok {
			messageCache[payloadChecksum] = true
			return handleMessageFunc.Handler(msg)
		}

		return nil
	}

	cachingHandleMessageFunc := net.HandleMessageFunc{
		Type:    handleMessageFunc.Type,
		Handler: cachingHandler,
	}

	return rbc.delegate.Recv(cachingHandleMessageFunc)
}

func calculatePayloadChecksum(message net.Message) (string, error) {
	payload, ok := message.Payload().(net.TaggedMarshaler)
	if !ok {
		return "", fmt.Errorf("could not cast message payload")
	}

	payloadBytes, err := payload.Marshal()
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(payloadBytes)
	return hex.EncodeToString(sum[:]), nil
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
