package ethutil_test

import (
	"fmt"
	"testing"

	"github.com/keep-network/keep-core/pkg/chain/ethereum/ethutil"
)

func TestKeyFileDecryption(t *testing.T) {
	goodKeyFile := "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	badKeyFile := "./testdata/nonexistent-file.booyan"

	tests := map[string]struct {
		keyFile      string
		password     string
		errorMessage string
	}{
		"good password": {
			keyFile:      goodKeyFile,
			password:     "password",
			errorMessage: "",
		},
		"bad file": {
			keyFile:  badKeyFile,
			password: "",
			errorMessage: fmt.Sprintf(
				"unable to read KeyFile %v [open %v: no such file or directory]",
				badKeyFile,
				badKeyFile,
			),
		},
		"bad password": {
			keyFile:  goodKeyFile,
			password: "nanananana",
			errorMessage: fmt.Sprintf(
				"unable to decrypt %v [could not decrypt key with given passphrase]",
				goodKeyFile,
			),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			_, err := ethutil.DecryptKeyFile(test.keyFile, test.password)
			message := ""
			if err != nil {
				message = err.Error()
			}

			if message != test.errorMessage {
				t.Errorf(
					"\nexpected: [%v]\nactual:   [%v]",
					test.errorMessage,
					err,
				)
			}
		})
	}
}
