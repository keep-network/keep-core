package entry

import (
	"fmt"
	"keep-network/beacon/relay"
)

// Request represents a request for an entry in the threshold relay.
type Request struct {
}

type entryProcessingState int

const (
	waitingForRequest entryProcessingState = iota
	// Upon joining, we may go to any of the following states
	generatingSigShare
	verifyingSigShares
	submittingSig
)

// ServeRequests kicks off the relay request monitoring/response publishing loop.
func ServeRequests(currentState relay.NodeState) {
	processingState := waitingForRequest
	// FIXME Probably best passed in from outside.
	thinger := make(chan Request)
	for request := range thinger {
		if isNodeResponsible(currentState) {
			processingState = generatingSigShare
			fmt.Println(processingState)
			generateSigShare(request)
		}
		processingState = waitingForRequest
	}
}

func isNodeResponsible(currentState relay.NodeState) bool {
	return currentState.IsNextGroup()
}

func generateSigShare(request Request) {
	fmt.Println(request)
}
