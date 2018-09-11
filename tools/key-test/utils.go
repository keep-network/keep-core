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
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli"
)

// promptPassphrase prompts the user for a passphrase.  Set confirmation to true
// to require the user to confirm the passphrase.
func promptPassphrase(confirmation, isOldPassword bool) string {
	passphrase, err := console.Stdin.PromptPassword("Passphrase: ")
	if err != nil {
		Fatalf(2, "Failed to read passphrase: %v", err)
	}

	if confirmation {
		confirm, err := console.Stdin.PromptPassword("Repeat passphrase: ")
		if err != nil {
			Fatalf(2, "Failed to read passphrase confirmation: %v", err)
		}
		if passphrase != confirm {
			Fatalf(2, "Passphrases do not match")
		}
	}
	passphrase = strings.TrimRight(passphrase, " \r\n\f\t")
	passphrase = strings.TrimLeft(passphrase, " \r\n\f\t")

	if CheckPassword(passphrase, isOldPassword) == false {
		Fatalf(3, "Password is a known to be leeked (pwned) password - you will need to use a different password\n")
	}

	return passphrase
}

func ReadFileAsLines(filePath string) (lines []string, err error) {

	fp, err := os.Open(filePath)
	if err != nil {
		return lines, err
	}
	defer fp.Close()

	// ... for each line ...
	line_no := 0
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text() // no \n or \r\n on line - already chomped - os independent:w
		line_no++

		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return lines, fmt.Errorf("Error scanning file: %s", err.Error())
	}
	return
}

var getPassphraseInit = false
var getPassphraseLines []string
var getPassphraseNthLine = 0
var getPassphraseFile string

func getNextLine() (rv string) {
	if getPassphraseNthLine < len(getPassphraseLines) {
		rv = getPassphraseLines[getPassphraseNthLine]
		getPassphraseNthLine++
		return
	}
	Fatalf(2, "Ran out of lines in passphrase file '%s': %d", getPassphraseFile, getPassphraseNthLine)
	return ""
}

// getPassphrase obtains a passphrase given by the user.  It first checks the
// --passfile command line flag and ultimately prompts the user for a
// passphrase.
func getPassphrase(ctx *cli.Context, isOldPassword bool) string {
	var err error
	pwFn := ctx.String(passphraseFlag.Name)
	genRandom := ctx.Bool(randomPassFlag.Name)
	if genRandom {
		b, err := GenRandBytes(12)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		// newPhrase := hex.EncodeToString(b)             // could do better!!!!!!!!!!!!!!!!!!!!!!!!!!! PJS xyzzy FIXME
		str := base64.StdEncoding.EncodeToString(b) // could still do better
		// fmt.Printf("Password: %s\n", newPhrase)
		// See: https://tyler.io/generating-strong-user-friendly-passwords-in-php/
		// Also: https://sourceforge.net/projects/pwgen/ , pwgen-2.0.8.tar.gz
		str = strings.Replace(str, "0", "2", -1)
		str = strings.Replace(str, "1", "d", -1)
		str = strings.Replace(str, "Z", "v", -1)
		str = strings.Replace(str, "l", "3", -1)
		str = strings.Replace(str, "1", "h", -1)
		str = strings.Replace(str, "=", "u", -1)
		newPhrase := str
		return newPhrase
	} else if pwFn != "" {
		if !getPassphraseInit {
			getPassphraseInit = true
			getPassphraseFile = pwFn
			getPassphraseLines, err = ReadFileAsLines(pwFn)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
		}
		//data, err := ioutil.ReadFile(pwFn)
		//if err != nil {
		//	Fatalf(2, "Failed to read passphrase file '%s': %v", pwFn, err)
		//}
		//newPhrase := string(data)
		newPhrase := getNextLine()
		newPhrase = strings.TrimRight(newPhrase, " \r\n\f\t")
		newPhrase = strings.TrimLeft(newPhrase, " \r\n\f\t")
		if CheckPassword(newPhrase, isOldPassword) == false {
			Fatalf(3, "Password is a known to be leeked (pwned) password - you will need to use a different password\n")
		}
		return newPhrase
	}

	// Otherwise prompt the user for the passphrase.
	return promptPassphrase(false, isOldPassword)
}

// signHash is a helper function that calculates a hash for the given message
// that can be safely used to calculate a signature from.
//
// The hash is calulcated as
//   keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

// mustPrintJSON prints the JSON encoding of the given object and
// exits the program with an error message when the marshaling fails.
func mustPrintJSON(jsonObject interface{}) {
	jj, err := json.MarshalIndent(jsonObject, "", "  ")
	if err != nil {
		Fatalf(2, "Failed to marshal JSON object: %v", err)
	}
	fmt.Printf("%s", jj)
}

// db01 	dump keyfile
// db02 	early exit
// db03 	print public/private keys as items, X, Y, D -- TODO - xyzzy
var DbMap = make(map[string]bool)

// SetDebugFlags sets a map with debug flags.  Extra flags are ignored.  Valid ones are listed above.
func SetDebugFlags(debugFlags string) {
	if debugFlags == "" {
		return
	}
	flags := strings.Split(debugFlags, ",")
	fmt.Printf("Debug Flags Are: %s\n", FormatAsJSON(flags))
	for _, f := range flags {
		DbMap[f] = true
	}
}

func getMessage(ctx *cli.Context, msgarg int) []byte {
	if file := ctx.String("msgfile"); file != "" {
		if len(ctx.Args()) > msgarg {
			Fatalf(2, "Can't use --msgfile and message argument at the same time.")
		}
		msg, err := ioutil.ReadFile(file)
		if err != nil {
			Fatalf(2, "Can't read message file: %v", err)
		}
		return msg
	} else if len(ctx.Args()) == msgarg+1 {
		return []byte(ctx.Args().Get(msgarg))
	}
	Fatalf(2, "Invalid number of arguments: want %d, got %d", msgarg+1, len(ctx.Args()))
	return nil
}

// FormatAsJSON return the JSON encoded version of the data with tab indentation.
func FormatAsJSON(v interface{}) string {
	// s, err := json.Marshal ( v )
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("Error:%s", err)
	} else {
		return string(s)
	}
}
