package handshake

import (
	"encoding/binary"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
)

// Marshal converts this Act1Message to a byte array suitable for network
// communication.
func (am *Act1Message) Marshal() ([]byte, error) {
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, am.nonce1)
	return (&pb.Act1Message{Nonce: nonceBytes}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a Act1Message.
func (am *Act1Message) Unmarshal(bytes []byte) error {
	pbAct1 := pb.Act1Message{}
	if err := pbAct1.Unmarshal(bytes); err != nil {
		return err
	}
	am.nonce1 = binary.LittleEndian.Uint64(pbAct1.Nonce)

	return nil
}

// Marshal converts this Act2Message to a byte array suitable for network
// communication.
func (am *Act2Message) Marshal() ([]byte, error) {
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, am.nonce2)
	return (&pb.Act2Message{Nonce: nonceBytes, Challenge: am.challenge[:]}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a Act2Message.
func (am *Act2Message) Unmarshal(bytes []byte) error {
	pbAct2 := pb.Act2Message{}
	if err := pbAct2.Unmarshal(bytes); err != nil {
		return err
	}
	am.nonce2 = binary.LittleEndian.Uint64(pbAct2.Nonce)
	copy(am.challenge[:], pbAct2.Challenge[:32])

	return nil
}

// Marshal converts this Act3Message to a byte array suitable for network
// communication.
func (am *Act3Message) Marshal() ([]byte, error) {
	return (&pb.Act3Message{Challenge: am.challenge[:]}).Marshal()
}

// Unmarshal converts a byte array produced by Marshal to a Act3Message.
func (am *Act3Message) Unmarshal(bytes []byte) error {
	pbAct3 := pb.Act3Message{}
	if err := pbAct3.Unmarshal(bytes); err != nil {
		return err
	}
	copy(am.challenge[:], pbAct3.Challenge[:32])

	return nil
}
