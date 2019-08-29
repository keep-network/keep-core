// Package watchtower continuously monitors the on-chain stake of all connected
// peers, and disconnects peers which fall below the minimum stake.
package watchtower

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
)

var logger = log.Logger("keep-net-watchtower")

// Guard contains the state necessary to make connection pruning decisions.
type Guard struct {
	duration time.Duration

	stakeMonitor chain.StakeMonitor

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
	checking, _ := g.peerCrossList[peer]
	g.peerCrossListLock.Unlock()
	return checking
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
			connectedPeers := g.connectionManager.ConnectedPeers()
			logger.Debugf(
				"connected to [%v] peers: %v\n",
				len(connectedPeers),
				connectedPeers,
			)

			for _, connectedPeer := range connectedPeers {
				if g.currentlyChecking(connectedPeer) {
					continue
				}

				// Ensure we mark the peer as being checked before
				// executing the async stake check.
				g.markAsChecking(connectedPeer)
				go g.manageConnectionByStake(ctx, connectedPeer)
			}
		}
	}
}

func (g *Guard) manageConnectionByStake(ctx context.Context, peer string) {
	defer g.completedCheck(peer)

	newContext, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	peerPublicKey, err := g.getPeerPublicKey(peer)
	if err != nil {
		// if we error while getting the peer's public key, the peer's id
		// or key may be malformed/unknown; disconnect them immediately.
		logger.Errorf(
			"dropping the connection - could not get public key for peer [%v]: [%v]",
			peer,
			err,
		)
		g.connectionManager.DisconnectPeer(peer)
		return
	}

	hasMinimumStake, err := g.validatePeerStake(
		newContext, peerPublicKey,
	)
	if err != nil {
		// network issues with geth shouldn't cause disconnects from the
		// network. Rather we'll abort the check and try again later.
		logger.Warningf("error validating peer stake, retrying later: [%v].", err)
		return
	}

	if !hasMinimumStake {
		// if a peer doesn't have at least the min stake, disconnect them.
		logger.Warningf(
			"dropping the connection - peer [%v] has no minimal stake",
			peer,
		)
		g.connectionManager.DisconnectPeer(peer)
	}
}

func (g *Guard) getPeerPublicKey(peer string) (*key.NetworkPublic, error) {
	peerPublicKey, err := g.connectionManager.GetPeerPublicKey(peer)
	if err != nil {
		return nil, err
	}

	if peerPublicKey == nil {
		return nil, fmt.Errorf(
			"failed to resolve valid public key for peer %s", peer,
		)
	}
	return peerPublicKey, nil
}

func (g *Guard) validatePeerStake(
	ctx context.Context,
	peerPublicKey *key.NetworkPublic,
) (bool, error) {
	hasMinimumStake, err := g.stakeMonitor.HasMinimumStake(
		key.NetworkPubKeyToEthAddress(peerPublicKey),
	)
	if err != nil {
		return false, fmt.Errorf(
			"Failed to get stake information for key [%s] with error: [%v]",
			peerPublicKey,
			err,
		)
	}

	return hasMinimumStake, nil
}
