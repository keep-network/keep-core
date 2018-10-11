package dkg

import "math/big"

// MemberCommitmentsMessage is a message payload that carries the sender's
// shares calculates for each group member and commitments during distributed
// key generation.
// It is expected to be encrypted before broadcast.
type MemberCommitmentsMessage struct {
	senderID *big.Int

	sharesS     map[*big.Int]*big.Int
	sharesT     map[*big.Int]*big.Int
	commitments []*big.Int
}
