package miner

import (
	"fmt"
	"time"

	"github.com/ipfs/go-log"
)

// Persistence defines the expected interface for storing and loading generated
// and not-yet-used parameters on persistent storage. Generating parameters is
// usually a computationally expensive operation and generated parameters should
// survive client restarts.
type Persistence[T any] interface {
	Save(*T) error
	Delete(*T) error
	ReadAll() ([]*T, error)
}

// ParameterPool autogenerates parameters based on the provided generation
// function up to the pool size. Parameters are stored in the cache and
// persisted using the provided persistence layer to survive client restarts.
// When a parameter is pulled from the pool, the pool starts generating a new
// parameter automatically. The pool submits the work to the provided miner
// instance and can be controlled by the miner.
type ParameterPool[T any] struct {
	persistence Persistence[T]
	pool        chan *T
}

// NewParameterPool creates a new instance of ParameterPool.
func NewParameterPool[T any](
	logger log.StandardLogger,
	miner *Miner,
	persistence Persistence[T],
	targetSize int,
	generateFn func() *T,
	generateDelay time.Duration,
) *ParameterPool[T] {
	pool := make(chan *T, targetSize)

	all, err := persistence.ReadAll()
	if err != nil {
		logger.Errorf("failed to read parameters from persistence: [%w]", err)
	}
	for _, parameter := range all {
		pool <- parameter
	}

	miner.Mine(func() {
		start := time.Now()

		generated := generateFn()

		err := persistence.Save(generated)
		if err != nil {
			logger.Errorf(
				"failed to persist generated parameter: [%w]",
				err,
			)
		}
		pool <- generated

		logger.Infof(
			"generated new parameters, took: [%s] current pool size: [%d]",
			time.Since(start),
			len(pool),
		)

		// Wait some time after delivering the result regardless if the delivery
		// took some time or not. We want to ensure all other processes of the
		// client receive access to CPU.
		time.Sleep(generateDelay)
	})

	return &ParameterPool[T]{
		persistence: persistence,
		pool:        pool,
	}
}

// Get returns a new parameter. It is fetched from the pool or generated if the
// pool is empty.
func (pp *ParameterPool[T]) Get() (*T, error) {
	generated := <-pp.pool
	err := pp.persistence.Delete(generated)
	if err != nil {
		return nil, fmt.Errorf("could not delete persisted parameter: [%w]", err)
	}
	return generated, nil
}

// CurrentSize returns the current size of the pool - the number of available
// parameters.
func (pp *ParameterPool[T]) CurrentSize() int {
	return len(pp.pool)
}
