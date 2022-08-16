package dkg

import (
	"time"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/miner"
)

// PreParams represents tECDSA DKG pre-parameters that were not yet consumed
// by DKG protocol execution.
type PreParams struct {
	data keygen.LocalPreParams
}

// NewPreParams constructs a new instance of tECDSA DKG pre-parameters based on
// the generated numbers.
func NewPreParams(data keygen.LocalPreParams) *PreParams {
	return &PreParams{data}
}

// tssPreParamsPool is a pool holding TSS pre parameters. It autogenerates
// entries up to the pool size. When an entry is pulled from the pool it
// will generate a new entry.
type tssPreParamsPool struct {
	*miner.ParameterPool[keygen.LocalPreParams]
	logger log.StandardLogger
}

// newTssPreParamsPool initializes a new TSS pre-parameters pool.
func newTssPreParamsPool(
	logger log.StandardLogger,
	poolSize int,
	generationTimeout time.Duration,
	generationDelay time.Duration,
	generationConcurrency int,
) *tssPreParamsPool {
	logger.Infof(
		"TSS pre-parameters target pool size is [%d], generation timeout is [%s] "+
			"generation delay is [%v], and concurrency level is [%d]",
		poolSize,
		generationTimeout,
		generationDelay,
		generationConcurrency,
	)

	newPreParamsFn := func() *keygen.LocalPreParams {
		preParams, err := keygen.GeneratePreParams(
			generationTimeout,
			generationConcurrency,
		)
		if err != nil {
			logger.Errorf("failed to generate TSS pre-params: [%v]", err)
		}

		return preParams
	}

	return &tssPreParamsPool{
		miner.NewParameterPool[keygen.LocalPreParams](
			logger,
			&miner.Miner{},   // TODO: pass as a parameter
			&noPersistence{}, // TODO: replace with a real persistence
			poolSize,
			newPreParamsFn,
			generationDelay,
		),
		logger,
	}
}

// TODO: temporary solution, will be replaced with a real persistence
type noPersistence struct{}

func (np *noPersistence) Save(pp *keygen.LocalPreParams) error {
	return nil
}
func (np *noPersistence) Delete(pp *keygen.LocalPreParams) error {
	return nil
}
func (np *noPersistence) ReadAll() ([]*keygen.LocalPreParams, error) {
	return []*keygen.LocalPreParams{}, nil
}
