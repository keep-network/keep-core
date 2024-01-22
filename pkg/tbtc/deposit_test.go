package tbtc

import (
	"encoding/hex"
	"testing"

	"github.com/keep-network/keep-core/pkg/chain"

	"github.com/keep-network/keep-core/internal/testutils"
)

func TestDeposit_Script(t *testing.T) {
	hexToSlice := func(hexString string) []byte {
		bytes, err := hex.DecodeString(hexString)
		if err != nil {
			t.Fatalf("error while converting [%v]: [%v]", hexString, err)
		}
		return bytes
	}

	var tests = map[string]struct {
		depositor           string
		blindingFactor      string
		walletPublicKeyHash string
		refundPublicKeyHash string
		refundLocktime      string
		extraData           string
		expectedScript      string
	}{
		"no extra data": {
			depositor:           "934b98637ca318a4d6e7ca6ffd1690b8e77df637",
			blindingFactor:      "f9f0c90d00039523",
			walletPublicKeyHash: "8db50eb52063ea9d98b3eac91489a90f738986f6",
			refundPublicKeyHash: "28e081f285138ccbe389c1eb8985716230129f89",
			refundLocktime:      "60bcea61",
			extraData:           "",
			expectedScript: "14934b98637ca318a4d6e7ca6ffd1690b8e77df637750" +
				"8f9f0c90d000395237576a9148db50eb52063ea9d98b3eac91489a90f" +
				"738986f68763ac6776a91428e081f285138ccbe389c1eb89857162301" +
				"29f89880460bcea61b175ac68",
		},
		"with extra data": {
			depositor:           "934b98637ca318a4d6e7ca6ffd1690b8e77df637",
			blindingFactor:      "f9f0c90d00039523",
			walletPublicKeyHash: "8db50eb52063ea9d98b3eac91489a90f738986f6",
			refundPublicKeyHash: "28e081f285138ccbe389c1eb8985716230129f89",
			refundLocktime:      "60bcea61",
			extraData: "a9b38ea6435c8941d6eda6a46b68e3e2117196995bd154ab55" +
				"196396b03d9bda",
			expectedScript: "14934b98637ca318a4d6e7ca6ffd1690b8e77df637752" +
				"0a9b38ea6435c8941d6eda6a46b68e3e2117196995bd154ab55196396" +
				"b03d9bda7508f9f0c90d000395237576a9148db50eb52063ea9d98b3e" +
				"ac91489a90f738986f68763ac6776a91428e081f285138ccbe389c1eb" +
				"8985716230129f89880460bcea61b175ac68",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Fill only the fields relevant for script computation.
			d := new(Deposit)
			d.Depositor = chain.Address(test.depositor)
			copy(d.BlindingFactor[:], hexToSlice(test.blindingFactor))
			copy(d.WalletPublicKeyHash[:], hexToSlice(test.walletPublicKeyHash))
			copy(d.RefundPublicKeyHash[:], hexToSlice(test.refundPublicKeyHash))
			copy(d.RefundLocktime[:], hexToSlice(test.refundLocktime))

			if len(test.extraData) > 0 {
				var extraData [32]byte
				copy(extraData[:], hexToSlice(test.extraData))
				d.ExtraData = &extraData
			}

			script, err := d.Script()
			if err != nil {
				t.Fatal(err)
			}

			expectedScript := hexToSlice(test.expectedScript)

			testutils.AssertBytesEqual(t, expectedScript, script)
		})
	}
}
