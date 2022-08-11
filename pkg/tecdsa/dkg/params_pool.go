package dkg

import (
	"fmt"
	"time"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/ipfs/go-log"
)

// tssPreParamsPool is a pool holding TSS pre parameters. It autogenerates
// entries up to the pool size. When an entry is pulled from the pool it
// will generate a new entry.
type tssPreParamsPool struct {
	logger         log.StandardLogger
	pool           chan *keygen.LocalPreParams
	newPreParamsFn func() (*keygen.LocalPreParams, error)
}

// newTssPreParamsPool initializes a new TSS pre-parameters pool.
func newTssPreParamsPool(
	logger log.StandardLogger,
	poolSize int,
	generationTimeout time.Duration,
	generationConcurrency int,
) *tssPreParamsPool {
	logger.Infof(
		"TSS pre-parameters target pool size is [%v], generation timeout is [%v], and concurrency level is [%v]",
		poolSize,
		generationTimeout,
		generationConcurrency,
	)

	newPreParamsFn := func() (*keygen.LocalPreParams, error) {
		preParams, err := keygen.GeneratePreParams(
			generationTimeout,
			generationConcurrency,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to generate TSS pre-params: [%v]",
				err,
			)
		}

		return preParams, nil
	}

	tssPreParamsPool := &tssPreParamsPool{
		logger:         logger,
		pool:           make(chan *keygen.LocalPreParams, poolSize),
		newPreParamsFn: newPreParamsFn,
	}

	go tssPreParamsPool.pumpPool()

	return tssPreParamsPool
}

func (t *tssPreParamsPool) pumpPool() {
	for {
		t.logger.Info("generating new tss pre parameters")

		start := time.Now()

		params, err := t.newPreParamsFn()
		if err != nil {
			t.logger.Warningf(
				"failed to generate tss pre parameters after [%s]: [%v]",
				time.Since(start),
				err,
			)
			continue
		}

		t.logger.Infof(
			"generated new tss pre parameters, took: [%s], "+
				"current pool size: [%d]",
			time.Since(start),
			len(t.pool)+1,
		)

		t.pool <- params
	}
}

// get returns TSS pre parameters from the pool. It pumps the pool after getting
// an entry. If the pool is empty it will wait for a new entry to be generated.
func (t *tssPreParamsPool) get() *keygen.LocalPreParams {
	return <-t.pool
}
