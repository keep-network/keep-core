package relay

import (
	"math/big"

	"github.com/keep-network/keep-core/pkg/altbn128"
	"github.com/keep-network/keep-core/pkg/beacon/relay/entry"
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

const privateKey = "123"

// GenesisRelayEntry generates genesis relay entry.
func GenesisRelayEntry() *big.Int {
	// genesisEntryValue is the initial entry value for the network.
	// The n digits of pi that fit into a big.Int represent a "nothing up our
	// sleeve" value that all consumers of this network can verify.
	genesisEntryValue := bigFromBase10(piAsString)

	// genesisSeed value is the seed value used for the initial entry for the
	// network. The n digits of e that fit into a big.Int represent a "nothing
	// up our sleeve" value that all consumers of this network can verify.
	genesisSeedValue := bigFromBase10(eAsString)

	combinedEntryToSign := entry.CombineToSign(
		genesisEntryValue,
		genesisSeedValue,
	)

	// BLS signature for provided genesisEntryValue and seed signed with
	// the privateKey
	genesisGroupSignature := altbn128.G1Point{
		G1: bls.Sign(bigFromBase10(privateKey), combinedEntryToSign),
	}.Compress()

	return new(big.Int).SetBytes(genesisGroupSignature)
}

// bigFromBase10 returns a big number from its string representation.
func bigFromBase10(s string) *big.Int {
	value, success := new(big.Int).SetString(s, 10)
	if !success {
		panic("failed to parse big number from string")
	}
	return value
}
