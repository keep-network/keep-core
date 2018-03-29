package BeaconConfig

import (
	"os"
	"testing"
)

func Test_ReadInConfig(t *testing.T) {

	os.Setenv("wallet_password", "password")

	bc := GetBeaconConfig("./test/config1.json")

	if bc.KeyFilePassword != "password" {
		t.Errorf("Failed to pull password form environment")
	}
	if bc.GroupSize != 10 {
		t.Errorf("Failed to get a default GroupSize")
	}
	if bc.BeaconRelayContractAddress != "0xe705ab560794cf4912960e5069d23ad6420acde7" {
		t.Errorf("Failed to read in value for contract address")
	}
}
