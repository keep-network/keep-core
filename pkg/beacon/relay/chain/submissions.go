package chain

import "math/big"

// Submissions - PHASE 14
type Submissions struct {
	requestID   *big.Int
	submissions []*Submission
}

// Submission - PHASE 14
type Submission struct {
	DKGResult *DKGResult
	Votes     int
}

// Lead returns a submission with the highest number of votes.
func (s *Submissions) Lead() *Submission {
	// TODO Implement
	return s.submissions[0]
}

func (s *Submissions) Contains(result *DKGResult) bool {
	// TODO Implement
	return false
}
