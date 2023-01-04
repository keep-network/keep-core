package bitcoin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestTransaction_SerializeRoundtrip(t *testing.T) {
	transaction := transactionFixture(t)

	witnessBytes := transaction.Serialize(Witness)

	expectedWitnessBytes := hexToSlice(
		t,
		"010000000001036896f9abcac13ce6bd2b80d125bedf997ff6330e999f2f605e"+
			"a15ea542f2eaf80000000000ffffffffed0ae94da996c6f3b89dfe967675d48"+
			"08251db93e81022ae9e038d06f92efed400000000c948304502210092327ddf"+
			"f69a2b8c7ae787c5d590a2f14586089e6339e942d56e82aa42052cd902204c0"+
			"d1700ba1ac617da27fee032a57937c9607f0187199ed3c46954df845643d701"+
			"2103989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8"+
			"581d94c5c14934b98637ca318a4d6e7ca6ffd1690b8e77df6377508f9f0c90d"+
			"000395237576a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6"+
			"776a914e257eccafbc07c381642ce6e7e55120fb077fbed8804e0250162b175"+
			"ac68ffffffffe37f552fc23fa0032bfd00c8eef5f5c22bf85fe4c6e73585771"+
			"9ff8a4ff66eb80000000000ffffffff0180ed0000000000001600148db50eb5"+
			"2063ea9d98b3eac91489a90f738986f602483045022100baf754252d0d6a49a"+
			"ceba7eb0ec40b4cc568e8c659e168b96598a11cf56dc078022051117466ee99"+
			"8a3fc72221006817e8cfe9c2e71ad622ff811a0bf100d888d49c012103989d2"+
			"53b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d90003"+
			"473044022014a535eb334656665ac69a678dbf7c019c4f13262e9ea4d195c61"+
			"a00cd5f698d022023c0062913c4614bdff07f94475ceb4c585df53f71611776"+
			"c3521ed8f8785913012103989d253b17a6a0f41838b84ff0d20e8898f9d7b1a"+
			"98f2564da4cc29dcf8581d95c14934b98637ca318a4d6e7ca6ffd1690b8e77d"+
			"f6377508f9f0c90d000395237576a9148db50eb52063ea9d98b3eac91489a90"+
			"f738986f68763ac6776a914e257eccafbc07c381642ce6e7e55120fb077fbed"+
			"8804e0250162b175ac6800000000",
	)
	testutils.AssertBytesEqual(t, expectedWitnessBytes, witnessBytes)

	standardBytes := transaction.Serialize(Standard)
	expectedStandardBytes := hexToSlice(
		t,
		"01000000036896f9abcac13ce6bd2b80d125bedf997ff6330e999f2f60"+
			"5ea15ea542f2eaf80000000000ffffffffed0ae94da996c6f3b89dfe967675d"+
			"4808251db93e81022ae9e038d06f92efed400000000c948304502210092327d"+
			"dff69a2b8c7ae787c5d590a2f14586089e6339e942d56e82aa42052cd902204"+
			"c0d1700ba1ac617da27fee032a57937c9607f0187199ed3c46954df845643d7"+
			"012103989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dc"+
			"f8581d94c5c14934b98637ca318a4d6e7ca6ffd1690b8e77df6377508f9f0c9"+
			"0d000395237576a9148db50eb52063ea9d98b3eac91489a90f738986f68763a"+
			"c6776a914e257eccafbc07c381642ce6e7e55120fb077fbed8804e0250162b1"+
			"75ac68ffffffffe37f552fc23fa0032bfd00c8eef5f5c22bf85fe4c6e735857"+
			"719ff8a4ff66eb80000000000ffffffff0180ed0000000000001600148db50e"+
			"b52063ea9d98b3eac91489a90f738986f600000000",
	)
	testutils.AssertBytesEqual(t, expectedStandardBytes, standardBytes)

	deserialized := new(Transaction)
	err := deserialized.Deserialize(witnessBytes)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(transaction, deserialized) {
		t.Errorf("unexpected deserialized transaction")
	}
}

func TestTransaction_Hash(t *testing.T) {
	hash := transactionFixture(t).Hash()

	expectedHash := "435d4aff6d4bc34134877bd3213c17970142fdd04d4113d534120033b9eecb2e"

	testutils.AssertStringsEqual(
		t,
		"hash",
		expectedHash,
		hash.Hex(ReversedByteOrder),
	)
}

func TestTransaction_WitnessHash(t *testing.T) {
	hash := transactionFixture(t).WitnessHash()

	expectedHash := "6131ce6056c8c76eb92f17c64516cf71143bb289b55f9ab0cacb5b1c8f2bd94a"

	testutils.AssertStringsEqual(
		t,
		"hash",
		expectedHash,
		hash.Hex(ReversedByteOrder),
	)
}

