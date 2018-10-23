package validaiton

/* Haskel Data Type

data Result = NoFaultFailure
            | FailureDQ { disqualified :: Array Bool }
            | PerfectSuccess { pubkey :: BeaconPubkey }
            | SuccessIA { pubkey :: BeaconPubkey
                        , inactive :: Array Bool }
            | SuccessDQ { pubkey :: BeaconPubkey
                        , disqualified :: Array Bool }
            | MixedSuccess { pubkey :: BeaconPubkey
                           , inactive :: Array Bool
                           , disqualified :: Array Bool }

*/

type ResultType int

const (
	NoFaultFailure ResultType = iota
	FailureDQ
	PerfectSuccess
	SuccessIA
	SuccessDQ
	MixedSuccess
)

type BeaconPubkey struct {
}

type Result struct {
	TypeIs       ResultType
	Disqualified []bool
	Pubkey       BeaconPubkey
	Inactive     []bool
}
