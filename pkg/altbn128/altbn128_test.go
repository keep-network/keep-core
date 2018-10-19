package altbn128

import (
	"bytes"
	"crypto/rand"
	"testing"
	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func TestCompressG1(t *testing.T) {
	for i := 0; i<100; i++ {
		_, p, err := bn256.RandomG1(rand.Reader)

		if err != nil {
			t.Errorf("Error generating random point on G1")
		}

		buffer := Compress(p)
		assertEqual(t, len(buffer), 32, "Compressed G1 should be 32 bytes")
	}
}

func TestDecompressG1(t *testing.T) {
	errorSeen := false
	for i := 0; i<100; i++ {
		buffer := make([]byte, 32)
		_, err := rand.Read(buffer)
		if err == nil {
			_, err2 := Decompress(buffer)


			if err2 == nil {
				errorSeen = true
			}
		}
	}
	if !errorSeen {
		t.Errorf("No errors seen decompressing random points on G1. Highly unlikely")
	}
}

func TestCompressG1Invertibility(t *testing.T) {
	for i := 0; i<100; i++ {
		_, p1, err1 := bn256.RandomG1(rand.Reader)

		if err1 != nil {
			continue
		}

		buffer := Compress(p1)

		t.Logf("Compressed G1 to [%v]", buffer)

		p2, _ := Decompress(buffer)

		assertPointsEqual(t, p1, p2, "Decompressing a compressed point should give the same point.")
	}
}

func assertPointsEqual(t *testing.T, p1 *bn256.G1, p2 *bn256.G1, msg string) {
	if p1 != p2 && !bytes.Equal(p1.Marshal(), p2.Marshal()) {
		t.Errorf("%v: [%v] != [%v]", msg, p1, p2)
	}
}

func assertEqual(t *testing.T, n int, n2 int, msg string) {
	if n != n2 {
		t.Errorf("%v: [%v] != [%v]", msg, n, n2)
	}
}
