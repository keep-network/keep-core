package generator

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log"
)

// ErrEmptyPool is returned by GetNow when the pool is empty.
var ErrEmptyPool = fmt.Errorf("pool is empty")

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
// parameter automatically. The pool submits the work to the provided scheduler
// instance and can be controlled by the scheduler.
type ParameterPool[T any] struct {
	persistence Persistence[T]
	pool        chan *T
}

// NewParameterPool creates a new instance of ParameterPool.
// The generateFn may return nil when the context passed to it has been
// cancelled or timed out during computations.
func NewParameterPool[T any](
	logger log.StandardLogger,
	scheduler *Scheduler,
	persistence Persistence[T],
	targetSize int,
	generateFn func(context.Context) *T,
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

	scheduler.compute(func(ctx context.Context) {
		start := time.Now()

		generated := generateFn(ctx)

		// The generateFn returns nil when the context is done. We should not
		// add nil element to the pool.
		if generated == nil {
			return
		}

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

// GetNow returns a new parameter from the pool. Returns ErrEmptyPool when the
// pool is empty.
func (pp *ParameterPool[T]) GetNow() (*T, error) {
	select {
	case generated := <-pp.pool:
		err := pp.persistence.Delete(generated)
		if err != nil {
			return nil, fmt.Errorf(
				"could not delete persisted parameter: [%w]",
				err,
			)
		}

		return generated, nil
	default:
		return nil, ErrEmptyPool
	}
}

// CurrentSize returns the current size of the pool - the number of available
// parameters.
func (pp *ParameterPool[T]) CurrentSize() int {
	return len(pp.pool)
}
