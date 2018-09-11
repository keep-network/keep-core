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
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore" //
	"github.com/ethereum/go-ethereum/crypto"            //

	"github.com/pborman/uuid"
	"github.com/urfave/cli"
)

/*
	"www.2c-why.com/Corp-Reg/MidGeth/Signature/crypto"   // "github.com/ethereum/go-ethereum/crypto"            //
	"www.2c-why.com/Corp-Reg/MidGeth/Signature/keystore" // "github.com/ethereum/go-ethereum/accounts/keystore" //
*/

type outputGenerate struct {
	Address      string
	AddressEIP55 string
}

var commandGenerate = cli.Command{
	Name: "generate",
	Usage: `generate new keyfile
					./key-test generate --default-name --random-pass
`,
	ArgsUsage: "[ <keyfile> ]",
	Description: `
Generate a new keyfile.

If you want to encrypt an existing private key, it can be specified by setting
--privatekey with the location of the file containing the private key.
`,
	Flags: []cli.Flag{
		passphraseFlag,
		jsonFlag,
		cli.StringFlag{
			Name:  "privatekey",
			Usage: "file containing a raw private key to encrypt",
		},
		defaultNameFlag,
		randomPassFlag,
		logFileFlag,
	},
	Action: ActionGenerate,
}

// ActionGenerate will generate a new key file.
func ActionGenerate(ctx *cli.Context) error {
	// Check if keyfile path given and make sure it doesn't already exist.
	keyfilepath := ctx.Args().First()

	// fmt.Printf("PJS: batch=%d\n", batch) // batch is 0 if not set.
	// if batch > 0 and --passphrase <file> then read 1 line for each password
	// if batch > 0 and --random-pass then generate random passwrods and print them out.

	var privateKey *ecdsa.PrivateKey
	var err error

	debugFlags := ctx.String(debugFlag.Name)
	SetDebugFlags(debugFlags)

	privateKeyFileName := ctx.String("privatekey")

	// ----------------------------------------------------------------------------------------------------------

	if file := privateKeyFileName; file != "" {
		// Load private key from file.
		privateKey, err = crypto.LoadECDSA(file)
		if err != nil {
			Fatalf(2, "Can't load private key: %v", err)
		}
	} else {
		// If not loaded, generate random.
		privateKey, err = crypto.GenerateKey()
		if err != nil {
			Fatalf(2, "Failed to generate random private key: %v", err)
		}
	}

	// Create the keyfile object with a random UUID.
	id := uuid.NewRandom()
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}

	// Encrypt key with passphrase.
	// passphrase := promptPassphrase(true, false)
	passphrase := getPassphrase(ctx, false)
	keyjson, err := keystore.EncryptKey(key, passphrase, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		Fatalf(2, "Error encrypting key: %v", err)
	}

	if DbMap["db03"] {
		fmt.Printf("PJS: Key file will be stored in: %s\n", filepath.Dir(keyfilepath))
	}

	address := key.Address.Hex()

	newKeyFileDir := filepath.Dir(keyfilepath)
	newname := filepath.Base(keyfilepath)
	defaultName := ctx.Bool(defaultNameFlag.Name)
	if defaultName && newname != "." {
		Fatalf(2, "Can not supply both a --default-name flag and a name for the file. defaultName=true newname=[%s]\n", newname)
	} else if defaultName {
		newname = fmt.Sprintf("UTC--%s--%s", time.Now().UTC().Format("2006-01-02T15-04-05.9999999999Z"), address[2:])
	}

	// Store the file to disk.
	if err := os.MkdirAll(newKeyFileDir, 0700); err != nil {
		Fatalf(2, "Could not create directory %s", newKeyFileDir)
	}

	path := filepath.Join(newKeyFileDir, newname)

	// check if file already exists - if so then cowerdly refuse to ovewrite it.
	if Exists(path) {
		Fatalf(2, "File [%s] already exists - will not overwrite\n", path)
	}

	// Output the file.
	if err := ioutil.WriteFile(path, keyjson, 0600); err != nil {
		Fatalf(2, "Failed to write keyfile to %s: %v", path, err)
	}

	// Output some information.
	out := outputGenerate{
		Address: address,
	}
	if ctx.Bool(jsonFlag.Name) {
		mustPrintJSON(out)
	} else {
		genRandom := ctx.Bool(randomPassFlag.Name)
		if genRandom {
			fmt.Printf("Password: %s\n", passphrase)
		}
		fmt.Println("Address:", out.Address)
		fmt.Println("File Name:", path)
		lf := ctx.String(logFileFlag.Name)
		if lf != "" {
			AppendToLog(lf, fmt.Sprintf("Password: %s\nAddress: %s\nFile Name:%s\n\n", passphrase, out.Address, path))
		}
	}
	return nil
}
