package diagnostics

import (
	"encoding/json"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/diagnostics"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
)

var logger = log.Logger("keep-diagnostics")

// Initialize set up the diagnostics registry and enables diagnostics server.
func Initialize(
	port int,
) (*diagnostics.DiagnosticsRegistry, bool) {
	if port == 0 {
		return nil, false
	}

	registry := diagnostics.NewRegistry()

	registry.EnableServer(port)

	return registry, true
}

// Registers diagnostics source providing information about connected peers
func RegisterConnectedPeersSource(registry *diagnostics.DiagnosticsRegistry, netProvider net.Provider) {
	registry.RegisterSource("connected_peers", func() string {
		connectionManager := netProvider.ConnectionManager()
		connectedPeers := connectionManager.ConnectedPeers()

		peersList := make([]map[string]interface{}, len(connectedPeers))
		for i := 0; i < len(connectedPeers); i++ {
			peer := connectedPeers[i]
			peerPublicKey, err := connectionManager.GetPeerPublicKey(peer)
			if err != nil {
				logger.Error("error on getting peer public key: [%v]", err)
				continue
			}

			peersList[i] = map[string]interface{}{
				"PeerId":        peer,
				"PeerPublicKey": key.NetworkPubKeyToEthAddress(peerPublicKey),
			}
		}

		bytes, err := json.Marshal(peersList)
		if err != nil {
			logger.Error("error on serializing peers list to JSON: [%v]", err)
			return ""
		}

		return string(bytes)
	})
}
