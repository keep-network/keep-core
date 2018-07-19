package tecdsa

import "testing"

// In the second round, signer reveals values for which he committed to in the
// first round. We check if the commitment verification succeeds and if the
// ZKP produced in the second round evaluates to true.
func TestSignRound1And2(t *testing.T) {
	signers, err := initializeNewSignerGroup()
	if err != nil {
		t.Fatal(err)
	}

	round1Signer, round1Message, err := signers[0].SignRound1()
	if err != nil {
		t.Fatal(err)
	}

	round2Signer, round2Message, err := round1Signer.SignRound2()
	if err != nil {
		t.Fatal(err)
	}

	if !round1Message.randomFactorCommitment.Verify(
		round2Message.randomFactorDecommitmentKey,
		round2Message.randomFactorShare.C.Bytes(),
		round2Message.secretKeyMultiple.C.Bytes(),
	) {
		t.Fatal("Round2Message commitment verification failed")
	}

	if !round2Message.secretKeyFactorProof.Verify(
		round2Message.secretKeyMultiple,
		round2Signer.dsaKey.secretKey,
		round2Message.randomFactorShare,
		round2Signer.zkpParameters,
	) {
		t.Fatal("Round2Message ZKP verification failed")
	}
}

// Crates and initializes a new group of `Signer`s with T-ECDSA key set and
// ready for signing.
func initializeNewSignerGroup() ([]*Signer, error) {
	localGroup, key, err := initializeNewLocalGroupWithFullKey()
	if err != nil {
		return nil, err
	}

	signers := make([]*Signer, len(localGroup))
	for i, localSigner := range localGroup {
		signers[i] = &Signer{
			dsaKey:     key,
			signerCore: localSigner.signerCore,
		}
	}

	return signers, nil
}
