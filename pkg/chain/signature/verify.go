package signature

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// VerifiedSignatureData is the set of data returend when
// a signature is verified.  If the signature fails to verify
// IsValid is false.
type VerifiedSignatureData struct {
	IsValid            bool   // True if the signature is a valid signature.
	RecoveredAddress   string // The hex address in EIP-55 format
	RecoveredPublicKey string // The public key for this address
}

// VerifySignature checks a signature and message and verifies the signature.
//
// Parameters:
//	addr		The address in hex that will be checked to validate the signature.
//	sig			The signature in hex that will be used.
//	msg			The message in hex that will be checked.
//
// Returns a VerifiedSignatureData if no error occures and if the signature validates.
func VerifySignature(addr, sig, msg string) (*VerifiedSignatureData, error) {
	message, err := hex.DecodeString(msg)
	if err != nil {
		return nil, fmt.Errorf("unabgle to decode message (invalid hex data) [%v]", err)
	}
	if !common.IsHexAddress(addr) {
		return nil, fmt.Errorf("invalid address: %s", addr)
	}
	address := common.HexToAddress(addr)
	signature, err := hex.DecodeString(sig)
	if err != nil {
		return nil, fmt.Errorf("signature is not valid hex [%v]", err)
	}

	recoveredPubkey, err := crypto.SigToPub(signHash(message), signature)
	if err != nil || recoveredPubkey == nil {
		return nil, fmt.Errorf("signature verification failed [%v]", err)
	}
	recoveredPublicKey := hex.EncodeToString(crypto.FromECDSAPub(recoveredPubkey))
	rawRecoveredAddress := crypto.PubkeyToAddress(*recoveredPubkey)
	isValid := address == rawRecoveredAddress
	if !isValid {
		return &VerifiedSignatureData{}, nil
	}
	recoveredAddress := rawRecoveredAddress.Hex()
	return &VerifiedSignatureData{
		IsValid:            isValid,
		RecoveredAddress:   recoveredAddress,
		RecoveredPublicKey: recoveredPublicKey,
	}, nil
}
