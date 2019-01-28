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
	"github.com/urfave/cli"
)

const (
	defaultGroupSize    int = 5
	defaultThreshold    int = 3
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

const smokeTestDescription = `The smoke-test command creates a local threshold group of the
   specified size and with the specified threshold and simulates a
   distributed key generation process with an in-process broadcast
   channel and chain implementation. Once the process is complete,
   a threshold signature is executed, once again with an in-process
   broadcast channel and chain, and the final signature is verified
   by each member of the group.`

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
		createNode(context, chainHandle, groupSize, threshold)
	}

	// Give the nodes a sec to get going.
	<-time.NewTimer(time.Second).C

	chainHandle.ThresholdRelay().RequestRelayEntry(
		&big.Int{},
		big.NewInt(12),
	)

	fmt.Println("Wait until initial group is ready...")
	blockCounter, err := chainHandle.BlockCounter()
	if err != nil {
		panic(err)
	}
	err = blockCounter.WaitForBlocks(45)
	if err != nil {
		panic(err)
	}

	chainHandle.ThresholdRelay().SubmitRelayEntry(&event.Entry{
		RequestID:     big.NewInt(0),
		Value:         big.NewInt(1111),
		Seed:          big.NewInt(13),
		PreviousValue: big.NewInt(101),
	})

	select {
	case <-context.Done():
		fmt.Println("All done!")
		return context.Err()
	}
}

func createNode(
	context context.Context,
	chainHandle chain.Handle,
	groupSize int,
	threshold int,
) {
	toEthereumAddress := func(id string) string {
		return common.BytesToAddress([]byte(id)).Hex()
	}

	chainCounter, err := chainHandle.BlockCounter()
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to run setup chainHandle.BlockCounter: [%v].",
			err,
		))
	}

	stakeMonitor, err := chainHandle.StakeMonitor()
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to run setup chainHandle.StakeMonitor: [%v].",
			err,
		))
	}

	netProvider := netlocal.Connect()

	go beacon.Initialize(
		context,
		toEthereumAddress(netProvider.ID().String()),
		chainHandle.ThresholdRelay(),
		chainCounter,
		stakeMonitor,
		netProvider,
	)
}
