package lib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// HashPublicKeyToAddress takes the bytes for the public key and returns
// an Ethereum address.  commmon.Address is a [32]byte or [AddressLength]byte
func HashPublicKeyToAddress(pk []byte) (*common.Address, error) {
	// pubKey := ecdsa.PublicKey{}
	pubKey, err := UnmarshalPubkey(pk)
	// func UnmarshalPubkey(pub []byte) (*ecdsa.PublicKey, error) {
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(*pubKey)
	return &address, nil
}

// ValidateSignature checks a signature that has been signed with a public key to see if
// it is authentic.  Return true if it is a valid signature.
func ValidateSignature(publicKey *ecdsa.PublicKey, rawSignature []byte) bool {
	r, s, signature, err := UnPackSignature(rawSignature)
	if err != nil {
		log.Printf("unpack of signature resulted in an error. %s", err)
		return false
	}
	return ecdsa.Verify(publicKey, signature, r, s)
}

// PackSignature packs the 'r' and 's' values in with the signature data in same
// format used in Bitcoin.  This is with a 4 byte hex length, then the 'r', then
// a 4 byte hex length, then the 's', then a 4 byte hex and the hex of the
// signature data.
func PackSignature(r, s *big.Int, sigData []byte) []byte {
	rHex := toHexInt(r)
	sHex := toHexInt(s)
	sdHex := fmt.Sprintf("%x", sigData)
	return []byte(fmt.Sprintf("%04x%s%04x%s%04x%x", len(rHex), rHex, len(sHex), sHex, len(sdHex), sdHex))
}

func UnPackSignature(packedSignature []byte) (*big.Int, *big.Int, []byte, error) {
	rLenStr := string(packedSignature[0:4])
	rLen, err := strconv.ParseInt(rLenStr, 16, 32)
	if err != nil {
		return nil, nil, []byte{}, err
	}
	rStr := string(packedSignature[4 : 4+rLen])
	r, ok := big.NewInt(0).SetString(rStr, 16)
	if !ok {
		return nil, nil, []byte{}, fmt.Errorf("xyzzy")
	}
	cur := 4 + rLen
	sLenStr := string(packedSignature[cur : cur+4])
	sLen, err := strconv.ParseInt(sLenStr, 16, 32)
	if err != nil {
		return nil, nil, []byte{}, err
	}
	sStr := string(packedSignature[cur+4 : cur+4+sLen])
	s, ok := big.NewInt(0).SetString(sStr, 16)
	if !ok {
		return nil, nil, []byte{}, fmt.Errorf("xyzzy")
	}
	cur += 4 + sLen
	sigLenStr := string(packedSignature[cur : cur+4])
	sigLen, err := strconv.ParseInt(sigLenStr, 16, 32)
	if err != nil {
		return nil, nil, []byte{}, err
	}
	sigData := packedSignature[cur+4 : cur+4+sigLen]
	return r, s, sigData, nil
}

func toHexInt(n *big.Int) string {
	return fmt.Sprintf("%x", n)
}

// UnmarshalPubkey converts bytes to a secp256k1 public key.
func UnmarshalPubkey(pub []byte) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(S256(), pub)
	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return &ecdsa.PublicKey{Curve: S256(), X: x, Y: y}, nil
}

// S256 returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return secp256k1.S256()
}
