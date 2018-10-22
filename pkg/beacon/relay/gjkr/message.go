package gjkr

import "math/big"

// MemberCommitmentsMessage is a message payload that carries the sender's
// commitments to polynomial coefficients during distributed key generation.
//
// It is expected to be broadcast.
type MemberCommitmentsMessage struct {
	senderID int

	commitments []*big.Int
}

// PeerSharesMessage is a message payload that carries the sender's shares `s_ij`
// and `t_ij` for the recipient during distributed key generation.
//
// It is expected to be communicated in encrypted fashion to the recipient over
// a broadcast channel.
type PeerSharesMessage struct {
	senderID   int
	receiverID int

	shareS *big.Int
	shareT *big.Int
}

// SecretSharesAccusationsMessage is a message payload that carries all of the
// sender's accusations against other members of the threshold group.
// If all other members behaved honestly from the sender's point of view, this
// message should be broadcast but with an empty slice of `accusedIDs`.
//
// It is expected to be broadcast.
type SecretSharesAccusationsMessage struct {
	senderID int

	accusedIDs []int
}

// MemberPublicCoefficientsMessage is a message payload that carries she sender's
// public coefficients.
// It is expected to be broadcast.
type MemberPublicCoefficientsMessage struct {
	senderID int

	publicCoefficients []*big.Int
}

// CoefficientsAccusationsMessage is a message payload that carries all of the
// sender's accusations against other members of the threshold group after public
// coefficients validation.
// If all other members behaved honestly from the sender's point of view, this
// message should be broadcast but with an empty slice of `accusedIDs`.
// It is expected to be broadcast.
type CoefficientsAccusationsMessage struct {
	senderID int

	accusedIDs []int
}
