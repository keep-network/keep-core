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

	"math/big"
	mathrand "math/rand"

	"github.com/keep-network/keep-core/pkg/tecdsa/commitment"
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

type signerCore struct {
	ID string

	paillierKey *paillier.ThresholdPrivateKey

	groupParameters *PublicParameters
	zkpParameters   *zkp.PublicParameters
}

// LocalSigner represents T-ECDSA group member during the initialization
// phase. It is responsible for constructing a broadcast
// PublicKeyShareCommitmentMessage containing public DSA key share commitment
// and a KeyShareRevealMessage revealing in a Paillier-encrypted way generated
// secret DSA key share and an unencrypted public key share.
type LocalSigner struct {
	signerCore

	dsaKeyShare *dsaKeyShare

	// Intermediate value stored between first and second round of
	// key generation. In the first round, `LocalSigner` commits to the chosen
	// public key share. In the second round, it reveals the public key share
	// along with the decommitment key.
	publicDsaKeyShareDecommitmentKey *commitment.DecommitmentKey
}

// Signer represents T-ECDSA group member in a fully initialized state,
// ready for signing. Each Signer has a reference to a ThresholdDsaKey used
// in a signing process. It represents a (t, n) threshold sharing of the
// underlying DSA key.
type Signer struct {
	signerCore

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

func (pp *PublicParameters) curveCardinality() *big.Int {
	return pp.curve.Params().N
}

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

// InitializeDsaKeyShares initializes key generation process by generating DSA
// key shares and publishing PublicKeyShareCommitmentMessage which is
// broadcasted to all other `Signer`s in the group and contains signer's public
// DSA key share commitment.
func (ls *LocalSigner) InitializeDsaKeyShares() (
	*PublicKeyShareCommitmentMessage,
	error,
) {
	keyShare, err := ls.generateDsaKeyShare()
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate DSA key shares [%v]", err,
		)
	}

	commitment, decommitmentKey, err := commitment.Generate(
		keyShare.publicKeyShare.Bytes(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate DSA public key commitment [%v]", err,
		)
	}

	ls.dsaKeyShare = keyShare
	ls.publicDsaKeyShareDecommitmentKey = decommitmentKey

	return &PublicKeyShareCommitmentMessage{
		signerID:   ls.ID,
		commitment: commitment,
	}, nil
}

// RevealDsaKeyShares produces a KeyShareRevealMessage and should be called
// when `PublicKeyShareCommitmentMessage`s from all group members are gathered.
//
// `KeyShareRevealMessage` contains signer's public DSA key share, decommitment
// key for this share (used to validate the commitment published in the previous
// `PublicKeyShareCommitmentMessage` message), encrypted secret DSA key share
// and ZKP for the secret key share correctness.
//
// Secret key share is encrypted with an additively homomorphic encryption
// scheme and sent to all other Signers in the group along with the public key
// share.
func (ls *LocalSigner) RevealDsaKeyShares() (*KeyShareRevealMessage, error) {
	paillierRandomness, err := paillier.GetRandomNumberInMultiplicativeGroup(
		ls.paillierKey.N, rand.Reader,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not generate random r for Paillier [%v]", err,
		)
	}

	encryptedSecretKeyShare, err := ls.paillierKey.EncryptWithR(
		ls.dsaKeyShare.secretKeyShare, paillierRandomness,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not encrypt secret key share [%v]", err,
		)
	}

	rangeProof, err := zkp.CommitDsaPaillierKeyRange(
		ls.dsaKeyShare.secretKeyShare,
		ls.dsaKeyShare.publicKeyShare,
		encryptedSecretKeyShare,
		paillierRandomness,
		ls.zkpParameters,
		rand.Reader,
	)

	return &KeyShareRevealMessage{
		signerID:                 ls.ID,
		secretKeyShare:           encryptedSecretKeyShare,
		publicKeyShare:           ls.dsaKeyShare.publicKeyShare,
		publicKeyDecommitmentKey: ls.publicDsaKeyShareDecommitmentKey,
		secretKeyProof:           rangeProof,
	}, nil
}

