package maintainer

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
)

// RetargetEvent represents an invocation of the Retarget method.
type RetargetEvent struct {
	oldDifficulty, newDifficulty uint32
}

type localBitcoinDifficultyChain struct {
	currentEpoch uint64
	proofLength  uint64

	retargetEvents []*RetargetEvent
}

// Ready checks whether the relay is active (i.e. genesis has been
// performed).
func (lbdc *localBitcoinDifficultyChain) Ready() (bool, error) {
	panic("unsupported")
}

// IsAuthorizationRequired checks whether the relay requires the address
// submitting a retarget to be authorised in advance by governance.
func (lbdc *localBitcoinDifficultyChain) IsAuthorizationRequired() (bool, error) {
	panic("unsupported")
}

// IsAuthorized checks whether the given address has been authorised to
// submit a retarget by governance.
func (lbdc *localBitcoinDifficultyChain) IsAuthorized(address chain.Address) (bool, error) {
	panic("unsupported")
}

// Signing returns the signing associated with the chain.
func (lbdc *localBitcoinDifficultyChain) Signing() chain.Signing {
	panic("unsupported")
}

// Retarget adds a new epoch to the Bitcoin difficulty relay by providing
// a proof of the difficulty before and after the retarget.
func (lbdc *localBitcoinDifficultyChain) Retarget(headers []*bitcoin.BlockHeader) error {
	// For simplicity, store block header bits instead of their difficulty
	// targets.
	retargetEvent := &RetargetEvent{
		oldDifficulty: headers[len(headers)/2-1].Bits,
		newDifficulty: headers[len(headers)/2].Bits,
	}
	lbdc.retargetEvents = append(lbdc.retargetEvents, retargetEvent)

	return nil
}

// CurrentEpoch returns the number of the latest epoch whose difficulty is
// proven to the relay. If the genesis epoch's number is set correctly, and
// retargets along the way have been legitimate, the current epoch equals
// the height of the block starting the most recent epoch, divided by 2016.
func (lbdc *localBitcoinDifficultyChain) CurrentEpoch() (uint64, error) {
	return lbdc.currentEpoch, nil
}

// ProofLength returns the number of blocks required for each side of a
// retarget proof: a retarget must provide `proofLength` blocks before
// the retarget and `proofLength` blocks after it.
func (lbdc *localBitcoinDifficultyChain) ProofLength() (uint64, error) {
	return lbdc.proofLength, nil
}

// SetCurrentEpoch sets the current proven epoch in the chain.
func (lbdc *localBitcoinDifficultyChain) SetCurrentEpoch(currentEpoch uint64) {
	lbdc.currentEpoch = currentEpoch
}

// SetCurrentEpoch sets the proof length needed for a retarget.
func (lbdc *localBitcoinDifficultyChain) SetProofLength(proofLength uint64) {
	lbdc.proofLength = proofLength
}

// RetargetEvents returns all invocations of the Retarget method.
func (lbdc *localBitcoinDifficultyChain) RetargetEvents() []*RetargetEvent {
	return lbdc.retargetEvents
}
