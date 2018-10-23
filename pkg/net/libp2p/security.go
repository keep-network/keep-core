package libp2p

import (
	"github.com/keep-network/keep-core/pkg/net/security/handshake"
	libp2p "github.com/libp2p/go-libp2p"
)

// DefaultSecurity is the default security option.
//
// Precedes every connection
var DefaultSecurity = libp2p.Security(handshake.ID, handshake.New)
