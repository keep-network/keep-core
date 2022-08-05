package tecdsatest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
)

const (
	testDataDirFormat                 = "%s/testdata"
	privateKeyShareTestDataFileFormat = "private_key_share_data_%d.json"
)

// LoadPrivateKeyShareTestFixtures loads tECDSA private key share test data.
// Code copied from:
//   https://github.com/binance-chain/tss-lib/blob/master/ecdsa/keygen/test_utils.go
// Test data JSON files copied from:
//   https://github.com/binance-chain/tss-lib/tree/master/test/_ecdsa_fixtures
func LoadPrivateKeyShareTestFixtures(count int) (
	[]keygen.LocalPartySaveData,
	error,
) {
	makeTestFixtureFilePath := func(partyIndex int) string {
		_, callerFileName, _, _ := runtime.Caller(0)
		srcDirName := filepath.Dir(callerFileName)
		fixtureDirName := fmt.Sprintf(testDataDirFormat, srcDirName)
		return fmt.Sprintf(
			"%s/"+privateKeyShareTestDataFileFormat,
			fixtureDirName,
			partyIndex,
		)
	}

	shares := make([]keygen.LocalPartySaveData, 0, count)

	for j := 0; j < count; j++ {
		fixtureFilePath := makeTestFixtureFilePath(j)

		// #nosec G304 (file path provided as taint input)
		// This line is used to read a test fixture file.
		// There is no user input.
		bz, err := ioutil.ReadFile(fixtureFilePath)
		if err != nil {
			return nil, fmt.Errorf(
				"could not open the test fixture for party [%d] "+
					"in the expected location [%s]: [%w]",
				j,
				fixtureFilePath,
				err,
			)
		}
		var share keygen.LocalPartySaveData
		if err = json.Unmarshal(bz, &share); err != nil {
			return nil, fmt.Errorf(
				"could not unmarshal fixture data for party [%d] "+
					"located at [%s]: [%w]",
				j,
				fixtureFilePath,
				err,
			)
		}
		shares = append(shares, share)
	}
	return shares, nil
}
