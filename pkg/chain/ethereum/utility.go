package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

func (euc *ethereumUtilityChain) Genesis() error {
	_, err := euc.keepRandomBeaconOperatorContract.Genesis(big.NewInt(0))
	return err
}

func (euc *ethereumUtilityChain) RequestRelayEntry(seed *big.Int) *async.EventRequestPromise {
	promise := &async.EventRequestPromise{}

	_, err := euc.keepRandomBeaconServiceContract.RequestRelayEntry(seed, common.BytesToAddress([]byte{}), "", big.NewInt(1), big.NewInt(0))
	if err != nil {
		promise.Fail(err)
	}

	return promise
}
