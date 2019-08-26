package gjkr

import (
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func (mcm *MemberCommitmentsMessage) SetCommitment(
	index int,
	commitment *bn256.G1,
) {
	mcm.commitments[index] = commitment
}
