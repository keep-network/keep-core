package chain

import "math/big"

// DKGSubmissions is the set of submissions for a requestID.
// DKGSubmissions tracks the number of votes for each unique
// submission and associates this set with a requestID.
type DKGSubmissions struct {
	requestID      *big.Int
	DKGSubmissions []*DKGSubmission
}

// DKGSubmission is an individual submission that counts the number of votes.
type DKGSubmission struct {
	DKGResult *DKGResult
	Votes     int
}

// Lead returns a submission with the highest number of votes.  If there are
// no submissions it returns nil.
func (s *DKGSubmissions) Lead() *DKGSubmission {
	if len(s.DKGSubmissions) == 0 {
		return nil
	}
	topSubmission := s.DKGSubmissions[0]
	for pos, submission := range s.DKGSubmissions {
		if topSubmission.Votes < submission.Votes {
			topSubmission = s.DKGSubmissions[pos]
		}
	}
	return topSubmission
}

// Contains returns true if 'result' is in the set of submissions.
func (s *DKGSubmissions) Contains(result *DKGResult) bool {
	for _, submission := range s.DKGSubmissions {
		if result.Equals(submission.DKGResult) {
			return true
		}
	}
	return false
}
