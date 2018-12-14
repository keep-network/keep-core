package ethereum

import (
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
)

// GetDKGSubmissions returns the set of submissions.
func (ec *ethereumChain) GetDKGSubmissions(requestID *big.Int) *relaychain.Submissions {
	// TODO -- implement GetDKGSubmissions
	return nil
}

// Vote increments the vote for requestID
func (ec *ethereumChain) Vote(requestID *big.Int, dkgResultHash []byte) {
	// TODO -- implement Vote
}

// OnDKGResultVote register a function for a callback when a vote occures.
func (ec *ethereumChain) OnDKGResultVote(func(dkgResultVote *event.DKGResultVote)) {
	// TODO -- implement OnDKGResultVote
}
