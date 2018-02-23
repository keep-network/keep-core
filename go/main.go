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
	memberChannel := make(chan *thresholdgroup.Member, beaconConfig.GroupSize)
	type empty struct{}
	groupSemaphore := make(chan empty, beaconConfig.GroupSize);
	for i := 0; i < beaconConfig.GroupSize; i++ {
		go func(i int) {
			member, err := relay.ExecuteDKG(chainCounter, channel, beaconConfig.GroupSize, beaconConfig.Threshold)
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
			groupSemaphore <- empty{};
		}(i)
	}
	// Wait for goroutines to finish
	for i := 0; i < beaconConfig.GroupSize; i++ { <-groupSemaphore }
	/// Populate members list
	for member := range memberChannel {
		members = append(members, member)
		if len(members) == beaconConfig.GroupSize { close(memberChannel) }
	}

	if len(members) < beaconConfig.GroupSize {
		fmt.Println("len(members)", len(members), "beaconConfig.GroupSize", beaconConfig.GroupSize)
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
