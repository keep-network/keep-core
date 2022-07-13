package firewall

import (
	"fmt"
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"github.com/keep-network/keep-core/pkg/operator"
	"reflect"
	"testing"
	"time"

	"github.com/keep-network/keep-common/pkg/cache"
)

const cachingPeriod = time.Second

func TestValidate_OperatorNotRecognized_NoApplications(t *testing.T) {
	policy := &anyApplicationPolicy{
		applications:        []Application{},
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	_, peerOperatorPublicKey, err := operator.GenerateKeyPair(
		local_v1.DefaultCurve,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = policy.Validate(peerOperatorPublicKey)

	expectedError := fmt.Errorf(
		"remote peer has not been recognized by any application",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestValidate_OperatorNotRecognized_MultipleApplications(t *testing.T) {
	_, peerOperatorPublicKey, err := operator.GenerateKeyPair(
		local_v1.DefaultCurve,
	)

	// Use two application: the first returns error during check, the second
	// does not recognize the operator
	application := newMockApplication()
	application.setIsRecognized(peerOperatorPublicKey, result{
		isRecognized: false,
		err:          fmt.Errorf("dummy error"),
	})
	policy := &anyApplicationPolicy{
		applications: []Application{
			newMockApplication(),
			newMockApplication()},
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	if err != nil {
		t.Fatal(err)
	}

	err = policy.Validate(peerOperatorPublicKey)

	expectedError := fmt.Errorf(
		"remote peer has not been recognized by any application",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestValidate_OperatorRecognized_FirstApplicationRecognizes(t *testing.T) {
	_, peerOperatorPublicKey, err := operator.GenerateKeyPair(
		local_v1.DefaultCurve,
	)
	if err != nil {
		t.Fatal(err)
	}

	application := newMockApplication()
	application.setIsRecognized(peerOperatorPublicKey, result{
		isRecognized: true,
		err:          nil,
	})

	policy := &anyApplicationPolicy{
		applications: []Application{
			application,
			newMockApplication()},
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	err = policy.Validate(peerOperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidate_OperatorRecognized_SecondApplicationRecognizes(t *testing.T) {
	_, peerOperatorPublicKey, err := operator.GenerateKeyPair(
		local_v1.DefaultCurve,
	)
	if err != nil {
		t.Fatal(err)
	}

	application := newMockApplication()
	application.setIsRecognized(peerOperatorPublicKey, result{
		isRecognized: true,
		err:          nil,
	})

	policy := &anyApplicationPolicy{
		applications: []Application{
			newMockApplication(),
			application},
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	err = policy.Validate(peerOperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidate_OperatorRecognized_FirstApplicationReturnedError(t *testing.T) {
	_, peerOperatorPublicKey, err := operator.GenerateKeyPair(
		local_v1.DefaultCurve,
	)
	if err != nil {
		t.Fatal(err)
	}

	// First application returns error during operator recognition check, but
	// the second recognizes the operator.
	application1 := newMockApplication()
	application1.setIsRecognized(peerOperatorPublicKey, result{
		isRecognized: false,
		err:          fmt.Errorf("dummy error"),
	})

	application2 := newMockApplication()
	application2.setIsRecognized(peerOperatorPublicKey, result{
		isRecognized: true,
		err:          nil,
	})

	policy := &anyApplicationPolicy{
		applications: []Application{
			application1,
			application2},
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	err = policy.Validate(peerOperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidate_OperatorRecognized_Cached(t *testing.T) {
	_, peerOperatorPublicKey, err := operator.GenerateKeyPair(
		local_v1.DefaultCurve,
	)
	if err != nil {
		t.Fatal(err)
	}

	application := newMockApplication()
	application.setIsRecognized(peerOperatorPublicKey, result{
		isRecognized: true,
		err:          nil,
	})

	policy := &anyApplicationPolicy{
		applications:        []Application{application},
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	err = policy.Validate(peerOperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure the application does not recognize the operator anymore.
	// Validation should still succeed, since the cached result should be used.
	application.setIsRecognized(peerOperatorPublicKey, result{
		isRecognized: false,
		err:          nil,
	})

	err = policy.Validate(peerOperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidate_OperatorNotRecognized_CacheEmptied(t *testing.T) {
	_, peerOperatorPublicKey, err := operator.GenerateKeyPair(
		local_v1.DefaultCurve,
	)
	if err != nil {
		t.Fatal(err)
	}

	application := newMockApplication()
	application.setIsRecognized(peerOperatorPublicKey, result{
		isRecognized: true,
		err:          nil,
	})

	policy := &anyApplicationPolicy{
		applications:        []Application{application},
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	err = policy.Validate(peerOperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure the application does not recognize the operator anymore.
	// Wait for the caching period to end. Validation should fail, as the
	// operator has been removed from the cache.
	application.setIsRecognized(peerOperatorPublicKey, result{
		isRecognized: false,
		err:          nil,
	})

	time.Sleep(cachingPeriod)

	err = policy.Validate(peerOperatorPublicKey)

	expectedError := fmt.Errorf(
		"remote peer has not been recognized by any application",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestValidate_OperatorNotRecognized_Cached(t *testing.T) {
	_, peerOperatorPublicKey, err := operator.GenerateKeyPair(
		local_v1.DefaultCurve,
	)
	if err != nil {
		t.Fatal(err)
	}

	application := newMockApplication()
	application.setIsRecognized(peerOperatorPublicKey, result{
		isRecognized: false,
		err:          nil,
	})

	policy := &anyApplicationPolicy{
		applications:        []Application{application},
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	err = policy.Validate(peerOperatorPublicKey)
	expectedError := fmt.Errorf(
		"remote peer has not been recognized by any application",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}

	// Ensure the application recognizes the operator, but the validation should
	// fail since the result from the previous call has been cached.
	application.setIsRecognized(peerOperatorPublicKey, result{
		isRecognized: true,
		err:          nil,
	})

	err = policy.Validate(peerOperatorPublicKey)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}
}

func TestValidate_OperatorRecognized_CacheEmptied(t *testing.T) {
	_, peerOperatorPublicKey, err := operator.GenerateKeyPair(
		local_v1.DefaultCurve,
	)
	if err != nil {
		t.Fatal(err)
	}

	application := newMockApplication()
	application.setIsRecognized(peerOperatorPublicKey, result{
		isRecognized: false,
		err:          nil,
	})

	policy := &anyApplicationPolicy{
		applications:        []Application{application},
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	err = policy.Validate(peerOperatorPublicKey)
	expectedError := fmt.Errorf(
		"remote peer has not been recognized by any application",
	)
	if !reflect.DeepEqual(expectedError, err) {
		t.Errorf(
			"unexpected error\nexpected: %v\nactual:   %v\n",
			expectedError,
			err,
		)
	}

	// Ensure the application recognizes the operator. Wait for the caching
	// period to elapse. The validation should pass since the result from the
	// previous call has been removed from the cache.
	application.setIsRecognized(peerOperatorPublicKey, result{
		isRecognized: true,
		err:          nil,
	})

	time.Sleep(cachingPeriod)

	err = policy.Validate(peerOperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}
}

func newMockApplication() *mockApplication {
	return &mockApplication{
		results: make(map[*operator.PublicKey]result),
	}
}

type result struct {
	isRecognized bool
	err          error
}

type mockApplication struct {
	results map[*operator.PublicKey]result
}

func (ma *mockApplication) setIsRecognized(
	operatorPublicKey *operator.PublicKey,
	result result) {
	ma.results[operatorPublicKey] = result
}

func (ma *mockApplication) IsRecognized(operatorPublicKey *operator.PublicKey) (
	bool, error) {
	result := ma.results[operatorPublicKey]
	return result.isRecognized, result.err
}
