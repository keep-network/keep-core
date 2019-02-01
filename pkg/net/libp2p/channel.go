package libp2p

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/internal"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-pubsub"
)

type channel struct {
	name string

	clientIdentity *identity
	peerStore      peerstore.Peerstore

	pubsubMutex sync.Mutex
	pubsub      *pubsub.PubSub

	subscription *pubsub.Subscription

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
	return c.doSend(nil, message)
}

func (c *channel) SendTo(
	recipientIdentifier net.ProtocolIdentifier,
	message net.TaggedMarshaler,
) error {
	return c.doSend(recipientIdentifier, message)
}

// doSend attempts to send a message, from a sender, to all members of a
// broadcastChannel, or optionally to a specific recipient. If recipient
// is nil (the typical case), then all messages of the broadcast channel
// should receive the message. Otherwise, given a valid recipient, we will
// address the message specifically to them.
func (c *channel) doSend(
	recipient net.ProtocolIdentifier,
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
	// Transform net.TaggedMarshaler to a protobuf message, sign, and wrap
	// in an envelope.
	envelopeBytes, err := c.envelopeProto(transportRecipient, message)
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

	handlers := 0
	for i, mh := range c.messageHandlers {
		// filter out the handlerType
		if mh.Type != handlerType {
			c.messageHandlers[i] = mh
			handlers++
		}
	}
	c.messageHandlers = c.messageHandlers[:handlers]

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

func (c *channel) messageProto(
	recipient net.TransportIdentifier,
	message net.TaggedMarshaler,
) ([]byte, error) {
	payloadBytes, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	senderIdentityBytes, err := c.clientIdentity.Marshal()
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

	return (&pb.NetworkMessage{
		Payload:   payloadBytes,
		Sender:    senderIdentityBytes,
		Recipient: recipientIdentityBytes,
		Type:      []byte(message.Type()),
	}).Marshal()
}

func (c *channel) sealEnvelope(
	recipient net.TransportIdentifier,
	message net.TaggedMarshaler,
) (*pb.NetworkEnvelope, error) {
	messageBytes, err := c.messageProto(recipient, message)
	if err != nil {
		return nil, err
	}
	signature, err := c.sign(messageBytes)
	if err != nil {
		return nil, err
	}

	return &pb.NetworkEnvelope{
		Message:   messageBytes,
		Signature: signature,
	}, nil
}

func (c *channel) envelopeProto(
	recipient net.TransportIdentifier,
	message net.TaggedMarshaler,
) ([]byte, error) {
	envelope, err := c.sealEnvelope(recipient, message)
	if err != nil {
		return nil, err
	}

	return envelope.Marshal()
}

func (c *channel) sign(messageBytes []byte) ([]byte, error) {
	return c.clientIdentity.privKey.Sign(messageBytes)
}

func (c *channel) verify(sender peer.ID, messageBytes []byte, signature []byte) error {
	return verifyEnvelope(sender, messageBytes, signature)
}

func verifyEnvelope(sender peer.ID, messageBytes []byte, signature []byte) error {
	pubKey, err := sender.ExtractPublicKey()
	if err != nil {
		return fmt.Errorf(
			"failed to extract public key from peer [%v]",
			sender,
		)
	}

	ok, err := pubKey.Verify(messageBytes, signature)
	if err != nil {
		return fmt.Errorf(
			"failed to verify signature [0x%v] for sender [%v] with err [%v]",
			hex.EncodeToString(signature),
			sender.Pretty(),
			err,
		)
	}

	if !ok {
		return fmt.Errorf(
			"invalid signature [0x%v] on message from sender [%v] ",
			hex.EncodeToString(signature),
			sender.Pretty(),
		)
	}

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

			if err := c.processPubsubMessage(msg); err != nil {
				// TODO: handle error - different error types
				// result in different outcomes. Print err is very noisy.
				fmt.Println(err)
				continue
			}
		}
	}
}

func (c *channel) processPubsubMessage(pubsubMessage *pubsub.Message) error {
	var envelope pb.NetworkEnvelope
	if err := proto.Unmarshal(pubsubMessage.Data, &envelope); err != nil {
		return err
	}

	if err := c.verify(
		pubsubMessage.GetFrom(),
		envelope.GetMessage(),
		envelope.GetSignature(),
	); err != nil {
		return err
	}

	var protoMessage pb.NetworkMessage
	if err := proto.Unmarshal(envelope.Message, &protoMessage); err != nil {
		return err
	}

	return c.processContainerMessage(pubsubMessage.GetFrom(), protoMessage)
}

func (c *channel) processContainerMessage(
	proposedSender peer.ID,
	message pb.NetworkMessage,
) error {
	// The protocol type is on the envelope; let's pull that type
	// from our map of unmarshallers.
	unmarshaled, err := c.getUnmarshalingContainerByType(string(message.Type))
	if err != nil {
		return err
	}

	if err := unmarshaled.Unmarshal(message.GetPayload()); err != nil {
		return err
	}

	// Construct an identifier from the sender.
	senderIdentifier := &identity{}
	if err := senderIdentifier.Unmarshal(message.Sender); err != nil {
		return err
	}

	// Ensure the sender wasn't tampered by:
	//     Test that the proposed sender (outer layer) matches the
	//     sender identifier we grab from the message (inner layer).
	if proposedSender != senderIdentifier.id {
		return fmt.Errorf(
			"Outer layer sender [%v] does not match inner layer sender [%v]",
			proposedSender,
			senderIdentifier,
		)
	}

	// Get the associated protocol identifier from an association map.
	protocolIdentifier, err := c.getProtocolIdentifier(senderIdentifier)
	if err != nil {
		return err
	}

	if message.Recipient != nil {
		// Construct an identifier from the Recipient.
		recipientIdentifier := &identity{}
		if err := recipientIdentifier.Unmarshal(message.Recipient); err != nil {
			return err
		}

		if recipientIdentifier.id.String() != c.clientIdentity.id.String() {
			return fmt.Errorf(
				"message not for intended recipient %s",
				recipientIdentifier.id.String(),
			)
		}
	}

	// Fire a message back to the protocol.
	protocolMessage := internal.BasicMessage(
		networkIdentity(senderIdentifier.id),
		protocolIdentifier,
		unmarshaled,
		string(message.Type),
	)

	return c.deliver(protocolMessage)
}

func (c *channel) getUnmarshalingContainerByType(messageType string) (net.TaggedUnmarshaler, error) {
	c.unmarshalersMutex.Lock()
	defer c.unmarshalersMutex.Unlock()

	unmarshaler, found := c.unmarshalersByType[messageType]
	if !found {
		return nil, fmt.Errorf(
			"couldn't find unmarshaler for type %s", messageType,
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
		}(message, handler)
	}

	return nil
}
