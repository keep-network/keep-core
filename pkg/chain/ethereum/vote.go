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

// PHASE 14 pt 2
func (c *localChain) MapRequestIDToGroupPubKey(requestID, groupPubKey *big.Int) error {
	// TODO -- implement  MapRequestIDToGroupPubKey(requestID, groupPubKey *big.Int) error {
	return nil
}

// PHASE 14 pt 2
func (c *localChain) GetGroupPubKeyForRequestID(requestID *big.Int) (*big.Int, error) {
	// TODO -- implement  GetGroupPubKeyForRequestID(requestID *big.Int) (*big.Int, error) {
	return nil, nil
}
