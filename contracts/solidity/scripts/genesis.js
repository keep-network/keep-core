const KeepRandomBeaconOperator = artifacts.require('KeepRandomBeaconOperator.sol')

module.exports = async function () {

  const contract = await KeepRandomBeaconOperator.deployed()
  const dkgGas = await contract.dkgGasEstimate()
  const fluctuationMargin = await contract.fluctuationMargin()
  const priceFeedEstimate = web3.utils.toBN(20).mul(web3.utils.toBN(10**9)) // 20 Gwei
  const gasPriceWithFluctuationMargin = priceFeedEstimate.add(priceFeedEstimate.mul(fluctuationMargin).div(web3.utils.toBN(100)));
  const dkgFee = dkgGas.mul(gasPriceWithFluctuationMargin)

  try {
    await contract.genesis({value: dkgFee})
    console.log('Genesis successfully triggered.')
  } catch(error) {
    console.error('Could not trigger genesis', error)
  }

  process.exit()
}
