package tbtc

import (
	"encoding/hex"
	"fmt"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
)

// depositScriptFormat is the format of the deposit P2(W)SH Bitcoin script
// The specific placeholders are: depositor, blindingFactor, walletPublicKeyHash,
// refundPublicKeyHash and refundLocktime. For reference see:
// https://github.com/keep-network/tbtc-v2/blob/83310bdc9ed934e286bc9ea5091cc16979950134/solidity/contracts/bridge/Deposit.sol#L172
const depositScriptFormat = "14%v7508%v7576a914%v8763ac6776a914%v8804%vb175ac68"

// deposit represents a tBTC deposit.
type deposit struct {
	// utxo is the unspent output of the deposit funding transaction that
	// represents the deposit on the Bitcoin chain.
	utxo *bitcoin.UnspentTransactionOutput
	// depositor is the depositor's address on the host chain.
	depositor [20]byte
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

// script constructs the deposit P2(W)SH Bitcoin script. This function
// assumes the deposit's fields are correctly set.
func (d *deposit) script() ([]byte, error) {
	script := fmt.Sprintf(
		depositScriptFormat,
		d.depositor,
		d.blindingFactor,
		d.walletPublicKeyHash,
		d.refundPublicKeyHash,
		d.refundLocktime,
	)

	return hex.DecodeString(script)
}
