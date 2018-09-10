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

	"github.com/urfave/cli"

	"github.com/ethereum/go-ethereum/accounts/keystore" //
	"github.com/ethereum/go-ethereum/crypto"            //
)

/*
	"www.2c-why.com/Corp-Reg/MidGeth/Signature/crypto"   // "github.com/ethereum/go-ethereum/crypto"            //
	"www.2c-why.com/Corp-Reg/MidGeth/Signature/keystore" // "github.com/ethereum/go-ethereum/accounts/keystore" //
*/

type outputInspect struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

var commandInspect = cli.Command{
	Name: "inspect",
	Usage: `inspect a keyfile
                 key-test inspect --priave <KeyFile.json> 
                 Print out the Address, Public Key, Private Key for a keyfile.
`,
	ArgsUsage: "<KeyFile.json>",
	Description: `
Print the address and the publick key from the keyfile.   Optionally print the private key.

Private key information can be printed by using the --private flag;
Use this feature with great caution!  It is your **private** unencrypted key.

Status 
	0 (sucdess) if a valid password is provided
	1 if an invalid password is provide
	2 if some other error occures (missing file etc.)
`,
	Flags: []cli.Flag{
		passphraseFlag,
		jsonFlag,
		cli.BoolFlag{
			Name:  "private",
			Usage: "include the private key in the output",
		},
		debugFlag, // from common-flags.go
	},
	Action: ActionInspect,
}

func ActionInspect(ctx *cli.Context) error {
	keyfilepath := ctx.Args().First()

	debugFlags := ctx.String(debugFlag.Name)
	SetDebugFlags(debugFlags)

	// Read key from file.
	keyjson, err := ioutil.ReadFile(keyfilepath)
	if err != nil {
		Fatalf(2, "Failed to read the keyfile at '%s': %v", keyfilepath, err)
	}

	// Decrypt key with passphrase.
	passphrase := getPassphrase(ctx, true)
	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		Fatalf(1, "Error decrypting key: %v", err)
	}

	// Output all relevant information we can retrieve.
	showPrivate := ctx.Bool("private")
	out := outputInspect{
		Address:   key.Address.Hex(),
		PublicKey: hex.EncodeToString(crypto.FromECDSAPub(&key.PrivateKey.PublicKey)),
	}
	if showPrivate {
		out.PrivateKey = hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))
	}

	if ctx.Bool(jsonFlag.Name) {
		mustPrintJSON(out)
	} else {
		fmt.Println("Address:       ", out.Address)
		fmt.Println("Public key:    ", out.PublicKey)
		if showPrivate {
			fmt.Println("Private key:   ", out.PrivateKey)
		}
	}
	return nil
}
