package dkg

import "math/big"

// MemberCommitmentsMessage is a message payload that carries the sender's
// commitments to polynomial coefficients during distributed key generation.
// It is expected to be broadcast.
type MemberCommitmentsMessage struct {
	senderID *big.Int

	commitments []*big.Int
}

// PeerSharesMessage is a message payload that carries the sender's shares `s_ij`
// and `t_ij` for the recipient during distributed key generation.
// It is expected to be communicated in encrypted fashion to the recipient over
// a broadcast channel.
type PeerSharesMessage struct {
	senderID   *big.Int
	receiverID *big.Int

	shareS *big.Int
	shareT *big.Int
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
