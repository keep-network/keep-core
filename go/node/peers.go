// Copyright 2018 The Keep Authors.  See LICENSE.md for details.
package node

import (
	"context"
	"log"
	"strings"
	peer "github.com/libp2p/go-libp2p-peer"
	host "github.com/libp2p/go-libp2p-host"
	ma "github.com/multiformats/go-multiaddr"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

import "bufio"

//AddPeers opens a connection to current peers and requests their peer lists be merged with the initiating host.
func AddPeers(ha host.Host) {
	peers := ha.Peerstore().Peers()
	for _, p := range peers {
		if ha.ID().String() != p.String() {
			stream, err := ha.NewStream(context.Background(), p, "/add/1.0.0")
			if err != nil {
				continue
			}

			for _, addr := range ha.Addrs() {
				if addr.String() != "" {
					_, err = stream.Write([]byte(addr.String() + "/ipfs/" + ha.ID().Pretty() + "\n"))
					if err != nil {
						log.Println(err)
					}
				}
			}
			buf := bufio.NewReader(stream)
			str, err := buf.ReadString('\n')
			// The following code extracts target's the peer ID from the
			// given multiaddress
			t := strings.TrimSpace(str)
			addresses := strings.Split(t, ",")
			for _, address := range addresses {

				ipfsaddr, err := ma.NewMultiaddr(address)
				if err != nil {
					log.Println(err)
				}
				pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
				if err != nil {
					log.Println(err)
				}

				peerid, err := peer.IDB58Decode(pid)
				if err != nil {
					log.Println(err)
				}
				// Decapsulate the /ipfs/<peerID> part from the target
				// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
				targetAddr := ipfsaddr.Decapsulate(ipfsaddr)
				ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)
			}
		}
	}
}

