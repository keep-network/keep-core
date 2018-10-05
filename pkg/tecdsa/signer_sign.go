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
	"crypto/rand"
	"errors"
	"fmt"

	"math/big"

	"github.com/keep-network/keep-core/pkg/tecdsa/commitment"
	"github.com/keep-network/keep-core/pkg/tecdsa/curve"
	"github.com/keep-network/keep-core/pkg/tecdsa/zkp"
	"github.com/keep-network/paillier"
)

// Round1Signer represents state of signer after executing the first round
// of signing algorithm.
type Round1Signer struct {
	Signer

	// Intermediate values stored between the first and second round of signing.
	secretKeyFactorShare                 *big.Int                               // ρ_i
	encryptedSecretKeyFactorShare        *paillier.Cypher                       // u_i = E(ρ_i)
	secretKeyMultipleShare               *paillier.Cypher                       // v_i = E(ρ_i * x)
	secretKeyFactorShareDecommitmentKeys map[string]*commitment.DecommitmentKey // D_1i
	paillierRandomness                   *big.Int
}

// SignRound1 executes the first round of T-ECDSA signing as described in
// [GGN 16], section 4.3.
//
// In the first round, each signer generates a secret key factor share `ρ_i`,
// encodes it with Paillier key `u_i = E(ρ_i)`, multiplies it with secret ECDSA
// key `v_i = E(ρ_i * x)` and publishes commitments for both those values
// `Com(u_i, v_i)`. Individual commitment is calculated for each peer signer.
func (s *Signer) SignRound1() (*Round1Signer, []*SignRound1Message, error) {
	// Choosing random ρ_i from Z_q
	secretKeyFactorShare, err := rand.Int(
		rand.Reader,
		s.publicParameters.curveCardinality(),
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
	encryptedSecretKeyFactorShare, err := s.paillierKey.EncryptWithR(
		secretKeyFactorShare, paillierRandomness,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 1 of signing [%v]", err,
		)
	}

	// v_i = E(ρ_i * x)
	secretKeyMultiple := s.paillierKey.Mul(
		s.ecdsaKey.secretKey,
		secretKeyFactorShare,
	)

	commitments := make(map[string]*commitment.MultiTrapdoorCommitment)
	decommitmentKeys := make(map[string]*commitment.DecommitmentKey)
	round1Messages := make([]*SignRound1Message, s.signerGroup.PeerSignerCount())

	for i, peerSignerID := range s.peerSignerIDs() {
		// [C_1i, D_1i] = Com([u_i, v_i])
		commitments[peerSignerID], decommitmentKeys[peerSignerID], err = commitment.Generate(
			s.protocolParameters[peerSignerID].commitmentMasterPublicKey,
			encryptedSecretKeyFactorShare.C.Bytes(),
			secretKeyMultiple.C.Bytes(),
		)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"could not execute round 1 of signing [%v]", err,
			)
		}

		round1Messages[i] = &SignRound1Message{
			senderID:   s.ID,
			receiverID: peerSignerID,

			secretKeyFactorShareCommitment: commitments[peerSignerID],
		}
	}

	round1Signer := &Round1Signer{
		Signer:                               *s,
		secretKeyFactorShare:                 secretKeyFactorShare,
		encryptedSecretKeyFactorShare:        encryptedSecretKeyFactorShare,
		secretKeyMultipleShare:               secretKeyMultiple,
		secretKeyFactorShareDecommitmentKeys: decommitmentKeys,
		paillierRandomness:                   paillierRandomness,
	}

	return round1Signer, round1Messages, nil
}

// Round2Signer represents state of signer after executing the second round
// of signing algorithm.
type Round2Signer struct {
	*Round1Signer
}

