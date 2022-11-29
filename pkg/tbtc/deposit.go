package tbtc

import (
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
)

// deposit represents a tBTC deposit.
type deposit struct {
	// utxo is the unspent output of the deposit funding transaction that
	// represents the deposit on the Bitcoin chain.
	utxo *bitcoin.UnspentTransactionOutput
	// depositor is the depositor's address on the host chain.
	depositor chain.Address
	// blindingFactor is an 8-byte arbitrary value that allows to distinguish
	// deposits from the same depositor.
	blindingFactor [8]byte
	// walletPublicKeyHash is a 20-byte hash of the target wallet public key.
	walletPublicKeyHash [20]byte
	// refundPublicKeyHash is a 20-byte hash of the refund public key.
	refundPublicKeyHash [20]byte
	// refundLocktime is a 4-byte value representing the refund locktime.
	refundLocktime [4]byte
	// vault is an optional field that holds the host chain address of the
	// target vault.
	vault chain.Address
}
