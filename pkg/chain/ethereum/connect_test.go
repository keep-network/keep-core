package ethereum

import "testing"

func TestResolveContractByValidName(t *testing.T) {
	config := Config{
		ContractAddresses: make(map[string]string),
	}

	conquerorAddress := "0x1895e4A71d0956553cf80f2ccac69642a2f1bFF4"
	barbarianAddress := "0x0b185C37E1C9D01437c800a8B60fA0845742c271"

	config.ContractAddresses["ConanTheConqueror"] = conquerorAddress
	config.ContractAddresses["ConanTheBarbarian"] = barbarianAddress

	address, err := addressForContract(config, "ConanTheConqueror")
	if err != nil {
		t.Fatal(err)
	}
	if address.Hex() != conquerorAddress {
		t.Fatalf(
			"Unexpected contract address\nExpected: [%v]\nActual:   [%v]\n",
			conquerorAddress,
			address.Hex(),
		)
	}

	address, err = addressForContract(config, "ConanTheBarbarian")
	if err != nil {
		t.Fatal(err)
	}
	if address.Hex() != barbarianAddress {
		t.Fatalf(
			"Unexpected contract address\nExpected: [%v]\nActual:   [%v]\n",
			barbarianAddress,
			address.Hex(),
		)
	}
}
