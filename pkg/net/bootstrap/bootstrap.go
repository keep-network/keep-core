package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	logging "github.com/ipfs/go-log"
	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
	periodicproc "github.com/jbenet/goprocess/periodic"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/routing"
)

var logger = logging.Logger("bootstrap")

// ErrNotEnoughBootstrapPeers signals that we do not have enough bootstrap
// peers to bootstrap correctly.
var ErrNotEnoughBootstrapPeers = errors.New("not enough bootstrap peers to bootstrap")

// BootstrapConfig specifies parameters used in the network bootstrapping process.
type BootstrapConfig struct {
	// Period governs the periodic interval at which the node will
	// attempt to bootstrap. The bootstrap process is not very expensive, so
	// this threshold can afford to be small (<=30s).
	Period time.Duration

	// ConnectionTimeout determines how long to wait for a bootstrap
	// connection attempt before cancelling it.
	ConnectionTimeout time.Duration

	// BootstrapPeers is a function that returns a set of bootstrap peers
	// for the bootstrap process to use. This makes it possible for clients
	// to control the peers the process uses at any moment.
	BootstrapPeers func() []peer.AddrInfo
}

// DefaultBootstrapConfig specifies default sane parameters for bootstrapping.
var DefaultBootstrapConfig = BootstrapConfig{
	Period:            30 * time.Second,
	ConnectionTimeout: (30 * time.Second) / 3, // Perod / 3
}

func BootstrapConfigWithPeers(pis []peer.AddrInfo) BootstrapConfig {
	cfg := DefaultBootstrapConfig
	cfg.BootstrapPeers = func() []peer.AddrInfo {
		return pis
	}
	return cfg
}

// Bootstrap kicks off bootstrapping. This function will periodically
// check the number of open connections and -- if there are too few -- initiate
// connections to well-known bootstrap peers. It also kicks off subsystem
// bootstrapping (i.e. routing).
func Bootstrap(
	id peer.ID,
	host host.Host,
	rt routing.Routing,
	cfg BootstrapConfig,
) (io.Closer, error) {
	// make a signal to wait for one bootstrap round to complete.
	doneWithRound := make(chan struct{})

	// the periodic bootstrap function -- the connection supervisor
	periodic := func(worker goprocess.Process) {
		ctx := goprocessctx.OnClosingContext(worker)

		if err := bootstrapRound(ctx, host, cfg); err != nil {
			logger.Warningf("bootstrap round error: [%v]", err)
		}

		<-doneWithRound
	}

	// kick off the node's periodic bootstrapping
	proc := periodicproc.Tick(cfg.Period, periodic)
	proc.Go(periodic) // run one right now.

	// kick off Routing.Bootstrap
	if rt != nil {
		ctx := goprocessctx.OnClosingContext(proc)
		if err := rt.Bootstrap(ctx); err != nil {
			proc.Close()
			return nil, err
		}
	}

	doneWithRound <- struct{}{}
	close(doneWithRound) // it no longer blocks periodic
	return proc, nil
}

func bootstrapRound(
	ctx context.Context,
	host host.Host,
	cfg BootstrapConfig,
) error {
	ctx, cancel := context.WithTimeout(ctx, cfg.ConnectionTimeout)
	defer cancel()

	logger.Debugf("starting bootstrap round")

	// get bootstrap peers from config. retrieving them here makes
	// sure we remain observant of changes to client configuration.
	peers := cfg.BootstrapPeers()

	if len(peers) == 0 {
		logger.Debugf("bootstrap round skipped; no bootstrap peers in config")
		return nil
	}

	// filter out bootstrap nodes we are already connected to
	var notConnected []peer.AddrInfo
	for _, p := range peers {
		if host.Network().Connectedness(p.ID) != network.Connected {
			notConnected = append(notConnected, p)
		}
	}

	// if connected to all bootstrap peer candidates, exit
	if len(notConnected) < 1 {
		logger.Debugf(
			"bootstrap round skipped; " +
				"connected to all bootstrap peers from config",
		)
		return nil
	}

	logger.Debugf("bootstrapping to nodes: [%v]", notConnected)

	return bootstrapConnect(ctx, host, notConnected)
}

func bootstrapConnect(
	ctx context.Context,
	ph host.Host,
	peers []peer.AddrInfo,
) error {
	if len(peers) < 1 {
		return ErrNotEnoughBootstrapPeers
	}

	errs := make(chan error, len(peers))
	var wg sync.WaitGroup
	for _, p := range peers {

		// performed asynchronously because when performed synchronously, if
		// one `Connect` call hangs, subsequent calls are more likely to
		// fail/abort due to an expiring context.
		// Also, performed asynchronously for dial speed.

		wg.Add(1)
		go func(p peer.AddrInfo) {
			defer wg.Done()

			logger.Debugf("trying to establish connection with bootstrap peer [%v]", p.ID)

			ph.Peerstore().AddAddrs(p.ID, p.Addrs, peerstore.PermanentAddrTTL)

			if err := ph.Connect(ctx, p); err != nil {
				logger.Warningf(
					"could not establish connection with bootstrap peer [%v]: [%v]",
					p.ID,
					err,
				)
				errs <- err
				return
			}

			logger.Debugf("established connection with bootstrap peer [%v]", p.ID)
		}(p)
	}
	wg.Wait()

	// our failure condition is when no connection attempt succeeded.
	// So drain the errs channel, counting the results.
	close(errs)
	count := 0
	var err error
	for err = range errs {
		if err != nil {
			count++
		}
	}
	if count == len(peers) {
		return fmt.Errorf("all bootstrap attempts failed")
	}
	return nil
}

type Peers []peer.AddrInfo

func (p Peers) ToPeerInfos() []peer.AddrInfo {
	pinfos := make(map[peer.ID]*peer.AddrInfo)
	for _, bootstrap := range p {
		pinfo, ok := pinfos[bootstrap.ID]
		if !ok {
			pinfo = new(peer.AddrInfo)
			pinfos[bootstrap.ID] = pinfo
			pinfo.ID = bootstrap.ID
		}

		pinfo.Addrs = append(pinfo.Addrs, bootstrap.Addrs...)
	}

	var peers []peer.AddrInfo
	for _, pinfo := range pinfos {
		peers = append(peers, *pinfo)
	}

	return peers
}
