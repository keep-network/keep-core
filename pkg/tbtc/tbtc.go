package tbtc

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/ipfs/go-log"

	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/clientinfo"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/sortition"
)

// TODO: Unit tests for `tbtc.go`.

var logger = log.Logger("keep-tbtc")

// ProtocolName denotes the name of the protocol defined by this package.
const ProtocolName = "tbtc"

const (
	DefaultPreParamsPoolSize              = 1000
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
	keyStorePersistence persistence.ProtectedHandle,
	workPersistence persistence.BasicHandle,
	scheduler *generator.Scheduler,
	config Config,
	clientInfo *clientinfo.Registry,
) error {
	node := newNode(chain, netProvider, keyStorePersistence, workPersistence, scheduler, config)
	deduplicator := newDeduplicator()

	if clientInfo != nil {
		// only if client info endpoint is configured
		clientInfo.ObserveApplicationSource(
			"tbtc",
			map[string]clientinfo.Source{
				"pre_params_count": func() float64 {
					return float64(node.dkgExecutor.PreParamsCount())
				},
			},
		)
	}

	err := sortition.MonitorPool(
		ctx,
		logger,
		chain,
		sortition.DefaultStatusCheckTick,
		sortition.NewConjunctionPolicy(
			sortition.NewBetaOperatorPolicy(chain, logger),
			&enoughPreParamsInPoolPolicy{
				node:   node,
				config: config,
			},
		),
	)
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
				logger.Warnf(
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
			messagesDigests := make([]string, len(event.Messages))
			for i, message := range event.Messages {
				bytes := message.Bytes()
				messagesDigests[i] = fmt.Sprintf(
					"0x%x...%x",
					bytes[:2],
					bytes[len(bytes)-2:],
				)
			}

			logger.Infof(
				"signature of messages [%s] requested from "+
					"wallet [0x%x] at block [%v]",
				strings.Join(messagesDigests, ", "),
				event.WalletPublicKey,
				event.BlockNumber,
			)

			controller, err := node.createSigningGroupController(
				unmarshalPublicKey(event.WalletPublicKey),
			)
			if err != nil {
				logger.Errorf("cannot get signing group controller: [%v]", err)
				return
			}

			signatures, err := controller.signBatch(
				context.TODO(),
				event.Messages,
				event.BlockNumber,
			)
			if err != nil {
				logger.Errorf("cannot sign batch: [%v]", err)
				return
			}

			logger.Infof(
				"generated [%v] sigantures for messages [%s] as "+
					"requested from wallet [0x%x] at block [%v]",
				len(signatures),
				strings.Join(messagesDigests, ", "),
				event.WalletPublicKey,
				event.BlockNumber,
			)
		}()
	})

	return nil
}

// enoughPreParamsInPoolPolicy is a policy that enforces the sufficient size
// of the DKG pre-parameters pool before joining the sortition pool.
type enoughPreParamsInPoolPolicy struct {
	node   *node
	config Config
}

func (eppip *enoughPreParamsInPoolPolicy) ShouldJoin() bool {
	paramsInPool := eppip.node.dkgExecutor.PreParamsCount()
	poolSize := eppip.config.PreParamsPoolSize
	return paramsInPool >= poolSize
}
