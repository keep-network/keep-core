package cmd

import (
	"fmt"
	"math/big"
	"os"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	netlocal "github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/thresholdgroup"
	"github.com/urfave/cli"
)

// SmokeTest simulates a DKG with a GroupSize of 10 and Threshold of 4
func SmokeTest(c *cli.Context) {

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
		fmt.Printf(
			"[member:%v] Did we get it? %v\n",
			member.BlsID.GetHexString(),
			member.VerifySignature(shares, message))
	}

}
