package electrum

import (
	"testing"
)

func TestParseProtocol(t *testing.T) {
	tests := map[string]struct {
		value          string
		expectedResult Protocol
		expectedOk     bool
	}{
		"empty": {
			value:          "",
			expectedResult: Unknown,
			expectedOk:     false,
		},
		"lower case: tcp": {
			value:          "tcp",
			expectedResult: TCP,
			expectedOk:     true,
		},
		"upper case: TCP": {
			value:          "TCP",
			expectedResult: TCP,
			expectedOk:     true,
		},
		"mixed case: Tcp": {
			value:          "Tcp",
			expectedResult: TCP,
			expectedOk:     true,
		},
		"lower case: ssl": {
			value:          "ssl",
			expectedResult: SSL,
			expectedOk:     true,
		},
		"upper case: SSL": {
			value:          "SSL",
			expectedResult: SSL,
			expectedOk:     true,
		},
		"mixed case: SsL": {
			value:          "SsL",
			expectedResult: SSL,
			expectedOk:     true,
		},
		"invalid: tls": {
			value:          "tls",
			expectedResult: Unknown,
			expectedOk:     false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			actualResult, actualOk := ParseProtocol(test.value)

			if test.expectedResult != actualResult {
				t.Errorf(
					"unexpected result\nexpected: [%+v]\nactual:   [%+v]",
					test.expectedResult,
					actualResult,
				)
			}

			if test.expectedOk != actualOk {
				t.Errorf(
					"unexpected ok\nexpected: [%+v]\nactual:   [%+v]",
					test.expectedOk,
					actualOk,
				)
			}

		})
	}
}
