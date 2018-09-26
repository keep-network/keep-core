package signature

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// GetPublicKey takes an address + message + signature and produces a public key.
func GetPublicKey(addr, sig, msg string) (string, error) {
	tmp, err := VerifySignature(addr, sig, msg)
	if err != nil {
		return "", err
	}
	return tmp.RecoveredPublicKey, nil
}

func GetPublicKeyECDSA(sig, msg string) (*ecdsa.PublicKey, error) {
	message, err := hex.DecodeString(msg)
	if err != nil {
		return nil, fmt.Errorf("unabgle to decode message (invalid hex data) [%v]", err)
	}

	signature, err := hex.DecodeString(sig)
	if err != nil {
		return nil, fmt.Errorf("signature is not valid hex [%v]", err)
	}

	recoveredPubkey, err := crypto.SigToPub(signHash(message), signature)
	if err != nil || recoveredPubkey == nil {
		return nil, fmt.Errorf("signature verification failed [%v]", err)
	}
	return recoveredPubkey, nil
}

// MessageHasValidSignature takes an address + signature + message and returns
// true iff the signature is valid.
func MessageHasValidSignature(addr, sig, msg string) bool {
	tmp, err := VerifySignature(addr, sig, msg)
	if err != nil {
		return false
	}
	return tmp.IsValid
}

// VerifySignatureWithPubKey takes a ecdsa.PublicKey, a signature and a message and
// verifies that it is a valid signed message for the specified key.
func VerifySignatureWithPubKey(pubkey *ecdsa.PublicKey, sig, msg string) (bool, error) {
	message, err := hex.DecodeString(msg)
	if err != nil {
		return false, fmt.Errorf("unabgle to decode message (invalid hex data) [%v]", err)
	}

	signature, err := hex.DecodeString(sig)
	if err != nil {
		return false, fmt.Errorf("signature is not valid hex [%v]", err)
	}

	recoveredPubkey, err := crypto.SigToPub(signHash(message), signature)
	if err != nil || recoveredPubkey == nil {
		return false, fmt.Errorf("signature verification failed [%v]", err)
	}
	rawRecoveredAddress := crypto.PubkeyToAddress(*recoveredPubkey)
	isValid := PublicKeyToAddress(pubkey) == rawRecoveredAddress
	return isValid, nil
}

// PublicKeyToAddress converts form a public key to an Ethereum address.
func PublicKeyToAddress(p *ecdsa.PublicKey) common.Address {
	pubBytes := crypto.FromECDSAPub(p)
	return common.BytesToAddress(crypto.Keccak256(pubBytes[1:])[12:])
}
