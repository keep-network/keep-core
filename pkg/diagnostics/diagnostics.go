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

// All available protocols https://github.com/multiformats/multiaddr#protocols
// Eligible protocols are used to fetch the address to call the /diagnostics
// endpoint
var eligibleIP4 = "ip4"
var eligibleDNS4 = "dns4"

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
		connectedPeersAddrInfo := connectionManager.ConnectedPeersAddrInfo()

		peersList := make([]map[string]interface{}, len(connectedPeers))
		for i := 0; i < len(connectedPeers); i++ {
			peer := connectedPeers[i]
			peersAddrInfo := connectedPeersAddrInfo[i].Addrs

			var peerAddr string
			for _, peerAddrInfo := range peersAddrInfo {
				// Peer may contain public and local addresses. We need to fetch a
				// public address only. Range of local addresses starting from 127.* are
				// skipped.
				if strings.Contains(peerAddrInfo.String(), eligibleDNS4) ||
					(strings.Contains(peerAddrInfo.String(), eligibleIP4) &&
						!strings.Contains(peerAddrInfo.String(), "/127.")) {
					// address is formatted as follows /<protocol_code>/<ip_address>/<connection_protocol>/<port>
					// Address is  the 1st index
					peerAddr = strings.Split(strings.Trim(peerAddrInfo.String(), "/"), "/")[1]
				}
			}

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
				"address":       peerAddr,
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

		clientInfo := map[string]interface{}{
			"network_id":    clientID,
			"chain_address": clientChainAddress.String(),
			"version":       clientVersion,
		}

		bytes, err := json.Marshal(clientInfo)
		if err != nil {
			logger.Error("error on serializing client info to JSON: [%v]", err)
			return ""
		}

		return string(bytes)
	})
}

// RegisterApplicationSource registers the diagnostics source providing
// information about the application.
func RegisterApplicationSource(registry *diagnostics.Registry, application string, fetchApplicationDiagnostics func() map[string]interface{}) {
	registry.RegisterSource(application, func() string {
		bytes, err := json.Marshal(fetchApplicationDiagnostics())
		if err != nil {
			logger.Error("error on serializing peers list to JSON: [%v]", err)
			return ""
		}

		return string(bytes)
	})
}
