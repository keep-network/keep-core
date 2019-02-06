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
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
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
}

func (c *channel) Name() string {
	return c.name
}

func (c *channel) Send(message net.TaggedMarshaler) error {
	// Transform net.TaggedMarshaler to a protobuf message, sign, and wrap
	// in an envelope.
	envelopeBytes, err := c.envelopeProto(message)
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

	return (&pb.NetworkMessage{
		Payload: payloadBytes,
		Sender:  senderIdentityBytes,
		Type:    []byte(message.Type()),
	}).Marshal()
}

func (c *channel) sealEnvelope(
	message net.TaggedMarshaler,
) (*pb.NetworkEnvelope, error) {
	messageBytes, err := c.messageProto(message)
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
	message net.TaggedMarshaler,
) ([]byte, error) {
	envelope, err := c.sealEnvelope(message)
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

	protocolIdentifier := senderIdentifier

	// Fire a message back to the protocol.
	protocolMessage := internal.BasicMessage(
		protocolIdentifier.id,
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
