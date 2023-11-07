package clientinfo

import (
	"encoding/json"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"

	"github.com/keep-network/keep-core/pkg/net"
)

// Diagnostics describes data structure returned by the diagnostics endpoint.
type Diagnostics struct {
	ClientInfo     Client `json:"client_info"`
	ConnectedPeers []Peer `json:"connected_peers"`
	EthChainInfo   Chain  `json:"eth_chain_info"`
	BtcChainInfo   Chain  `json:"btc_chain_info"`
}

// Client describes data structure of client information.
type Client struct {
	ChainAddress string `json:"chain_address"`
	NetworkID    string `json:"network_id"`
	Version      string `json:"version"`
	Revision     string `json:"revision"`
}

// Peer describes data structure of peer information.
type Peer struct {
	ChainAddress          string   `json:"chain_address"`
	NetworkID             string   `json:"network_id"`
	NetworkMultiAddresses []string `json:"multiaddrs"`
}

// Chain describes data structure of a chain information.
type Chain struct {
	LatestBlockNumber uint `json:"latest_block_number"`
}

// ApplicationInfo describes data structure of application information.
type ApplicationInfo map[string]interface{}

// RegisterConnectedPeersSource registers the diagnostics source providing
// information about connected peers.
func (r *Registry) RegisterConnectedPeersSource(
	netProvider net.Provider,
	signing chain.Signing,
) {
	r.RegisterDiagnosticSource("connected_peers", func() string {
		connectionManager := netProvider.ConnectionManager()
		connectedPeersAddrInfo := connectionManager.ConnectedPeersAddrInfo()

		var peersList []Peer
		for peerNetworkID, multiaddrs := range connectedPeersAddrInfo {
			peerPublicKey, err := connectionManager.GetPeerPublicKey(peerNetworkID)
			if err != nil {
				logger.Errorf("error on getting peer public key: [%v]", err)
				continue
			}

			peerChainAddress, err := signing.PublicKeyToAddress(
				peerPublicKey,
			)
			if err != nil {
				logger.Errorf("error on getting peer chain address: [%v]", err)
				continue
			}

			peersList = append(peersList, Peer{
				NetworkID:             peerNetworkID,
				ChainAddress:          peerChainAddress.String(),
				NetworkMultiAddresses: multiaddrs,
			})
		}

		bytes, err := json.Marshal(peersList)
		if err != nil {
			logger.Errorf("error on serializing peers list to JSON: [%v]", err)
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
	r.RegisterDiagnosticSource("client_info", func() string {
		connectionManager := netProvider.ConnectionManager()

		clientID := netProvider.ID().String()
		clientPublicKey, err := connectionManager.GetPeerPublicKey(clientID)
		if err != nil {
			logger.Errorf("error on getting client public key: [%v]", err)
			return ""
		}

		clientChainAddress, err := signing.PublicKeyToAddress(
			clientPublicKey,
		)
		if err != nil {
			logger.Errorf("error on getting peer chain address: [%v]", err)
			return ""
		}

		clientInfo := Client{
			NetworkID:    clientID,
			ChainAddress: clientChainAddress.String(),
			Version:      clientVersion,
			Revision:     clientRevision,
		}

		bytes, err := json.Marshal(clientInfo)
		if err != nil {
			logger.Errorf("error on serializing client info to JSON: [%v]", err)
			return ""
		}

		return string(bytes)
	})
}

// RegisterBtcChainInfoSource registers the diagnostics source providing
// information about btc chain.
func (r *Registry) RegisterBtcChainInfoSource(
	btcChain bitcoin.Chain,
) {
	r.RegisterDiagnosticSource("btc_chain_info", func() string {
		btcLatestBlock, err := btcChain.GetLatestBlockHeight()
		if err != nil {
			logger.Errorf("error on getting Bitcoin latest block number: [%v]", err)
		}
		chainInfo := Chain{
			LatestBlockNumber: btcLatestBlock,
		}

		bytes, err := json.Marshal(chainInfo)
		if err != nil {
			logger.Errorf("error on serializing btc chain info to JSON: [%v]", err)
			return ""
		}

		return string(bytes)
	})
}

// RegisterEthChainInfoSource registers the diagnostics source providing
// information about eth chain.
func (r *Registry) RegisterEthChainInfoSource(
	blockCounter chain.BlockCounter,
) {
	r.RegisterDiagnosticSource("eth_chain_info", func() string {
		ethLatestBlock, err := blockCounter.CurrentBlock()
		if err != nil {
			logger.Errorf("error on getting Ethereum latest block number: [%v]", err)
		}
		chainInfo := Chain{
			LatestBlockNumber: uint(ethLatestBlock),
		}

		bytes, err := json.Marshal(chainInfo)
		if err != nil {
			logger.Errorf("error on serializing eth chain info to JSON: [%v]", err)
			return ""
		}

		return string(bytes)
	})
}

// RegisterApplicationSource registers the diagnostics source providing
// information about the application.
func (r *Registry) RegisterApplicationSource(
	application string,
	fetchApplicationDiagnostics func() ApplicationInfo,
) {
	r.RegisterDiagnosticSource(application, func() string {
		bytes, err := json.Marshal(fetchApplicationDiagnostics())
		if err != nil {
			logger.Errorf("error on serializing peers list to JSON: [%v]", err)
			return ""
		}

		return string(bytes)
	})
}
