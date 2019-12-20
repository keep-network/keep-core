package libp2p

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"sync"

	"github.com/btcsuite/btcd/btcec"

	"github.com/gogo/protobuf/proto"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/internal"
	"github.com/keep-network/keep-core/pkg/net/key"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
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

	retransmitter *retransmitter
}

func (c *channel) Name() string {
	return c.name
}

func (c *channel) Send(message net.TaggedMarshaler, ctx ...context.Context) error {
	// Transform net.TaggedMarshaler to a protobuf message
	messageProto, err := c.messageProto(message)
	if err != nil {
		return err
	}

	if len(ctx) > 0 {
		c.retransmitter.scheduleRetransmissions(ctx[0], messageProto, c.publishToPubSub)
	}
	return c.publishToPubSub(messageProto)
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

	removedCount := 0

	// updated slice shares the same backing array and capacity as the original,
	// so the storage is reused for the filtered slice.
	updated := c.messageHandlers[:0]

	for _, mh := range c.messageHandlers {
		if mh.Type != handlerType {
			updated = append(updated, mh)
		} else {
			removedCount++
		}
	}

	c.messageHandlers = updated[:len(c.messageHandlers)-removedCount]

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
) (*pb.NetworkMessage, error) {
	payloadBytes, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	senderIdentityBytes, err := c.clientIdentity.Marshal()
	if err != nil {
		return nil, err
	}

	return &pb.NetworkMessage{
		Payload: payloadBytes,
		Sender:  senderIdentityBytes,
		Type:    []byte(message.Type()),
	}, nil
}

func (c *channel) publishToPubSub(message *pb.NetworkMessage) error {
	messageBytes, err := message.Marshal()
	if err != nil {
		return err
	}

	c.pubsubMutex.Lock()
	defer c.pubsubMutex.Unlock()

	return c.pubsub.Publish(c.name, messageBytes)
}

func (c *channel) handleMessages(ctx context.Context) {
	defer c.subscription.Cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			message, err := c.subscription.Next(ctx)
			if err != nil {
				// TODO: handle error - different error types
				// result in different outcomes. Print err is very noisy.
				logger.Error(err)
				continue
			}

			// Every message should be independent from any other message.
			go func(msg *pubsub.Message) {
				if err := c.processPubsubMessage(msg); err != nil {
					// TODO: handle error - different error types
					// result in different outcomes. Print err is very noisy.
					logger.Error(err)
					return
				}
			}(message)
		}
	}
}

func (c *channel) processPubsubMessage(pubsubMessage *pubsub.Message) error {
	var messageProto pb.NetworkMessage
	if err := proto.Unmarshal(pubsubMessage.Data, &messageProto); err != nil {
		return err
	}

	onFirstTimeReceived := func() error {
		return c.processContainerMessage(pubsubMessage.GetFrom(), messageProto)
	}

	return c.retransmitter.receive(&messageProto, onFirstTimeReceived)
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

	networkKey := key.Libp2pKeyToNetworkKey(senderIdentifier.pubKey)
	if networkKey == nil {
		return fmt.Errorf(
			"sender [%v] with key [%v] is not of correct type",
			senderIdentifier.id,
			senderIdentifier.pubKey,
		)
	}

	// Fire a message back to the protocol.
	protocolMessage := internal.BasicMessage(
		senderIdentifier.id,
		unmarshaled,
		string(message.Type),
		key.Marshal(networkKey),
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
				logger.Error(err)
				return
			}
		}(message, handler)
	}

	return nil
}

func (c *channel) AddFilter(filter net.BroadcastChannelFilter) error {
	c.pubsubMutex.Lock()
	defer c.pubsubMutex.Unlock()

	return c.pubsub.RegisterTopicValidator(c.name, createTopicValidator(filter))
}

func createTopicValidator(filter net.BroadcastChannelFilter) pubsub.Validator {
	return func(_ context.Context, _ peer.ID, message *pubsub.Message) bool {
		authorPublicKey, err := extractPublicKey(message.GetFrom())
		if err != nil {
			logger.Warningf(
				"could not retrieve message author public key: [%v]",
				err,
			)
			return false
		}
		return filter(authorPublicKey)
	}
}

func extractPublicKey(peer peer.ID) (*ecdsa.PublicKey, error) {
	publicKey, err := peer.ExtractPublicKey()
	if err != nil {
		return nil, err
	}

	secp256k1PublicKey, ok := publicKey.(*crypto.Secp256k1PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is of type other than Secp256k1")
	}

	return (*btcec.PublicKey)(secp256k1PublicKey).ToECDSA(), nil
}
