package main

// List of all of the flags that are commonly used in this program

import (
	"github.com/urfave/cli"
)

var passphraseFlag = cli.StringFlag{
	Name:  "passwordfile",
	Usage: "the file that contains the passphrase for the keyfile",
}

var jsonFlag = cli.BoolFlag{
	Name:  "json",
	Usage: "output JSON instead of human-text format",
}

var newPassphraseFlag = cli.StringFlag{
	Name:  "newpasswordfile",
	Usage: "the file that contains the new passphrase for the keyfile",
}

var debugFlag = cli.StringFlag{
	Name:  "debug",
	Usage: "turn on debug flags",
}

var msgfileFlag = cli.StringFlag{
	Name:  "msgfile",
	Usage: "file containing the message to sign/verify",
}

var defaultNameFlag = cli.BoolFlag{
	Name:  "default-name",
	Usage: "Use a default name for the output file, `DateTime--AccountNo`.",
}

var genMsgFlag = cli.BoolFlag{
	Name:  "gen-msg",
	Usage: "randomly generate a message to sign",
}

var randomPassFlag = cli.BoolFlag{
	Name:  "random-pass",
	Usage: "randomly generate passwords",
}

var logFileFlag = cli.StringFlag{
	Name:  "log-file",
	Usage: "file where info will be logged",
}

var newNameFlag = cli.StringFlag{
	Name:  "newname",
	Usage: "new name of keyfile, set to '-' to overwrite the existing file.",
}
