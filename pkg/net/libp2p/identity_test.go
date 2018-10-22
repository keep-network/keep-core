package libp2p

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	peer "github.com/libp2p/go-libp2p-peer"
)

// Create ECDSA public key using Ethereum's secp256k1 curve as defined in the
// `go-ethereum` library and see if the created peer ID has an expected value.
//
// `geth` client uses `go-ethereum` library to generate ECDSA key. The library
// uses its definition of secp256k1 curve. LibP2P's `peer.IDFromPublicKey`
// currently expects that the curve must be from Go standard library and keys
// created with curve defined by `go-ethereum` can not be turned into bytes as
// this operation yields an error.
//
// If this test fails, please see if libp2p strategy for multihashing key
// bytes did not change. The current rule is as follows:
//
// * When `len(bytes) <= peer.MaxInlineKeyLength`, the peer ID is the identity
//   multihash of bytes.
// * When `len(bytes) > peer.MaxInlineKeyLength`, the peer ID is the sha2-256
//   multihash of bytes.
//
// Currently, `MaxInlineKeyLength = 42`.
//
// Please bear in mind that key's multihash is a part of libp2p multiaddress.
// Any change to the multihash entails the need of updating peer's multiaddress.
func TestCreatePeerIdFromEthereumKey(t *testing.T) {
	ethCurve := secp256k1.S256()
	pubKeyX, pubKeyY := ethCurve.ScalarBaseMult(big.NewInt(1337).Bytes())
	ethKey := &ecdsa.PublicKey{
		Curve: secp256k1.S256(), // go-ethereum curve
		X:     pubKeyX,
		Y:     pubKeyY,
	}

	peerID := peerIDFromPublicKey(ethKey)

	expectedID := "QmZZEiMy2DPLnPhhG7M9tDh2HcF1g38jZuAU4QKHuqDtBo"
	if peerID.Pretty() != expectedID {
		t.Fatalf(
			"Unexpected peer ID\nExpected: %v\nActual: %v",
			expectedID,
			peerID.Pretty(),
		)
	}
}

// When `len(bytes) <= peer.MaxInlineKeyLength`,
// the peer ID is the identity multihash of bytes.
func TestBytesUnderMaxInlineKeyLengthToId(t *testing.T) {
	bytes := [3]byte{0x01, 0x02, 0x03}
	peerID := bytesToPeerID(bytes[:])

	expectedID := "15TJUr"
	if peerID.Pretty() != expectedID {
		t.Fatalf(
			"Unexpected peer ID\nExpected: %v\nActual: %v",
			expectedID,
			peerID.Pretty(),
		)
	}
}

// When `len(bytes) > peer.MaxInlineKeyLength`,
// the peer ID is the sha2-256 multihash of bytes.
func TestBytesOverMaxInlineKeyLengthToId(t *testing.T) {
	bytes := make([]byte, peer.MaxInlineKeyLength+10)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = 0x99
	}

	peerID := bytesToPeerID(bytes)

	expectedID := "Qmca2aYUUWAFy2mAvzL9Yj9FjrcT1qmuk9nAsjaa5S4jVF"
	if peerID.Pretty() != expectedID {
		t.Fatalf(
			"Unexpected peer ID\nExpected: %v\nActual: %v",
			expectedID,
			peerID.Pretty(),
		)
	}
}
