package tbtc

import (
	"encoding/hex"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"testing"
)

func TestDeposit_Script(t *testing.T) {
	hexToSlice := func(hexString string) []byte {
		bytes, err := hex.DecodeString(hexString)
		if err != nil {
			t.Fatalf("error while converting [%v]: [%v]", hexString, err)
		}
		return bytes
	}

	// Fill only the fields relevant for script computation.
	d := new(deposit)
	copy(d.depositor[:], hexToSlice("934b98637ca318a4d6e7ca6ffd1690b8e77df637"))
	copy(d.blindingFactor[:], hexToSlice("f9f0c90d00039523"))
	copy(d.walletPublicKeyHash[:], hexToSlice("8db50eb52063ea9d98b3eac91489a90f738986f6"))
	copy(d.refundPublicKeyHash[:], hexToSlice("28e081f285138ccbe389c1eb8985716230129f89"))
	copy(d.refundLocktime[:], hexToSlice("60bcea61"))

	script, err := d.script()
	if err != nil {
		t.Fatal(err)
	}

	expectedScript := hexToSlice(
		"14934b98637ca318a4d6e7ca6ffd1690b8e77df6377508f9f0c90d0003" +
			"95237576a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a" +
			"91428e081f285138ccbe389c1eb8985716230129f89880460bcea61b175ac68",
	)

	testutils.AssertBytesEqual(t, expectedScript, script)
}
