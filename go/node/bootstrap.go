// Copyright 2018 The Keep Authors.  See LICENSE.md for details.
package node

import (
	"fmt"
	"log"

	host "github.com/libp2p/go-libp2p-host"
	ma "github.com/multiformats/go-multiaddr"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	. "github.com/keep-network/keep-core/libs/config"
)

//BootstrapConnect connects to a predefined list of peers to "bootstrap" the cluster.
func BootstrapConnect(ha host.Host) {
	log.Println("Bootstrapping peers...")
	for _, p := range Config.Bootstrap.Nodes {
		// The following code extracts target's the peer ID from the given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(p)
		if err != nil {
			log.Fatalln(err)
		}
		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
		}
		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
		}
		// Decapsulate the /ipfs/<peerID> part from the target  /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)
		haId := ha.ID().String()
		log.Println("haId", haId)
		if ha.ID().String() != peerid.String() {
			// We have a peer ID and a targetAddr so we add it to the peerstore so LibP2P knows how to contact it
			ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)
		}
	}
}
