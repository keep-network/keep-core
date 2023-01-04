package firewall

import (
	"fmt"
	"time"

	"github.com/keep-network/keep-common/pkg/cache"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/operator"
)

// Application defines functionalities for operator verification in the firewall.
type Application interface {
	// IsRecognized returns true if the application recognizes the operator
	// as one participating in the application.
	IsRecognized(operatorPublicKey *operator.PublicKey) (bool, error)
}

// Disabled is an empty Firewall implementation enforcing no rules
// on the connection.
var Disabled = &noFirewall{}

type noFirewall struct{}

func (nf *noFirewall) Validate(remotePeerPublicKey *operator.PublicKey) error {
	return nil
}

// AllowList represents a list of operator public keys that are not checked
// against the firewall rules and are always valid peers.
type AllowList struct {
	allowedPublicKeys map[string]bool
}

// NewAllowList creates a new firewall's allowlist based on the given public
// key list.
func NewAllowList(operatorPublicKeys []*operator.PublicKey) *AllowList {
	allowedPublicKeys := make(map[string]bool, len(operatorPublicKeys))

	for _, operatorPublicKey := range operatorPublicKeys {
		allowedPublicKeys[operatorPublicKey.String()] = true
	}

	return &AllowList{allowedPublicKeys}
}

func (al *AllowList) Contains(operatorPublicKey *operator.PublicKey) bool {
	return al.allowedPublicKeys[operatorPublicKey.String()]
}

// EmptyAllowList represents an empty firewall allowlist.
var EmptyAllowList = NewAllowList([]*operator.PublicKey{})

const (
	// PositiveIsRecognizedCachePeriod is the time period the cache maintains
	// the positive result of the last `IsRecognized` checks.
	// We use the cache to minimize calls to the on-chain client.
	PositiveIsRecognizedCachePeriod = 12 * time.Hour

	// NegativeIsRecognizedCachePeriod is the time period the cache maintains
	// the negative result of the last `IsRecognized` checks.
	// We use the cache to minimize calls to the on-chain client.
	NegativeIsRecognizedCachePeriod = 1 * time.Hour
)

var errNotRecognized = fmt.Errorf(
	"remote peer has not been recognized by any application",
)

func AnyApplicationPolicy(
	applications []Application,
	allowList *AllowList,
) net.Firewall {
	return &anyApplicationPolicy{
		applications:        applications,
		allowList:           allowList,
		positiveResultCache: cache.NewTimeCache(PositiveIsRecognizedCachePeriod),
		negativeResultCache: cache.NewTimeCache(NegativeIsRecognizedCachePeriod),
	}
}

type anyApplicationPolicy struct {
	applications        []Application
	allowList           *AllowList
	positiveResultCache *cache.TimeCache
	negativeResultCache *cache.TimeCache
}

// Validate checks whether the given operator meets the conditions to join
// the network. The operator can join the network if it is an allowlisted node
// or it is a non-allowlisted node but it is recognized as eligible by any of
// the applications. Nil is returned on a successful validation, error otherwise.
// Due to performance reasons, the results of validations for non-allowlisted
// nodes are stored in a cache for a certain amount of time.
func (aap *anyApplicationPolicy) Validate(
	remotePeerPublicKey *operator.PublicKey,
) error {
	// If the peer is on the allowlist, consider it validated.
	if aap.allowList.Contains(remotePeerPublicKey) {
		return nil
	}

	// First, check in the in-memory time caches to minimize hits to the ETH client.
	// If the Keep client with the given chain address is in the positive result
	// cache it means it has been recognized when the last `IsRecognized` was
	// executed and caching period has not elapsed yet. Similarly, if the client
	// is in the negative result cache it means it hasn't been recognized.
	//
	// If the caching period elapsed, cache checks will return false and we
	// have to ask the chain about the current status.
	aap.positiveResultCache.Sweep()
	aap.negativeResultCache.Sweep()

	remotePeerPublicKeyHex := remotePeerPublicKey.String()

	if aap.positiveResultCache.Has(remotePeerPublicKeyHex) {
		return nil
	}

	if aap.negativeResultCache.Has(remotePeerPublicKeyHex) {
		return errNotRecognized
	}

	validationSuccessful := false
	for _, application := range aap.applications {
		isRecognized, err := application.IsRecognized(remotePeerPublicKey)
		if err != nil {
			return fmt.Errorf(
				"could not validate if remote peer is recognized by application: [%w]",
				err,
			)
		}
		if isRecognized {
			validationSuccessful = true
			break
		}
	}

	if !validationSuccessful {
		// Add this address to the negative result cache.
		// `IsRecognized` will not be called again for the entire caching period.
		aap.negativeResultCache.Add(remotePeerPublicKeyHex)
		return errNotRecognized
	}

	// Add this address to the positive result cache.
	// `IsRecognized` will not be called again for the entire caching period.
	aap.positiveResultCache.Add(remotePeerPublicKeyHex)

	return nil
}