// SignRound2 executes the second round of T-ECDSA signing as described in
// [GGN 16], section 4.3.
//
// In the second round, encrypted secret key factor share `u_i = E(ρ_i)` and
// secret DSA key multiple `v_i = E(ρ_i * x)` are revealed along with
// a decommitment key `D_1i` allowing to check revealed values against the
// commitment published in the first round.
// Moreover, message produced in the second round contains a ZKP allowing to
// verify correctness of the revealed values.
func (s *Round1Signer) SignRound2() (*Round2Signer, []*SignRound2Message, error) {
	zkp, err := zkp.CommitEcdsaPaillierSecretKeyFactorRange(
		s.secretKeyMultipleShare,
		s.ecdsaKey.secretKey,
		s.encryptedSecretKeyFactorShare,
		s.secretKeyFactorShare,
		s.paillierRandomness,
		s.zkpParameters,
		rand.Reader,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 2 of signing [%v]", err,
		)
	}

	signer := &Round2Signer{s}

	round2Messages := make([]*SignRound2Message, s.signerGroup.PeerSignerCount())

	for i, peerSignerID := range s.peerSignerIDs() {
		round2Messages[i] = &SignRound2Message{
			senderID:   s.ID,
			receiverID: peerSignerID,

			secretKeyFactorShare:                s.encryptedSecretKeyFactorShare,
			secretKeyMultipleShare:              s.secretKeyMultipleShare,
			secretKeyFactorShareDecommitmentKey: s.secretKeyFactorShareDecommitmentKeys[peerSignerID],
			secretKeyFactorProof:                zkp,
		}
	}

	return signer, round2Messages, nil
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
	secretKeyFactor *paillier.Cypher,
	secretKeyMultiple *paillier.Cypher,
	err error,
) {
	peerSignerCount := s.signerGroup.PeerSignerCount()

	if len(round1Messages) != peerSignerCount {
		return nil, nil, fmt.Errorf(
			"round 1 messages required from all group peer members; got %v, expected %v",
			len(round1Messages),
			peerSignerCount,
		)
	}

	if len(round2Messages) != peerSignerCount {
		return nil, nil, fmt.Errorf(
			"round 2 messages required from all group peer members; got %v, expected %v",
			len(round2Messages),
			peerSignerCount,
		)
	}

	secretKeyFactorShares := make([]*paillier.Cypher, peerSignerCount)
	secretKeyMultipleShares := make([]*paillier.Cypher, peerSignerCount)

	for i, round1Message := range round1Messages {
		foundMatchingRound2Message := false

		for _, round2Message := range round2Messages {
			if round1Message.senderID == round2Message.senderID {
				foundMatchingRound2Message = true

				if round2Message.isValid(
					s.selfProtocolParameters().commitmentMasterPublicKey,
					round1Message.secretKeyFactorShareCommitment,
					s.ecdsaKey.secretKey,
					s.zkpParameters,
				) {
					secretKeyFactorShares[i] = round2Message.secretKeyFactorShare
					secretKeyMultipleShares[i] = round2Message.secretKeyMultipleShare
				} else {
					return nil, nil, errors.New("round 2 message rejected")
				}
			}
		}

		if !foundMatchingRound2Message {
			return nil, nil, fmt.Errorf(
				"no matching round 2 message for signer with ID = %v",
				round1Message.senderID,
			)
		}
	}

	// Add signer's own shares
	secretKeyFactorShares = append(secretKeyFactorShares, s.encryptedSecretKeyFactorShare)
	secretKeyMultipleShares = append(secretKeyMultipleShares, s.secretKeyMultipleShare)

	secretKeyFactor = s.paillierKey.Add(secretKeyFactorShares...)
	secretKeyMultiple = s.paillierKey.Add(secretKeyMultipleShares...)
	err = nil

	return
}

// Round3Signer represents state of signer after executing the third round
// of signing algorithm.
type Round3Signer struct {
	Signer

	secretKeyFactor   *paillier.Cypher // u = E(ρ)
	secretKeyMultiple *paillier.Cypher // v = E(ρx)

	// Intermediate values stored between the third and fourth round of signing
	signatureFactorSecretShare           *big.Int                               // k_i
	signatureFactorPublicShare           *curve.Point                           // r_i = g^{k_i}
	signatureFactorMaskShare             *big.Int                               // c_i
	signatureUnmaskShare                 *paillier.Cypher                       // w_i = E(k_i * ρ + c_i * q)
	signatureFactorShareDecommitmentKeys map[string]*commitment.DecommitmentKey // Com(r_i, w_i)
	paillierRandomness                   *big.Int
}

