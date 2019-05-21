// Package watchtower is a packge with introduces Guard, a type that takes a
// stakemonitor and libp2p host. The purpose of this package is to
// continuously monitor the on-chain stake of a connected peer, and to
// disconnect peers which fall below the minimum stake.
package watchtower

import (
	"context"
	"fmt"
	"sync"
	"time"

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

	peerCrossListLock sync.Mutex
	peerCrossList     map[peer.ID]bool
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
		stakeMonitor:  stakeMonitor,
		host:          host,
		peerCrossList: make(map[peer.ID]bool),
	}
	go guard.start(ctx)
	return guard
}

func (g *Guard) currentlyChecking(peer peer.ID) bool {
	g.peerCrossListLock.Lock()
	_, inProcess := g.peerCrossList[peer]
	g.peerCrossListLock.Unlock()
	return inProcess
}

func (g *Guard) markAsChecking(peer peer.ID) {
	g.peerCrossListLock.Lock()
	g.peerCrossList[peer] = true
	g.peerCrossListLock.Unlock()
}

func (g *Guard) completedCheck(peer peer.ID) {
	g.peerCrossListLock.Lock()
	g.peerCrossList[peer] = false
	g.peerCrossListLock.Unlock()
}

// start executes the connection management background worker. If it receives a
// signal to stop the execution of the client, it kills this task.
func (g *Guard) start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, connectedPeer := range g.host.Network().Peers() {
				if g.currentlyChecking(connectedPeer) {
					continue
				}

				go func(ctx context.Context, inProcessPeer peer.ID) {
					g.markAsChecking(inProcessPeer)
					defer g.completedCheck(inProcessPeer)

					newContext, cancel := context.WithTimeout(ctx, 10*time.Second)
					defer cancel()

					hasMinimumStake, err := g.validatePeerStake(
						newContext, inProcessPeer,
					)
					if err != nil {
						fmt.Println(err)
						return
					}

					if !hasMinimumStake {
						g.disconnectPeer(inProcessPeer)
					}
				}(ctx, connectedPeer)
			}
		}
	}
}

func (g *Guard) validatePeerStake(ctx context.Context, peer peer.ID) (bool, error) {
	peerPublicKey, err := peer.ExtractPublicKey()
	if err != nil {
		return false, fmt.Errorf(
			"Failed to extract peer [%s] public key with error [%v]",
			peer,
			err,
		)
	}

	g.stakeMonitorLock.Lock()
	hasMinimumStake, err := g.stakeMonitor.HasMinimumStake(
		key.NetworkPubKeyToEthAddress(
			peerPublicKey.(*key.NetworkPublic),
		),
	)
	if err != nil {
		g.stakeMonitorLock.Unlock()
		return false, fmt.Errorf(
			"Failed to get stake information for peer [%s] with error [%v]",
			peer,
			err,
		)
	}
	g.stakeMonitorLock.Unlock()
	return hasMinimumStake, nil
}

func (g *Guard) disconnectPeer(peer peer.ID) {
	connections := g.host.Network().ConnsToPeer(peer)
	for _, connection := range connections {
		connection.Close()
	}
}
