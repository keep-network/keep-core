package ethereum

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	erpc "github.com/ethereum/go-ethereum/rpc"
)

func TestMain(m *testing.M) {
	// fmt.Printf("%sAT:%s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
	var envFn = os.Getenv("KEEP_TEST_CFG")
	if envFn == "" {
		envFn = "int_cfg.json"
	}
	err := ReadTestConfig(envFn)
	if err != nil {
		fmt.Printf("FAIL - test did not read configuration file correctly, %s -- no test were run\n", err)
		os.Exit(1)
	}
	ec, err := ConnectTestToGeth(TestConfig)
	if err != nil {
		fmt.Printf("FAIL - test did not connect to Geth, %s -- no test were run\n", err)
		os.Exit(1)
	}
	EthConn = ec
	retCode := m.Run()
	os.Exit(retCode)
}

// ConnectTest makes the network connection to the Ethereum network.  Note: for
// other things to work correctly the configuration will need to reference a
// websocket, "ws://", or local IPC connection.
func ConnectTestToGeth(cfg TestConfigType) (*ethereumChain, error) {
	client, err := ethclient.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf(
			"error Connecting to Geth Server: %s [%v]",
			cfg.URL,
			err,
		)
	}

	clientws, err := erpc.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf(
			"error Connecting to Geth Server: %s [%v]",
			cfg.URL,
			err,
		)
	}

	clientrpc, err := erpc.Dial(cfg.URLRPC)
	if err != nil {
		return nil, fmt.Errorf(
			"error Connecting to Geth Server: %s [%v]",
			cfg.URL,
			err,
		)
	}

	eCfg := Config{
		URL:               cfg.URL,
		URLRPC:            cfg.URLRPC,
		ContractAddresses: make(map[string]string),
		Account: Account{
			Address:         cfg.Address,
			KeyFile:         cfg.KeyFile,
			KeyFilePassword: cfg.KeyFilePassword,
		},
	}

	// Setup to use contracts through proxies - addresses are set to proxy addresses.
	eCfg.ContractAddresses["KeepGroup"] = cfg.ContractAddress["KeepGroup"]
	eCfg.ContractAddresses["KeepGroupImplV1"] = cfg.ContractAddress["KeepGroup"]
	eCfg.ContractAddresses["KeepRandomBeacon"] = cfg.ContractAddress["KeepRandomBeacon"]
	eCfg.ContractAddresses["KeepRandomBeaconImplV1"] = cfg.ContractAddress["KeepRandomBeacon"]

	pv := &ethereumChain{
		config:    eCfg,
		client:    client,
		clientRPC: clientrpc,
		clientWS:  clientws,
	}

	keepRandomBeaconContract, err := newKeepRandomBeacon(pv)
	if err != nil {
		return nil, fmt.Errorf(
			"error attaching to KeepRandomBeacon contract: [%v]",
			err,
		)
	}
	pv.keepRandomBeaconContract = keepRandomBeaconContract

	keepGroupContract, err := newKeepGroup(pv)
	if err != nil {
		return nil, fmt.Errorf("error attaching to KeepGroup contract: [%v]", err)
	}
	pv.keepGroupContract = keepGroupContract

	return pv, nil
}

// SetAddressToCallImpl changes the config to allow calls directly to the implementation contracts.
func SetAddressToCallImpl(cfg TestConfigType, pv *ethereumChain) error {
	pv.config.ContractAddresses["KeepGroup"] = cfg.ContractAddress["KeepGroupImplV1"]
	pv.config.ContractAddresses["KeepGroupImplV1"] = cfg.ContractAddress["KeepGroupImplV1"]
	pv.config.ContractAddresses["KeepRandomBeacon"] = cfg.ContractAddress["KeepRandomBeaconImplV1"]
	pv.config.ContractAddresses["KeepRandomBeaconImplV1"] = cfg.ContractAddress["KeepRandomBeaconImplV1"]
	keepRandomBeaconContract, err := newKeepRandomBeacon(pv)
	if err != nil {
		return fmt.Errorf(
			"error attaching to KeepRandomBeacon contract: [%v]",
			err,
		)
	}
	pv.keepRandomBeaconContract = keepRandomBeaconContract

	keepGroupContract, err := newKeepGroup(pv)
	if err != nil {
		return fmt.Errorf("error attaching to KeepGroup contract: [%v]", err)
	}
	pv.keepGroupContract = keepGroupContract
	return nil
}

// SetAddressToCallProxy changes the config to allow calls to proxy contracts.
func SetAddressToCallProxy(cfg TestConfigType, pv *ethereumChain) error {
	pv.config.ContractAddresses["KeepGroup"] = cfg.ContractAddress["KeepGroup"]
	pv.config.ContractAddresses["KeepGroupImplV1"] = cfg.ContractAddress["KeepGroup"]
	pv.config.ContractAddresses["KeepRandomBeacon"] = cfg.ContractAddress["KeepRandomBeacon"]
	pv.config.ContractAddresses["KeepRandomBeaconImplV1"] = cfg.ContractAddress["KeepRandomBeacon"]
	keepRandomBeaconContract, err := newKeepRandomBeacon(pv)
	if err != nil {
		return fmt.Errorf(
			"error attaching to KeepRandomBeacon contract: [%v]",
			err,
		)
	}
	pv.keepRandomBeaconContract = keepRandomBeaconContract

	keepGroupContract, err := newKeepGroup(pv)
	if err != nil {
		return fmt.Errorf("error attaching to KeepGroup contract: [%v]", err)
	}
	pv.keepGroupContract = keepGroupContract
	return nil
}

func ReadTestConfig(fn string) error {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &TestConfig)
	if err != nil {
		return err
	}
	return nil
}
