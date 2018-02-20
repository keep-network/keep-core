// Copyright 2018 The Keep Authors.  See LICENSE.md for details.
package node

import (
	"log"
	"strings"
	"time"
	"bufio"
	ma "github.com/multiformats/go-multiaddr"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	"fmt"
)

var (
	store          = datastore.NewMapDatastore()
	ps             = pstore.NewPeerstore()
)

func InitTestHandler(host host.Host) {
	log.Println(fmt.Sprintf(">> InitTestHandler called for %s", host.ID().String()))
}


func InitGetHandler(host host.Host) {
	log.Println(fmt.Sprintf(">> InitGetHandler called for %s", host.ID().String()))
	host.SetStreamHandler("/get/1.0.0", func(s net.Stream) {
		log.Println(">> /get/1.0.0  endpoint called")
		peers := host.Peerstore().Peers()
		peerStrings := []string{}
		for _, peer := range peers {
			addrs := host.Peerstore().Addrs(peer)
			for _, addr := range addrs {
				if peer.Pretty() != "" && addr.String() != "" {
					peerStrings = append(peerStrings, addr.String()+"/ipfs/"+peer.Pretty())
				}
			}
		}
		data := strings.Join(peerStrings, ",") + "\n"
		log.Println(data)
		s.Write([]byte(data))
	})
}

func InitAddHandler(host host.Host) {
	log.Println(fmt.Sprintf(">> InitAddHandler called for %s", host.ID().String()))

	//Adds the current peerlist from the connecting peer and sends back this host's peerlist.
	host.SetStreamHandler("/add/1.0.0", func(s net.Stream) {
		log.Println(">> /add/1.0.0 endpoint called")
		buf := bufio.NewReader(s)
		str, err := buf.ReadString('\n')
		// The following code extracts target's the peer ID from the
		// given multiaddress
		t := strings.TrimSpace(str)

		ipfsaddr, err := ma.NewMultiaddr(t)
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
		if peerid.String() != host.ID().String() {
			if peerid.String() != "" {
				ps.AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)
			}
		}

		peers := ps.Peers()
		peerStrings := []string{}
		for _, peer := range peers {
			addrs := ps.Addrs(peer)
			for _, addr := range addrs {
				if peer.Pretty() != "" && addr.String() != "" {
					peerStrings = append(peerStrings, addr.String()+"/ipfs/"+peer.Pretty())
				}
			}
		}
		data := strings.Join(peerStrings, ",") + "\n"
		s.Write([]byte(data))
	})

	BootstrapConnect(host)

	go func() {
		for {
			log.Println(fmt.Sprintf("Host IP (%v) ID:", host.Network()), host.ID(), host.Peerstore().Peers())
			time.Sleep(time.Duration(5) * time.Second)
			AddPeers(host)
		}
	}()

	//sleep long enough for peerlists to get built.
	time.Sleep(time.Duration(6) * time.Second)
	Test(host)
	select {}
}

