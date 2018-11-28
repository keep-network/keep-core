package gjkr

import "testing"

func TestFinedMostVotes(t *testing.T) {
	vc := &ChallengeState{}
	vgs := []ResultVotes{
		{Votes: 4},
		{Votes: 2},
		{Votes: 6},
		{Votes: 1},
	}
	n := vc.findMostVotes(vgs)
	if n != 2 {
		t.Errorf("\nexpected: %d\nactual:   %d\n", 2, n)
	}
}
