package gjkr

// Phase 14 code - PJS

import (
	"bytes"
	"fmt"
	"math/big"
	"sync"

	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
	"github.com/keep-network/keep-core/pkg/chain"
)

type ValidationGroupState struct {
	TFirst           int             // T_first - Block height for the group - when the first result event occurred block height
	AllResults       []result.Result // Set of all results
	AllVotes         []ResultVotes   // Set of all results
	LeadResult       int             // Position of lead result
	AlreadySubmitted bool            //
}

type ResultVotes struct {
	Votes       int    // How many votes this has
	BlockHeight int    // Block height of this results
	Group       string // the PubKey of the group
}

type ControlType struct {
	Cmd            string
	Hash           []byte
	GroupPublicKey *big.Int
}

type CompleteType struct {
	FinishState string
}

var Submission chan result.Result
var Control chan ControlType
var Complete chan CompleteType

func init() {
	Submission = make(chan result.Result)
	Control = make(chan ControlType)
	Complete = make(chan CompleteType, 30)
}

/*
while resultPublished and not finished:
  allResults = getSubmissions()
  leadResult = allResults.mostVotes

  T_now = getCurrentBlockHeight()
  T_first = allResults.earliest.submitTime

  if T_now > T_first + T_conflict or leadResult.votes > M_max:
    finished = True

  elif correctResult = leadResult or alreadySubmitted:
    wait()

  elif correctResult in allResults:
    submit(sign(resultHash))
    alreadySubmitted = True

  else:
    submit(correctResult)
    alreadySubmitted = True

*/

// Vote Adds 1 to vote
func (vc *DKG) Vote(resultHash []byte, groupPublicKey *big.Int) error {
	group := fmt.Sprintf("%s", groupPublicKey) // Public Key for Group - big number
	m1 := vc.PerGroup[group]
	// Find resultHash in set of all results or error
	pos, err := FindHash(resultHash, m1.AllResults)
	if err != nil {
		return fmt.Errorf("invalid hash not found [%v]", err)
	}
	m1.AllVotes[pos].Votes++
	// Increment vote
	Control <- ControlType{
		Cmd:            "vote",
		Hash:           resultHash,
		GroupPublicKey: groupPublicKey,
	}
	return nil
}

func (vc *DKG) Challenge(resultHash []byte, correctResult int) {

	GroupPublicKeyActive := make(map[string]*big.Int)
	var mutex = &sync.Mutex{}

	var addGroup = func(groupPublicKey *big.Int) {
		mutex.Lock()
		sK := fmt.Sprintf("%s", groupPublicKey)
		GroupPublicKeyActive[sK] = groupPublicKey
		mutex.Unlock()
	}

	var doneGroup = func(groupPublicKey *big.Int) {
		mutex.Lock()
		sK := fmt.Sprintf("%s", groupPublicKey)
		delete(GroupPublicKeyActive, sK)
		mutex.Unlock()
	}

	var process = func(result *result.Result, groupPublicKey *big.Int, isNew bool) string {
		addGroup(groupPublicKey)
		group := fmt.Sprintf("%s", groupPublicKey) // Public Key for Group - big number
		m1 := vc.PerGroup[group]

		vc.TNow = getCurrentBlockHeight() // has to have geth connect info

		if isNew {
			m1.AllResults = append(m1.AllResults, *result) // Should this merge the solutions together

			resultPkg := ResultVotes{
				Votes:       1,
				BlockHeight: vc.TNow,
				Group:       group,
			}
			m1.AllVotes = append(m1.AllVotes, resultPkg) // Should this merge the solutions together
		}

		leadResult := FindMostVotes(m1.AllVotes) // Summarize the votes for a particular solution
		vc.TFirst = m1.AllVotes[0].BlockHeight
		vc.PerGroup[group] = m1
		if vc.TNow > (vc.TFirst + vc.TConflict) {
			Complete <- CompleteType{FinishState: "success-timeout"}
			doneGroup(groupPublicKey)
			return "done"
		} else if m1.AllVotes[leadResult].Votes > vc.TMax {
			Complete <- CompleteType{FinishState: "success-enough-votes"}
			doneGroup(groupPublicKey)
			return "done"
		} else if correctResult == leadResult {
			vc.PerGroup[group] = m1
		} else if m1.AlreadySubmitted {
			vc.PerGroup[group] = m1
		} else if inResult(correctResult, m1.AllVotes, m1.AllResults) {
			Complete <- CompleteType{FinishState: "success-signed"}
			SubmitRessult(SignResult(resultHash))
			m1.AlreadySubmitted = true
		} else {
			SubmitRessult(m1.AllResults[correctResult])
			m1.AlreadySubmitted = true
		}
		vc.PerGroup[group] = m1
		return "loop"
	}

	for {
		select {
		case result := <-Submission:
			// GroupPublicKey = result.GroupPublicKey
			result.HashValue = resultHash
			if process(&result, result.GroupPublicKey, true) == "done" {
				return
			}

		case ctrl := <-Control:
			switch ctrl.Cmd {
			case "block-no":
				// iterate over all the groups that are active
				for _, GroupPublicKey := range GroupPublicKeyActive {
					if process(nil, GroupPublicKey, false) == "done" {
						return
					}
				}

			case "vote":
				if process(nil, ctrl.GroupPublicKey, false) == "done" {
					return
				}

			}
		}
	}

}

