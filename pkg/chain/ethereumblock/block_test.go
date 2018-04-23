package ethereumblock_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"

	"github.com/keep-network/keep-core/pkg/chain/ethereumblock" // /Users/corwin/go/src/github.com/keep-network/keep-core/pkg/chain/ethereumblock
)

type CfgType struct {
	Server string
}

var gCfg CfgType
var initTest bool
var client *rpc.Client

func Test_BlockTest01(t *testing.T) {

	// ----------------------------------------------------------------------------------------------------
	// Test code that needs to do a timeout!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// ----------------------------------------------------------------------------------------------------

	// get connection info - cfg.json
	err := InitTest(t)
	if err != nil {
		t.Fatalf("Initialization Failed, error: %s", err)
		return
	}

	godebug.Printf(db02, "AT: %s\n", godebug.LF())

	// func BlockCounter(client *rpc.Client) chain.BlockCounter {
	countWait := ethereumblock.BlockCounter(client)
	if db01 {
		fmt.Printf("Before Wait\n")
	}
	start := time.Now()
	countWait.WaitForBlocks(1)
	tm := time.Now()
	elapsed := tm.Sub(start)
	if elapsed < 1000000000 {
		t.Errorf("Did not wait\n")
	}
	if db01 {
		fmt.Printf("After Wait, %d\n", elapsed)
	}

	godebug.Printf(db02, "AT: %s\n", godebug.LF())
	start = time.Now()
	countWait.WaitForBlocks(2)
	tm = time.Now()
	elapsed = tm.Sub(start)
	if elapsed < 2000000000 {
		t.Errorf("Did not wait\n")
	}

	godebug.Printf(db02, "AT: %s\n", godebug.LF())
	if db01 {
		fmt.Printf("Before test #3 , %d\n", elapsed)
	}
	start = time.Now()
	countWait.WaitForBlocks(0)
	tm = time.Now()
	elapsed = tm.Sub(start)
	if elapsed < 1000 {
		t.Errorf("Did not wait\n")
	}
	godebug.Printf(db02, "AT: %s\n", godebug.LF())
}

func InitTest(t *testing.T) (err error) {
	if !initTest {
		initTest = true

		var buf []byte

		godebug.Printf(db02, "AT: %s\n", godebug.LF())

		// get connection info - cfg.json
		buf, err = ioutil.ReadFile("cfg.json")
		if err != nil {
			t.Fatalf("Setup incorrect - unable to read cfg.json: error: %s", err)
			return
		}

		err = json.Unmarshal(buf, &gCfg)
		if err != nil {
			t.Fatalf("Setup incorrect - Unable to parse cfg.json, error: %s", err)
			return
		}

		godebug.Printf(db02, "AT: %s\n", godebug.LF())

		// use that to setup connection to Geth
		client, err = rpc.Dial(gCfg.Server)
		if err != nil {
			fmt.Printf("Error Connecting to Server: %s server %s\n", err, gCfg.Server)
			t.Fatalf("Setup incorrect - failed to connect to server")
			return
		}

		godebug.Printf(db02, "%sSuccessfully connected to Geth: AT: %s%s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
	}
	return
}

const db01 = false
const db02 = true