// transactionFixture returns a real testnet transaction:
// https://live.blockcypher.com/btc-testnet/tx/435d4aff6d4bc34134877bd3213c17970142fdd04d4113d534120033b9eecb2e
func transactionFixture(t *testing.T) *Transaction {
	tx := new(Transaction)

	tx.Version = 1

	tx.Inputs = append(tx.Inputs, &TransactionInput{
		Outpoint: &TransactionOutpoint{
			TransactionHash: hexToHash(
				t,
				"f8eaf242a55ea15e602f9f990e33f67f99dfbe25d1802bbde6"+
					"3cc1caabf99668",
			),
			OutputIndex: 0,
		},
		SignatureScript: []byte{},
		Witness: [][]byte{
			hexToSlice(
				t,
				"3045022100baf754252d0d6a49aceba7eb0ec40b4cc568e8c6"+
					"59e168b96598a11cf56dc078022051117466ee998a3fc7222100681"+
					"7e8cfe9c2e71ad622ff811a0bf100d888d49c01",
			),
			hexToSlice(
				t,
				"03989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564"+
					"da4cc29dcf8581d9",
			),
		},
		Sequence: 0xffffffff,
	})
	tx.Inputs = append(tx.Inputs, &TransactionInput{
		Outpoint: &TransactionOutpoint{
			TransactionHash: hexToHash(
				t,
				"d4fe2ef9068d039eae2210e893db518280d4757696fe9db8f3"+
					"c696a94de90aed",
			),
			OutputIndex: 0,
		},
		SignatureScript: hexToSlice(
			t,
			"48304502210092327ddff69a2b8c7ae787c5d590a2f14586089e63"+
				"39e942d56e82aa42052cd902204c0d1700ba1ac617da27fee032a57937c"+
				"9607f0187199ed3c46954df845643d7012103989d253b17a6a0f41838b8"+
				"4ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d94c5c14934b98637ca"+
				"318a4d6e7ca6ffd1690b8e77df6377508f9f0c90d000395237576a9148d"+
				"b50eb52063ea9d98b3eac91489a90f738986f68763ac6776a914e257ecc"+
				"afbc07c381642ce6e7e55120fb077fbed8804e0250162b175ac68",
		),
		Witness:  [][]byte{},
		Sequence: 0xffffffff,
	})
	tx.Inputs = append(tx.Inputs, &TransactionInput{
		Outpoint: &TransactionOutpoint{
			TransactionHash: hexToHash(
				t,
				"b86ef64f8aff19778535e7c6e45ff82bc2f5f5eec800fd2b03"+
					"a03fc22f557fe3",
			),
			OutputIndex: 0,
		},
		SignatureScript: []byte{},
		Witness: [][]byte{
			hexToSlice(
				t,
				"3044022014a535eb334656665ac69a678dbf7c019c4f13262e"+
					"9ea4d195c61a00cd5f698d022023c0062913c4614bdff07f94475ce"+
					"b4c585df53f71611776c3521ed8f878591301",
			),
			hexToSlice(
				t,
				"03989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564"+
					"da4cc29dcf8581d9",
			),
			hexToSlice(
				t,
				"14934b98637ca318a4d6e7ca6ffd1690b8e77df6377508f9f0"+
					"c90d000395237576a9148db50eb52063ea9d98b3eac91489a90f738"+
					"986f68763ac6776a914e257eccafbc07c381642ce6e7e55120fb077"+
					"fbed8804e0250162b175ac68",
			),
		},
		Sequence: 0xffffffff,
	})

	tx.Outputs = append(tx.Outputs, &TransactionOutput{
		Value: 60800,
		PublicKeyScript: hexToSlice(
			t,
			"00148db50eb52063ea9d98b3eac91489a90f738986f6",
		),
	})

	tx.Locktime = 0

	return tx
}

func transactionFrom(t *testing.T, hex string) *Transaction {
	tx := new(Transaction)
	err := tx.Deserialize(hexToSlice(t, hex))
	if err != nil {
		t.Fatalf("error while converting [%v]: [%v]", hex, err)
	}
	return tx
}

func hexToSlice(t *testing.T, hexString string) []byte {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		t.Fatalf("error while converting [%v]: [%v]", hexString, err)
	}
	return bytes
}

func hexToHash(t *testing.T, hex string) Hash {
	hash, err := NewHashFromString(hex, ReversedByteOrder)
	if err != nil {
		t.Fatalf("error while converting [%v]: [%v]", hex, err)
	}
	return hash
}

func hexToPublicKet(t *testing.T, hex string) *ecdsa.PublicKey {
	x, y := elliptic.Unmarshal(btcec.S256(), hexToSlice(t, hex))
	return &ecdsa.PublicKey{
		Curve: btcec.S256(),
		X:     x,
		Y:     y,
	}
}
