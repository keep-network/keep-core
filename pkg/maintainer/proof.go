package maintainer

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

// AssembleTransactionProof assembles a proof that a given transaction was
// included in the blockchain and has accumulated the required number of
// confirmations.
func AssembleTransactionProof(
	transactionHash bitcoin.Hash,
	requiredConfirmations uint,
	bitcoinClient bitcoin.Chain,
) (*bitcoin.Transaction, *bitcoin.Proof, error) {
	transaction, err := bitcoinClient.GetTransaction(transactionHash)
	if err != nil {
		return nil, nil, err
	}

	confirmations, err := bitcoinClient.GetTransactionConfirmations(transactionHash)
	if err != nil {
		return nil, nil, err
	}

	if confirmations < requiredConfirmations {
		return nil, nil, fmt.Errorf(
			"transaction confirmations number[%v] is not enough, required [%v]",
			confirmations,
			requiredConfirmations,
		)
	}

	latestBlockHeight, err := bitcoinClient.GetLatestBlockHeight()
	if err != nil {
		return nil, nil, err
	}

	txBlockHeight := latestBlockHeight - confirmations + 1

	headersChain, err := getHeadersChain(
		bitcoinClient,
		txBlockHeight,
		requiredConfirmations-1,
	)
	if err != nil {
		return nil, nil, err
	}

	merkleBranch, err := bitcoinClient.GetTransactionMerkleProof(
		transactionHash,
		txBlockHeight,
	)
	if err != nil {
		return nil, nil, err
	}

	merkleProof, err := CreateMerkleProof(merkleBranch)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Merkle proof [%w]", err)
	}

	proof := &bitcoin.Proof{
		MerkleProof:    merkleProof,
		TxIndexInBlock: merkleBranch.Position,
		BitcoinHeaders: headersChain,
	}

	return transaction, proof, nil
}

// CreateMerkleProof creates a proof of transaction inclusion in the block by
// concatenating 32-byte-long hash values. The values are converted to the
// little endian form. The branch of a Merkle tree leading to a transaction
// needs to be provided. The transaction inclusion proof in hexadecimal form is
// returned.
func CreateMerkleProof(txMerkleBranch *bitcoin.TransactionMerkleProof) (
	string,
	error,
) {
	var proof bytes.Buffer

	for _, item := range txMerkleBranch.MerkleNodes {
		hashBytes, err := hex.DecodeString(item)
		if err != nil {
			return "", err
		}
		reversedHash := reverseBytes(hashBytes)
		proof.Write(reversedHash)
	}
	return hex.EncodeToString(proof.Bytes()), nil
}

// reverseBytes reverses the order of bytes in a byte slice.
func reverseBytes(b []byte) []byte {
	length := len(b)
	reversed := make([]byte, length)
	for i := 0; i < length; i++ {
		reversed[i] = b[length-1-i]
	}
	return reversed
}

// getHeadersChain gets a chain of Bitcoin block headers that starts at the
// provided block height and has the specified number of subsequent headers.
func getHeadersChain(
	bitcoinClient bitcoin.Chain,
	blockHeight uint,
	chainLength uint,
) ([]*bitcoin.BlockHeader, error) {
	// TODO: Consider modifying the Bitcoin chain so that it can return
	//       multiple headers
	var blockHeaders []*bitcoin.BlockHeader

	for i := blockHeight; i <= blockHeight+chainLength; i++ {
		blockHeader, err := bitcoinClient.GetBlockHeader(i)
		if err != nil {
			return nil, err
		}
		blockHeaders = append(blockHeaders, blockHeader)
	}
	return blockHeaders, nil
}
