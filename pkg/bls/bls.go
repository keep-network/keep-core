package bls

import (
	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

// AggregateG1Points aggregates array of G1 points into a single G1 point.
func AggregateG1Points(points []*bn256.G1) *bn256.G1 {
	result := new(bn256.G1)
	for _, point := range points {
		result.Add(result, point)
	}
	return result
}

// AggregateG2Points aggregates array of G2 points into a single G2 point.
func AggregateG2Points(points []*bn256.G2) *bn256.G2 {
	result := new(bn256.G2)
	for _, point := range points {
		result.Add(result, point)
	}
	return result
}
