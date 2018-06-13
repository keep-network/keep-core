package bls

import (
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func TestCompressG1(t *testing.T) {
	for i := 0; i<100; i++ {
		buffer := p.Compress()
		assertEqual(len(buffer), 32)
	}
}

func TestDecompressG1(t *testing.T) {
	for i := 0; i<100; i++ {
		buffer, err := rand.Read(32)
		if err == nil {
			p := Decompress(buffer)
		}
	}
}

func TestCompressG1Invertibility(t *testing.T) {
	for i := 0; i<100; i++ {
		_, p1, err1 := bn256.RandomG1()

		if err1 != nil {
			continue
		}

		buffer, err2 := p1.Compress()

		if err2 != nil {
			t.Errorf("Error compressing point [%v]", p1)
		}

		p2, err3 := Decompress(buffer)

		if err3 != nil {
			t.Errorf("Error decompressing point [%v]", p1)
		}

		assertEqual(t, p1, p2, "Decompressing a compressed point should give the same point.")
	}
}

func assertEqual(t *testing.T, p1 *bn256.G1, p2 *bn256.G1, msg string) {
	if !p1.p.x.Equals(p2.p.x) || !p1.p.y.Equals(p2.p.y) {
		t.Errorf("%v: [%v] != [%v]", msg, p1, p2)
	}
}

func assertEqual(t *testing.T, n int, n2 int, msg string) {
	if n != n2 {
		t.Errorf("%v: [%v] != [%v]", msg, n, n2)
	}
}
