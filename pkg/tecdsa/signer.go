// Package tecdsa contains the code that implements Threshold ECDSA signatures.
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

// PublicParameters for T-ECDSA. Defines how many Signers are in the group
// and what is a group signing threshold.
//
// If we consider an honest-but-curious adversary, i.e. an adversary that learns
// all the secret data of compromised server but does not change their code,
// then [GGN 16] protocol produces signature with n = t + 1 players in the
// network (since all players will behave honestly, even the corrupted ones).
// But in the presence of a malicious adversary, who can force corrupted players
// to shut down or send incorrect messages, one needs at least n = 2t + 1
// players in total to guarantee robustness, i.e. the ability to generate
// signatures even in the presence of malicous faults.
//
// Threshold is just for signing. If anything goes wrong during key generation,
// e.g. one of ZKP fails or any commitment opens incorrectly, key generation
// protocol terminates without an output.
//
// The Curve specified in the PublicParameters is the one used for signing and
// all intermedite constructions during initialization and signing process.
type PublicParameters struct {
	groupSize int
	threshold int

	curve elliptic.Curve
}

// LocalSigner represents T-ECDSA group member prior to the initialisation
// phase. It is responsible for constructing a broadcast InitMessage containing
// DSA key coproducts. Each LocalSigner has a reference to a threshold Paillier
// key used for encrypting part of the InitMessage.
type LocalSigner struct {
	ID               string
	publicParameters *PublicParameters
	paillerKey       *paillier.ThresholdPrivateKey
}

// Signer represents T-ECDSA group member in a fully initialized state,
// ready for signing. Each Signer has a reference to a ThresholdDsaKey used
// for a signing process. It represents a (t, n) threshold sharing of the
// underlying DSA key.
type Signer struct {
	LocalSigner

	dsaKey *ThresholdDsaKey
}

// In order for the [GGN 16] protocol to be correct, all the homomorphic
// operations over the ciphertexts (which are modulo N) must not conflict with
// the operations modulo q of the DSA algorithms. Becase of that, [GGN 16]
// requires that N > q^8.
//
// secp256k1 cardinality q is a 256 bit number, so we must have at least
// 2048 bit Paillier modulus.
// TODO: Boost prime generator performance and switch to 2048
const paillierModulusBitLength = 256

// generateDsaKeyShare generates a DSA key share coproducts xi and yi and puts
// them into dsaKeyShare. xi is a random integer from Z_q where q is the
// cardinality of Elliptic Curve and yi is a random point on the Curve.
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

// InitializeDsaKeyGen initializes key generation process by generating DSA key
// coproducts and putting them into the InitMessage which is broadcasted to all
// other Signers in the group.
//
// Each LocalSigner i selects a random value xi from Z_q, where q is the
// cardinality of used Elliptic Curve, compute yi = g^xi and E(x) where E is an
// additively homomorphic scheme encryption. In our case, it's Paillier.
//
// E(x) and yi are put into the InitMessage and sent to all other Signers in the
// group.
func (s *LocalSigner) InitializeDsaKeyGen() (*InitMessage, error) {
	keyShare, err := s.generateDsaKeyShare()
	if err != nil {
		return nil, fmt.Errorf("could not initialize DSA key genration [%v]", err)
	}

	exi, err := s.paillerKey.Encrypt(keyShare.xi, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("could not initialize DSA key genration [%v]", err)
	}

	return &InitMessage{
		xi: exi,
		yi: keyShare.yi,
	}, nil
}

// CombineDsaKeyShares combines all group InitMessages into a ThresholdDsaKey.
// ThresholdDsaKey is a (t, n) threshold sharing of an underlying (x, y) DSA
// key. Shares are combined in the following way:
//
// E(x) = E(x1) + E(x2) + ... + E(xn)
// y = y1 + y2 + ... + yn
//
// E is an additively homomorphic encryption scheme, hence + operation is
// possible. Each E(xn) share comes from InitMessage that was created by each
// LocalSigner of the signing group.
//
// y is a sum of all yn EllipticCurve points which were points generated by
// each LocalMember of the signinig group along with E(xn).
func (s *LocalSigner) CombineDsaKeyShares(shares []*InitMessage) (*ThresholdDsaKey, error) {
	// TODO: check ZKPs and required shares number

	xiShares := make([]*paillier.Cypher, len(shares))
	for i, share := range shares {
		xiShares[i] = share.xi
	}
	x := s.paillerKey.Add(xiShares...)

	yx := shares[0].yi.x
	yy := shares[0].yi.y
	for _, share := range shares[1:] {
		yx, yy = s.publicParameters.curve.Add(yx, yy, share.yi.x, share.yi.y)
	}

	return &ThresholdDsaKey{
		x: x,
		y: &CurvePoint{
			x: yx,
			y: yy,
		},
	}, nil
}

// newGroup generates a new signing group backed by a threshold Paillier key.
// This implementation works in an oracle mode - one party is responsible for
// generating Paillier keys and distributing them. Be careful please.
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
