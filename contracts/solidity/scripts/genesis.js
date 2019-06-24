const KeepRandomBeaconProxy = artifacts.require('KeepRandomBeacon.sol')
const KeepRandomBeacon = artifacts.require("KeepRandomBeaconImplV1")

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
const groupPubKey = "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d0"
const groupSignature = web3.utils.toBN('10920102476789591414949377782104707130412218726336356788412941355500907533021')

module.exports = async function () {

  const keepRandomBeaconProxy = await KeepRandomBeaconProxy.deployed()
  let contract = await KeepRandomBeacon.at(keepRandomBeaconProxy.address)
  try {
    await contract.relayEntry(1, groupSignature, groupPubKey, previousEntry, seed)
    console.log('Genesis entry successfully submitted.')
  } catch(error) {
    console.error('Genesis entry submission failed with', error)
  }

  process.exit()
}
