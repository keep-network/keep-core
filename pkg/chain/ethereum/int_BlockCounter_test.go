package ethereum

import (
	"fmt"
	"testing"
	"time"
)

func TestBlockCounter(t *testing.T) {

	return // PJS - quick return

	ebc, err := EthConn.BlockCounter()
	if err != nil {
		t.Errorf("\nUnable to setup BlockCounter()\n")
		return
	}
	BlockNo := ebc.GetBlockNo()
	fmt.Printf("Sleeping 60 seconds - waiting for blocks to occur on chain\n")
	time.Sleep(60 * time.Second)
	nebc, err := EthConn.BlockCounter() // (chain.BlockCounter, error) {
	if err != nil {
		t.Errorf("\nUnable to setup BlockCounter()\n")
		return
	}
	NewBlockNo := nebc.GetBlockNo()

	if BlockNo >= NewBlockNo {
		t.Errorf(
			"\nexpected: [%v] to be larger than\nprevious: [%v]\n",
			BlockNo,
			NewBlockNo,
		)
	}
}
