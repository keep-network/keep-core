package diagnostics

import (
	"encoding/json"
	"strings"

	"github.com/keep-network/keep-core/pkg/chain"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/diagnostics"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-diagnostics")

// Config stores diagnostics-related configuration.
type Config struct {
	Port int
}

// Initialize sets up the diagnostics registry and enables diagnostics server.
func Initialize(port int) (*diagnostics.Registry, bool) {
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
	registry *diagnostics.Registry,
	netProvider net.Provider,
	signing chain.Signing,
) {
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

			peerChainAddress, err := signing.PublicKeyToAddress(
				peerPublicKey,
			)
			if err != nil {
				logger.Error("error on getting peer chain address: [%v]", err)
				continue
			}

			peersList[i] = map[string]interface{}{
				"network_id":    peer,
				"chain_address": peerChainAddress.String(),
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
	registry *diagnostics.Registry,
	netProvider net.Provider,
	signing chain.Signing,
	clientVersion string,
) {
	registry.RegisterSource("client_info", func() string {
		connectionManager := netProvider.ConnectionManager()

		clientID := netProvider.ID().String()
		clientPublicKey, err := connectionManager.GetPeerPublicKey(clientID)
		if err != nil {
			logger.Error("error on getting client public key: [%v]", err)
			return ""
		}

		clientChainAddress, err := signing.PublicKeyToAddress(
			clientPublicKey,
		)
		if err != nil {
			logger.Error("error on getting peer chain address: [%v]", err)
			return ""
		}

		clientVersionArr := strings.Split(clientVersion, "revision")
		clientInfo := map[string]interface{}{
			"network_id":    clientID,
			"chain_address": clientChainAddress.String(),
			"version": strings.TrimSpace(clientVersionArr[0]),
			"revision": strings.TrimSpace(clientVersionArr[1]),
		}

		bytes, err := json.Marshal(clientInfo)
		if err != nil {
			logger.Error("error on serializing client info to JSON: [%v]", err)
			return ""
		}

		return string(bytes)
	})
}
