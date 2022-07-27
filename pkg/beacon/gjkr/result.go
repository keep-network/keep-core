package gjkr

import (
	"fmt"
	"math/big"
	"sync"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// Result of distributed key generation protocol.
type Result struct {
	// Group represents the group state, including members, disqualified,
	// and inactive members.
	Group *group.Group
	// Group public key generated by protocol execution.
	GroupPublicKey *bn256.G2
	// Share of the group private key. It is used for signing and should never
	// be revealed publicly.
	GroupPrivateKeyShare *big.Int

	groupPublicKeySharesMutex   sync.Mutex
	groupPublicKeySharesChannel <-chan map[group.MemberIndex]*bn256.G2
	groupPublicKeyShares        map[group.MemberIndex]*bn256.G2
}

// GroupPublicKeyBytes returns marshalled group public key.
func (r *Result) GroupPublicKeyBytes() ([]byte, error) {
	if r.GroupPublicKey == nil {
		return nil, fmt.Errorf("group public key is nil")
	}

	return r.GroupPublicKey.Marshal(), nil
}

// GroupPublicKeyShares returns shares of the group public key for each
// individual member of the group. They are used for verification of signatures
// received from other members created using their respective group private
// key share.
func (r *Result) GroupPublicKeyShares() map[group.MemberIndex]*bn256.G2 {
	r.groupPublicKeySharesMutex.Lock()
	defer r.groupPublicKeySharesMutex.Unlock()

	if r.groupPublicKeyShares == nil {
		r.groupPublicKeyShares = <-r.groupPublicKeySharesChannel
	}

	return r.groupPublicKeyShares
}
