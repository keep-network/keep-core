package main

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli"
)

var commandVerifyMessage = cli.Command{
	Name: "verify-message",
	Usage: `verify the signature of a signed message
                ./key-test very-message Address Signature "message"
                    or
                ./key-test very-message --json Address Signature "message"
`,
	ArgsUsage: "<address> <signature> <message>",
	Description: `
Verify the signature of the message.
It is possible to refer to a file containing the message.`,
	Flags: []cli.Flag{
		jsonFlag,
		msgfileFlag,
		debugFlag,
	},
	Action: ActionVerifyMessage,
}

// ActinVerifyMessage takes an address, signature and the original messsage on the command line and verifies the
// signature for that message.
func ActionVerifyMessage(ctx *cli.Context) error {
	addressStr := ctx.Args().First()
	signatureHex := ctx.Args().Get(1)
	message := getMessage(ctx, 2) // xyzzy - change to return a string!

	debugFlags := ctx.String(debugFlag.Name)
	SetDebugFlags(debugFlags)

	// message is hex, need to convert back to byte, line 61 dose this.

	rAddr, rPubKey, err := VerifySignature(addressStr, signatureHex, string(message))
	if err != nil {
		// fmt.Printf("Signature verification failed!\n")
		return fmt.Errorf("signature did not verify")
	}
	fmt.Printf("Signature verified\n")
	fmt.Printf("Recovered public key: %s\n", rPubKey)
	// See: https://github.com/ethereum/EIPs/blob/master/EIPS/eip-55.md
	fmt.Printf("Recovered address (in EIP-55 format): %s\n", rAddr)

	return nil
}

// 2. Find out where to verify this - VerifySignature ( Address, Signature, Message ) -> bool

// VerifySignature takes hex encoded addr, sig and msg and verifies that the signature matches with the address.
//
func VerifySignature(addr, sig, msg string) (recoveredAddress, recoveredPublicKey string, err error) {
	message, err := hex.DecodeString(msg)
	if err != nil {
		return "", "", fmt.Errorf("unabgle to decode message (invalid hex data) [%v]", err)
	}
	if !common.IsHexAddress(addr) {
		return "", "", fmt.Errorf("invalid address: %s", addr)
	}
	address := common.HexToAddress(addr)
	signature, err := hex.DecodeString(sig)
	if err != nil {
		return "", "", fmt.Errorf("signature is not valid hex [%v]", err)
	}

	recoveredPubkey, err := crypto.SigToPub(signHash(message), signature)
	if err != nil || recoveredPubkey == nil {
		return "", "", fmt.Errorf("signature verification failed [%v]", err)
	}
	recoveredPublicKey = hex.EncodeToString(crypto.FromECDSAPub(recoveredPubkey))
	rawRecoveredAddress := crypto.PubkeyToAddress(*recoveredPubkey)
	if address != rawRecoveredAddress {
		return "", "", fmt.Errorf("signature did not verify, addresses did not match")
	}
	recoveredAddress = rawRecoveredAddress.Hex()
	return
}
