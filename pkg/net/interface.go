package net

import (
	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/gogo/protobuf/proto"
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

// ClientIdentifier represents the identity of a recipient for a message.
type ClientIdentifier string

// HandleMessageFunc is the type of function called for each Message m furnished
// by the BroadcastChannel. If there is a problem handling the Message, the
// incoming error will describe the problem and the function can decide how to
// handle that error. If an error is returned, processing stops.
type HandleMessageFunc func(m interface{}) error

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
	SendTo(recipient ClientIdentifier, m TaggedMarshaler) error

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
