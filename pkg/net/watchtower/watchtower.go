// Package watchtower is a packge with introduces Guard, a type that takes a
// stakemonitor and libp2p host. The purpose of this package is to
// continuously monitor the on-chain stake of a connected peer, and to
// disconnect peers which fall below the minimum stake.
package watchtower

import (
	"context"
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net/key"
	host "github.com/libp2p/go-libp2p-host"
	peer "github.com/libp2p/go-libp2p-peer"
)

// Guard contains the state necessary to make connection pruning decisions.
type Guard struct {
	stakeMonitorLock sync.Mutex
	stakeMonitor     chain.StakeMonitor

	host host.Host
}

// NewGuard returns a new instance of Guard. Can only be called once.
// Instantiating a new instance of Guard automatically runs the Guard instance in
// the background for the lifetime of the client.
func NewGuard(
	ctx context.Context,
	stakeMonitor chain.StakeMonitor,
	host host.Host,
) *Guard {
	guard := &Guard{
		stakeMonitor: stakeMonitor,
		host:         host,
	}
	go guard.start(ctx)
	return guard
}

// start executes the connection management background worker. If it receives a
// signal to stop the execution of the client, it kills this task.
func (g *Guard) start(ctx context.Context) {
	// use a timer or you're gonna blow out the cpu
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, connectedPeer := range g.host.Network().Peers() {
				go func(ctx context.Context, connectedPeer peer.ID) {
					_, cancel := context.WithCancel(ctx)
					defer cancel()

					peerPublicKey, err := connectedPeer.ExtractPublicKey()
					if err != nil {
						fmt.Printf(
							"Failed to extract peer [%s] public key with error [%v]",
							connectedPeer,
							err,
						)
						return
					}

					g.stakeMonitorLock.Lock()
					hasMinimumStake, err := g.stakeMonitor.HasMinimumStake(
						key.NetworkPubKeyToEthAddress(
							peerPublicKey.(*key.NetworkPublic),
						),
					)
					if err != nil {
						g.stakeMonitorLock.Unlock()
						fmt.Printf(
							"Failed to get stake information for peer [%s] with error [%v]",
							connectedPeer,
							err,
						)
						return
					}
					g.stakeMonitorLock.Unlock()

					if !hasMinimumStake {
						connections := g.host.Network().ConnsToPeer(connectedPeer)
						for _, connection := range connections {
							connection.Close()
						}
					}
				}(ctx, connectedPeer)
			}
		}
	}
}