// SignRound3 executes the third round of T-ECDSA signing as described in
// [GGN 16], section 4.3.
//
// Before it executes all computations described in [GGN 16], it's required to
// combine messages from the previous two rounds in order to combine
// secret key factor shares and secret key multiple shares:
// u = u_1 + u_2 + ... + u_n = E(ρ_1) + E(ρ_2) + ... + E(ρ_n)
// v = v_1 + v_2 + ... + v_n = E(ρ_1 * x) + E(ρ_2 * x) + ... + E(ρ_n * x)
//
// To do that, please execute `CombineRound2Messages`` function and pass the
// returned values as an arguments to `SignRound3`.
func (s *Round2Signer) SignRound3(
	secretKeyFactor *paillier.Cypher, // u = E(ρ)
	secretKeyMultiple *paillier.Cypher, // v = E(ρx)
) (
	*Round3Signer, []*SignRound3Message, error,
) {
	// k_i = rand(Z_q)
	signatureFactorSecretShare, err := rand.Int(
		rand.Reader,
		s.publicParameters.curveCardinality(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 3 of signing [%v]", err,
		)
	}

	// r_i = g^{k_i}
	signatureFactorPublicShare := curve.NewPoint(
		s.publicParameters.Curve.ScalarBaseMult(
			signatureFactorSecretShare.Bytes(),
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
	signatureFactorMaskShare, err := rand.Int(
		rand.Reader,
		new(big.Int).Exp(
			s.publicParameters.curveCardinality(),
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
			signatureFactorMaskShare,
			s.publicParameters.curveCardinality(),
		),
		paillierRandomness,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 3 of signing [%v]", err,
		)
	}
	signatureUnmaskShare := s.paillierKey.Add(
		s.paillierKey.Mul(secretKeyFactor, signatureFactorSecretShare),
		maskShareMulCardinality,
	)

	commitments := make(map[string]*commitment.MultiTrapdoorCommitment)
	decommitmentKeys := make(map[string]*commitment.DecommitmentKey)
	round3Messages := make([]*SignRound3Message, s.signerGroup.PeerSignerCount())

	for i, peerSignerID := range s.peerSignerIDs() {
		// [C_2i, D_2i] = Com(r_i, w_i)
		commitments[peerSignerID], decommitmentKeys[peerSignerID], err =
			commitment.Generate(
				s.protocolParameters[peerSignerID].commitmentMasterPublicKey,
				signatureFactorPublicShare.Bytes(),
				signatureUnmaskShare.C.Bytes(),
			)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"could not execute round 3 of signing [%v]", err,
			)
		}

		round3Messages[i] = &SignRound3Message{
			senderID:   s.ID,
			receiverID: peerSignerID,

			signatureFactorShareCommitment: commitments[peerSignerID],
		}
	}

	signer := &Round3Signer{
		Signer: s.Signer,

		secretKeyFactor:   secretKeyFactor,
		secretKeyMultiple: secretKeyMultiple,

		signatureFactorSecretShare:           signatureFactorSecretShare,
		signatureFactorPublicShare:           signatureFactorPublicShare,
		signatureFactorMaskShare:             signatureFactorMaskShare,
		signatureUnmaskShare:                 signatureUnmaskShare,
		signatureFactorShareDecommitmentKeys: decommitmentKeys,
		paillierRandomness:                   paillierRandomness,
	}

	return signer, round3Messages, nil
}

// Round4Signer represents state of signer after executing the fourth round
// of signing algorithm.
type Round4Signer struct {
	*Round3Signer
}

// SignRound4 executes the fourth round of T-ECDSA signing as described in
// [GGN 16], section 4.3.
//
// In the round 4, signer reveals signature factor public share
// (`r_i`), signature unmask share (`w_i`) evaluated in the previous round,
// decommitment key allowing to validate commitment to those values
// (published in the previous round) as well as ZKP allowing to check their
// correctness.
func (s *Round3Signer) SignRound4() (*Round4Signer, []*SignRound4Message, error) {
	zkp, err := zkp.CommitEcdsaSignatureFactorRangeProof(
		s.signatureFactorPublicShare,
		s.signatureUnmaskShare,
		s.secretKeyFactor,
		s.signatureFactorSecretShare,
		s.signatureFactorMaskShare,
		s.paillierRandomness,
		s.zkpParameters,
		rand.Reader,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not execute round 4 of signing [%v]", err,
		)
	}

	signer := &Round4Signer{s}

	round4Messages := make([]*SignRound4Message, s.signerGroup.PeerSignerCount())

	for i, peerSignerID := range s.peerSignerIDs() {
		round4Messages[i] = &SignRound4Message{
			senderID:   s.ID,
			receiverID: peerSignerID,

			signatureFactorPublicShare:          s.signatureFactorPublicShare,
			signatureUnmaskShare:                s.signatureUnmaskShare,
			signatureFactorShareDecommitmentKey: s.signatureFactorShareDecommitmentKeys[peerSignerID],

			signatureFactorProof: zkp,
		}
	}

	return signer, round4Messages, nil
}

