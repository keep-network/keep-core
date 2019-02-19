package chain

import (
	"math/big"
	"reflect"
	"testing"
)

func TestLead(t *testing.T) {
	dkgResult1 := &DKGResult{
		Success:        true,
		GroupPublicKey: []byte{10, 1},
		Disqualified:   []byte{},
		Inactive:       []byte{},
	}
	dkgResult2 := &DKGResult{
		Success:        true,
		GroupPublicKey: []byte{10, 02},
		Disqualified:   []byte{},
		Inactive:       []byte{},
	}
	tests := map[string]struct {
		currentSubmissions *DKGSubmissions
		expectedResult     *DKGSubmission
	}{
		"empty set of submissions": {
			currentSubmissions: &DKGSubmissions{},
			expectedResult:     nil, // &DKGSubmission{},
		},
		"2nd submission has high votes": {
			currentSubmissions: &DKGSubmissions{
				requestID: big.NewInt(100),
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
			expectedResult: &DKGSubmission{
				DKGResult: dkgResult2,
				Votes:     2,
			},
		},
		"1st submission has high votes": {
			currentSubmissions: &DKGSubmissions{
				requestID: big.NewInt(100),
				DKGSubmissions: []*DKGSubmission{
					{
						DKGResult: dkgResult1,
						Votes:     3,
					},
					{
						DKGResult: dkgResult2,
						Votes:     2,
					},
				},
			},
			expectedResult: &DKGSubmission{
				DKGResult: dkgResult1,
				Votes:     3,
			},
		},
		"submission has same votes": {
			currentSubmissions: &DKGSubmissions{
				requestID: big.NewInt(100),
				DKGSubmissions: []*DKGSubmission{
					{
						DKGResult: dkgResult1,
						Votes:     1,
					},
					{
						DKGResult: dkgResult2,
						Votes:     1,
					},
				},
			},
			expectedResult: &DKGSubmission{
				DKGResult: dkgResult1,
				Votes:     1,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := test.currentSubmissions.Lead()
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
				requestID: big.NewInt(100),
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
				requestID: big.NewInt(100),
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
				requestID: big.NewInt(100),
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
