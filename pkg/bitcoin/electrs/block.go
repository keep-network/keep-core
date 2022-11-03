package electrs

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

type blockHeader struct {
	Version           int32  `json:"version"`
	Timestamp         uint32 `json:"timestamp"`
	MerkleRoot        string `json:"merkle_root"`
	PreviousBlockHash string `json:"previousblockhash"`
	Nonce             uint32 `json:"nonce"`
	Bits              uint32 `json:"bits"`
}

func (b *blockHeader) convert() (bitcoin.BlockHeader, error) {
	result := bitcoin.BlockHeader{
		Version: b.Version,
		Time:    b.Timestamp,
		Bits:    b.Bits,
		Nonce:   b.Nonce,
	}

	previousBlockHeaderHash, err := bitcoin.NewHashFromString(b.PreviousBlockHash, bitcoin.ReversedByteOrder)
	if err != nil {
		return result, fmt.Errorf(
			"failed to decode previous block hash from [%s]: %w",
			b.PreviousBlockHash,
			err,
		)
	}
	result.PreviousBlockHeaderHash = previousBlockHeaderHash

	merkleRootHash, err := bitcoin.NewHashFromString(b.MerkleRoot, bitcoin.ReversedByteOrder)
	if err != nil {
		return result, fmt.Errorf(
			"failed to decode merkle root hash from [%s]: %w",
			b.MerkleRoot,
			err,
		)
	}
	result.MerkleRootHash = merkleRootHash

	return result, nil
}
