package gjkr

import (
	"fmt"
	"math/big"

	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

func (vc *DKG) ChallengeStateChange(
	EventName string,
	group *big.Int,
	currentResult *result.Result,
	resultHash []byte,
	correctResult int,
) {

	AddGrupIfNotExists := func(group *big.Int) {
		groupAsKey := fmt.Sprintf("%s", group)
		_, found := vc.PerGroup[groupAsKey]
		if found {
			return
		}
		newGroup := ValidationGroupState{
			TFirst:           vc.TNow,                  // T_first - Block height for the group - when the first result event occurred block height
			AllResults:       make([]result.Result, 0), // []result.Result // Set of all results
			AllVotes:         make([]ResultVotes, 0),   // []ResultVotes   // Set of all results
			LeadResult:       correctResult,            // Position of lead result
			AlreadySubmitted: false,                    //
		}
		vc.PerGroup[groupAsKey] = newGroup
	}

	RmGrup := func(group *big.Int) {
	}

	switch EventName {
	case "Challenge":
		// If group not in set of groups add it.
		vc.TNow = getCurrentBlockHeight() // has to have geth connect info
		AddGrupIfNotExists(group)

		groupAsKey := fmt.Sprintf("%s", group)
		m1 := vc.PerGroup[groupAsKey]
		// If 'resultHash' is new then Initialize - In this group, take
		// TFirst and stash it.  This is the staring block, save in teh set
		// of results.
		_, err := FindHash(resultHash, m1.AllResults)
		if err != nil {
			fmt.Printf("%sREACHED: 00000 Add result, %s%s\n", MiscLib.ColorYellow, godebug.LF(), MiscLib.ColorReset)
			m1.AllResults = append(m1.AllResults, *currentResult) // Should this merge the solutions together
			resultPkg := ResultVotes{
				Votes:       1,
				BlockHeight: vc.TNow,
				Group:       groupAsKey,
			}
			m1.AllVotes = append(m1.AllVotes, resultPkg) // Should this merge the solutions together
		}
		leadResult := FindMostVotes(m1.AllVotes) // Summarize the votes for a particular solution
		vc.PerGroup[groupAsKey] = m1

		if vc.TNow > (vc.TFirst + vc.TConflict) {
			fmt.Printf("%sREACHED: 00001 Criteria Met, %s%s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
			m1.CurrentState = 1
		} else if m1.AllVotes[leadResult].Votes > vc.TMax {
			fmt.Printf("%sREACHED: 00002 Criteria Met: Votes for lead larger than TMax, %s%s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
			m1.CurrentState = 2
		} else if correctResult == leadResult {
			fmt.Printf("%sREACHED: 00003 Criteria Met: Correct result is lead result, %s%s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
			m1.CurrentState = 3
		} else if m1.AlreadySubmitted {
			m1.CurrentState = 5
			fmt.Printf("%sREACHED: 00005 AllreadySubmitted is true. %s%s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
		} else if inResult(correctResult, m1.AllVotes, m1.AllResults) {
			m1.CurrentState = 7
			fmt.Printf("%sREACHED: 00007 Submitted Result!!!!!! Yea !!!!!!!. %s%s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
			SubmitRessult(SignResult(resultHash))
			m1.AlreadySubmitted = true
		} else {
			m1.CurrentState = 11
			fmt.Printf("%sREACHED: 00011 how do you get to this point???. %s%s\n", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
			m1.AlreadySubmitted = true
		}
		vc.PerGroup[groupAsKey] = m1

	case "Vote":
		vc.TNow = getCurrentBlockHeight()
		AddGrupIfNotExists(group)
		vc.VoteChangeState(resultHash, group)

	case "Block":
		// If too mahy blocks have passed then toss and close.

		// go thorugh all the pending stuff -
		//	For each one
		//	  for each group
		//		If ....

	case "CreateGroup":
		AddGrupIfNotExists(group)

	case "TerminateGroup":
		RmGrup(group)

	}
}

// Vote Adds 1 to vote
func (vc *DKG) VoteChangeState(resultHash []byte, groupPublicKey *big.Int) error {
	group := fmt.Sprintf("%s", groupPublicKey) // Public Key for Group - big number
	m1 := vc.PerGroup[group]
	// Find resultHash in set of all results or error
	pos, err := FindHash(resultHash, m1.AllResults)
	if err != nil {
		return fmt.Errorf("invalid hash not found [%v]", err)
	}
	m1.AllVotes[pos].Votes++
	return nil
}
