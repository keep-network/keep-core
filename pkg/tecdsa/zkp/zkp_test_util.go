package zkp

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// mockRandReader is an implementation of `io.Reader` allowing to get
// predictable random numbers in your tests. Each new generated number is larger
// by 1 from the previous one starting from counter seed provided when
// constructing mockRandReader.
//
// We use `mockRandomReader` to test commitment phase of ZKPs defined in this
// package where we need predictable values instead of random ones.
//
// mockRandom := &mockRandReader{ counter: big.NewInt(1) }
// r1, _ := rand.Int(mockRandom, big.NewInt(10000)) // r1=1
// r2, _ := rand.Int(mockRandom, big.NewInt(10000)) // r2=2
// r3, _ := rand.Int(mockRandom, big.NewInt(10000)) // r3=3
type mockRandReader struct {
	counter *big.Int
}

func (r *mockRandReader) Read(b []byte) (int, error) {
	cb := r.counter.Bytes()

	for i := range b {
		// iterate backwards
		bIdx := len(b) - i - 1
		cbIdx := len(cb) - i - 1

		if cbIdx >= 0 {
			b[bIdx] = cb[cbIdx]
		}
	}

	r.counter = new(big.Int).Add(r.counter, big.NewInt(1))
	return len(b), nil
}

func generateTestPublicParams() *PublicParameters {
	return &PublicParameters{
		N:      big.NewInt(1081),  // 23 * 47
		NTilde: big.NewInt(25651), // 23 * 11

		h1: big.NewInt(20535),
		h2: big.NewInt(20919),

		q:     secp256k1.S256().Params().N,
		curve: secp256k1.S256(),
	}
}
