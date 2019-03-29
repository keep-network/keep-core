package thresholdsignature

import "github.com/keep-network/keep-core/pkg/beacon/relay/member"

// SignatureShareMessage is a message payload that carries the sender's
// signature share for the given message.
type SignatureShareMessage struct {
	senderID   member.MemberIndex
	ShareBytes []byte
}
