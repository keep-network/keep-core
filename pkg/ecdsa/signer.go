package ecdsa

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"

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

// Signer represents T-ECDSA group member in a fully initialized state,
// ready for signing.
type Signer struct {
	LocalSigner
}

// LocalSigner represents T-ECDSA group member prior to the initialisation
// phase.
type LocalSigner struct {
	ID         string
	paillerKey *paillier.ThresholdPrivateKey
}

func newGroup(parameters *PublicParameters) ([]*LocalSigner, error) {
	paillierKeyGen := paillier.GetThresholdKeyGenerator(
		paillierModulusBitLength,
		parameters.groupSize,
		parameters.dishonestThreshold,
		crand.Reader,
	)

	paillierKeys, err := paillierKeyGen.Generate()
	if err != nil {
		return nil, fmt.Errorf(
			"Could not generate threshold Paillier keys [%v]", err,
		)
	}

	members := make([]*LocalSigner, len(paillierKeys))
	for i := 0; i < len(members); i++ {
		members[i] = &LocalSigner{
			ID:         generateMemberID(),
			paillerKey: paillierKeys[i],
		}
	}

	return members, nil
}

func generateMemberID() string {
	memberID := "0"
	for memberID = fmt.Sprintf("%v", mrand.Int31()); memberID == "0"; {
	}
	return memberID
}
