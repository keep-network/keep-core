package bitcoin

import (
	"encoding/hex"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

// Transactions holds details of the transactions used as test vectors.
var Transactions = map[string]struct {
	TxHash      bitcoin.Hash
	BlockHeight uint
	BitcoinTx   bitcoin.Transaction
}{
	// Transactions data taken from TBTCv2 Deposit Sweeps
	// See: https://github.com/keep-network/tbtc-v2/blob/8b9d2629bf4333e650a54f32a4da7cf86bf6785e/solidity/test/data/deposit-sweep.ts

	// https://blockstream.info/testnet/api/tx/c580e0e352570d90e303d912a506055ceeb0ee06f97dce6988c69941374f5479
	"input: P2PKH, output: P2SH, P2WPKH": {
		hashFromString("c580e0e352570d90e303d912a506055ceeb0ee06f97dce6988c69941374f5479"),
		2135049,
		bitcoin.Transaction{
			Version: 1,
			Inputs: []*bitcoin.TransactionInput{
				{
					Outpoint: &bitcoin.TransactionOutpoint{
						TransactionHash: hashFromString("e788a344a86f7e369511fe37ebd1d74686dde694ee99d06db5db3d4a14719b1d"),
						OutputIndex:     1,
					},
					SignatureScript: decodeString("47304402206f8553c07bcdc0c3b906311888103d623ca9096ca0b28b7d04650a029a01fcf9022064cda02e39e65ace712029845cfcf58d1b59617d753c3fd3556f3551b609bbb00121039d61d62dcd048d3f8550d22eb90b4af908db60231d117aeede04e7bc11907bfa"),
					Sequence:        4294967295,
				},
			},
			Outputs: []*bitcoin.TransactionOutput{
				{
					PublicKeyScript: decodeString("a9143ec459d0f3c29286ae5df5fcc421e2786024277e87"),
					Value:           20000,
				},
				{
					PublicKeyScript: decodeString("0014e257eccafbc07c381642ce6e7e55120fb077fbed"),
					Value:           1360550,
				},
			},
			Locktime: 0,
		},
	},
	// https://blockstream.info/testnet/api/tx/f5b9ad4e8cd5317925319ebc64dc923092bef3b56429c6b1bc2261bbdc73f351
	"input: P2SH, output: P2WPKH": {
		hashFromString("f5b9ad4e8cd5317925319ebc64dc923092bef3b56429c6b1bc2261bbdc73f351"),
		2135502,
		bitcoin.Transaction{
			Version: 1,
			Inputs: []*bitcoin.TransactionInput{
				{
					Outpoint: &bitcoin.TransactionOutpoint{
						TransactionHash: hashFromString("c580e0e352570d90e303d912a506055ceeb0ee06f97dce6988c69941374f5479"),
						OutputIndex:     0,
					},
					SignatureScript: decodeString("47304402205eff3ae003a5903eb33f32737e3442b6516685a1addb19339c2d02d400cf67ce0220707435fc2a0577373c63c99d242c30bea5959ec180169978d43ece50618fe0ff012103989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d94c5c14934b98637ca318a4d6e7ca6ffd1690b8e77df6377508f9f0c90d000395237576a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a914e257eccafbc07c381642ce6e7e55120fb077fbed8804e0250162b175ac68"),
					Sequence:        4294967295,
				},
			},
			Outputs: []*bitcoin.TransactionOutput{
				{
					PublicKeyScript: decodeString("00148db50eb52063ea9d98b3eac91489a90f738986f6"),
					Value:           18500,
				},
			},
			Locktime: 0,
		},
	},
	// https://blockstream.info/testnet/api/tx/c1082c460527079a84e39ec6481666db72e5a22e473a78db03b996d26fd1dc83
	"input: P2WPKH, output: P2WSH + P2WPKH": {
		hashFromString("c1082c460527079a84e39ec6481666db72e5a22e473a78db03b996d26fd1dc83"),
		2137779,
		bitcoin.Transaction{
			Version: 1,
			Inputs: []*bitcoin.TransactionInput{
				{
					Outpoint: &bitcoin.TransactionOutpoint{
						TransactionHash: hashFromString("e2131bdd5017d078ec2c17d463c9bc17abf79a9c8a37746f032b2d48ac2ff189"),
						OutputIndex:     1,
					},
					Sequence:        4294967295,
					SignatureScript: []byte{},
					Witness: [][]byte{
						decodeString("304402205e28ad48e4b128ce8b30dae8c98c8422a5a1e9aa079c0aa9d21cae999831851d02204603961ea369acfdff28a5fee1b095a9ee6a338d5c13cf8775023418e1e7c4d801"),
						decodeString("02ee067a0273f2e3ba88d23140a24fdb290f27bbcd0f94117a9c65be3911c5c04e"),
					},
				},
			},
			Outputs: []*bitcoin.TransactionOutput{
				{
					PublicKeyScript: decodeString("0020ef0b4d985752aa5ef6243e4c6f6bebc2a007e7d671ef27d4b1d0db8dcc93bc1c"),
					Value:           80000,
				},
				{
					PublicKeyScript: decodeString("00147ac2d9378a1c47e589dfb8095ca95ed2140d2726"),
					Value:           2741370,
				},
			},
			Locktime: 0,
		},
	},
	// https://blockstream.info/testnet/api/tx/9efc9d555233e12e06378a35a7b988d54f7043b5c3156adc79c7af0a0fd6f1a0
	"input: P2WSH, output: P2WPKH": {
		hashFromString("9efc9d555233e12e06378a35a7b988d54f7043b5c3156adc79c7af0a0fd6f1a0"),
		2137780,
		bitcoin.Transaction{
			Version: 1,
			Inputs: []*bitcoin.TransactionInput{
				{
					Outpoint: &bitcoin.TransactionOutpoint{
						TransactionHash: hashFromString("c1082c460527079a84e39ec6481666db72e5a22e473a78db03b996d26fd1dc83"),
						OutputIndex:     0,
					},
					Sequence:        4294967295,
					SignatureScript: []byte{},
					Witness: [][]byte{
						decodeString("3045022100bcb5b2fa3fab8d24d5ef4f601d6bc0374319162b0f534e905ffaec7abee1c69902202c25189466157797cdc5ec5049f7a2122afb89be49172f3b8c176a0bc6caf02801"),
						decodeString("03989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d9"),
						decodeString("14f4292022f75add9b079b0573d0fd63c376a85f417508b0bb0e4d6083951d7576a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a914056514a7032b0b486e56a607fb434756c61d1f74880438421962b175ac68"),
					},
				},
			},
			Outputs: []*bitcoin.TransactionOutput{
				{
					PublicKeyScript: decodeString("00148db50eb52063ea9d98b3eac91489a90f738986f6"),
					Value:           78000,
				},
			},
			Locktime: 0,
		},
	},
	// https://blockstream.info/testnet/api/tx/4459881f4964ee08dd298a12dfc1f461bf35cca8a105974d8baf0955c830d836
	"multiple inputs": {
		hashFromString("4459881f4964ee08dd298a12dfc1f461bf35cca8a105974d8baf0955c830d836"),
		2137896,
		bitcoin.Transaction{
			Version: 1,
			Inputs: []*bitcoin.TransactionInput{
				{
					Outpoint: &bitcoin.TransactionOutpoint{
						TransactionHash: hashFromString("2a5d5f472e376dc28964e1b597b1ca5ee5ac042101b5199a3ca8dae2deec3538"),
						OutputIndex:     0,
					},
					Sequence:        4294967295,
					SignatureScript: []byte{},
					Witness: [][]byte{
						decodeString("3045022100cdd1df1d2a4e15fa6824dc7a028fc0613af78fb40e2174abea22317ea5f69bcc02206dec476a49ed4e7ac900a924ef9b424f06c7d800ec15d126c0280fa5aa6535a201"),
						decodeString("03989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d9"),
					},
				},
				{
					Outpoint: &bitcoin.TransactionOutpoint{
						TransactionHash: hashFromString("71b13c7b1e2968f869c832ccdb72bbdccd35d64b78826d251d350d79a7a32f30"),
						OutputIndex:     0,
					},
					Sequence:        4294967295,
					SignatureScript: []byte{},
					Witness: [][]byte{
						decodeString("30450221009494cfbe0cd015182c05be8618fd144e4cd6db7ba9adea3909720741d530ca9502207bb2637c066af408ea0feb8021858741e542c05407322f2cd3a4703305e5bd0501"),
						decodeString("03989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d9"),
						decodeString("14208ff63189df8749780917cb5901183075dbabc175088bdbb150483eb2f27576a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a91473f3252d5e6b9f501dfafbfbca40836cc1f505f78804b80f1762b175ac68"),
					},
				},
				{
					Outpoint: &bitcoin.TransactionOutpoint{
						TransactionHash: hashFromString("68f4041f6bbddb146f672d31e4a2cce6431e1583bb24a33a2c836a7f238625d3"),
						OutputIndex:     0,
					},
					Sequence:        4294967295,
					SignatureScript: decodeString("483045022100afeb157db4284ab218a3d27b6962aabe1905eb205c6c6216dfad7e76615c0bb702205ffd88f2d2dea7509b7ea3b01910002544a785efa93c7ecd1cabafbdec508d3f012103989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d94c5c1435d54bc29e0a5170c3ac73e64c7fa539a867f0fe7508dfe75a3a6ed52db67576a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a91411d6c57c31ea78b48020dcbf42c34ccd60d92c8c880428531862b175ac68"),
					Witness:         [][]byte{},
				},
				{
					Outpoint: &bitcoin.TransactionOutpoint{
						TransactionHash: hashFromString("468e0be44cf5b2a529f22c49d8006fb29a147a4f1b6a54326a8c181208560ec6"),
						OutputIndex:     0,
					},
					Sequence:        4294967295,
					SignatureScript: decodeString("47304402200abefbc8d4d6bbe668c97ee305fde12f3c6c796ab6fbf84f00289ad5910ed8ac02200b81dcd12d45a83237569d53bcc629db559ce8c2cfd62d11fe5c58d501f785e0012103989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d94c5c142219eac966fbc0454c4a2e122717e4429dd7608f7508251c7239917eae297576a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a914032a5188c34f2fb56a4228b2bb2b7165a797eb95880488c61762b175ac68"),
					Witness:         [][]byte{},
				},
				{
					Outpoint: &bitcoin.TransactionOutpoint{
						TransactionHash: hashFromString("8c535793b98f1dbd638773e7ee07ebbbc5f86a55b5ae31ba91f63a67682e95aa"),
						OutputIndex:     0,
					},
					Sequence:        4294967295,
					SignatureScript: []byte{},
					Witness: [][]byte{
						decodeString("3045022100be74b99f0b3a616ee650a980a536ad4ba08d121ea11f15d7f51445347105dad102201f5c5becb32d2545839554fe1076fb4e6911f225f136b17232aad022fb4a5cd901"),
						decodeString("03989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d9"),
						decodeString("14462418b7495561bf2872a0786109a11f5d494aa27508eca429ef209bf5007576a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a91446c5760250ab89b3d4b956cee325561fa7effff888046c4b1862b175ac68"),
					},
				},
				{
					Outpoint: &bitcoin.TransactionOutpoint{
						TransactionHash: hashFromString("85eb466ed605916ea764860ceda68fa05e7448cc772558c866a409366b997a85"),
						OutputIndex:     0,
					},
					Sequence:        4294967295,
					SignatureScript: []byte{},
					Witness: [][]byte{
						decodeString("3045022100d94df77c599c3b443203735c966396ded29db08f3538ad60a50dc7c2c0d685f802205a3d7e5c0534a4aeb6d9a4fad4133abfa465dd814e9ac1e27d12eaffe0c6963a01"),
						decodeString("03989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d9"),
						decodeString("147f62cdde8a86328d63b9517bc70b255017f25eea75081d5c0a1bc9528ea27576a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a91464c2b58db5259ecc3c169b76c6bd83f3a94210908804e8fb1862b175ac68"),
					},
				},
			},
			Outputs: []*bitcoin.TransactionOutput{
				{
					PublicKeyScript: decodeString("00148db50eb52063ea9d98b3eac91489a90f738986f6"),
					Value:           4145001,
				},
			},
			Locktime: 0,
		},
	},
}

// TxMerkleProof holds details of the transaction Merkle proof data used as a
// test vector.
var TxMerkleProof = struct {
	TxHash      bitcoin.Hash
	BlockHeigh  uint
	MerkleProof bitcoin.TransactionMerkleProof
}{
	// https://blockstream.info/testnet/api/tx/72e7fd57c2adb1ed2305c4247486ff79aec363296f02ec65be141904f80d214e
	TxHash: hashFromString(
		"72e7fd57c2adb1ed2305c4247486ff79aec363296f02ec65be141904f80d214e",
	),
	BlockHeigh: 1569342,
	MerkleProof: bitcoin.TransactionMerkleProof{
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
}

var AssembleProof = map[string]struct {
	RequiredConfirmations uint
	BitcoinChainData      struct {
		TransactionHash            bitcoin.Hash
		Transaction                bitcoin.Transaction
		AccumulatedTxConfirmations uint
		HeadersChain               map[uint]*bitcoin.BlockHeader
		TransactionMerkleProof     *bitcoin.TransactionMerkleProof
	}
	ExpectedProof       *bitcoin.Proof
	ExpectedTransaction bitcoin.Transaction
}{
	"single input": {
		RequiredConfirmations: 6,
		BitcoinChainData: struct {
			TransactionHash            bitcoin.Hash
			Transaction                bitcoin.Transaction
			AccumulatedTxConfirmations uint
			HeadersChain               map[uint]*bitcoin.BlockHeader
			TransactionMerkleProof     *bitcoin.TransactionMerkleProof
		}{
			TransactionHash: hashFromString(
				"44c568bc0eac07a2a9c2b46829be5b5d46e7d00e17bfb613f506a75ccf86a473",
			),
			Transaction: bitcoin.Transaction{
				Version: 1, // TODO: provide proper version
				Inputs: []*bitcoin.TransactionInput{
					{
						Outpoint: &bitcoin.TransactionOutpoint{
							TransactionHash: hashFromString(
								"8ee67b585eeb682bf6907ea311282540ee53edf605e0f09757226a4dc3e72a67",
							),
							OutputIndex: 0,
						},
						SignatureScript: decodeString(""),
					},
				},
				Outputs: []*bitcoin.TransactionOutput{
					{
						Value: 8400,
						PublicKeyScript: decodeString(
							"00148db50eb52063ea9d98b3eac91489a90f738986f6",
						),
					},
				},
			},
			AccumulatedTxConfirmations: 7,
			HeadersChain: map[uint]*bitcoin.BlockHeader{
				2164152: &bitcoin.BlockHeader{
					Version: 536928260,
					PreviousBlockHeaderHash: hashFromStringInternal(
						"732d33ea35d62f9488cff5d64c0d702afd5d88092230ddfcc45f000000000000",
					),
					MerkleRootHash: hashFromStringInternal(
						"196283ba24a3f5bad91ef95338aa6d214c934f2c1392e39a0447377fe5b0a04b",
					),
					Time:  1646051559,
					Bits:  486604799,
					Nonce: 655015664,
				},
				2164153: &bitcoin.BlockHeader{
					Version: 536870916,
					PreviousBlockHeaderHash: hashFromStringInternal(
						"6c318b23e5c42e86ef3edd080e50c9c233b9f0b6d186bd57e413000000000000",
					),
					MerkleRootHash: hashFromStringInternal(
						"21fb8cda200bff4fec1338d85a1e005bb4d729d908a7c5c232ecd0713231d044",
					),
					Time:  1646051678,
					Bits:  436420333,
					Nonce: 1850098555,
				},
				2164154: &bitcoin.BlockHeader{
					Version: 536928260,
					PreviousBlockHeaderHash: hashFromStringInternal(
						"f416898d79d4a46fa6c54f190ad3d502bad8aa3afdec0714aa00000000000000",
					),
					MerkleRootHash: hashFromStringInternal(
						"0603a5cc15e5906cb4eac9f747869fdc9be856e76a110b4f87da90db20f9fbe2",
					),
					Time:  1646051727,
					Bits:  436420333,
					Nonce: 3687046933,
				},
				2164155: &bitcoin.BlockHeader{
					Version: 536870916,
					PreviousBlockHeaderHash: hashFromStringInternal(
						"642125b3910fdaead521b57955e28893d89f8ce7fd3ba1dd6d01000000000000",
					),
					MerkleRootHash: hashFromStringInternal(
						"f9e17a266a2267ee02d5ab82a75a76805db821a13abd2e80e0950d883311e535",
					),
					Time:  1646051933,
					Bits:  436420333,
					Nonce: 3288530142,
				},
				2164156: &bitcoin.BlockHeader{
					Version: 536870916,
					PreviousBlockHeaderHash: hashFromStringInternal(
						"5b6de55e069be71b21a62cd140dc7031225f7258dc758f19ea01000000000000",
					),
					MerkleRootHash: hashFromStringInternal(
						"139966d27d9ed0c0c1ed9162c2fea2ccf0ba212706f6bc421d0a2b6211de040d"),
					Time:  1646052378,
					Bits:  436420333,
					Nonce: 2404591175,
				},
				2164157: &bitcoin.BlockHeader{
					Version: 536928260,
					PreviousBlockHeaderHash: hashFromStringInternal(
						"8475e15e0314635d32abf04c761fee528d6a3f2db3b3d1379800000000000000",
					),
					MerkleRootHash: hashFromStringInternal(
						"2a3fa06fecd9dd4bf2e25e22a95d4f65435d5c5b42bcf498b4e756f9f4ea67ce",
					),
					Time:  1646052769,
					Bits:  436420333,
					Nonce: 2901638045,
				},
				2164158: &bitcoin.BlockHeader{
					Version: 536870912,
					PreviousBlockHeaderHash: hashFromStringInternal(
						"3f16d450c51853a4cd9569d225028aa08ab6139eee31f4f67a01000000000000",
					),
					MerkleRootHash: hashFromStringInternal(
						"4cda79bc48b970de2fb29c3f38626eb9d70d8bae7b92aad09f2a0ad2d2f334d3",
					),
					Time:  1646053979,
					Bits:  486604799,
					Nonce: 398626564,
				},
			},
			TransactionMerkleProof: &bitcoin.TransactionMerkleProof{
				BlockHeight: 2164152,
				MerkleNodes: []string{
					"7bffaff2c61291861276da41cf6c3842fad555af97dd1ff98ce41c61a0072b12",
					"7a5876ddee8e553ff0650c739b2ec66e192d8afe5fc0ce763bf810457aea330c",
					"2d17b67d5519bc39fbef8650afd3fe11fdfb3f471434a5b551cfa9a41441901f",
					"1376d102b677591ce2fa62553e2a57ab5919022b03036521facfce93a0338026",
					"43ad3aadad675e398c59eb846a8e037cf7de8ba3b38f3388175f25d84b777c80",
					"6969c227128793b3c9e99c05f20fb9b91fdb73458fd53151b5fe29d30c10cf9a",
					"0a76bc4d8c3d532357be4d188ba89e9ae364a7d3c365e690e3cb07359b86129c",
				},
				Position: 11,
			},
		},
		ExpectedProof: &bitcoin.Proof{
			MerkleProof: decodeString(
				"122b07a0611ce48cf91fdd97af55d5fa42386ccf41da7612869112c6f2afff7b0c" +
					"33ea7a4510f83b76cec05ffe8a2d196ec62e9b730c65f03f558eeedd76587a1f90" +
					"4114a4a9cf51b5a53414473ffbfd11fed3af5086effb39bc19557db6172d268033" +
					"a093cecffa216503032b021959ab572a3e5562fae21c5977b602d17613807c774b" +
					"d8255f1788338fb3a38bdef77c038e6a84eb598c395e67adad3aad439acf100cd3" +
					"29feb55131d58f4573db1fb9b90ff2059ce9c9b393871227c269699c12869b3507" +
					"cbe390e665c3d3a764e39a9ea88b184dbe5723533d8c4dbc760a",
			),
			TxIndexInBlock: 11,
			BitcoinHeaders: decodeString(
				"04e00020732d33ea35d62f9488cff5d64c0d702afd5d88092230ddfcc45f000000" +
					"000000196283ba24a3f5bad91ef95338aa6d214c934f2c1392e39a0447377fe5b0" +
					"a04be7c01c62ffff001df0be0a27040000206c318b23e5c42e86ef3edd080e50c9" +
					"c233b9f0b6d186bd57e41300000000000021fb8cda200bff4fec1338d85a1e005b" +
					"b4d729d908a7c5c232ecd0713231d0445ec11c62ed3e031a7b43466e04e00020f4" +
					"16898d79d4a46fa6c54f190ad3d502bad8aa3afdec0714aa000000000000000603" +
					"a5cc15e5906cb4eac9f747869fdc9be856e76a110b4f87da90db20f9fbe28fc11c" +
					"62ed3e031a15dfc3db04000020642125b3910fdaead521b57955e28893d89f8ce7" +
					"fd3ba1dd6d01000000000000f9e17a266a2267ee02d5ab82a75a76805db821a13a" +
					"bd2e80e0950d883311e5355dc21c62ed3e031adefc02c4040000205b6de55e069b" +
					"e71b21a62cd140dc7031225f7258dc758f19ea01000000000000139966d27d9ed0" +
					"c0c1ed9162c2fea2ccf0ba212706f6bc421d0a2b6211de040d1ac41c62ed3e031a" +
					"4726538f04e000208475e15e0314635d32abf04c761fee528d6a3f2db3b3d13798" +
					"000000000000002a3fa06fecd9dd4bf2e25e22a95d4f65435d5c5b42bcf498b4e7" +
					"56f9f4ea67cea1c51c62ed3e031a9d7bf3ac",
			),
		},
	},
	"multiple inputs": {
		RequiredConfirmations: 6,
		BitcoinChainData: struct {
			TransactionHash            bitcoin.Hash
			Transaction                bitcoin.Transaction
			AccumulatedTxConfirmations uint
			HeadersChain               map[uint]*bitcoin.BlockHeader
			TransactionMerkleProof     *bitcoin.TransactionMerkleProof
		}{
			TransactionHash: hashFromString(
				"5083822ed0b8d0bc661362b778e666cb572ff6d5152193992dd69d3207995753",
			),
			Transaction: bitcoin.Transaction{
				Version: 1, // TODO: provide proper version
				Inputs: []*bitcoin.TransactionInput{
					{
						Outpoint: &bitcoin.TransactionOutpoint{
							TransactionHash: hashFromString(
								"ea4d9e45f8c1b8a187c007f36ba1e9b201e8511182c7083c4edcaf9325b2998f",
							),
							OutputIndex: 0,
						},
						SignatureScript: decodeString(""),
					},
					{
						Outpoint: &bitcoin.TransactionOutpoint{
							TransactionHash: hashFromString(
								"c844ff4c1781c884bb5e80392398b81b984d7106367ae16675f132bd1a7f33fd",
							),
							OutputIndex: 0,
						},
						SignatureScript: decodeString(""),
					}, {
						Outpoint: &bitcoin.TransactionOutpoint{
							TransactionHash: hashFromString(
								"44c568bc0eac07a2a9c2b46829be5b5d46e7d00e17bfb613f506a75ccf86a473",
							),
							OutputIndex: 0,
						},
						SignatureScript: decodeString(""),
					}, {
						Outpoint: &bitcoin.TransactionOutpoint{
							TransactionHash: hashFromString(
								"f548c00e464764e112826450a00cf005ca771a6108a629b559b6c60a519e4378",
							),
							OutputIndex: 0,
						},
						SignatureScript: decodeString(""),
					},
				},
				Outputs: []*bitcoin.TransactionOutput{
					{
						Value: 39800,
						PublicKeyScript: decodeString(
							"00148db50eb52063ea9d98b3eac91489a90f738986f6",
						),
					},
				},
			},
			AccumulatedTxConfirmations: 6,
			HeadersChain: map[uint]*bitcoin.BlockHeader{
				2164155: &bitcoin.BlockHeader{
					Version: 536870916,
					PreviousBlockHeaderHash: hashFromStringInternal(
						"642125b3910fdaead521b57955e28893d89f8ce7fd3ba1dd6d01000000000000",
					),
					MerkleRootHash: hashFromStringInternal(
						"f9e17a266a2267ee02d5ab82a75a76805db821a13abd2e80e0950d883311e535",
					),
					Time:  1646051933,
					Bits:  436420333,
					Nonce: 3288530142,
				},
				2164156: &bitcoin.BlockHeader{
					Version: 536870916,
					PreviousBlockHeaderHash: hashFromStringInternal(
						"5b6de55e069be71b21a62cd140dc7031225f7258dc758f19ea01000000000000",
					),
					MerkleRootHash: hashFromStringInternal(
						"139966d27d9ed0c0c1ed9162c2fea2ccf0ba212706f6bc421d0a2b6211de040d",
					),
					Time:  1646052378,
					Bits:  436420333,
					Nonce: 2404591175,
				},
				2164157: &bitcoin.BlockHeader{
					Version: 536928260,
					PreviousBlockHeaderHash: hashFromStringInternal(
						"8475e15e0314635d32abf04c761fee528d6a3f2db3b3d1379800000000000000",
					),
					MerkleRootHash: hashFromStringInternal(
						"2a3fa06fecd9dd4bf2e25e22a95d4f65435d5c5b42bcf498b4e756f9f4ea67ce",
					),
					Time:  1646052769,
					Bits:  436420333,
					Nonce: 2901638045,
				},
				2164158: &bitcoin.BlockHeader{
					Version: 536870912,
					PreviousBlockHeaderHash: hashFromStringInternal(
						"3f16d450c51853a4cd9569d225028aa08ab6139eee31f4f67a01000000000000",
					),
					MerkleRootHash: hashFromStringInternal(
						"4cda79bc48b970de2fb29c3f38626eb9d70d8bae7b92aad09f2a0ad2d2f334d3",
					),
					Time:  1646053979,
					Bits:  486604799,
					Nonce: 398626564,
				},
				2164159: &bitcoin.BlockHeader{
					Version: 536870912,
					PreviousBlockHeaderHash: hashFromStringInternal(
						"687e487acbf5eb375c631a15127fbf7d80ca084461e7f26f92c509b600000000",
					),
					MerkleRootHash: hashFromStringInternal(
						"6fad33bd7c8d651bd6dc86c286f0a99340b668f019b9e97a59fd392c36c4f469",
					),
					Time:  1646055184,
					Bits:  486604799,
					Nonce: 2863431488,
				},
				2164160: &bitcoin.BlockHeader{
					Version: 536870916,
					PreviousBlockHeaderHash: hashFromStringInternal(
						"40f4c65610f26f06c4365305b956934501713e01c2fc08b919e0bc1b00000000",
					),
					MerkleRootHash: hashFromStringInternal(
						"e401a6a884ba015e83c6fe2cd363e877ef03982e81eaff4e2c95af1e23a670f4",
					),
					Time:  1646056455,
					Bits:  486604799,
					Nonce: 407750232,
				},
			},
			TransactionMerkleProof: &bitcoin.TransactionMerkleProof{
				BlockHeight: 2164155,
				MerkleNodes: []string{
					"322cfdf3ca53cf597b6f08e93489b9a1cfa1f5958c3657474b0d8f5efb5ca92e",
					"82aedffef6c9670375effee25740fecce143d21f8abf98307235b7ebd31ad4d1",
					"837fa041b9a8f5b42353fdf8981e3b7a78c61858852e43058bfe6cacf9eab5a3",
					"a51612d3f3f857e95803a4d86aa6dbbe2e756dc2ed6cc0e04630e8baf597e377",
					"a00501650e0c4f8a1e07a5d6d5bc5e75e4c75de61a65f0410cce354bbae78686",
				},
				Position: 6,
			},
		},
		ExpectedProof: &bitcoin.Proof{
			MerkleProof: decodeString(
				"2ea95cfb5e8f0d4b4757368c95f5a1cfa1b98934e9086f7b59cf53caf3fd2c32d1d" +
					"41ad3ebb735723098bf8a1fd243e1ccfe4057e2feef750367c9f6fedfae82a3b5ea" +
					"f9ac6cfe8b05432e855818c6787a3b1e98f8fd5323b4f5a8b941a07f8377e397f5b" +
					"ae83046e0c06cedc26d752ebedba66ad8a40358e957f8f3d31216a58686e7ba4b35" +
					"ce0c41f0651ae65dc7e4755ebcd5d6a5071e8a4f0c0e650105a0",
			),
			TxIndexInBlock: 6,
			BitcoinHeaders: decodeString(
				"04000020642125b3910fdaead521b57955e28893d89f8ce7fd3ba1dd6d010000000" +
					"00000f9e17a266a2267ee02d5ab82a75a76805db821a13abd2e80e0950d883311e5" +
					"355dc21c62ed3e031adefc02c4040000205b6de55e069be71b21a62cd140dc70312" +
					"25f7258dc758f19ea01000000000000139966d27d9ed0c0c1ed9162c2fea2ccf0ba" +
					"212706f6bc421d0a2b6211de040d1ac41c62ed3e031a4726538f04e000208475e15" +
					"e0314635d32abf04c761fee528d6a3f2db3b3d13798000000000000002a3fa06fec" +
					"d9dd4bf2e25e22a95d4f65435d5c5b42bcf498b4e756f9f4ea67cea1c51c62ed3e0" +
					"31a9d7bf3ac000000203f16d450c51853a4cd9569d225028aa08ab6139eee31f4f6" +
					"7a010000000000004cda79bc48b970de2fb29c3f38626eb9d70d8bae7b92aad09f2" +
					"a0ad2d2f334d35bca1c62ffff001d048fc21700000020687e487acbf5eb375c631a" +
					"15127fbf7d80ca084461e7f26f92c509b6000000006fad33bd7c8d651bd6dc86c28" +
					"6f0a99340b668f019b9e97a59fd392c36c4f46910cf1c62ffff001d407facaa0400" +
					"002040f4c65610f26f06c4365305b956934501713e01c2fc08b919e0bc1b0000000" +
					"0e401a6a884ba015e83c6fe2cd363e877ef03982e81eaff4e2c95af1e23a670f407" +
					"d41c62ffff001d58c64d18",
			),
		},
	},
}

func hashFromString(s string) bitcoin.Hash {
	hash, err := bitcoin.NewHashFromString(
		s,
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		panic(err)
	}

	return hash
}

// TODO: Replace with the above function
func hashFromStringInternal(s string) bitcoin.Hash {
	hash, err := bitcoin.NewHashFromString(
		s,
		bitcoin.InternalByteOrder,
	)
	if err != nil {
		panic(err)
	}

	return hash
}

func decodeString(s string) []byte {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return bytes
}
