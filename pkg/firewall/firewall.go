package firewall

import (
	"fmt"
	"time"

	"github.com/keep-network/keep-common/pkg/cache"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/operator"
)

type Application interface {
	IsRecognized(operatorPublicKey *operator.PublicKey) (bool, error)
}

// Disabled is an empty Firewall implementation enforcing no rules
// on the connection.
var Disabled = &noFirewall{}

type noFirewall struct{}

func (nf *noFirewall) Validate(remotePeerPublicKey *operator.PublicKey) error {
	return nil
}

const (
	// PositiveMinimumStakeCachePeriod is the time period the cache maintains
	// the positive result of the last `IsRecognized` checks.
	// We use the cache to minimize calls to the on-chain client.
	PositiveMinimumStakeCachePeriod = 12 * time.Hour

	// NegativeMinimumStakeCachePeriod is the time period the cache maintains
	// the negative result of the last `IsRecognized` checks.
	// We use the cache to minimize calls to the on-chain client.
	NegativeMinimumStakeCachePeriod = 1 * time.Hour
)

var errNotRecognized = fmt.Errorf(
	"remote peer has not been recognized by any application",
)

func AnyApplicationPolicy(applications []Application) net.Firewall {
	return &anyApplicationPolicy{
		applications:        applications,
		positiveResultCache: cache.NewTimeCache(PositiveMinimumStakeCachePeriod),
		negativeResultCache: cache.NewTimeCache(NegativeMinimumStakeCachePeriod),
	}
}

type anyApplicationPolicy struct {
	applications        []Application
	positiveResultCache *cache.TimeCache
	negativeResultCache *cache.TimeCache
}

func (msp *anyApplicationPolicy) Validate(
	remotePeerPublicKey *operator.PublicKey,
) error {
	remotePeerPublicKeyHex := remotePeerPublicKey.String()

	// First, check in the in-memory time caches to minimize hits to the on client.
	// If the Keep client with the given chain address is in the positive result
	// cache it means it has been recognized when the last `IsRecognized` was
	// executed and caching period has not elapsed yet. Similarly, if the client
	// is in the negative result cache it means it hasn't been recognized.
	//
	// If the caching period elapsed, cache checks will return false and we
	// have to ask the chain about the current status.
	msp.positiveResultCache.Sweep()
	msp.negativeResultCache.Sweep()

	if msp.positiveResultCache.Has(remotePeerPublicKeyHex) {
		return nil
	}

	if msp.negativeResultCache.Has(remotePeerPublicKeyHex) {
		return errNotRecognized
	}

	validationSuccessful := false
	for _, application := range msp.applications {
		isRecognized, err := application.IsRecognized(remotePeerPublicKey)
		if err == nil && isRecognized {
			validationSuccessful = true
			break
		}
	}

	if !validationSuccessful {
		// Add this address to the negative result cache.
		// `IsRecognized` will not be called again for the entire caching period.
		msp.negativeResultCache.Add(remotePeerPublicKeyHex)
		return errNotRecognized
	}

	// Add this address to the positive result cache.
	// `IsRecognized` will not be called again for the entire caching period.
	msp.positiveResultCache.Add(remotePeerPublicKeyHex)

	return nil
}
