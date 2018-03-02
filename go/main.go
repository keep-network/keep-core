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
	gsTestChannel := make(chan error)
	for i := 0; i < len(members); i++ {
		go func(i int) {
			gsTestChannel <- relay.ExecuteGroupSignature(message, chainCounter, channel, members[i])
		}(i)
	}

	// go through the channel so we can see whats going on running relay.ExecuteGroupSignature via goroutine
	i := 0
	for gs := range gsTestChannel {
		if gs == nil {
			fmt.Printf("[member:%v] Exit ExecuteGroupSignature \n", members[i].ID)
		}
		i++
		if i == beaconConfig.GroupSize {
			break
		}
	}

	// Verify group signature
	for _, member := range members {
		fmt.Printf(
			"[member:%v] Verifying signature with all received group signature shares: %v\n",
			member.ID,
			member.VerifySignature(member.GetReceivedGroupSignatureShares(), message))
	}
}
