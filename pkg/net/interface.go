package net

import (
	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/gogo/protobuf/proto"

	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

// GroupIdentity contains the Group's public key as created by the dkg process
// and a list of Members that belong to the group.
// TODO: move to a more appropriate package; revise which fields we need
type GroupIdentity struct {
	// Group names are isomorphic to channel names
	// Channel names are Keccak(StakingPubKey1 || ... || StakingPubKeyN) of
	// all valid group members.
	Name string
	// Public key for the group; created through DKG.
	// Verified from the on-chain Group Registry
	GroupPublicKey *bls.PublicKey
	// The final list of qualified Group Members; empty if not yet computed
	// Verified from the on-chain Group Registry
	Members []bls.PublicKey
}

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
	// does not have an associated ClientIdentifier, or if it is neither a
	// ProtocolIdentifier nor a ClientIdentifier.
	SendTo(recipientIdentifier interface{}, m TaggedMarshaler) error

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

// An ID corresponds to the identification of a member in a peer-to-peer network.
type ID peer.ID

// PubKey is a type alias for the underlying PublicKey implementation we choose.
type PubKey = ci.PubKey

// Identity represents a group member's network level identity. A valid group
// member will generate or provide a keypair, which will correspond to a network
// ID. Consumers of the net package require an ID to register with protocol level
// ID's, as well as a public key for authentication.
type Identity interface {
	ID() ID
	PubKeyFromID(ID) (PubKey, error)
}
