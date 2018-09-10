// Copyright 2017 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore" //
	"github.com/ethereum/go-ethereum/crypto"            //

	"github.com/urfave/cli"
)

/*
	"www.2c-why.com/Corp-Reg/MidGeth/Signature/crypto"   // "github.com/ethereum/go-ethereum/crypto"            //
	"www.2c-why.com/Corp-Reg/MidGeth/Signature/keystore" // "github.com/ethereum/go-ethereum/accounts/keystore" //
*/

var commandSignMessage = cli.Command{
	Name: "sign-message",
	Usage: `sign a message
						./key-test sign-message <KeyFile.json> "the message to sign"
							or
						./key-test sign-message --gen-msg <KeyFile.json> 
							or
						./key-test sign-message --msgfile <MessageInAFile> <KeyFile.json> 
							or
						./key-test sign-message --gen-msg <KeyFile.json> 
`,
	ArgsUsage: "<keyfile> <message>",
	Description: `
Sign the message with a keyfile.

To sign a message contained in a file, use the --msgfile flag.
`,
	Flags: []cli.Flag{
		passphraseFlag,
		jsonFlag,
		msgfileFlag,
		debugFlag,
		genMsgFlag,
	},
	Action: ActionSignMessage,
}

// TODO - add ability to generate a cryptographically strong message.

func ActionSignMessage(ctx *cli.Context) error {
	var message []byte
	var signature string
	var err error

	genMsg := ctx.Bool(genMsgFlag.Name)
	if !genMsg {
		message = getMessage(ctx, 1)
	}

	debugFlags := ctx.String(debugFlag.Name)
	SetDebugFlags(debugFlags)

	keyFile := ctx.Args().First()
	password := getPassphrase(ctx, false)
	messageStr, signature, err := GenerateSignature(keyFile, password, message)
	if err != nil {
		Fatalf(2, "Failed to sign message: %v", err)
	}

	if ctx.Bool("json") {
		fmt.Printf("{\"Signature\":%q,\"Message\":%q}", signature, messageStr)
	} else {
		fmt.Printf("MessageLengthInBytes: %d\n", len(messageStr))
		fmt.Printf("Message: %s\n", messageStr)
		fmt.Printf("Signature: %s\n", signature)
	}

	return nil
}

// GenerateSignature uses a keyfile and password to sign a message.  If the input message is "" then a random message
// will be generated.  The messgae and the signature are returned.
//
// 1. Find out where to call this - to gen signature GenerateSignature ( keyFile, keyFilePassword ) -> { message, sig }
func GenerateSignature(keyFile, password string, inMessage []byte) (message, signature string, err error) {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return "", "", fmt.Errorf("unable to read keyfile %s [%v]", keyFile, err)
	}
	key, err := keystore.DecryptKey(data, password)
	if err != nil {
		return "", "", fmt.Errorf("unable to decrypt %s [%v]", keyFile, err)
	}
	if len(inMessage) == 0 {
		inMessage, err = GenRandBytes(20)
		if err != nil {
			return "", "", fmt.Errorf("unable to generate random message [%v]", err)
		}
	}
	message = hex.EncodeToString(inMessage)
	rawSignature, err := crypto.Sign(signHash(inMessage), key.PrivateKey) // Sign Raw Bytes, Return hex of Raw Bytes
	if err != nil {
		return "", "", fmt.Errorf("unable to sign message [%v]", err)
	}
	signature = hex.EncodeToString(rawSignature)
	return message, signature, nil
}
