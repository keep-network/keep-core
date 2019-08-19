package gjkr

import (
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func (mcm *MemberCommitmentsMessage) Commitments() []*bn256.G1 {
	return mcm.commitments
}

func (mcm *MemberCommitmentsMessage) SetCommitments(commitments []*bn256.G1) {
	mcm.commitments = commitments
}
