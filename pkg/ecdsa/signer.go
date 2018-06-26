package ecdsa

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/keep-network/paillier"
)

// N is of length 2048, making the operations mod N^2 of length 4096
const paillierModulusBitLength = 2048

// PublicParameters for T-ECDSA
type PublicParameters struct {
	groupSize          int
	dishonestThreshold int

	q *big.Int // EC cardinality
}

// InitializedSigner represents T-ECDSA group member in a fully initialized
// state, ready for signing.
type InitializedSigner struct {
}

// InitializingSigner represents T-ECDSA group member in the initialisation
// phase.
type InitializingSigner struct {
	paillerKey *paillier.ThresholdPrivateKey
}

func newGroup(parameters *PublicParameters) ([]*InitializingSigner, error) {
	paillierKeyGen := paillier.GetThresholdKeyGenerator(
		paillierModulusBitLength,
		parameters.groupSize,
		parameters.dishonestThreshold,
		rand.Reader,
	)

	paillierKeys, err := paillierKeyGen.Generate()
	if err != nil {
		return nil, fmt.Errorf(
			"Could not generate threshold Paillier keys [%v]", err,
		)
	}

	members := make([]*InitializingSigner, len(paillierKeys))
	for i := 0; i < len(members); i++ {
		members[i] = &InitializingSigner{
			paillerKey: paillierKeys[i],
		}
	}

	return members, nil
}
