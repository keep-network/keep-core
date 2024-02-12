package libp2p

import (
	"fmt"
	"github.com/ipfs/go-ipfs-config"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/test"
	"testing"
)

func TestMultipleAddrsPerPeer(t *testing.T) {
	var bsps []peer.AddrInfo
	for i := 0; i < 10; i++ {
		pid, err := test.RandPeerID()
		if err != nil {
			t.Fatal(err)
		}

		addr := fmt.Sprintf("/ip4/127.0.0.1/tcp/5001/ipfs/%s", pid.String())
		bsp1, err := config.ParseBootstrapPeers([]string{addr})
		if err != nil {
			t.Fatal(err)
		}

		addr = fmt.Sprintf("/ip4/127.0.0.1/udp/5002/utp/ipfs/%s", pid.String())
		bsp2, err := config.ParseBootstrapPeers([]string{addr})
		if err != nil {
			t.Fatal(err)
		}

		bsp1Addr, err := peer.AddrInfoFromP2pAddr(bsp1[0].Multiaddr())
		if err != nil {
			t.Fatal(err)
		}

		bsp2Addr, err := peer.AddrInfoFromP2pAddr(bsp2[0].Multiaddr())
		if err != nil {
			t.Fatal(err)
		}

		bsps = append(bsps, *bsp1Addr, *bsp2Addr)
	}

	pinfos := peers.toPeerInfos(bsps)
	if len(pinfos) != len(bsps)/2 {
		t.Fatal("expected fewer peers")
	}
}
