package chain

// DKGResultsVotes is a map of votes for each DKG Result.
type DKGResultsVotes map[DKGResultHash]int

// Contains returns true if result hash is in the set of DKG results votes.
func (d DKGResultsVotes) Contains(resultHash DKGResultHash) bool {
	for dkgResultHash := range d {
		if resultHash == dkgResultHash {
			return true
		}
	}
	return false
}

// Leads returns hashes of DKG results with the highest number of registered votes.
// If the set of DKG results votes is empty it returns nil.
func (d DKGResultsVotes) Leads() []DKGResultHash {
	leadingResultsHashes := make([]DKGResultHash, 0)

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

// IsOnlyLead checks if given result hash is the only leading DKG result. If the
// set of DKG results votes is empty it returns false.
func (d DKGResultsVotes) IsOnlyLead(resultHash DKGResultHash) bool {
	leadingDKGResults := d.Leads()

	if len(leadingDKGResults) == 1 {
		if leadingDKGResults[0] == resultHash {
			return true
		}
	}

	return false
}
