package gjkr

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func initializeSymmetricKeyMembersGroup(
	threshold int,
	groupSize int,
	dkg *DKG,
) ([]*SymmetricKeyGeneratingMember, error) {
	group := &Group{
		groupSize:          groupSize,
		dishonestThreshold: threshold,
	}

	var members []*SymmetricKeyGeneratingMember

	for i := 1; i <= groupSize; i++ {
		id := MemberID(i)
		members = append(members, &SymmetricKeyGeneratingMember{
			memberCore: &memberCore{
				ID:             id,
				group:          group,
				protocolConfig: dkg,
			},
			ephemeralKeyPairs: make(map[MemberID]*ephemeral.KeyPair),
			symmetricKeys:     make(map[MemberID]ephemeral.SymmetricKey),
		})
		group.RegisterMemberID(id)
	}

	// generate ephemeral key pairs for all other members of the group
	for _, member1 := range members {
		for _, member2 := range members {
			if member1.ID != member2.ID {

				keyPair, err := ephemeral.GenerateKeyPair()
				if err != nil {
					return nil, fmt.Errorf(
						"SymmetricKeyGeneratingMember initialization failed [%v]",
						err,
					)
				}
				member1.ephemeralKeyPairs[member2.ID] = keyPair
			}
		}
	}

	// generate symmetric keys with all other members of the group
	for _, member1 := range members {
		for _, member2 := range members {
			if member1.ID != member2.ID {

				privKey := member1.ephemeralKeyPairs[member2.ID].PrivateKey
				pubKey := member2.ephemeralKeyPairs[member1.ID].PublicKey
				member1.symmetricKeys[member2.ID] = privKey.Ecdh(pubKey)
			}
		}
	}

	return members, nil
}
