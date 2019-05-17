package watchtower

import (
	"context"
	"sync"

	"github.com/keep-network/keep-core/pkg/chain"
	host "github.com/libp2p/go-libp2p-host"
)

// The watchtower takes a stakemonitor,
// runs the monitor in a loop, checking all potential connections

type Guard struct {
	stakeMonitorLock sync.Mutex
	stakeMonitor     chain.StakeMonitor

	// networkLock sync.Mutex
	// network     swarm.Dialer
	host host.Host
}

func NewGuard(
	ctx context.Context,
	stakeMonitor chain.StakeMonitor,
	host host.Host,
) *Guard {
	guard := &Guard{
		stakeMonitor: stakeMonitor,
		host:         host,
	}
	go guard.Start(ctx)
	return guard
}

func (g *Guard) Start(ctx context.Context) {
	// use a timer or you're gonna blow out the cpu
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// get the connected peers
			// g.host.Peerstore().Peers()
			// do the stake check
			//
		}
	}
}
