package ethereum

import (
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
)

// PHASE 14
func (ec *ethereumChain) GetDKGSubmissions(requestID *big.Int) *relaychain.Submissions {
	// TODO -- implement GetDKGSubmissions
	return nil
}

// PHASE 14
func (ec *ethereumChain) Vote(requestID *big.Int, dkgResultHash []byte) {
	// TODO -- implement Vote
}

// PHASE 14
func (ec *ethereumChain) OnDKGResultVote(func(dkgResultVote *event.DKGResultVote)) {
	// TODO -- implement OnDKGResultVote
}
