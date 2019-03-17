package relay

import (
	"math/big"
	"time"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/altbn128"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/bls"
)

// The data below should match genesis relay request data defined on contract
// initialization i.e. in 2_deploy_contracts.js. Successfull genesis entry will
// trigger creation of the first group that will be chosen to respond on the
// next relay request, resulting another relay entry with creation of another
// group and so on.

// https://www.wolframalpha.com/input/?i=pi+to+78+digits
const piAsString = "31415926535897932384626433832795028841971693993751058209749445923078164062862"

// https://www.wolframalpha.com/input/?i=e+to+78+digits
const eAsString = "27182818284590452353602874713526624977572470936999595749669676277240766303535"

const secretKey = "123"

// GenesisRelayEntry generates genesis relay entry.
func GenesisRelayEntry() *event.Entry {
	// genesisEntryValue is the initial entry value for the network.
	// The n digits of pi that fit into a big.Int represent a "nothing up our
	// sleeve" value that all consumers of this network can verify.
	genesisEntryValue := bigFromBase10(piAsString)

	// genesisSeed value is the seed value used for the initial entry for the
	// network. The n digits of e that fit into a big.Int represent a "nothing
	// up our sleeve" value that all consumers of this network can verify.
	genesisSeedValue := bigFromBase10(eAsString)

	combinedEntryToSign := make([]byte, 0)
	combinedEntryToSign = append(combinedEntryToSign, genesisEntryValue.Bytes()...)
	combinedEntryToSign = append(combinedEntryToSign, genesisSeedValue.Bytes()...)

	// BLS signature for provided genesisEntryValue signed with secretKey
	genesisGroupSignature := altbn128.G1Point{
		G1: bls.Sign(bigFromBase10(secretKey), combinedEntryToSign),
	}.Compress()

	// public key for provided secretKey
	genesisGroupPubKey := altbn128.G2Point{
		G2: new(bn256.G2).ScalarBaseMult(bigFromBase10(secretKey)),
	}.Compress()

	return &event.Entry{
		RequestID:     big.NewInt(int64(1)),
		Value:         new(big.Int).SetBytes(genesisGroupSignature),
		GroupPubKey:   genesisGroupPubKey,
		PreviousEntry: genesisEntryValue,
		Timestamp:     time.Now().UTC(),
		Seed:          genesisSeedValue,
	}
}

// bigFromBase10 returns a big number from its string representation.
func bigFromBase10(s string) *big.Int {
	value, success := new(big.Int).SetString(s, 10)
	if !success {
		panic("failed to parse big number from string")
	}
	return value
}
