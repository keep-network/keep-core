package signature

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// RecoverPublicKey use a signature and a mesage to recover the publick key.
func RecoverPublicKey(sig, msg string) (*ecdsa.PublicKey, error) {
	message, err := hex.DecodeString(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex message to bytes: [%v]", err)
	}

	signature, err := hex.DecodeString(sig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex signature to bytes: [%v]", err)
	}

	recoveredPubkey, err := crypto.SigToPub(signHash(message), signature)
	if err != nil || recoveredPubkey == nil {
		return nil, fmt.Errorf("failed to verify signature: [%v]", err)
	}
	return recoveredPubkey, nil
}

// PublicKeyToAddress converts form a public key to an Ethereum address.
func PublicKeyToAddress(p *ecdsa.PublicKey) common.Address {
	if p == nil {
		return common.Address{0}
	}
	pubBytes := crypto.FromECDSAPub(p)
	return common.BytesToAddress(crypto.Keccak256(pubBytes[1:])[12:])
}
