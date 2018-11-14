package ethereum

import "github.com/keep-network/keep-core/pkg/beacon/relay/result"

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
