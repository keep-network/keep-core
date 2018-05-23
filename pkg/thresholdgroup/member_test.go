package thresholdgroup

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/dfinity/go-dfinity-crypto/bls"
)

func TestMain(m *testing.M) {
	bls.Init(bls.CurveSNARK1)

	os.Exit(m.Run())
}

var (
	defaultID        = "12345"
	defaultThreshold = 4
	defaultGroupSize = 12
)

func TestLocalMemberCreation(t *testing.T) {
	id := fmt.Sprintf("%x", rand.Int31())

	member, err := NewMember(id, defaultThreshold, defaultGroupSize)
	if err != nil {
		t.Fatalf("unexpected error [%v]", err)
	}

	if member == nil {
		t.Fatal("expected: non-nil\nactual: nil")
	}

	var propertyTests = map[string]struct {
		propertyFunc func(*LocalMember) string
		expected     string
	}{
		"id": {
			propertyFunc: func(lm *LocalMember) string { return lm.ID },
			expected:     fmt.Sprintf("0x%010s", id),
		},
		"BLS ID": {
			propertyFunc: func(lm *LocalMember) string { return lm.BlsID.GetHexString() },
			expected:     id,
		},
		"threshold": {
			propertyFunc: func(lm *LocalMember) string { return fmt.Sprintf("%v", lm.threshold) },
			expected:     fmt.Sprintf("%v", defaultThreshold),
		},
		"group size": {
			propertyFunc: func(lm *LocalMember) string { return fmt.Sprintf("%v", lm.groupSize) },
			expected:     fmt.Sprintf("%v", defaultGroupSize),
		},
	}

	for testName, test := range propertyTests {
		t.Run(testName, func(t *testing.T) {
			property := test.propertyFunc(member)
			if test.expected != property {
				t.Errorf("\nexpected: %s\nactual:   %s", test.expected, property)
			}
		})
	}
}

func TestLocalMemberFailsForHighThreshold(t *testing.T) {
	member, err := NewMember(defaultID, defaultGroupSize/2, defaultGroupSize)
	if err == nil {
		t.Fatal("\nexpected error but got nil")
	}
	if member != nil {
		t.Fatalf("\nexpected nil member for errored instantiation\ngot [%v]", member)
	}
}

func TestLocalMemberCommitments(t *testing.T) {
	member, _ := NewMember(defaultID, defaultThreshold, defaultGroupSize)

	if len(member.Commitments()) != defaultThreshold {
		t.Errorf(
			"\nexpected: %v commitments\nactual:   %v commitments",
			defaultThreshold,
			len(member.Commitments()),
		)
	}

	// Smoke tests: all commitments are initialized, no two commitments in a row
	// are the same.
	uninitializedCommitment := bls.PublicKey{}
	previousCommitment := uninitializedCommitment
	for i, commitment := range member.Commitments() {
		if commitment.IsEqual(&uninitializedCommitment) {
			t.Fatalf(
				"at index %v\nexpected: initialized commitment\nactual:   uninitialized commitment",
				i,
			)
		} else if commitment.IsEqual(&previousCommitment) {
			t.Errorf(
				"at index %v\nexpected: different commitment\nactual:   same as previous commitment",
				i,
			)
		}

		previousCommitment = commitment
	}
}

func TestLocalMemberRegistration(t *testing.T) {
	member, _ := NewMember(defaultID, defaultThreshold, defaultGroupSize)
	completeMemberCount := defaultGroupSize - 1

	member.RegisterMemberID(&member.BlsID)
	for i := 0; i < completeMemberCount; i++ {
		if member.MemberListComplete() {
			t.Fatalf(
				"\nmember list complete after %v instead of %v members",
				i+1,
				completeMemberCount,
			)
		}

		id := bls.ID{}
		id.SetDecString(fmt.Sprintf("%v", i))
		member.RegisterMemberID(&id)
	}

	if !member.MemberListComplete() {
		t.Errorf(
			"\nexpected: member list complete after %v members\nactual:   member list incomplete",
			completeMemberCount,
		)
	}
}
