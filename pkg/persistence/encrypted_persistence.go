package persistence

import (
	"crypto/sha256"

	"github.com/keep-network/keep-core/pkg/encryption"
)

// KeyLength represents the byte size of the key.
const KeyLength = encryption.KeyLength

type encryptedPersistence struct {
	delegate Handle
	box      encryption.Box
}

// NewEncryptedPersistence creates an adapter for the disk persistence to store data
// in an encrypted format
func NewEncryptedPersistence(handle Handle, password string) Handle {
	passwordBytes := []byte(password)

	return &encryptedPersistence{
		delegate: handle,
		box:      encryption.NewBox(sha256.Sum256(passwordBytes)),
	}
}

func (ep *encryptedPersistence) Save(data []byte, directory string, name string) error {
	encrypted, err := ep.box.Encrypt(data)
	if err != nil {
		return err
	}

	return ep.delegate.Save(encrypted, directory, name)
}

func (ep *encryptedPersistence) ReadAll() (<-chan DataDescriptor, <-chan error) {
	outputData := make(chan DataDescriptor)
	outputErrors := make(chan error)

	inputData, inputErrors := ep.delegate.ReadAll()

	// pass thru all errors from the input to the output channel without
	// changing anything
	go func() {
		defer close(outputErrors)
		for err := range inputErrors {
			outputErrors <- err
		}
	}()

	// pipe input data descriptor channel to the output data descriptor channel
	// decorading the descriptor passed so that the content is decrypted on read
	go func() {
		defer close(outputData)
		for descriptor := range inputData {
			// capture shared loop variable's value for the closure
			d := descriptor

			outputData <- &dataDescriptor{
				name:      d.Name(),
				directory: d.Directory(),
				readFunc: func() ([]byte, error) {
					content, err := d.Content()
					if err != nil {
						return nil, err
					}
					return ep.box.Decrypt(content)
				},
			}
		}
	}()

	return outputData, outputErrors
}

func (ep *encryptedPersistence) Archive(directory string) error {
	return ep.delegate.Archive(directory)
}
