package relay

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	chainLocal "github.com/keep-network/keep-core/pkg/chain/local"
)

var address = "0x65ea55c1f10491038425725dc00dffeab2a1e28a"
var relayEntryTimeout = uint64(15)

func TestMonitorRelayEntryOnChain_WhenResultWasSubmitted(t *testing.T) {
	chain := chainLocal.Connect(5, 3, big.NewInt(200))
	blockCounter, err := chain.BlockCounter()
	if err != nil {
		fmt.Printf("failed to setup a block counter")
	}

	stakeMonitor := chainLocal.NewStakeMonitor(big.NewInt(200))
	staker, err := stakeMonitor.StakerFor(address)
	if err != nil {
		fmt.Printf("failed to stake an address")
	}
	node := NewNode(
		staker,
		nil,
		blockCounter,
		nil,
		nil,
	)

	relayChain := chain.ThresholdRelay()
	startBlockHeight, _ := blockCounter.CurrentBlock()
	chainConfig := &config.Chain{
		RelayEntryTimeout: uint64(relayEntryTimeout),
	}
	go node.MonitorRelayEntryOnChain(
		relayChain,
		uint64(startBlockHeight),
		chainConfig,
	)

	backwardBlocks := uint64(5)
	relayEntryResultWindow := startBlockHeight + relayEntryTimeout - backwardBlocks
	err = blockCounter.WaitForBlockHeight(relayEntryResultWindow)
	if err != nil {
		fmt.Printf("failed to wait for a block: [%v]", relayEntryResultWindow)
	}

	relayEntryPromise := chain.ThresholdRelay().SubmitRelayEntry(big.NewInt(1))

	done := make(chan *event.Entry)
	relayEntryPromise.OnSuccess(func(entry *event.Entry) {
		done <- entry
	}).OnFailure(func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	})

	blockCounter.WaitForBlockHeight(startBlockHeight + relayEntryTimeout)

	count := chain.ReportRelayEntryTimeoutCount()

	if count != 0 {
		t.Fatalf(
			"\nexpected: [%v]\nactual:   [%v]",
			0,
			count,
		)
	}
}

func TestMonitorRelayEntryOnChain_WhenResultWasNotDelivered(t *testing.T) {
	chain := chainLocal.Connect(5, 3, big.NewInt(200))
	blockCounter, err := chain.BlockCounter()
	if err != nil {
		fmt.Printf("failed to setup a block counter")
	}

	stakeMonitor := chainLocal.NewStakeMonitor(big.NewInt(200))
	staker, err := stakeMonitor.StakerFor(address)
	if err != nil {
		fmt.Printf("failed to stake an address")
	}
	node := NewNode(
		staker,
		nil,
		blockCounter,
		nil,
		nil,
	)

	relayChain := chain.ThresholdRelay()
	currentBlock, _ := blockCounter.CurrentBlock()
	chainConfig := &config.Chain{
		RelayEntryTimeout: uint64(relayEntryTimeout),
	}
	go node.MonitorRelayEntryOnChain(
		relayChain,
		uint64(currentBlock),
		chainConfig,
	)

	relayEntryTimeoutExtended := relayEntryTimeout + 5
	blockCounter.WaitForBlockHeight(currentBlock + relayEntryTimeoutExtended)

	count := chain.ReportRelayEntryTimeoutCount()

	if count != 1 {
		t.Fatalf(
			"Reported timeout count does not match\nexpected: [%v]\nactual:   [%v]",
			1,
			count,
		)
	}
}
