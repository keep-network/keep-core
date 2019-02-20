package chain

import "math/big"

// DKGSubmissions is the set of submissions for a requestID. It tracks the number
// of votes for each unique submission and associates this set with a requestID.
type DKGSubmissions struct {
	requestID      *big.Int
	DKGSubmissions []*DKGSubmission
}

// DKGSubmission is an individual submission that counts the number of votes.
type DKGSubmission struct {
	DKGResult *DKGResult
	Votes     int
}

// Leads returns submissions with the highest number of votes. If there are
// no submissions it returns nil.
func (d *DKGSubmissions) Leads() []*DKGSubmission {
	leadingSubmissions := make([]*DKGSubmission, 0)

	if len(d.DKGSubmissions) == 0 {
		return leadingSubmissions
	}

	highestVotes := d.DKGSubmissions[0].Votes
	for _, submission := range d.DKGSubmissions {
		if submission.Votes > highestVotes {
			highestVotes = submission.Votes
		}
	}

	for _, submission := range d.DKGSubmissions {
		if submission.Votes == highestVotes {
			leadingSubmissions = append(leadingSubmissions, submission)
		}
	}

	return leadingSubmissions
}

// Contains returns true if 'result' is in the set of submissions.
func (d *DKGSubmissions) Contains(result *DKGResult) bool {
	for _, submission := range d.DKGSubmissions {
		if result.Equals(submission.DKGResult) {
			return true
		}
	}
	return false
}

// IsOnlyLead checks if given result is the only leading submission. If the
// submissions set is empty it returns false.
func (d *DKGSubmissions) IsOnlyLead(result *DKGResult) bool {
	leadingSubmissions := d.Leads()

	if len(leadingSubmissions) == 1 {
		if leadingSubmissions[0].DKGResult.Equals(result) {
			return true
		}
	}

	return false
}
