package sortition

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

func (cp *ConjunctionPolicy) ShouldJoin() bool {
	for _, policy := range cp.policies {
		if !policy.ShouldJoin() {
			return false
		}
	}

	return true
}
