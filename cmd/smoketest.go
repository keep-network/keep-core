package cmd

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/keep-network/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
	netlocal "github.com/keep-network/keep-core/pkg/net/local"
	"github.com/urfave/cli"
)

const (
	defaultGroupSize int = 10
	defaultThreshold int = 4
)

// SmokeTestCommand contains the definition of the smoke-test command-line
// subcommand.
var SmokeTestCommand cli.Command

const (
	groupSizeFlag  = "group-size"
	groupSizeShort = "g"
	thresholdFlag  = "threshold"
	thresholdShort = "t"
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
		},
	}
}

// SmokeTest sets up a set of local virtual nodes and launches the beacon on
// them, simulating some relay entries and requests.
func SmokeTest(c *cli.Context) error {
	groupSize := c.Int(groupSizeFlag)
	threshold := c.Int(thresholdFlag)

	chainHandle := local.Connect(groupSize, threshold)
	context, contextCancel := context.WithTimeout(context.Background(), 30*time.Second)

	for i := 0; i < groupSize; i++ {
		createNode(context, chainHandle, groupSize, threshold)
	}

	// Give the nodes a sec to get going.
	<-time.NewTimer(time.Second).C

	requestID := big.NewInt(135)

	chainHandle.ThresholdRelay().SubmitRelayEntry(&event.Entry{
		RequestID:     requestID,
		Value:         big.NewInt(154),
		GroupID:       big.NewInt(168),
		PreviousEntry: &big.Int{},
		Seed:          big.NewInt(101),
	})

	// TODO Add validations when DKG Phase 14 is implemented.
	chainHandle.ThresholdRelay().OnDKGResultPublished(func(dkgResultPublication *event.DKGResultPublication) {
		if dkgResultPublication.RequestID.Cmp(requestID) != 0 {
			panic(fmt.Sprintf("unexpected request ID for published result\nexpected: %v\nactual:   %v\n",
				requestID,
				dkgResultPublication.RequestID,
			))
		}

		// TODO We can cancel the context after we are sure that all validation passed
		// Need to revisit this part after Phase 14 is implemented. Currently
		// `OnGroupRegistered` is not called so the context is cancelled here.
		contextCancel()
	})

	chainHandle.ThresholdRelay().
		OnGroupRegistered(func(registration *event.GroupRegistration) {
			// Give the nodes a sec to all get registered.
			<-time.NewTimer(time.Second).C
			chainHandle.ThresholdRelay().RequestRelayEntry(&big.Int{}, &big.Int{})
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
