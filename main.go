package main

import (
	"fmt"
	"os"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/chain/gen"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	netlocal "github.com/keep-network/keep-core/pkg/net/local"
	"github.com/keep-network/keep-core/pkg/thresholdgroup"
)

func main() {
	bls.Init(bls.CurveFp382_1)

	chainHandle := local.InitLocal()
	chainCounter := chainHandle.BlockCounter()
	channel := netlocal.Channel("test")

	dkg.Init(channel)

	_ = pb.GossipMessage{}

	beaconConfig := chainHandle.RandomBeacon().GetConfig()

	members := make([]*thresholdgroup.Member, 0, beaconConfig.GroupSize)
	memberChannel := make(chan *thresholdgroup.Member)
	for i := 0; i < beaconConfig.GroupSize; i++ {
		go func(i int) {
			member, err := dkg.ExecuteDKG(chainCounter, channel, beaconConfig.GroupSize, beaconConfig.Threshold)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"[member:%v] Failed to run DKG: [%s] (index %d).",
					member.BlsID.GetHexString(),
					err,
					i)
				memberChannel <- nil
				return
			}

			memberChannel <- member
		}(i)
	}

	seenMembers := 0
	for member := range memberChannel {
		seenMembers++
		if member != nil {
			members = append(members, member)
			if len(members) == beaconConfig.GroupSize {
				break
			}
		}

		if seenMembers == beaconConfig.GroupSize {
			break
		}
	}

	if len(members) < beaconConfig.GroupSize {
		panic("Failed to reach group size during DKG, aborting.")
	}

	message := "This is a message!"
	shares := make(map[bls.ID][]byte, 0)
	for _, member := range members {
		shares[member.BlsID] = member.SignatureShare(message)
	}

	for _, member := range members {
		fmt.Printf(
			"[member:%v] Did we get it? %v\n",
			member.BlsID.GetHexString(),
			member.VerifySignature(shares, message))
	}
}

// verifyCompilationOfGeneraedCode is a compile time check that verifies all of the
// abigen generated code for contracts.
func verifyCompilationOfGeneraedCode() {
	_ = gen.KeepRandomBeaconImplV1{}
	_ = gen.KeepToken{}
	_ = gen.StakingProxy{}
	_ = gen.TokenGrant{}
	_ = gen.TokenStaking{}
}
