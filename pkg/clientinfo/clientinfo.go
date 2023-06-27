package clientinfo

import (
	"context"
	"time"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-common/pkg/clientinfo"
)

var logger = log.Logger("keep-clientinfo")

// Config stores configuration for the client info.
type Config struct {
	Port                int
	NetworkMetricsTick  time.Duration
	EthereumMetricsTick time.Duration
	BitcoinMetricsTick  time.Duration
}

// Registry wraps keep-common clientinfo registry and exposes additional
// functions for registering client-custom metrics and diagnostics
type Registry struct {
	*clientinfo.Registry

	ctx context.Context
}

// Initialize set up the client info registry and enables metrics and
// diagnostics server.
func Initialize(
	ctx context.Context,
	port int,
) (*Registry, bool) {
	if port == 0 {
		return nil, false
	}

	registry := &Registry{clientinfo.NewRegistry(), ctx}

	registry.EnableServer(port)

	return registry, true
}