// CombineDsaKeyShares combines all group `PublicKeyShareCommitmentMessage`s and
// `KeyShareRevealMessage`s into a `ThresholdDsaKey` which is a (t, n) threshold
// sharing of an underlying secret DSA key. Secret and public
// DSA key shares are combined in the following way:
//
// E(secretKey) = E(secretKeyShare_1) + E(secretKeyShare_2) + ... + E(secretKeyShare_n)
// publicKey = publicKeyShare_1 + publicKeyShare_2 + ... + publicKeyShare_n
//
// E is an additively homomorphic encryption scheme, hence `+` operation is
// possible. Each key share share comes from the `KeyShareRevealMessage` that
// was sent by each `LocalSigner` of the signing group.
//
// Before shares are combined, messages are validated - we check whether
// the published public key share is what the signer originally committed to
// as well as we check validity of the secret key share using the provided ZKP.
//
// Every `PublicKeyShareCommitmentMessage` should have a corresponding
// `KeyShareRevealMessage`. They are matched by a signer ID contained in
// each of the messages.
func (ls *LocalSigner) CombineDsaKeyShares(
	shareCommitments []*PublicKeyShareCommitmentMessage,
	revealedShares []*KeyShareRevealMessage,
) (*ThresholdDsaKey, error) {
	if len(shareCommitments) != ls.groupParameters.groupSize {
		return nil, fmt.Errorf(
			"commitments required from all group members; got %v, expected %v",
			len(shareCommitments),
			ls.groupParameters.groupSize,
		)
	}

	if len(revealedShares) != ls.groupParameters.groupSize {
		return nil, fmt.Errorf(
			"all group members should reveal shares; Got %v, expected %v",
			len(revealedShares),
			ls.groupParameters.groupSize,
		)
	}

	secretKeyShares := make([]*paillier.Cypher, ls.groupParameters.groupSize)
	publicKeyShares := make([]*curve.Point, ls.groupParameters.groupSize)

	for i, commitmentMsg := range shareCommitments {
		foundMatchingRevealMessage := false

		for _, revealedSharesMsg := range revealedShares {

			if commitmentMsg.signerID == revealedSharesMsg.signerID {
				foundMatchingRevealMessage = true

				if revealedSharesMsg.isValid(
					commitmentMsg.commitment, ls.zkpParameters,
				) {
					secretKeyShares[i] = revealedSharesMsg.secretKeyShare
					publicKeyShares[i] = revealedSharesMsg.publicKeyShare
				} else {
					return nil, errors.New("KeyShareRevealMessage rejected")
				}
			}
		}

		if !foundMatchingRevealMessage {
			return nil, fmt.Errorf(
				"no matching share reveal message for signer with ID=%v",
				commitmentMsg.signerID,
			)
		}
	}

	secretKey := ls.paillierKey.Add(secretKeyShares...)
	publicKey := publicKeyShares[0]
	for _, share := range publicKeyShares[1:] {
		publicKey = curve.NewPoint(ls.groupParameters.curve.Add(
			publicKey.X, publicKey.Y, share.X, share.Y,
		))
	}

	return &ThresholdDsaKey{secretKey, publicKey}, nil
}

func generateMemberID() string {
	memberID := "0"
	for memberID = fmt.Sprintf("%v", mathrand.Int31()); memberID == "0"; {
	}
	return memberID
}

// Round1Signer represents state of `Signer` after executing the first round
// of signing algorithm.
type Round1Signer struct {
	Signer

	// Intermediate values stored between the first and second round of signing.
	randomFactorShare           *big.Int                    // ρ_i
	encryptedRandomFactorShare  *paillier.Cypher            // u_i = E(ρ_i)
	secretKeyMultiple           *paillier.Cypher            // v_i = E(ρ_i * x)
	randomFactorDecommitmentKey *commitment.DecommitmentKey // D_1i
	paillierRandomness          *big.Int
}

