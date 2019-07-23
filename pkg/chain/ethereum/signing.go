package ethereum

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/keep-network/keep-core/pkg/chain"
)

type ethereumSigning struct {
	operatorKey *ecdsa.PrivateKey
}

// SignatureSize is a byte size of a signature calculated by Ethereum with
// recovery-id, V, included.
const SignatureSize = 65

func (ec *ethereumChain) Signing() chain.Signing {
	return &ethereumSigning{ec.accountKey.PrivateKey}
}

func (es *ethereumSigning) Sign(hash []byte) ([]byte, error) {
	prefixedHash := crypto.Keccak256(
		[]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%v", len(hash))),
		hash,
	)

	return crypto.Sign(prefixedHash, es.operatorKey)
}

func (es *ethereumSigning) Verify(hash []byte, signature []byte) (bool, error) {
	sig := signature

	// Convert the operator's static key into an uncompressed public key
	// which should be 65 bytes in length.
	uncompressedPubKey := crypto.FromECDSAPub(&es.operatorKey.PublicKey)
	// If our sig is in the [R || S || V] format, ensure we strip out
	// the Ethereum-specific recovery-id, V, if it already hasn't been done.
	if len(sig) == SignatureSize {
		sig = sig[:len(sig)-1]
	}

	// The sig should be now 64 bytes long.
	if len(sig) != 64 {
		return false, fmt.Errorf(
			"signature should have 64 bytes; has: [%v]",
			len(sig),
		)
	}

	prefixedHash := crypto.Keccak256(
		[]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%v", len(hash))),
		hash,
	)

	return crypto.VerifySignature(
		uncompressedPubKey,
		prefixedHash,
		sig[:],
	), nil
}
