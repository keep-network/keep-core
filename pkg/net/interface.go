package net

import (
	"github.com/gogo/protobuf/proto"
)

// TransportIdentifier represents the identity of a participant at the transport
// layer (e.g., libp2p).
type TransportIdentifier interface {
	// Returns a string name of the network provider. Expected to be purely
	// informational.
	ProviderName() string
}

// ProtocolIdentifier represents a protocol-level identifier. It is an opaque
// type to the network layer.
type ProtocolIdentifier interface{}

// Message represents a message exchanged within the network layer. It carries
// a sender id for the transport layer and, if available, for the protocol
// layer. It also carries an unmarshaled payload.
type Message interface {
	TransportSenderID() TransportIdentifier
	ProtocolSenderID() ProtocolIdentifier
	Payload() interface{}
}

// HandleMessageFunc is the type of function called for each Message m furnished
// by the BroadcastChannel. If there is a problem handling the Message, the
// incoming error will describe the problem and the function can decide how to
// handle that error. If an error is returned, processing stops.
type HandleMessageFunc func(m Message) error

// TaggedMarshaler is an interface that includes the proto.Marshaler interface,
// but also provides a string type for the marshalable object.
type TaggedMarshaler interface {
	proto.Marshaler
	Type() string
}

// Provider represents an entity that can provide network access.
//
// Currently only 3 methods are exposed by providers: the ability to get a
// named BroadcastChannel, the ability to return a provider type, which is
// an informational string indicating what type of provider this is, and
// ListenIPAddresses for this node.
type Provider interface {
	ChannelFor(name string) (BroadcastChannel, error)
	Type() string
	ListenIPAddresses(port int) ([]string, error)
}

// TaggedUnmarshaler is an interface that includes the proto.Unmarshaler
// interface, but also provides a string type for the unmarshalable object. The
// Type() method is expected to be invokable on a just-initialized instance of
// the unmarshaler (i.e., before unmarshaling is completed).
type TaggedUnmarshaler interface {
	proto.Unmarshaler
	Type() string
}

// BroadcastChannel represents a named pubsub channel. It allows Group Members
// to send messages on the channel (via Send), and to access a low-level receive chan
// that furnishes messages sent onto the BroadcastChannel.
type BroadcastChannel interface {
	// Name returns the name of this broadcast channel.
	Name() string

	// Given a message m that can marshal itself to protobuf, broadcast m to
	// members of the Group through the BroadcastChannel.
	Send(m TaggedMarshaler) error
	// Given a recipient and a message m that can marshal itself to protobouf,
	// send the message to the recipient over the broadcast channel such that
	// only the recipient can understand it.
	//
	// The recipient should be a ProtocolIdentifier registered using
	// RegisterIdentifier, or a ClientIdentifier used by the network layer.
	//
	// Returns an error if the recipient identifier is a ProtocolIdentifier that
	// does not have an associated ClientIdentifier, or if it is not a
	// ProtocolIdentifier.
	SendTo(recipientIdentifier ProtocolIdentifier, m TaggedMarshaler) error

	// RegisterIdentifier associates the given network identifier with a
	// protocol-specific identifier that will be passed to the receiving code
	// in HandleMessageFunc.
	//
	// Returns an error if either identifier already has an association for
	// this channel.
	RegisterIdentifier(
		networkIdentifier TransportIdentifier,
		protocolIdentifier ProtocolIdentifier,
	) error

	// Recv takes a HandleMessageFunc and returns an error. This function should
	// be retried.
	Recv(h HandleMessageFunc) error

	// RegisterUnmarshaler registers an unmarshaler that will unmarshal a given
	// type to a concrete object that can be passed to and understood by any
	// registered message handling functions. The unmarshaler should be a
	// function that returns a fresh object of type proto.TaggedUnmarshaler,
	// ready to read in the bytes for an object marked as tpe.
	//
	// The string type associated with the unmarshaler is the result of calling
	// Type() on a raw unmarshaler.
	RegisterUnmarshaler(unmarshaler func() TaggedUnmarshaler) error
}
