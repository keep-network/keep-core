package main

import "testing"

func TestReadConfig(t *testing.T) {
	cfg, err := readConfig("./test/config.toml")
	if err != nil {
		t.Errorf("Error:%s\n", err)
	}
	expectedURL := "ws://192.168.0.157:8546"
	if cfg.Ethereum.URL != expectedURL {
		t.Errorf("Error: Did not correctly read in ./test/config.toml, Expected [%s] Got [%s]\n", expectedURL, cfg.Ethereum.URL)
	}
	expectedAddress := "0x639deb0dd975af8e4cc91fe9053a37e4faf37649"
	if _, ok := cfg.Ethereum.ContractAddresses["KeepRandomBeacon"]; !ok {
		t.Errorf("Error: Did not correctly read in ./test/config.toml, expecte to have a key in the map for [KeepRandomBeacon].  Did not find one.\n")
	} else if vv, ok2 := cfg.Ethereum.ContractAddresses["KeepRandomBeacon"]; ok2 && vv != expectedAddress {
		t.Errorf("Error: Did not correctly read in ./test/config.toml, for KeepRandomBeacon expecte to have an address of [%s] got [%s]\n", expectedAddress, vv)
	}
}
