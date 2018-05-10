package ethereum

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

type CfgType struct {
	GethConnectionString string
}

var (
	gCfg   CfgType
	client *rpc.Client
)

func TestEthereumBlockTest(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		wait         int
		want         time.Duration
		errorMessage string
	}{
		"does wait for a block": {
			wait:         1,
			want:         time.Duration(100000000),
			errorMessage: "Failed to wait for a single block",
		},
		"waited for a longer time": {
			wait:         2,
			want:         time.Duration(200000000),
			errorMessage: "Failed to wait for 2 blocks",
		},
		"doesn't wait if 0 blocks": {
			wait:         0,
			want:         time.Duration(1000),
			errorMessage: "Failed for a 0 block wait",
		},
		"invalid value": {
			wait:         -1,
			want:         time.Duration(0),
			errorMessage: "Waiting for a time when it should have errored",
		},
	}

	var e chan interface{}

	tim := 240 // Force test to fail if not completed in 4 minutes
	tick := time.NewTimer(time.Duration(tim) * time.Second)

	go func() {
		select {
		case e <- tick:
			t.Fatal("Test ran too long - it failed")
		}
	}()

	waitForBlock := BlockCounter(client)

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			start := time.Now().UTC()
			waitForBlock.WaitForBlocks(test.wait)
			end := time.Now().UTC()

			elapsed := end.Sub(start)
			if elapsed < test.want {
				t.Error(test.errorMessage)
			}
		})
	}
}

func TestMain(m *testing.M) {

	fn := "cfg.json" // TODO - pick this up from the environment.

	buf, err := ioutil.ReadFile(fn)
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
