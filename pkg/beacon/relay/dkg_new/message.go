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

// FirstAccusationsMessage is a message payload that carries all of the sender's
// accusations against other members of the threshold group.
// If all other members behaved honestly from the sender's point of view, this message should
// be broadcast but with an empty slice of `accusedIDs`.
// It is expected to be broadcast.
type FirstAccusationsMessage struct {
	senderID *big.Int

	accusedIDs []*big.Int
}
