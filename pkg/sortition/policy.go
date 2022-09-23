package sortition

import "github.com/ipfs/go-log"

// JoinPolicy determines how the client is supposed to join to the sortition
// pool. The policy can encapsulate special conditions that the client want
// to fulfill before joining the sortition pool.
type JoinPolicy interface {
	// ShouldJoin indicates whether the joining condition is fulfilled and
	// the client should join the pool.
	ShouldJoin() bool
}

// UnconditionalJoinPolicy is a policy that doesn't enforce any conditions
// for joining the sortition pool.
var UnconditionalJoinPolicy = &unconditionalJoinPolicy{}

type unconditionalJoinPolicy struct{}

func (ujp *unconditionalJoinPolicy) ShouldJoin() bool {
	return true
}

// ConjunctionPolicy is a JoinPolicy implementation requiring all the provided
// policies to pass before allowing to join the sortition pool.
type ConjunctionPolicy struct {
	policies []JoinPolicy
}

func NewConjunctionPolicy(
	policies ...JoinPolicy,
) *ConjunctionPolicy {
	return &ConjunctionPolicy{policies}
}

func (cp *ConjunctionPolicy) ShouldJoin() bool {
	for _, policy := range cp.policies {
		if !policy.ShouldJoin() {
			return false
		}
	}

	return true
}

// BetaOperatorPolicy is a JoinPolicy implementation checking chaosnet and
// beta operator status. If chaosnet has been deactivated, the operator is
// allowed to join the pool. If chaosnet is active and the operator is beta
// operator, the operator is allowed to join the pool. If chaosnet is active and
// the operator is not beta operator, the operator is not allowed to join the
// pool.
type BetaOperatorPolicy struct {
	chain  Chain
	logger log.StandardLogger
}

func NewBetaOperatorPolicy(
	chain Chain,
	logger log.StandardLogger,
) *BetaOperatorPolicy {
	return &BetaOperatorPolicy{
		chain,
		logger,
	}
}

func (bop *BetaOperatorPolicy) ShouldJoin() bool {
	isChaosnetActive, err := bop.chain.IsChaosnetActive()
	if err != nil {
		bop.logger.Errorf("could not determine if chaosnet is active: [%v]", err)
		return false
	}

	if !isChaosnetActive {
		return true
	}

	isBetaOperator, err := bop.chain.IsBetaOperator()
	if err != nil {
		bop.logger.Errorf("could not determine beta operator status: [%v]", err)
		return false
	}

	return isBetaOperator
}
