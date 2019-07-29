const KeepRandomBeaconOperator = artifacts.require('KeepRandomBeaconOperator.sol')

// The data below should match genesis relay request data defined on contract
// initialization i.e. in 2_deploy_contracts.js. Successful genesis entry will
// trigger creation of the first group that will be chosen to respond on the
// next relay request, resulting another relay entry with creation of another
// group and so on.

// https://www.wolframalpha.com/input/?i=pi+to+78+digits
const previousEntry = web3.utils.toBN('31415926535897932384626433832795028841971693993751058209749445923078164062862')

// https://www.wolframalpha.com/input/?i=e+to+78+digits
const seed = web3.utils.toBN('27182818284590452353602874713526624977572470936999595749669676277240766303535')

// Data generated using client keep-core/pkg/bls package signing previous entry using master secret key '123'
const groupSignature = web3.utils.toBN('10920102476789591414949377782104707130412218726336356788412941355500907533021')

module.exports = async function () {

  const contract = await KeepRandomBeaconOperator.deployed()
  try {
    await contract.relayEntry(groupSignature, previousEntry, seed)
    console.log('Genesis entry successfully submitted.')
  } catch(error) {
    console.error('Genesis entry submission failed with', error)
  }

  process.exit()
}
