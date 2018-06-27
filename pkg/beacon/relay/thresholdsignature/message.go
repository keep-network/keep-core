package thresholdsignature

import "github.com/dfinity/go-dfinity-crypto/bls"

// SignatureShareMessage is a message payload that carries the sender's
// signature share for the given message.
type SignatureShareMessage struct {
	ID         *bls.ID
	ShareBytes []byte
}
