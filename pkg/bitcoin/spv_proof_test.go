package bitcoin

import (
	"reflect"
	"testing"

	"encoding/hex"
)

// SpvProofData holds details of the transaction proof data used as a test
// vector.
var SpvProofData = map[string]struct {
	RequiredConfirmations uint
	BitcoinChainData      struct {
		TransactionHex             string
		AccumulatedTxConfirmations uint
		HeadersChain               map[uint]*BlockHeader
		TransactionMerkleProof     *TransactionMerkleProof
	}
	ExpectedProof       *SpvProof
	ExpectedTransaction Transaction
}{
	// https://blockstream.info/testnet/api/tx/44c568bc0eac07a2a9c2b46829be5b5d46e7d00e17bfb613f506a75ccf86a473
	"single input": {
		RequiredConfirmations: 6,
		BitcoinChainData: struct {
			TransactionHex             string
			AccumulatedTxConfirmations uint
			HeadersChain               map[uint]*BlockHeader
			TransactionMerkleProof     *TransactionMerkleProof
		}{
			TransactionHex:             "01000000000101672ae7c34d6a225797f0e005f6ed53ee40252811a37e90f62b68eb5e587be68e0000000000ffffffff01d0200000000000001600148db50eb52063ea9d98b3eac91489a90f738986f603483045022100b12afadf68ad9781600f065e0b09e22058ca2293aa86ac38add3ca7cfb01b3b7022009ecce0c1c3ebd26569c6b0d60e15b4675860737487d1b7c782439acf4709bdf012103989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d95c14934b98637ca318a4d6e7ca6ffd1690b8e77df6377508f9f0c90d000395237576a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a914e257eccafbc07c381642ce6e7e55120fb077fbed8804e0250162b175ac6800000000",
			AccumulatedTxConfirmations: 7,
			HeadersChain: map[uint]*BlockHeader{
				2164152: {
					Version: 536928260,
					PreviousBlockHeaderHash: hashFromString(
						"0000000000005fc4fcdd302209885dfd2a700d4cd6f5cf88942fd635ea332d73",
					),
					MerkleRootHash: hashFromString(
						"4ba0b0e57f3747049ae392132c4f934c216daa3853f91ed9baf5a324ba836219",
					),
					Time:  1646051559,
					Bits:  486604799,
					Nonce: 655015664,
				},
				2164153: {
					Version: 536870916,
					PreviousBlockHeaderHash: hashFromString(
						"00000000000013e457bd86d1b6f0b933c2c9500e08dd3eef862ec4e5238b316c",
					),
					MerkleRootHash: hashFromString(
						"44d0313271d0ec32c2c5a708d929d7b45b001e5ad83813ec4fff0b20da8cfb21",
					),
					Time:  1646051678,
					Bits:  436420333,
					Nonce: 1850098555,
				},
				2164154: {
					Version: 536928260,
					PreviousBlockHeaderHash: hashFromString(
						"00000000000000aa1407ecfd3aaad8ba02d5d30a194fc5a66fa4d4798d8916f4",
					),
					MerkleRootHash: hashFromString(
						"e2fbf920db90da874f0b116ae756e89bdc9f8647f7c9eab46c90e515cca50306",
					),
					Time:  1646051727,
					Bits:  436420333,
					Nonce: 3687046933,
				},
				2164155: {
					Version: 536870916,
					PreviousBlockHeaderHash: hashFromString(
						"000000000000016ddda13bfde78c9fd89388e25579b521d5eada0f91b3252164",
					),
					MerkleRootHash: hashFromString(
						"35e51133880d95e0802ebd3aa121b85d80765aa782abd502ee67226a267ae1f9",
					),
					Time:  1646051933,
					Bits:  436420333,
					Nonce: 3288530142,
				},
				2164156: {
					Version: 536870916,
					PreviousBlockHeaderHash: hashFromString(
						"00000000000001ea198f75dc58725f223170dc40d12ca6211be79b065ee56d5b",
					),
					MerkleRootHash: hashFromString(
						"0d04de11622b0a1d42bcf6062721baf0cca2fec26291edc1c0d09e7dd2669913",
					),
					Time:  1646052378,
					Bits:  436420333,
					Nonce: 2404591175,
				},
				2164157: {
					Version: 536928260,
					PreviousBlockHeaderHash: hashFromString(
						"000000000000009837d1b3b32d3f6a8d52ee1f764cf0ab325d6314035ee17584",
					),
					MerkleRootHash: hashFromString(
						"ce67eaf4f956e7b498f4bc425b5c5d43654f5da9225ee2f24bddd9ec6fa03f2a",
					),
					Time:  1646052769,
					Bits:  436420333,
					Nonce: 2901638045,
				},
				2164158: {
					Version: 536870912,
					PreviousBlockHeaderHash: hashFromString(
						"000000000000017af6f431ee9e13b68aa08a0225d26995cda45318c550d4163f",
					),
					MerkleRootHash: hashFromString(
						"d334f3d2d20a2a9fd0aa927bae8b0dd7b96e62383f9cb22fde70b948bc79da4c",
					),
					Time:  1646053979,
					Bits:  486604799,
					Nonce: 398626564,
				},
			},
			TransactionMerkleProof: &TransactionMerkleProof{
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
		ExpectedProof: &SpvProof{
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
	// https://blockstream.info/testnet/api/tx/5083822ed0b8d0bc661362b778e666cb572ff6d5152193992dd69d3207995753
	"multiple inputs": {
		RequiredConfirmations: 6,
		BitcoinChainData: struct {
			TransactionHex             string
			AccumulatedTxConfirmations uint
			HeadersChain               map[uint]*BlockHeader
			TransactionMerkleProof     *TransactionMerkleProof
		}{
			TransactionHex:             "010000000001048f99b22593afdc4e3c08c7821151e801b2e9a16bf307c087a1b8c1f8459e4dea00000000c9483045022100bb54f2717647b2f2c5370b5f12b55e27f97a6e2009dcd21fca08527df949e1fd022058bc3cd1dd739b89b9e4cda43b13bc59cfb15663b80cbfa3edb4539107bba35d012103989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d94c5c14934b98637ca318a4d6e7ca6ffd1690b8e77df6377508f9f0c90d000395237576a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a914e257eccafbc07c381642ce6e7e55120fb077fbed8804e0250162b175ac68fffffffffd337f1abd32f17566e17a3606714d981bb8982339805ebb84c881174cff44c80000000000ffffffff73a486cf5ca706f513b6bf170ed0e7465d5bbe2968b4c2a9a207ac0ebc68c5440000000000ffffffff78439e510ac6b659b529a608611a77ca05f00ca050648212e16447460ec048f50000000000ffffffff01789b0000000000001600148db50eb52063ea9d98b3eac91489a90f738986f6000347304402205199b28a3b4a81579fe4ea99925380b298e28ca38a3b14e50f12daec87945449022065c5034f96ed785aa10b3817c501ecc59f1abf329fad07229170c3dd5f53bc91012103989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d95c14934b98637ca318a4d6e7ca6ffd1690b8e77df6377508f9f0c90d000395237576a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a914e257eccafbc07c381642ce6e7e55120fb077fbed8804e0250162b175ac680247304402201b2a3b03a1088c6bbc406e96a6017e52ce86c0897541c9bb59d94179daa84f8702204b1e665bd43bbe968e1d89b15c5f0b5551011fa4caf2fbb7eb22d89a38fad04d012103989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d903473044022007ce54f21a2f5233bd046c4600bcb1c874aaf9053b1d7ee7d47eb134b695fbed022002e8684548b7a3cdaecb8c6393244c396c15e1ebaedb53e2fcc6c5cea7310490012103989d253b17a6a0f41838b84ff0d20e8898f9d7b1a98f2564da4cc29dcf8581d95c14934b98637ca318a4d6e7ca6ffd1690b8e77df6377508f9f0c90d000395237576a9148db50eb52063ea9d98b3eac91489a90f738986f68763ac6776a914e257eccafbc07c381642ce6e7e55120fb077fbed8804e0250162b175ac6800000000",
			AccumulatedTxConfirmations: 6,
			HeadersChain: map[uint]*BlockHeader{
				2164155: {
					Version: 536870916,
					PreviousBlockHeaderHash: hashFromString(
						"000000000000016ddda13bfde78c9fd89388e25579b521d5eada0f91b3252164",
					),
					MerkleRootHash: hashFromString(
						"35e51133880d95e0802ebd3aa121b85d80765aa782abd502ee67226a267ae1f9",
					),
					Time:  1646051933,
					Bits:  436420333,
					Nonce: 3288530142,
				},
				2164156: {
					Version: 536870916,
					PreviousBlockHeaderHash: hashFromString(
						"00000000000001ea198f75dc58725f223170dc40d12ca6211be79b065ee56d5b",
					),
					MerkleRootHash: hashFromString(
						"0d04de11622b0a1d42bcf6062721baf0cca2fec26291edc1c0d09e7dd2669913",
					),
					Time:  1646052378,
					Bits:  436420333,
					Nonce: 2404591175,
				},
				2164157: {
					Version: 536928260,
					PreviousBlockHeaderHash: hashFromString(
						"000000000000009837d1b3b32d3f6a8d52ee1f764cf0ab325d6314035ee17584",
					),
					MerkleRootHash: hashFromString(
						"ce67eaf4f956e7b498f4bc425b5c5d43654f5da9225ee2f24bddd9ec6fa03f2a",
					),
					Time:  1646052769,
					Bits:  436420333,
					Nonce: 2901638045,
				},
				2164158: {
					Version: 536870912,
					PreviousBlockHeaderHash: hashFromString(
						"000000000000017af6f431ee9e13b68aa08a0225d26995cda45318c550d4163f",
					),
					MerkleRootHash: hashFromString(
						"d334f3d2d20a2a9fd0aa927bae8b0dd7b96e62383f9cb22fde70b948bc79da4c",
					),
					Time:  1646053979,
					Bits:  486604799,
					Nonce: 398626564,
				},
				2164159: {
					Version: 536870912,
					PreviousBlockHeaderHash: hashFromString(
						"00000000b609c5926ff2e7614408ca807dbf7f12151a635c37ebf5cb7a487e68",
					),
					MerkleRootHash: hashFromString(
						"69f4c4362c39fd597ae9b919f068b64093a9f086c286dcd61b658d7cbd33ad6f",
					),
					Time:  1646055184,
					Bits:  486604799,
					Nonce: 2863431488,
				},
				2164160: {
					Version: 536870916,
					PreviousBlockHeaderHash: hashFromString(
						"000000001bbce019b908fcc2013e7101459356b9055336c4066ff21056c6f440",
					),
					MerkleRootHash: hashFromString(
						"f470a6231eaf952c4effea812e9803ef77e863d32cfec6835e01ba84a8a601e4",
					),
					Time:  1646056455,
					Bits:  486604799,
					Nonce: 407750232,
				},
			},
			TransactionMerkleProof: &TransactionMerkleProof{
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
		ExpectedProof: &SpvProof{
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

func hashFromString(s string) Hash {
	hash, err := NewHashFromString(
		s,
		ReversedByteOrder,
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

func TestAssembleTransactionProof(t *testing.T) {
	for testName, test := range SpvProofData {
		t.Run(testName, func(t *testing.T) {
			transaction := transactionFrom(t, test.BitcoinChainData.TransactionHex)
			transactionHash := transaction.Hash()
			requiredConfirmations := test.RequiredConfirmations
			accumulatedConfirmations := test.BitcoinChainData.AccumulatedTxConfirmations
			blockHeaders := test.BitcoinChainData.HeadersChain
			transactionMerkleProof := test.BitcoinChainData.TransactionMerkleProof
			expectedProof := test.ExpectedProof
			expectedTx := transaction

			bitcoinChain := newLocalChain()
			bitcoinChain.addTransaction(transaction)
			bitcoinChain.addTransactionConfirmations(
				transactionHash,
				accumulatedConfirmations,
			)
			for blockNumber, blockHeader := range blockHeaders {
				bitcoinChain.addBlockHeader(blockNumber, blockHeader)
			}
			bitcoinChain.addTransactionMerkleProof(
				transactionHash,
				transactionMerkleProof,
			)

			tx, proof, err := AssembleSpvProof(
				transactionHash,
				requiredConfirmations,
				bitcoinChain,
			)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(expectedProof, proof) {
				t.Errorf(
					"unexpected proof\nexpected: %v\nactual:   %v\n",
					expectedProof,
					proof,
				)
			}
			if !reflect.DeepEqual(expectedTx, tx) {
				t.Errorf(
					"unexpected transaction\nexpected: %v\nactual:   %v\n",
					expectedTx,
					tx,
				)
			}
		})
	}
}
