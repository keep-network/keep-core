package libp2p

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	peer "github.com/libp2p/go-libp2p-peer"
)

func TestCreatePeerIDFromKey(t *testing.T) {
	_, pub, err := identityKeyPair(generateEthereumKey())
	if err != nil {
		t.Fatal(err)
	}

	peerID, err := peer.IDFromPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}

	expectedID := "QmZZEiMy2DPLnPhhG7M9tDh2HcF1g38jZuAU4QKHuqDtBo"
	if peerID.Pretty() != expectedID {
		t.Fatalf(
			"Unexpected peer ID\nExpected: %v\nActual: %v",
			expectedID,
			peerID.Pretty(),
		)
	}
}

func TestSignAndVerify(t *testing.T) {
	priv, pub, err := identityKeyPair(generateEthereumKey())
	if err != nil {
		t.Fatal(err)
	}

	msgDigest := big.NewInt(1410977).Bytes()

	signature, err := priv.Sign(msgDigest)
	if err != nil {
		t.Fatal(err)
	}

	isValid, err := pub.Verify(msgDigest, signature)
	if err != nil {
		t.Fatal(err)
	}

	if !isValid {
		t.Fatal("invalid signature")
	}
}

func TestVerifyInvalidSignature(t *testing.T) {
	priv, pub, err := identityKeyPair(generateEthereumKey())
	if err != nil {
		t.Fatal(err)
	}

	msgDigest := big.NewInt(1410977).Bytes()
	signature, err := priv.Sign(msgDigest)
	if err != nil {
		t.Fatal(err)
	}

	anotherMsgDigest := big.NewInt(2410995).Bytes()
	isValid, err := pub.Verify(anotherMsgDigest, signature)
	if err != nil {
		t.Fatal(err)
	}

	if isValid {
		t.Fatal("invalid signature expected")
	}
}

func generateEthereumKey() *ecdsa.PrivateKey {
	ethCurve := secp256k1.S256()
	privateKey := big.NewInt(1337)

	pubKeyX, pubKeyY := ethCurve.ScalarBaseMult(privateKey.Bytes())

	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: secp256k1.S256(), // go-ethereum curve
			X:     pubKeyX,
			Y:     pubKeyY,
		},
		D: privateKey,
	}
}
