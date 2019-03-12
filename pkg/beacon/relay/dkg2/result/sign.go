package result

import (
	"github.com/ethereum/go-ethereum/crypto"
	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// TODO: This file is just temporary util signatures are properly implemented in
// `static` package.

// Signature is 65-byte signature
type Signature = []byte // TODO: Remove and replace with static.Signature

// TODO: Remove and replace with static.Sign
func (fm *ResultSigningMember) sign(resultHash relayChain.DKGResultHash) (Signature, error) {
	// TODO: Implement

	sig, err := crypto.Sign(resultHash[:], fm.privateKey)
	return sig, err
}

// TODO: Remove and replace with static.VerifySignature
func (fm *ResultSigningMember) verifySignature(
	otherMemberIndex MemberIndex,
	hash relayChain.DKGResultHash,
	signature Signature,
) bool {
	if len(signature) != 65 {
		return false
	}

	// TODO: Implement
	return crypto.VerifySignature(
		crypto.CompressPubkey(fm.otherMembersPublicKeys[otherMemberIndex]),
		hash[:],
		signature[:64],
	)
}
