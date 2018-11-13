package altbn128

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func TestCompressG1(t *testing.T) {
	for i := 0; i < 100; i++ {
		_, p, err := bn256.RandomG1(rand.Reader)

		if err != nil {
			t.Errorf("Error generating random point on G1")
		}

		buffer := G1Point{p}.Compress()
		assertEqual(t, len(buffer), 32, "Compressed G1 should be 32 bytes")
	}
}

func TestDecompressG1(t *testing.T) {
	errorSeen := false
	for i := 0; i < 100; i++ {
		buffer := make(compressedPoint, 32)
		_, err := rand.Read(buffer)
		if err == nil {
			_, err2 := buffer.DecompressToG1()

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
	for i := 0; i < 100; i++ {
		_, p1, err1 := bn256.RandomG1(rand.Reader)

		if err1 != nil {
			continue
		}

		buffer := G1Point{p1}.Compress()

		t.Logf("Compressed G1 to [%v]", buffer)

		p2, _ := buffer.DecompressToG1()

		assertPointsEqual(t, p1, p2, "Decompressing a compressed point should give the same point.")
	}
}

func TestCompressG2Invertibility(t *testing.T) {
	for i := 0; i < 100; i++ {
		_, p1, err1 := bn256.RandomG2(rand.Reader)

		if err1 != nil {
			continue
		}

		buffer := G2Point{p1}.Compress()

		t.Logf("Compressed G2 to [%v]", buffer)

		p2, _ := buffer.DecompressToG2()

		assertG2PointsEqual(t, p1, p2, "Decompressing a compressed point should give the same point.")
	}
}

func assertPointsEqual(t *testing.T, p1 *bn256.G1, p2 *bn256.G1, msg string) {
	if p1 != p2 && !bytes.Equal(p1.Marshal(), p2.Marshal()) {
		t.Errorf("%v: [%v] != [%v]", msg, p1, p2)
	}
}

func assertG2PointsEqual(t *testing.T, p1 *bn256.G2, p2 *bn256.G2, msg string) {
	if p1 != p2 && !bytes.Equal(p1.Marshal(), p2.Marshal()) {
		t.Errorf("%v: [%v] != [%v]", msg, p1, p2)
	}
}

func assertEqual(t *testing.T, n int, n2 int, msg string) {
	if n != n2 {
		t.Errorf("%v: [%v] != [%v]", msg, n, n2)
	}
}
