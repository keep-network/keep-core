package tecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func TestCustomSignatureVerification(t *testing.T) {
	message := sha256.Sum256([]byte("test message"))
	hash := message[:32]
	curve := secp256k1.S256()

	publicKeyX, _ := new(big.Int).SetString("75085108144671174812571143296826646504308369233975693874959462551479380474454", 10)
	publicKeyY, _ := new(big.Int).SetString("65794673518851427144172906713993768445535249202674357005551972276007493788366", 10)

	r, _ := new(big.Int).SetString("4098436129741618965855800457503623013651109799740130558508100814753065300475", 10)
	s, _ := new(big.Int).SetString("46452259181239170314122798282699560779513018398282723793779377265488882483371", 10)
	v := byte(1) // Recovery ID - 0 or 1

	if result, err := verifySignatureInBitcoin(curve, hash, publicKeyX, publicKeyY, r, s); !result || err != nil {
		t.Fatalf("Signature verification in btcec failed [%s]", err)
	}

	if result, err := verifySignatureInEthereum(curve, hash, publicKeyX, publicKeyY, r, s, v); !result || err != nil {
		t.Fatalf("Signature verification in ethereum failed [%s]", err)
	}
}

func verifySignatureInBitcoin(
	curve elliptic.Curve,
	hash []byte,
	publicKeyX,
	publicKeyY,
	signatureR *big.Int,
	signatureS *big.Int,
) (bool, error) {
	publicKey := &btcec.PublicKey{
		Curve: curve,
		X:     publicKeyX,
		Y:     publicKeyY,
	}

	sig := &btcec.Signature{R: signatureR, S: signatureS}

	// Verify if signature is valid for given hash and public key
	if !sig.Verify(hash, publicKey) {
		return false, fmt.Errorf("Signature verification failed")
	}

	// Serialize type Signature {R,S} to DER format supported by Bitcoin:
	// 0x30 <length> 0x02 <length r> r 0x02 <length s> s
	sigSerialized := sig.Serialize()

	// Deserialize signature in DER format to a Signature type {R,S}
	sigDeserialized, err := btcec.ParseDERSignature(sigSerialized, curve)
	if err != nil {
		return false, err
	}

	// Validate deserialized signature matches original signature
	if !sigDeserialized.IsEqual(sig) {
		return false, fmt.Errorf("Signatures are not equal")
	}
	return true, nil
}

func verifySignatureInEthereum(
	curve elliptic.Curve,
	hash []byte,
	publicKeyX,
	publicKeyY,
	signatureR *big.Int,
	signatureS *big.Int,
	signatureV byte,
) (bool, error) {
	publicKey := &ecdsa.PublicKey{
		Curve: curve,
		X:     publicKeyX,
		Y:     publicKeyY,
	}

	signatureRS := append(signatureR.Bytes(), signatureS.Bytes()...)

	// Verify Signature
	if !crypto.VerifySignature(crypto.CompressPubkey(publicKey), hash, signatureRS) {
		return false, fmt.Errorf("Signature verification failed")
	}
	signatureRSV := append(signatureRS, signatureV)

	recoveredPublicKey, err := crypto.SigToPub(hash, signatureRSV)
	if err != nil {
		return false, fmt.Errorf("Recovering public key failed [%s]", err)
	}
	if !reflect.DeepEqual(recoveredPublicKey, publicKey) {
		return false, fmt.Errorf("Recovered Public Key doesn't match expected")
	}
	return true, nil
}

// We might not need this test
func TestBitcoinCompactSignature(t *testing.T) {
	hash := []byte("test message")
	curve := btcec.S256()

	// Generate Private and Public keys
	privateKey, _ := btcec.NewPrivateKey(curve)
	publicKey := privateKey.PubKey()
	isPubKeyCompressed := false

	// Sign the message
	sig, err := btcec.SignCompact(curve, privateKey, hash, isPubKeyCompressed)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v", sig)

	recoveredKey, wasCompressed, err := btcec.RecoverCompact(curve, sig, hash)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(publicKey, recoveredKey) {
		t.Error("Recovered key doesn't match original")
	}
	if wasCompressed != isPubKeyCompressed {
		t.Errorf("recovered pubkey doesn't match compressed state (%v vs %v)", isPubKeyCompressed, wasCompressed)
		return
	}
}

// We might not need this test
func TestBtecSignatureBtcecVerification(t *testing.T) {
	hash := []byte("test message")
	curve := btcec.S256()

	// Generate Private and Public keys
	privateKey, _ := btcec.NewPrivateKey(curve)
	publicKey := privateKey.PubKey()

	// Sign the message
	sig, err := privateKey.Sign(hash)
	if err != nil {
		t.Fatal(err)
	}

	if result, err := verifySignatureInBitcoin(curve, hash, publicKey.X, publicKey.Y, sig.R, sig.S); !result || err != nil {
		t.Fatalf("Signature verification in btcec failed [%s]", err)
	}
}
