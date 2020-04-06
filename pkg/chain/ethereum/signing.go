package ethereum

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/chain"
)

// SignatureSize is a byte size of a signature calculated by Ethereum with
// recovery-id, V, included. The signature consist of three values (R,V,S)
// in the following order:
// R = [0:31]
// S = [32:63]
// V = [64]
const SignatureSize = 65

type ethereumSigning struct {
	operatorKey *ecdsa.PrivateKey
}

func (ec *ethereumChain) Signing() chain.Signing {
	return &ethereumSigning{ec.accountKey.PrivateKey}
}

func (es *ethereumSigning) PublicKey() []byte {
	publicKey := es.operatorKey.PublicKey
	return elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)
}

func (es *ethereumSigning) Sign(message []byte) ([]byte, error) {
	prefixedHash := crypto.Keccak256(
		[]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%v", len(message))),
		message,
	)

	signature, err := crypto.Sign(prefixedHash, es.operatorKey)
	if err != nil {
		return nil, err
	}

	if len(signature) == SignatureSize {
		// go-ethereum/crypto produces signature with v={0, 1} and we need to add
		// 27 to v-part (signature[64]) to conform wtih the on-chain signature
		// validation code that accepts v={27, 28} as specified in the
		// Appendix F of the Ethereum Yellow Paper
		// https://ethereum.github.io/yellowpaper/paper.pdf
		signature[len(signature)-1] = signature[len(signature)-1] + 27
	}

	return signature, nil
}

func (es *ethereumSigning) Verify(message []byte, signature []byte) (bool, error) {
	return verifySignature(message, signature, &es.operatorKey.PublicKey)
}

func (es *ethereumSigning) VerifyWithPublicKey(
	message []byte,
	signature []byte,
	publicKey []byte,
) (bool, error) {
	unmarshalledPubKey, err := unmarshalPublicKey(
		publicKey,
		es.operatorKey.Curve,
	)
	if err != nil {
		return false, err
	}

	return verifySignature(message, signature, unmarshalledPubKey)
}

func verifySignature(
	message []byte,
	signature []byte,
	publicKey *ecdsa.PublicKey,
) (bool, error) {
	// Convert the operator's static key into an uncompressed public key
	// which should be 65 bytes in length.
	uncompressedPubKey := crypto.FromECDSAPub(publicKey)
	// If our signature is in the [R || S || V] format, ensure we strip out
	// the Ethereum-specific recovery-id, V, if it already hasn't been done.
	if len(signature) == SignatureSize {
		signature = signature[:len(signature)-1]
	}

	// The signature should be now 64 bytes long.
	if len(signature) != 64 {
		return false, fmt.Errorf(
			"signature should have 64 bytes; has: [%v]",
			len(signature),
		)
	}

	prefixedHash := crypto.Keccak256(
		[]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%v", len(message))),
		message,
	)

	return crypto.VerifySignature(
		uncompressedPubKey,
		prefixedHash,
		signature,
	), nil
}

func unmarshalPublicKey(
	bytes []byte,
	curve elliptic.Curve,
) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(curve, bytes)
	if x == nil {
		return nil, fmt.Errorf(
			"invalid public key bytes",
		)
	}
	ecdsaPublicKey := &ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	return (*ecdsa.PublicKey)(ecdsaPublicKey), nil
}

func (es *ethereumSigning) PublicKeyToAddress(publicKey ecdsa.PublicKey) []byte {
	return crypto.PubkeyToAddress(publicKey).Bytes()
}

func (es *ethereumSigning) PublicKeyBytesToAddress(publicKey []byte) []byte {
	// Does the same as crypto.PubkeyToAddress but directly on public key bytes.
	return crypto.Keccak256(publicKey[1:])[12:]
}
