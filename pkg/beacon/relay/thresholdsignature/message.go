package thresholdsignature

import "github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"

// SignatureShareMessage is a message payload that carries the sender's
// signature share for the given message.
type SignatureShareMessage struct {
	senderID   gjkr.MemberID
	ShareBytes []byte
}
