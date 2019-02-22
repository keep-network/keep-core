package chain

import (
	"reflect"
	"testing"
)

func TestLeads(t *testing.T) {
	dkgResult1 := &DKGResult{
		GroupPublicKey: []byte{10, 1},
	}
	dkgResult2 := &DKGResult{
		GroupPublicKey: []byte{10, 02},
	}
	dkgResult3 := &DKGResult{
		GroupPublicKey: []byte{10, 03},
	}
	dkgResult4 := &DKGResult{
		GroupPublicKey: []byte{10, 04},
	}

	tests := map[string]struct {
		currentSubmissions *DKGSubmissions
		expectedResult     []*DKGSubmission
	}{
		"empty set of submissions": {
			currentSubmissions: &DKGSubmissions{},
			expectedResult:     []*DKGSubmission{},
		},
		"only one submission in the set": {
			currentSubmissions: &DKGSubmissions{
				DKGSubmissions: []*DKGSubmission{
					{
						DKGResult: dkgResult1,
						Votes:     1,
					},
				},
			},
			expectedResult: []*DKGSubmission{
				&DKGSubmission{
					DKGResult: dkgResult1,
					Votes:     1,
				},
			},
		},
		"1st submission has high votes": {
			currentSubmissions: &DKGSubmissions{
				DKGSubmissions: []*DKGSubmission{
					{
						DKGResult: dkgResult1,
						Votes:     3,
					},
					{
						DKGResult: dkgResult2,
						Votes:     2,
					},
					{
						DKGResult: dkgResult3,
						Votes:     2,
					},
					{
						DKGResult: dkgResult4,
						Votes:     1,
					},
				},
			},
			expectedResult: []*DKGSubmission{
				&DKGSubmission{
					DKGResult: dkgResult1,
					Votes:     3,
				},
			},
		},
		"2nd submission has highest votes": {
			currentSubmissions: &DKGSubmissions{
				DKGSubmissions: []*DKGSubmission{
					{
						DKGResult: dkgResult1,
						Votes:     1,
					},
					{
						DKGResult: dkgResult2,
						Votes:     3,
					},
					{
						DKGResult: dkgResult3,
						Votes:     2,
					},
					{
						DKGResult: dkgResult4,
						Votes:     1,
					},
				},
			},
			expectedResult: []*DKGSubmission{
				&DKGSubmission{
					DKGResult: dkgResult2,
					Votes:     3,
				},
			},
		},
		"1st and 3rd submissions has highest votes": {
			currentSubmissions: &DKGSubmissions{
				DKGSubmissions: []*DKGSubmission{
					{
						DKGResult: dkgResult1,
						Votes:     3,
					},
					{
						DKGResult: dkgResult2,
						Votes:     2,
					},
					{
						DKGResult: dkgResult3,
						Votes:     3,
					},
					{
						DKGResult: dkgResult4,
						Votes:     1,
					},
				},
			},
			expectedResult: []*DKGSubmission{
				&DKGSubmission{
					DKGResult: dkgResult1,
					Votes:     3,
				},
				&DKGSubmission{
					DKGResult: dkgResult3,
					Votes:     3,
				},
			},
		},
		"2nd and 4th submissions has highest votes": {
			currentSubmissions: &DKGSubmissions{
				DKGSubmissions: []*DKGSubmission{
					{
						DKGResult: dkgResult1,
						Votes:     2,
					},
					{
						DKGResult: dkgResult2,
						Votes:     4,
					},
					{
						DKGResult: dkgResult3,
						Votes:     1,
					},
					{
						DKGResult: dkgResult4,
						Votes:     4,
					},
				},
			},
			expectedResult: []*DKGSubmission{
				&DKGSubmission{
					DKGResult: dkgResult2,
					Votes:     4,
				},
				&DKGSubmission{
					DKGResult: dkgResult4,
					Votes:     4,
				},
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := test.currentSubmissions.Leads()
			if !reflect.DeepEqual(test.expectedResult, actualResult) {
				t.Errorf(
					"\nexpected: %+v\nactual:   %+v",
					test.expectedResult,
					actualResult,
				)
			}
		})
	}
}

func TestContains(t *testing.T) {
	dkgResult1 := &DKGResult{
		Success:        true,
		GroupPublicKey: []byte{10, 01},
		Disqualified:   []byte{},
		Inactive:       []byte{},
	}
	dkgResult2 := &DKGResult{
		Success:        true,
		GroupPublicKey: []byte{10, 02},
		Disqualified:   []byte{},
		Inactive:       []byte{},
	}
	dkgResult3 := &DKGResult{
		Success:        true,
		GroupPublicKey: []byte{10, 32},
		Disqualified:   []byte{},
		Inactive:       []byte{},
	}

	tests := map[string]struct {
		currentSubmissions *DKGSubmissions
		lookFor            *DKGResult
		expectedResult     bool
	}{
		"empty set of submissions": {
			currentSubmissions: &DKGSubmissions{},
			lookFor:            dkgResult1,
			expectedResult:     false,
		},
		"1st submission is a match": {
			currentSubmissions: &DKGSubmissions{
				DKGSubmissions: []*DKGSubmission{
					{
						DKGResult: dkgResult1,
						Votes:     1,
					},
					{
						DKGResult: dkgResult2,
						Votes:     2,
					},
				},
			},
			lookFor:        dkgResult1,
			expectedResult: true,
		},
		"2nd submission is a match": {
			currentSubmissions: &DKGSubmissions{
				DKGSubmissions: []*DKGSubmission{
					{
						DKGResult: dkgResult1,
						Votes:     1,
					},
					{
						DKGResult: dkgResult2,
						Votes:     2,
					},
				},
			},
			lookFor:        dkgResult2,
			expectedResult: true,
		},
		"not found - with current submissions": {
			currentSubmissions: &DKGSubmissions{
				DKGSubmissions: []*DKGSubmission{
					{
						DKGResult: dkgResult1,
						Votes:     1,
					},
					{
						DKGResult: dkgResult2,
						Votes:     2,
					},
				},
			},
			lookFor:        dkgResult3,
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := test.currentSubmissions.Contains(test.lookFor)
			if test.expectedResult != actualResult {
				t.Errorf(
					"\nexpected: %v\nactual:   %v",
					test.expectedResult,
					actualResult,
				)
			}
		})
	}
}
