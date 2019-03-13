package relay

import "math/big"
import bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
import "github.com/keep-network/keep-core/pkg/altbn128"
import "github.com/keep-network/keep-core/pkg/bls"

// The data below should match genesis relay request data defined on contract
// initialization i.e. in 2_deploy_contracts.js. Successfull genesis entry will
// trigger creation of the first group that will be chosen to respond on the next
// relay request, resulting another relay entry with creation of another group and so on.

// https://www.wolframalpha.com/input/?i=pi+to+78+digits
const piAsString = "31415926535897932384626433832795028841971693993751058209749445923078164062862"
const secretKey = "123"

// GenesisEntryValue is the seed value for the network. The n digits of pi that fit into a
// *big.Int represent a "nothing up our sleeve" value that all consumers of this
// network can verify.
func GenesisEntryValue() *big.Int {
	return bigFromBase10(piAsString)
}

// GenesisGroupPubKey is the public key for provided secretKey
func GenesisGroupPubKey() []byte {
	return altbn128.G2Point{new(bn256.G2).ScalarBaseMult(bigFromBase10(secretKey))}.Compress()
}

// GenesisGroupSignature is BLS signature for provided GenesisEntryValue signed with secretKey
func GenesisGroupSignature() []byte {
	return altbn128.G1Point{bls.Sign(bigFromBase10(secretKey), GenesisEntryValue().Bytes())}.Compress()
}

// bigFromBase10 returns a big number from it's string representation.
func bigFromBase10(s string) *big.Int {
	value, success := new(big.Int).SetString(s, 10)
	if !success {
		panic("failed to parse big number from string")
	}
	return value
}
