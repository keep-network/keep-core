package tbtctest

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"math/big"
)

// UnmarshalJSON implements a custom JSON unmarshaling logic to produce a
// proper DepositSweepTestScenario.
func (dsts *DepositSweepTestScenario) UnmarshalJSON(data []byte) error {
	type depositSweepTestScenario struct {
		Title            string
		WalletPublicKey  string
		WalletPrivateKey string
		WalletMainUtxo   *utxo
		Deposits         []struct {
			Utxo                utxo
			Depositor           string
			BlindingFactor      string
			WalletPublicKeyHash string
			RefundPublicKeyHash string
			RefundLocktime      string
			Vault               string
		}
		InputTransactions                   []string
		Fee                                 int64
		Signatures                          []signature
		ExpectedSigHashes                   []string
		ExpectedSweepTransaction            string
		ExpectedSweepTransactionHash        string
		ExpectedSweepTransactionWitnessHash string
	}

	var unmarshaled depositSweepTestScenario

	err := json.Unmarshal(data, &unmarshaled)
	if err != nil {
		return err
	}

	// Unmarshal title.
	dsts.Title = unmarshaled.Title

	// Unmarshal wallet public key.
	x, y := elliptic.Unmarshal(
		tecdsa.Curve,
		hexToSlice(unmarshaled.WalletPublicKey),
	)
	dsts.WalletPublicKey = &ecdsa.PublicKey{
		Curve: tecdsa.Curve,
		X:     x,
		Y:     y,
	}

	// Unmarshal wallet private key.
	dsts.WalletPrivateKey = new(big.Int).SetBytes(
		hexToSlice(unmarshaled.WalletPrivateKey),
	)

	// Unmarshal optional wallet main UTXO.
	if walletMainUtxo := unmarshaled.WalletMainUtxo; walletMainUtxo != nil {
		dsts.WalletMainUtxo = walletMainUtxo.convert()
	}

	// Unmarshal deposits.
	for _, deposit := range unmarshaled.Deposits {
		d := new(Deposit)

		d.Utxo = deposit.Utxo.convert()
		d.Depositor = chain.Address(deposit.Depositor)
		copy(d.BlindingFactor[:], hexToSlice(deposit.BlindingFactor))
		copy(d.WalletPublicKeyHash[:], hexToSlice(deposit.WalletPublicKeyHash))
		copy(d.RefundPublicKeyHash[:], hexToSlice(deposit.RefundPublicKeyHash))
		copy(d.RefundLocktime[:], hexToSlice(deposit.RefundLocktime))

		var vault *chain.Address
		if v := chain.Address(deposit.Vault); len(v.String()) > 0 {
			vault = &v
		}
		d.Vault = vault

		dsts.Deposits = append(dsts.Deposits, d)
	}

	// Unmarshal input transactions.
	for _, inputTransaction := range unmarshaled.InputTransactions {
		transaction := new(bitcoin.Transaction)
		err = transaction.Deserialize(hexToSlice(inputTransaction))
		if err != nil {
			return err
		}

		dsts.InputTransactions = append(dsts.InputTransactions, transaction)
	}

	// Unmarshal fee.
	dsts.Fee = unmarshaled.Fee

	// Unmarshal signatures.
	for _, s := range unmarshaled.Signatures {
		dsts.Signatures = append(
			dsts.Signatures,
			s.convert(dsts.WalletPublicKey),
		)
	}

	// Unmarshal expected signature hashes.
	for _, expectedSigHash := range unmarshaled.ExpectedSigHashes {
		dsts.ExpectedSigHashes = append(
			dsts.ExpectedSigHashes,
			new(big.Int).SetBytes(hexToSlice(expectedSigHash)),
		)
	}

	// Unmarshal expected sweep transaction.
	dsts.ExpectedSweepTransaction = new(bitcoin.Transaction)
	err = dsts.ExpectedSweepTransaction.Deserialize(
		hexToSlice(unmarshaled.ExpectedSweepTransaction),
	)
	if err != nil {
		return err
	}

	// Unmarshal expected sweep transaction hash.
	dsts.ExpectedSweepTransactionHash, err = bitcoin.NewHashFromString(
		unmarshaled.ExpectedSweepTransactionHash,
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		return err
	}

	// Unmarshal expected sweep transaction witness hash.
	dsts.ExpectedSweepTransactionWitnessHash, err = bitcoin.NewHashFromString(
		unmarshaled.ExpectedSweepTransactionWitnessHash,
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		return err
	}

	return nil
}

// utxo is a helper type used for unmarshal UTXO encoded as JSON.
type utxo struct {
	Outpoint struct {
		// TransactionHash is the funding transaction hash in the
		// bitcoin.ReversedByteOrder.
		TransactionHash string
		// OutputIndex is index of the funding output.
		OutputIndex uint32
	}
	// Value is the value of the given UTXO.
	Value int64
}

// convert is responsible to construct a proper bitcoin.UnspentTransactionOutput
// object based on the unmarshaled data.
func (u utxo) convert() *bitcoin.UnspentTransactionOutput {
	transactionHash, err := bitcoin.NewHashFromString(
		u.Outpoint.TransactionHash,
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		panic(err)
	}

	return &bitcoin.UnspentTransactionOutput{
		Outpoint: &bitcoin.TransactionOutpoint{
			TransactionHash: transactionHash,
			OutputIndex:     u.Outpoint.OutputIndex,
		},
		Value: u.Value,
	}
}

// signature is a helper type used for unmarshal signatures encoded as JSON.
type signature struct {
	R, S string
}

// convert is responsible to construct a proper bitcoin.SignatureContainer
// object based on the unmarshaled data.
func (s signature) convert(
	publicKey *ecdsa.PublicKey,
) *bitcoin.SignatureContainer {
	return &bitcoin.SignatureContainer{
		R:         new(big.Int).SetBytes(hexToSlice(s.R)),
		S:         new(big.Int).SetBytes(hexToSlice(s.S)),
		PublicKey: publicKey,
	}
}

func hexToSlice(hexString string) []byte {
	if len(hexString) == 0 {
		return []byte{}
	}

	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		panic(err)
	}

	return bytes
}
