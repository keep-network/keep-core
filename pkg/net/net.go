package net

import (
	"context"
	"crypto/ecdsa"

	"github.com/gogo/protobuf/proto"
	"github.com/keep-network/keep-core/pkg/net/key"
)

// TransportIdentifier represents a protocol-level identifier. It is an opaque
// type to the network layer.
type TransportIdentifier interface {
	String() string
}

// Message represents a message exchanged within the network layer. It carries
// a sender id for the transport layer and, if available, for the protocol
// layer. It also carries an unmarshaled payload.
type Message interface {
	TransportSenderID() TransportIdentifier
	Payload() interface{}
	Type() string
	SenderPublicKey() []byte
}

// HandleMessageFunc is the type of function called for each Message m furnished
// by the BroadcastChannel. If there is a problem handling the Message, the
// incoming error will describe the problem and the function can decide how to
// handle that error. If an error is returned, processing stops.
type HandleMessageFunc struct {
	Type    string
	Handler func(m Message) error
}

// TaggedMarshaler is an interface that includes the proto.Marshaler interface,
// but also provides a string type for the marshalable object.
type TaggedMarshaler interface {
	proto.Marshaler
	Type() string
}

// Provider represents an entity that can provide network access.
//
// Providers expose the ability to get a named BroadcastChannel, the ability to
// return a provider type, which is an informational string indicating what type
// of provider this is, the list of IP addresses on which it can listen, and
// known peers from peer discovery mechanims.
type Provider interface {
	ID() TransportIdentifier

	ChannelFor(name string) (BroadcastChannel, error)
	Type() string
	AddrStrings() []string

	// All known peers from the underlying PeerStore. This may include
	// peers we're not directly connected to.
	Peers() []string

	ConnectionManager() ConnectionManager
}

// ConnectionManager is an interface which exposes peers a client is connected
// to, and their individual identities, so that a client may forcibly disconnect
// from any given connected peer.
type ConnectionManager interface {
	ConnectedPeers() []string
	GetPeerPublicKey(connectedPeer string) (*key.NetworkPublic, error)
	DisconnectPeer(connectedPeer string)
}

// TaggedUnmarshaler is an interface that includes the proto.Unmarshaler
// interface, but also provides a string type for the unmarshalable object. The
// Type() method is expected to be invokable on a just-initialized instance of
// the unmarshaler (i.e., before unmarshaling is completed).
type TaggedUnmarshaler interface {
	proto.Unmarshaler
	Type() string
}

// BroadcastChannelFilter represents a filter which determine if the incoming
// message should be processed by the receivers. It takes the message author's
// public key as its argument and returns true if the message should be
// processed or false otherwise.
type BroadcastChannelFilter func(*ecdsa.PublicKey) bool

// BroadcastChannel represents a named pubsub channel. It allows Group Members
// to send messages on the channel (via Send), and to access a low-level receive chan
// that furnishes messages sent onto the BroadcastChannel. Messages are not
// guaranteed to be ordered at the pubsub level, though they will be at the
// underlying network protocol (ie. tcp, quic).
type BroadcastChannel interface {
	// Name returns the name of this broadcast channel.
	Name() string
	// Given a message m that can marshal itself to protobuf, broadcast m to
	// members of the group through the BroadcastChannel with possible
	// retransmissions for the lifetime of the provided context.
	//
	// If the BroadcastChannel implementation does not guarantee message
	// delivery it must implement message retransmissions.
	// All the BroadcastChannel implementations performing retransmissions must
	// filter out retransmissions of previously received message and ensure the
	// message is received by handlers registered in the channel only one time.
	//
	// If the BroadcastChannel implementation does guarantee message delivery,
	// retransmissions should not be performed.
	Send(ctx context.Context, m TaggedMarshaler) error
	// Recv takes a HandleMessageFunc and returns an error. This function should
	// be retried.
	Recv(h HandleMessageFunc) error
	// UnregisterRecv takes the type of HandleMessageFunc and returns an
	// error. This function should be defered.
	UnregisterRecv(handlerType string) error
	// RegisterUnmarshaler registers an unmarshaler that will unmarshal a given
	// type to a concrete object that can be passed to and understood by any
	// registered message handling functions. The unmarshaler should be a
	// function that returns a fresh object of type proto.TaggedUnmarshaler,
	// ready to read in the bytes for an object marked as tpe.
	//
	// The string type associated with the unmarshaler is the result of calling
	// Type() on a raw unmarshaler.
	RegisterUnmarshaler(unmarshaler func() TaggedUnmarshaler) error
	// AddFilter registers a broadcast channel filter which will be used
	// to determine if given broadcast channel message should be processed
	// by the receivers.
	AddFilter(filter BroadcastChannelFilter) error
}