// SignRound1 executes the first round of T-ECDSA signing as described in
// [GGN 16], section 4.3.
//
// In the first round, each signer generates a random factor share `ρ_i`,
// encodes it with Paillier key `u_i = E(ρ_i)`, multiplies it with secret ECDSA
// key `v_i = E(ρ_i * x)` and publishes commitment for both those values
// `Com(u_i, v_i)`.
func (s *Signer) SignRound1() (*Round1Signer, *SignRound1Message, error) {
	// Choosing random ρ_i from Z_q
	randomFactorShare, err := rand.Int(
		rand.Reader,
		s.groupParameters.curveCardinality(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 1 of signing [%v]", err,
		)
	}

	paillierRandomness, err := paillier.GetRandomNumberInMultiplicativeGroup(
		s.paillierKey.N, rand.Reader,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 1 of signing [%v]", err,
		)
	}

	// u_i = E(ρ_i)
	encryptedRandomFactorShare, err := s.paillierKey.EncryptWithR(
		randomFactorShare, paillierRandomness,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 1 of signing [%v]", err,
		)
	}

	// v_i = E(ρ_i * x)
	secretKeyMultiple := s.paillierKey.Mul(
		s.dsaKey.secretKey,
		randomFactorShare,
	)

	// [C_1i, D_1i] = Com([u_i, v_i])
	commitment, decommitmentKey, err := commitment.Generate(
		encryptedRandomFactorShare.C.Bytes(),
		secretKeyMultiple.C.Bytes(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 1 of signing [%v]", err,
		)
	}

	round1Signer := &Round1Signer{
		*s,
		randomFactorShare,
		encryptedRandomFactorShare,
		secretKeyMultiple,
		decommitmentKey,
		paillierRandomness,
	}

	round1Message := &SignRound1Message{
		signerID:               s.ID,
		randomFactorCommitment: commitment,
	}

	return round1Signer, round1Message, nil
}

// Round2Signer represents state of `Signer` after executing the second round
// of signing algorithm.
type Round2Signer struct {
	Signer
}

// SignRound2 executes the second round of T-ECDSA signing as described in
// [GGN 16], section 4.3.
//
// In the second round, encrypted secret factor share `u_i = E(ρ_i)` and
// secret DSA key multiple `v_i = E(ρ_i * x)` is revealed along with
// a decommitment key `D_1i` allowing to check revealed values against the
// commitment published in the first round.
// Moreover, message produced in the second round contains a ZKP allowing to
// verify correctness of revealed values.
func (s *Round1Signer) SignRound2() (*Round2Signer, *SignRound2Message, error) {
	zkp, err := zkp.CommitDsaPaillierSecretKeyFactorRange(
		s.secretKeyMultiple,
		s.dsaKey.secretKey,
		s.encryptedRandomFactorShare,
		s.randomFactorShare,
		s.paillierRandomness,
		s.zkpParameters,
		rand.Reader,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 2 of signing [%v]", err,
		)
	}

	signer := &Round2Signer{s.Signer}

	round2Message := &SignRound2Message{
		signerID:                    s.ID,
		randomFactorShare:           s.encryptedRandomFactorShare,
		secretKeyMultipleShare:      s.secretKeyMultiple,
		randomFactorDecommitmentKey: s.randomFactorDecommitmentKey,
		secretKeyFactorProof:        zkp,
	}

	return signer, round2Message, nil
}

