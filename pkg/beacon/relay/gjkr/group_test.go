package gjkr

import (
	"testing"
)

func TestDisqualifyMemberID(t *testing.T) {
	group := &Group{}

	id1 := MemberID(123)
	id2 := MemberID(321)

	if len(group.disqualifiedMemberIDs) != 0 {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			0,
			len(group.disqualifiedMemberIDs),
		)
	}

	// Disqualify a member.
	group.DisqualifyMemberID(id1)
	if len(group.disqualifiedMemberIDs) != 1 {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			1,
			len(group.disqualifiedMemberIDs),
		)
	}

	// Disqualify the same member for a second time.
	group.DisqualifyMemberID(id1)
	if len(group.disqualifiedMemberIDs) != 1 {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			1,
			len(group.disqualifiedMemberIDs),
		)
	}

	// Disqualify a next member.
	group.DisqualifyMemberID(id2)
	if len(group.disqualifiedMemberIDs) != 2 {
		t.Fatalf("\nexpected: %v\nactual:   %v\n",
			2,
			len(group.disqualifiedMemberIDs),
		)
	}
}
