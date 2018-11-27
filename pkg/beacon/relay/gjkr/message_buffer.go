package gjkr

import "github.com/keep-network/keep-core/pkg/net/ephemeral"

// For the sake of complaint resolution, group members need to have an access
// to messages exchanged between the accuser and the accused party. There are
// two situations in DKG protocol when group member generates values
// individually for each other group member:
//
// - Ephemeral key generation (phase 1) - each group member generates an
// ephemeral keypair for each other group member and broadcast all the ephemeral
// public keys publicly. In case of an accusation, members performing compliant
// resolution need to validate private ephemeral key revealed by the accuser.
// To perform the validation, members need to compare public ephemeral key
// published by the accuser in phase 1 with the private ephemeral key published
// by the accuser with the complaint.
//
// - Polynomial generation (phase 3) - each group member generates two sharing
// polynomials and calculates shares as points on these polynomials individually
// for each other group member. Shares are publicly broadcast, encrypted with
// a symmetric key established between the sender and receiver. In case of an
// accusation, members performing compliant resolution need to look at the
// shares sent by the accused party. To do that, they read the round 3 message
// from the buffer passing the symmetric key used between the accuser and
// accused party so that round 3 message from the accused party can be
// decrypted.
type messageBuffer interface {

	// ephemeralPublicKeyMessage returns the `EphemeralPublicKeyMessage`
	// broadcast in the first protocol round by the given sender for the
	// given receiver.
	ephemeralPublicKeyMessage(
		sender int,
		receiver int,
	) *EphemeralPublicKeyMessage

	// peerSharesMessage returns the `PeerShareMessage` broadcast in the third
	// protocol round by the given sender for the given receiver. It is required
	// to pass an `ephemeral.SymmetricKey` used to encrypt the communication
	// between the sender and receiver.
	peerSharesMessage(
		sender int,
		receiver int,
		key ephemeral.SymmetricKey,
	) *PeerSharesMessage
}
