package zkp

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/keep-network/paillier"
)

func TestRoundTrip(t *testing.T) {
	secret := big.NewInt(430)

	p, q, err := paillier.GenerateSafePrimes(256, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	privateKey := paillier.CreatePrivateKey(p, q)

	NTilde := new(big.Int).Mul(p, q)

	h1, h2, err := getH1H2(p, q, NTilde)
	if err != nil {
		t.Fatal(err)
	}

	params := &ZKPPublicParameters{
		N:      privateKey.N,
		N2:     privateKey.GetNSquare(),
		NTilde: NTilde,
		G:      new(big.Int).Add(privateKey.N, big.NewInt(1)),

		h1: h1,
		h2: h2,

		q:     secp256k1.S256().Params().N,
		curve: secp256k1.S256(),
	}

	eta, err := rand.Int(rand.Reader, params.q)
	if err != nil {
		t.Fatalf("could not generate eta [%v]", err)
	}

	r, err := paillier.GetRandomNumberInMultiplicativeGroup(params.N, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	c2, err := privateKey.EncryptWithR(secret, r)
	if err != nil {
		t.Fatal(err)
	}

	c1 := new(big.Int).Exp(c2.C, eta, params.N2)

	// c3, err := params.privateKey.EncryptWithR(eta, r)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// c3 = (Γ^η)*(r^N) mod N2
	c3 := new(big.Int).Mul(new(big.Int).Exp(params.G, eta, params.N2), new(big.Int).Exp(r, params.N, params.N2))

	zkp := new(PI1)
	// if params.N.Cmp(new(big.Int).Exp(params.q, big.NewInt(8), nil)) < 1 {
	// 	t.Fatalf("N is not bigger than q^8")
	// }

	zkp.Commit(eta, r, c1, c2.C, c3, params)

	if !zkp.Verify(c1, c2.C, c3, params) {
		t.Fatalf("ERROR")
	}
}

func getH1H2(p, q, NTilde *big.Int) (h1, h2 *big.Int, err error) {
	p, q, _ := paillier.GenerateSafePrimes(256/2, rand.Reader)

	NTilde := new(big.Int).Mul(p, q)

	// Fujisaki Osamoto : Chapter 3.1
	// Odd prime divisors
	pDivisor := new(big.Int).Div(
		new(big.Int).Add(p, big.NewInt(-1)),
		big.NewInt(2),
	)

	qDivisor := new(big.Int).Div(
		new(big.Int).Add(q, big.NewInt(-1)),
		big.NewInt(2),
	)

	pqDivisor := new(big.Int).Mul(pDivisor, qDivisor)

	// Fujisaki Osamoto : Chapter 3.1 - b0
	h2, err = randomFromMultiplicativeGroup(rand.Reader, NTilde)
	if err != nil {
		return nil, nil, err
	}

	// Fujisaki Osamoto : Chapter 3.1 - alpha
	x, err := rand.Int(rand.Reader, pqDivisor)
	if err != nil {
		return nil, nil, err
	}

	// Fujisaki Osamoto : Chapter 3.1 - b1 = (b0)^alpha mod N
	h1 = new(big.Int).Exp(h2, x, NTilde)
	return
}
