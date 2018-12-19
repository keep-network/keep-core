package chain

import "math/big"

// Submissions is the set of submissions for a requestID.
// Submissions tracks the number of votes for each unique
// submission and associates this set with a requestID.
type Submissions struct {
	requestID   *big.Int
	Submissions []*Submission
}

// Submission is an individual submission that counts the number of votes
type Submission struct {
	DKGResult *DKGResult
	Votes     int
}

// Lead returns a submission with the highest number of votes.  If there are
// no submissions it returns nil.
func (s *Submissions) Lead() *Submission {
	if len(s.Submissions) == 0 {
		return nil
	}
	topVote := -1
	topPos := 0
	for pos, submission := range s.Submissions {
		if topVote < submission.Votes {
			topPos = pos
			topVote = submission.Votes
		}
	}
	return s.Submissions[topPos]
}

// Contains returns true if 'result' is in the set of submissions.
func (s *Submissions) Contains(result *DKGResult) bool {
	for _, submission := range s.Submissions {
		if result.Equals(submission.DKGResult) {
			return true
		}
	}
	return false
}
