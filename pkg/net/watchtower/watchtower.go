// Package watchtower continuously monitors firewall rules compliance of all
// connected peers, and disconnects peers which do not comply to the rules.
package watchtower

import (
	"context"
	"fmt"
	"github.com/keep-network/keep-core/pkg/operator"
	"sync"
	"time"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-core/pkg/net"
)

// Guard contains the state necessary to make connection pruning decisions.
type Guard struct {
	logger log.StandardLogger

	duration time.Duration

	firewall net.Firewall

	connectionManager net.ConnectionManager

	peerCrossListLock sync.Mutex
	peerCrossList     map[string]bool
}

// NewGuard returns a new instance of Guard. Should only be called once per
// provider. Instantiating a new instance of Guard automatically runs it in the
// background for the lifetime of the client.
func NewGuard(
	ctx context.Context,
	logger log.StandardLogger,
	duration time.Duration,
	firewall net.Firewall,
	connectionManager net.ConnectionManager,
) *Guard {
	guard := &Guard{
		logger:            logger,
		duration:          duration,
		firewall:          firewall,
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
			g.logger.Debugf("starting firewall guard round")

			connectedPeers := g.connectionManager.ConnectedPeers()

			for _, connectedPeer := range connectedPeers {
				if g.currentlyChecking(connectedPeer) {
					continue
				}

				// Ensure we mark the peer as being checked before
				// executing the async stake check.
				g.markAsChecking(connectedPeer)
				go g.checkFirewallRules(connectedPeer)
			}
		}
	}
}

func (g *Guard) checkFirewallRules(peer string) {
	defer g.completedCheck(peer)

	peerPublicKey, err := g.getPeerPublicKey(peer)
	if err != nil {
		// if we error while getting the peer's public key, the peer's id
		// or key may be malformed/unknown; disconnect them immediately.
		g.logger.Errorf(
			"dropping the connection; "+
				"could not get public key for peer [%v]: [%v]",
			peer,
			err,
		)
		g.connectionManager.DisconnectPeer(peer)
		return
	}

	if err := g.firewall.Validate(peerPublicKey); err != nil {
		g.logger.Warningf(
			"dropping the connection; "+
				"firewall rules not satisfied for peer [%v]: [%v] ",
			peer,
			err,
		)
		g.connectionManager.DisconnectPeer(peer)
	}
}

func (g *Guard) getPeerPublicKey(peer string) (*operator.PublicKey, error) {
	peerPublicKey, err := g.connectionManager.GetPeerPublicKey(peer)
	if err != nil {
		return nil, err
	}

	if peerPublicKey == nil {
		return nil, fmt.Errorf(
			"failed to resolve valid public key for peer [%s]", peer,
		)
	}
	return peerPublicKey, nil
}
