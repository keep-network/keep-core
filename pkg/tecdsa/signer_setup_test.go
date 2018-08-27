package tecdsa

func setupGroup(group []*LocalSigner) error {
	var err error

	// Initialize master public key for multi-trapdoor commitment scheme.
	// Each signer generates a master public key share which is a point in
	// G2 abstract cyclic group of bn256 curve. The share is broadcasted in
	// MasterPublicKeyShareMessage.
	// The shares are combined by adding the points which results in a point
	// which is a master public key.
	masterPublicKeyShareMessages := make([]*MasterPublicKeyShareMessage, len(group))
	for i, signer := range group {
		masterPublicKeyShareMessages[i], err = signer.GenerateMasterPublicKeyShare()
		if err != nil {
			return err
		}
	}

	masterPublicKey, err := group[0].CombineMasterPublicKeyShares(masterPublicKeyShareMessages)
	if err != nil {
		return err
	}

	for _, signer := range group {
		signer.commitmentMasterPublicKey = masterPublicKey
	}

	return nil
}
