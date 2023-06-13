package bitcoin

import "github.com/keep-network/keep-core/pkg/bitcoin"

// TxMerkleProofs holds details of transaction merkle proof data used as test vectors.
var TxMerkleProofs = map[bitcoin.Network]struct {
	TxHash      bitcoin.Hash
	BlockHeight uint
	MerkleProof *bitcoin.TransactionMerkleProof
}{
	bitcoin.Testnet: {
		TxHash: hashFromString(
			"72e7fd57c2adb1ed2305c4247486ff79aec363296f02ec65be141904f80d214e",
		),
		BlockHeight: 1569342,
		MerkleProof: &bitcoin.TransactionMerkleProof{
			BlockHeight: 1569342,
			MerkleNodes: []string{
				"8b5bbb5bdf6727bf70fad4f46fe4eaab04c98119ffbd2d95c29adf32d26f8452",
				"53637bacb07965e4a8220836861d1b16c6da29f10ea9ab53fc4eca73074f98b9",
				"0267e738108d094ceb05217e2942e9c2a4c6389ac47f476f572c9a319ce4dfbc",
				"34e00deec50c48d99678ca2b52b82d6d5432326159c69e7233d0dde0924874b4",
				"7a53435e6c86a3620cdbae510901f17958f0540314214379197874ed8ed7a913",
				"6315dbb7ce350ceaa16cd4c35c5a147005e8b38ca1e9531bd7320629e8d17f5b",
				"40380cdadc0206646208871e952af9dcfdff2f104305ce463aed5eeaf7725d2f",
				"5d74bae6a71fd1cff2416865460583319a40343650bd4bb89de0a6ae82097037",
				"296ddccfc659e0009aad117c8ed15fb6ff81c2bade73fbc89666a22708d233f9",
			},
			Position: 176,
		},
	},
	bitcoin.Mainnet: {
		TxHash: hashFromString(
			"c11feea0f8b9e256e0a70ad58f88c7e3800d8ffbfb85edf7169629584a97a7ab",
		),
		BlockHeight: 793508,
		MerkleProof: &bitcoin.TransactionMerkleProof{
			BlockHeight: 793508,
			MerkleNodes: []string{
				"557900a6e94e751da2de3cc4c322a12d682b0285996187b4276c8c1936940efe",
				"910ae74a83dda3faa97dff964eac4d4102119161f8c850c578afff33a5e43c33",
				"0d3648b1fba7f375b54a268644e0f6f7f739a490156752916312dc945b3c5a5d",
				"8d19e661f8b3bfdb1d4228a5f8aa855808ef5a0a5ffa65e32551c8f3599ff01b",
				"ea1bbd7e974a8455f38bb4c047b2324fa62df49dc486018f3604c02a71473261",
				"b315a3bb6dfc35e53eaf8756214a521ece6b7c5323d7ed616e6e828de6d024e2",
				"d7272801d3f9d8dd459451cdf066cf4fab03cc2198391213e22f49b5a5563c2e",
				"baf373eb68aa0fece7e680d09183d7b441073a5b0745bbad92345fdd0447badb",
				"016c2eb0e2be1cdbed723690e6de31030a03bece34edf6944ba5d55ab64fa305",
				"dc74dc3b4b84293cce52b9d739ca29501e612ceb566ea391a5fde1341d46c950",
				"9c941fd7c44fc29403865d57b78941664a236f8506b7a8676d1a833eeea16ff3",
				"6ad9f32a66f3d376ea4bcc48d807f565e8b5228cf5b45eefb42879c270263d2d",
			},
			Position: 227,
		},
	},
}
