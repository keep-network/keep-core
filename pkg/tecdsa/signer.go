// Package tecdsa contains the code that implements threshold ECDSA signatures.
// The approach is based on [GGN 16].
//
//     [GGN 16]: Gennaro R., Goldfeder S., Narayanan A. (2016) Threshold-Optimal
//          DSA/ECDSA Signatures and an Application to Bitcoin Wallet Security.
//          In: Manulis M., Sadeghi AR., Schneider S. (eds) Applied Cryptography
//          and Network Security. ACNS 2016. Lecture Notes in Computer Science,
//          vol 9696. Springer, Cham
package tecdsa

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	mathrand "math/rand"

	"github.com/keep-network/paillier"
)

const paillierModulusBitLength = 128 // TODO: must be larger than q^8?

// PublicParameters for T-ECDSA
type PublicParameters struct {
	groupSize int
	threshold int

	curve elliptic.Curve
}

// Signer represents T-ECDSA group member in a fully initialized state,
// ready for signing.
type Signer struct {
	LocalSigner
}

// LocalSigner represents T-ECDSA group member prior to the initialisation
// phase.
type LocalSigner struct {
	ID               string
	publicParameters *PublicParameters
	paillerKey       *paillier.ThresholdPrivateKey
}

func (s *LocalSigner) generateDsaKeyShare() (*dsaKeyShare, error) {
	curveParams := s.publicParameters.curve.Params()

	xi, err := rand.Int(rand.Reader, curveParams.N)
	if err != nil {
		return nil, fmt.Errorf("could not generate DSA key share [%v]", err)
	}

	yxi, yyi := s.publicParameters.curve.ScalarBaseMult(xi.Bytes())

	return &dsaKeyShare{
		xi: xi,
		yi: &CurvePoint{
			x: yxi,
			y: yyi,
		},
	}, nil
}

func newGroup(parameters *PublicParameters) ([]*LocalSigner, error) {
	paillierKeyGen := paillier.GetThresholdKeyGenerator(
		paillierModulusBitLength,
		parameters.groupSize,
		parameters.threshold,
		rand.Reader,
	)

	paillierKeys, err := paillierKeyGen.Generate()
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate threshold Paillier keys [%v]", err,
		)
	}

	members := make([]*LocalSigner, len(paillierKeys))
	for i := 0; i < len(members); i++ {
		members[i] = &LocalSigner{
			ID:               generateMemberID(),
			paillerKey:       paillierKeys[i],
			publicParameters: parameters,
		}
	}

	return members, nil
}

func generateMemberID() string {
	memberID := "0"
	for memberID = fmt.Sprintf("%v", mathrand.Int31()); memberID == "0"; {
	}
	return memberID
}
