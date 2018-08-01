package tecdsa

import (
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

func TestKeysGeneration(t *testing.T) {
	curve := btcec.S256()
	networkParams := &chaincfg.TestNet3Params

	// Validate Private Key Generation
	privateKey, err := btcec.NewPrivateKey(curve)
	if err != nil {
		t.Fatal(err)
	}

	recoveredPrivateKey, _ := btcec.PrivKeyFromBytes(curve, privateKey.Serialize())
	if !reflect.DeepEqual(privateKey, recoveredPrivateKey) {
		t.Fatalf("Recovered Private Key doesn't match expected")
	}

	// Validate WIF Generation
	wif, err := btcutil.NewWIF(privateKey, networkParams, false)
	if err != nil {
		t.Fatal(err)
	}
	recoveredWif, err := btcutil.DecodeWIF(wif.String())
	if err != nil {
		t.Fatal(err)
	}
	if wif.String() != recoveredWif.String() {
		t.Fatalf("Recovered WIF doesn't match expected")
	}

	// Validate Address Generation
	address, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), networkParams)
	if err != nil {
		t.Fatal(err)
	}
	recoveredAddress, err := btcutil.DecodeAddress(address.EncodeAddress(), networkParams)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(address.EncodeAddress(), recoveredAddress.EncodeAddress()) {
		t.Fatalf("Recovered Address doesn't match expected")
	}
}
