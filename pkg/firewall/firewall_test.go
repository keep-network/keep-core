package firewall

import (
	"math/big"
	"testing"
	"time"

	"github.com/keep-network/keep-common/pkg/cache"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net/key"
)

var minimumStake = big.NewInt(1000)
var cachingPeriod = time.Second

func TestHasMinimumStake(t *testing.T) {
	stakeMonitor := local.NewStakeMonitor(minimumStake)
	policy := &minimumStakePolicy{
		stakeMonitor:        stakeMonitor,
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	_, remotePeerPublicKey, err := key.GenerateStaticNetworkKey()
	if err != nil {
		t.Fatal(err)
	}
	remotePeerAddress := key.NetworkPubKeyToEthAddress(remotePeerPublicKey)

	stakeMonitor.StakeTokens(remotePeerAddress)

	if err := policy.Validate(
		key.NetworkKeyToECDSAKey(remotePeerPublicKey),
	); err != nil {
		t.Fatalf("validation should pass: [%v]", err)
	}
}

func TestHasNoMinimumStake(t *testing.T) {
	stakeMonitor := local.NewStakeMonitor(minimumStake)
	policy := &minimumStakePolicy{
		stakeMonitor:        stakeMonitor,
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	_, remotePeerPublicKey, err := key.GenerateStaticNetworkKey()
	if err != nil {
		t.Fatal(err)
	}

	if err := policy.Validate(
		key.NetworkKeyToECDSAKey(remotePeerPublicKey),
	); err != errNoMinimumStake {
		t.Fatalf(
			"unexpected validation error\nactual:   [%v]\nexpected: [%v]",
			err,
			errNoMinimumStake,
		)
	}
}

func TestCachesActiveKeepMembers(t *testing.T) {
	stakeMonitor := local.NewStakeMonitor(minimumStake)
	policy := &minimumStakePolicy{
		stakeMonitor:        stakeMonitor,
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	_, remotePeerPublicKey, err := key.GenerateStaticNetworkKey()
	if err != nil {
		t.Fatal(err)
	}
	remotePeerAddress := key.NetworkPubKeyToEthAddress(remotePeerPublicKey)
	stakeMonitor.StakeTokens(remotePeerAddress)

	if err := policy.Validate(
		key.NetworkKeyToECDSAKey(remotePeerPublicKey),
	); err != nil {
		t.Fatalf("validation should pass: [%v]", err)
	}

	stakeMonitor.UnstakeTokens(remotePeerAddress)

	// still caching the old result
	if err := policy.Validate(
		key.NetworkKeyToECDSAKey(remotePeerPublicKey),
	); err != nil {
		t.Fatalf("validation should pass: [%v]", err)
	}

	time.Sleep(time.Second)

	// no longer caches the previous result
	if err := policy.Validate(
		key.NetworkKeyToECDSAKey(remotePeerPublicKey),
	); err != errNoMinimumStake {
		t.Fatalf(
			"unexpected validation error\nactual:   [%v]\nexpected: [%v]",
			err,
			errNoMinimumStake,
		)
	}
}

func TestCachesInactiveKeepMembers(t *testing.T) {
	stakeMonitor := local.NewStakeMonitor(minimumStake)
	policy := &minimumStakePolicy{
		stakeMonitor:        stakeMonitor,
		positiveResultCache: cache.NewTimeCache(cachingPeriod),
		negativeResultCache: cache.NewTimeCache(cachingPeriod),
	}

	_, remotePeerPublicKey, err := key.GenerateStaticNetworkKey()
	if err != nil {
		t.Fatal(err)
	}
	remotePeerAddress := key.NetworkPubKeyToEthAddress(remotePeerPublicKey)

	if err := policy.Validate(
		key.NetworkKeyToECDSAKey(remotePeerPublicKey),
	); err != errNoMinimumStake {
		t.Fatalf(
			"unexpected validation error\nactual:   [%v]\nexpected: [%v]",
			err,
			errNoMinimumStake,
		)
	}

	stakeMonitor.StakeTokens(remotePeerAddress)

	// still caching the old result
	if err := policy.Validate(
		key.NetworkKeyToECDSAKey(remotePeerPublicKey),
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
		key.NetworkKeyToECDSAKey(remotePeerPublicKey),
	); err != nil {
		t.Fatalf("validation should pass: [%v]", err)
	}
}
