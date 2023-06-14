package altbn128

import (
	"crypto/rand"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"github.com/keep-network/keep-core/internal/testutils"
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
		buffer := make([]byte, 32)
		_, err := rand.Read(buffer)
		if err == nil {
			_, err2 := DecompressToG1(buffer)

			if err2 == nil {
				errorSeen = true
			}
		}
	}
	if !errorSeen {
		t.Errorf("No errors seen decompressing random points on G1. Highly unlikely")
	}
}

func TestCompressDecompressGivesSameG1Point(t *testing.T) {
	for i := 0; i < 100; i++ {
		_, p1, err1 := bn256.RandomG1(rand.Reader)

		if err1 != nil {
			continue
		}

		buffer := G1Point{p1}.Compress()

		t.Logf("Compressed G1 to [%v]", buffer)

		p2, _ := DecompressToG1(buffer)

		testutils.AssertBytesEqual(t, p1.Marshal(), p2.Marshal())
	}
}

func TestCompressDecompressGivesSameG2Point(t *testing.T) {
	for i := 0; i < 100; i++ {
		_, p1, err1 := bn256.RandomG2(rand.Reader)

		if err1 != nil {
			continue
		}

		buffer := G2Point{p1}.Compress()

		t.Logf("Compressed G2 to [%v]", buffer)

		p2, _ := DecompressToG2(buffer)

		testutils.AssertBytesEqual(t, p1.Marshal(), p2.Marshal())
	}
}

func assertEqual(t *testing.T, n int, n2 int, msg string) {
	if n != n2 {
		t.Errorf("%v: [%v] != [%v]", msg, n, n2)
	}
}
