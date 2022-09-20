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

// Registry wraps keep-common registry for internal use of exposed keep-common
// registry methods.
type Registry struct {
	Registry *diagnostics.Registry
}

// Initialize sets up the diagnostics registry and enables diagnostics server.
func Initialize(port int) (*Registry, bool) {
	if port == 0 {
		return nil, false
	}

	registry := diagnostics.NewRegistry()

	registry.EnableServer(port)

	newRegistry := &Registry{
		Registry: registry,
	}

	return newRegistry, true
}

// RegisterConnectedPeersSource registers the diagnostics source providing
// information about connected peers.
func (r *Registry) RegisterConnectedPeersSource(
	netProvider net.Provider,
	signing chain.Signing,
) {
	r.Registry.RegisterSource("connected_peers", func() string {
		connectionManager := netProvider.ConnectionManager()
		connectedPeers := connectionManager.ConnectedPeers()
		connectedPeersAddrInfo := connectionManager.ConnectedPeersAddrInfo()

		peersList := make([]map[string]interface{}, len(connectedPeers))
		for i := 0; i < len(connectedPeers); i++ {
			peer := connectedPeers[i]
			peersAddrInfo := connectedPeersAddrInfo[i].Addrs

			var peerAddr string
			for _, peerAddrInfo := range peersAddrInfo {
				if strings.Contains(peerAddrInfo.String(), eligibleDNS4) ||
					(strings.Contains(peerAddrInfo.String(), eligibleIP4)) {
					// address is formatted as follows /<protocol_code>/<ip_address>/<connection_protocol>/<port>
					// Address is the 1st index
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
func (r *Registry) RegisterClientInfoSource(
	netProvider net.Provider,
	signing chain.Signing,
	clientVersion string,
	clientRevision string,
) {
	r.Registry.RegisterSource("client_info", func() string {
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
			"revision":      clientRevision,
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
func (r *Registry) RegisterApplicationSource(application string, fetchApplicationDiagnostics func() map[string]interface{}) {
	r.Registry.RegisterSource(application, func() string {
		bytes, err := json.Marshal(fetchApplicationDiagnostics())
		if err != nil {
			logger.Error("error on serializing peers list to JSON: [%v]", err)
			return ""
		}

		return string(bytes)
	})
}
