package tecdsa

func setupGroup(group []*LocalSigner) error {
	var err error

	// Initialize master public key for multi-trapdoor commitment scheme.
	// Each signer generates a master public key which is a point in
	// G2 abstract cyclic group of bn256 curve. The key is broadcasted in
	// CommitmentMasterPublicKeyMessage.
	commitmentMasterPublicKeyMessages := make(
		[]*CommitmentMasterPublicKeyMessage, len(group),
	)
	for i, signer := range group {
		commitmentMasterPublicKeyMessages[i], err =
			signer.GenerateCommitmentMasterPublicKey()
		if err != nil {
			return err
		}
	}

	for _, signer := range group {
		signer.ReceiveCommitmentMasterPublicKeys(
			commitmentMasterPublicKeyMessages,
		)
	}

	return nil
}
