package chain

import "math"

// OperatorID is a unique identifier of an operator assigned by the Sortition
// Pool when the operator enters the pool for the first time. ID is never
// changing for the given operator address.
type OperatorID = uint32

// MaxOperatorID is the maximum allowed value for OperatorID supported by
// Sortition Pool contract.
const MaxOperatorID = math.MaxUint32

// OperatorIDs is a list of OperatorID values.
type OperatorIDs []OperatorID
