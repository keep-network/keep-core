package handshake

import (
	"encoding/binary"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
)

func (am *act1Message) Marshal() ([]byte, error) {
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, am.nonce1)
	return (&pb.Act1Message{Nonce: nonceBytes}).Marshal()
}

func (am *act1Message) Unmarshal(bytes []byte) error {
	pbAct1 := pb.Act1Message{}
	if err := pbAct1.Unmarshal(bytes); err != nil {
		return err
	}

	am.nonce1 = binary.LittleEndian.Uint64(pbAct1.Nonce)
	return nil
}
