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
			"4db013214a88ac4bff07a7ef160c6d7cf82dea26f747c035c7aa6f3d7a7f4d4c",
		),
		BlockHeight: 792384,
		MerkleProof: &bitcoin.TransactionMerkleProof{
			BlockHeight: 792384,
			MerkleNodes: []string{
				"dde09ada70242eb53da60ad14966d56a2fb13305e3d6c91d450b345579810c16",
				"8bdf3bfac9c4ffc5949b18b2d381a954432bae1ebf419a08dd237b0e70f3c671",
				"6fdd8d8b8fd2c248600a4d0da2b3d06449962461a1590443dee6859470c5bb8a",
				"59a57e6bbd22a797960c0f690aa66d14ce3491a998c79b48ddd40fa5d183b27e",
				"befd298b231c3f957e26b8dcb3b2a8bc011c387d7fb0addbdde5c2fca408f046",
				"5738cefbca2c0d0f8042510055f660dca093407413c67033ca10ee4f9768ff91",
				"af1937b829566b0e684c22c60c8bc06ba762d914263259a3e37ae66955e9a3a4",
				"7a34a70b80fe1db7aa9d94efc3f47fb2a4fbf33185c17da06015e0ccf7e733be",
				"a2ca5d291a2c16e14abf9a6d828fb324c744a131434bcb5b15f1d6d7ded6ca0f",
				"b3f830dfbfbcda88e1740357cbee20d89e6c258acb740fd3ed6080cd54d4b54d",
				"86aab3b304c7612889565a85bc540f1635d6dd810d3d9c985e12c0945cba9f8c",
				"f7b94a757adc025ecb47bddbbc655dee3eb0916b4b4d2c5919cc9690f51835db",
			},
			Position: 6,
		},
	},
}
