const KeepRandomBeaconOperator = artifacts.require('KeepRandomBeaconOperator.sol')

module.exports = async function () {

  const contract = await KeepRandomBeaconOperator.deployed()
  const dkgGas = await contract.dkgGasEstimate()
  const gasPrice = web3.utils.toBN(20).mul(web3.utils.toBN(10**9)) // 20 Gwei
  const dkgFee = dkgGas.mul(gasPrice)

  try {
    await contract.genesis({value: dkgFee})
    console.log('Genesis successfully triggered.')
  } catch(error) {
    console.error('Could not trigger genesis', error)
  }

  process.exit()
}
