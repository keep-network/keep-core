package thresholdsignature

import "github.com/keep-network/keep-core/pkg/beacon/relay/group"

// SignatureShareMessage is a message payload that carries the sender's
// signature share for the given message.
type SignatureShareMessage struct {
	senderID   group.MemberIndex
	ShareBytes []byte
}
