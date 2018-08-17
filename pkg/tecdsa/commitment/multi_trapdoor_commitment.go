// Package commitment generates and validates trapdoor commitments using
// bn256 elliptic curve pairings as described in
// docs/cryptography/trapdoor-commitments.adoc
package commitment

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// DecommitmentKey allows to open a commitment and verify if the value is what
// we have really committed to.
type DecommitmentKey struct {
	r                   *big.Int
	commitmentSignature []byte
}

// TrapdoorCommitment is produced for each message we have committed to.
// It is usually revealed to the receiver immediately after it has been produced.
// Commitment lets to verify if the message revealed later by the sending party
// is really what that party has committed to.
// However, the Commitment itself is not enough for a verification.
// In order to perform verification, the interested party must receive
// the DecommitmentKey too.
//
// Usually the process happens in two phases:
// first, Commitment is evaluated and sent to receiver and then, after some time,
// secret value along with a DecommitmentKey is revealed and the receiver can
// check the secret value against the Commitment received earlier.
type TrapdoorCommitment struct {
	// Public key for a specific trapdoor commitment.
	pubKey *big.Int
	// Master trapdoor public key for the commitment family.
	h *bn256.G2
	// Calculated trapdoor commitment.
	commitment *bn256.G2

	verificationKey ecdsa.PublicKey
}

// Generate evaluates a commitment and decommitment key for the secret
// messages provided as an argument.
func Generate(secrets ...[]byte) (*TrapdoorCommitment, *DecommitmentKey, error) {
	secret := combineSecrets(secrets...)

	signatureSecretKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not generate multi-trapdoor commitment [%v]", err,
		)
	}
	signaturePublicKey := signatureSecretKey.PublicKey

	// pk = H(vk)
	commitmentPublicKey := hashPublicSignatureKey(signaturePublicKey)

	// Generate a decommitment key.
	r, _, err := bn256.RandomG1(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Generate random point.
	_, h, err := bn256.RandomG2(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	hash := sha256Sum(secret)
	digest := new(big.Int).Mod(hash, bn256.Order)

	// he = h + g * pubKey
	he := new(bn256.G2).Add(h, new(bn256.G2).ScalarBaseMult(commitmentPublicKey))

	// commitment = g * digest + he * r
	commitment := new(bn256.G2).Add(
		new(bn256.G2).ScalarBaseMult(digest),
		new(bn256.G2).ScalarMult(he, r),
	)

	commitmentSignature, err := signatureSecretKey.Sign(
		rand.Reader, commitment.Marshal(), nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"could not generate multi-trapdoor commitment [%v]", err,
		)
	}

	return &TrapdoorCommitment{
			pubKey:          commitmentPublicKey,
			h:               h,
			commitment:      commitment,
			verificationKey: signaturePublicKey,
		},
		&DecommitmentKey{
			r:                   r,
			commitmentSignature: commitmentSignature,
		},
		nil
}

// Verify checks received commitment against the revealed secret messages.
func (tc *TrapdoorCommitment) Verify(
	decommitmentKey *DecommitmentKey,
	secrets ...[]byte,
) bool {
	secret := combineSecrets(secrets...)

	hash := sha256Sum(secret)
	digest := new(big.Int).Mod(hash, bn256.Order)

	// a = g * r
	a := new(bn256.G1).ScalarBaseMult(decommitmentKey.r)

	// b = h + g * pubKey
	b := new(bn256.G2).Add(tc.h, new(bn256.G2).ScalarBaseMult(tc.pubKey))

	// c = commitment - g * digest
	c := new(bn256.G2).Add(
		tc.commitment,
		new(bn256.G2).Neg(new(bn256.G2).ScalarBaseMult(digest)),
	)

	// Get base point `g`
	g := new(bn256.G1).ScalarBaseMult(big.NewInt(1))

	if bn256.Pair(a, b).String() != bn256.Pair(g, c).String() {
		return false
	}
	return true
}

func hashPublicSignatureKey(publicSignatureKey ecdsa.PublicKey) *big.Int {
	return new(big.Int).Mod(
		sha256Sum(combineSecrets(
			publicSignatureKey.X.Bytes(),
			publicSignatureKey.Y.Bytes(),
		)),
		publicSignatureKey.Params().N,
	)
}

// sha256Sum calculates sha256 hash for the passed `secret`
// and converts it to `big.Int`.
func sha256Sum(secret []byte) *big.Int {
	hash := sha256.Sum256(secret)

	return new(big.Int).SetBytes(hash[:])
}

func combineSecrets(secrets ...[]byte) []byte {
	var combined []byte
	for _, secret := range secrets {
		combined = append(combined, secret...)
	}
	return combined
}
