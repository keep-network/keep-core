const KeepRandomBeaconOperator = artifacts.require('KeepRandomBeaconOperator.sol')

module.exports = async function () {

  const contract = await KeepRandomBeaconOperator.deployed()
  try {
    await contract.genesis()
    console.log('Genesis successfully triggered.')
  } catch(error) {
    console.error('Could not trigger genesis', error)
  }

  process.exit()
}
