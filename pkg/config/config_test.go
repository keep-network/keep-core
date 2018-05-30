package config

import (
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	os.Setenv("KEEP_ETHEREUM_PASSWORD", "not-my-password")

	cfg, err := ReadConfig("./test/config.toml")
	if err != nil {
		dir, _ = os.Getwd()
		t.Errorf("Error: %s, in %s, fn=./test/config.toml\n", err, dir)
	}
	expectedURL := "ws://192.168.0.157:8546"
	if cfg.Ethereum.URL != expectedURL {
		t.Errorf("Error: Did not correctly read in ./test/config.toml, Expected [%s] Got [%s]\n", expectedURL, cfg.Ethereum.URL)
	}
	expectedURLRPC := "http://192.168.0.157:8545"
	if cfg.Ethereum.URLRPC != expectedURLRPC {
		t.Errorf("Error: Did not correctly read in ./test/config.toml, Expected [%s] Got [%s]\n", expectedURL, cfg.Ethereum.URL)
	}
	expectedAddress := "0x639deb0dd975af8e4cc91fe9053a37e4faf37649"
	if vv, ok := cfg.Ethereum.ContractAddresses["KeepRandomBeacon"]; !ok {
		t.Errorf("failed read of test/config.toml, expected key in map [KeepRandomBeacon].  Key missing.\n")
	} else if vv != expectedAddress {
		t.Errorf("in test/config.toml\ngot address: %s\nwant address: %s\n", vv, expectedAddress)
	}
}
