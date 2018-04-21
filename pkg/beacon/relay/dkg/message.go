package dkg

import (
	"github.com/dfinity/go-dfinity-crypto/bls"
)

// JoinMessage is an empty message payload indicating a member has joined. The
// sender is the joining member. It is expected to be broadcast publicly.
type JoinMessage struct {
	id *bls.ID
}

// MemberCommitmentsMessage is a message payload that carries the sender's
// public commitments during distributed key generation. It is expected to be
// broadcast publicly.
type MemberCommitmentsMessage struct {
	id          *bls.ID
	Commitments []bls.PublicKey
}

// MemberShareMessage is a message payload that carries the sender's private
// share for the recipient during distributed key generation. It is expected to
// be communicated in encrypted fashion to the recipient over a broadcast
// channel.
type MemberShareMessage struct {
	id         *bls.ID
	receiverID *bls.ID
	Share      *bls.SecretKey
}

// AccusationsMessage is a message payload that carries all of the sender's
// accusations against other members of the threshold group. If all other
// members behaved honestly from the sender's point of view, this message should
// be broadcast but with an empty slice of `accusedIDs`. It is expected to be
// broadcast.
type AccusationsMessage struct {
	id         *bls.ID
	accusedIDs []bls.ID
}

// JustificationsMessage is a message payload that carries all of the sender's
// justifications in response to other threshold group members' accusations. If
// no other member accused the sender, this message should be broadcast but with
// an empty map of `justifications`. It is expected to be broadcast.
type JustificationsMessage struct {
	id             *bls.ID
	justifications map[bls.ID]bls.SecretKey
}
