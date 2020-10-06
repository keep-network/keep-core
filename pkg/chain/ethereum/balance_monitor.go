package ethereum

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type BalanceSource func(address common.Address) (*big.Int, error)

type BalanceMonitor struct {
	balanceSource BalanceSource
}

func (bm *BalanceMonitor) Observe(
	ctx context.Context,
	address common.Address,
	alertThreshold *big.Int,
	tick time.Duration,
) {
	check := func() {
		balance, err := bm.balanceSource(address)
		if err != nil {
			logger.Errorf("ethereum balance monitor error: [%v]", err)
			return
		}

		if balance.Cmp(alertThreshold) < -1 {
			logger.Errorf(
				"ethereum balance is below [%v] wei; "+
					"please fund your operator account",
				alertThreshold.Text(10),
			)
		}
	}

	go func() {
		ticker := time.NewTicker(tick)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				check()
			case <-ctx.Done():
				return
			}
		}
	}()
}
