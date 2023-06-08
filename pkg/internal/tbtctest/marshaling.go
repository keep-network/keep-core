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
	"time"
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

// UnmarshalJSON implements a custom JSON unmarshaling logic to produce a
// proper RedemptionTestScenario.
func (rts *RedemptionTestScenario) UnmarshalJSON(data []byte) error {
	type redemptionTestScenario struct {
		Title              string
		WalletPublicKey    string
		WalletPrivateKey   string
		WalletMainUtxo     *utxo
		RedemptionRequests []struct {
			Redeemer             string
			RedeemerOutputScript string
			RequestedAmount      uint64
			TreasuryFee          uint64
			TxMaxFee             uint64
			RequestedAt          int64
		}
		InputTransaction                         string
		Fee                                      int64
		Signature                                signature
		ExpectedSigHash                          string
		ExpectedRedemptionTransaction            string
		ExpectedRedemptionTransactionHash        string
		ExpectedRedemptionTransactionWitnessHash string
	}

	var unmarshaled redemptionTestScenario

	err := json.Unmarshal(data, &unmarshaled)
	if err != nil {
		return err
	}

	// Unmarshal title.
	rts.Title = unmarshaled.Title

	// Unmarshal wallet public key.
	x, y := elliptic.Unmarshal(
		tecdsa.Curve,
		hexToSlice(unmarshaled.WalletPublicKey),
	)
	rts.WalletPublicKey = &ecdsa.PublicKey{
		Curve: tecdsa.Curve,
		X:     x,
		Y:     y,
	}

	// Unmarshal wallet private key.
	rts.WalletPrivateKey = new(big.Int).SetBytes(
		hexToSlice(unmarshaled.WalletPrivateKey),
	)

	// Unmarshal wallet main UTXO.
	rts.WalletMainUtxo = unmarshaled.WalletMainUtxo.convert()

	// Unmarshal redemption requests.
	for _, request := range unmarshaled.RedemptionRequests {
		r := new(RedemptionRequest)

		r.Redeemer = chain.Address(request.Redeemer)
		r.RedeemerOutputScript = hexToSlice(request.RedeemerOutputScript)
		r.RequestedAmount = request.RequestedAmount
		r.TreasuryFee = request.TreasuryFee
		r.TxMaxFee = request.TxMaxFee
		r.RequestedAt = time.Unix(request.RequestedAt, 0)

		rts.RedemptionRequests = append(rts.RedemptionRequests, r)
	}

	// Unmarshal input transaction.
	rts.InputTransaction = new(bitcoin.Transaction)
	err = rts.InputTransaction.Deserialize(hexToSlice(unmarshaled.InputTransaction))
	if err != nil {
		return err
	}

	// Unmarshal fee.
	rts.Fee = unmarshaled.Fee

	// Unmarshal signature.
	rts.Signature = unmarshaled.Signature.convert(rts.WalletPublicKey)

	// Unmarshal expected signature hash.
	rts.ExpectedSigHash = new(big.Int).SetBytes(hexToSlice(unmarshaled.ExpectedSigHash))

	// Unmarshal expected redemption transaction.
	rts.ExpectedRedemptionTransaction = new(bitcoin.Transaction)
	err = rts.ExpectedRedemptionTransaction.Deserialize(
		hexToSlice(unmarshaled.ExpectedRedemptionTransaction),
	)
	if err != nil {
		return err
	}

	// Unmarshal expected redemption transaction hash.
	rts.ExpectedRedemptionTransactionHash, err = bitcoin.NewHashFromString(
		unmarshaled.ExpectedRedemptionTransactionHash,
		bitcoin.ReversedByteOrder,
	)
	if err != nil {
		return err
	}

	// Unmarshal expected redemption transaction witness hash.
	rts.ExpectedRedemptionTransactionWitnessHash, err = bitcoin.NewHashFromString(
		unmarshaled.ExpectedRedemptionTransactionWitnessHash,
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
