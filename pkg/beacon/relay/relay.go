package relay

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/beacon/entry"
)

type entryProcessingState int

const (
	waitingForRequest entryProcessingState = iota
	// Upon joining, we may go to any of the following states
	generatingSigShare
	verifyingSigShares
	submittingSig
)

type partialEntry struct {
	myShare             signatureShare
	verifiedOtherShares []signatureShare
}

type signatureShare struct {
	// groupID is the id of the node whose share this is within the relay group.
	groupID uint16
	// shareBytes is the actual bytes of the signature share.
	shareBytes []byte
}

// NodeState represents the current state of a relay node.
type NodeState struct {
	// groupCount is the total number of groups in the relay.
	groupCount uint32
	// group is the id of the relay group this node belongs to. 0 if none.
	// Necessarily less than groupCount.
	group uint32
	// groupId is the id of this node within its relay group. 0 if none.
	GroupID uint16
	// lastSeenEntry is the last relay entry this node is aware of.
	lastSeenEntry entry.Entry
}

// IsNextGroup returns true if the next group expected to generate a threshold
// signature is the same as the group the NodeState belongs to.
func (state NodeState) IsNextGroup() bool {
	return binary.BigEndian.Uint32(state.lastSeenEntry.Value[:])%state.groupCount == state.group
}

// EmptyState returns an empty NodeState with no group, zero group count, and
// a nil last seen entry.
func EmptyState() NodeState {
	return NodeState{groupCount: 0, group: 0, GroupID: 0, lastSeenEntry: entry.Entry{Value: [8]byte{}, Timestamp: time.Unix(0, 0)}}
}

// ServeRequests kicks off the relay request monitoring/response publishing loop.
func ServeRequests(currentState NodeState) {
	processingState := waitingForRequest
	// FIXME Probably best passed in from outside.
	thinger := make(chan entry.Request)
	// FIXME Best passed in from the outside; channel for broadcasting a
	//       generated share.
	broadcastShare := func(share signatureShare) error { return nil }
	// FIXME Best passed in from the outside; channel for receiving other group
	//       members' generated shares.
	groupShares := make(chan signatureShare)
	for request := range thinger {
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
			err = submitEntry(currentEntry)
			if err != nil {
				// FIXME Failing to submit an entry should probably be okay but
				// log a diagnostic message, but that remains to be decided.
				panic(fmt.Sprintf("Tried to submit relay entry but failed: [%v].", err))
			}
		}
		processingState = waitingForRequest
	}
}

func isNodeResponsible(currentState NodeState) bool {
	return currentState.IsNextGroup()
}

func generateSigShare(currentState NodeState, request entry.Request) signatureShare {
	mySigShare := blsSign(request.PreviousEntry.Value[:])

	return signatureShare{currentState.GroupID, mySigShare}
}

// groupThreshold is the number of valid signature shares we need in order to
// recover the correct signature for the group.
const groupThreshold = 5

func verifyIncomingGroupShares(request entry.Request, groupShares chan signatureShare) []signatureShare {
	previousValue := request.PreviousEntry.Value[:]
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

func submitEntry(partialEntry partialEntry) error {
	allShares := make([]signatureShare, len(partialEntry.verifiedOtherShares)+1 /* my share */)
	copy(allShares, partialEntry.verifiedOtherShares)
	allShares[len(partialEntry.verifiedOtherShares)] = partialEntry.myShare
	finalSignature := blsFinalSignatureFromShares(allShares)

	finalEntry := entry.Entry{Value: finalSignature, Timestamp: time.Now()}
	fmt.Println(fmt.Sprintf("fake-submitting entry [%v]", finalEntry))

	// FIXME Magically submit error to the chain.
	// FIXME Also probably want to return more than just error (e.g., were we
	//       the accepted entry?).
	return nil
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
