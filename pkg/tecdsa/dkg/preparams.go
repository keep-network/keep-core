package dkg

import (
	"context"
	"time"

	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/generator"
)

// PreParams represents tECDSA DKG pre-parameters that were not yet consumed
// by DKG protocol execution.
type PreParams struct {
	data *keygen.LocalPreParams
}

// NewPreParams constructs a new instance of tECDSA DKG pre-parameters based on
// the generated numbers.
func NewPreParams(data *keygen.LocalPreParams) *PreParams {
	return &PreParams{data}
}

// tssPreParamsPool is a pool holding TSS pre parameters. It autogenerates
// entries up to the pool size. When an entry is pulled from the pool it
// will generate a new entry.
type tssPreParamsPool struct {
	*generator.ParameterPool[PreParams]
	logger log.StandardLogger
}

// newTssPreParamsPool initializes a new TSS pre-parameters pool.
func newTssPreParamsPool(
	logger log.StandardLogger,
	scheduler *generator.Scheduler,
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

	newPreParamsFn := func(ctx context.Context) *PreParams {
		timingOutCtx, cancel := context.WithTimeout(ctx, generationTimeout)
		defer cancel()

		preParams, err := keygen.GeneratePreParamsWithContext(
			timingOutCtx,
			generationConcurrency,
		)
		// tss-lib returns generic errors saying "timeout or error while ...".
		// There are three possibilities:
		// 1. Pool canceled the parent `ctx`. This is normal and we should not
		//    log anything in this case.
		// 2. `timingOutCtx` timed out. It means the machine is not fast enough
		//    or that it was just unlucky. We should log a warning.
		// 3. There is some error from tss-lib generator. We log it as a warning
		//    because we'll re-attempt to generate parameters again.
		if err != nil && ctx.Err() == nil {
			logger.Warnf("failed to generate TSS pre-params: [%v]", err)
		}

		return &PreParams{preParams}
	}

	return &tssPreParamsPool{
		generator.NewParameterPool[PreParams](
			logger,
			scheduler,
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

func (np *noPersistence) Save(pp *PreParams) error {
	return nil
}
func (np *noPersistence) Delete(pp *PreParams) error {
	return nil
}
func (np *noPersistence) ReadAll() ([]*PreParams, error) {
	return []*PreParams{}, nil
}
