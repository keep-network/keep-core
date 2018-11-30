package gjkr

import "github.com/keep-network/keep-core/pkg/net/ephemeral"

// For complaint resolution, group members need to have access to messages
// exchanged between the accuser and the accused party. There are two situations
// in the DKG protocol where group members generate values for every other group
// member:
//
// - Ephemeral ECDH (phase 2) - after each group member generates an ephemeral
// keypair for each other group member and broadcasts those ephemeral public keys
// in the clear (phase 1), group members must ECDH those public keys with the
// ephemeral private key for that group member to derive a symmetric key.
// In the case of an accusation, members performing compliant resolution need to
// validate the private ephemeral key revealed by the accuser. To perform the
// validation, members need to compare public ephemeral key published by the
// accuser in phase 1 with the private ephemeral key published by the accuser.
//
// - Polynomial generation (phase 3) - each group member generates two sharing
// polynomials, and calculates shares as points on these polynomials individually
// for each other group member. Shares are publicly broadcast, encrypted with a
// symmetric key established between the sender and receiver. In the case of an
// accusation, members performing compliant resolution need to look at the shares
// sent by the accused party. To do this, they read the round 3 message from the
// buffer, passing the symmetric key used between the accuser and accused so that
// the round 3 message from the accused party can be decrypted.
type evidenceLog interface {
	// ephemeralPublicKeyMessage returns the `EphemeralPublicKeyMessage`
	// broadcast in the first protocol round by the given sender for the
	// given receiver.
	ephemeralPublicKeyMessage(
		sender MemberID,
		receiver MemberID,
	) *EphemeralPublicKeyMessage

	// peerSharesMessage returns the `PeerShareMessage` broadcast in the third
	// protocol round by the given sender for the given receiver.
	peerSharesMessage(
		sender MemberID,
		receiver MemberID,
	) *PeerSharesMessage
}
