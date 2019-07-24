package cmd

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
	netlocal "github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/persistence"
	"github.com/urfave/cli"
)

const (
	defaultGroupSize    int = 10
	defaultThreshold    int = 4
	defaultMinimumStake int = 2000000
)

// SmokeTestCommand contains the definition of the smoke-test command-line
// subcommand.
var SmokeTestCommand cli.Command

const (
	groupSizeFlag     = "group-size"
	groupSizeShort    = "g"
	thresholdFlag     = "threshold"
	thresholdShort    = "t"
	minimumStakeFlag  = "minimum-stake"
	minimumStakeShort = "s"
)

const smokeTestDescription = `The smoke-test command creates a local threshold
   group of the specified size and with the specified threshold and simulates a
   distributed key generation process with an in-process broadcast channel and
   chain implementation. Once the process is complete, a threshold signature is
   executed, once again with an in-process broadcast channel and chain, and the
   final signature is verified by each member of the group.`

type noopPersistence struct {
}

func (np *noopPersistence) Save(data []byte, directory string, name string) error {
	// noop
	return nil
}

func (np *noopPersistence) ReadAll() (<-chan persistence.DataDescriptor, <-chan error) {
	dataChannel := make(chan persistence.DataDescriptor)
	errorChannel := make(chan error)

	close(dataChannel)
	close(errorChannel)

	return dataChannel, errorChannel
}

func (np *noopPersistence) Archive(directory string) error {
	// noop
	return nil
}

func init() {
	SmokeTestCommand = cli.Command{
		Name:        "smoke-test",
		Usage:       "Simulates Distributed Key Generation (DKG) and signature generation locally",
		Description: smokeTestDescription,
		Action:      SmokeTest,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  groupSizeFlag + "," + groupSizeShort,
				Value: defaultGroupSize,
			},
			&cli.IntFlag{
				Name:  thresholdFlag + "," + thresholdShort,
				Value: defaultThreshold,
			},
			&cli.IntFlag{
				Name:  minimumStakeFlag + "," + minimumStakeShort,
				Value: defaultMinimumStake,
			},
		},
	}
}

// SmokeTest sets up a set of local virtual nodes and launches the beacon on
// them, simulating some relay entries and requests.
func SmokeTest(c *cli.Context) error {
	groupSize := c.Int(groupSizeFlag)
	threshold := c.Int(thresholdFlag)
	minimumStake := c.Int(minimumStakeFlag)

	chainHandle := local.Connect(
		groupSize,
		threshold,
		big.NewInt(int64(minimumStake)),
	)

	context := context.Background()

	for i := 0; i < groupSize; i++ {
		operatorPrivateKey, _, err := operator.GenerateKeyPair()
		if err != nil {
			panic("failed to generate private key")
		}

		createNode(
			context, operatorPrivateKey, chainHandle, groupSize, threshold,
		)
	}

	// Give the nodes a sec to get going.
	<-time.NewTimer(time.Second).C

	chainHandle.ThresholdRelay().SubmitRelayEntry(&event.Entry{
		SigningId:     big.NewInt(0),
		Value:         big.NewInt(0),
		GroupPubKey:   big.NewInt(0).Bytes(),
		Seed:          big.NewInt(0),
		PreviousEntry: &big.Int{},
	})

	// TODO Add validations when DKG Phase 14 is implemented.

	select {
	case <-context.Done():
		fmt.Println("All done!")
		return context.Err()
	}
}

func createNode(
	context context.Context,
	operatorPrivateKey *operator.PrivateKey,
	chainHandle chain.Handle,
	groupSize int,
	threshold int,
) {
	toEthereumAddress := func(value string) string {
		return common.BytesToAddress(
			[]byte(value),
		).String()
	}

	stakeMonitor, err := chainHandle.StakeMonitor()
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to run setup chainHandle.StakeMonitor: [%v].",
			err,
		))
	}

	storage := &noopPersistence{}

	netProvider := netlocal.Connect()

	go func() {
		// Generate staker's ID. It needs to be a properly formatter ethereum
		// address. Address can be created from any string.
		stakingID := toEthereumAddress(netProvider.ID().String())

		localMonitor := stakeMonitor.(*local.StakeMonitor)
		localMonitor.StakeTokens(stakingID)

		err := beacon.Initialize(
			context,
			stakingID,
			chainHandle,
			netProvider,
			storage,
		)
		if err != nil {
			panic(fmt.Sprintf(
				"Failed to run beacon.Initialize: [%v].",
				err,
			))
		}
	}()
}
