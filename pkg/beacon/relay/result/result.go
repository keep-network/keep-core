package result

import (
	"fmt"
	"math/big"
)

type ResultType int

const (
	NoFaultFailure ResultType = iota
	FailureDQ
	PerfectSuccess
	SuccessIA
	SuccessDQ
	MixedSuccess
)

type Result struct {
	Type           ResultType
	GroupPublicKey *big.Int `json:"pubkey"`
	Disqualified   []int    `json:"disqualified"`
	Inactive       []int    `json:"inactive"`
}

func (r *Result) Hash() []byte {
	return []byte(fmt.Sprintf("%v", r))
}
