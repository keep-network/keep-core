package local

import (
	"context"
	"fmt"
	"github.com/keep-network/keep-core/pkg/operator"
	"sync"
	"sync/atomic"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/internal"
)

type unicastChannel struct {
	counter uint64

	structMutex *sync.RWMutex

	senderTransportID       net.TransportIdentifier
	senderOperatorPublicKey *operator.PublicKey

	receiverTransportID net.TransportIdentifier

	messageReceivers   []*unicastChannelRecv
	unmarshalersByType map[string]func() net.TaggedUnmarshaler
}

type unicastChannelRecv struct {
	ctx      context.Context
	handleFn func(message net.Message)
}

func newUnicastChannel(
	senderTransportID net.TransportIdentifier,
	senderOperatorPublicKey *operator.PublicKey,
	receiverTransportID net.TransportIdentifier,
) *unicastChannel {
	return &unicastChannel{
		structMutex:             &sync.RWMutex{},
		senderTransportID:       senderTransportID,
		senderOperatorPublicKey: senderOperatorPublicKey,
		receiverTransportID:     receiverTransportID,
		messageReceivers:        make([]*unicastChannelRecv, 0),
		unmarshalersByType:      make(map[string]func() net.TaggedUnmarshaler),
	}
}

func (uc *unicastChannel) nextSeqno() uint64 {
	return atomic.AddUint64(&uc.counter, 1)
}

func (uc *unicastChannel) Send(message net.TaggedMarshaler) error {
	marshalled, err := message.Marshal()
	if err != nil {
		return fmt.Errorf("could not marshal message [%v]", err)
	}

	return deliverMessage(
		uc.senderTransportID,
		uc.receiverTransportID,
		marshalled,
		message.Type(),
	)
}

func (uc *unicastChannel) Recv(
	ctx context.Context,
	handler func(message net.Message),
) {
	uc.structMutex.Lock()
	defer uc.structMutex.Unlock()

	uc.messageReceivers = append(
		uc.messageReceivers,
		&unicastChannelRecv{ctx: ctx, handleFn: handler},
	)
}

func (uc *unicastChannel) SetUnmarshaler(
	unmarshaler func() net.TaggedUnmarshaler,
) {
	uc.structMutex.Lock()
	defer uc.structMutex.Unlock()

	uc.unmarshalersByType[unmarshaler().Type()] = unmarshaler
}

func (uc *unicastChannel) receiveMessage(
	messagePayload []byte,
	messageType string,
) error {
	uc.structMutex.Lock()
	defer uc.structMutex.Unlock()

	unmarshaler, found := uc.unmarshalersByType[messageType]
	if !found {
		return fmt.Errorf(
			"couldn't find unmarshaler for type [%s]",
			messageType,
		)
	}

	unmarshaled := unmarshaler()
	err := unmarshaled.Unmarshal(messagePayload)
	if err != nil {
		return err
	}

	senderOperatorPublicKeyBytes := operator.MarshalUncompressed(uc.senderOperatorPublicKey)

	message := internal.BasicMessage(
		uc.senderTransportID,
		unmarshaled,
		messageType,
		senderOperatorPublicKeyBytes,
		uc.nextSeqno(),
	)

	i := 0
	for _, receiver := range uc.messageReceivers {
		// check if still active...
		if receiver.ctx.Err() == nil {
			// still active, should remain in the slice
			uc.messageReceivers[i] = receiver
			i++

			// firing handler asynchronously to
			// do not block the loop
			go receiver.handleFn(message)
		}
	}

	// cleaning up those no longer active
	uc.messageReceivers = uc.messageReceivers[:i]

	return nil
}
