package sortition

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
