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
	"github.com/keep-network/keep-core/pkg/tbtc/maintainer"
)

// TODO: Unit tests for `tbtc.go`.

var logger = log.Logger("keep-tbtc")

// ProtocolName denotes the name of the protocol defined by this package.
const ProtocolName = "tbtc"

// GroupParameters is a structure grouping TBTC group parameters.
type GroupParameters struct {
	// GroupSize is the target size of a group in TBTC.
	GroupSize int
	// GroupQuorum is the minimum number of active participants behaving
	// according to the protocol needed to generate a group in TBTC. This value
	// is smaller than the GroupSize and bigger than the HonestThreshold.
	GroupQuorum int
	// HonestThreshold is the minimum number of active participants behaving
	// according to the protocol needed to generate a signature.
	HonestThreshold int
}

// DishonestThreshold is the maximum number of misbehaving participants for
// which it is still possible to generate a signature. Misbehaviour is any
// misconduct to the protocol, including inactivity.
func (gp *GroupParameters) DishonestThreshold() int {
	return gp.GroupSize - gp.HonestThreshold
}

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

	// Maintainer defines the configuration for the tBTC maintainer.
	Maintainer maintainer.Config
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
	groupParameters := &GroupParameters{
		GroupSize:       100,
		GroupQuorum:     90,
		HonestThreshold: 51,
	}

	node, err := newNode(
		groupParameters,
		chain,
		netProvider,
		keyStorePersistence,
		workPersistence,
		scheduler,
		config,
	)
	if err != nil {
		return fmt.Errorf("cannot set up TBTC node: [%v]", err)
	}

	deduplicator := newDeduplicator()

	if clientInfo != nil {
		// only if client info endpoint is configured
		clientInfo.ObserveApplicationSource(
			"tbtc",
			map[string]clientinfo.Source{
				"pre_params_count": func() float64 {
					return float64(node.dkgExecutor.preParamsCount())
				},
			},
		)
	}

	err = sortition.MonitorPool(
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

	_ = chain.OnDKGResultSubmitted(func(event *DKGResultSubmittedEvent) {
		go func() {
			if ok := deduplicator.notifyDKGResultSubmitted(
				event.Seed,
				event.ResultHash,
				event.BlockNumber,
			); !ok {
				logger.Warnf(
					"Result with hash [0x%x] for DKG with seed [0x%x] "+
						"and starting block [%v] has been already processed",
					event.ResultHash,
					event.Seed,
					event.BlockNumber,
				)
				return
			}

			logger.Infof(
				"Result with hash [0x%x] for DKG with seed [0x%x] "+
					"submitted at block [%v]",
				event.ResultHash,
				event.Seed,
				event.BlockNumber,
			)

			node.validateDKG(
				event.Seed,
				event.BlockNumber,
				event.Result,
				event.ResultHash,
			)
		}()
	})

	_ = chain.OnHeartbeatRequested(func(event *HeartbeatRequestedEvent) {
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
				"heartbeat [%s] requested from "+
					"wallet [0x%x] at block [%v]",
				strings.Join(messagesDigests, ", "),
				event.WalletPublicKey,
				event.BlockNumber,
			)

			executor, ok, err := node.getSigningExecutor(
				unmarshalPublicKey(event.WalletPublicKey),
			)
			if err != nil {
				logger.Errorf("cannot get signing executor: [%v]", err)
				return
			}
			if !ok {
				logger.Infof(
					"node does not control signers of wallet "+
						"with public key [0x%x]",
					event.WalletPublicKey,
				)
				return
			}

			signatures, err := executor.signBatch(
				context.TODO(),
				event.Messages,
				event.BlockNumber,
			)
			if err != nil {
				logger.Errorf("cannot sign batch: [%v]", err)
				return
			}

			logger.Infof(
				"generated [%v] signatures for heartbeat [%s] as "+
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
	paramsInPool := eppip.node.dkgExecutor.preParamsCount()
	poolSize := eppip.config.PreParamsPoolSize
	return paramsInPool >= poolSize
}