// CombineRound4Messages takes all messages from the third and fourth signing
// round, validates and combines them together in order to evaluate public
// signature factor `R` and signature unmask parameter `w`:
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
	signatureFactorPublic *curve.Point, // R
	err error,
) {
	peerSignerCount := s.signerGroup.PeerSignerCount()

	if len(round3Messages) != peerSignerCount {
		return nil, nil, fmt.Errorf(
			"round 3 messages required from all group peer members; got %v, expected %v",
			len(round3Messages),
			peerSignerCount,
		)
	}

	if len(round4Messages) != peerSignerCount {
		return nil, nil, fmt.Errorf(
			"round 4 messages required from all group peer members; got %v, expected %v",
			len(round4Messages),
			peerSignerCount,
		)
	}

	signatureUnmaskShares := make([]*paillier.Cypher, peerSignerCount)
	signatureFactorPublicShares := make([]*curve.Point, peerSignerCount)

	for i, round3Message := range round3Messages {
		foundMatchingRound4Message := false

		for _, round4Message := range round4Messages {
			if round3Message.senderID == round4Message.senderID {
				foundMatchingRound4Message = true

				if round4Message.isValid(
					s.selfProtocolParameters().commitmentMasterPublicKey,
					round3Message.signatureFactorShareCommitment,
					s.secretKeyFactor,
					s.zkpParameters,
				) {
					signatureUnmaskShares[i] = round4Message.signatureUnmaskShare
					signatureFactorPublicShares[i] = round4Message.signatureFactorPublicShare
				} else {
					return nil, nil, errors.New("round 4 message rejected")
				}
			}
		}

		if !foundMatchingRound4Message {
			return nil, nil, fmt.Errorf(
				"no matching round 4 message for signer with ID = %v",
				round3Message.senderID,
			)
		}
	}

	// Add signer's own shares
	signatureUnmaskShares = append(signatureUnmaskShares, s.signatureUnmaskShare)
	signatureFactorPublicShares = append(signatureFactorPublicShares, s.signatureFactorPublicShare)

	// w = w_1 + w_2 + ... + w_n
	signatureUnmask = s.paillierKey.Add(signatureUnmaskShares...)

	// R = r_i + r_2 + ... + r_n
	signatureFactorPublic = signatureFactorPublicShares[0]
	for _, share := range signatureFactorPublicShares[1:] {
		signatureFactorPublic = curve.NewPoint(
			s.publicParameters.Curve.Add(
				signatureFactorPublic.X,
				signatureFactorPublic.Y,
				share.X,
				share.Y,
			))
	}

	err = nil

	return
}

// Round5Signer represents state of `Signer` after executing the fifth round
// of signing algorithm.
type Round5Signer struct {
	Signer

	secretKeyFactor           *paillier.Cypher // u = E(ρ)
	secretKeyMultiple         *paillier.Cypher // v = E(ρx)
	signatureUnmask           *paillier.Cypher // w
	signatureFactorPublic     *curve.Point     // R
	signatureFactorPublicHash *big.Int         // r = H'(R)
}

// SignRound5 executes the fifth round of signing. In the fifth round, signers
// jointly decrypt signature unmask `w` as well as compute hash of the signature
// factor public parameter. Both values will be used in round six when
// evaluating the final signature.
func (s *Round4Signer) SignRound5(
	signatureUnmask *paillier.Cypher, // w
	signatureFactorPublic *curve.Point, // R
) (
	*Round5Signer, *SignRound5Message, error,
) {

	// TDec(w) share
	signatureUnmaskPartialDecryption := s.paillierKey.Decrypt(signatureUnmask.C)

	// r = H'(R)
	//
	// According to [GGN 16], H' is a hash function defined from `G` to `Z_q`.
	// It does not have to be a cryptographic hash function, so we use the
	// simplest possible form here.
	signatureFactorPublicHash := new(big.Int).Mod(
		signatureFactorPublic.X,
		s.publicParameters.curveCardinality(),
	)

	message := &SignRound5Message{
		senderID: s.ID,

		signatureUnmaskPartialDecryption: signatureUnmaskPartialDecryption,
	}

	signer := &Round5Signer{
		Signer: s.Signer,

		secretKeyFactor:           s.secretKeyFactor,
		secretKeyMultiple:         s.secretKeyMultiple,
		signatureFactorPublic:     signatureFactorPublic,
		signatureFactorPublicHash: signatureFactorPublicHash,
	}

	return signer, message, nil
}

