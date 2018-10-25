package handshake

import (
	"encoding/binary"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
)

// Marshal converts this act1Message to a byte array suitable for network
// communication.
func (am *act1Message) Marshal() ([]byte, error) {
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, am.nonce1)
	return (&pb.Act1Message{Nonce: nonceBytes}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a act1Message.
func (am *act1Message) Unmarshal(bytes []byte) error {
	pbAct1 := pb.Act1Message{}
	if err := pbAct1.Unmarshal(bytes); err != nil {
		return err
	}
	am.nonce1 = binary.LittleEndian.Uint64(pbAct1.Nonce)

	return nil
}

// Marshal converts this act2Message to a byte array suitable for network
// communication.
func (am *act2Message) Marshal() ([]byte, error) {
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, am.nonce2)
	return (&pb.Act2Message{Nonce: nonceBytes, Challenge: am.challenge[:]}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a act2Message.
func (am *act2Message) Unmarshal(bytes []byte) error {
	pbAct2 := pb.Act2Message{}
	if err := pbAct2.Unmarshal(bytes); err != nil {
		return err
	}

	am.nonce2 = binary.LittleEndian.Uint64(pbAct2.Nonce)
	copy(am.challenge[:], pbAct2.Challenge[:32])

	return nil
}
