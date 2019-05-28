// Package watchtower continuously monitors the on-chain stake of all connected
// peers, and disconnects peers which fall below the minimum stake.
package watchtower

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
)

// Guard contains the state necessary to make connection pruning decisions.
type Guard struct {
	duration time.Duration

	stakeMonitorLock sync.Mutex
	stakeMonitor     chain.StakeMonitor

	connectionManager net.ConnectionManager

	peerCrossListLock sync.Mutex
	peerCrossList     map[string]bool
}

// NewGuard returns a new instance of Guard. Should only be called once per
// provider. Instantiating a new instance of Guard automatically runs it in the
// background for the lifetime of the client.
func NewGuard(
	ctx context.Context,
	duration time.Duration,
	stakeMonitor chain.StakeMonitor,
	connectionManager net.ConnectionManager,
) *Guard {
	guard := &Guard{
		duration:          duration,
		stakeMonitor:      stakeMonitor,
		connectionManager: connectionManager,
		peerCrossList:     make(map[string]bool),
	}
	go guard.start(ctx)
	return guard
}

func (g *Guard) currentlyChecking(peer string) bool {
	g.peerCrossListLock.Lock()
	_, inProcess := g.peerCrossList[peer]
	g.peerCrossListLock.Unlock()
	return inProcess
}

func (g *Guard) markAsChecking(peer string) {
	g.peerCrossListLock.Lock()
	g.peerCrossList[peer] = true
	g.peerCrossListLock.Unlock()
}

func (g *Guard) completedCheck(peer string) {
	g.peerCrossListLock.Lock()
	g.peerCrossList[peer] = false
	g.peerCrossListLock.Unlock()
}

// start executes the connection management background worker. If it receives a
// signal to stop the execution of the client, it kills this task.
func (g *Guard) start(ctx context.Context) {
	ticker := time.NewTicker(g.duration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, connectedPeer := range g.connectionManager.ConnectedPeers() {
				if g.currentlyChecking(connectedPeer) {
					continue
				}

				go g.manageConnectionByStake(ctx, connectedPeer)
			}
		}
	}
}

func (g *Guard) manageConnectionByStake(ctx context.Context, peer string) {
	g.markAsChecking(peer)
	defer g.completedCheck(peer)

	newContext, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	hasMinimumStake, err := g.validatePeerStake(
		newContext, peer,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !hasMinimumStake {
		g.connectionManager.DisconnectPeer(peer)
	}
}

func (g *Guard) validatePeerStake(ctx context.Context, peer string) (bool, error) {
	peerPublicKey, err := g.connectionManager.GetPeerPublicKey(peer)
	if err != nil {
		return false, err
	}

	if peerPublicKey == nil {
		return false, fmt.Errorf(
			"failed to resolve valid public key for peer %s", peer,
		)
	}

	g.stakeMonitorLock.Lock()
	hasMinimumStake, err := g.stakeMonitor.HasMinimumStake(
		key.NetworkPubKeyToEthAddress(peerPublicKey),
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
