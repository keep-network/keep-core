package bitcoin

import "github.com/keep-network/keep-core/pkg/bitcoin"

// Blocks holds details of block header data used as test vectors.
var Blocks = map[bitcoin.Network]struct {
	BlockHeight    uint
	BlockHeader    *bitcoin.BlockHeader
	CoinbaseTxHash bitcoin.Hash
}{
	bitcoin.Testnet: {
		BlockHeight: 2135502,
		BlockHeader: &bitcoin.BlockHeader{
			Version:                 536870916,
			PreviousBlockHeaderHash: hashFromString("000000000066450030efdf72f233ed2495547a32295deea1e2f3a16b1e50a3a5"),
			MerkleRootHash:          hashFromString("1251774996b446f85462d5433f7a3e384ac1569072e617ab31e86da31c247de2"),
			Time:                    1641914003,
			Bits:                    436256810,
			Nonce:                   778087099,
		},
		CoinbaseTxHash: hashFromString("1f523d1ce7553ec609bae104812dede95aa38eb13d2c2c6b64ffe868bbc1a54c"),
	},
	bitcoin.Mainnet: {
		BlockHeight: 792379,
		BlockHeader: &bitcoin.BlockHeader{
			Version:                 547356672,
			PreviousBlockHeaderHash: hashFromString("000000000000000000035d7345203315d8c73c37da256eebbe58a683331b60f5"),
			MerkleRootHash:          hashFromString("a406a1495c4c4bed2d458f06f393cb841a545cbf7cc04b1ee3d9ae9e6eccb402"),
			Time:                    1685622977,
			Bits:                    386236009,
			Nonce:                   2812580720,
		},
		CoinbaseTxHash: hashFromString("5a21bae1a78a06519166b3a26d6b20e6985a91fc66a48384982451211b520234"),
	},
}