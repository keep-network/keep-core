package result

import (
	"github.com/ethereum/go-ethereum/crypto"
	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/operator"
)

// // TODO: This file is just temporary util signatures are properly implemented in
// // `operator` package.

// Signature is 65-byte signature
type Signature = []byte // TODO: Remove and replace with operator.Signature

// TODO: Remove and replace with operator.Sign
func (fm *SigningMember) sign(resultHash relayChain.DKGResultHash) (Signature, error) {
	// TODO: Implement

	sig, err := crypto.Sign(resultHash[:], fm.privateKey)
	return sig, err
}

// TODO: Remove and replace with operator.VerifySignature
func (fm *SigningMember) verifySignature(
	hash relayChain.DKGResultHash,
	signature Signature,
	publicKey *operator.PublicKey,
) bool {
	if len(signature) != 65 {
		return false
	}

	// TODO: Implement
	return crypto.VerifySignature(
		crypto.CompressPubkey(publicKey),
		hash[:],
		signature[:64],
	)
}
