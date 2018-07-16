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
	"errors"
	"fmt"

	mathrand "math/rand"

	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/keep-core/pkg/tecdsa/zkp"
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
// signatures even in the presence of malicious faults.
//
// Threshold is just for signing. If anything goes wrong during key generation,
// e.g. one of ZKP fails or any commitment opens incorrectly, key generation
// protocol terminates without an output.
//
// The Curve specified in the PublicParameters is the one used for signing and
// all intermediate constructions during initialization and signing process.
type PublicParameters struct {
	groupSize int
	threshold int

	curve elliptic.Curve
}

// LocalSigner represents T-ECDSA group member prior to the initialization
// phase. It is responsible for constructing a broadcast InitMessage containing
// DSA key shares. Each LocalSigner has a reference to a threshold Paillier
// key used for encrypting part of the InitMessage.
type LocalSigner struct {
	ID              string
	groupParameters *PublicParameters
	zkpParameters   *zkp.PublicParameters
	paillierKey     *paillier.ThresholdPrivateKey
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
// the operations modulo q of the DSA algorithms. Because of that, [GGN 16]
// requires that N > q^8.
//
// secp256k1 cardinality q is a 256 bit number, so we must have at least
// 2048 bit Paillier modulus.
// TODO: Boost prime generator performance and switch to 2048
const paillierModulusBitLength = 256

// generateDsaKeyShare generates a DSA public and secret key shares and puts
// them into `dsaKeyShare`. Secret key share is a random integer from Z_q where
// `q` is the cardinality of Elliptic Curve and public key share is a point
// on the Curve g^secretKeyShare.
func (ls *LocalSigner) generateDsaKeyShare() (*dsaKeyShare, error) {
	curveParams := ls.groupParameters.curve.Params()

	secretKeyShare, err := rand.Int(rand.Reader, curveParams.N)
	if err != nil {
		return nil, fmt.Errorf("could not generate DSA key share [%v]", err)
	}

	publicKeyShare := curve.NewPoint(
		ls.groupParameters.curve.ScalarBaseMult(secretKeyShare.Bytes()),
	)

	return &dsaKeyShare{
		secretKeyShare: secretKeyShare,
		publicKeyShare: publicKeyShare,
	}, nil
}

// InitializeDsaKeyGen initializes key generation process by generating DSA key
// shares and putting them into the `InitMessage` which is broadcasted to all
// other `Signer`s in the group.
//
// Secret key share is encrypted with an additively homomorphic encryption
// scheme and sent to all other Signers in the group along with the public key
// share.
//
// Along with secret and public key share, we ship a zero knowledge argument
// allowing to validate received shares.
func (ls *LocalSigner) InitializeDsaKeyGen() (*KeyShareRevealMessage, error) {
	keyShare, err := ls.generateDsaKeyShare()
	if err != nil {
		return nil, fmt.Errorf(
			"could not initialize DSA key generation [%v]", err,
		)
	}

	paillierRandomness, err := paillier.GetRandomNumberInMultiplicativeGroup(
		ls.paillierKey.N, rand.Reader,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not initialize DSA key generation [%v]", err,
		)
	}

	encryptedSecretKeyShare, err := ls.paillierKey.EncryptWithR(
		keyShare.secretKeyShare, paillierRandomness,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not initialize DSA key generation [%v]", err,
		)
	}

	rangeProof, err := zkp.CommitDsaPaillierKeyRange(
		keyShare.secretKeyShare,
		keyShare.publicKeyShare,
		encryptedSecretKeyShare,
		paillierRandomness,
		ls.zkpParameters,
		rand.Reader,
	)

	return &KeyShareRevealMessage{
		secretKeyShare: encryptedSecretKeyShare,
		publicKeyShare: keyShare.publicKeyShare,
		rangeProof:     rangeProof,
	}, nil
}

// CombineDsaKeyShares combines all group `InitMessages` into a
// `ThresholdDsaKey` which is a (t, n) threshold sharing of an underlying secret
// and public DSA key shares. Secret and public DSA key shares are combined in
// the following way:
//
// E(secretKey) = E(secretKeyShare_1) + E(secretKeyShare_2) + ... + E(secretKeyShare_n)
// publicKey = publicKeyShare_1 + publicKeyShare_2 + ... + publicKeyShare_n
//
// E is an additively homomorphic encryption scheme, hence `+` operation is
// possible. `Each E(secretKeyShare_i)` share comes from `InitMessage` that was
// created by each `LocalSigner` of the signing group.
func (ls *LocalSigner) CombineDsaKeyShares(
	shares []*KeyShareRevealMessage) (*ThresholdDsaKey, error) {
	if len(shares) != ls.groupParameters.groupSize {
		return nil, fmt.Errorf(
			"InitMessages required from all group members; Got %v, expected %v",
			len(shares),
			ls.groupParameters.groupSize,
		)
	}

	for _, share := range shares {
		if !share.IsValid(ls.zkpParameters) {
			return nil, errors.New("Invalid InitMessage - ZKP rejected")
		}
	}

	secretKeyShares := make([]*paillier.Cypher, len(shares))
	for i, share := range shares {
		secretKeyShares[i] = share.secretKeyShare
	}
	secretKey := ls.paillierKey.Add(secretKeyShares...)

	publicKeyShareX := shares[0].publicKeyShare.X
	publicKeyShareY := shares[0].publicKeyShare.Y
	for _, share := range shares[1:] {
		publicKeyShareX, publicKeyShareY = ls.groupParameters.curve.Add(
			publicKeyShareX, publicKeyShareY,
			share.publicKeyShare.X, share.publicKeyShare.Y,
		)
	}

	return &ThresholdDsaKey{
		secretKey: secretKey,
		publicKey: &curve.Point{
			X: publicKeyShareX,
			Y: publicKeyShareY,
		},
	}, nil
}

// newGroup generates a new signing group backed by a threshold Paillier key
// and ZKP public parameters built from the generated Paillier key.
// This implementation works in an oracle mode - one party is responsible for
// generating Paillier keys and distributing them. Be careful, please.
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

	zkpParameters, err := zkp.GeneratePublicParameters(
		paillierKeys[0].N,
		parameters.curve,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate public ZKP parameters [%v]", err,
		)
	}

	members := make([]*LocalSigner, len(paillierKeys))
	for i := 0; i < len(members); i++ {
		members[i] = &LocalSigner{
			ID:              generateMemberID(),
			paillierKey:     paillierKeys[i],
			groupParameters: parameters,
			zkpParameters:   zkpParameters,
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
