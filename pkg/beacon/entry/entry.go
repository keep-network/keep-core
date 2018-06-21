package entry

import (
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/relay"
	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// Request represents a request for an entry in the threshold relay.
type Request struct {
	previousEntry relay.Entry
}

type partialEntry struct {
	myShare             signatureShare
	verifiedOtherShares []signatureShare
}

type entryProcessingState int

type signatureShare struct {
	// groupID is the id of the node whose share this is within the relay group.
	groupID uint16
	// shareBytes is the actual bytes of the signature share.
	shareBytes []byte
}

const (
	waitingForRequest entryProcessingState = iota
	// Upon joining, we may go to any of the following states
	generatingSigShare
	verifyingSigShares
	submittingSig
)

// ServeRequests kicks off the relay request monitoring/response publishing loop.
func ServeRequests(relayChain relaychain.Interface, currentState *relay.NodeState) {
	processingState := waitingForRequest
	// FIXME Probably best passed in from outside.
	requestChan := make(chan Request)
	// FIXME Best passed in from the outside; channel for broadcasting a
	//       generated share.
	broadcastShare := func(share signatureShare) error { return nil }
	// FIXME Best passed in from the outside; channel for receiving other group
	//       members' generated shares.
	groupShares := make(chan signatureShare)
	for request := range requestChan {
		if isNodeResponsible(currentState) {
			processingState = generatingSigShare

			fmt.Println(processingState)
			nextShare := generateSigShare(currentState, request)

			processingState = verifyingSigShares
			err := broadcastShare(nextShare)
			if err != nil {
				// FIXME Need to figure out what failing to broadcast a share
				// (after retries, which broadcastShare should encapsulate)
				// triggers. Should we try to notify the host somehow beyond
				// crashing?
				panic(fmt.Sprintf("Tried to broadcast share but failed: [%v].", err))
			}

			finalShares := verifyIncomingGroupShares(request, groupShares)

			processingState = submittingSig
			currentEntry := partialEntry{nextShare, finalShares}
			err = submitEntry(currentEntry, relayChain)
			if err != nil {
				// FIXME Failing to submit an entry should probably be okay but
				// log a diagnostic message, but that remains to be decided.
				panic(fmt.Sprintf("Tried to submit relay entry but failed: [%v].", err))
			}
		}
		processingState = waitingForRequest
	}
}

func isNodeResponsible(currentState *relay.NodeState) bool {
	return currentState.IsNextGroup()
}

func generateSigShare(currentState *relay.NodeState, request Request) signatureShare {
	mySigShare := blsSign(request.previousEntry.Value[:])

	return signatureShare{currentState.GroupID, mySigShare}
}

// groupThreshold is the number of valid signature shares we need in order to
// recover the correct signature for the group.
const groupThreshold = 5

func verifyIncomingGroupShares(request Request, groupShares chan signatureShare) []signatureShare {
	previousValue := request.previousEntry.Value[:]
	verifiedShares := make([]signatureShare, groupThreshold-1 /* we already have our share */)
	currentShare := 0
	for share := range groupShares {
		// FIXME This will require a bit more info about the group's setup:
		// > Anyone can verify that share σ_i is valid by checking that
		// > (g_2,u_i,H(M),σ_i) is a co-Diffie-Hellman tuple.
		if !blsVerifyShare(previousValue, share) {
			// FIXME Need to broadcast accusation, perhaps trigger group
			// dissolution, rather than panicking.
			panic("Got invalid share, bailing!")
		}

		verifiedShares[currentShare] = share

		currentShare++
		if currentShare >= groupThreshold {
			break
		}
	}

	return verifiedShares
}

func submitEntry(entry partialEntry, relayChain relaychain.Interface) error {
	allShares := make([]signatureShare, len(entry.verifiedOtherShares)+1 /* my share */)
	copy(allShares, entry.verifiedOtherShares)
	allShares[len(entry.verifiedOtherShares)] = entry.myShare
	finalSignature := blsFinalSignatureFromShares(allShares)

	finalEntry := relay.Entry{Value: finalSignature, Timestamp: time.Now()}
	fmt.Println(fmt.Sprintf("attempting to submit entry [%v]", finalEntry))

	// FIXME Also probably want to return more than just error (e.g., were we
	//       the accepted entry?).
	return relayChain.SubmitRelayEntryCandidate(finalEntry)
}

// FIXME Actually sign instead of doubling all bytes...
func blsSign(previousValue []byte) []byte {
	newValue := make([]byte, len(previousValue))
	for i, bte := range previousValue {
		newValue[i] = bte * 2
	}

	return newValue
}

// FIXME Actually verify instead of always failing.
func blsVerifyShare(previousValue []byte, share signatureShare) bool {
	return false
}

// FIXME Actually build final signature instead of concatenating all the shares.
func blsFinalSignatureFromShares(shares []signatureShare) [8]byte {
	fullSignature := make([]byte, 0)
	for _, share := range shares {
		fullSignature = append(fullSignature, share.shareBytes...)
	}

	// truncate to fit into return size
	rightLengthSignature := [8]byte{}
	copy(rightLengthSignature[:], fullSignature)
	return rightLengthSignature
}
