package gjkr

import "github.com/keep-network/keep-core/pkg/beacon/relay/result"

type ChallengeState struct {
	TFirst           int             // T_first - Block height for the group - when the first result event occurred block height
	AllResults       []result.Result // Set of all results
	AllVotes         []ResultVotes   // Set of all results
	LeadResult       int             // Position of lead result
	AlreadySubmitted bool            //
	TConflict        int             // == votingDuration
	TMax             int             // == dishonestThreshold
	dkg              *DKG
}
