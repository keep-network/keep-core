package signature

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// VerifySignature checks a signature and message and verifies the signature.
//
// Parameters:
//	addr		The address in hex that will be checked to validate the signature.
//	sig			The signature in hex that will be used.
//	msg			The message in hex that will be checked.
//
// Returns: If no error occures and if the signature validates.
//	recoveredAddress	The hex address in EIP-55 format
//	recoveredPublicKey	The public key for this address
//	isValid				True if the signature validated, false otherwize
//	err					Any errors from the process.
//
func VerifySignature(addr, sig, msg string) (string, string, bool, error) {
	message, err := hex.DecodeString(msg)
	if err != nil {
		return "", "", false, fmt.Errorf("unabgle to decode message (invalid hex data) [%v]", err)
	}
	if !common.IsHexAddress(addr) {
		return "", "", false, fmt.Errorf("invalid address: %s", addr)
	}
	address := common.HexToAddress(addr)
	signature, err := hex.DecodeString(sig)
	if err != nil {
		return "", "", false, fmt.Errorf("signature is not valid hex [%v]", err)
	}

	recoveredPubkey, err := crypto.SigToPub(signHash(message), signature)
	if err != nil || recoveredPubkey == nil {
		return "", "", false, fmt.Errorf("signature verification failed [%v]", err)
	}
	recoveredPublicKey := hex.EncodeToString(crypto.FromECDSAPub(recoveredPubkey))
	rawRecoveredAddress := crypto.PubkeyToAddress(*recoveredPubkey)
	isValid := address == rawRecoveredAddress
	if !isValid {
		return "", "", isValid, nil
	}
	recoveredAddress := rawRecoveredAddress.Hex()
	return recoveredAddress, recoveredPublicKey, isValid, nil
}
