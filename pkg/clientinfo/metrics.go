package clientinfo

import (
	"fmt"
	"time"

	"github.com/keep-network/keep-common/pkg/clientinfo"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

type Source func() float64

// Names under which metrics are exposed.
const (
	ConnectedPeersCountMetricName     = "connected_peers_count"
	ConnectedBootstrapCountMetricName = "connected_bootstrap_count"
	EthConnectivityMetricName         = "eth_connectivity"
	ClientInfoMetricName              = "client_info"
)

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
		ConnectedPeersCountMetricName,
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
		ConnectedBootstrapCountMetricName,
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
		EthConnectivityMetricName,
		input,
		validateTick(tick, DefaultEthereumMetricsTick),
	)
}

// ObserveApplicationSource triggers an observation process of
// application-specific metrics.
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

// RegisterMetricClientInfo registers static client information labels for metrics.
func (r *Registry) RegisterMetricClientInfo(version string) {
	_, err := r.NewMetricInfo(
		ClientInfoMetricName,
		[]clientinfo.Label{
			clientinfo.NewLabel("version", version),
		},
	)
	if err != nil {
		logger.Warnf("could not register metric client info: [%v]", err)
	}
}

func (r *Registry) observe(
	name string,
	input Source,
	tick time.Duration,
) {
	observer, err := r.NewMetricGaugeObserver(name, clientinfo.MetricObserverInput(input))
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
