package cmd

import (
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	netlocal "github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/thresholdgroup"
	"github.com/urfave/cli"
)

// SmokeTest performs a simulated distributed key generation and verifyies that the members can do a threshold signature
func SmokeTest(c *cli.Context) error {

	local.GroupSize = c.Int("group-size")
	local.Threshold = c.Int("threshold")

	header(fmt.Sprintf("Smoke test for DKG - GroupSize (%d), Threshold (%d)", local.GroupSize, local.Threshold))
	chainHandle := local.Connect()
	chainCounter := chainHandle.BlockCounter()

	_ = pb.GossipMessage{}

	beaconConfig := chainHandle.RandomBeacon().GetConfig()

	memberChannel := make(chan *thresholdgroup.Member)
	for i := 0; i < beaconConfig.GroupSize; i++ {
		channel := netlocal.Channel("test")
		dkg.Init(channel)

		go func(i int) {
			member, err := dkg.ExecuteDKG(chainCounter, channel, beaconConfig.GroupSize, beaconConfig.Threshold)
			if err != nil {
				panic(fmt.Sprintf("Failed to run DKG [%v].", err))
			}

			chainHandle.ThresholdRelay().OnGroupPublicKeySubmitted(
				func(groupID string, activationBlock *big.Int) {
					if groupID == "test" {
						memberChannel <- member
					}
				})
			chainHandle.ThresholdRelay().OnGroupPublicKeySubmissionFailed(
				func(groupID string, errorMsg string) {
					if groupID == "test" {
						fmt.Fprintf(
							os.Stderr,
							"[member:%s] Failed to submit group public key: [%s]\n",
							member.BlsID.GetHexString(),
							err,
						)
						memberChannel <- nil
					}
				})

			err = chainHandle.ThresholdRelay().SubmitGroupPublicKey(
				"test",
				member.GroupPublicKeyBytes(),
			)
		}(i)
	}

	seenMembers := make(map[*bls.ID]*thresholdgroup.Member)
	for member := range memberChannel {
		if _, alreadySeen := seenMembers[&member.BlsID]; !alreadySeen {
			seenMembers[&member.BlsID] = member
		}

		if len(seenMembers) == beaconConfig.GroupSize {
			break
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

func header(header string) {
	dashes := strings.Repeat("-", len(header))
	fmt.Printf("\n%s\n%s\n%s\n", dashes, header, dashes)
}
