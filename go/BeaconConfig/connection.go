package BeaconConfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// BeaconConfig contains configuration for the threshold relay beacon, typically
// from the underlying blockchain.
//
// KeyFilePassword        	If '$$ENV:<<name>>' will look for <<name>> in environment
//							and pick password form environment intead of using the value
//							in the file.
//
// BlockTimeout 			This is the number of seconds to wait for a block to arrive
//							before giving up and just clearing the "wait".  This is intended
//							for testing with "truffle" where no actual mining of new blocks
//							occures.  A 0 value disables this.
//
type BeaconConfig struct {
	GroupSize                  int    `json:"group_size"`
	Threshold                  int    `json:"threshold"`
	BlockTimeout               int    `json:"block_timeout"`
	KeyFile                    string `json:"key_file"`
	KeyFilePassword            string `json:"key_file_password"` // See above.
	BeaconRelayContractAddress string `json:"beacon_relay_contract_address"`
	FromAddress                string `json:"from_address"` // Address in KeyFile
	GethServer                 string `json:"geth_server"`
	// more stuff for p2p config to be put at this location
}

// GetBeaconConfig Get the latest threshold relay beacon configuration.
// TODO Make this actually look up/update from chain information.
func GetBeaconConfig(fn string) BeaconConfig {
	tBc := BeaconConfig{
		GroupSize:    10, // Establish some defaults before reading in file
		Threshold:    4,  //
		BlockTimeout: 0,  // Assume that we are connecting at "Geth" to truffle test rpc - no timeout on blocks
	}
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatalf("could not open %s error: %s", fn, err)
	}
	err = json.Unmarshal(data, &tBc)
	if err != nil {
		log.Fatalf("could not JSON parrse %s error: %s", fn, err)
	}

	if strings.HasPrefix(tBc.KeyFilePassword, "$$ENV:") {
		if len(tBc.KeyFilePassword) <= 6 {
			log.Fatalf("Malformed environemt variable reqeust for password, need name of environment variable after $$ENV:")
		}
		env := tBc.KeyFilePassword[6:]
		val := os.Getenv(env)
		if val == "" {
			log.Fatalf("Missing environment variable %s - for password, Look at key_file_password in %s - set and export password", env, fn)
		}
		tBc.KeyFilePassword = val
	}

	return tBc
}

// NewEthConnection creates a websocket (ws://) or IPC connection to the geth client.
// This can also create a HTTP connection - but most of what we are working on will
// require the websocket/ipc type connections.
func NewRpcConnection(GethServer string, timeout int) (client *rpc.Client, err error) {
	// HTTP/ws:/ipc - setup connection to Geth
	client, err = rpc.Dial(GethServer)
	if err != nil {
		log.Fatalf("could not create %s client: %v", GethServer, err)
	}
	return
}

func OpenRpcConnection(cfgData BeaconConfig) (client *rpc.Client, err error) {
	client, err = NewRpcConnection(cfgData.GethServer, cfgData.BlockTimeout)
	return
}

func NewEthConnection(GethServer string, timeout int) (conn *ethclient.Client, err error) {
	// HTTP/ws:/ipc - setup connection to Geth

	//	conn, err := ethclient.Dial(*GethServer)
	//	if err != nil {
	//		log.Fatalf("Failed to connect to the Ethereum client: %v at address: %s", err, *GethServer)
	//	}
	conn, err = ethclient.Dial(GethServer)
	if err != nil {
		log.Fatalf("could not create %s client: %v", GethServer, err)
	}
	return
}

func OpenEthConnection(cfgData BeaconConfig) (conn *ethclient.Client, err error) {
	conn, err = NewEthConnection(cfgData.GethServer, cfgData.BlockTimeout)
	return
}
