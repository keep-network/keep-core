// Copyright 2018 The Keep Authors.  See LICENSE.md for details.
package node

import (
	"context"
	"fmt"
	"io"
	"log"
	"crypto/rand"
	mrand "math/rand"
	yamux "github.com/whyrusleeping/go-smux-yamux"
	swarm "github.com/libp2p/go-libp2p-swarm"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	crypto2 "github.com/libp2p/go-libp2p-crypto"
	msmux "github.com/whyrusleeping/go-smux-multistream"
	peer "github.com/libp2p/go-libp2p-peer"
	ma "github.com/multiformats/go-multiaddr"
	host "github.com/libp2p/go-libp2p-host"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	"net"
)

var MyIpAddress string

func init() {
	MyIpAddress = GetOutboundIP()
}

//func StartWorkers() {
//	go StatsWorker()
//}


// Get preferred outbound ip of this container
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return fmt.Sprintf("%v", localAddr.IP)
}


// makeBasicHost creates a LibP2P host with a random peer ID listening on the given multiaddress.
// It will use secio if secio is true.
func MakeBasicHost(ps pstore.Peerstore, listenPort int, secio bool, randseed int64) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a deterministic
	// randomness source to make generated keys stay the same across multiple runs.
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it at least to obtain a valid host ID.
	priv, pub, err := crypto2.GenerateKeyPairWithReader(crypto2.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	// Obtain Peer ID from public key
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	multiAddress := fmt.Sprintf("/ip4/%s/tcp/%d", MyIpAddress, listenPort)

	println(" My IP is " + MyIpAddress)
	// multiAddress example: /ip4/127.0.0.1/tcp/2700
	addr, err := ma.NewMultiaddr(multiAddress)

	if err != nil {
		return nil, err
	}

	// If using secio, we add the keys to the peerstore
	// for this peer ID.
	if secio {
		ps.AddPrivKey(pid, priv)
		ps.AddPubKey(pid, pub)
	}

	// Set up stream multiplexer
	tpt := msmux.NewBlankTransport()
	tpt.AddTransport("/yamux/1.0.0", yamux.DefaultTransport)

	// Create swarm (implements libP2P Network)
	swrm, err := swarm.NewSwarmWithProtector(
		context.Background(),
		[]ma.Multiaddr{addr},
		pid,
		ps,
		nil,
		tpt,
		nil,
	)
	if err != nil {
		return nil, err
	}

	netw := (*swarm.Network)(swrm)

	basicHost := bhost.New(netw)

	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host by encapsulating both addresses:
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)

	//StartWorkers()

	return basicHost, nil
}

