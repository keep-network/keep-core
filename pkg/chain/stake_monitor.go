package chain

import "github.com/keep-network/keep-core/pkg/operator"

// StakeMonitor is an interface that provides ability to check and monitor
// the stake for the provided address.
type StakeMonitor interface {
	// HasMinimumStake checks if the specified account has enough active stake
	// to become network operator and that the operator contract the client is
	// working with has been authorized for potential slashing.
	//
	// Having the required minimum of active stake makes the operator eligible
	// to join the network. If the active stake is not currently undelegating,
	// operator is also eligible for work selection.
	HasMinimumStake(operatorPublicKey *operator.PublicKey) (bool, error)
}
