package signature

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

// GenerateSignature uses a key file and password to sign a message.
// If the input message is a zero length byte slice then a random message
// 20 bytes long will be generated.  The message, encoded in hex, and the
// signature are returned.
func GenerateSignature(
	keyFile, password string,
	inMessage []byte,
) ([]byte, string, error) {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return []byte{}, "", fmt.Errorf("unable to read keyfile %s [%v]", keyFile, err)
	}
	key, err := keystore.DecryptKey(data, password)
	if err != nil {
		return []byte{}, "", fmt.Errorf("unable to decrypt %s [%v]", keyFile, err)
	}
	var message []byte
	if len(inMessage) == 0 {
		message, err = genRandBytes(20)
		if err != nil {
			return []byte{}, "", fmt.Errorf("unable to generate random message [%v]", err)
		}
		tmp := hex.EncodeToString(message)
		message = []byte(tmp)
	} else {
		message = inMessage
	}
	rawSignature, err := crypto.Sign(signHash(message), key.PrivateKey)
	if err != nil {
		return []byte{}, "", fmt.Errorf("unable to sign message [%v]", err)
	}
	signature := hex.EncodeToString(rawSignature)
	return message, signature, nil
}

// genRandBytes generates `n` random bytes of data using the cryptographically
// strong random generator.
func genRandBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// signHash is a helper function that calculates a hash for the given message.
func signHash(data []byte) []byte {
	return crypto.Keccak256(data)
}
