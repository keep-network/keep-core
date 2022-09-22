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

const (
	// DefaultNetworkMetricsTick is the default duration of the
	// observation tick for network metrics.
	DefaultNetworkMetricsTick = 1 * time.Minute
	// DefaultEthereumMetricsTick is the default duration of the
	// observation tick for Ethereum metrics.
	DefaultEthereumMetricsTick = 10 * time.Minute
)

// Config stores meta-info about metrics.
type Config struct {
	Port                int
	NetworkMetricsTick  time.Duration
	EthereumMetricsTick time.Duration
}

// Registry wraps keep-common metrics registry and exposes additional functions
// for registering client-custom metrics.
type Registry struct {
	*metrics.Registry
}

// Initialize set up the metrics registry and enables metrics server.
func Initialize(
	port int,
) (*Registry, bool) {
	if port == 0 {
		return nil, false
	}

	registry := &Registry{metrics.NewRegistry()}

	registry.EnableServer(port)

	return registry, true
}

// ObserveConnectedPeersCount triggers an observation process of the
// connected_peers_count metric.
func (r *Registry) ObserveConnectedPeersCount(
	ctx context.Context,
	netProvider net.Provider,
	tick time.Duration,
) {
	input := func() float64 {
		connectedPeers := netProvider.ConnectionManager().ConnectedPeers()
		return float64(len(connectedPeers))
	}

	r.observe(
		ctx,
		"connected_peers_count",
		input,
		validateTick(tick, DefaultNetworkMetricsTick),
	)
}

// ObserveConnectedBootstrapCount triggers an observation process of the
// connected_bootstrap_count metric.
func (r *Registry) ObserveConnectedBootstrapCount(
	ctx context.Context,
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

	r.observe(
		ctx,
		"connected_bootstrap_count",
		input,
		validateTick(tick, DefaultNetworkMetricsTick),
	)
}

// ObserveEthConnectivity triggers an observation process of the
// eth_connectivity metric.
func (r *Registry) ObserveEthConnectivity(
	ctx context.Context,
	blockCounter chain.BlockCounter,
	tick time.Duration,
) {
	input := func() float64 {
		_, err := blockCounter.CurrentBlock()

		if err != nil {
			return 0
		}

		return 1
	}

	r.observe(
		ctx,
		"eth_connectivity",
		input,
		validateTick(tick, DefaultEthereumMetricsTick),
	)
}

func (r *Registry) observe(
	ctx context.Context,
	name string,
	input metrics.ObserverInput,
	tick time.Duration,
) {
	observer, err := r.NewGaugeObserver(name, input)
	if err != nil {
		logger.Warnf("could not create gauge observer [%v]", name)
		return
	}

	observer.Observe(ctx, tick)

	logger.Infof("observing %s with [%s] tick", name, tick)
}

func validateTick(tick time.Duration, defaultTick time.Duration) time.Duration {
	if tick > 0 {
		return tick
	}

	return defaultTick
}
