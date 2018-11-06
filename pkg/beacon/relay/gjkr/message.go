package gjkr

import "math/big"

// MemberCommitmentsMessage is a message payload that carries the sender's
// commitments to polynomial coefficients during distributed key generation.
//
// It is expected to be broadcast.
type MemberCommitmentsMessage struct {
	senderID int

	commitments []*big.Int // slice of `C_ik`
}

// PeerSharesMessage is a message payload that carries shares `s_ij` and `t_ij`
// calculated by the sender `i` for the recipient `j` during distributed key
// generation.
//
// It is expected to be communicated in an encrypted fashion to the selected
// recipient.
type PeerSharesMessage struct {
	senderID   int // i
	receiverID int // j

	shareS *big.Int // s_ij
	shareT *big.Int // t_ij
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

// MemberPublicCoefficientsMessage is a message payload that carries the sender's
// public coefficients.
// It is expected to be broadcast.
type MemberPublicCoefficientsMessage struct {
	senderID int

	publicCoefficients []*big.Int // A_ik = g^{a_ik} mod p
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
