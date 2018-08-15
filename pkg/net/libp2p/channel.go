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
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
)

type channel struct {
	name string

	clientIdentity *identity
	peerStore      peerstore.Peerstore

	pubsubMutex sync.Mutex
	pubsub      *floodsub.PubSub

	subscription *floodsub.Subscription

	messageHandlersMutex sync.Mutex
	messageHandlers      []net.HandleMessageFunc

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
	return c.doSend(nil, c.clientIdentity, message)
}

func (c *channel) SendTo(
	recipientIdentifier net.ProtocolIdentifier,
	message net.TaggedMarshaler,
) error {
	return c.doSend(recipientIdentifier, c.clientIdentity, message)
}

// doSend attempts to send a message, from a sender, to all members of a
// broadcastChannel, or optionally to a specific recipient. If recipient
// is nil (the typical case), then all messages of the broadcast channel
// should receive the message. Otherwise, given a valid recipient, we will
// address the message specifically to them.
func (c *channel) doSend(
	recipient net.ProtocolIdentifier,
	sender *identity,
	message net.TaggedMarshaler,
) error {
	var transportRecipient net.TransportIdentifier
	if recipient != nil {
		c.identifiersMutex.Lock()
		if transportID, ok := c.protoToTransportIdentifiers[recipient]; ok {
			transportRecipient = transportID
		}
		c.identifiersMutex.Unlock()
	}
	// Transform net.TaggedMarshaler to a protobuf message
	envelopeBytes, err := c.envelopeProto(transportRecipient, sender, message)
	if err != nil {
		return err
	}

	c.pubsubMutex.Lock()
	defer c.pubsubMutex.Unlock()

	// Publish the proto to the network
	return c.pubsub.Publish(c.name, envelopeBytes)
}

func (c *channel) Recv(handler net.HandleMessageFunc) error {
	c.messageHandlersMutex.Lock()
	c.messageHandlers = append(c.messageHandlers, handler)
	c.messageHandlersMutex.Unlock()

	return nil
}

func (c *channel) UnregisterRecv(handlerType string) error {
	c.messageHandlersMutex.Lock()
	defer c.messageHandlersMutex.Unlock()

	for i, mh := range c.messageHandlers {
		if mh.Type == handlerType {
			if len(c.messageHandlers) == 1 {
				c.messageHandlers = c.messageHandlers[:i]
				return nil
			}

			// If the underlying type changes to a pointer, this is a memory leak
			c.messageHandlers = append(c.messageHandlers[:i], c.messageHandlers[i+1:]...)
		}
	}

	return nil
}

func (c *channel) envelopeProto(
	recipient net.TransportIdentifier,
	sender *identity,
	message net.TaggedMarshaler,
) ([]byte, error) {
	payloadBytes, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	sig, err := c.clientIdentity.privKey.Sign(payloadBytes)
	if err != nil {
		return nil, err
	}

	senderIdentityBytes, err := sender.Marshal()
	if err != nil {
		return nil, err
	}

	var recipientIdentityBytes []byte
	if recipient != nil {
		recipientIdentity := &identity{id: peer.ID(recipient.(networkIdentity))}
		recipientIdentityBytes, err = recipientIdentity.Marshal()
		if err != nil {
			return nil, err
		}
	}

	return (&pb.Envelope{
		Payload:   payloadBytes,
		Signature: sig,
		Sender:    senderIdentityBytes,
		Recipient: recipientIdentityBytes,
		Type:      []byte(message.Type()),
	}).Marshal()
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

	if existingProtocolIdentifier, exists := c.transportToProtoIdentifiers[transportIdentifier]; exists {
		if existingProtocolIdentifier != protocolIdentifier {
			return fmt.Errorf(
				"protocol identifier in channel [%s] already associated with [%v]",
				c.name, transportIdentifier,
			)
		}
	}

	if existingTransportIdentifier, exists := c.protoToTransportIdentifiers[protocolIdentifier]; exists {
		if existingTransportIdentifier != transportIdentifier {
			return fmt.Errorf(
				"transport identifier in channel [%s] already associated with [%v]",
				c.name, protocolIdentifier,
			)
		}
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
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.subscription.Next(ctx)
			if err != nil {
				// TODO: handle error - different error types
				// result in different outcomes. Print err is very noisy.
				fmt.Println(err)
				continue
			}

			if err := c.processMessage(msg); err != nil {
				// TODO: handle error - different error types
				// result in different outcomes. Print err is very noisy.
				fmt.Println(err)
				continue
			}
		}
	}
}

func (c *channel) processMessage(message *floodsub.Message) error {
	var envelope pb.Envelope
	if err := proto.Unmarshal(message.Data, &envelope); err != nil {
		return err
	}

	// TODO: handle authentication, etc

	// The protocol type is on the envelope; let's pull that type
	// from our map of unmarshallers.
	unmarshaled, err := c.getUnmarshalingContainerByType(string(envelope.Type))
	if err != nil {
		return err
	}

	if err := unmarshaled.Unmarshal(envelope.GetPayload()); err != nil {
		return err
	}

	if envelope.Recipient != nil {
		// Construct an identifier from the Recipient
		recipientIdentifier := &identity{}
		if err := recipientIdentifier.Unmarshal(envelope.Recipient); err != nil {
			return err
		}

		if recipientIdentifier.id.String() != c.clientIdentity.id.String() {
			return fmt.Errorf(
				"message not for intended recipient %s",
				recipientIdentifier.id.String(),
			)
		}
	}

	// Construct an identifier from the sender
	senderIdentifier := &identity{}
	if err := senderIdentifier.Unmarshal(envelope.Sender); err != nil {
		return err
	}

	// Get the associated protocol identifier from an association map
	protocolIdentifier, err := c.getProtocolIdentifier(senderIdentifier)
	if err != nil {
		return err
	}

	// Fire a message back to the protocol
	protocolMessage := internal.BasicMessage(
		networkIdentity(senderIdentifier.id),
		protocolIdentifier,
		unmarshaled,
		string(envelope.Type),
	)

	return c.deliver(protocolMessage)
}

func (c *channel) getUnmarshalingContainerByType(envelopeType string) (net.TaggedUnmarshaler, error) {
	c.unmarshalersMutex.Lock()
	defer c.unmarshalersMutex.Unlock()

	unmarshaler, found := c.unmarshalersByType[envelopeType]
	if !found {
		return nil, fmt.Errorf(
			"couldn't find unmarshaler for type %s", envelopeType,
		)
	}

	return unmarshaler(), nil
}

func (c *channel) getProtocolIdentifier(senderIdentifier *identity) (net.ProtocolIdentifier, error) {
	c.identifiersMutex.Lock()
	defer c.identifiersMutex.Unlock()

	return c.transportToProtoIdentifiers[networkIdentity(senderIdentifier.id)], nil
}

func (c *channel) deliver(message net.Message) error {
	c.messageHandlersMutex.Lock()
	snapshot := make([]net.HandleMessageFunc, len(c.messageHandlers))
	copy(snapshot, c.messageHandlers)
	c.messageHandlersMutex.Unlock()

	for _, handler := range snapshot {
		go func(msg net.Message, handler net.HandleMessageFunc) {
			if err := handler.Handler(msg); err != nil {
				// TODO: handle error
				fmt.Println(err)
			}
			return
		}(message, handler)
	}

	return nil
}
