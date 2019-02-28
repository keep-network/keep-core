package dkg2

import "github.com/keep-network/keep-core/pkg/beacon/relay/chain"

type dkgResultsVotes chain.DKGResultsVotes
type dkgResultHash = chain.DKGResultHash

// contains returns true if result hash is in the set of DKG results votes.
func (d dkgResultsVotes) contains(resultHash dkgResultHash) bool {
	_, contains := d[resultHash]
	return contains
}

// leads returns hashes of DKG results with the highest number of registered votes.
// In case there are more than one result registered with the same highest number
// of votes the function will return multiple DKR result hashes. If the set of
// DKG results votes is empty it returns empty slice.
func (d dkgResultsVotes) leads() []dkgResultHash {
	leadingResultsHashes := make([]dkgResultHash, 0)

	if len(d) == 0 {
		return leadingResultsHashes
	}

	highestVotes := 0
	for _, votes := range d {
		if votes > highestVotes {
			highestVotes = votes
		}
	}

	for resultHash, votes := range d {
		if votes == highestVotes {
			leadingResultsHashes = append(leadingResultsHashes, resultHash)
		}
	}

	return leadingResultsHashes
}

// isStrictlyLeading checks if given result hash is the only leading DKG result.
// If the set of DKG results votes is empty it returns false.
func (d dkgResultsVotes) isStrictlyLeading(resultHash dkgResultHash) bool {
	leadingDKGResults := d.leads()

	if len(leadingDKGResults) == 1 {
		if leadingDKGResults[0] == resultHash {
			return true
		}
	}

	return false
}