// CombineRound2Messages takes all messages from the first and second signing
// round, validates and combines them together in order to evaluate random
// secret key factor `u` and secret key multiple `v`:
//
// u = u_1 + u_2 + ... + u_n = E(ρ_1) + E(ρ_2) + ... + E(ρ_n)
// v = v_1 + v_2 + ... + v_n = E(ρ_1 * x) + E(ρ_2 * x) + ... + E(ρ_n * x)
//
// This function should be called before the `SignRound3` and the returned
// values should be used as a parameters to `SignRound3`.
func (s *Round2Signer) CombineRound2Messages(
	round1Messages []*SignRound1Message,
	round2Messages []*SignRound2Message,
) (
	randomFactor *paillier.Cypher,
	secretKeyMultiple *paillier.Cypher,
	err error,
) {
	groupSize := s.groupParameters.groupSize

	if len(round1Messages) != groupSize {
		return nil, nil, fmt.Errorf(
			"round 1 messages required from all group members; got %v, expected %v",
			len(round1Messages),
			groupSize,
		)
	}

	if len(round2Messages) != groupSize {
		return nil, nil, fmt.Errorf(
			"round 2 messages required from all group members; got %v, expected %v",
			len(round2Messages),
			groupSize,
		)
	}

	randomFactorShares := make([]*paillier.Cypher, groupSize)
	secretKeyMultipleShares := make([]*paillier.Cypher, groupSize)

	for i, round1Message := range round1Messages {
		foundMatchingRound2Message := false

		for _, round2Message := range round2Messages {
			if round1Message.signerID == round2Message.signerID {
				foundMatchingRound2Message = true

				if round2Message.isValid(
					round1Message.randomFactorCommitment,
					s.dsaKey.secretKey,
					s.zkpParameters,
				) {
					randomFactorShares[i] = round2Message.randomFactorShare
					secretKeyMultipleShares[i] = round2Message.secretKeyMultipleShare
				} else {
					return nil, nil, errors.New("round 2 message rejected")
				}
			}
		}

		if !foundMatchingRound2Message {
			return nil, nil, fmt.Errorf(
				"no matching round 2 message for signer with ID = %v",
				round1Message.signerID,
			)
		}
	}

	randomFactor = s.paillierKey.Add(randomFactorShares...)
	secretKeyMultiple = s.paillierKey.Add(secretKeyMultipleShares...)
	err = nil

	return
}

// Round3Signer represents state of `Signer` after executing the third round
// of signing algorithm.
type Round3Signer struct {
	Signer

	secretKeyRandomFactor *paillier.Cypher // u = E(ρ)
	secretKeyMultiple     *paillier.Cypher // v = E(ρx)

	signatureRandomMultipleSecretShare *big.Int                    // k_i
	signatureRandomMultiplePublicShare *curve.Point                // r_i = g^{k_i}
	signatureRandomMultipleMaskShare   *big.Int                    // c_i
	signatureUnmaskShare               *paillier.Cypher            // w_i = E(k_i * ρ + c_i * q)
	signatureFactorDecommitmentKey     *commitment.DecommitmentKey // Com(r_i, w_i)
	paillierRandomness                 *big.Int
}

