package cmd

import (
	"fmt"
	"os"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	netlocal "github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/thresholdgroup"
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

// SmokeTest performs a simulated distributed key generation and verifyies that the members can do a threshold signature
func SmokeTest(c *cli.Context) error {
	groupSize := c.Int(groupSizeFlag)
	threshold := c.Int(thresholdFlag)

	chainHandle := local.Connect(groupSize, threshold)
	chainCounter, err := chainHandle.BlockCounter()
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to run setup chainHandle.BlockCounter: [%v].",
			err,
		))
	}

	_ = pb.Envelope{}

	beaconConfig, err := chainHandle.ThresholdRelay().GetConfig()
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to run get configuration: [%v].",
			err,
		))
	}

	memberChannel := make(chan *thresholdgroup.Member)
	for i := 0; i < beaconConfig.GroupSize; i++ {
		channel := netlocal.Channel("test")
		dkg.Init(channel)

		go func(i int) {
			member, err := dkg.ExecuteDKG(
				i+1,
				chainCounter,
				channel,
				beaconConfig.GroupSize,
				beaconConfig.Threshold,
			)
			if err != nil {
				panic(fmt.Sprintf("Failed to run DKG [%v].", err))
			}

			chainHandle.ThresholdRelay().SubmitGroupPublicKey(
				"test",
				member.GroupPublicKeyBytes(),
			).OnSuccess(func(data *event.GroupRegistration) {
				if string(data.GroupPublicKey) == "test" {
					memberChannel <- member
				} else {
					fmt.Fprintf(
						os.Stderr,
						"[member:%s] incorrect data, expected 'test' got '%s', activation block: %s\n",
						member.BlsID.GetHexString(),
						string(data.GroupPublicKey),
						data.ActivationBlockHeight,
					)
					memberChannel <- nil
				}
			}).OnFailure(func(err error) {
				fmt.Fprintf(
					os.Stderr,
					"[member:%s] Failed to submit group public key: [%s]\n",
					member.BlsID.GetHexString(),
					err,
				)
				memberChannel <- nil
			})

		}(i)
	}

	seenMembers := make(map[*bls.ID]*thresholdgroup.Member)
	for member := range memberChannel {
		if member != nil {
			if _, alreadySeen := seenMembers[&member.BlsID]; !alreadySeen {
				seenMembers[&member.BlsID] = member
			}

			if len(seenMembers) == beaconConfig.GroupSize {
				break
			}
		} else {
			fmt.Printf("nil member\n")
		}
	}

	if len(seenMembers) < beaconConfig.GroupSize {
		panic("Failed to reach group size during DKG, aborting.")
	}

	message := "This is a message!"
	shares := make(map[bls.ID][]byte, 0)
	for _, member := range seenMembers {
		shares[member.BlsID] = member.SignatureShare(message)
	}

	for _, member := range seenMembers {
		validSignature, err := member.VerifySignature(shares, message)
		if err != nil {
			fmt.Printf(
				"[member:0x%010s] Error verifying signature: [%v].\n",
				member.BlsID.GetHexString(),
				err,
			)
		}

		fmt.Printf(
			"[member:0x%010s] Did we get it? %v\n",
			member.BlsID.GetHexString(),
			validSignature,
		)
	}

	return nil
}
