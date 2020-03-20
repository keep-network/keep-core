package libp2p

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	libp2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
)

const signPrefix = "keep-unicast:"

func signMessage(
	message *pb.UnicastNetworkMessage,
	privateKey libp2pcrypto.PrivKey,
) error {
	messageCopy := *message
	messageCopy.Signature = nil

	bytes, err := messageCopy.Marshal()
	if err != nil {
		return err
	}

	bytes = withSignPrefix(bytes)

	signature, err := privateKey.Sign(bytes)
	if err != nil {
		return err
	}

	message.Signature = signature

	return nil
}

func verifyMessageSignature(
	message *pb.UnicastNetworkMessage,
	publicKey libp2pcrypto.PubKey,
) error {
	messageCopy := *message
	messageCopy.Signature = nil

	bytes, err := messageCopy.Marshal()
	if err != nil {
		return err
	}

	bytes = withSignPrefix(bytes)

	valid, err := publicKey.Verify(bytes, message.Signature)
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("invalid message signature")
	}

	return nil
}

func withSignPrefix(bytes []byte) []byte {
	return append([]byte(signPrefix), bytes...)
}