/*
Q's
	1. What is the data type of T_now, T_first - is it an integer from the block number of the chain.
		T_now referees to getting the current block height,
		T_first refers to "submitTime" - do we need a time or a block height?

	2. "leadResult" is a 256 bit integer (big.Int?)

	3a. in "submit(sign(resultHash))" - what public/private key is doing the signing?
	3b. in "submit(sign(resultHash))" - what is the "resultHash" - is that the hash of the leadResult.
		If so - then how is it converted - is that convert the leadResult to hex, then take the hash
		or convert to decimal and take the hash - or use as a packed-binary and take the hash.
		Which hash should be used?

	4. submit() referes to submitting the signature to the chain?

	5. Is this process per-group that is creating a signature?

	6. in "elif correctResult == leadresult" what determines correctResult

	7. Is correctResult a *big.Int

	8. In "submit(correctResut)" why is this not submit(hash(correctResult))" ?

#phase14questions

*/

func inResult(correctResult int, allVotes []ResultVotes, allResults []result.Result) bool {
	cr := allResults[correctResult]
	for _, val := range allResults {
		if cr.Equal(val) {
			return true
		}
	}
	return false
}

func SubmitRessult(aResult interface{}) {
	// TODO
	fmt.Printf("*********************************************\n* Submitted Result: %s\n*********************************************\n*\n", aResult)
}

func SignResult(resultHash []byte) []byte {
	// TODO
	return []byte("SignedResultHash:" + string(resultHash))
}

//////////////////////////////////////////////////////////////////

var getCurrentBlockHeight func() int

func (pm *PublishingMember) setupCurrentBlockHeight() error {
	blockCounter, err := pm.protocolConfig.chain.BlockCounter()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	getCurrentBlockHeight = func() int {
		t_now, err := blockCounter.CurrentBlock()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return 0
		}
		return t_now
	}
	return nil
}

func setup2CurrentBlockHeight(chain chain.Handle) error {
	blockCounter, err := chain.BlockCounter()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	getCurrentBlockHeight = func() int {
		t_now, err := blockCounter.CurrentBlock()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return 0
		}
		return t_now
	}
	return nil
}

func FindMostVotes(vgs []ResultVotes) int {
	maxPos := 0
	nVote := -1
	for pos := range vgs {
		if vgs[pos].Votes > nVote {
			nVote = vgs[pos].Votes
			maxPos = pos
		}
	}
	return maxPos
}

func FindHash(resultHash []byte, all []result.Result) (int, error) {
	for pos := range all {
		if bytes.Equal(resultHash, all[pos].HashValue) {
			return pos, nil
		}
	}
	return -1, fmt.Errorf("Did not find resultHash")
}
