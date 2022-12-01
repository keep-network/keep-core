package maintainer

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"github.com/keep-network/keep-core/pkg/operator"
)

// RetargetEvent represents an invocation of the Retarget method.
type RetargetEvent struct {
	oldDifficulty, newDifficulty uint32
}

type localBitcoinDifficultyChain struct {
	operatorPrivateKey *operator.PrivateKey

	currentEpoch uint64
	proofLength  uint64

	ready                 bool
	authorizationRequired bool
	authorizedOperators   map[chain.Address]bool

	retargetEvents []*RetargetEvent
}

// Ready checks whether the relay is active (i.e. genesis has been performed).
func (lbdc *localBitcoinDifficultyChain) Ready() (bool, error) {
	return lbdc.ready, nil
}

// AuthorizationRequired checks whether the relay requires the address
// submitting a retarget to be authorised in advance by governance.
func (lbdc *localBitcoinDifficultyChain) AuthorizationRequired() (bool, error) {
	return lbdc.authorizationRequired, nil
}

// IsAuthorized checks whether the given address has been authorised by
// governance to submit a retarget.
func (lbdc *localBitcoinDifficultyChain) IsAuthorized(
	address chain.Address,
) (bool, error) {
	return lbdc.authorizedOperators[address], nil
}

// Signing returns the signing associated with the chain.
func (lbdc *localBitcoinDifficultyChain) Signing() chain.Signing {
	return local_v1.NewSigner(lbdc.operatorPrivateKey)
}

// Retarget adds a new epoch to the relay by providing a proof
// of the difficulty before and after the retarget.
func (lbdc *localBitcoinDifficultyChain) Retarget(
	headers []*bitcoin.BlockHeader,
) error {
	// For simplicity, store block header bits instead of their difficulty
	// targets.
	retargetEvent := &RetargetEvent{
		oldDifficulty: headers[len(headers)/2-1].Bits,
		newDifficulty: headers[len(headers)/2].Bits,
	}
	lbdc.retargetEvents = append(lbdc.retargetEvents, retargetEvent)

	lbdc.currentEpoch++

	return nil
}

// CurrentEpoch returns the number of the latest epoch whose difficulty is
// proven to the relay.
func (lbdc *localBitcoinDifficultyChain) CurrentEpoch() (uint64, error) {
	return lbdc.currentEpoch, nil
}

// ProofLength returns the number of blocks required for each side of a
// retarget proof.
func (lbdc *localBitcoinDifficultyChain) ProofLength() (uint64, error) {
	return lbdc.proofLength, nil
}

// SetReady sets chain's status as either ready or not.
func (lbdc *localBitcoinDifficultyChain) SetReady(ready bool) {
	lbdc.ready = ready
}

// SetAuthorizationRequired sets chain's authorization requirement to true
// or false.
func (lbdc *localBitcoinDifficultyChain) SetAuthorizationRequired(required bool) {
	lbdc.authorizationRequired = required
}

// SetAuthorizedOperator sets the given operator address as either authorized or
// unauthorized.
func (lbdc *localBitcoinDifficultyChain) SetAuthorizedOperator(
	operatorAddress chain.Address,
	authorized bool,
) {
	lbdc.authorizedOperators[operatorAddress] = authorized
}

// SetCurrentEpoch sets the current proven epoch in the chain.
func (lbdc *localBitcoinDifficultyChain) SetCurrentEpoch(currentEpoch uint64) {
	lbdc.currentEpoch = currentEpoch
}

// SetProofLength sets the proof length needed for a retarget.
func (lbdc *localBitcoinDifficultyChain) SetProofLength(proofLength uint64) {
	lbdc.proofLength = proofLength
}

// RetargetEvents returns all invocations of the Retarget method.
func (lbdc *localBitcoinDifficultyChain) RetargetEvents() []*RetargetEvent {
	return lbdc.retargetEvents
}

// ConnectLocal connects to the local Bitcoin difficulty chain and returns
// a chain handle.
func ConnectLocal() *localBitcoinDifficultyChain {
	operatorPrivateKey, _, err := operator.GenerateKeyPair(local_v1.DefaultCurve)
	if err != nil {
		panic(err)
	}

	return &localBitcoinDifficultyChain{
		operatorPrivateKey:  operatorPrivateKey,
		authorizedOperators: make(map[chain.Address]bool),
	}
}
