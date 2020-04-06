const KeepRandomBeaconOperator = artifacts.require('KeepRandomBeaconOperator.sol')

module.exports = async function () {
  try {
    const contract = await KeepRandomBeaconOperator.deployed()
  
    const dkgGas = await contract.dkgGasEstimate()
    const gasPrice = await contract.gasPriceCeiling()
    const dkgFee = dkgGas.mul(gasPrice)

    await contract.genesis({value: dkgFee})
    console.log('Genesis successfully triggered.')    
  } catch(error) {
    console.error('Could not trigger genesis', error)
  }

  process.exit()
}
