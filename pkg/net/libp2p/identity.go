package libp2p

import (
	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

func PubKeyFromID(pi peerIdentifier) (ci.PubKey, error) {
	pid := peer.ID(pi)
	return pid.ExtractPublicKey()
}

// AddIdentityToStore takes a peerIdentity and notifies the addressbook of the
// existance of a new client joining the network.
func addIdentityToStore(pi peerIdentifier) (pstore.Peerstore, error) {
	// TODO: investigate a generic store interface that gives us a unified interface
	// to our address book (peerstore in libp2p) from secure storage (dht)
	ps := pstore.NewPeerstore()
	_ = peer.ID(pi)

	// if err := ps.AddPrivKey(id, pi.privKey); err != nil {
	// 	return nil, fmt.Errorf("failed to add PrivateKey with error %s", err)
	// }
	// if err := ps.AddPubKey(id, pi.privKey.GetPublic()); err != nil {
	// 	return nil, fmt.Errorf("failed to add PubKey with error %s", err)
	// }
	return ps, nil
}
