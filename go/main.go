package main

import (
	"fmt"
	"os"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/go/beacon/broadcast"
	"github.com/keep-network/keep-core/go/beacon/chain"
	"github.com/keep-network/keep-core/go/beacon/relay"
	"github.com/keep-network/keep-core/go/thresholdgroup"
)

func main() {
	bls.Init(bls.CurveFp382_1)

	beaconConfig := chain.GetBeaconConfig()

	channel := broadcast.LocalChannel("test")
	chainCounter := chain.LocalBlockCounter()

	members := make([]*thresholdgroup.Member, 0, beaconConfig.GroupSize)
	memberChannel := make(chan *thresholdgroup.Member)
	for i := 0; i < beaconConfig.GroupSize; i++ {
		go func() {
			member, err := relay.ExecuteDKG(chainCounter, channel, beaconConfig.GroupSize, beaconConfig.Threshold)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to run DKG for member %v: [%s].", i, err)
				memberChannel <- nil
				return
			}

			memberChannel <- member
		}()
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

	fmt.Printf("Members! %v\n", members)
}
