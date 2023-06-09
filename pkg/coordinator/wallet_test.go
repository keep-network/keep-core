package coordinator

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

var walletPublicKeyHashTests = []struct {
	input          string
	expectedResult WalletPublicKeyHash
	wantErr        error // if set, decoding must fail
}{
	// invalid
	{input: ``, wantErr: fmt.Errorf("empty hex string")},
	{input: `01`, wantErr: fmt.Errorf("invalid bytes length: [1], expected: [20]")},
	{input: `0x01`, wantErr: fmt.Errorf("invalid bytes length: [1], expected: [20]")},
	{input: `0x5bee2805df9fcea4691c442fe4c1a33f7288e2`, wantErr: fmt.Errorf("invalid bytes length: [19], expected: [20]")},
	{input: `000f4224b6858eee7f8999e6299c056c6405bbede0`, wantErr: fmt.Errorf("invalid bytes length: [21], expected: [20]")},
	// valid
	{input: `0x48b88e1074c33c7a934f781220e1a4523f1768c0`, expectedResult: WalletPublicKeyHash{72, 184, 142, 16, 116, 195, 60, 122, 147, 79, 120, 18, 32, 225, 164, 82, 63, 23, 104, 192}},
	{input: `0x00008e1074c33c7a934f781220e1a4523f1768c0`, expectedResult: WalletPublicKeyHash{00, 00, 142, 16, 116, 195, 60, 122, 147, 79, 120, 18, 32, 225, 164, 82, 63, 23, 104, 192}},
	{input: `0x48b88e1074c33c7a934f781220e1a4523f000000`, expectedResult: WalletPublicKeyHash{72, 184, 142, 16, 116, 195, 60, 122, 147, 79, 120, 18, 32, 225, 164, 82, 63, 00, 00, 00}},
}

func TestWalletPublicKeyHashRoundTrip(t *testing.T) {
	for _, test := range walletPublicKeyHashTests {
		t.Run(test.input, func(t *testing.T) {
			actualResult, err := NewWalletPublicKeyHash(test.input)
			if !reflect.DeepEqual(err, test.wantErr) {
				t.Fatalf("unexpected error\nexpected: %v\nactual:   %v", test.wantErr, err)
			}

			testutils.AssertBytesEqual(t, test.expectedResult[:], actualResult[:])

			if test.wantErr == nil {
				if actualResult.String() != test.input {
					t.Fatalf(
						"unexpected string\nexpected: %v\nactual:   %v",
						test.input,
						actualResult.String(),
					)
				}
			}
		})
	}
}
