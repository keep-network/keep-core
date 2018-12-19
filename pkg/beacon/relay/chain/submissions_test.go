package chain

import (
	"math/big"
	"testing"
)

// Test Lead()
func TestLead(t *testing.T) {
	tests := map[string]struct {
		currentSubmissions *Submissions
		expectedResult     *Submission
	}{
		"empty set of submissions": {
			currentSubmissions: &Submissions{},
			expectedResult:     nil, // &Submission{},
		},
		"2nd submission has high votes": {
			currentSubmissions: &Submissions{
				requestID: big.NewInt(100),
				Submissions: []*Submission{
					{
						DKGResult: &DKGResult{
							Success:        true,
							GroupPublicKey: big.NewInt(1001),
							Disqualified:   []bool{},
							Inactive:       []bool{},
						},
						Votes: 1,
					},
					{
						DKGResult: &DKGResult{
							Success:        true,
							GroupPublicKey: big.NewInt(1002),
							Disqualified:   []bool{},
							Inactive:       []bool{},
						},
						Votes: 2,
					},
				},
			},
			expectedResult: &Submission{
				DKGResult: &DKGResult{
					Success:        true,
					GroupPublicKey: big.NewInt(1002),
					Disqualified:   []bool{},
					Inactive:       []bool{},
				},
				Votes: 2,
			},
		},
		"1st submission has high votes": {
			currentSubmissions: &Submissions{
				requestID: big.NewInt(100),
				Submissions: []*Submission{
					{
						DKGResult: &DKGResult{
							Success:        true,
							GroupPublicKey: big.NewInt(1001),
							Disqualified:   []bool{},
							Inactive:       []bool{},
						},
						Votes: 3,
					},
					{
						DKGResult: &DKGResult{
							Success:        true,
							GroupPublicKey: big.NewInt(1002),
							Disqualified:   []bool{},
							Inactive:       []bool{},
						},
						Votes: 2,
					},
				},
			},
			expectedResult: &Submission{
				DKGResult: &DKGResult{
					Success:        true,
					GroupPublicKey: big.NewInt(1001),
					Disqualified:   []bool{},
					Inactive:       []bool{},
				},
				Votes: 3,
			},
		},
		"submission has same votes": {
			currentSubmissions: &Submissions{
				requestID: big.NewInt(100),
				Submissions: []*Submission{
					{
						DKGResult: &DKGResult{
							Success:        true,
							GroupPublicKey: big.NewInt(1001),
							Disqualified:   []bool{},
							Inactive:       []bool{},
						},
						Votes: 1,
					},
					{
						DKGResult: &DKGResult{
							Success:        true,
							GroupPublicKey: big.NewInt(1002),
							Disqualified:   []bool{},
							Inactive:       []bool{},
						},
						Votes: 1,
					},
				},
			},
			expectedResult: &Submission{
				DKGResult: &DKGResult{
					Success:        true,
					GroupPublicKey: big.NewInt(1001),
					Disqualified:   []bool{},
					Inactive:       []bool{},
				},
				Votes: 1,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := test.currentSubmissions.Lead()
			if test.expectedResult == nil {
				if actualResult != nil {
					t.Errorf(
						"\nexpected: [nil]\nactual:   [%+v]",
						test.expectedResult,
					)
				}
			} else {
				if !test.expectedResult.DKGResult.Equals(actualResult.DKGResult) {
					t.Errorf(
						"\nexpected: [%+v]\nactual:   [%+v]",
						actualResult,
						test.expectedResult,
					)
				}
			}
		})
	}
}

// Test Contains
func TestContains(t *testing.T) {
	tests := map[string]struct {
		currentSubmissions *Submissions
		// func (s *Submissions) Contains(result *DKGResult) bool {
		lookFor        *DKGResult
		expectedResult bool
	}{
		"empty set of submissions": {
			currentSubmissions: &Submissions{},
			lookFor: &DKGResult{
				Success:        true,
				GroupPublicKey: big.NewInt(1001),
				Disqualified:   []bool{},
				Inactive:       []bool{},
			},
			expectedResult: false,
		},
		"1st submission is a match": {
			currentSubmissions: &Submissions{
				requestID: big.NewInt(100),
				Submissions: []*Submission{
					{
						DKGResult: &DKGResult{
							Success:        true,
							GroupPublicKey: big.NewInt(1001),
							Disqualified:   []bool{},
							Inactive:       []bool{},
						},
						Votes: 1,
					},
					{
						DKGResult: &DKGResult{
							Success:        true,
							GroupPublicKey: big.NewInt(1002),
							Disqualified:   []bool{},
							Inactive:       []bool{},
						},
						Votes: 2,
					},
				},
			},
			lookFor: &DKGResult{
				Success:        true,
				GroupPublicKey: big.NewInt(1001),
				Disqualified:   []bool{},
				Inactive:       []bool{},
			},
			expectedResult: true,
		},
		"2nd submission is a match": {
			currentSubmissions: &Submissions{
				requestID: big.NewInt(100),
				Submissions: []*Submission{
					{
						DKGResult: &DKGResult{
							Success:        true,
							GroupPublicKey: big.NewInt(1001),
							Disqualified:   []bool{},
							Inactive:       []bool{},
						},
						Votes: 1,
					},
					{
						DKGResult: &DKGResult{
							Success:        true,
							GroupPublicKey: big.NewInt(1002),
							Disqualified:   []bool{},
							Inactive:       []bool{},
						},
						Votes: 2,
					},
				},
			},
			lookFor: &DKGResult{
				Success:        true,
				GroupPublicKey: big.NewInt(1002),
				Disqualified:   []bool{},
				Inactive:       []bool{},
			},
			expectedResult: true,
		},
		"not found - with current submissions": {
			currentSubmissions: &Submissions{
				requestID: big.NewInt(100),
				Submissions: []*Submission{
					{
						DKGResult: &DKGResult{
							Success:        true,
							GroupPublicKey: big.NewInt(1001),
							Disqualified:   []bool{},
							Inactive:       []bool{},
						},
						Votes: 1,
					},
					{
						DKGResult: &DKGResult{
							Success:        true,
							GroupPublicKey: big.NewInt(1002),
							Disqualified:   []bool{},
							Inactive:       []bool{},
						},
						Votes: 2,
					},
				},
			},
			lookFor: &DKGResult{
				Success:        true,
				GroupPublicKey: big.NewInt(1032),
				Disqualified:   []bool{},
				Inactive:       []bool{},
			},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := test.currentSubmissions.Contains(test.lookFor)
			if test.expectedResult != actualResult {
				t.Errorf(
					"\nexpected: [%+v]\nactual:   [%+v]",
					actualResult,
					test.expectedResult,
				)
			}
		})
	}
}
