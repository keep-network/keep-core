package ethereum

import (
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/gen/async"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
)

// GetDKGSubmissions returns the set of submissions.
func (ec *ethereumChain) GetDKGSubmissions(requestID *big.Int) *relaychain.DKGSubmissions {
	// TODO -- implement GetDKGSubmissions
	panic("NOT IMPLEMENTED")
}

// DKGResultVote increments the vote for requestID
func (ec *ethereumChain) DKGResultVote(
	requestID *big.Int,
	dkgResultHash []byte,
) *async.DKGResultVotePromise {
	// TODO -- implement Vote
	dkgResultVotePromise := &async.DKGResultVotePromise{}
	dkgResultVotePromise.Fail(fmt.Errorf("function not implemented"))
	return dkgResultVotePromise
}

// OnDKGResultVote register a function for a callback when a vote occurs.
func (ec *ethereumChain) OnDKGResultVote(func(dkgResultVote *event.DKGResultVote)) {
	// TODO -- implement OnDKGResultVote
}
