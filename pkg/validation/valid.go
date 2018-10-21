package validaiton

import "math/big"

type ValidationConfig struct {
	TMax      int // M_Max
	TConflict     // T_conflict
}

type ValidationState struct {
	TNow     int // Current block height at the current time
	PerGroup map[string]ValidationGroupState
}

type ValidationGroupState struct {
	TFirst     int          // T_first - Block height for the group - when the first result event occurred block height
	AllResults []ResultType // Set of all results
	LeadResult int          // Position of lead result
}

type ResultType struct {
	Votes       int      // How many votes this has
	BlockHeight int      // Block height of this results
	Result      *big.Int // The random value
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

func Phase14(vs *ValidationState, vc ValidationConfig) {

	for {
		select {
		case result := <-submission:
			vs.AllResults = append(vs.AllResults, result)
			leadResult = FindMostVotes(vs.AllResults)
			vs.TNow = getCurrentBlockHeight() // has to have geth connect info
			vs.TFirst = vs.AllResults[0].BlockHeight
			if vs.TNow > (vs.TFirst+vc.TConflict) || leadResults.votes > vc.TMax {
				return true
			} else if correctResult == leadResult || vs.AllResults[group].alreadySubmitted {
				continue
			} else if inResult(correctResult, vs.AllResults) {
				SubmitRessult(SignResult(resultHash))
				vs.AllResults[group].alreadySubmitted = true
			} else {
				SubmitRessult(correctResult)
				vs.AllResults[group].alreadySubmitted = true
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

func inResult(correctResult *big.Int, allResults []ResultType) {
	for ii := range ResultType {
		if big.Cmp(correctResult.Result, correctResult) == 0 {
			return true
		}
	}
	return false
}

func SubmitRessult(aResult interface{}) {
}

func SignResult(resultHash interface{}) {
}

func getCurrentBlockHeight() int {
	// TODO - xyzzy - work to dd.
	return 42
}

func FindMostVotes(vgs []ValidationGroupState) int {
	maxPos := 0
	nVote := -1
	for pos := range vgs {
		if vgs[pos].Votes > nVoge {
			nVote = vgs[pos].Votes
			maxPos = pos
		}
	}
	return maxPos
}
