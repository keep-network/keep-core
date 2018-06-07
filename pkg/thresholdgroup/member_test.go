package thresholdgroup

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMemberReturnsFormattedID(t *testing.T) {
	member := buildSharingMember("face")

	expected := "0x000000face"
	if member.MemberID() != expected {
		t.Errorf(
			"\nexpected: %s\nactual:   %s",
			expected,
			member.MemberID(),
		)
	}
}

func TestMemberFailsOnInvalidID(t *testing.T) {
	member, err := NewMember("z", 5, 12)
	expected := fmt.Errorf("err mclBnFr_setStr -1")
	if !reflect.DeepEqual(err, expected) {
		t.Errorf(
			"\nexpected: %v\nactual:   %v",
			expected,
			err,
		)
	}
	if member != nil {
		t.Errorf("\nexpected: nil member\nactual:   %v", member)
	}
}

func TestMemberFailsOnInvalidThreshold(t *testing.T) {
	member, err := NewMember("z", 12, 12)
	expected := fmt.Errorf("threshold 12 >= 12 / 2, so group security cannot be guaranteed")
	if !reflect.DeepEqual(err, expected) {
		t.Errorf(
			"\nexpected: %v\nactual:   %v",
			expected,
			err,
		)
	}
	if member != nil {
		t.Errorf("\nexpected: nil member\nactual:   %v", member)
	}
}
