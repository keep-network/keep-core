package gjkr

import (
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func (mcm *MemberCommitmentsMessage) SetCommitment(
	index int,
	commitment *bn256.G1,
) {
	mcm.commitments[index] = commitment
}

func (ssam *SecretSharesAccusationsMessage) SetAccusedMembersKeys(
	accusedMembersKeys map[group.MemberIndex]*ephemeral.PrivateKey,
) {
	ssam.accusedMembersKeys = accusedMembersKeys
}
