package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/urfave/cli"
)

var commandChangePassphrase = cli.Command{
	Name: "change-passphrase",
	Usage: `change the passphrase on a keyfile
                 key-test changepassphrase --newname <newKeyFileName.json> <ExistingKeyFile.json> 
                 if <newKeyFileName.json> is "-" then the existing file will be overwritten with the new one.
                 You will be prompted for the passwords - unles you supply the --password and --newPassword flags.
				 key-test changepassphrase --newname <NewFileName.json> --password <FileWithOldPassword> --newPassword <fileWithNewPw> <ExistingKeyFile.json>
`,
	ArgsUsage: "<keyfile>",
	Description: `
Change the passphrase/password of a keyfile.  
`,
	Flags: []cli.Flag{
		passphraseFlag,    // --passwordfile <filename> - from common-flags.go
		newPassphraseFlag, // --newpasswordfile <filename> - from common-flags.go
		newNameFlag,
		defaultNameFlag,
		debugFlag, // --debug Flag1,Flag2,... - from common-flags.go
	},
	Action: ActionChangeKeyfilePassword,
}

var commandChangePassword = cli.Command{
	Name: "change-password",
	Usage: `change the passphrase on a keyfile
                 key-test change-password --newname <newKeyFileName.json> <ExistingKeyFile.json> 
                 if <newKeyFileName.json> is "-" then the existing file will be overwritten with the new one.
                 You will be prompted for the passwords - unles you supply the --password and --newPassword flags.
				 key-test change-password --newname <NewFileName.json> --password <FileWithOldPassword> --newPassword <fileWithNewPw> <ExistingKeyFile.json>
`,
	ArgsUsage: "<keyfile>",
	Description: `
Change the passphrase/password of a keyfile.  
`,
	Flags: []cli.Flag{
		passphraseFlag,    // --passwordfile <filename> - from common-flags.go
		newPassphraseFlag, // --newpasswordfile <filename> - from common-flags.go
		debugFlag,         // --debug Flag1,Flag2,... - from common-flags.go
	},
	Action: ActionChangeKeyfilePassword,
}

// ActionChangeKeyfilePassword processes the changepassphrase and changepassword commands on the command line.
func ActionChangeKeyfilePassword(ctx *cli.Context) error {
	keyfilepath := ctx.Args().First()

	debugFlags := ctx.String(debugFlag.Name)
	SetDebugFlags(debugFlags)

	if !Exists(keyfilepath) {
		Fatalf(2, "input key file missing: %s\n", keyfilepath)
	}

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

	// dump the keyfile to stdout - non-encrypted.
	if DbMap["db01"] {
		fmt.Printf("Dump of keyfile - unencrypted: %s\n", FormatAsJSON(key))
	}

	address := key.Address.Hex()

	newKeyFileName := keyfilepath
	newname := ctx.String(newNameFlag.Name)
	defaultName := ctx.Bool(defaultNameFlag.Name)
	if defaultName && newname != "" {
		Fatalf(2, "Can only supply one of --newname FN and --defaultName\n")
	} else if defaultName {
		newname = fmt.Sprintf("UTC--%s--%s\n", time.Now().UTC().Format("2006-01-02T15:04:05.999999999Z"), address)
	} else if newname == "" {
		Fatalf(2, "Must supply a new name for the key file or --defaultName option must be set\n")
	} else if newname == "-" { // Then overwrite, default set above.
	} else {
		newKeyFileName = newname
		if Exists(newKeyFileName) {
			Fatalf(2, "Will not overwrite an existing keyfile - remove existing keyfile first: %s\n", newKeyFileName)
		}
	}

	// Get the new password
	var newPhrase string
	if passFile := ctx.String(newPassphraseFlag.Name); passFile != "" {
		content, err := ioutil.ReadFile(passFile)
		if err != nil {
			Fatalf(3, "Failed to read new passphrase file '%s': %v", passFile, err)
		}
		newPhrase = string(content)
		// FIXME - TODO - should check that the password passes muster with remote site at this point
	} else {
		fmt.Printf("Input a new passphrase (leading and trailing blanks are ignored)\n")
		newPhrase = promptPassphrase(true, true)
	}
	newPhrase = strings.TrimRight(newPhrase, " \r\n\f\t")
	newPhrase = strings.TrimLeft(newPhrase, " \r\n\f\t")

	// Encrypt the key with the new passphrase.
	newJson, err := keystore.EncryptKey(key, newPhrase, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		Fatalf(2, "Error encrypting with new passphrase: %v", err)
	}

	// Then write the new keyfile in place of the old one.
	if err := ioutil.WriteFile(newKeyFileName, newJson, 600); err != nil {
		Fatalf(2, "Error writing new keyfile to disk: %v", err)
	}

	// Don't print anything.  Just return successfully, producing a positive exit code.
	return nil
}
