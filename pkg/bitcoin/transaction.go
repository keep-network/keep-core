package bitcoin

const (
	// TransactionVersion is the current latest supported transaction version.
	TransactionVersion = 1
	// MaxTransactionInputSequence is the maximum sequence number the sequence
	// field of a transaction input can be.
	MaxTransactionInputSequence uint32 = 0xffffffff
)

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

// NewTransaction constructs an empty Transaction instance.
func NewTransaction() *Transaction {
	return &Transaction{
		Version:  TransactionVersion,
		Inputs:   make([]*TransactionInput, 0),
		Outputs:  make([]*TransactionOutput, 0),
		Locktime: 0,
	}
}

// AddInput adds a new unsigned TransactionInput pointing to the provided
// TransactionOutpoint.
func (t *Transaction) AddInput(outpoint *TransactionOutpoint) {
	t.Inputs = append(t.Inputs, &TransactionInput{
		Outpoint:        outpoint,
		SignatureScript: nil,
		Witness:         nil,
		Sequence:        MaxTransactionInputSequence,
	})
}

// AddOutput adds a new TransactionOutput of the given value and locked
// using the provided publicKeyScript.
func (t *Transaction) AddOutput(publicKeyScript []byte, value int64) {
	t.Outputs = append(t.Outputs, &TransactionOutput{
		Value:           value,
		PublicKeyScript: publicKeyScript,
	})
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
	return nil
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
	// This slice MUST NOT start with the byte-length of the script encoded as
	// a CompactSizeUint as this is done during transaction serialization.
	// This field is not set (nil or empty) for SegWit transaction inputs.
	// That means it is mutually exclusive with the below Witness field.
	SignatureScript []byte
	// Witness holds the witness data for the given input. It should be interpreted
	// as a stack with one or many elements. Individual elements MUST NOT start
	// with the byte-length of the script encoded as a CompactSizeUint as this
	// is done during transaction serialization. This field is not set
	// (nil or empty) for non-SegWit transaction inputs. That means it is
	// mutually exclusive with the above SignatureScript field.
	Witness [][]byte
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
	// this output. This slice MUST NOT start with the byte-length of the script
	// encoded as CompactSizeUint as this is done during transaction serialization.
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
