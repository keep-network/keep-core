package tecdsa

type signerGroup struct {
	// IDs of all signers in active signer's group, including the signer itself.
	signerIDs []string
}

// RegisterSignerID adds a signer to the list of signers the local signer
// knows about.
func (sg *signerGroup) RegisterSignerID(ID string) {
	// TODO Validate if signer ID is unique, add trim
	sg.signerIDs = append(sg.signerIDs, ID)
}

// RemoveSignerID removes a signer from the list of signers the local signer
// knows about.
func (sg *signerGroup) RemoveSignerID(ID string) {
	for i := 0; i < len(sg.signerIDs); i++ {
		if sg.signerIDs[i] == ID {
			sg.signerIDs = append(sg.signerIDs[:i], sg.signerIDs[i+1:]...)
		}
	}
}

// IsActiveSigner checks if a signer with given ID is one of the signers the local
// signer knows about.
func (sg *signerGroup) IsActiveSigner(ID string) bool {
	for i := 0; i < len(sg.signerIDs); i++ {
		if sg.signerIDs[i] == ID {
			return true
		}
	}
	return false
}