// SignRound3 executes the third round of T-ECDSA signing as described in
// [GGN 16], section 4.3.
//
// Before it executes all computations described in [GGN 16], it's required to
// combine messages from the previous two rounds in order to combine
// secret key random factor shares and secret key multiple shares:
// u = u_1 + u_2 + ... + u_n = E(ρ_1) + E(ρ_2) + ... + E(ρ_n)
// v = v_1 + v_2 + ... + v_n = E(ρ_1 * x) + E(ρ_2 * x) + ... + E(ρ_n * x)
//
// To do that, please execute `CombineRound2Messages`` function and pass the
// returned values as an arguments to `SignRound3`.
func (s *Round2Signer) SignRound3(
	randomFactor *paillier.Cypher,
	secretKeyMultiple *paillier.Cypher,
) (
	*Round3Signer, *SignRound3Message, error,
) {
	// k_i = rand(Z_q)
	signatureRandomMultipleSecretShare, err := rand.Int(
		rand.Reader,
		s.groupParameters.curveCardinality(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 3 of signing [%v]", err,
		)
	}

	// r_i = g^{k_i}
	signatureRandomMultiplePublicShare := curve.NewPoint(
		s.groupParameters.curve.ScalarBaseMult(
			signatureRandomMultipleSecretShare.Bytes(),
		),
	)

	// c_i = rand[0, q^6)
	//
	// According to [GGN 16], `c_i` should be randomly chosen from
	// `[-q^6, q^6]`. Since `k_i` is chosen from [0, q), it means that in
	// a lot of cases, signature unmask will be a negative integer, since
	// `D(w) = k_i * rho + c_i * q`.
	// However, keep in mind, that Paillier encryption scheme does not allow for
	// encrypting negative integers by default since they are out of the allowed
	// plaintext space `[0, N)` where `N` is the Paillier modulus.
	// If we pick a negative integer as `c_i`, there is a high probability the
	// signature ZKP and final T-ECDSA signature will fail.
	// That's the reason why we decided to pick a random element from [0, q^6)
	// instead of from `[-q^6, q^6]`.
	signatureRandomMultipleMaskShare, err := rand.Int(
		rand.Reader,
		new(big.Int).Exp(
			s.groupParameters.curveCardinality(),
			big.NewInt(6),
			nil,
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 3 of signing [%v]", err,
		)
	}

	// w_i = E(k_i * ρ + c_i * q)
	paillierRandomness, err := paillier.GetRandomNumberInMultiplicativeGroup(
		s.paillierKey.N, rand.Reader,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 3 of signing [%v]", err,
		)
	}
	maskShareMulCardinality, err := s.paillierKey.EncryptWithR(
		new(big.Int).Mul(
			signatureRandomMultipleMaskShare,
			s.groupParameters.curveCardinality(),
		),
		paillierRandomness,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 3 of signing [%v]", err,
		)
	}
	signatureUnmaskShare := s.paillierKey.Add(
		s.paillierKey.Mul(randomFactor, signatureRandomMultipleSecretShare),
		maskShareMulCardinality,
	)

	// [C_2i, D_2i] = Com(r_i, w_i)
	signatureFactorCommitment, signatureFactorDecommitmentKey, err :=
		commitment.Generate(
			signatureRandomMultiplePublicShare.Bytes(),
			signatureUnmaskShare.C.Bytes(),
		)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 3 of signing [%v]", err,
		)
	}

	signer := &Round3Signer{
		Signer: s.Signer,

		secretKeyRandomFactor: randomFactor,
		secretKeyMultiple:     secretKeyMultiple,

		signatureRandomMultipleSecretShare: signatureRandomMultipleSecretShare,
		signatureRandomMultiplePublicShare: signatureRandomMultiplePublicShare,
		signatureRandomMultipleMaskShare:   signatureRandomMultipleMaskShare,
		signatureUnmaskShare:               signatureUnmaskShare,
		signatureFactorDecommitmentKey:     signatureFactorDecommitmentKey,
		paillierRandomness:                 paillierRandomness,
	}

	round3Message := &SignRound3Message{
		signerID:                  s.ID,
		signatureFactorCommitment: signatureFactorCommitment,
	}

	return signer, round3Message, nil
}

// Round4Signer represents state of `Signer` after executing the fourth round
// of signing algorithm.
type Round4Signer struct {
	Signer

	secretKeyRandomFactor *paillier.Cypher // u = E(ρ)
}

// SignRound4 executes the fourth round of T-ECDSA signing as described in
// [GGN 16], section 4.3.
//
// In the round 4, signer reveals signature random multiple public share
// (`r_i`), signature unmask share (`w_i`) evaluated in the previous round,
// decommitment key allowing to validate commitment to those values
// (published in the previous round) as well as ZKP allowing to check their
// correctness.
func (s *Round3Signer) SignRound4() (*Round4Signer, *SignRound4Message, error) {
	zkp, err := zkp.CommitEcdsaSignatureFactorRangeProof(
		s.signatureRandomMultiplePublicShare,
		s.signatureUnmaskShare,
		s.secretKeyRandomFactor,
		s.signatureRandomMultipleSecretShare,
		s.signatureRandomMultipleMaskShare,
		s.paillierRandomness,
		s.zkpParameters,
		rand.Reader,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 4 of signing [%v]", err,
		)
	}

	signer := &Round4Signer{
		Signer: s.Signer,

		secretKeyRandomFactor: s.secretKeyRandomFactor,
	}

	round4Message := &SignRound4Message{
		signerID: s.ID,

		signatureRandomMultiplePublicShare: s.signatureRandomMultiplePublicShare,
		signatureUnmaskShare:               s.signatureUnmaskShare,
		signatureFactorDecommitmentKey:     s.signatureFactorDecommitmentKey,

		signatureFactorProof: zkp,
	}

	return signer, round4Message, nil
}

