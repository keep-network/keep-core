package signature

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

// Sign uses a key to sign a message. The message, encoded in hex.
// The signature and message are returned.
func Sign(
	key *keystore.Key,
	message []byte,
) (string, string, error) {
	messageStr := hex.EncodeToString(message)
	signature, err := crypto.Sign(signHash(message), key.PrivateKey)
	if err != nil {
		return "", "", fmt.Errorf("unable to sign message [%v]", err)
	}
	signatureStr := hex.EncodeToString(signature)
	return messageStr, signatureStr, nil
}

// signHash is a helper function that calculates a hash for the given message.
func signHash(data []byte) []byte {
	return crypto.Keccak256(data)
}
