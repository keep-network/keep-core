package chain

import (
	"bytes"
	"reflect"
	"sort"
	"testing"
)

func TestContains(t *testing.T) {
	dkgResult1 := &DKGResult{GroupPublicKey: []byte{10, 01}}
	dkgResult2 := &DKGResult{GroupPublicKey: []byte{10, 02}}
	dkgResult3 := &DKGResult{GroupPublicKey: []byte{10, 32}}
	dkgResult1Hash, _ := dkgResult1.Hash()
	dkgResult2Hash, _ := dkgResult2.Hash()
	dkgResult3Hash, _ := dkgResult3.Hash()

	tests := map[string]struct {
		currentDKGResultsVotes DKGResultsVotes
		lookFor                DKGResultHash
		expectedResult         bool
	}{
		"empty set of results votes": {
			currentDKGResultsVotes: DKGResultsVotes{},
			lookFor:                dkgResult1Hash,
			expectedResult:         false,
		},
		"only one result votes": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 1,
			},
			lookFor:        dkgResult1Hash,
			expectedResult: true,
		},
		"1st result is a match": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 1,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult1Hash,
			expectedResult: true,
		},
		"2nd result is a match": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 1,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult2Hash,
			expectedResult: true,
		},
		"result not found in current results votes": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 1,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult3Hash,
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := test.currentDKGResultsVotes.Contains(test.lookFor)
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
	dkgResult1Hash, _ := dkgResult1.Hash()
	dkgResult2Hash, _ := dkgResult2.Hash()
	dkgResult3Hash, _ := dkgResult3.Hash()
	dkgResult4Hash, _ := dkgResult4.Hash()

	tests := map[string]struct {
		currentDKGResultsVotes DKGResultsVotes
		expectedResultHash     []DKGResultHash
	}{
		"empty set of results votes": {
			currentDKGResultsVotes: DKGResultsVotes{},
			expectedResultHash:     []DKGResultHash{},
		},
		"only one result vote in the set": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 1,
			},
			expectedResultHash: []DKGResultHash{
				dkgResult1Hash,
			},
		},
		"1st result hash has highest votes": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 3,
				dkgResult2Hash: 2,
				dkgResult3Hash: 2,
				dkgResult4Hash: 1,
			},
			expectedResultHash: []DKGResultHash{
				dkgResult1Hash,
			},
		},
		"2nd result hash has highest votes": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 1,
				dkgResult2Hash: 3,
				dkgResult3Hash: 2,
				dkgResult4Hash: 1,
			},
			expectedResultHash: []DKGResultHash{
				dkgResult2Hash,
			},
		},
		"1st and 3rd results hashes has highest votes": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 3,
				dkgResult2Hash: 2,
				dkgResult3Hash: 3,
				dkgResult4Hash: 1,
			},
			expectedResultHash: []DKGResultHash{
				dkgResult1Hash,
				dkgResult3Hash,
			},
		},
		"2nd and 4th results hashes has highest votes": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 2,
				dkgResult2Hash: 4,
				dkgResult3Hash: 1,
				dkgResult4Hash: 4,
			},
			expectedResultHash: []DKGResultHash{
				dkgResult2Hash,
				dkgResult4Hash,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := test.currentDKGResultsVotes.Leads()

			sortDKGResultHashes(actualResult)
			sortDKGResultHashes(test.expectedResultHash)

			if !reflect.DeepEqual(test.expectedResultHash, actualResult) {
				t.Errorf(
					"\nexpected: %+v\nactual:   %+v",
					test.expectedResultHash,
					actualResult,
				)
			}
		})
	}
}

func sortDKGResultHashes(slice []DKGResultHash) {
	sort.SliceStable(
		slice,
		func(i, j int) bool { return (bytes.Compare(slice[i][:], slice[j][:]) < 0) },
	)
}

func TestIsOnlyLead(t *testing.T) {
	dkgResult1 := &DKGResult{GroupPublicKey: []byte{10, 01}}
	dkgResult2 := &DKGResult{GroupPublicKey: []byte{10, 02}}
	dkgResult3 := &DKGResult{GroupPublicKey: []byte{10, 32}}
	dkgResult1Hash, _ := dkgResult1.Hash()
	dkgResult2Hash, _ := dkgResult2.Hash()
	dkgResult3Hash, _ := dkgResult3.Hash()

	tests := map[string]struct {
		currentDKGResultsVotes DKGResultsVotes
		lookFor                DKGResultHash
		expectedResult         bool
	}{
		"empty set of results votes": {
			currentDKGResultsVotes: DKGResultsVotes{},
			lookFor:                dkgResult1Hash,
			expectedResult:         false,
		},
		"only one result votes": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 1,
			},
			lookFor:        dkgResult1Hash,
			expectedResult: true,
		},
		"two leading results": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 2,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult1Hash,
			expectedResult: false,
		},
		"result is not not leading": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 3,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult2Hash,
			expectedResult: false,
		},
		"result is strictly leading": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 1,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult2Hash,
			expectedResult: true,
		},
		"result is not registered": {
			currentDKGResultsVotes: DKGResultsVotes{
				dkgResult1Hash: 1,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult3Hash,
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := test.currentDKGResultsVotes.IsOnlyLead(test.lookFor)
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
