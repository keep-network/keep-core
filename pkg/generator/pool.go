package generator

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log/v2"
)

// ErrEmptyPool is returned by GetNow when the pool is empty.
var ErrEmptyPool = fmt.Errorf("pool is empty")

// Persistence defines the expected interface for storing and loading generated
// and not-yet-used parameters on persistent storage. Generating parameters is
// usually a computationally expensive operation and generated parameters should
// survive client restarts.
type Persistence[T any] interface {
	Save(*T) (*Persisted[T], error)
	Delete(*Persisted[T]) error
	ReadAll() ([]*Persisted[T], error)
}

// Persisted is a wrapper for the data that are stored, it adds an identifier.
// The identifier can be used in `Delete` function implementation to determine
// which entry should be removed from the persistent storage.
type Persisted[S any] struct {
	Data S
	ID   string
}

// ParameterPool autogenerates parameters based on the provided generation
// function up to the pool size. Parameters are stored in the cache and
// persisted using the provided persistence layer to survive client restarts.
// When a parameter is pulled from the pool, the pool starts generating a new
// parameter automatically. The pool submits the work to the provided scheduler
// instance and can be controlled by the scheduler.
type ParameterPool[T any] struct {
	persistence Persistence[T]
	pool        chan *Persisted[T]
}

// NewParameterPool creates a new instance of ParameterPool.
// The generateFn may return nil when the context passed to it has been
// cancelled or timed out during computations.
func NewParameterPool[T any](
	logger log.StandardLogger,
	scheduler *Scheduler,
	persistence Persistence[T],
	poolSize int,
	generateFn func(context.Context) *T,
	generateDelay time.Duration,
) *ParameterPool[T] {
	pool := make(chan *Persisted[T], poolSize)

	all, err := persistence.ReadAll()
	if err != nil {
		logger.Errorf("failed to read parameters from persistence: [%w]", err)
	}

	logger.Debugf("read [%d] parameters from persistence", len(all))

	for i, parameter := range all {
		// Load to the pool only the number of the parameters read from the persistence
		// that can fit within the pool's size, to avoid locking on writing to the
		// channel.
		if i >= poolSize {
			break
		}

		pool <- parameter
	}

	logger.Infof("loaded [%d] parameters from persistence", len(pool))

	scheduler.compute(func(ctx context.Context) {
		if len(pool) < poolSize {
			start := time.Now()

			generated := generateFn(ctx)

			// The generateFn returns nil when the context is done. We should not
			// add nil element to the pool.
			if generated == nil {
				return
			}

			persisted, err := persistence.Save(generated)
			if err != nil {
				logger.Errorf(
					"failed to persist generated parameter: [%w]",
					err,
				)
			}
			pool <- persisted

			logger.Infof(
				"generated new parameters, took: [%s] current pool size: [%d]",
				time.Since(start),
				len(pool),
			)
		}

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

		return &generated.Data, nil
	default:
		return nil, ErrEmptyPool
	}
}

// ParametersCount returns the number of parameters in the pool.
func (pp *ParameterPool[T]) ParametersCount() int {
	return len(pp.pool)
}
