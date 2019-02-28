package dkg2

import (
	"bytes"
	"reflect"
	"sort"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

func TestContains(t *testing.T) {
	dkgResult1Hash := chain.DKGResultHash{10}
	dkgResult2Hash := chain.DKGResultHash{20}
	dkgResult3Hash := chain.DKGResultHash{30}

	tests := map[string]struct {
		currentDKGResultsVotes dkgResultsVotes
		lookFor                chain.DKGResultHash
		expectedResult         bool
	}{
		"empty set of results votes": {
			currentDKGResultsVotes: dkgResultsVotes{},
			lookFor:                dkgResult1Hash,
			expectedResult:         false,
		},
		"only one result votes": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 1,
			},
			lookFor:        dkgResult1Hash,
			expectedResult: true,
		},
		"1st result is a match": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 1,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult1Hash,
			expectedResult: true,
		},
		"2nd result is a match": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 1,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult2Hash,
			expectedResult: true,
		},
		"result not found in current results votes": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 1,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult3Hash,
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := test.currentDKGResultsVotes.contains(test.lookFor)
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
	dkgResult1Hash := chain.DKGResultHash{10}
	dkgResult2Hash := chain.DKGResultHash{20}
	dkgResult3Hash := chain.DKGResultHash{30}
	dkgResult4Hash := chain.DKGResultHash{40}

	tests := map[string]struct {
		currentDKGResultsVotes dkgResultsVotes
		expectedLeads          []chain.DKGResultHash
	}{
		"empty set of results votes": {
			currentDKGResultsVotes: dkgResultsVotes{},
			expectedLeads:          []chain.DKGResultHash{},
		},
		"only one result vote in the set": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 1,
			},
			expectedLeads: []chain.DKGResultHash{
				dkgResult1Hash,
			},
		},
		"1st result hash has the highest number of votes": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 3,
				dkgResult2Hash: 2,
				dkgResult3Hash: 2,
				dkgResult4Hash: 1,
			},
			expectedLeads: []chain.DKGResultHash{
				dkgResult1Hash,
			},
		},
		"2nd result hash has the highest number of votes": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 1,
				dkgResult2Hash: 3,
				dkgResult3Hash: 2,
				dkgResult4Hash: 1,
			},
			expectedLeads: []chain.DKGResultHash{
				dkgResult2Hash,
			},
		},
		"1st and 3rd results hashes have the highest number of votes": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 3,
				dkgResult2Hash: 2,
				dkgResult3Hash: 3,
				dkgResult4Hash: 1,
			},
			expectedLeads: []chain.DKGResultHash{
				dkgResult1Hash,
				dkgResult3Hash,
			},
		},
		"2nd and 4th results hashes have the highest number of votes": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 2,
				dkgResult2Hash: 4,
				dkgResult3Hash: 1,
				dkgResult4Hash: 4,
			},
			expectedLeads: []chain.DKGResultHash{
				dkgResult2Hash,
				dkgResult4Hash,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := test.currentDKGResultsVotes.leads()

			sortDKGResultHashes(actualResult)
			sortDKGResultHashes(test.expectedLeads)

			if !reflect.DeepEqual(test.expectedLeads, actualResult) {
				t.Errorf(
					"\nexpected: %+v\nactual:   %+v",
					test.expectedLeads,
					actualResult,
				)
			}
		})
	}
}

func sortDKGResultHashes(slice []chain.DKGResultHash) {
	sort.SliceStable(
		slice,
		func(i, j int) bool { return (bytes.Compare(slice[i][:], slice[j][:]) < 0) },
	)
}

func TestIsStrictlyLeading(t *testing.T) {
	dkgResult1Hash := chain.DKGResultHash{10}
	dkgResult2Hash := chain.DKGResultHash{20}
	dkgResult3Hash := chain.DKGResultHash{30}

	tests := map[string]struct {
		currentDKGResultsVotes dkgResultsVotes
		lookFor                chain.DKGResultHash
		expectedResult         bool
	}{
		"empty set of results votes": {
			currentDKGResultsVotes: dkgResultsVotes{},
			lookFor:                dkgResult1Hash,
			expectedResult:         false,
		},
		"only one result votes": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 1,
			},
			lookFor:        dkgResult1Hash,
			expectedResult: true,
		},
		"two leading results": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 2,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult1Hash,
			expectedResult: false,
		},
		"result is not not leading": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 3,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult2Hash,
			expectedResult: false,
		},
		"result is strictly leading": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 1,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult2Hash,
			expectedResult: true,
		},
		"result is not registered": {
			currentDKGResultsVotes: dkgResultsVotes{
				dkgResult1Hash: 1,
				dkgResult2Hash: 2,
			},
			lookFor:        dkgResult3Hash,
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult := test.currentDKGResultsVotes.isStrictlyLeading(test.lookFor)
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
