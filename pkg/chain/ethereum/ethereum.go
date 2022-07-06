package ethereum

import (
	"context"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/rate"
)

var logger = log.Logger("keep-chain-ethereum")

// Chain represents a base, non-application-specific chain handle. It
// provides the implementation of generic features like balance monitor,
// block counter and similar.
type Chain struct {
	client ethutil.EthereumClient
}

// NewChain construct a new instance of the Ethereum chain handle.
func NewChain(
	ctx context.Context,
	config *ethereum.Config,
	client ethutil.EthereumClient,
) *Chain {
	return &Chain{
		client: wrapClientAddons(config, client),
	}
}

// wrapClientAddons wraps the client instance with add-ons like logging, rate
// limiting and so on.
func wrapClientAddons(
	config *ethereum.Config,
	client ethutil.EthereumClient,
) ethutil.EthereumClient {
	loggingClient := ethutil.WrapCallLogging(logger, client)

	if config.RequestsPerSecondLimit > 0 || config.ConcurrencyLimit > 0 {
		logger.Infof(
			"enabled ethereum client request rate limiter; "+
				"rps limit [%v]; "+
				"concurrency limit [%v]",
			config.RequestsPerSecondLimit,
			config.ConcurrencyLimit,
		)

		return ethutil.WrapRateLimiting(
			loggingClient,
			&rate.LimiterConfig{
				RequestsPerSecondLimit: config.RequestsPerSecondLimit,
				ConcurrencyLimit:       config.ConcurrencyLimit,
			},
		)
	}

	return loggingClient
}
