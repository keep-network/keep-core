package ethereum

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
)

// Values related with balance monitoring.
// defaultBalanceAlertThreshold determines the alert threshold below which
// the alert should be triggered.
var defaultBalanceAlertThreshold = ethereum.WrapWei(
	big.NewInt(500000000000000000),
)

// defaultBalanceMonitoringTick determines how often the monitoring
// check should be triggered.
const defaultBalanceMonitoringTick = 10 * time.Minute

// defaultBalanceMonitoringRetryTimeout determines the timeout for balance check
// at each tick.
const defaultBalanceMonitoringRetryTimeout = 5 * time.Minute

// initializeBalanceMonitoring sets up the balance monitoring process
func (c *Chain) initializeBalanceMonitoring(ctx context.Context) {
	balanceMonitor, err := c.balanceMonitor()
	if err != nil {
		logger.Errorf("could not get balance monitor [%v]", err)
		return
	}

	alertThreshold := defaultBalanceAlertThreshold
	if value := c.config.BalanceAlertThreshold; value != nil {
		alertThreshold = value
	}

	balanceMonitor.Observe(
		ctx,
		c.Address(),
		alertThreshold,
		defaultBalanceMonitoringTick,
		defaultBalanceMonitoringRetryTimeout,
	)

	logger.Infof(
		"started balance monitoring for address [%v] "+
			"with the alert threshold set to [%v] wei",
		c.Address().Hex(),
		alertThreshold,
	)
}

// balanceMonitor returns a balance monitor.
func (c *Chain) balanceMonitor() (*ethutil.BalanceMonitor, error) {
	weiBalanceOf := func(
		address common.Address,
	) (*ethereum.Wei, error) {
		ctx, cancelCtx := context.WithTimeout(
			context.Background(),
			30*time.Second,
		)
		defer cancelCtx()

		balance, err := c.client.BalanceAt(ctx, address, nil)
		if err != nil {
			return nil, err
		}

		return ethereum.WrapWei(balance), nil
	}

	return ethutil.NewBalanceMonitor(weiBalanceOf), nil
}
