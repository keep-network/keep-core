package ethereum

import (
	"math/big"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
)

func (ec *ethereumChain) GetDKGSubmissions(requestID *big.Int) *relaychain.Submissions {
	// TODO -- implement GetDKGSubmissions
	return nil
}

func (ec *ethereumChain) Vote(requestID *big.Int, dkgResultHash []byte) {
	// TODO -- implement Vote
}

func (ec *ethereumChain) OnDKGResultVote(func(dkgResultVote *event.DKGResultVote)) {
	// TODO -- implement OnDKGResultVote
}

func (c *localChain) MapRequestIDToGroupPubKey(requestID, groupPubKey *big.Int) error {
	// TODO -- implement  MapRequestIDToGroupPubKey(requestID, groupPubKey *big.Int) error {
	return nil
}

func (c *localChain) GetGroupPubKeyForRequestID(requestID *big.Int) (*big.Int, error) {
	// TODO -- implement  GetGroupPubKeyForRequestID(requestID *big.Int) (*big.Int, error) {
	return nil, nil
}
