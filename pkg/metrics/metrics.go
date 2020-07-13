package metrics

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/metrics"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-metrics")

const (
	// DefaultNetworkMetricsTick is the default duration of the
	// observation tick for network metrics.
	DefaultNetworkMetricsTick = 1 * time.Minute
	// DefaultEthereumMetricsTick is the default duration of the
	// observation tick for Ethereum metrics.
	DefaultEthereumMetricsTick = 10 * time.Minute
)

// Initialize set up the metrics registry and enables metrics server.
func Initialize(
	port int,
) (*metrics.Registry, bool) {
	if port == 0 {
		return nil, false
	}

	registry := metrics.NewRegistry()

	registry.EnableServer(port)

	return registry, true
}

// ObserveConnectedPeersCount triggers an observation process of the
// connected_peers_count metric.
func ObserveConnectedPeersCount(
	ctx context.Context,
	registry *metrics.Registry,
	netProvider net.Provider,
	tick time.Duration,
) {
	input := func() float64 {
		connectedPeers := netProvider.ConnectionManager().ConnectedPeers()

		peersList := make([]PeerInfo, len(connectedPeers))
		for i := 0; i < len(connectedPeers); i++ {
			peerAddress, err := netProvider.ConnectionManager().GetPeerPublicKey(connectedPeers[i])
			if err == nil {
				peersList[i] = PeerInfo{PeerId: connectedPeers[i], PeerAddress: peerAddress.X.String()}
			}
		}

		bytes, err := json.Marshal(peersList)
		if err == nil {
			logger.Debug("nodes list: ", string(bytes))

			label := metrics.NewLabel("list", string(bytes))
			registry.UpdateInfo("peerList", []metrics.Label{label})
		}

		return float64(len(connectedPeers))
	}

	observe(
		ctx,
		"connected_peers_count",
		input,
		registry,
		validateTick(tick, DefaultNetworkMetricsTick),
	)
}

type PeerInfo struct {
	PeerId      string
	PeerAddress string
}

// ObserveConnectedBootstrapCount triggers an observation process of the
// connected_bootstrap_count metric.
func ObserveConnectedBootstrapCount(
	ctx context.Context,
	registry *metrics.Registry,
	netProvider net.Provider,
	bootstraps []string,
	tick time.Duration,
) {
	input := func() float64 {
		currentCount := 0

		for _, address := range bootstraps {
			if netProvider.ConnectionManager().IsConnected(address) {
				currentCount++
			}
		}

		return float64(currentCount)
	}

	observe(
		ctx,
		"connected_bootstrap_count",
		input,
		registry,
		validateTick(tick, DefaultNetworkMetricsTick),
	)
}

// ObserveEthConnectivity triggers an observation process of the
// eth_connectivity metric.
func ObserveEthConnectivity(
	ctx context.Context,
	registry *metrics.Registry,
	stakeMonitor chain.StakeMonitor,
	address string,
	tick time.Duration,
) {
	input := func() float64 {
		_, err := stakeMonitor.HasMinimumStake(address)

		if err != nil {
			return 0
		}

		return 1
	}

	observe(
		ctx,
		"eth_connectivity",
		input,
		registry,
		validateTick(tick, DefaultEthereumMetricsTick),
	)
}

func observe(
	ctx context.Context,
	name string,
	input metrics.ObserverInput,
	registry *metrics.Registry,
	tick time.Duration,
) {
	observer, err := registry.NewGaugeObserver(name, input)
	if err != nil {
		logger.Warningf("could not create gauge observer [%v]", name)
		return
	}

	observer.Observe(ctx, tick)
}

func validateTick(tick time.Duration, defaultTick time.Duration) time.Duration {
	if tick > 0 {
		return tick
	}

	return defaultTick
}

// ExposeLibP2PInfo provides some basic information about libp2p config.
func ExposeLibP2PInfo(
	registry *metrics.Registry,
	netProvider net.Provider,
) {
	name := "libp2p_info"

	id := metrics.NewLabel("id", netProvider.ID().String())

	_, err := registry.NewInfo(name, []metrics.Label{id})
	if err != nil {
		logger.Warningf("could not create info metric [%v]", name)
		return
	}

	name = "peerList"
	id = metrics.NewLabel("list", "")

	_, err = registry.NewInfo(name, []metrics.Label{id})
	if err != nil {
		logger.Warningf("could not create info metric [%v]", name)
		return
	}
}
