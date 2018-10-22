package validaiton

import "math/big"

type ValidationConfig struct {
	TMax      int // M_Max
	TConflict int // T_conflict
}

type ValidationState struct {
	TNow     int // Current block height at the current time
	TFirst   int //
	PerGroup map[string]ValidationGroupState
}

type ValidationGroupState struct {
	TFirst           int          // T_first - Block height for the group - when the first result event occurred block height
	AllResults       []ResultType // Set of all results
	LeadResult       int          // Position of lead result
	alreadySubmitted bool         //
}

type ResultType struct {
	Votes       int      // How many votes this has
	BlockHeight int      // Block height of this results
	Result      *big.Int // The random value
	group       string
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

var submission chan ResultType

func Phase14(vs *ValidationState, vc ValidationConfig) {

	correctResult := 0   // ResultType{} // Where is this from?
	resultHash := "TODO" // where is this from?
	for {
		select {
		case result := <-submission:
			group := result.group
			m1 := vs.PerGroup[group]
			m1.AllResults = append(m1.AllResults, result)
			leadResult := FindMostVotes(m1.AllResults)
			vs.TNow = getCurrentBlockHeight() // has to have geth connect info
			vs.TFirst = m1.AllResults[0].BlockHeight
			if vs.TNow > (vs.TFirst+vc.TConflict) || m1.AllResults[leadResult].Votes > vc.TMax {
				return
			} else if correctResult == leadResult || m1.alreadySubmitted {
				continue
			} else if inResult(correctResult, m1.AllResults) {
				SubmitRessult(SignResult(resultHash))
				m1.alreadySubmitted = true
			} else {
				SubmitRessult(correctResult)
				m1.alreadySubmitted = true
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

*/

func inResult(correctResult int, allResults []ResultType) bool {
	cr := allResults[correctResult].Result
	for _, val := range allResults {
		if cr.Cmp(val.Result) == 0 {
			return true
		}
	}
	return false
}

func SubmitRessult(aResult interface{}) string {
	return ""
}

func SignResult(resultHash interface{}) string {
	return ""
}

func getCurrentBlockHeight() int {
	// TODO - xyzzy - work to dd.
	return 42
}

func FindMostVotes(vgs []ResultType) int {
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
