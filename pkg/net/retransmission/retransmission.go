package retransmission

import (
	"github.com/keep-network/keep-core/pkg/net"
)

// NetworkMessage enhances net.Message with functions needed to perform message
// retransmission. Specifically, we need to know if the given message is
// a retransmission and know a fingerprint of the message to filter out
// duplicates in the message handler.
type NetworkMessage interface {
	net.Message

	Fingerprint() string
	Retransmission() uint32
}
