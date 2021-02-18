package diagnostics

import (
	"encoding/json"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/diagnostics"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
	manet "github.com/multiformats/go-multiaddr/net"
)

var logger = log.Logger("keep-diagnostics")

// Initialize sets up the diagnostics registry and enables diagnostics server.
func Initialize(port int) (*diagnostics.DiagnosticsRegistry, bool) {
	if port == 0 {
		return nil, false
	}

	registry := diagnostics.NewRegistry()

	registry.EnableServer(port)

	return registry, true
}

// RegisterConnectedPeersSource registers the diagnostics source providing
// information about connected peers.
func RegisterConnectedPeersSource(
	registry *diagnostics.DiagnosticsRegistry,
	netProvider net.Provider,
) {
	registry.RegisterSource("connected_peers", func() string {
		connectionManager := netProvider.ConnectionManager()
		peers := connectionManager.ConnectedPeersAddrInfo()
		logger.Infof("get peers number: %d", len(peers))

		peersList := make([]map[string]interface{}, len(peers))
		for i := 0; i < len(peers); i++ {
			peer := peers[i]
			peerPublicKey, err := connectionManager.GetPeerPublicKey(peer.ID)
			if err != nil {
				logger.Errorf("error on getting peer public key: [%v]", err)
				continue
			}
			Addr, err := manet.ToNetAddr(peer.Addr)
			if err != nil {
				logger.Errorf("error on getting peer net addr: [%v]", err)
				continue
			}

			peersList[i] = map[string]interface{}{
				"network_id":       peer.ID,
				"network_addr":     Addr.String(),
				"ethereum_address": key.NetworkPubKeyToEthAddress(peerPublicKey),
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

// RegisterClientInfoSource registers the diagnostics source providing
// information about the client itself.
func RegisterClientInfoSource(
	registry *diagnostics.DiagnosticsRegistry,
	netProvider net.Provider,
) {
	registry.RegisterSource("client_info", func() string {
		logger.Warning("get client_info")

		connectionManager := netProvider.ConnectionManager()

		clientID := netProvider.ID().String()
		clientPublicKey, err := connectionManager.GetPeerPublicKey(clientID)
		if err != nil {
			logger.Error("error on getting client public key: [%v]", err)
			return ""
		}

		clientInfo := map[string]interface{}{
			"network_id":       clientID,
			"ethereum_address": key.NetworkPubKeyToEthAddress(clientPublicKey),
			"network_addrs":    connectionManager.NetAddrStrings(),
			"datetime":         time.Now().Format("2006-01-02 15:04:05"),
		}

		bytes, err := json.Marshal(clientInfo)
		if err != nil {
			logger.Error("error on serializing client info to JSON: [%v]", err)
			return ""
		}

		return string(bytes)
	})
}
