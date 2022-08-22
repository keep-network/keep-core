package generator

import (
	"context"
	"math/big"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

// TestGet covers the most basic path - calling `Get()` function multiple times
// and making sure result is always returned, assuming there are no errors from
// the persistence layer. The pool size is lower than the number of parameters
// fetched from the pool so this test also ensures the pool does not stop
// generating numbers for future `Get()` calls.
func TestGet(t *testing.T) {
	pool, scheduler, _ := newTestPool(5)
	defer scheduler.stop()

	for i := 0; i < 70; i++ {
		e, err := pool.Get()
		if err != nil {
			t.Errorf("unexpected error: [%v]", err)
		}
		if e == nil {
			t.Errorf("expected not-nil parameter")
		}
	}
}

// TestStop ensures the pool honors the stop signal send to the scheduler and it
// does not keep generating params in some internal loop.
func TestStop(t *testing.T) {
	pool, scheduler, _ := newTestPool(50000)

	// give some time to generate parameters and stop
	time.Sleep(25 * time.Millisecond)
	scheduler.stop()

	// give some time for the generation process to stop and capture the number
	// of parameters generated
	time.Sleep(10 * time.Millisecond)
	size := pool.CurrentSize()

	// wait some time and make sure no new parameters are generated
	time.Sleep(20 * time.Millisecond)
	if size != pool.CurrentSize() {
		t.Errorf("expected no new parameters to be generated")
	}
}

// TestPersist ensures parameters generated by the pool are persisted.
func TestPersist(t *testing.T) {
	pool, scheduler, persistence := newTestPool(50000)

	// give some time to generate parameters and stop
	time.Sleep(25 * time.Millisecond)
	scheduler.stop()

	// give some time for the generation process to stop
	time.Sleep(10 * time.Millisecond)

	if pool.CurrentSize() != persistence.parameterCount() {
		t.Errorf("not all parameters have been persisted")
	}
}

// TestReadAll ensures pool reads parameters from the persistence before
// generating new ones.
func TestReadAll(t *testing.T) {
	persistence := &mockPersistence{storage: map[string]*big.Int{
		"100": big.NewInt(100),
		"200": big.NewInt(200),
	}}

	pool, scheduler := newTestPoolWithPersistence(100, persistence)
	defer scheduler.stop()

	e, err := pool.Get()
	if err != nil {
		t.Errorf("unexpected error: [%v]", err)
	}
	testutils.AssertBigIntsEqual(t, "parameter value", big.NewInt(100), e)

	e, err = pool.Get()
	if err != nil {
		t.Errorf("unexpected error: [%v]", err)
	}
	testutils.AssertBigIntsEqual(t, "parameter value", big.NewInt(200), e)

}

// TestDelete ensures parameters fetched from the pool are deleted from the
// persistence layer.
func TestDelete(t *testing.T) {
	persistence := &mockPersistence{storage: map[string]*big.Int{
		"100": big.NewInt(100),
	}}

	pool, scheduler := newTestPoolWithPersistence(100, persistence)
	defer scheduler.stop()

	e, err := pool.Get()
	if err != nil {
		t.Errorf("unexpected error: [%v]", err)
	}
	if persistence.isPresent(e) {
		t.Errorf("element should be deleted from persistence: [%v]", e)
	}
}

func newTestPool(targetSize int) (*ParameterPool[big.Int], *Scheduler, *mockPersistence) {
	persistence := &mockPersistence{storage: make(map[string]*big.Int)}
	pool, scheduler := newTestPoolWithPersistence(targetSize, persistence)
	return pool, scheduler, persistence
}

func newTestPoolWithPersistence(
	targetSize int,
	persistence *mockPersistence,
) (*ParameterPool[big.Int], *Scheduler) {
	generateFn := func(context.Context) *big.Int {
		time.Sleep(5 * time.Millisecond)
		return big.NewInt(time.Now().UnixMilli())
	}

	scheduler := &Scheduler{}

	return NewParameterPool[big.Int](
		logger,
		scheduler,
		persistence,
		targetSize,
		generateFn,
		time.Duration(0), // no delay
	), scheduler
}

type mockPersistence struct {
	storage map[string]*big.Int
	mutex   sync.RWMutex
}

func (mp *mockPersistence) Save(element *big.Int) error {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()

	mp.storage[element.String()] = element
	return nil
}

func (mp *mockPersistence) Delete(element *big.Int) error {
	mp.mutex.Lock()
	defer mp.mutex.Unlock()

	delete(mp.storage, element.String())
	return nil
}

func (mp *mockPersistence) ReadAll() ([]*big.Int, error) {
	mp.mutex.RLock()
	defer mp.mutex.RUnlock()

	all := make([]*big.Int, 0, len(mp.storage))
	for _, v := range mp.storage {
		all = append(all, v)
	}
	// sorting is needed for TestReadAll
	sort.Slice(all, func(i, j int) bool {
		return all[i].Cmp(all[j]) < 0
	})
	return all, nil
}

func (mp *mockPersistence) parameterCount() int {
	mp.mutex.RLock()
	defer mp.mutex.RUnlock()

	return len(mp.storage)
}

func (mp *mockPersistence) isPresent(element *big.Int) bool {
	mp.mutex.RLock()
	defer mp.mutex.RUnlock()

	_, ok := mp.storage[element.String()]
	return ok
}