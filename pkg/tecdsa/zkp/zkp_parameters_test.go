package zkp

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func TestGeneratePublicParameters(t *testing.T) {
	paillierModulus := big.NewInt(13)
	curve := secp256k1.S256()

	params, err := GeneratePublicParameters(paillierModulus, curve)
	if err != nil {
		t.Fatal(err)
	}

	if params.NTilde.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("NTilde should be greater than 0. Actual: %v", params.NTilde)
	}
	if params.h1.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("h1 should be greater than 0. Actual: %v", params.h1)
	}
	if params.h2.Cmp(big.NewInt(0)) <= 0 {
		t.Errorf("h2 should be greater than 0. Actual: %v", params.h2)
	}
	if params.q.Cmp(secp256k1.S256().Params().N) != 0 {
		t.Errorf(
			"Unexpected q\nExpected: %v\nActual: %v",
			secp256k1.S256().Params().N,
			params.q,
		)
	}
	if params.N != paillierModulus {
		t.Errorf("Unexpected N\nExpected: %v\nActual: %v",
			paillierModulus,
			params.N,
		)
	}
	if params.curve != curve {
		t.Errorf("Unexpected Elliptic Curve used")
	}
}
