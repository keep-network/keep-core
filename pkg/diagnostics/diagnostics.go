package diagnostics

import (
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/diagnostics"
)

var logger = log.Logger("keep-diagnostics")

// Initialize set up the diagnostics registry and enables diagnostics server.
func Initialize(
	port int,
) (*diagnostics.DiagnosticsRegistry, bool) {
	if port == 0 {
		return nil, false
	}

	registry := diagnostics.NewRegistry()

	registry.EnableServer(port)

	return registry, true
}
