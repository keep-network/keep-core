package net

import (
	"context"
	"github.com/keep-network/keep-core/pkg/internal/pb"
	"github.com/keep-network/keep-core/pkg/operator"
)

// TransportIdentifier represents a protocol-level identifier. It is an opaque
// type to the network layer.
type TransportIdentifier interface {
	String() string
}

// MessageType is an alias representing the network message type.
type MessageType = string

// Message represents a message exchanged within the network layer. It carries
// a sender id for the transport layer and, if available, for the protocol
// layer. It also carries an unmarshaled payload.
type Message interface {
	TransportSenderID() TransportIdentifier
	SenderPublicKey() []byte

	Payload() interface{}

	Type() MessageType
	Seqno() uint64
}

// TypedPayload is a constraint type representing an abstract network payload
// that describes its own type. Despite the Message interface extends TypedPayload,
// there is no guarantee that Message.Payload() does the same. The TypedPayload
// should be used in the context of Message.Payload() only in cases where the
// underlying Message.Payload() type actually implements TypedPayload.
type TypedPayload interface {
	Type() MessageType
}

// TypedPayloads is a collection of items conforming the TypedPayload constraint.
type TypedPayloads[T TypedPayload] []T

// TypedMessagesSupplier is an abstract supplier of Messages whose MessageType
// returned by Message.Type() matches the concrete type behind the TypedPayload
// constraint.
type TypedMessagesSupplier[T TypedPayload] func(MessageType) []Message

// Payloads uses the TypedMessagesSupplier to extract TypedPayloads that
// represent a concrete type conforming the TypedPayload constraint.
func (tms TypedMessagesSupplier[T]) Payloads() TypedPayloads[T] {
	var template T

	payloads := make([]T, 0)
	for _, message := range tms(template.Type()) {
		payload, ok := message.Payload().(T)
		if !ok {
			continue
		}

		payloads = append(payloads, payload)
	}

	return payloads
}

// Deduplicate returns a deduplicated copy of tp. The deduplication logic
// determines duplicates using the provided keyFn that is used to produce
// a key for the given payload. Payloads using the same key are considered
// equal. Only the first payload with the given key is put to the resulting
// slice.
func (tp TypedPayloads[T]) Deduplicate(keyFn func(T) string) TypedPayloads[T] {
	keys := make(map[string]bool)
	result := make([]T, 0)

	for _, payload := range tp {
		key := keyFn(payload)
		if _, exists := keys[key]; !exists {
			keys[key] = true
			result = append(result, payload)
		}
	}

	return result
}

// TaggedMarshaler is an interface that includes the proto.Marshaler interface,
// but also provides a string type for the marshalable object.
type TaggedMarshaler interface {
	pb.Marshaler
	Type() string
}

// Provider represents an entity that can provide network access.
//
// Providers expose the ability to get a named BroadcastChannel, the ability to
// return a provider type, which is an informational string indicating what type
// of provider this is, the list of IP addresses on which it can listen, and
// known peers from peer discovery mechanims.
type Provider interface {
	// ID returns provider identifier.
	ID() TransportIdentifier
	// Type gives an information about provider type.
	Type() string

	// BroadcastChannelFor provides a broadcast channel instance for given
	// channel name.
	BroadcastChannelFor(name string) (BroadcastChannel, error)

	// ConnectionManager returns the connection manager used by the provider.
	ConnectionManager() ConnectionManager

	// CreateTransportIdentifier creates a transport identifier based on the
	// provided operator public key.
	CreateTransportIdentifier(
		operatorPublicKey *operator.PublicKey,
	) (TransportIdentifier, error)

	// BroadcastChannelForwarderFor creates a message relay for given channel name.
	BroadcastChannelForwarderFor(name string)
}

// ConnectionManager is an interface which exposes peers a client is connected
// to, and their individual identities, so that a client may forcibly disconnect
// from any given connected peer.
type ConnectionManager interface {
	ConnectedPeers() []string
	ConnectedPeersAddrInfo() map[string][]string
	GetPeerPublicKey(connectedPeer string) (*operator.PublicKey, error)
	DisconnectPeer(connectedPeer string)

	// AddrStrings returns all listen addresses of the provider.
	AddrStrings() []string

	IsConnected(address string) bool
}

// TaggedUnmarshaler is an interface that includes the proto.Unmarshaler
// interface, but also provides a string type for the unmarshalable object. The
// Type() method is expected to be invokable on a just-initialized instance of
// the unmarshaler (i.e., before unmarshaling is completed).
type TaggedUnmarshaler interface {
	pb.Unmarshaler
	Type() string
}

// BroadcastChannel represents a named pubsub channel. It allows group members
// to broadcast and receive messages. BroadcastChannel implements strategy
// for the retransmission of broadcast messages and handle duplicates before
// passing the received message to the client.
type BroadcastChannel interface {
	// Name returns the name of this broadcast channel.
	Name() string
	// Send function publishes a message m to the channel. Message m needs to
	// conform to the marshalling interface. Message will be periodically
	// retransmitted by the channel for the lifetime of the provided context.
	Send(ctx context.Context, m TaggedMarshaler) error
	// Recv installs a message handler that will receive messages from the
	// channel for the entire lifetime of the provided context.
	// When the context is done, handler is automatically unregistered and
	// receives no more messages. Already received message retransmissions are
	// filtered out before calling the handler.
	Recv(ctx context.Context, handler func(m Message))
	// SetUnmarshaler set an unmarshaler that will unmarshal a given
	// type to a concrete object that can be passed to and understood by any
	// registered message handling functions. The unmarshaler should be a
	// function that returns a fresh object of type proto.TaggedUnmarshaler,
	// ready to read in the bytes for an object marked as tpe.
	//
	// The string type associated with the unmarshaler is the result of calling
	// Type() on a raw unmarshaler.
	SetUnmarshaler(unmarshaler func() TaggedUnmarshaler)
	// SetFilter registers a broadcast channel filter which will be used
	// to determine if given broadcast channel message should be processed
	// by the receivers.
	SetFilter(filter BroadcastChannelFilter) error
}

// BroadcastChannelFilter represents a filter which determine if the incoming
// message should be processed by the receivers. It takes the message author's
// public key as its argument and returns true if the message should be
// processed or false otherwise.
type BroadcastChannelFilter func(*operator.PublicKey) bool

// Firewall represents a set of rules the remote peer has to conform to so that
// a connection with that peer can be approved.
type Firewall interface {

	// Validate takes the remote peer public key and executes all the checks
	// needed to decide whether the connection with the remote peer can be
	// approved.
	// If expectations are not met, this function should return an error
	// describing what is wrong.
	Validate(remotePeerPublicKey *operator.PublicKey) error
}
