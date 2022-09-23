package sortition

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestConjunctionPolicy(t *testing.T) {
	var tests = map[string]struct {
		policies       []JoinPolicy
		expectedResult bool
	}{
		"empty policy": {
			policies:       []JoinPolicy{},
			expectedResult: true,
		},
		"one positive policy": {
			policies:       []JoinPolicy{&mockPolicy{true}},
			expectedResult: true,
		},
		"one negative policy": {
			policies:       []JoinPolicy{&mockPolicy{false}},
			expectedResult: false,
		},
		"two policies: both negative": {
			policies:       []JoinPolicy{&mockPolicy{false}, &mockPolicy{false}},
			expectedResult: false,
		},
		"two policies: both positive": {
			policies:       []JoinPolicy{&mockPolicy{true}, &mockPolicy{true}},
			expectedResult: true,
		},
		"two policies: positive and negative": {
			policies:       []JoinPolicy{&mockPolicy{true}, &mockPolicy{false}},
			expectedResult: false,
		},
		"two policies: negative and positive": {
			policies:       []JoinPolicy{&mockPolicy{false}, &mockPolicy{true}},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			policy := ConjunctionPolicy{test.policies}

			actualResult := policy.ShouldJoin()
			testutils.AssertBoolsEqual(
				t,
				"ShouldJoin() result",
				test.expectedResult,
				actualResult,
			)
		})
	}
}

type mockPolicy struct {
	shouldJoin bool
}

func (mp *mockPolicy) ShouldJoin() bool {
	return mp.shouldJoin
}
