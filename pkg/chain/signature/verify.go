package signature

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

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
