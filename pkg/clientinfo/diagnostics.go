package clientinfo

import (
	"encoding/json"

	"github.com/keep-network/keep-core/pkg/chain"

	"github.com/keep-network/keep-core/pkg/net"
)

// RegisterConnectedPeersSource registers the diagnostics source providing
// information about connected peers.
func (r *Registry) RegisterConnectedPeersSource(
	netProvider net.Provider,
	signing chain.Signing,
) {
	r.Registry.RegisterDiagnosticSource("connected_peers", func() string {
		connectionManager := netProvider.ConnectionManager()
		connectedPeersAddrInfo := connectionManager.ConnectedPeersAddrInfo()

		var peersList []map[string]interface{}
		for peer, multiaddrs := range connectedPeersAddrInfo {
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

			peerInfo := map[string]interface{}{
				"network_id":    peer,
				"chain_address": peerChainAddress.String(),
				"multiaddrs":    multiaddrs,
			}
			peersList = append(peersList, peerInfo)
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
	r.Registry.RegisterDiagnosticSource("client_info", func() string {
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
func (r *Registry) RegisterApplicationSource(
	application string,
	fetchApplicationDiagnostics func() map[string]interface{},
) {
	r.Registry.RegisterDiagnosticSource(application, func() string {
		bytes, err := json.Marshal(fetchApplicationDiagnostics())
		if err != nil {
			logger.Error("error on serializing peers list to JSON: [%v]", err)
			return ""
		}

		return string(bytes)
	})
}