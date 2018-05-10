package config

import "testing"

func TestReadConfig(t *testing.T) {
	cfg, err := ReadConfig("./test/config.toml")
	if err != nil {
		t.Errorf("Error:%s\n", err)
	}
	if cfg.Ethereum.URL != "ws://192.168.0.157:8546" {
		t.Errorf("Error: Did not correctly read in ./test/config.toml\n")
	}
	if vv, ok := cfg.Ethereum.ContractAddresses["KeepRandomBeacon"]; !ok || vv != "0x639deb0dd975af8e4cc91fe9053a37e4faf37649" {
		t.Errorf("Error: Did not correctly read in ./test/config.toml, invalid map\n")
	}
}
