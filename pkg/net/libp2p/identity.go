package libp2p

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/keep-network/keep-core/pkg/operator"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"

	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

// identity represents a group member's network level identity. It
// implements the net.TransportIdentifier interface. A valid group member will
// generate or provide a keypair, which will correspond to a network ID.
//
// Consumers of the net package require an ID to register with protocol level
// IDs, as well as a public key for authentication.
type identity struct {
	id      peer.ID
	pubKey  libp2pcrypto.PubKey
	privKey libp2pcrypto.PrivKey
}

type networkIdentity peer.ID

func generateIdentity() (*identity, error) {
	operatorPrivateKey, _, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot generate operator key pair: [%v]",
			err,
		)
	}

	networkPrivateKey, _, err := operatorPrivateKeyToNetworkKeyPair(
		operatorPrivateKey,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot convert operator private key to network key pair: [%v]",
			err,
		)
	}

	return createIdentity(networkPrivateKey)
}

func createIdentity(privateKey libp2pcrypto.PrivKey) (*identity, error) {
	peerID, err := peer.IDFromPublicKey(privateKey.GetPublic())
	if err != nil {
		return nil, fmt.Errorf(
			"could not transform public key to peer's identity: [%v]",
			err,
		)
	}

	return &identity{peerID, privateKey.GetPublic(), privateKey}, nil
}

func (ni networkIdentity) String() string {
	return peer.ID(ni).String()
}

func (i *identity) Marshal() ([]byte, error) {
	var (
		err error
	)

	pubKey := i.pubKey
	if pubKey == nil {
		pubKey, err = i.id.ExtractPublicKey()
		if err != nil {
			return nil, err
		}
	}

	pubKeyBytes, err := libp2pcrypto.MarshalPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(&pb.Identity{PubKey: pubKeyBytes})
}

func (i *identity) Unmarshal(bytes []byte) error {
	var (
		err        error
		pid        peer.ID
		pbIdentity pb.Identity
	)

	if err = proto.Unmarshal(bytes, &pbIdentity); err != nil {
		return fmt.Errorf("unmarshalling failed: [%v]", err)
	}
	i.pubKey, err = libp2pcrypto.UnmarshalPublicKey(pbIdentity.PubKey)
	if err != nil {
		return err
	}
	pid, err = peer.IDFromPublicKey(i.pubKey)
	if err != nil {
		return fmt.Errorf(
			"failed to generate valid libp2p identity: [%v]",
			err,
		)
	}
	i.id = pid

	return nil
}
