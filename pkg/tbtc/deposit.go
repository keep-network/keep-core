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

// Deposit represents a tBTC deposit.
type Deposit struct {
	// Utxo is the unspent output of the deposit funding transaction that
	// represents the deposit on the Bitcoin chain.
	Utxo *bitcoin.UnspentTransactionOutput
	// Depositor is the depositor's address on the host chain.
	Depositor [20]byte
	// BlindingFactor is an 8-byte arbitrary value that allows to distinguish
	// deposits from the same depositor.
	BlindingFactor [8]byte
	// WalletPublicKeyHash is a 20-byte hash of the target wallet public key.
	WalletPublicKeyHash [20]byte
	// RefundPublicKeyHash is a 20-byte hash of the refund public key.
	RefundPublicKeyHash [20]byte
	// RefundLocktime is a 4-byte value representing the refund locktime.
	RefundLocktime [4]byte
	// Vault is an optional field that holds the host chain address of the
	// target vault.
	Vault chain.Address
}

// Script constructs the deposit P2(W)SH Bitcoin script. This function
// assumes the deposit's fields are correctly set.
func (d *Deposit) Script() ([]byte, error) {
	script := fmt.Sprintf(
		depositScriptFormat,
		hex.EncodeToString(d.Depositor[:]),
		hex.EncodeToString(d.BlindingFactor[:]),
		hex.EncodeToString(d.WalletPublicKeyHash[:]),
		hex.EncodeToString(d.RefundPublicKeyHash[:]),
		hex.EncodeToString(d.RefundLocktime[:]),
	)

	return hex.DecodeString(script)
}
