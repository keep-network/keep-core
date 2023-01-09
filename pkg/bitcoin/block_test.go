package bitcoin

import (
	"encoding/hex"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestBlockHeaderSerialize(t *testing.T) {
	// Test data comes from a Bitcoin testnet block:
	// https://live.blockcypher.com/btc-testnet/block/000000000000002af10911b8db32ed34dc6ea6515f84af5f7b82973c9a839e6d/

	previousBlockHeaderHash, err := NewHashFromString(
		"000000000066450030efdf72f233ed2495547a32295deea1e2f3a16b1e50a3a5",
		ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	merkleRootHash, err := NewHashFromString(
		"1251774996b446f85462d5433f7a3e384ac1569072e617ab31e86da31c247de2",
		ReversedByteOrder,
	)
	if err != nil {
		t.Fatal(err)
	}

	blockHeader := BlockHeader{
		Version:                 536870916,
		PreviousBlockHeaderHash: previousBlockHeaderHash,
		MerkleRootHash:          merkleRootHash,
		Time:                    1641914003,
		Bits:                    436256810,
		Nonce:                   778087099,
	}

	actualSerializedHeader := blockHeader.Serialize()

	expectedSerializedHeader, err := hex.DecodeString(
		"04000020a5a3501e6ba1f3e2a1ee5d29327a549524ed33f272dfef30004566000000" +
			"0000e27d241ca36de831ab17e6729056c14a383e7a3f43d56254f846b4964977" +
			"5112939edd612ac0001abbaa602e",
	)
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertBytesEqual(
		t,
		expectedSerializedHeader,
		actualSerializedHeader[:],
	)
}
