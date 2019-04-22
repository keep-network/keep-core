package relay

import (
	"fmt"
	"math/big"
	"testing"
)

func TestNextGroupIndex(t *testing.T) {
	var tests = map[string]struct {
		previousEntry  int
		numberOfGroups int
		expectedIndex  int
	}{
		"zero groups": {
			previousEntry:  12,
			numberOfGroups: 0,
			expectedIndex:  0,
		},
		"fewer groups than the previous entry value": {
			previousEntry:  13,
			numberOfGroups: 4,
			expectedIndex:  1,
		},
		"more groups than the previous entry value": {
			previousEntry:  3,
			numberOfGroups: 12,
			expectedIndex:  3,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			bigPreviousEntry := big.NewInt(int64(test.previousEntry))
			bigNumberOfGroups := big.NewInt(int64(test.numberOfGroups))
			expectedInex := big.NewInt(int64(test.expectedIndex))

			bigIndex := nextGroupIndex(bigPreviousEntry, bigNumberOfGroups)

			fmt.Printf("%v ≟ %v: %v\n", expectedInex, bigIndex, test)

			if bigIndex.Cmp(expectedInex) != 0 {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]\n",
					expectedInex,
					bigIndex,
				)
			}
		})
	}
}
