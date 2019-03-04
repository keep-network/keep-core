package dkg2

import "github.com/keep-network/keep-core/pkg/beacon/relay/chain"

type dkgResultsVotes chain.DKGResultsVotes

// contains returns true if result hash is in the set of DKG results votes.
func (d dkgResultsVotes) contains(resultHash chain.DKGResultHash) bool {
	_, contains := d[resultHash]
	return contains
}

// leads returns hashes of DKG results with the highest number of registered votes.
// In case there are more than one result registered with the same highest number
// of votes the function will return multiple DKR result hashes. If the set of
// DKG results votes is empty it returns empty slice.
func (d dkgResultsVotes) leads() []chain.DKGResultHash {
	leadingResultsHashes := make([]chain.DKGResultHash, 0)

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
// It returns true if the given DKG result hash is the only result with the highest
// number of votes. It means that if the given result does not have the highest
// number of votes or there are two or more results on the leading position it
// will return false. If the set of DKG results votes is empty it returns false.
func (d dkgResultsVotes) isStrictlyLeading(resultHash chain.DKGResultHash) bool {
	leadingDKGResults := d.leads()

	if len(leadingDKGResults) == 1 {
		if leadingDKGResults[0] == resultHash {
			return true
		}
	}

	return false
}

// leadHasEnoughVotes checks if leading results' number of votes for is greater
// than threshold value.
func (d dkgResultsVotes) leadHasEnoughVotes(dishonestThreshold int) bool {
	if leads := d.leads(); len(leads) > 0 {
		return d[leads[0]] > dishonestThreshold
	}

	return false
}
