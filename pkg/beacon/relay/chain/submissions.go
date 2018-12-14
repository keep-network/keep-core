package chain

import "math/big"

// Submissions - PHASE 14
type Submissions struct {
	requestID   *big.Int
	Submissions []*Submission
}

// Submission - PHASE 14
type Submission struct {
	DKGResult *DKGResult
	Votes     int
}

// Lead returns a submission with the highest number of votes.
func (s *Submissions) Lead() *Submission {
	if len(s.Submissions) == 0 {
		return nil
	}
	top := -1
	topPos := 0
	for pos, aSubmission := range s.Submissions {
		if top < aSubmission.Votes {
			topPos = pos
			top = aSubmission.Votes
		}
	}
	return s.Submissions[topPos]
}

// Contains returns true if 'result' is in the set of submissions.
func (s *Submissions) Contains(result *DKGResult) bool {
	for _, aSubmission := range s.Submissions {
		if result.Equals(aSubmission.DKGResult) {
			return true
		}
	}
	return false
}
