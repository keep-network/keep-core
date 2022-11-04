package bitcoin

// Transaction represents a Bitcoin transaction. For reference, see:
// https://developer.bitcoin.org/reference/transactions.html#raw-transaction-format
type Transaction struct {
	// Version is the transaction version number. Usually, version 1 or 2.
	Version int32
	// Inputs is a slice holding all inputs of the transaction.
	Inputs []*TransactionInput
	// Outputs is a slice holding all outputs of the transaction.
	Outputs []*TransactionOutput
	// Locktime is the transaction locktime being a Unix epoch time or
	// block number, interpreted according to the locktime parsing rules:
	// https://developer.bitcoin.org/devguide/transactions.html#locktime-and-sequence-number
	Locktime uint32
}

// Serialize serializes the transaction to a byte array using the traditional
// serialization format: [version][inputs][outputs][locktime].
func (t *Transaction) Serialize() []byte {
	// TODO: Implementation of the Serialize function that consists of the following
	//       (see https://hongchao.me/anatomy-of-raw-bitcoin-transaction for reference):
	//       1. Serialize t.Version to an InternalByteOrder byte array.
	//       2. Serialize t.Inputs as follows:
	//          2.1. Serialize each input separately. All numbers should be
	//               serialized to an InternalByteOrder byte array.
	//          2.2. Concatenate serialized inputs into a single array
	//               preserving the inputs order.
	//          2.3. Prepend the concatenation with its length encoded
	//               as an CompactSizeUint.
	//       3. Serialize t.Outputs just like t.Inputs.
	//       4. Serialize t.Locktime to an InternalByteOrder byte array.
	return []byte("01000000011d9b71144a3ddbb56dd099ee94e6dd8646d7d1eb37fe1195367e6fa844a388e7010000006a47304402206f8553c07bcdc0c3b906311888103d623ca9096ca0b28b7d04650a029a01fcf9022064cda02e39e65ace712029845cfcf58d1b59617d753c3fd3556f3551b609bbb00121039d61d62dcd048d3f8550d22eb90b4af908db60231d117aeede04e7bc11907bfaffffffff02204e00000000000017a9143ec459d0f3c29286ae5df5fcc421e2786024277e87a6c2140000000000160014e257eccafbc07c381642ce6e7e55120fb077fbed00000000")
}

// Hash calculates the transaction's hash as the double SHA-256 of the
// traditional serialization format: [version][inputs][outputs][locktime].
// The outcome is equivalent to the txid field defined in the Bitcoin
// specification. Do not confuse it with the wtxid field defined by BIP 141.
// For reference, see:
// https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki#transaction-id
func (t *Transaction) Hash() Hash {
	// TODO: Implementation of the Hash function that consists of the following:
	//       1. Call t.Serialize() to get the serialized transaction.
	//       2. Compute the double SHA-256 over the serialized transaction.
	//       3. Construct the Hash instance appropriately.
	return Hash{}
}

// TransactionOutpoint represents a Bitcoin transaction outpoint.
// For reference, see:
// https://developer.bitcoin.org/reference/transactions.html#outpoint-the-specific-part-of-a-specific-output
type TransactionOutpoint struct {
	// TransactionHash is the hash of the transaction holding the output
	// to spend.
	TransactionHash Hash
	// OutputIndex is the zero-based index of the output to spend from the
	// specified transaction.
	OutputIndex uint32
}

// TransactionInput represents a Bitcoin transaction input. For reference, see:
// https://developer.bitcoin.org/reference/transactions.html#txin-a-transaction-input-non-coinbase
type TransactionInput struct {
	// Outpoint is the previous transaction outpoint being spent by this input.
	Outpoint *TransactionOutpoint
	// SignatureScript is a script-language script that satisfies the conditions
	// placed in the outpoint's public key script (see TransactionOutput.PublicKeyScript).
	// This slice must start with the byte-length of the script encoded as a
	// CompactSizeUint. This field is not set (nil or empty) for SegWit
	// transactions.
	SignatureScript []byte
	// Sequence is the sequence number for this input. Default value
	// is 0xffffffff. For reference, see:
	// https://developer.bitcoin.org/devguide/transactions.html#locktime-and-sequence-number
	Sequence uint32
}

// TransactionOutput represents a Bitcoin transaction output. For reference, see:
// https://developer.bitcoin.org/reference/transactions.html#txout-a-transaction-output
type TransactionOutput struct {
	// Value denotes the number of satoshis to spend. Zero is a valid value.
	Value int64
	// PublicKeyScript defines the conditions that must be satisfied to spend
	// this output. This slice must start with the byte-length of the script
	// encoded as CompactSizeUint
	PublicKeyScript []byte
}

// UnspentTransactionOutput represents an unspent output (UTXO) of a Bitcoin
// transaction.
type UnspentTransactionOutput struct {
	// Outpoint is the transaction outpoint this UTXO points to.
	Outpoint *TransactionOutpoint
	// Value denotes the number of unspent satoshis.
	Value int64
}
