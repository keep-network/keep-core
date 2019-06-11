package watchtower

import (
	"context"
	"math/big"
	"testing"
	"time"

	localChain "github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net/key"
	localNetwork "github.com/keep-network/keep-core/pkg/net/local"
)

func TestDisconnectPeerBelowMinStake(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	minStake := big.NewInt(200)
	stakeMonitor := localChain.NewStakeMonitor(minStake)

	peer1PubKey, peer1Address, err := createNewPeerIdentity()
	if err != nil {
		t.Fatal(err)
	}

	if err := stakeNetworkPeer(stakeMonitor, peer1Address); err != nil {
		t.Fatal(err)
	}

	// kick off the network
	peer1Provider := localNetwork.Connect()
	// add self to our own map of connections
	peer1Provider.AddPeer(peer1Provider.ID().String(), peer1PubKey)

	// set watchtower for bootstrap peer
	_ = NewGuard(ctx, 1*time.Second, stakeMonitor, peer1Provider.ConnectionManager())

	// initialize second peer
	peer2PubKey, peer2Address, err := createNewPeerIdentity()
	if err != nil {
		t.Fatal(err)
	}

	if err := stakeNetworkPeer(stakeMonitor, peer2Address); err != nil {
		t.Fatal(err)
	}

	// kick off the network
	peer2Provider := localNetwork.Connect()
	peer2Provider.AddPeer(peer2Provider.ID().String(), peer2PubKey)

	// set watchtower for our second peer
	_ = NewGuard(ctx, 1*time.Second, stakeMonitor, peer2Provider.ConnectionManager())

	// Make sure they add each other
	peer1Provider.AddPeer(peer2Provider.ID().String(), peer2PubKey)
	peer2Provider.AddPeer(peer1Provider.ID().String(), peer1PubKey)

	// drop our second peer below the min stake
	if err := stakeMonitor.UnstakeTokens(peer2Address); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)

	// make sure the connection that the first peer has to the second peer has been untethered.
	for _, peer := range peer1Provider.ConnectionManager().ConnectedPeers() {
		if peer == peer2Provider.ID().String() {
			t.Fatal("did not expect connection to peer with min stake")
		}
	}
}

// we need the peer pubkey, the peer address, and provider
func createNewPeerIdentity() (*key.NetworkPublic, string, error) {
	_, peerPublicKey, err := key.GenerateStaticNetworkKey()
	if err != nil {
		return nil, "", err
	}

	return peerPublicKey, key.NetworkPubKeyToEthAddress(peerPublicKey), nil
}

func stakeNetworkPeer(
	stakeMonitor *localChain.StakeMonitor,
	address string,
) error {
	staker, err := stakeMonitor.StakerFor(address)
	if err != nil {
		return err
	}
	if err := stakeMonitor.StakeTokens(address); err != nil {
		return err
	}
	_, err = staker.Stake()
	return err

}
