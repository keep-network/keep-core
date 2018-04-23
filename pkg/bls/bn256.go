package bls

import (
	"crypto/sha256"
    "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare/bn256"
	"math/big"
)

func sum(ints ...*big.Int) *big.Int {
	acc := big.NewInt(0)
	for _, num := range ints {
		acc.Add(acc, num)
	}
	return acc
}

func product(ints ...*big.Int) *big.Int {
	acc := big.NewInt(1)
	for _, num := range ints {
		acc.Mul(acc, num)
	}
	return acc
}

func negate(i *big.Int) *big.Int {
	return new(big.Int).Neg(i)
}

func mod(i, m *big.Int) *big.Int {
	return new(big.Int).Mod(i, m)
}

func modSqrt(i, m *big.Int) *big.Int {
	return new(big.Int).ModSqrt(i, m)
}

func modInverse(i, m *big.Int) *big.Int {
	return new(big.Int).ModInverse(i, m)
}

func G1HashToPoint(m []byte) (*big.Int, *big.Int) {

	zero, one := big.NewInt(0), big.NewInt(1)

	sqrtNeg3 := modSqrt(sum(P, negate(big.NewInt(3))))
	inverse2 := modInverse(big.NewInt(2), bn256.P)
	// TODO get b from curve
	b := one

	h := sha256.Sum256(m)

	t := mod(new(big.Int).SetBytes(h[:]), bn256.P)

	if t.Cmp(zero) == 0 {
		// TODO handle odd case
	}

	t2ModP := mod(product(t, t), bn256.P)

	chi_t := Legendre(t)

	tInvMod := modInverse(sum(one, b, t2ModP), bn256.P)
	w := product(sqrtNeg3, t, tInvMod)

	x1 := mod(sum(product(sum(sqrtNeg3, negate(one)), inverse2), product(t, w)), bn256.P)
	x1CubedPlusB := sum(product(x1, x1, x1), b)

	if Legendre(x1CubedPlusB).Comp(1) == 0 {
		x1Sqrt := modSqrt(x1CubedPlusB, bn256.P)
		// TODO return
		return big.NewInt(0), big.NewInt(0)
	}

	x2 := mod(sum(negate(one), negate(x1)), bn256.P)
	x2CubedPlusB := sum(product(x2, x2, x2), b)

	if Legendre(x2CubedPlusB).Comp(1) == 0 {
		x2Sqrt := modSqrt(x2CubedPlusB, bn256.P)
		// TODO return
		return big.NewInt(0), big.NewInt(0)
	}

	x3 := sum(one, modInverse(product(w, w)))
	x3CubedPlusB := sum(product(x3, x3, x3), b)

	if Legendre(x3CubedPlusB).Comp(1) == 0 {
		x3Sqrt := modSqrt(x3CubedPlusB, bn256.P)
		// TODO return
		return big.NewInt(0), big.NewInt(0)
	}
	// panic
}

func Legendre(a *big.Int) *big.Int {
	zero, one := big.NewInt(0), big.NewInt(1)
	pMinus := new(big.Int).Sub(P, one)
	x := new(big.Int).Exp(a, pMinus, bn256.P)
	if x.Cmp(one) == 0 || x.Cmp(zero) == 0 {
		return x
	}
	if x.Cmp(pMinus) == 0 {
		return big.NewInt(-1)
	}
	// assert should not happen
}
