package result

import (
	"crypto/ecdsa"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/chain"
)

// MemberIndex is an index of a participant in the group.
// Indexing starts with `1`.
type MemberIndex uint32

// SigningMember represents a member sharing preferred DKG result hash
// and signature over this hash with peer members.
type SigningMember struct {
	index MemberIndex

	chainHandle chain.Handle

	// Keys used for signing the DKG result hash.
	privateKey             *ecdsa.PrivateKey                // TODO: Change to static.PrivateKey
	otherMembersPublicKeys map[MemberIndex]*ecdsa.PublicKey // TODO: Change to static.PrivateKey

	// Hash of DKG result preferred by the current participant.
	preferredDKGResultHash relayChain.DKGResultHash
	// Received valid signatures supporting the same DKG result as current's
	// participants prefer. Contains also current's participant's signature.
	receivedValidResultSignatures map[MemberIndex]Signature // TODO: Change to static.Signature
}
