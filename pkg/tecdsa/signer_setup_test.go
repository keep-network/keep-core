package tecdsa

import (
	"errors"
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
		t.Fatalf("peer signer IDs doesn't match expected\nactual: %s\nexpected: %s",
			signer.peerSignerIDs(),
			expectedPeerIDs)
	}
}

func TestGenerateCommitmentMasterPublicKey(t *testing.T) {
	signer := &LocalSigner{}

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

func TestReceiveCommitmentMasterPublicKeys(t *testing.T) {
	var tests = map[string]struct {
		updateSignerKeys func(signerKeys map[string]string)
		expectedError    error
	}{
		"positive": {
			expectedError: nil,
		},
		"negative validation - not enough messages": {
			updateSignerKeys: func(signerKeys map[string]string) {
				delete(signerKeys, "1003")
			},
			expectedError: errors.New("master public key messages required from all group peer members; got 2, expected 3"),
		},
		"negative validation - too many messages": {
			updateSignerKeys: func(signerKeys map[string]string) {
				signerKeys["1005"] = "key1005"
			},
			expectedError: errors.New("master public key messages required from all group peer members; got 4, expected 3"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			// Given
			signerID := "1001"
			signerIDs := []string{"1001", "1002", "1003", "1004"}

			signer := &LocalSigner{
				signerCore: signerCore{
					ID: signerID,
					signerGroup: &signerGroup{
						signerIDs: signerIDs,
					},
					protocolParameters: make(map[string]*protocolParameters, 0),
				},
			}

			signerKeys := make(map[string]string)
			signerKeys["1002"] = "key1002"
			signerKeys["1003"] = "key1003"
			signerKeys["1004"] = "key1004"

			if test.updateSignerKeys != nil {
				test.updateSignerKeys(signerKeys)
			}

			messages := make([]*CommitmentMasterPublicKeyMessage, 0)
			for k, v := range signerKeys {
				messages = append(messages, &CommitmentMasterPublicKeyMessage{
					senderID:        k,
					masterPublicKey: []byte(v),
				})
			}

			// When
			err := signer.ReceiveCommitmentMasterPublicKeys(messages)

			// Then
			if !reflect.DeepEqual(test.expectedError, err) {
				t.Fatalf(
					"unexpected error\nexpected: %v\nactual: %v",
					test.expectedError,
					err,
				)
			}
			if test.expectedError == nil {
				for k, v := range signerKeys {
					if reflect.DeepEqual(
						signer.protocolParameters[k].commitmentMasterPublicKey.Marshal(),
						[]byte(v),
					) {
						t.Fatal("peer's commitmentMasterPublicKey doesn't match received in a message")
					}
				}
			}
		})
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
		messages := peerSignersCommitmentMasterPublicKeys(
			commitmentMasterPublicKeyMessages,
			signer.ID)
		signer.ReceiveCommitmentMasterPublicKeys(messages)
	}

	return nil
}

func peerSignersCommitmentMasterPublicKeys(
	messages []*CommitmentMasterPublicKeyMessage,
	signerID string,
) []*CommitmentMasterPublicKeyMessage {
	filtered := make([]*CommitmentMasterPublicKeyMessage, 0)
	for _, message := range messages {
		if message.senderID != signerID {
			filtered = append(filtered, message)
		}
	}
	return filtered
}
