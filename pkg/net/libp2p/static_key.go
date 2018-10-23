package libp2p

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
)

func loadSecp256k1Key(ethereumKey *keystore.Key) libp2pcrypto.PrivKey {
	privKey, _ := btcec.PrivKeyFromBytes(
		btcec.S256(), ethereumKey.PrivateKey.D.Bytes(),
	)

	return (*libp2pcrypto.Secp256k1PrivateKey)(privKey)
}
