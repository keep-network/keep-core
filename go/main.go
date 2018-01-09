package main

import (
	"fmt"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/dfinity/go-dfinity-crypto/groupsig"
)

func main() {
	fmt.Printf("Starting!")
	groupsig.Init(bls.CurveFp382_1)

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
