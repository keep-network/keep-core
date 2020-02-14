package libp2p

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
)

func signMessage(
	message *pb.NetworkMessage,
	privateKey libp2pcrypto.PrivKey,
) error {
	bytes, err := message.Marshal()
	if err != nil {
		return err
	}

	signature, err := privateKey.Sign(bytes)
	if err != nil {
		return err
	}

	message.Signature = signature

	return nil
}

func verifyMessageSignature(
	message *pb.NetworkMessage,
	publicKey libp2pcrypto.PubKey,
) error {
	messageCopy := *message
	messageCopy.Signature = nil

	bytes, err := messageCopy.Marshal()
	if err != nil {
		return err
	}

	valid, err := publicKey.Verify(bytes, message.Signature)
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("invalid message signature")
	}

	return nil
}
