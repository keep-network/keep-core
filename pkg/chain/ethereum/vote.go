package ethereum

import (
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
)

// GetDKGSubmissions returns the set of submissions.
func (ec *ethereumChain) GetDKGSubmissions(requestID *big.Int) *relaychain.DKGSubmissions {
	// TODO -- implement GetDKGSubmissions
	panic("NOT IMPLEMENTED")
	return nil
}

// DKGResultVote increments the vote for requestID
func (ec *ethereumChain) DKGResultVote(requestID *big.Int, dkgResultHash []byte) {
	// TODO -- implement Vote
	panic("NOT IMPLEMENTED")
}

// OnDKGResultVote register a function for a callback when a vote occurs.
func (ec *ethereumChain) OnDKGResultVote(func(dkgResultVote *event.DKGResultVote)) {
	// TODO -- implement OnDKGResultVote
	panic("NOT IMPLEMENTED")
}