// CombineRound5Messages combines together all `SignRound5Message`s produced by
// signers. Each message contains a partial decryption for signature unmask
// parameter `w`. Function combines them together and returns a final decrypted
// value of signature unmask.
func (s *Round5Signer) CombineRound5Messages(
	round5Messages []*SignRound5Message,
) (
	signatureUnmask *big.Int, // TDec(w)
	err error,
) {
	groupSize := s.signerGroup.InitialGroupSize

	if len(round5Messages) != groupSize {
		return nil, fmt.Errorf(
			"round 5 messages required from all group members; got %v, expected %v",
			len(round5Messages),
			groupSize,
		)
	}

	partialDecryptions := make([]*paillier.PartialDecryption, groupSize)
	for i, round5Message := range round5Messages {
		partialDecryptions[i] = round5Message.signatureUnmaskPartialDecryption
	}

	signatureUnmask, err = s.paillierKey.CombinePartialDecryptions(
		partialDecryptions,
	)
	if err != nil {
		err = fmt.Errorf(
			"could not combine signature unmask partial decryptions [%v]",
			err,
		)
	}

	return
}

// SignRound6 executes the sixth round of signing. In the sixth round, all
// parameters signers evaluates so far are combined together in order to produce
// a final signature. The final signature is in a Paillier-encrypted form, so
// a threshold decode action must be performed.
func (s *Round5Signer) SignRound6(
	signatureUnmask *big.Int, // TDec(w)
	messageHash []byte, // m
) (*SignRound6Message, error) {
	if len(messageHash) != 32 {
		return nil, fmt.Errorf(
			"message hash is required to be exactly 32 bytes and it's %d bytes",
			len(messageHash),
		)
	}

	signatureCypher := s.paillierKey.Mul(
		s.paillierKey.Add(
			s.paillierKey.Mul(
				s.secretKeyFactor,
				new(big.Int).SetBytes(messageHash[:]),
			),
			s.paillierKey.Mul(
				s.secretKeyMultiple,
				s.signatureFactorPublicHash,
			),
		),
		new(big.Int).ModInverse(
			signatureUnmask,
			s.publicParameters.curveCardinality(),
		),
	)

	return &SignRound6Message{
		signaturePartialDecryption: s.paillierKey.Decrypt(signatureCypher.C),
	}, nil
}

// Signature represents a final T-ECDSA signature
type Signature struct {
	R *big.Int
	S *big.Int
}

// CombineRound6Messages combines together all partial decryptions of signature
// generated in the sixth round of signing. It outputs a final T-ECDSA signature
// in an unencrypted form.
func (s *Round5Signer) CombineRound6Messages(
	round6Messages []*SignRound6Message,
) (*Signature, error) {
	groupSize := s.signerGroup.InitialGroupSize

	if len(round6Messages) != groupSize {
		return nil, fmt.Errorf(
			"round 6 messages required from all group members; got %v, expected %v",
			len(round6Messages),
			groupSize,
		)
	}

	partialDecryptions := make([]*paillier.PartialDecryption, groupSize)
	for i, round6Message := range round6Messages {
		partialDecryptions[i] = round6Message.signaturePartialDecryption
	}

	sign, err := s.paillierKey.CombinePartialDecryptions(
		partialDecryptions,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not combine signature partial decryptions [%v]",
			err,
		)
	}

	sign = new(big.Int).Mod(sign, s.publicParameters.curveCardinality())

	// Inherent ECDSA signature malleability
	// BTC and ETH require that the S value inside ECDSA signatures is at most
	// the curve order divided by 2 (essentially restricting this value to its
	// lower half range).
	if sign.Cmp(s.publicParameters.halfCurveCardinality()) == 1 {
		sign = new(big.Int).Sub(s.publicParameters.curveCardinality(), sign)
	}

	return &Signature{
		R: s.signatureFactorPublicHash,
		S: sign,
	}, nil
}
