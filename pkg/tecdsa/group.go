package tecdsa

import "fmt"

type signerGroup struct {
	// IDs of all signers in active signer's group, including the signer itself.
	signerIDs []string
}

// AddSignerID adds a signer ID to the list of signers in the group.
func (sg *signerGroup) AddSignerID(ID string) {
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

// Size return number of signers in the signing group.
func (sg *signerGroup) Size() int {
	return len(sg.signerIDs)
}

// IsSignerGroupComplete checks if a number of signers in a group matches initial
// signers group size.
func (sc *signerCore) IsSignerGroupComplete() (bool, error) {
	if sc.signerGroup.Size() != sc.groupParameters.GroupSize {
		return false, fmt.Errorf("current signers group size %v doesn't match expected size %v",
			sc.signerGroup.Size(),
			sc.groupParameters.GroupSize,
		)
	}
	return true, nil
}
