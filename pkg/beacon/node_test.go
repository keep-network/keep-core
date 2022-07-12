package beacon

import (
	"fmt"
	"math/big"
	"testing"

	chainLocal "github.com/keep-network/keep-core/pkg/chain/local"
)

var address = "0x65ea55c1f10491038425725dc00dffeab2a1e28a"
var relayEntryTimeout = uint64(15)

func TestMonitorRelayEntryOnChain_EntrySubmitted(t *testing.T) {
	chain := chainLocal.Connect(5, 3, big.NewInt(200))

	node := &node{
		beaconChain: chain,
	}

	blockCounter, err := node.beaconChain.BlockCounter()
	if err != nil {
		fmt.Printf("failed to setup a block counter: [%v]", err)
	}

	startBlockHeight, err := blockCounter.CurrentBlock()
	if err != nil {
		t.Fatal(err)
	}

	go node.MonitorRelayEntry(startBlockHeight)

	// the window to get a relay entry is from currentBlock to (currentBlock+relayEntryTimeout)
	// we subtract arbitarly 5 blocks to be within this window. Ex. 0 + 15 - 5
	relayEntrySubmissionWindow := startBlockHeight + relayEntryTimeout - 5
	err = blockCounter.WaitForBlockHeight(relayEntrySubmissionWindow)
	if err != nil {
		fmt.Printf(
			"failed to wait for a block: [%v]: [%v]",
			relayEntrySubmissionWindow,
			err,
		)
	}

	chain.ThresholdRelay().SubmitRelayEntry(big.NewInt(1).Bytes()).
		OnFailure(func(err error) {
			if err != nil {
				t.Fatal(err)
			}
		})

	blockCounter.WaitForBlockHeight(startBlockHeight + relayEntryTimeout)

	timeoutsReport := chain.GetRelayEntryTimeoutReports()
	numberOfReports := len(timeoutsReport)

	if numberOfReports != 0 {
		t.Fatalf(
			"expected 0 relay entry timeout reports; has: [%v]",
			numberOfReports,
		)
	}
}

func TestMonitorRelayEntryOnChain_EntryNotSubmitted(t *testing.T) {
	chain := chainLocal.Connect(5, 3, big.NewInt(200))

	node := &node{
		beaconChain: chain,
	}

	blockCounter, err := node.beaconChain.BlockCounter()
	if err != nil {
		fmt.Printf("failed to setup a block counter: [%v]", err)
	}

	startBlockHeight, err := blockCounter.CurrentBlock()
	if err != nil {
		t.Fatal(err)
	}

	go node.MonitorRelayEntry(startBlockHeight)

	relayEntryTimeoutFromStart := startBlockHeight + relayEntryTimeout

	// we want to exceed the relay entry timeout to report that a relay entry
	// was not submitted. 5 is an arbitrary number to exceed relayEntryTimeout.
	err = blockCounter.WaitForBlockHeight(relayEntryTimeoutFromStart + 5)
	if err != nil {
		t.Fatal(err)
	}

	timeoutsReport := chain.GetRelayEntryTimeoutReports()
	numberOfReports := len(timeoutsReport)

	if numberOfReports != 1 {
		t.Fatalf(
			"Number of timeout reports does not match\nexpected: [%v]\nactual:   [%v]",
			1,
			numberOfReports,
		)
	}

	if timeoutsReport[0] != relayEntryTimeoutFromStart {
		t.Fatalf(
			"Timeout reporting must happen only after a relay entry timeout\nexpected: [%v]\nactual:   [%v]",
			relayEntryTimeoutFromStart,
			timeoutsReport[0],
		)
	}
}
