package ethereum_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
)

type CfgType struct {
	GethConnectionString string
}

var (
	gCfg   CfgType
	client *rpc.Client
)

func Test_BlockTest(t *testing.T) {

	// ----------------------------------------------------------------------------------------------------
	// Test code that needs to do a timeout!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// ----------------------------------------------------------------------------------------------------

	// func BlockCounter(client *rpc.Client) chain.BlockCounter {
	countWait := ethereum.BlockCounter(client)
	start := time.Now()
	countWait.WaitForBlocks(1)
	tm := time.Now()
	elapsed := tm.Sub(start)
	if elapsed < 1000000000 {
		t.Errorf("Did not wait\n")
	}

	start = time.Now()
	countWait.WaitForBlocks(2)
	tm = time.Now()
	elapsed = tm.Sub(start)
	if elapsed < 2000000000 {
		t.Errorf("Did not wait\n")
	}

	start = time.Now()
	countWait.WaitForBlocks(0)
	tm = time.Now()
	elapsed = tm.Sub(start)
	if elapsed < 1000 {
		t.Errorf("Did not wait\n")
	}
}

func TestMain(m *testing.M) {

	var buf []byte
	fn := "cfg.json"

	buf, err = ioutil.ReadFile(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: Setup incorrect - unable to read %s: Error: %s", fn, err)
		return
	}

	err = json.Unmarshal(buf, &gCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: Setup incorrect - Unable to parse %s: data -->>%s<<-- Error: %s", fn, buf, err)
		return
	}

	// use that to setup connection to Geth
	client, err = rpc.Dial(gCfg.GethConnectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL: Error Connecting to Geth Server: %s server %s\n", err, gCfg.GethConnectionString)
		return
	}

	os.Exit(m.Run())
}
