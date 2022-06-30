package bootstrap

import (
	"fmt"
	"github.com/ipfs/go-ipfs-config"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/test"
	"testing"
)

func TestMultipleAddrsPerPeer(t *testing.T) {
	var bsps []peer.AddrInfo
	for i := 0; i < 10; i++ {
		pid, err := test.RandPeerID()
		if err != nil {
			t.Fatal(err)
		}

		addr := fmt.Sprintf("/ip4/127.0.0.1/tcp/5001/ipfs/%s", pid.Pretty())
		bsp1, err := config.ParseBootstrapPeers([]string{addr})
		if err != nil {
			t.Fatal(err)
		}

		addr = fmt.Sprintf("/ip4/127.0.0.1/udp/5002/utp/ipfs/%s", pid.Pretty())
		bsp2, err := config.ParseBootstrapPeers([]string{addr})
		if err != nil {
			t.Fatal(err)
		}

		bsps = append(bsps, bsp1[0], bsp2[0])
	}

	pinfos := Peers.ToPeerInfos(bsps)
	if len(pinfos) != len(bsps)/2 {
		t.Fatal("expected fewer peers")
	}
}
