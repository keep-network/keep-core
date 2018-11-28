package gjkr

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/keep-network/keep-core/pkg/beacon/relay/result"
)

func (dkg *DKG) NewChallengeState() (*ChallengeState, error) {
	TNow, err := dkg.GetChain().CurrentBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to get current block [%v]", err)
	}
	return &ChallengeState{
		TNow:             0,
		TFirst:           TNow,
		AllResults:       make([]result.Result, 0, 1),
		AllVotes:         make([]ResultVotes, 0, 1),
		LeadResult:       0,
		AlreadySubmitted: false,
		TConflict:        3,
		TMax:             5,
		dkg:              dkg,
	}, nil
}

func (vc *ChallengeState) GetChain() *Chain {
	return vc.dkg.GetChain()
}

func (vc *ChallengeState) ChallengeStateChange(
	currentResult *result.Result,
	resultHash []byte,
	correctResult int,
) error {
	TNow, err := vc.GetChain().CurrentBlock()
	if err != nil {
		return fmt.Errorf("failed to get current block [%v]", err)
	}
	vc.TNow = TNow
	_, err = vc.findHash(resultHash, vc.AllResults)
	if err != nil {
		vc.AllResults = append(vc.AllResults, *currentResult)
		resultPkg := ResultVotes{
			Votes:       1,
			BlockHeight: vc.TNow,
		}
		vc.AllVotes = append(vc.AllVotes, resultPkg)
	}
	leadResult := vc.findMostVotes(vc.AllVotes)
	vc.LeadResult = leadResult

	if vc.TNow > (vc.TFirst + vc.TConflict) {
		vc.stateReached(1)
		return nil
	} else if vc.AllVotes[leadResult].Votes > vc.TMax {
		vc.stateReached(2)
		return nil
	} else if correctResult == leadResult {
		vc.stateReached(3)
		return nil
	} else if vc.AlreadySubmitted {
		vc.stateReached(5)
		return nil
	} else if vc.inResult(correctResult, vc.AllVotes, vc.AllResults) {
		vc.stateReached(7)
		SubmitRessult(SignResult(resultHash))
		vc.AlreadySubmitted = true
		return nil
	} else {
		vc.stateReached(11)
		vc.AlreadySubmitted = true
		return nil
	}
	return nil
}

// Vote Adds 1 to vote for a specific result.
func (vc *ChallengeState) VoteForHash(resultHash []byte) error {
	if vc.TNow == 0 {
		return fmt.Errorf("did not see a challenge before a vote")
	}
	TNow, err := vc.GetChain().CurrentBlock()
	if err != nil {
		return fmt.Errorf("failed to get current block [%v]", err)
	}
	vc.TNow = TNow
	pos, err := vc.findHash(resultHash, vc.AllResults)
	if err != nil {
		return fmt.Errorf("invalid hash not found [%v]", err)
	}
	vc.AllVotes[pos].Votes++
	return nil
}

var stateReachedValue int

func (vc *ChallengeState) stateReached(stateNo int) {
	stateReachedValue = stateNo
	Complete <- CompleteType{
		StateReached:   stateNo,
		LineNoFileName: LF(2),
	}
}

func (vc *ChallengeState) findMostVotes(vgs []ResultVotes) int {
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

func (vc *ChallengeState) findHash(resultHash []byte, all []result.Result) (int, error) {
	for pos := range all {
		if bytes.Equal(resultHash, all[pos].HashValue) {
			return pos, nil
		}
	}
	return -1, fmt.Errorf("Did not find resultHash")
}

func (vc *ChallengeState) inResult(correctResult int, allVotes []ResultVotes, allResults []result.Result) bool {
	cr := allResults[correctResult]
	for _, val := range allResults {
		if reflect.DeepEqual(cr, val) {
			return true
		}
	}
	return false
}

func SubmitRessult(signedResult []byte) {
	Complete <- CompleteType{
		StateReached:   7,
		LineNoFileName: LF(2),
		Result:         signedResult,
	}
}

func SignResult(resultHash []byte) []byte {
	// TODO -- Implement Signature of this. -- Whith what key to sign?
	return []byte("SignedResultHash:" + fmt.Sprintf("%x", string(resultHash)))
}

type ResultVotes struct {
	Votes       int // How many votes this has
	BlockHeight int // Block height of this results
}

type CompleteType struct {
	StateReached   int
	LineNoFileName string
	Result         []byte
}

type VoteType struct {
	resultHash []byte
}

type BlockType struct {
	blockNo int
}

type ChallengeType struct {
	result        *result.Result
	resultHash    []byte
	correctResult int
}

var Challenge chan ChallengeType
var Vote chan VoteType
var Block chan BlockType
var Complete chan CompleteType

func init() {
	Challenge = make(chan ChallengeType)
	Vote = make(chan VoteType)
	Block = make(chan BlockType)
	Complete = make(chan CompleteType, 200)
}

func EventDispatch(vc *ChallengeState) {
	for {
		select {
		case challenge := <-Challenge:
			err := vc.ChallengeStateChange(challenge.result, challenge.resultHash, challenge.correctResult)
			if err != nil {
				fmt.Printf("error creating challenge: [%v]\n", err)
			} else {
				fmt.Printf("message processed: Challenge\n")
			}

		case vote := <-Vote:
			err := vc.VoteForHash(vote.resultHash)
			if err != nil {
				fmt.Printf("error with vote: [%v]\n", err)
			} else {
				fmt.Printf("message processed: Challenge\n")
			}

		case block := <-Block:
			_ = block
		}
	}
}
