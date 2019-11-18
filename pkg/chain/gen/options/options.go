package options

import "math/big"

// TransactionOptions represents custom transaction options which will be
// used while invoking contracts methods.
type TransactionOptions struct {
	GasLimit uint64
	GasPrice *big.Int
}
