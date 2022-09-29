package tecdsatest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
)

const (
	testDataDirFormat                 = "%s/testdata"
	privateKeyShareTestDataFileFormat = "private_key_share_data_%d.json"
)

// LoadPrivateKeyShareTestFixtures loads tECDSA private key share test data.
// Code copied from:
//
//	https://github.com/bnb-chain/tss-lib/blob/v1.3.3/ecdsa/keygen/test_utils.go#L36
//
// Test data JSON files were generated using the tECDSA DKG protocol from
// the `pkg/tecdsa/dkg` package and represent a signing group 3-of-5.
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
