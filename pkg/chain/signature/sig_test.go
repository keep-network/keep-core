package signature

import "testing"

func TestSignature(t *testing.T) {

	tests := []struct {
		in                  []byte
		addr                string
		expectedMsg         string
		expectedEIP55Addr   string
		expectedPubKey      string
		keyFile             string
		password            string
		expectError         bool
		expectNotToValidate bool
	}{
		{
			in:                  []byte{01, 02, 03, 04},
			addr:                "6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			expectedMsg:         "01020304",
			expectedEIP55Addr:   "0x6FFBA2D0F4C8FD7961F516af43C55fe2d56f6044",
			expectedPubKey:      "0483bb5756ae8c2e9a4345682e38d585f76a769f5ba3e08505c1a1338c05edf800baf45ad8d256aeb74ee2fa6f52aa4a02621a95e208c263884beca60d8543bc4e",
			keyFile:             "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			password:            "password",
			expectError:         false,
			expectNotToValidate: false,
		},
		{
			in:                []byte{01, 02, 03, 04},
			addr:              "6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			expectedMsg:       "01020304",
			expectedEIP55Addr: "0x6FFBA2D0F4C8FD7961F516af43C55fe2d56f6044",
			expectedPubKey:    "0483bb5756ae8c2e9a4345682e38d585f76a769f5ba3e08505c1a1338c05edf800baf45ad8d256aeb74ee2fa6f52aa4a02621a95e208c263884beca60d8543bc4e",
			keyFile:           "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			password:          "nanananana",
			expectError:       true,
		},
		{
			in:                  []byte{01, 02, 03, 04},
			addr:                "9ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			expectedMsg:         "01020304",
			expectedEIP55Addr:   "0x6FFBA2D0F4C8FD7961F516af43C55fe2d56f6044",
			expectedPubKey:      "0483bb5756ae8c2e9a4345682e38d585f76a769f5ba3e08505c1a1338c05edf800baf45ad8d256aeb74ee2fa6f52aa4a02621a95e208c263884beca60d8543bc4e",
			keyFile:             "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
			password:            "password",
			expectError:         false,
			expectNotToValidate: true,
		},
	}

	for ii, test := range tests {
		msg, sig, err := GenerateSignature(test.keyFile, test.password, test.in)
		if test.expectError {
			if err == nil {
				t.Errorf("Test %d, failed to returne an error [%v] \n", ii, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("Test %d, returned an error [%v] \n", ii, err)
		}
		if msg != test.expectedMsg {
			t.Errorf("Test %d, expected %s got %s\n", ii, test.expectedMsg, msg)
		}
		ra, pk, sigValid, err := VerifySignature(test.addr, sig, msg)
		if test.expectNotToValidate {
			if sigValid {
				t.Errorf("Test %d, should not have validated but did\n", ii)
			}
			continue
		}
		if err != nil {
			t.Errorf("Test %d, valied to verify [%v] \n", ii, err)
		}
		if ra != test.expectedEIP55Addr {
			t.Errorf("Test %d, expected %s got %s\n", ii, test.expectedEIP55Addr, ra)
		}
		if pk != test.expectedPubKey {
			t.Errorf("Test %d, expected %s got %s\n", ii, test.expectedEIP55Addr, ra)
		}
	}

}
