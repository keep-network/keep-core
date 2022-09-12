package tbtc

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/ipfs/go-log"
	commonDiagnostics "github.com/keep-network/keep-common/pkg/diagnostics"
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/diagnostics"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/sortition"
)

// TODO: Unit tests for `tbtc.go`.

var logger = log.Logger("keep-tbtc")

// ProtocolName denotes the name of the protocol defined by this package.
const ProtocolName = "tbtc"

const (
	DefaultPreParamsPoolSize              = 3000
	DefaultPreParamsGenerationTimeout     = 2 * time.Minute
	DefaultPreParamsGenerationDelay       = 10 * time.Second
	DefaultPreParamsGenerationConcurrency = 1
)

var DefaultKeyGenerationConcurrency = runtime.GOMAXPROCS(0)

// Config carries the config for tBTC protocol.
type Config struct {
	// The size of the pre-parameters pool for tECDSA.
	PreParamsPoolSize int
	// Timeout for pre-parameters generation for tECDSA.
	PreParamsGenerationTimeout time.Duration
	// The delay between generating new pre-params for tECDSA.
	PreParamsGenerationDelay time.Duration
	// Concurrency level for pre-parameters generation for tECDSA.
	PreParamsGenerationConcurrency int
	// Concurrency level for key-generation for tECDSA.
	KeyGenerationConcurrency int
}

// Initialize kicks off the TBTC by initializing internal state, ensuring
// preconditions like staking are met, and then kicking off the internal TBTC
// implementation. Returns an error if this failed.
func Initialize(
	ctx context.Context,
	chain Chain,
	netProvider net.Provider,
	persistence persistence.Handle,
	scheduler *generator.Scheduler,
	config Config,
	registry *commonDiagnostics.Registry,
) error {
	node := newNode(chain, netProvider, persistence, scheduler, config)
	deduplicator := newDeduplicator()

	assembleTbtcDiagnostics := func() map[string]interface{} {
		return map[string]interface{}{
			"preParamsPoolSize": node.dkgExecutor.PreParamsPool().CurrentSize(),
		}
	}
	diagnostics.RegisterApplicationSource(registry, "tbtc", assembleTbtcDiagnostics)

	err := sortition.MonitorPool(ctx, logger, chain, sortition.DefaultStatusCheckTick)
	if err != nil {
		return fmt.Errorf(
			"could not set up sortition pool monitoring: [%v]",
			err,
		)
	}

	_ = chain.OnDKGStarted(func(event *DKGStartedEvent) {
		go func() {
			if ok := deduplicator.notifyDKGStarted(
				event.Seed,
			); !ok {
				logger.Warningf(
					"DKG started event with seed [0x%x] and "+
						"starting block [%v] has been already processed",
					event.Seed,
					event.BlockNumber,
				)
				return
			}

			logger.Infof(
				"DKG started with seed [0x%x] at block [%v]",
				event.Seed,
				event.BlockNumber,
			)

			node.joinDKGIfEligible(
				event.Seed,
				event.BlockNumber,
			)
		}()
	})

	// TODO: This is a temporary signing loop trigger that should be removed
	//       once the client is integrated with real on-chain contracts.
	_ = chain.OnSignatureRequested(func(event *SignatureRequestedEvent) {
		go func() {
			// There is no need to deduplicate. Test loop events are unique.

			logger.Infof(
				"signature of message [%v] requested from "+
					"wallet [0x%x] at block [%v]",
				event.Message,
				event.WalletPublicKey,
				event.BlockNumber,
			)

			node.joinSigningIfEligible(
				event.Message,
				unmarshalPublicKey(event.WalletPublicKey),
				event.BlockNumber,
			)
		}()
	})

	return nil
}
