package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/metrics"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-metrics")

type Source func() float64

const (
	// DefaultNetworkMetricsTick is the default duration of the
	// observation tick for network metrics.
	DefaultNetworkMetricsTick = 1 * time.Minute
	// DefaultEthereumMetricsTick is the default duration of the
	// observation tick for Ethereum metrics.
	DefaultEthereumMetricsTick = 10 * time.Minute
	// The duration of the observation tick for all application-specific
	// metrics.
	ApplicationMetricsTick = 1 * time.Minute
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

	ctx context.Context
}

// Initialize set up the metrics registry and enables metrics server.
func Initialize(
	ctx context.Context,
	port int,
) (*Registry, bool) {
	if port == 0 {
		return nil, false
	}

	registry := &Registry{metrics.NewRegistry(), ctx}

	registry.EnableServer(port)

	return registry, true
}

// ObserveConnectedPeersCount triggers an observation process of the
// connected_peers_count metric.
func (r *Registry) ObserveConnectedPeersCount(
	netProvider net.Provider,
	tick time.Duration,
) {
	input := func() float64 {
		connectedPeers := netProvider.ConnectionManager().ConnectedPeers()
		return float64(len(connectedPeers))
	}

	r.observe(
		"connected_peers_count",
		input,
		validateTick(tick, DefaultNetworkMetricsTick),
	)
}

// ObserveConnectedBootstrapCount triggers an observation process of the
// connected_bootstrap_count metric.
func (r *Registry) ObserveConnectedBootstrapCount(
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
		"connected_bootstrap_count",
		input,
		validateTick(tick, DefaultNetworkMetricsTick),
	)
}

// ObserveEthConnectivity triggers an observation process of the
// eth_connectivity metric.
func (r *Registry) ObserveEthConnectivity(
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
		"eth_connectivity",
		input,
		validateTick(tick, DefaultEthereumMetricsTick),
	)
}

// ObserveEthConnectivity triggers an observation process of application-specific
// metrics.
func (r *Registry) ObserveApplicationSource(
	application string,
	inputs map[string]Source,
) {
	for k, v := range inputs {
		r.observe(
			fmt.Sprintf("%s_%s", application, k),
			v,
			ApplicationMetricsTick,
		)
	}
}

func (r *Registry) observe(
	name string,
	input Source,
	tick time.Duration,
) {
	observer, err := r.NewGaugeObserver(name, metrics.ObserverInput(input))
	if err != nil {
		logger.Warnf("could not create gauge observer [%v]", name)
		return
	}

	observer.Observe(r.ctx, tick)

	logger.Infof("observing %s with [%s] tick", name, tick)
}

func validateTick(tick time.Duration, defaultTick time.Duration) time.Duration {
	if tick > 0 {
		return tick
	}

	return defaultTick
}