// CombineRound4Messages takes all messages from the third and fourth signing
// round, validates and combines them together in order to evaluate public
// signature random multiple `R` and signature unmask parameter `w`:
//
// w = w_1 + w_2 + ... + w_n = E(kρ + cq)
// R = r_1 + r_2 + ... + r_n = g^k
//
// This function should be called before the `SignRound5` and the returned
// values should be used as a parameters to `SignRound5`.
func (s *Round4Signer) CombineRound4Messages(
	round3Messages []*SignRound3Message,
	round4Messages []*SignRound4Message,
) (
	signatureUnmask *paillier.Cypher, // w
	signatureRandomMultiplePublic *curve.Point, // R
	err error,
) {
	groupSize := s.groupParameters.groupSize

	if len(round3Messages) != groupSize {
		return nil, nil, fmt.Errorf(
			"round 3 messages required from all group members; got %v, expected %v",
			len(round3Messages),
			groupSize,
		)
	}

	if len(round4Messages) != groupSize {
		return nil, nil, fmt.Errorf(
			"round 4 messages required from all group members; got %v, expected %v",
			len(round4Messages),
			groupSize,
		)
	}

	signatureUnmaskShares := make([]*paillier.Cypher, groupSize)
	signatureRandomMultiplePublicShares := make([]*curve.Point, groupSize)

	for i, round3Message := range round3Messages {
		foundMatchingRound4Message := false

		for _, round4Message := range round4Messages {
			if round3Message.signerID == round4Message.signerID {
				foundMatchingRound4Message = true

				if round4Message.isValid(
					round3Message.signatureFactorCommitment,
					s.secretKeyRandomFactor,
					s.zkpParameters,
				) {
					signatureRandomMultiplePublicShares[i] = round4Message.signatureRandomMultiplePublicShare
					signatureUnmaskShares[i] = round4Message.signatureUnmaskShare
				} else {
					return nil, nil, errors.New("round 4 message rejected")
				}
			}
		}

		if !foundMatchingRound4Message {
			return nil, nil, fmt.Errorf(
				"no matching round 4 message for signer with ID = %v",
				round3Message.signerID,
			)
		}
	}

	// w = w_1 + w_2 + ... + w_n
	signatureUnmask = s.paillierKey.Add(signatureUnmaskShares...)

	// R = r_i + r_2 + ... + r_n
	signatureRandomMultiplePublicX := big.NewInt(0)
	signatureRandomMultiplePublicY := big.NewInt(0)
	for _, share := range signatureRandomMultiplePublicShares {
		signatureRandomMultiplePublicX.Add(
			signatureRandomMultiplePublicX, share.X,
		)
		signatureRandomMultiplePublicY.Add(
			signatureRandomMultiplePublicY, share.Y,
		)
	}
	signatureRandomMultiplePublic = &curve.Point{
		X: signatureRandomMultiplePublicX,
		Y: signatureRandomMultiplePublicY,
	}

	err = nil

	return
}

// Round5Signer represents state of `Signer` after executing the fifth round
// of signing algorithm.
type Round5Signer struct {
	Signer

	signatureUnmask               *paillier.Cypher // w
	signatureRandomMultiplePublic *curve.Point     // R
}

// SignRound5 executes the fifth round of signing
func (s *Round4Signer) SignRound5(
	signatureUnmask *paillier.Cypher, // w
	signatureRandomMultiplePublic *curve.Point, // R
) (
	*Round5Signer, error,
) {
	signer := &Round5Signer{
		Signer: s.Signer,

		signatureUnmask:               signatureUnmask,
		signatureRandomMultiplePublic: signatureRandomMultiplePublic,
	}

	return signer, nil
}
