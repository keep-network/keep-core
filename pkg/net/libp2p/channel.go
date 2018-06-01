package libp2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/internal"
	floodsub "github.com/libp2p/go-floodsub"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

type channel struct {
	name string

	identity *identity
	store    pstore.Peerstore

	pubsubLock   sync.Mutex
	pubsub       *floodsub.PubSub
	subscription *floodsub.Subscription

	messageLock sync.RWMutex
	messageBus  []net.Message

	unmarshalersMutex  sync.Mutex
	unmarshalersByType map[string]func() net.TaggedUnmarshaler

	identifiersMutex            sync.Mutex
	transportToProtoIdentifiers map[net.TransportIdentifier]net.ProtocolIdentifier
	protoToTransportIdentifiers map[net.ProtocolIdentifier]net.TransportIdentifier
}

func (c *channel) Name() string {
	return c.name
}

func (c *channel) Send(message net.TaggedMarshaler) error {
	return c.doSend(message, c.identity)
}

func envelopeProto(message net.TaggedMarshaler, sender *identity) ([]byte, error) {
	payloadBytes, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	identityBytes, err := sender.Marshal()
	if err != nil {
		return nil, err
	}

	env := &pb.Envelope{
		Payload: payloadBytes,
		Sender:  identityBytes,
		Type:    []byte(message.Type()),
	}

	return env.Marshal()
}

func (c *channel) doSend(message net.TaggedMarshaler, sender *identity) error {
	// 1. Change a net.TaggedMarshaler to a protobuf message
	envelopeBytes, err := envelopeProto(message, sender)
	if err != nil {
		return err
	}

	// 2. get that representation and put it onto the publish
	c.pubsubLock.Lock()
	if err := c.pubsub.Publish(c.name, envelopeBytes); err != nil {
		return err
	}
	c.pubsubLock.Unlock()

	return nil
}

func (c *channel) SendTo(
	recipientIdentifier interface{},
	message net.TaggedMarshaler,
) error {
	return nil
}

func (c *channel) Recv(h net.HandleMessageFunc) error {
	c.messageLock.RLock()
	snapshot := make([]net.Message, len(c.messageBus))
	copy(snapshot, c.messageBus)
	c.messageLock.RUnlock()

	for _, message := range snapshot {
		if err := h(message); err != nil {
			return err
		}
	}

	return nil
}

func (c *channel) RegisterIdentifier(
	transportIdentifier net.TransportIdentifier,
	protocolIdentifier net.ProtocolIdentifier,
) error {
	c.identifiersMutex.Lock()
	defer c.identifiersMutex.Unlock()

	if _, ok := transportIdentifier.(networkIdentity); !ok {
		return fmt.Errorf(
			"incorrect type for transportIdentifier: [%v] in channel [%s]",
			transportIdentifier, c.name,
		)
	}

	if _, exists := c.transportToProtoIdentifiers[transportIdentifier]; exists {
		return fmt.Errorf(
			"protocol identifier in channel [%s] already associated with [%v]",
			c.name, transportIdentifier,
		)
	}

	if _, exists := c.protoToTransportIdentifiers[protocolIdentifier]; exists {
		return fmt.Errorf(
			"transport identifier in channel [%s] already associated with [%v]",
			c.name, protocolIdentifier,
		)
	}

	c.transportToProtoIdentifiers[transportIdentifier] = protocolIdentifier
	c.protoToTransportIdentifiers[protocolIdentifier] = transportIdentifier

	return nil
}

func (c *channel) RegisterUnmarshaler(unmarshaler func() net.TaggedUnmarshaler) error {
	tpe := unmarshaler().Type()

	c.unmarshalersMutex.Lock()
	defer c.unmarshalersMutex.Unlock()

	if _, exists := c.unmarshalersByType[tpe]; exists {
		return fmt.Errorf("type %s already has an associated unmarshaler", tpe)
	}

	c.unmarshalersByType[tpe] = unmarshaler
	return nil
}

func (c *channel) handleMessages(ctx context.Context) {
	defer c.subscription.Cancel()

	for {
		// TODO: thread in a context with cancel
		msg, err := c.subscription.Next(ctx)
		if err != nil {
			// TODO: handle error - different error types
			// result in different outcomes
			fmt.Println(err)
			return
		}

		if err := c.processMessage(msg); err != nil {
			// TODO: handle error - different error types
			// result in different outcomes
			fmt.Println(err)
			return
		}

		select {
		case <-ctx.Done():
			return
		}
	}
}

func (c *channel) processMessage(message *floodsub.Message) error {
	// 1. Unmarshall the message in to the Envelope
	var envelope pb.Envelope
	if err := proto.Unmarshal(message.Data, &envelope); err != nil {
		return err
	}

	// TODO: 2. the whole is the receiver the senderIdentifier thing

	// 3. Since the protocol type is on the envelope, let's
	//    pull that type from our map of unmarshallers.
	c.unmarshalersMutex.Lock()
	unmarshaler, found := c.unmarshalersByType[string(envelope.Type)]
	if !found {
		return fmt.Errorf(
			"Couldn't find unmarshaler for type %s",
			string(envelope.Type),
		)
	}
	c.unmarshalersMutex.Unlock()

	unmarshaled := unmarshaler()
	if err := unmarshaled.Unmarshal(envelope.GetPayload()); err != nil {
		return err
	}

	// 4. Construct an identifier from the sender (on the message)
	senderIdentifier := &identity{}
	if err := senderIdentifier.Unmarshal(envelope.Sender); err != nil {
		return err
	}

	// 5. Get the associated protocol identifier from an association map
	c.identifiersMutex.Lock()
	protocolIdentifier, found := c.transportToProtoIdentifiers[senderIdentifier.id]
	if !found {
		return fmt.Errorf(
			"Couldn't find protocol identifier for sender identifier %+v",
			senderIdentifier,
		)
	}
	c.identifiersMutex.Unlock()

	// 6. Construct an internal.BasicMessage to fire back to the protocol
	protocolMessage := internal.BasicMessage(senderIdentifier.id,
		protocolIdentifier, unmarshaled,
	)

	// 7. Slap our internal message onto a channel from which we can
	//    pull off later.
	c.messageLock.Lock()
	c.messageBus = append(c.messageBus, protocolMessage)
	c.messageLock.Unlock()

	return nil
}
