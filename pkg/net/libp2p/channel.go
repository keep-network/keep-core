package libp2p

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/btcsuite/btcd/btcec"

	"github.com/gogo/protobuf/proto"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/keep-network/keep-core/pkg/net/internal"
	"github.com/keep-network/keep-core/pkg/net/key"
	"github.com/keep-network/keep-core/pkg/net/retransmission"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const subscriptionWorkersCount = 32

type channel struct {
	// channel-scoped atomic counter for sequence numbers
	//
	// Must be declared at the top of the struct!
	// See: https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	counter uint64

	name string

	clientIdentity *identity
	peerStore      peerstore.Peerstore

	pubsubMutex sync.Mutex
	pubsub      *pubsub.PubSub

	subscription *pubsub.Subscription

	messageHandlersMutex sync.Mutex
	messageHandlers      []*messageHandler

	unmarshalersMutex  sync.Mutex
	unmarshalersByType map[string]func() net.TaggedUnmarshaler

	retransmissionTicker *retransmission.Ticker
}

type messageHandler struct {
	ctx     context.Context
	channel chan retransmission.NetworkMessage
}

func (c *channel) nextSeqno() uint64 {
	return atomic.AddUint64(&c.counter, 1)
}

func (c *channel) Name() string {
	return c.name
}

func (c *channel) Send(ctx context.Context, message net.TaggedMarshaler) error {
	messageProto, err := c.messageProto(message)
	if err != nil {
		return err
	}

	retransmission.ScheduleRetransmissions(
		ctx,
		c.retransmissionTicker,
		messageProto,
		c.publishToPubSub,
	)

	return c.publishToPubSub(messageProto)
}

func (c *channel) Recv(ctx context.Context, handler func(m net.Message)) {
	messageHandler := &messageHandler{
		ctx:     ctx,
		channel: make(chan retransmission.NetworkMessage),
	}

	c.messageHandlersMutex.Lock()
	c.messageHandlers = append(c.messageHandlers, messageHandler)
	c.messageHandlersMutex.Unlock()

	handleWithRetransmissions := retransmission.WithRetransmissionSupport(handler)

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Debug("context is done, removing handler")
				c.removeHandler(messageHandler)
				return

			case msg := <-messageHandler.channel:
				// Go language specification says that if one or more of the
				// communications in the select statement can proceed, a single
				// one that will proceed is chosen via a uniform pseudo-random
				// selection.
				// Thus, it can happen this communication is called when ctx is
				// already done. Since we guarantee in the network channel API
				// that handler is not called after ctx is done (client code
				// could e.g. perform come cleanup), we need to double-check
				// the context state here.
				if messageHandler.ctx.Err() != nil {
					continue
				}

				handleWithRetransmissions(msg)
			}
		}
	}()
}

func (c *channel) removeHandler(handler *messageHandler) {
	c.messageHandlersMutex.Lock()
	defer c.messageHandlersMutex.Unlock()

	for i, h := range c.messageHandlers {
		if h.channel == handler.channel {
			c.messageHandlers[i] = c.messageHandlers[len(c.messageHandlers)-1]
			c.messageHandlers = c.messageHandlers[:len(c.messageHandlers)-1]
		}
	}
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
	for i := 0; i < subscriptionWorkersCount; i++ {
		go c.subscriptionWorker(ctx)
	}

	<-ctx.Done()
	c.subscription.Cancel()
}

func (c *channel) subscriptionWorker(ctx context.Context) {
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

	return c.processContainerMessage(pubsubMessage.GetFrom(), messageProto)
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

	fingerprint := retransmission.CalculateFingerprint(
		senderIdentifier.id,
		message.GetPayload(),
	)

	c.deliver(retransmission.NewNetworkMessage(
		internal.BasicMessage(
			senderIdentifier.id,
			unmarshaled,
			string(message.Type),
			key.Marshal(networkKey),
			c.nextSeqno(),
		),
		fingerprint,
		message.Retransmission,
	))

	return nil
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

func (c *channel) deliver(message retransmission.NetworkMessage) {
	c.messageHandlersMutex.Lock()
	snapshot := make([]*messageHandler, len(c.messageHandlers))
	copy(snapshot, c.messageHandlers)
	c.messageHandlersMutex.Unlock()

	for _, handler := range snapshot {
		go func(message retransmission.NetworkMessage, handler *messageHandler) {
			select {
			case handler.channel <- message:
			// Nothing to do here; we block until the message is handled
			// or until the context gets closed.
			// This way we don't lose any message but also don't stay
			// with any dangling goroutines if there is no longer anyone
			// to receive messages.
			case <-handler.ctx.Done():
				return
			}
		}(message, handler)
	}
}

func (c *channel) SetFilter(filter net.BroadcastChannelFilter) error {
	c.pubsubMutex.Lock()
	defer c.pubsubMutex.Unlock()

	c.pubsub.UnregisterTopicValidator(c.name)

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
