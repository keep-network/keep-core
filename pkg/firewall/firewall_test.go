package firewall

import (
	"github.com/keep-network/keep-core/pkg/chain/local_v1"
	"github.com/keep-network/keep-core/pkg/operator"
	"math/big"
	"testing"
	"time"

	"github.com/keep-network/keep-common/pkg/cache"
)

var minimumStake = big.NewInt(1000)
var cachingPeriod = time.Second

func TestHasMinimumStake(t *testing.T) {
	stakeMonitor := local_v1.NewStakeMonitor(minimumStake)
	policy := &minimumStakePolicy{
		stakeMonitor:        stakeMonitor,
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	_, remotePeerOperatorPublicKey, err := operator.GenerateKeyPair(local_v1.DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	err = stakeMonitor.StakeTokens(remotePeerOperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	if err := policy.Validate(
		remotePeerOperatorPublicKey,
	); err != nil {
		t.Fatalf("validation should pass: [%v]", err)
	}
}

func TestHasNoMinimumStake(t *testing.T) {
	stakeMonitor := local_v1.NewStakeMonitor(minimumStake)
	policy := &minimumStakePolicy{
		stakeMonitor:        stakeMonitor,
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	_, remotePeerOperatorPublicKey, err := operator.GenerateKeyPair(local_v1.DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	if err := policy.Validate(
		remotePeerOperatorPublicKey,
	); err != errNoMinimumStake {
		t.Fatalf(
			"unexpected validation error\nactual:   [%v]\nexpected: [%v]",
			err,
			errNoMinimumStake,
		)
	}
}

func TestCachesHasMinimumStake(t *testing.T) {
	stakeMonitor := local_v1.NewStakeMonitor(minimumStake)
	policy := &minimumStakePolicy{
		stakeMonitor:        stakeMonitor,
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	_, remotePeerOperatorPublicKey, err := operator.GenerateKeyPair(local_v1.DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	err = stakeMonitor.StakeTokens(remotePeerOperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	if err := policy.Validate(
		remotePeerOperatorPublicKey,
	); err != nil {
		t.Fatalf("validation should pass: [%v]", err)
	}

	err = stakeMonitor.UnstakeTokens(remotePeerOperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	// still caching the old result
	if err := policy.Validate(
		remotePeerOperatorPublicKey,
	); err != nil {
		t.Fatalf("validation should pass: [%v]", err)
	}

	time.Sleep(time.Second)

	// no longer caches the previous result
	if err := policy.Validate(
		remotePeerOperatorPublicKey,
	); err != errNoMinimumStake {
		t.Fatalf(
			"unexpected validation error\nactual:   [%v]\nexpected: [%v]",
			err,
			errNoMinimumStake,
		)
	}
}

func TestCachesHasNoMinimumStake(t *testing.T) {
	stakeMonitor := local_v1.NewStakeMonitor(minimumStake)
	policy := &minimumStakePolicy{
		stakeMonitor:        stakeMonitor,
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	_, remotePeerOperatorPublicKey, err := operator.GenerateKeyPair(local_v1.DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	if err := policy.Validate(
		remotePeerOperatorPublicKey,
	); err != errNoMinimumStake {
		t.Fatalf(
			"unexpected validation error\nactual:   [%v]\nexpected: [%v]",
			err,
			errNoMinimumStake,
		)
	}

	err = stakeMonitor.StakeTokens(remotePeerOperatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	// still caching the old result
	if err := policy.Validate(
		remotePeerOperatorPublicKey,
	); err != errNoMinimumStake {
		t.Fatalf(
			"unexpected validation error\nactual:   [%v]\nexpected: [%v]",
			err,
			errNoMinimumStake,
		)
	}

	time.Sleep(time.Second)

	// no longer caches the previous result
	if err := policy.Validate(
		remotePeerOperatorPublicKey,
	); err != nil {
		t.Fatalf("validation should pass: [%v]", err)
	}
}
