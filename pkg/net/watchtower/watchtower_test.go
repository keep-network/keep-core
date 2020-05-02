package watchtower

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/net/key"
	localNetwork "github.com/keep-network/keep-core/pkg/net/local"
)

func TestDisconnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, peer1PublicKey, err := key.GenerateStaticNetworkKey()
	if err != nil {
		t.Fatal(err)
	}
	_, peer2PublicKey, err := key.GenerateStaticNetworkKey()
	if err != nil {
		t.Fatal(err)
	}

	firewall := newMockFirewall()
	firewall.updatePeer(peer1PublicKey, true)
	firewall.updatePeer(peer2PublicKey, true)

	// setup the first peer
	peer1Provider := localNetwork.Connect()
	_ = NewGuard(ctx, 1*time.Second, firewall, peer1Provider.ConnectionManager())

	// setup the second peer
	peer2Provider := localNetwork.Connect()
	_ = NewGuard(ctx, 1*time.Second, firewall, peer2Provider.ConnectionManager())

	// connect them with each other
	peer1Provider.AddPeer(peer2Provider.ID().String(), peer2PublicKey)
	peer2Provider.AddPeer(peer1Provider.ID().String(), peer1PublicKey)

	// make sure they are connected
	if len(peer1Provider.ConnectionManager().ConnectedPeers()) != 1 {
		t.Fatal("peer 1 not connected properly with peer 2")
	}
	if len(peer2Provider.ConnectionManager().ConnectedPeers()) != 1 {
		t.Fatal("peers 2 not connected properly with peer 1")
	}

	// cut off the second peer in the firewall
	firewall.updatePeer(peer2PublicKey, false)

	// two seconds to run the validation loop
	time.Sleep(2 * time.Second)

	// peer 1 should drop the connection with peer 2
	if len(peer1Provider.ConnectionManager().ConnectedPeers()) != 0 {
		t.Fatal("peer 1 should drop the connection with peer 2")
	}
}

func newMockFirewall() *mockFirewall {
	return &mockFirewall{
		meetsCriteria: make(map[uint64]bool),
	}
}

type mockFirewall struct {
	meetsCriteria map[uint64]bool
}

func (mf *mockFirewall) Validate(remotePeerPublicKey *ecdsa.PublicKey) error {
	if !mf.meetsCriteria[remotePeerPublicKey.X.Uint64()] {
		return fmt.Errorf("remote peer does not meet firewall criteria")
	}
	return nil
}

func (mf *mockFirewall) updatePeer(
	remotePeerPublicKey *key.NetworkPublic,
	meetsCriteria bool,
) {
	x := key.NetworkKeyToECDSAKey(remotePeerPublicKey).X.Uint64()
	mf.meetsCriteria[x] = meetsCriteria
}
