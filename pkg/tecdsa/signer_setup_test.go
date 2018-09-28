package tecdsa

import (
	"reflect"
	"testing"
)

func TestPeerSignerIDs(t *testing.T) {
	signerID := "1003"
	signerIDs := []string{"1001", "1002", "1003", "1004"}
	expectedPeerIDs := []string{"1001", "1002", "1004"}

	signer := &LocalSigner{
		signerCore: signerCore{
			ID: signerID,
			signerGroup: &signerGroup{
				signerIDs: signerIDs,
			},
		},
	}

	if !reflect.DeepEqual(signer.peerSignerIDs(), expectedPeerIDs) {
		t.Fatalf("peer signer IDs doesn't match expected\nactual: %s\nexpected: %s", signer.peerSignerIDs(), expectedPeerIDs)
	}
}

func TestGenerateCommitmentMasterPublicKey(t *testing.T) {
	signers, _, err := generateNewLocalGroup()
	if err != nil {
		t.Fatal(err)
	}
	signer := signers[0]

	message, err := signer.GenerateCommitmentMasterPublicKey()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(
		signer.selfProtocolParameters().commitmentMasterPublicKey.Marshal(),
		message.masterPublicKey,
	) {
		t.Fatal("signer's commitmentMasterPublicKey doesn't match the one from the message")
	}

}

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
		messages := commitmentMasterPublicKeyMessagesNotFromSigner(commitmentMasterPublicKeyMessages, signer.ID)
		signer.ReceiveCommitmentMasterPublicKeys(messages)
	}

	return nil
}

func commitmentMasterPublicKeyMessagesNotFromSigner(
	messages []*CommitmentMasterPublicKeyMessage,
	signerID string,
) []*CommitmentMasterPublicKeyMessage {
	filtered := make([]*CommitmentMasterPublicKeyMessage, 0)
	for _, message := range messages {
		if message.signerID != signerID {
			filtered = append(filtered, message)
		}
	}
	return filtered
}
