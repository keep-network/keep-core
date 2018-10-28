package handshake

import (
	"encoding/binary"

	"github.com/gogo/protobuf/proto"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
)

// Proto converts this Act1Message to a proto.Message suitable for network
// communication.
func (am *Act1Message) Proto() proto.Message {
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, am.nonce1)
	return &pb.Act1Message{Nonce: nonceBytes}
}

// Act1MessageFromProto converts a pb.Act1Message produced by Proto to a
// Act1Message.
func Act1MessageFromProto(pbAct1 pb.Act1Message) *Act1Message {
	return &Act1Message{nonce1: binary.LittleEndian.Uint64(pbAct1.Nonce)}
}

// Proto converts this Act2Message to a proto.Message suitable for network
// communication.
func (am *Act2Message) Proto() proto.Message {
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, am.nonce2)
	return &pb.Act2Message{Nonce: nonceBytes, Challenge: am.challenge[:]}
}

// Act2MessageFromProto converts a pb.Act2Message produced by Proto to a
// Act2Message.
func Act2MessageFromProto(pbAct2 pb.Act2Message) *Act2Message {
	am := &Act2Message{nonce2: binary.LittleEndian.Uint64(pbAct2.Nonce)}
	copy(am.challenge[:], pbAct2.Challenge[:])
	return am
}

// Proto converts this Act3Message to a proto.Message suitable for network
// communication.
func (am *Act3Message) Proto() proto.Message {
	return &pb.Act3Message{Challenge: am.challenge[:]}
}

// Act3MessageFromProto converts a pb.Act3Message produced by Proto to a
// Act3Message.
func Act3MessageFromProto(pbAct3 pb.Act3Message) *Act3Message {
	am := &Act3Message{}
	copy(am.challenge[:], pbAct3.Challenge[:])
	return am
}
