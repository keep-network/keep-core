package thresholdsignature

// SignatureShareMessage is a message payload that carries the sender's
// signature share for the given message.
type SignatureShareMessage struct {
	ID         string
	ShareBytes []byte
}
