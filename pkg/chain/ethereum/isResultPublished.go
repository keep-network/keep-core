package ethereum

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

func (ec *ethereumChain) IsResultPublished(result *result.Result) bool {
	resultHash := result.Hash()

	// Placeholder FIXME.
	/*
		for _, r := range c.submittedResults {
			if reflect.DeepEqual(r, resultHash) {
				return true
			}
		}
	*/

	return false
}

func (ec *ethereumChain) SubmitResult(publisherID int, result *result.Result) *async.ResultPublishPromise {
	/*
		c.submittedResultsMutex.Lock()
		defer c.submittedResultsMutex.Unlock()
	*/

	resultPublishPromise := &async.ResultPublishPromise{}

	resultHash := result.Hash()

	// Placeholder FIXME.
	/*
		for _, r := range c.submittedResults {
			if reflect.DeepEqual(r, resultHash) {
				resultPublishPromise.Fail(fmt.Errorf("Result already submitted"))
				return resultPublishPromise
			}
		}
	*/

	resultPublishPromise.Fulfill(&event.PublishedResult{
		PublisherID: publisherID,
		Hash:        resultHash,
	})

	// Placeholder FIXME.
	/*
		c.handlerMutex.Lock()
		c.submittedResults = append(c.submittedResults, resultHash)
		c.handlerMutex.Unlock()
	*/

	return resultPublishPromise
}
