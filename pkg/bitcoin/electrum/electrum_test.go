package electrum

import (
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"testing"
)

func TestConvertBtcKbToSatVByte(t *testing.T) {
	var tests = map[string]struct {
		btcPerKbFee            float32
		expectedSatPerVByteFee int64
	}{
		"BTC/KB is negative": {
			btcPerKbFee:            -1,
			expectedSatPerVByteFee: 1,
		},
		"BTC/KB is 0": {
			btcPerKbFee:            0,
			expectedSatPerVByteFee: 1,
		},
		"BTC/KB is 0.000001": {
			btcPerKbFee:            0.000001,
			expectedSatPerVByteFee: 1,
		},
		"BTC/KB is 0.00001": {
			btcPerKbFee:            0.00001,
			expectedSatPerVByteFee: 1,
		},
		"BTC/KB is 0.00002": {
			btcPerKbFee:            0.00002,
			expectedSatPerVByteFee: 2,
		},
		"BTC/KB is 0.0001": {
			btcPerKbFee:            0.0001,
			expectedSatPerVByteFee: 10,
		},
		"BTC/KB is 0.001": {
			btcPerKbFee:            0.001,
			expectedSatPerVByteFee: 100,
		},
		"BTC/KB is 0.0012350": {
			btcPerKbFee:            0.0012350,
			expectedSatPerVByteFee: 123,
		},
		"BTC/KB is 0.0012351": {
			btcPerKbFee:            0.0012351,
			expectedSatPerVByteFee: 124,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			satPerVByteFee := convertBtcKbToSatVByte(test.btcPerKbFee)

			testutils.AssertIntsEqual(
				t,
				"sat/vbyte fee",
				int(test.expectedSatPerVByteFee),
				int(satPerVByteFee),
			)
		})
	}
}
