import Contract from 'truffle-contract'

import KeepRandomBeaconService from "contracts/KeepRandomBeaconService-deployed.json"
import KeepRandomBeaconServiceImplV1 from "contracts/KeepRandomBeaconServiceImplV1-deployed.json"

const addContracts = async (drizzle) => {
  // Proxy the KeepRandomBeaconServiceImplV1 contract through the KeepRandomBeaconService
  const keepRandomBeaconServicContract = Contract(KeepRandomBeaconService)
  keepRandomBeaconServicContract.setProvider(window.web3.currentProvider)
  const keepRandomBeaconServiceInstance = await keepRandomBeaconServicContract.deployed()

  // Keeping in case I need it
  // const keepRandomBeaconServiceImplV1Contract = Contract(KeepRandomBeaconServiceImplV1)
  // keepRandomBeaconServiceImplV1Contract.setProvider(window.web3.currentProvider)
  // const keepRandomBeaconServiceImplV1Instance = await keepRandomBeaconServiceImplV1Contract.at(keepRandomBeaconServiceInstance.address)

  // Address from the KeepRandomBeaconService
  const address = keepRandomBeaconServiceInstance.address

  // All other data from KeepRandomBeaconServiceImplV1
  const contractName = KeepRandomBeaconServiceImplV1.contractName
  const abi = KeepRandomBeaconServiceImplV1.abi
  const bytecode = KeepRandomBeaconServiceImplV1.deployedBytecode

  const contract = {
    contractName,
    web3Contract: new drizzle.web3.eth.Contract(abi, address, { data: bytecode })
  }

  drizzle.addContract(contract)
}

export default {
  addContracts
}