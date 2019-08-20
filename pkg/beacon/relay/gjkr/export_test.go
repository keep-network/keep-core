package gjkr

import (
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func (mcm *MemberCommitmentsMessage) Commitments() []*bn256.G1 {
	return mcm.commitments
}

func (mcm *MemberCommitmentsMessage) SetCommitments(commitments []*bn256.G1) {
	mcm.commitments = commitments
}

func (ssam *SecretSharesAccusationsMessage) SetAccusedMembersKeys(
	accusedMembersKeys map[group.MemberIndex]*ephemeral.PrivateKey,
) {
	ssam.accusedMembersKeys = accusedMembersKeys
}
