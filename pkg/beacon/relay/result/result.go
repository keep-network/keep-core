package result

import (
	"fmt"
	"math/big"
)

// Result of distributed key generation protocol.
//
// Success means that the protocol execution finished with acceptable number of
// disqualified or inactive members. The group of remaining members should be
// added to the signing groups for the threshold relay.
//
// Failure means that the group creation could not finish, due to either the number
// of inactive or disqualified participants, or the presented results being
// disputed in a way where the correct outcome cannot be ascertained.
type Result struct {
	// Result type of the protocol execution. True if success, false if failure.
	Success bool
	// Group public key generated by protocol execution.
	GroupPublicKey *big.Int `json:"pubkey"`
	// Disqualified members IDs.
	Disqualified []int `json:"disqualified"`
	// Inactive members IDs.
	Inactive []int `json:"inactive"`
	// HashValue the hash // xyzzy
	HashValue []byte // xyzzy
}

// Bytes returns the result as a byte slice.
// TODO: How should we send it to the chain? Should it be sha256 hash, result
// serialized to json or something else?
func (r *Result) Bytes() []byte {
	return []byte(fmt.Sprintf("%v", r))
}
