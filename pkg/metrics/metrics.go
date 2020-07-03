package metrics

import (
	"context"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/metrics"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-metrics")

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
		return float64(len(connectedPeers))
	}

	observe(
		ctx,
		"connected_peers_count",
		input,
		registry,
		tick,
	)
}

// ObserveConnectedBootstrapPercentage triggers an observation process of the
// connected_bootstrap_percentage metric.
func ObserveConnectedBootstrapPercentage(
	ctx context.Context,
	registry *metrics.Registry,
	netProvider net.Provider,
	bootstraps []string,
	tick time.Duration,
) {
	input := func() float64 {
		maxCount := len(bootstraps)
		if maxCount == 0 {
			return 0
		}

		currentCount := 0
		for _, address := range bootstraps {
			if netProvider.ConnectionManager().IsConnected(address) {
				currentCount++
			}
		}

		return float64((currentCount * 100) / maxCount)
	}

	observe(
		ctx,
		"connected_bootstrap_percentage",
		input,
		registry,
		tick,
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
		tick,
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
}
