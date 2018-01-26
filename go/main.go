package main

import (
	"fmt"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/dfinity/go-dfinity-crypto/groupsig"
)

type GroupMember struct {
	id                  bls.ID
	receivedShares      []groupsig.Seckey
	verificationVectors [][]*groupsig.Pubkey
}

func idsFromGroupMembers(members []GroupMember) []*bls.ID {
	memberIds := make([]*bls.ID, 0)
	for _, member := range members {
		memberIds = append(memberIds, &member.id)
	}
	return memberIds
}

// Each player generates a commitment, which consists of a set of secret
// shares, which are used in turn to generate a verification vector alongside a
// set of contributions, one per member in the system.
//
// Note that there are t verifications, where t is the threshold of honest
// players. However, there are n contributions, where n is the total number of
// players.
func generateCommitment(memberIDs []*bls.ID, threshold int) ([]*bls.PublicKey, []*bls.SecretKey) {
	verificationVector := make([]*bls.PublicKey, threshold)
	secretShares := make([]bls.SecretKey, len(memberIDs))
	contribution := make([]*bls.SecretKey, 0)

	for i := 0; i < threshold; i++ {
		secretShares[i].SetByCSPRNG()
		verificationVector[i] = secretShares[i].GetPublicKey()
	}

	for _, memberID := range memberIDs {
		contributionPiece := bls.SecretKey{}
		contributionPiece.Set(secretShares, memberID)
		contribution := append(contribution, &contributionPiece)
	}

	return verificationVector, contribution
}

func main() {
	fmt.Printf("Starting!")
	groupsig.Init(bls.CurveFp382_1)

	threshold := 4

	memberNumbers := []int{0, 1, 2, 3, 4, 5, 6}
	members := [7]*GroupMember{}
	for _, number := range memberNumbers {
		member := GroupMember{bls.ID{}, make([]groupsig.Seckey, len(members)), make([][]groupsig.Pubkey, len(members))}
		member.id.SetDecString(string(number))
	}

	for _, member := range members {
		contribution := generateContribution(memberIds, threshold)
	}

	message := []byte("Booyan booyanescu")
	seckeys := []groupsig.Seckey{
		*groupsig.NewSeckeyFromInt(1),
		*groupsig.NewSeckeyFromInt(2),
		*groupsig.NewSeckeyFromInt(3),
		*groupsig.NewSeckeyFromInt(4),
		*groupsig.NewSeckeyFromInt(5),
	}
	pubkeys := []groupsig.Pubkey{
		*groupsig.NewPubkeyFromSeckey(seckeys[0]),
		*groupsig.NewPubkeyFromSeckey(seckeys[1]),
		*groupsig.NewPubkeyFromSeckey(seckeys[2]),
		*groupsig.NewPubkeyFromSeckey(seckeys[3]),
		*groupsig.NewPubkeyFromSeckey(seckeys[4]),
	}
	signatures := []groupsig.Signature{
		groupsig.Sign(seckeys[0], message),
		groupsig.Sign(seckeys[1], message),
		groupsig.Sign(seckeys[2], message),
		groupsig.Sign(seckeys[3], message),
		groupsig.Sign(seckeys[4], message),
	}

	master := groupsig.AggregateSeckeys(seckeys)
	masterPub := groupsig.NewPubkeyFromSeckey(*master)
	signature := groupsig.Sign(*master, message)
	aggregatedSig := groupsig.AggregateSigs(signatures)

	verification := groupsig.VerifySig(*masterPub, message, signature) == groupsig.VerifySig(*masterPub, message, aggregatedSig)
	aggregateVerification := groupsig.VerifyAggregateSig(pubkeys, message, signature)
	batchVerification := groupsig.BatchVerify(pubkeys, message, signatures)

	fmt.Printf(
		"%v = %v\nVerified: %v\nVerified in aggregate: %v\nVerified by batch: %v\n",
		signature.GetHexString(),
		aggregatedSig.GetHexString(),
		verification,
		aggregateVerification,
		batchVerification,
	)
}
